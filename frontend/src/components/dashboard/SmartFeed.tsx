import React from 'react';
import { AlertTriangle, ArrowRight, CheckCircle2, Clock } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import Link from 'next/link';

// Define the interface for our Smart Email item
export interface SmartEmailItem {
  id: string;
  subject: string;
  sender: string;
  summary: string; // AI Summary
  riskLevel: 'High' | 'Medium' | 'Low';
  suggestedAction: 'Approve' | 'Review' | 'Reply' | 'None';
  receivedAt: string;
  category: string;
}

// Mock Data Generator
const mockSmartFeed: SmartEmailItem[] = [
  {
    id: '1',
    subject: 'Q4 Budget Approval Request',
    sender: 'finance@company.com',
    summary: 'Finance dept requests approval for Q4 marketing budget increase (+15%). Deadline: Friday.',
    riskLevel: 'Medium',
    suggestedAction: 'Approve',
    receivedAt: '2h ago',
    category: 'Finance'
  },
  {
    id: '2',
    subject: 'Urgent: Client X Escalation',
    sender: 'support@company.com',
    summary: 'Key Client X is threatening to churn due to SLA breach. Immediate executive intervention required.',
    riskLevel: 'High',
    suggestedAction: 'Review',
    receivedAt: '30m ago',
    category: 'Support'
  },
  {
    id: '3',
    subject: 'Partnership Proposal - TechCorp',
    sender: 'bizdev@techcorp.com',
    summary: 'Proposal for strategic API integration. Potential revenue impact: $500k/yr.',
    riskLevel: 'Low',
    suggestedAction: 'Reply',
    receivedAt: '5h ago',
    category: 'Partnership'
  }
];

export function SmartFeed() {
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between mb-2">
        <h3 className="text-lg font-bold text-slate-800 flex items-center gap-2">
          <span className="relative flex h-3 w-3">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75"></span>
            <span className="relative inline-flex rounded-full h-3 w-3 bg-blue-500"></span>
          </span>
          Smart Feed
          <span className="text-xs font-normal text-slate-400 ml-2 bg-slate-100 px-2 py-0.5 rounded-full">Priority Only</span>
        </h3>
        <Link href="/dashboard/inbox?filter=smart" className="text-sm font-medium text-blue-600 hover:text-blue-700 flex items-center gap-1">
          View All <ArrowRight className="w-4 h-4" />
        </Link>
      </div>

      <div className="grid gap-4">
        {mockSmartFeed.map((item) => (
          <SmartFeedCard key={item.id} item={item} />
        ))}
      </div>
    </div>
  );
}

function SmartFeedCard({ item }: { item: SmartEmailItem }) {
  const isHighRisk = item.riskLevel === 'High';

  return (
    <div className={`group relative bg-white rounded-xl border p-5 transition-all duration-200 hover:shadow-md
      ${isHighRisk ? 'border-red-100 bg-red-50/10' : 'border-slate-100'}
    `}>
      {/* Header: Sender & Meta */}
      <div className="flex justify-between items-start mb-3">
        <div className="flex items-center gap-3">
          <div className={`w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold
            ${isHighRisk ? 'bg-red-100 text-red-600' : 'bg-blue-100 text-blue-600'}
          `}>
            {item.sender[0].toUpperCase()}
          </div>
          <div>
            <h4 className="font-semibold text-slate-800 leading-tight">{item.subject}</h4>
            <p className="text-xs text-slate-500 mt-0.5">{item.sender} • {item.receivedAt}</p>
          </div>
        </div>
        
        {/* Risk Badge */}
        {item.riskLevel !== 'Low' && (
          <span className={`text-[10px] font-bold px-2 py-1 rounded-full uppercase tracking-wide flex items-center gap-1
            ${item.riskLevel === 'High' ? 'bg-red-100 text-red-600' : 'bg-orange-100 text-orange-600'}
          `}>
            {item.riskLevel === 'High' && <AlertTriangle className="w-3 h-3" />}
            {item.riskLevel} Risk
          </span>
        )}
      </div>

      {/* AI Summary Body */}
      <div className="ml-13 pl-0">
        <div className="bg-slate-50/80 rounded-lg p-3 text-sm text-slate-700 leading-relaxed border border-slate-100/50 relative">
            <div className="absolute top-3 left-0 w-1 h-full bg-blue-500 rounded-l-lg opacity-0"></div> {/* Decorative bar if needed */}
            <span className="font-semibold text-blue-600/80 mr-1">AI Summary:</span>
            {item.summary}
        </div>

        {/* Action Footer */}
        <div className="mt-4 flex items-center justify-between">
            <div className="flex gap-2">
                 {item.suggestedAction === 'Approve' && (
                     <Button className="h-8 bg-green-600 hover:bg-green-700 text-white gap-1.5 shadow-sm shadow-green-200 px-3">
                        <CheckCircle2 className="w-3.5 h-3.5" /> Approve
                     </Button>
                 )}
                 <Button className="h-8 text-slate-600 border border-slate-200 hover:bg-slate-50 bg-transparent px-3">
                    Reply with AI
                 </Button>
            </div>
            
            {/* Suggested Action Label */}
            <div className="text-xs font-medium text-slate-400 flex items-center gap-1.5">
                <Clock className="w-3.5 h-3.5" />
                Suggested: <span className="text-slate-600">{item.suggestedAction}</span>
            </div>
        </div>
      </div>
    </div>
  );
}
