package services

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StatsService struct {
	db *mongo.Database
}

func NewStatsService(db *mongo.Database) *StatsService {
	return &StatsService{db: db}
}

type PlayerStats struct {
	PlayerID   primitive.ObjectID `json:"player_id"`
	PlayerName string             `json:"player_name"`

	// Game Statistics
	TotalGames   int     `json:"total_games"`
	GamesWon     int     `json:"games_won"`
	GamesLost    int     `json:"games_lost"`
	GameWinRate  float64 `json:"game_win_rate"`
	GameLossRate float64 `json:"game_loss_rate"`

	// Point Statistics
	TotalPoints   int     `json:"total_points"`
	PointsWon     int     `json:"points_won"`
	PointsLost    int     `json:"points_lost"`
	PointWinRate  float64 `json:"point_win_rate"`
	PointLossRate float64 `json:"point_loss_rate"`
}

// ComputeStats calculates statistics for all players in a group
func (s *StatsService) ComputeStats(ctx context.Context, groupName string) ([]PlayerStats, error) {
	matchesColl := s.db.Collection("matches")
	detailsColl := s.db.Collection("matchdetails")
	playersColl := s.db.Collection("players")

	// Get all completed matches for the group
	matches, err := matchesColl.Find(ctx, bson.M{
		"group_name": groupName,
		"status":     "completed",
	})
	if err != nil {
		return nil, err
	}
	defer matches.Close(ctx)

	// Initialize player stats map
	playerStatsMap := make(map[primitive.ObjectID]*PlayerStats)

	// Process each match
	for matches.Next(ctx) {
		var match struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		if err := matches.Decode(&match); err != nil {
			continue
		}

		// Get match details
		var detail struct {
			Team1      []primitive.ObjectID `bson:"team1"`
			Team2      []primitive.ObjectID `bson:"team2"`
			ScoreTeam1 int                  `bson:"score_team1"`
			ScoreTeam2 int                  `bson:"score_team2"`
		}
		if err := detailsColl.FindOne(ctx, bson.M{"match_id": match.ID}).Decode(&detail); err != nil {
			continue
		}

		// Update stats for Team 1 players
		for _, playerID := range detail.Team1 {
			if _, exists := playerStatsMap[playerID]; !exists {
				playerStatsMap[playerID] = &PlayerStats{PlayerID: playerID}
			}
			stats := playerStatsMap[playerID]
			stats.TotalGames++
			stats.TotalPoints += detail.ScoreTeam1
			stats.PointsWon += detail.ScoreTeam1
			stats.PointsLost += detail.ScoreTeam2
			if detail.ScoreTeam1 > detail.ScoreTeam2 {
				stats.GamesWon++
			} else {
				stats.GamesLost++
			}
		}

		// Update stats for Team 2 players
		for _, playerID := range detail.Team2 {
			if _, exists := playerStatsMap[playerID]; !exists {
				playerStatsMap[playerID] = &PlayerStats{PlayerID: playerID}
			}
			stats := playerStatsMap[playerID]
			stats.TotalGames++
			stats.TotalPoints += detail.ScoreTeam2
			stats.PointsWon += detail.ScoreTeam2
			stats.PointsLost += detail.ScoreTeam1
			if detail.ScoreTeam2 > detail.ScoreTeam1 {
				stats.GamesWon++
			} else {
				stats.GamesLost++
			}
		}
	}

	// Get player names and calculate rates
	var result []PlayerStats
	for playerID, stats := range playerStatsMap {
		// Get player name
		var player struct {
			Name string `bson:"name"`
		}
		if err := playersColl.FindOne(ctx, bson.M{"_id": playerID}).Decode(&player); err != nil {
			continue
		}

		// Calculate rates
		if stats.TotalGames > 0 {
			stats.GameWinRate = float64(stats.GamesWon) / float64(stats.TotalGames) * 100
			stats.GameLossRate = float64(stats.GamesLost) / float64(stats.TotalGames) * 100
		}
		if stats.TotalPoints > 0 {
			stats.PointWinRate = float64(stats.PointsWon) / float64(stats.TotalPoints) * 100
			stats.PointLossRate = float64(stats.PointsLost) / float64(stats.TotalPoints) * 100
		}

		stats.PlayerName = player.Name
		result = append(result, *stats)
	}

	return result, nil
}
