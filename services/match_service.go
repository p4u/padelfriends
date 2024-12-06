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

// CreateMatch starts a new match record.
func (s *MatchService) CreateMatch(ctx context.Context, groupID primitive.ObjectID, playerIDs []primitive.ObjectID) (models.Match, models.MatchDetail, error) {
	if len(playerIDs) != 4 {
		return models.Match{}, models.MatchDetail{}, errors.New("exactly 4 players required for a match")
	}

	matchesColl := s.db.Collection("matches")
	detailsColl := s.db.Collection("matchdetails")

	// Here we’d randomize teams. Simplify for now: first two = team1, last two = team2
	match := models.Match{
		GroupID:   groupID,
		Timestamp: time.Now(),
	}
	res, err := matchesColl.InsertOne(ctx, match)
	if err != nil {
		return models.Match{}, models.MatchDetail{}, err
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
		return models.Match{}, models.MatchDetail{}, err
	}

	return match, detail, nil
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

// ListMatches returns all matches for a group.
func (s *MatchService) ListMatches(ctx context.Context, groupID primitive.ObjectID) ([]models.Match, error) {
	matchesColl := s.db.Collection("matches")

	cur, err := matchesColl.Find(ctx, bson.M{"group_id": groupID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var matches []models.Match
	if err := cur.All(ctx, &matches); err != nil {
		return nil, err
	}
	return matches, nil
}
