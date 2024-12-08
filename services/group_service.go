package services

import (
	"context"
	"fmt"
	"time"

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
		CreatedAt:    time.Now(),
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

// ListGroupDetails retrieves the ID, name, and creation time of all groups, sorted alphabetically by name.
func (s *GroupService) ListGroups(ctx context.Context) ([]struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	CreatedAt time.Time          `bson:"created_at"`
}, error) {
	groupsColl := s.db.Collection("groups")

	// Query to fetch ID, name, and created_at fields, sorted by name
	findOptions := options.Find().
		SetSort(bson.D{{Key: "name", Value: 1}}). // Sort by name ascending
		SetProjection(bson.D{
			{Key: "_id", Value: 1},
			{Key: "name", Value: 1},
			{Key: "created_at", Value: 1},
		}) // Include only ID, name, and created_at

	cursor, err := groupsColl.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Extract the relevant fields from the result
	var groups []struct {
		ID        primitive.ObjectID `bson:"_id"`
		Name      string             `bson:"name"`
		CreatedAt time.Time          `bson:"created_at"`
	}
	if err := cursor.All(ctx, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

// CheckPassword verifies if the provided password matches the stored hash.
func CheckPassword(password, hash string) bool {
	return models.CheckPasswordHash(password, hash)
}
