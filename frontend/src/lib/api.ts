import axios from 'axios';
import { useAuthStore } from '@/store/auth';

const apiClient = axios.create({
    baseURL: process.env.NEXT_PUBLIC_API_URL || '/api/v1',
    headers: {
        'Content-Type': 'application/json',
    },
});

apiClient.interceptors.request.use(
    (config) => {
        const token = useAuthStore.getState().token;
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
);

apiClient.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            useAuthStore.getState().logout();
            // Optional: Redirect to login page if not handled by AuthGuard
            // window.location.href = '/login'; 
        }
        return Promise.reject(error);
    }
);

// Types
export interface SearchResult {
    email_id: string;
    subject: string;
    snippet: string;
    sender: string;
    date: string;
    score: number;
}

export interface SearchResponse {
    query: string;
    results: SearchResult[];
    count: number;
}

// API Functions
export interface SearchFilters {
    sender?: string;
    startDate?: string;
    endDate?: string;
}

export const searchEmails = async (query: string, filters: SearchFilters = {}, limit: number = 10): Promise<SearchResponse> => {
    const params: Record<string, string | number> = { q: query, limit };
    if (filters.sender) params.sender = filters.sender;
    if (filters.startDate) params.start_date = filters.startDate;
    if (filters.endDate) params.end_date = filters.endDate;

    const response = await apiClient.get<SearchResponse>('/search', {
        params,
    });
    return response.data;
};

export default apiClient;
