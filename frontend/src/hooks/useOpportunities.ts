import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from '@/lib/api';
import { Opportunity, CreateOpportunityRequest, UpdateOpportunityRequest } from '@/types/opportunity';

export const useOpportunities = (filters?: {
  status?: string;
  type?: string;
  limit?: number;
  offset?: number;
}) => {
  return useQuery({
    queryKey: ['opportunities', filters],
    queryFn: async (): Promise<Opportunity[]> => {
      const params = new URLSearchParams();

      if (filters?.status) params.append('status', filters.status);
      if (filters?.type) params.append('type', filters.type);
      if (filters?.limit) params.append('limit', filters.limit.toString());
      if (filters?.offset) params.append('offset', filters.offset.toString());

      const { data } = await api.get(`/opportunities?${params}`);
      return data;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchInterval: 2 * 60 * 1000, // Refetch every 2 minutes
  });
};

export const useOpportunity = (id: string) => {
  return useQuery({
    queryKey: ['opportunity', id],
    queryFn: async (): Promise<Opportunity> => {
      const { data } = await api.get(`/opportunities/${id}`);
      return data;
    },
    enabled: !!id,
    staleTime: 2 * 60 * 1000,
  });
};

export const useCreateOpportunity = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateOpportunityRequest): Promise<Opportunity> => {
      const { data: created } = await api.post('/opportunities', data);
      return created;
    },
    onSuccess: (newOpportunity) => {
      queryClient.invalidateQueries({ queryKey: ['opportunities'] });
      queryClient.setQueryData(['opportunity', newOpportunity.id], newOpportunity);
    },
    onError: (error) => {
      console.error('Failed to create opportunity:', error);
    },
  });
};

export const useUpdateOpportunity = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: UpdateOpportunityRequest }): Promise<Opportunity> => {
      const { data: updated } = await api.patch(`/opportunities/${id}`, data);
      return updated;
    },
    onSuccess: (updatedOpportunity, { id }) => {
      queryClient.invalidateQueries({ queryKey: ['opportunities'] });
      queryClient.setQueryData(['opportunity', id], updatedOpportunity);
    },
    onError: (error) => {
      console.error('Failed to update opportunity:', error);
    },
  });
};

export const useDeleteOpportunity = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string): Promise<void> => {
      await api.delete(`/opportunities/${id}`);
    },
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ['opportunities'] });
      queryClient.removeQueries({ queryKey: ['opportunity', id] });
    },
    onError: (error) => {
      console.error('Failed to delete opportunity:', error);
    },
  });
};