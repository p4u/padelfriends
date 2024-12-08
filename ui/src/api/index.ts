import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
});

// Add response interceptor for better error handling
api.interceptors.response.use(
  response => response,
  error => {
    let errorMessage = 'An unexpected error occurred';
    
    if (error.response) {
      switch (error.response.status) {
        case 401:
          errorMessage = 'Unauthorized access';
          break;
        case 404:
          errorMessage = 'Resource not found';
          break;
        case 500:
          errorMessage = 'Server error';
          break;
      }
    } else if (error.request) {
      errorMessage = 'Network error - please check your connection';
    }
    
    return Promise.reject(errorMessage);
  }
);

export const groupApi = {
  create: (name: string, password: string) => 
    api.post('/group', JSON.stringify({ name, password })),
  
  list: () => api.get('/groups'),
  
  getByName: (name: string, password: string) =>
    api.get(`/group/byname/${name}`, { params: { password } }),
  
  addPlayer: (groupId: string, password: string, name: string) =>
    api.post(`/group/${groupId}/players`, JSON.stringify({ name }), { params: { password } }),
  
  getPlayers: (groupId: string, password: string) =>
    api.get(`/group/${groupId}/players`, { params: { password } }),
  
  createMatch: (groupId: string, password: string, playerIds: string[]) =>
    api.post(`/group/${groupId}/matches`, JSON.stringify({ player_ids: playerIds }), { params: { password } }),
  
  submitResults: (groupId: string, matchId: string, password: string, scoreTeam1: number, scoreTeam2: number) =>
    api.post(`/group/${groupId}/matches/${matchId}/results`, 
      JSON.stringify({ score_team1: scoreTeam1, score_team2: scoreTeam2 }), 
      { params: { password } }
    ),
  
  getMatches: (groupId: string, password: string) =>
    api.get(`/group/${groupId}/matches`, { params: { password } }),
  
  getStatistics: (groupId: string, password: string) =>
    api.get(`/group/${groupId}/statistics`, { params: { password } }),
};