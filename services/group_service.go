package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/p4u/padelfriends/models"
	"go.mongodb.org/mongo-driver/bson"
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

	_, err = groupsColl.InsertOne(ctx, group)
	if err != nil {
		return models.Group{}, err
	}

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

type GroupDetails struct {
	Name      string    `bson:"name" json:"name"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

// ListGroupDetails retrieves the name and creation time of all groups, sorted alphabetically by name.
func (s *GroupService) ListGroups(ctx context.Context) ([]*GroupDetails, error) {
	groupsColl := s.db.Collection("groups")

	findOptions := options.Find().
		SetSort(bson.D{{Key: "name", Value: 1}}).
		SetProjection(bson.D{
			{Key: "name", Value: 1},
			{Key: "created_at", Value: 1},
		})

	cursor, err := groupsColl.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var groups []*GroupDetails
	if err := cursor.All(ctx, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

// ExportGroupMatchesCSV exports all matches for a group in CSV format
func (s *GroupService) ExportGroupMatchesCSV(ctx context.Context, groupName string) (string, error) {
	matchesColl := s.db.Collection("matches")
	matchDetailsColl := s.db.Collection("matchdetails")
	playersColl := s.db.Collection("players")

	// First get all completed matches
	cursor, err := matchesColl.Find(ctx, bson.M{
		"group_name": groupName,
		"status":     "completed",
	}, options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}))
	if err != nil {
		return "", fmt.Errorf("error finding matches: %v", err)
	}
	defer cursor.Close(ctx)

	var matches []models.Match
	if err := cursor.All(ctx, &matches); err != nil {
		return "", fmt.Errorf("error decoding matches: %v", err)
	}

	log.Printf("Found %d completed matches", len(matches))

	// CSV header
	csvLines := []string{"Date,Team 1 Player 1,Team 1 Player 2,Team 2 Player 1,Team 2 Player 2,Score Team 1,Score Team 2"}

	for _, match := range matches {
		// Get match details
		var detail models.MatchDetail
		err := matchDetailsColl.FindOne(ctx, bson.M{
			"match_id": match.ID,
		}).Decode(&detail)
		if err != nil {
			log.Printf("Error finding match details for match %s: %v", match.ID.Hex(), err)
			continue
		}

		log.Printf("Processing match %s from %s", match.ID.Hex(), match.Timestamp)

		// Get player names for team 1
		var team1Names []string
		for _, playerID := range detail.Team1 {
			var player models.Player
			err = playersColl.FindOne(ctx, bson.M{"_id": playerID}).Decode(&player)
			if err != nil {
				log.Printf("Error finding player: %v", err)
				team1Names = append(team1Names, "Unknown")
			} else {
				team1Names = append(team1Names, player.Name)
			}
		}

		// Get player names for team 2
		var team2Names []string
		for _, playerID := range detail.Team2 {
			var player models.Player
			err = playersColl.FindOne(ctx, bson.M{"_id": playerID}).Decode(&player)
			if err != nil {
				log.Printf("Error finding player: %v", err)
				team2Names = append(team2Names, "Unknown")
			} else {
				team2Names = append(team2Names, player.Name)
			}
		}

		// Ensure we have 2 players per team
		for len(team1Names) < 2 {
			team1Names = append(team1Names, "Unknown")
		}
		for len(team2Names) < 2 {
			team2Names = append(team2Names, "Unknown")
		}

		line := fmt.Sprintf("%s,%s,%s,%s,%s,%d,%d",
			match.Timestamp.Format("2006-01-02 15:04:05"),
			team1Names[0],
			team1Names[1],
			team2Names[0],
			team2Names[1],
			detail.ScoreTeam1,
			detail.ScoreTeam2,
		)
		csvLines = append(csvLines, line)

		log.Printf("Added line: %s", line)
	}

	return strings.Join(csvLines, "\n"), nil
}

// CheckPassword verifies if the provided password matches the stored hash.
func CheckPassword(password, hash string) bool {
	return models.CheckPasswordHash(password, hash)
}
