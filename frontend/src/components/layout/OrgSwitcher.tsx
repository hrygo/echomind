'use client';

import { useState } from 'react';
import { useOrganizationStore, Organization } from '@/lib/store/organization';
import { Button } from '@/components/ui/Button';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/DropdownMenu';
import { PlusCircle, Check, ChevronsUpDown } from 'lucide-react';
import { cn } from '@/lib/utils';
import { CreateOrganizationModal } from './CreateOrganizationModal';

export function OrgSwitcher() {
  const { organizations, currentOrgId, setCurrentOrg, addOrganization } = useOrganizationStore();
  const [isModalOpen, setIsModalOpen] = useState(false);

  const currentOrg = organizations.find((org: Organization) => org.id === currentOrgId);

  return (
    <div className="flex flex-col gap-2 p-2">
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            variant="ghost"
            role="combobox"
            aria-expanded={false} // State managed by DropdownMenu
            className="w-full justify-between pr-2"
          >
            {currentOrg ? currentOrg.name : "Select Organization..."}
            <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-[var(--radix-dropdown-menu-trigger-width)]">
          {organizations.map((org: Organization) => (
            <DropdownMenuItem
              key={org.id}
              onSelect={() => setCurrentOrg(org.id)}
              className="flex items-center justify-between"
            >
              {org.name}
              <Check
                className={cn(
                  "ml-auto h-4 w-4",
                  currentOrgId === org.id ? "opacity-100" : "opacity-0"
                )}
              />
            </DropdownMenuItem>
          ))}
          <DropdownMenuItem onSelect={() => setIsModalOpen(true)}>
            <PlusCircle className="mr-2 h-4 w-4" />
            Create New Organization
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      <CreateOrganizationModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onOrganizationCreated={(org) => {
          addOrganization(org);
          setCurrentOrg(org.id);
          setIsModalOpen(false);
        }}
      />
    </div>
  );
}
