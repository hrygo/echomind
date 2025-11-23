import { api } from '../api';

export interface Context {
  ID: string;
  Name: string;
  Color: string;
  Keywords: string[]; // JSON array from backend
  Stakeholders: string[]; // JSON array from backend
  CreatedAt: string;
  UpdatedAt: string;
}

export interface ContextInput {
  name: string;
  color: string;
  keywords: string[];
  stakeholders: string[];
}

export const ContextAPI = {
  list: async (): Promise<Context[]> => {
    const response = await api.get('/contexts');
    return response.data.map((ctx: any) => ({
      ...ctx,
      // Parse JSON fields if they come as strings, though usually axios handles JSON response
      // But GORM might send them as actual JSON objects if configured right, 
      // or strings if using datatypes.JSON. Let's assume standard JSON response.
      Keywords: typeof ctx.Keywords === 'string' ? JSON.parse(ctx.Keywords) : ctx.Keywords,
      Stakeholders: typeof ctx.Stakeholders === 'string' ? JSON.parse(ctx.Stakeholders) : ctx.Stakeholders,
    }));
  },

  create: async (data: ContextInput): Promise<Context> => {
    const response = await api.post('/contexts', data);
    return response.data;
  },

  update: async (id: string, data: ContextInput): Promise<Context> => {
    const response = await api.patch(`/contexts/${id}`, data);
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/contexts/${id}`);
  },
};
