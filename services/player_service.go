package services

import (
	"context"
	"errors"

	"github.com/p4u/padelfriends/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PlayerService struct {
	db *mongo.Database
}

func NewPlayerService(db *mongo.Database) *PlayerService {
	return &PlayerService{db: db}
}

// AddPlayer adds a player to a group if not duplicate.
func (s *PlayerService) AddPlayer(ctx context.Context, groupName string, name string) (models.Player, error) {
	playersColl := s.db.Collection("players")

	// Check duplicate
	err := playersColl.FindOne(ctx, bson.M{"group_name": groupName, "name": name}).Err()
	if err == nil {
		return models.Player{}, errors.New("player already exists in this group")
	}
	if err != mongo.ErrNoDocuments {
		return models.Player{}, err
	}

	p := models.Player{
		GroupName: groupName,
		Name:      name,
	}

	res, err := playersColl.InsertOne(ctx, p)
	if err != nil {
		return models.Player{}, err
	}
	p.ID = res.InsertedID.(primitive.ObjectID)
	return p, nil
}

// ListPlayers lists all players for a given group.
func (s *PlayerService) ListPlayers(ctx context.Context, groupName string) ([]models.Player, error) {
	playersColl := s.db.Collection("players")

	cur, err := playersColl.Find(ctx, bson.M{"group_name": groupName}, options.Find().SetSort(bson.M{"name": 1}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var players []models.Player
	if err := cur.All(ctx, &players); err != nil {
		return nil, err
	}
	return players, nil
}
