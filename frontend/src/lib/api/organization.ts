import { api } from '../api';
import { Organization } from '@/lib/store/organization';
export const organizationApi = {
  getOrganizations: async (): Promise<Organization[]> => {
    const response = await api.get<'' | Organization[]>(`/orgs`);
    // Handle empty string response from API for no organizations
    if (typeof response.data === 'string') {
      return [];
    }
    return response.data;
  },
  createOrganization: async (name: string): Promise<Organization> => {
    const response = await api.post<Organization>(`/orgs`, { name });
    return response.data;
  },
};
