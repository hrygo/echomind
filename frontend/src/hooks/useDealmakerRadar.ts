import { useQuery } from '@tanstack/react-query';
import api from '@/lib/api';

// Radar data interface matching what the component expects
export interface RadarDataItem {
  category: string;
  value: number;
  fullMark: number;
}

export const useDealmakerRadar = () => {
  return useQuery({
    queryKey: ['dealmaker-radar'],
    queryFn: async (): Promise<RadarDataItem[]> => {
      const { data } = await api.get('/insights/dealmaker/radar');
      return data;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchInterval: 10 * 60 * 1000, // Refetch every 10 minutes
  });
};