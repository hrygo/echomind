import { create } from 'zustand';
import { Context, ContextAPI, ContextInput } from '../api/contexts';

interface ContextState {
  contexts: Context[];
  isLoading: boolean;
  error: string | null;
  
  fetchContexts: () => Promise<void>;
  addContext: (data: ContextInput) => Promise<void>;
  updateContext: (id: string, data: ContextInput) => Promise<void>;
  deleteContext: (id: string) => Promise<void>;
}

export const useContextStore = create<ContextState>((set) => ({
  contexts: [],
  isLoading: false,
  error: null,

  fetchContexts: async () => {
    set({ isLoading: true, error: null });
    try {
      const contexts = await ContextAPI.list();
      set({ contexts, isLoading: false });
    } catch {
      set({ error: 'Failed to fetch contexts', isLoading: false });
    }
  },

  addContext: async (data: ContextInput) => {
    set({ isLoading: true, error: null });
    try {
      const newContext = await ContextAPI.create(data);
      set(state => ({ 
        contexts: [newContext, ...state.contexts],
        isLoading: false 
      }));
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to create context';
      set({ error: errorMessage, isLoading: false });
      throw err;
    }
  },

  updateContext: async (id: string, data: ContextInput) => {
    set({ isLoading: true, error: null });
    try {
      const updatedContext = await ContextAPI.update(id, data);
      set(state => ({
        contexts: state.contexts.map(c => c.ID === id ? updatedContext : c),
        isLoading: false
      }));
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to update context';
      set({ error: errorMessage, isLoading: false });
      throw err;
    }
  },

  deleteContext: async (id: string) => {
    set({ isLoading: true, error: null });
    try {
      await ContextAPI.delete(id);
      set(state => ({
        contexts: state.contexts.filter(c => c.ID !== id),
        isLoading: false
      }));
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to delete context';
      set({ error: errorMessage, isLoading: false });
      throw err;
    }
  },
}));
