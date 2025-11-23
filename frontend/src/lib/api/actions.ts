import { api } from '../api';

export const ActionAPI = {
  approve: async (emailId: string): Promise<void> => {
    await api.post('/actions/approve', { email_id: emailId });
  },

  snooze: async (emailId: string, duration: string): Promise<{ status: string; until: string }> => {
    const response = await api.post('/actions/snooze', { email_id: emailId, duration });
    return response.data;
  },

  dismiss: async (emailId: string): Promise<void> => {
    await api.post('/actions/dismiss', { email_id: emailId });
  },
};
