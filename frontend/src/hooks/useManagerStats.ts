import { useQuery } from '@tanstack/react-query'
import api from '@/lib/api'

interface ManagerStats {
  activeTasksCount: number
  overdueTasksCount: number
  completedTodayCount: number
  teamProductivity: number
  urgentEmailsCount: number
}

export const useManagerStats = () => {
  return useQuery({
    queryKey: ['manager-stats'],
    queryFn: async (): Promise<ManagerStats> => {
      const { data } = await api.get('/insights/manager/stats')
      return data
    },
    staleTime: 2 * 60 * 1000, // 2 minutes
    refetchInterval: 5 * 60 * 1000, // Refetch every 5 minutes
  })
}