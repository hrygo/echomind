import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { api } from '@/lib/api'; // Import api client

interface User {
    id: string;
    email: string;
    name?: string;
    role?: string; // Add role
    has_account?: boolean; // Add has_account
}

interface AuthState {
    user: User | null;
    token: string | null;
    isAuthenticated: boolean;
    isHydrated: boolean;
    setAuth: (token: string, user: User) => void;
    logout: () => void;
    setHydrated: () => void;
    login: (email: string, password: string) => Promise<void>; // Add login method
    register: (name: string, email: string, password: string) => Promise<void>; // Add register method
    updateUser: (updates: Partial<User>) => void; // Add updateUser method
}

export const useAuthStore = create<AuthState>()(
    persist(
        (set) => ({
            user: null,
            token: null,
            isAuthenticated: false,
            isHydrated: false,
            setAuth: (token, user) => set({ token, user, isAuthenticated: true }),
            logout: () => set({ token: null, user: null, isAuthenticated: false }),
            setHydrated: () => set({ isHydrated: true }),
            updateUser: (updates) => set((state) => ({
                user: state.user ? { ...state.user, ...updates } : null
            })),
            login: async (email, password) => {
                try {
                    const response = await api.post('/auth/login', { email, password });
                    const { token, user } = response.data;
                    set({ token, user, isAuthenticated: true });
                } catch (error) {
                    // Re-throw to be caught by AuthForm
                    throw error;
                }
            },
            register: async (name, email, password) => {
                try {
                    const response = await api.post('/auth/register', { name, email, password });
                    const { token, user } = response.data;
                    set({ token, user, isAuthenticated: true });
                } catch (error) {
                    // Re-throw to be caught by AuthForm
                    throw error;
                }
            },
        }),
        {
            name: 'auth-storage',
            onRehydrateStorage: () => (state) => {
                state?.setHydrated();
            },
        }
    )
);
