import { api } from '../api';

export interface Email {
  ID: string;
  Subject: string;
  Sender: string;
  Snippet: string;
  BodyText: string;
  Date: string;
  Summary: string;
  Category: string;
  Sentiment: string;
  Urgency: string;
  IsRead: boolean;
  ActionItems: string[]; // JSON
  SmartActions: any; // JSON
}

export const EmailAPI = {
  list: async (params?: { limit?: number; offset?: number; context_id?: string; folder?: string; category?: string; filter?: string }): Promise<Email[]> => {
    const response = await api.get('/emails', { params });
    return response.data;
  },

  get: async (id: string): Promise<Email> => {
    const response = await api.get(`/emails/${id}`);
    return response.data;
  },
};
