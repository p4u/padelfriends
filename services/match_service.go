package services

import (
	"context"
	"errors"
	"time"

	"github.com/p4u/padelfriends/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	// Check for duplicate players
	if hasDuplicatePlayers(playerIDs) {
		return models.MatchResponse{}, errors.New("duplicate players are not allowed in a match")
	}

	matchesColl := s.db.Collection("matches")
	detailsColl := s.db.Collection("matchdetails")

	// Create the match
	match := models.Match{
		GroupName: groupName,
		Timestamp: time.Now(),
	}
	res, err := matchesColl.InsertOne(ctx, match)
	if err != nil {
		return models.MatchResponse{}, err
	}
	match.ID = res.InsertedID.(primitive.ObjectID)

	// Create match details
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

	// Get player information
	team1Players, err := s.getPlayersInfo(ctx, detail.Team1)
	if err != nil {
		return models.MatchResponse{}, err
	}

	team2Players, err := s.getPlayersInfo(ctx, detail.Team2)
	if err != nil {
		return models.MatchResponse{}, err
	}

	// Create response
	response := models.MatchResponse{
		ID:         match.ID,
		GroupName:  match.GroupName,
		Timestamp:  match.Timestamp,
		Team1:      team1Players,
		Team2:      team2Players,
		ScoreTeam1: detail.ScoreTeam1,
		ScoreTeam2: detail.ScoreTeam2,
		Status:     "pending",
	}

	return response, nil
}

// SubmitResults updates the match detail with final scores.
func (s *MatchService) SubmitResults(ctx context.Context, matchID primitive.ObjectID, scoreTeam1, scoreTeam2 int) error {
	if scoreTeam1 < 0 || scoreTeam1 > 8 || scoreTeam2 < 0 || scoreTeam2 > 8 {
		return errors.New("invalid scores")
	}
	detailsColl := s.db.Collection("matchdetails")

	_, err := detailsColl.UpdateOne(ctx, bson.M{"match_id": matchID},
		bson.M{"$set": bson.M{"score_team1": scoreTeam1, "score_team2": scoreTeam2}})
	return err
}

// ListMatches returns all matches for a group with player details.
func (s *MatchService) ListMatches(ctx context.Context, groupName string) ([]models.MatchResponse, error) {
	matchesColl := s.db.Collection("matches")
	detailsColl := s.db.Collection("matchdetails")

	// Get all matches
	cur, err := matchesColl.Find(ctx, bson.M{"group_name": groupName})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var matches []models.Match
	if err := cur.All(ctx, &matches); err != nil {
		return nil, err
	}

	// Build response with details and player information
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

		status := "pending"
		if detail.ScoreTeam1 > 0 || detail.ScoreTeam2 > 0 {
			status = "completed"
		}

		response := models.MatchResponse{
			ID:         match.ID,
			GroupName:  match.GroupName,
			Timestamp:  match.Timestamp,
			Team1:      team1Players,
			Team2:      team2Players,
			ScoreTeam1: detail.ScoreTeam1,
			ScoreTeam2: detail.ScoreTeam2,
			Status:     status,
		}
		responses = append(responses, response)
	}

	return responses, nil
}
