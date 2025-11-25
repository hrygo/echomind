export interface Opportunity {
  id: string;
  title: string;
  description: string;
  company: string;
  value: string;
  type: 'buying' | 'partnership' | 'renewal' | 'strategic';
  status: 'new' | 'active' | 'won' | 'lost' | 'on_hold';
  confidence: number;
  user_id: string;
  team_id: string;
  org_id: string;
  source_email_id?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateOpportunityRequest {
  title: string;
  description?: string;
  company: string;
  value?: string;
  type?: 'buying' | 'partnership' | 'renewal' | 'strategic';
  confidence?: number;
  source_email_id?: string;
}

export interface UpdateOpportunityRequest {
  title?: string;
  description?: string;
  value?: string;
  status?: 'new' | 'active' | 'won' | 'lost' | 'on_hold';
  confidence?: number;
}

// UI Helper Types
export type OpportunityType = Opportunity['type'];
export type OpportunityStatus = Opportunity['status'];

// Mock data for DealmakerView
export interface MockOpportunity {
  id: number;
  title: string;
  company: string;
  value: string;
  confidence: number;
  type: OpportunityType;
}