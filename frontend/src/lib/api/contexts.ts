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

export interface RawContext extends Omit<Context, 'Keywords' | 'Stakeholders'> {
  Keywords: string | string[];
  Stakeholders: string | string[];
}

export const ContextAPI = {
  list: async (): Promise<Context[]> => {
    const response = await api.get('/contexts');
    return response.data.map((ctx: RawContext) => ({
      ...ctx,
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
