import { useQuery } from '@tanstack/react-query';
import api from '@/lib/api';

// Risk item interface
export interface RiskItem {
  id: string;
  title: string;
  severity: 'high' | 'medium' | 'low';
  deadline?: string;
  description?: string;
}

// Risk response interface
export interface RiskData {
  highRiskItems: RiskItem[];
  mediumRiskItems: RiskItem[];
  lowRiskItems: RiskItem[];
  riskTrend: 'increasing' | 'decreasing' | 'stable';
  totalRiskCount: number;
}

export const useRisks = () => {
  return useQuery({
    queryKey: ['risks'],
    queryFn: async (): Promise<RiskData> => {
      const { data } = await api.get('/insights/risks');
      return data;
    },
    staleTime: 3 * 60 * 1000, // 3 minutes
    refetchInterval: 5 * 60 * 1000, // Refetch every 5 minutes
  });
};