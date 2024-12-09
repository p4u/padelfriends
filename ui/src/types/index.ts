export interface Group {
  id: string;
  name: string;
  created_at: string;
}

export interface Player {
  id: string;
  name: string;
  group_id: string;
  created_at: string;
}

export interface PlayerInfo {
  id: string;
  name: string;
}

export interface Match {
  id: string;
  group_name: string;
  timestamp: string;
  team1: PlayerInfo[];
  team2: PlayerInfo[];
  score_team1?: number;
  score_team2?: number;
  status: 'pending' | 'completed' | 'cancelled';
}

export interface Statistics {
  player_id: string;
  player_name: string;
  
  // Game Statistics
  total_games: number;
  games_won: number;
  games_lost: number;
  game_win_rate: number;
  game_loss_rate: number;
  
  // Point Statistics
  total_points: number;
  points_won: number;
  points_lost: number;
  point_win_rate: number;
  point_loss_rate: number;
}

export interface CreateMatchPayload {
  player_ids: string[];
}

export interface CreateBatchMatchesPayload {
  matches: string[][];
}

export interface SubmitScorePayload {
  score_team1: number;
  score_team2: number;
}

// Index signature for dynamic key access in statistics sorting
export interface IndexSignature {
  [key: string]: any;
}

// Extend Statistics to include index signature
export interface StatisticsWithIndex extends Statistics, IndexSignature {}
