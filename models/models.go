package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Group represents a padel group context.
type Group struct {
	Name         string    `bson:"name" json:"name"`
	PasswordHash string    `bson:"password_hash" json:"-"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
}

// Player represents a player within a group.
type Player struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GroupName string             `bson:"group_name" json:"group_name"`
	Name      string             `bson:"name" json:"name"`
}

// Match represents a single match played in a group.
type Match struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GroupName string             `bson:"group_name" json:"group_name"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}

// MatchDetail stores the details of a match (teams, scores).
type MatchDetail struct {
	MatchID    primitive.ObjectID   `bson:"match_id" json:"match_id"`
	Team1      []primitive.ObjectID `bson:"team1" json:"team1"`
	Team2      []primitive.ObjectID `bson:"team2" json:"team2"`
	ScoreTeam1 int                  `bson:"score_team1" json:"score_team1"`
	ScoreTeam2 int                  `bson:"score_team2" json:"score_team2"`
}

// MatchResponse combines Match and MatchDetail with player names for API responses
type MatchResponse struct {
	ID         primitive.ObjectID `json:"id"`
	GroupName  string             `json:"group_name"`
	Timestamp  time.Time          `json:"timestamp"`
	Team1      []PlayerInfo       `json:"team1"`
	Team2      []PlayerInfo       `json:"team2"`
	ScoreTeam1 int                `json:"score_team1"`
	ScoreTeam2 int                `json:"score_team2"`
	Status     string             `json:"status"`
}

// PlayerInfo contains the essential player information for responses
type PlayerInfo struct {
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
}

// HashPassword hashes the given password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash verifies if the given password matches the hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
