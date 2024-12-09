package services

import (
	"context"
	"errors"
	"time"

	"github.com/p4u/padelfriends/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MatchService struct {
	db *mongo.Database
}

func NewMatchService(db *mongo.Database) *MatchService {
	return &MatchService{db: db}
}

// getPlayerInfo retrieves player information by ID
func (s *MatchService) getPlayerInfo(ctx context.Context, playerID primitive.ObjectID) (models.PlayerInfo, error) {
	playersColl := s.db.Collection("players")
	var player models.Player
	err := playersColl.FindOne(ctx, bson.M{"_id": playerID}).Decode(&player)
	if err != nil {
		return models.PlayerInfo{}, err
	}
	return models.PlayerInfo{
		ID:   player.ID,
		Name: player.Name,
	}, nil
}

// getPlayersInfo retrieves information for multiple players
func (s *MatchService) getPlayersInfo(ctx context.Context, playerIDs []primitive.ObjectID) ([]models.PlayerInfo, error) {
	var players []models.PlayerInfo
	for _, id := range playerIDs {
		player, err := s.getPlayerInfo(ctx, id)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	return players, nil
}

// hasDuplicatePlayers checks if there are any duplicate player IDs
func hasDuplicatePlayers(playerIDs []primitive.ObjectID) bool {
	seen := make(map[primitive.ObjectID]bool)
	for _, id := range playerIDs {
		if seen[id] {
			return true
		}
		seen[id] = true
	}
	return false
}

// CreateMatch starts a new match record.
func (s *MatchService) CreateMatch(ctx context.Context, groupName string, playerIDs []primitive.ObjectID) (models.MatchResponse, error) {
	if len(playerIDs) != 4 {
		return models.MatchResponse{}, errors.New("exactly 4 players required for a match")
	}

	if hasDuplicatePlayers(playerIDs) {
		return models.MatchResponse{}, errors.New("duplicate players are not allowed in a match")
	}

	matchesColl := s.db.Collection("matches")
	detailsColl := s.db.Collection("matchdetails")

	match := models.Match{
		GroupName: groupName,
		Timestamp: time.Now(),
		Status:    "pending",
	}
	res, err := matchesColl.InsertOne(ctx, match)
	if err != nil {
		return models.MatchResponse{}, err
	}
	match.ID = res.InsertedID.(primitive.ObjectID)

	detail := models.MatchDetail{
		MatchID:    match.ID,
		Team1:      playerIDs[:2],
		Team2:      playerIDs[2:],
		ScoreTeam1: 0,
		ScoreTeam2: 0,
	}

	_, err = detailsColl.InsertOne(ctx, detail)
	if err != nil {
		return models.MatchResponse{}, err
	}

	team1Players, err := s.getPlayersInfo(ctx, detail.Team1)
	if err != nil {
		return models.MatchResponse{}, err
	}

	team2Players, err := s.getPlayersInfo(ctx, detail.Team2)
	if err != nil {
		return models.MatchResponse{}, err
	}

	response := models.MatchResponse{
		ID:         match.ID,
		GroupName:  match.GroupName,
		Timestamp:  match.Timestamp,
		Team1:      team1Players,
		Team2:      team2Players,
		ScoreTeam1: detail.ScoreTeam1,
		ScoreTeam2: detail.ScoreTeam2,
		Status:     match.Status,
	}

	return response, nil
}

// CreateMatches creates multiple matches at once
func (s *MatchService) CreateMatches(ctx context.Context, groupName string, matchesPlayerIDs [][]primitive.ObjectID) ([]models.MatchResponse, error) {
	var responses []models.MatchResponse
	for _, playerIDs := range matchesPlayerIDs {
		match, err := s.CreateMatch(ctx, groupName, playerIDs)
		if err != nil {
			return nil, err
		}
		responses = append(responses, match)
	}
	return responses, nil
}

// CancelMatch deletes a match and its details
func (s *MatchService) CancelMatch(ctx context.Context, matchID primitive.ObjectID) error {
	matchesColl := s.db.Collection("matches")
	detailsColl := s.db.Collection("matchdetails")

	session, err := s.db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		_, err := detailsColl.DeleteOne(sessCtx, bson.M{"match_id": matchID})
		if err != nil {
			return nil, err
		}

		result, err := matchesColl.DeleteOne(sessCtx, bson.M{
			"_id":    matchID,
			"status": "pending",
		})
		if err != nil {
			return nil, err
		}

		if result.DeletedCount == 0 {
			return nil, errors.New("match not found or already completed")
		}

		return nil, nil
	})

	return err
}

// GetRecentMatches returns the last 20 matches for a group
func (s *MatchService) GetRecentMatches(ctx context.Context, groupName string) ([]models.MatchResponse, error) {
	matchesColl := s.db.Collection("matches")
	detailsColl := s.db.Collection("matchdetails")

	// Get last 20 matches
	findOptions := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(20)

	cur, err := matchesColl.Find(ctx, bson.M{"group_name": groupName}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var matches []models.Match
	if err := cur.All(ctx, &matches); err != nil {
		return nil, err
	}

	var responses []models.MatchResponse
	for _, match := range matches {
		var detail models.MatchDetail
		err := detailsColl.FindOne(ctx, bson.M{"match_id": match.ID}).Decode(&detail)
		if err != nil {
			continue
		}

		team1Players, err := s.getPlayersInfo(ctx, detail.Team1)
		if err != nil {
			continue
		}

		team2Players, err := s.getPlayersInfo(ctx, detail.Team2)
		if err != nil {
			continue
		}

		response := models.MatchResponse{
			ID:         match.ID,
			GroupName:  match.GroupName,
			Timestamp:  match.Timestamp,
			Team1:      team1Players,
			Team2:      team2Players,
			ScoreTeam1: detail.ScoreTeam1,
			ScoreTeam2: detail.ScoreTeam2,
			Status:     match.Status,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// ListMatches returns all matches for a group with pagination
func (s *MatchService) ListMatches(ctx context.Context, groupName string, page, pageSize int) ([]models.MatchResponse, int, error) {
	matchesColl := s.db.Collection("matches")
	detailsColl := s.db.Collection("matchdetails")

	// Get total count
	totalCount, err := matchesColl.CountDocuments(ctx, bson.M{"group_name": groupName})
	if err != nil {
		return nil, 0, err
	}

	// Calculate skip value
	skip := (page - 1) * pageSize

	// Get matches with pagination
	findOptions := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	cur, err := matchesColl.Find(ctx, bson.M{"group_name": groupName}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var matches []models.Match
	if err := cur.All(ctx, &matches); err != nil {
		return nil, 0, err
	}

	var responses []models.MatchResponse
	for _, match := range matches {
		var detail models.MatchDetail
		err := detailsColl.FindOne(ctx, bson.M{"match_id": match.ID}).Decode(&detail)
		if err != nil {
			continue
		}

		team1Players, err := s.getPlayersInfo(ctx, detail.Team1)
		if err != nil {
			continue
		}

		team2Players, err := s.getPlayersInfo(ctx, detail.Team2)
		if err != nil {
			continue
		}

		response := models.MatchResponse{
			ID:         match.ID,
			GroupName:  match.GroupName,
			Timestamp:  match.Timestamp,
			Team1:      team1Players,
			Team2:      team2Players,
			ScoreTeam1: detail.ScoreTeam1,
			ScoreTeam2: detail.ScoreTeam2,
			Status:     match.Status,
		}
		responses = append(responses, response)
	}

	return responses, int(totalCount), nil
}

// SubmitResults updates the match detail with final scores.
func (s *MatchService) SubmitResults(ctx context.Context, matchID primitive.ObjectID, scoreTeam1, scoreTeam2 int) error {
	if scoreTeam1 < 0 || scoreTeam1 > 10 || scoreTeam2 < 0 || scoreTeam2 > 10 {
		return errors.New("invalid scores")
	}

	matchesColl := s.db.Collection("matches")
	detailsColl := s.db.Collection("matchdetails")

	session, err := s.db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		_, err := matchesColl.UpdateOne(
			sessCtx,
			bson.M{"_id": matchID, "status": "pending"},
			bson.M{"$set": bson.M{"status": "completed"}},
		)
		if err != nil {
			return nil, err
		}

		_, err = detailsColl.UpdateOne(
			sessCtx,
			bson.M{"match_id": matchID},
			bson.M{"$set": bson.M{
				"score_team1": scoreTeam1,
				"score_team2": scoreTeam2,
			}},
		)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}
