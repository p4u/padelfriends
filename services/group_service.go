package services

import (
	"context"
	"fmt"

	"github.com/p4u/padelfriends/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GroupService provides methods to interact with groups in the database.
type GroupService struct {
	db *mongo.Database
}

// NewGroupService creates a new GroupService.
func NewGroupService(db *mongo.Database) *GroupService {
	return &GroupService{db: db}
}

// CreateGroup creates a new group with a hashed password.
func (s *GroupService) CreateGroup(ctx context.Context, name, password string) (models.Group, error) {
	groupsColl := s.db.Collection("groups")

	// Check if a group with the same name already exists
	count, err := groupsColl.CountDocuments(ctx, bson.M{"name": name})
	if err != nil {
		return models.Group{}, err
	}
	if count > 0 {
		return models.Group{}, fmt.Errorf("duplicate key: group name '%s' already exists", name)
	}

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

// GetGroupByName retrieves a group by its name.
func (s *GroupService) GetGroupByName(ctx context.Context, name string) (models.Group, error) {
	groupsColl := s.db.Collection("groups")
	var g models.Group
	err := groupsColl.FindOne(ctx, bson.M{"name": name}).Decode(&g)
	if err != nil {
		return models.Group{}, err
	}
	return g, nil
}

// GetGroupByID retrieves a group by its ID.
func (s *GroupService) GetGroupByID(ctx context.Context, id primitive.ObjectID) (models.Group, error) {
	groupsColl := s.db.Collection("groups")
	var g models.Group
	err := groupsColl.FindOne(ctx, bson.M{"_id": id}).Decode(&g)
	if err != nil {
		return models.Group{}, err
	}
	return g, nil
}

// ListGroups retrieves all groups.
// Note: Implement proper authorization to restrict access as needed.
func (s *GroupService) ListGroups(ctx context.Context) ([]models.Group, error) {
	groupsColl := s.db.Collection("groups")

	// Optionally, implement pagination or filtering here
	findOptions := options.Find().SetSort(bson.D{{Key: "name", Value: 1}}) // Sort by name ascending

	cursor, err := groupsColl.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var groups []models.Group
	if err := cursor.All(ctx, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

// CheckPassword verifies if the provided password matches the stored hash.
func CheckPassword(password, hash string) bool {
	return models.CheckPasswordHash(password, hash)
}
