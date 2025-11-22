import { create } from 'zustand';
import { persist, devtools } from 'zustand/middleware';

export interface Organization {
  id: string;
  name: string;
  slug: string;
  owner_id: string;
  created_at: string;
  updated_at: string;
}

interface OrganizationState {
  organizations: Organization[];
  currentOrgId: string | null;
  setOrganizations: (orgs: Organization[]) => void;
  setCurrentOrg: (orgId: string) => void;
  addOrganization: (org: Organization) => void;
  clearOrganizations: () => void;
}

export const useOrganizationStore = create<OrganizationState>()(
  devtools(
    persist(
      (set) => ({
        organizations: [],
        currentOrgId: null,
        setOrganizations: (orgs) => set({ organizations: orgs, currentOrgId: orgs.length > 0 ? orgs[0].id : null }),
        setCurrentOrg: (orgId) => set({ currentOrgId: orgId }),
        addOrganization: (org) => set((state) => ({ organizations: [...state.organizations, org] })),
        clearOrganizations: () => set({ organizations: [], currentOrgId: null }),
      }),
      {
        name: 'organization-storage',
        // @ts-expect-error: localStorage is compatible but TypeScript struggles with the generic PersistStorage interface
        storage: localStorage,
      }
    )
  )
);

