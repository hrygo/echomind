import { useQuery } from '@tanstack/react-query';
import api from '@/lib/api';

// Executive Overview interface matching backend API
export interface ExecutiveOverview {
  totalConnections: number;
  activeProjects: number;
  teamCollaborationScore: number;
  productivityTrend: 'upward' | 'downward' | 'stable';
  criticalAlerts: number;
  upcomingDeadlines: number;
}

export const useExecutiveOverview = () => {
  return useQuery({
    queryKey: ['executive-overview'],
    queryFn: async (): Promise<ExecutiveOverview> => {
      const { data } = await api.get('/insights/executive/overview');
      return data;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchInterval: 10 * 60 * 1000, // Refetch every 10 minutes
  });
};