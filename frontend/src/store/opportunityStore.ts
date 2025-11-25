import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export interface Opportunity {
  id: string;
  title: string;
  description: string;
  company: string;
  value: string;
  type: 'buying' | 'partnership' | 'renewal' | 'strategic';
  status: 'new' | 'active' | 'won' | 'lost' | 'on_hold';
  confidence: number;
  user_id: string;
  team_id: string;
  org_id: string;
  source_email_id?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateOpportunityRequest {
  title: string;
  description?: string;
  company: string;
  value?: string;
  type?: 'buying' | 'partnership' | 'renewal' | 'strategic';
  confidence?: number;
  source_email_id?: string;
}

export interface UpdateOpportunityRequest {
  title?: string;
  description?: string;
  value?: string;
  status?: 'new' | 'active' | 'won' | 'lost' | 'on_hold';
  confidence?: number;
}

interface OpportunityState {
  opportunities: Opportunity[];
  isLoading: boolean;
  error: string | null;

  // Actions
  fetchOpportunities: (filters?: {
    status?: string;
    type?: string;
    limit?: number;
    offset?: number;
  }) => Promise<void>;
  createOpportunity: (data: CreateOpportunityRequest) => Promise<Opportunity>;
  updateOpportunity: (id: string, data: UpdateOpportunityRequest) => Promise<Opportunity>;
  deleteOpportunity: (id: string) => Promise<void>;
  getOpportunity: (id: string) => Promise<Opportunity>;
  clearError: () => void;
}

export const useOpportunityStore = create<OpportunityState>()(
  persist(
    (set) => ({
      opportunities: [],
      isLoading: false,
      error: null,

      fetchOpportunities: async (filters = {}) => {
        set({ isLoading: true, error: null });
        try {
          const params = new URLSearchParams();

          if (filters.status) params.append('status', filters.status);
          if (filters.type) params.append('type', filters.type);
          if (filters.limit) params.append('limit', filters.limit.toString());
          if (filters.offset) params.append('offset', filters.offset.toString());

          const response = await fetch(`/api/v1/opportunities?${params}`, {
            credentials: 'include',
          });

          if (!response.ok) {
            throw new Error('Failed to fetch opportunities');
          }

          const data = await response.json();
          set({ opportunities: data, isLoading: false });
        } catch (error) {
          set({
            error: error instanceof Error ? error.message : 'Failed to fetch opportunities',
            isLoading: false
          });
        }
      },

      createOpportunity: async (data) => {
        set({ isLoading: true, error: null });
        try {
          const response = await fetch('/api/v1/opportunities', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify(data),
          });

          if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Failed to create opportunity');
          }

          const newOpportunity = await response.json();
          set(state => ({
            opportunities: [...state.opportunities, newOpportunity],
            isLoading: false
          }));

          return newOpportunity;
        } catch (error) {
          set({
            error: error instanceof Error ? error.message : 'Failed to create opportunity',
            isLoading: false
          });
          throw error;
        }
      },

      updateOpportunity: async (id, data) => {
        set({ isLoading: true, error: null });
        try {
          const response = await fetch(`/api/v1/opportunities/${id}`, {
            method: 'PATCH',
            headers: {
              'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify(data),
          });

          if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Failed to update opportunity');
          }

          const updatedOpportunity = await response.json();
          set(state => ({
            opportunities: state.opportunities.map(opp =>
              opp.id === id ? updatedOpportunity : opp
            ),
            isLoading: false
          }));

          return updatedOpportunity;
        } catch (error) {
          set({
            error: error instanceof Error ? error.message : 'Failed to update opportunity',
            isLoading: false
          });
          throw error;
        }
      },

      deleteOpportunity: async (id) => {
        set({ isLoading: true, error: null });
        try {
          const response = await fetch(`/api/v1/opportunities/${id}`, {
            method: 'DELETE',
            credentials: 'include',
          });

          if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Failed to delete opportunity');
          }

          set(state => ({
            opportunities: state.opportunities.filter(opp => opp.id !== id),
            isLoading: false
          }));
        } catch (error) {
          set({
            error: error instanceof Error ? error.message : 'Failed to delete opportunity',
            isLoading: false
          });
          throw error;
        }
      },

      getOpportunity: async (id) => {
        set({ isLoading: true, error: null });
        try {
          const response = await fetch(`/api/v1/opportunities/${id}`, {
            credentials: 'include',
          });

          if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Failed to get opportunity');
          }

          const opportunity = await response.json();
          set({ isLoading: false });
          return opportunity;
        } catch (error) {
          set({
            error: error instanceof Error ? error.message : 'Failed to get opportunity',
            isLoading: false
          });
          throw error;
        }
      },

      clearError: () => set({ error: null }),
    }),
    {
      name: 'opportunity-storage',
      partialize: (state) => ({
        opportunities: state.opportunities
      }),
    }
  )
);