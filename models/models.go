package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Group represents a padel group context.
type Group struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	PasswordHash string             `bson:"password_hash"`
}

// Player represents a player within a group.
type Player struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	GroupID primitive.ObjectID `bson:"group_id"`
	Name    string             `bson:"name"`
}

// Match represents a single match played in a group.
type Match struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	GroupID   primitive.ObjectID `bson:"group_id"`
	Timestamp time.Time          `bson:"timestamp"`
}

// MatchDetail stores the details of a match (teams, scores).
type MatchDetail struct {
	MatchID    primitive.ObjectID   `bson:"match_id"`
	Team1      []primitive.ObjectID `bson:"team1"`
	Team2      []primitive.ObjectID `bson:"team2"`
	ScoreTeam1 int                  `bson:"score_team1"`
	ScoreTeam2 int                  `bson:"score_team2"`
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
