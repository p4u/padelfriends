import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json'
  }
});

export const groupApi = {
  getGroups: () => api.get('/groups'),
  
  getGroup: (name: string, password: string) => 
    api.get(`/groups/${name}`, { headers: { 'X-Group-Password': password } }),
  
  createGroup: (name: string, password: string) => 
    api.post('/groups', { name, password }),
  
  addPlayer: (groupName: string, password: string, playerName: string) => 
    api.post(`/groups/${groupName}/players`, 
      { name: playerName },
      { headers: { 'X-Group-Password': password } }
    ),
  
  getPlayers: (groupName: string, password: string) => 
    api.get(`/groups/${groupName}/players`,
      { headers: { 'X-Group-Password': password } }
    ),
  
  getMatches: (groupName: string, password: string, page: number = 1, pageSize: number = 10) => 
    api.get(`/groups/${groupName}/matches`,
      { 
        headers: { 'X-Group-Password': password },
        params: { page, page_size: pageSize }
      }
    ),
  
  createMatch: (groupName: string, password: string, playerIds: string[]) => 
    api.post(`/groups/${groupName}/matches`,
      { player_ids: playerIds },
      { headers: { 'X-Group-Password': password } }
    ),
  
  createBatchMatches: (groupName: string, password: string, matches: string[][]) => 
    api.post(`/groups/${groupName}/matches/batch`,
      { matches },
      { headers: { 'X-Group-Password': password } }
    ),
  
  submitResults: (groupName: string, matchId: string, password: string, scoreTeam1: number, scoreTeam2: number) => 
    api.post(`/groups/${groupName}/matches/${matchId}/results`,
      { score_team1: scoreTeam1, score_team2: scoreTeam2 },
      { headers: { 'X-Group-Password': password } }
    ),
  
  cancelMatch: (groupName: string, matchId: string, password: string) => 
    api.delete(`/groups/${groupName}/matches/${matchId}`,
      { headers: { 'X-Group-Password': password } }
    ),
  
  getStatistics: (groupName: string, password: string) => 
    api.get(`/groups/${groupName}/statistics`,
      { headers: { 'X-Group-Password': password } }
    )
};

export const healthApi = {
  check: () => api.get('/health')
};
