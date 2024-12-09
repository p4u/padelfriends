import axios from 'axios';
import type { CreateMatchPayload, CreateBatchMatchesPayload, SubmitScorePayload } from '../types';

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
    api.post('/group', { name, password }),
  
  list: () => api.get('/groups'),
  
  getByName: (name: string, password: string) =>
    api.get(`/group/byname/${name}`, { params: { password } }),
  
  addPlayer: (groupId: string, password: string, name: string) =>
    api.post(`/group/${groupId}/players`, { name }, { params: { password } }),
  
  getPlayers: (groupId: string, password: string) =>
    api.get(`/group/${groupId}/players`, { params: { password } }),
  
  createMatch: (groupId: string, password: string, playerIds: string[]) =>
    api.post(`/group/${groupId}/matches`, 
      { player_ids: playerIds }, 
      { params: { password } }
    ),

  createBatchMatches: (groupId: string, password: string, matches: string[][]) =>
    api.post(`/group/${groupId}/matches/batch`,
      { matches },
      { params: { password } }
    ),
  
  cancelMatch: (groupId: string, matchId: string, password: string) =>
    api.post(`/group/${groupId}/matches/${matchId}/cancel`,
      {},
      { params: { password } }
    ),

  submitResults: (groupId: string, matchId: string, password: string, scoreTeam1: number, scoreTeam2: number) =>
    api.post(`/group/${groupId}/matches/${matchId}/results`, 
      { 
        score_team1: parseInt(String(scoreTeam1)), 
        score_team2: parseInt(String(scoreTeam2)) 
      }, 
      { params: { password } }
    ),
  
  getRecentMatches: (groupId: string, password: string) =>
    api.get(`/group/${groupId}/matches`, { 
      params: { 
        password,
        recent: true
      } 
    }),

  getMatches: (groupId: string, password: string, page: number = 1, pageSize: number = 10) =>
    api.get(`/group/${groupId}/matches`, { 
      params: { 
        password,
        page,
        pageSize
      } 
    }),
  
  getStatistics: (groupId: string, password: string) =>
    api.get(`/group/${groupId}/statistics`, { params: { password } }),

  exportMatchesCSV: (groupId: string, password: string) =>
    api.get(`/group/${groupId}/export/csv`, {
      params: { password },
      responseType: 'blob'
    }),
};
