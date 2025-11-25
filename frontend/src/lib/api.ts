import axios from 'axios';
import { useAuthStore } from '@/store/auth';
import { useOrganizationStore } from '@/lib/store/organization';

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

        const currentOrgId = useOrganizationStore.getState().currentOrgId;
        if (currentOrgId) {
            config.headers['X-Organization-ID'] = currentOrgId;
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
            useOrganizationStore.getState().clearOrganizations();
            // Optional: Redirect to login page if not handled by AuthGuard
            // window.location.href = '/login'; 
        }
        return Promise.reject(error);
    }
);

// API Functions (Remaining existing functions, but export apiClient as `api` for consistency)
export const api = apiClient;
export default apiClient;

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
export const searchEmails = async (query: string, limit: number = 10): Promise<SearchResponse> => {
    const params: Record<string, string | number> = { q: query, limit };

    const response = await api.get<SearchResponse>('/search', {
        params,
    });
    return response.data;
};

