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
  status: 'pending' | 'completed';
}

export interface Statistics {
  player_id: string;
  player_name: string;
  matches_played: number;
  matches_won: number;
  win_rate: number;
}
