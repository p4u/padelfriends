package services

import (
	"context"

	"github.com/p4u/padelfriends/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GroupService struct {
	db *mongo.Database
}

func NewGroupService(db *mongo.Database) *GroupService {
	return &GroupService{db: db}
}

// CreateGroup creates a new group with a hashed password.
func (s *GroupService) CreateGroup(ctx context.Context, name, password string) (models.Group, error) {
	groupsColl := s.db.Collection("groups")

	// Hash the password
	hash, err := models.HashPassword(password)
	if err != nil {
		return models.Group{}, err
	}

	group := models.Group{
		Name:         name,
		PasswordHash: hash,
	}

	res, err := groupsColl.InsertOne(ctx, group)
	if err != nil {
		return models.Group{}, err
	}

	group.ID = res.InsertedID.(primitive.ObjectID)
	return group, nil
}

// GetGroupByName retrieves a group by name.
func (s *GroupService) GetGroupByName(ctx context.Context, name string) (models.Group, error) {
	groupsColl := s.db.Collection("groups")
	var g models.Group
	err := groupsColl.FindOne(ctx, bson.M{"name": name}).Decode(&g)
	if err != nil {
		return models.Group{}, err
	}
	return g, nil
}

// CheckGroupPassword verifies a group's password.
func (s *GroupService) CheckGroupPassword(ctx context.Context, name, password string) (bool, error) {
	g, err := s.GetGroupByName(ctx, name)
	if err != nil {
		return false, err
	}
	return models.CheckPasswordHash(password, g.PasswordHash), nil
}

// Add this method after GetGroupByName:
func (s *GroupService) GetGroupByID(ctx context.Context, id primitive.ObjectID) (models.Group, error) {
	groupsColl := s.db.Collection("groups")
	var g models.Group
	err := groupsColl.FindOne(ctx, bson.M{"_id": id}).Decode(&g)
	if err != nil {
		return models.Group{}, err
	}
	return g, nil
}

// Simple wrapper around models.CheckPasswordHash
func CheckPassword(password, hash string) bool {
	return models.CheckPasswordHash(password, hash)
}
