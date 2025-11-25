import { useQuery } from '@tanstack/react-query';
import api from '@/lib/api';

// Trend data point interface
export interface TrendDataPoint {
  date: string;
  value: number;
}

// Trends response interface
export interface TrendsData {
  productivity: TrendDataPoint[];
  collaboration: TrendDataPoint[];
  communication: TrendDataPoint[];
  weeklyInteraction: number;
  interactionChange: number; // percentage change
}

export const useTrends = () => {
  return useQuery({
    queryKey: ['trends'],
    queryFn: async (): Promise<TrendsData> => {
      const { data } = await api.get('/insights/trends');
      return data;
    },
    staleTime: 10 * 60 * 1000, // 10 minutes
    refetchInterval: 15 * 60 * 1000, // Refetch every 15 minutes
  });
};