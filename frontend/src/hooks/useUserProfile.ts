import { useQuery } from '@tanstack/react-query'
import api from '@/lib/api'

interface UserProfile {
  id: string
  name: string
  email: string
  role: 'executive' | 'manager' | 'dealmaker'
  avatar_url?: string
  preferences?: {
    language: string
    timezone: string
  }
}

export const useUserProfile = () => {
  return useQuery({
    queryKey: ['user-profile'],
    queryFn: async (): Promise<UserProfile> => {
      const { data } = await api.get('/users/me/profile')
      return data
    },
    staleTime: 10 * 60 * 1000, // 10 minutes
    refetchOnWindowFocus: false,
  })
}