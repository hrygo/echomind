import { create } from 'zustand';
import { Email, EmailAPI } from '../api/emails';

interface EmailState {
  emails: Email[];
  isLoading: boolean;
  error: string | null;
  
  fetchEmails: (params?: { limit?: number; offset?: number; context_id?: string; folder?: string; category?: string; filter?: string }) => Promise<void>;
}

export const useEmailStore = create<EmailState>((set, get) => ({
  emails: [],
  isLoading: false,
  error: null,

  fetchEmails: async (params) => {
    set({ isLoading: true, error: null });
    try {
      const data = await EmailAPI.list(params);
      set({ emails: data, isLoading: false });
    } catch (err: any) {
      set({ error: err.message || 'Failed to fetch emails', isLoading: false });
    }
  },
}));
