import React, { useEffect, useState } from 'react';
import { AlertTriangle, ArrowRight, CheckCircle2, MoreVertical, BellOff, Trash2, Sparkles, Copy, X } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import Link from 'next/link';
import { useLanguage } from "@/lib/i18n/LanguageContext";
import { useEmailStore } from '@/lib/store/emails';
import { Email } from '@/lib/api/emails';
import { formatDistanceToNow } from 'date-fns';
import { useActionStore } from '@/lib/store/actions';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/DropdownMenu';
import { useAIReply } from '@/hooks/useAIReply';

interface SmartFeedProps {
  contextId?: string | null;
}

export function SmartFeed({ contextId }: SmartFeedProps) {
  const { t } = useLanguage();
  const { emails, isLoading, fetchEmails } = useEmailStore();

  useEffect(() => {
    // Fetch emails filtered by context if provided, or just smart filter for main dashboard
    const params = contextId 
      ? { context_id: contextId, limit: 5 }
      : { filter: 'smart', limit: 5 };
      
    fetchEmails(params);
  }, [fetchEmails, contextId]);

  if (isLoading && emails.length === 0) {
    return <div className="p-4 text-center text-slate-400">Loading Smart Feed...</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between mb-2">
        <h3 className="text-lg font-bold text-slate-800 flex items-center gap-2">
          <span className="relative flex h-3 w-3">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75"></span>
            <span className="relative inline-flex rounded-full h-3 w-3 bg-blue-500"></span>
          </span>
          {contextId ? t('dashboard.contextFeed') : t('dashboard.smartFeed')}
          <span className="text-xs font-normal text-slate-400 ml-2 bg-slate-100 px-2 py-0.5 rounded-full">{t('dashboard.priorityOnly')}</span>
        </h3>
        <Link href={contextId ? `/dashboard/inbox?context=${contextId}` : "/dashboard/inbox?filter=smart"} className="text-sm font-medium text-blue-600 hover:text-blue-700 flex items-center gap-1">
          {t('dashboard.viewAll')} <ArrowRight className="w-4 h-4" />
        </Link>
      </div>

      <div className="grid gap-4">
        {emails.length === 0 ? (
          <div className="text-center py-8 bg-slate-50 rounded-xl border border-dashed border-slate-200 text-slate-500 text-sm">
             No urgent items found in this context.
          </div>
        ) : (
          emails.map((email) => (
            <SmartFeedCard key={email.ID} item={email} />
          ))
        )}
      </div>
    </div>
  );
}

function SmartFeedCard({ item }: { item: Email }) {
  const { t } = useLanguage();
  const { approveEmail, snoozeEmail, dismissEmail } = useActionStore();
  const { mutate: generateReply, isPending: isGeneratingReply } = useAIReply();
  const [aiReply, setAiReply] = useState<{ reply: string; confidence: number } | null>(null);
  const [showReplyDialog, setShowReplyDialog] = useState(false);

  const isHighRisk = item.Urgency === 'High';

  // Parse Date
  const timeAgo = item.Date ? formatDistanceToNow(new Date(item.Date), { addSuffix: true }) : '';

  // Handle AI reply generation
  const handleReplyWithAI = () => {
    generateReply(
      {
        emailId: item.ID,
        tone: 'professional',
        context: 'brief'
      },
      {
        onSuccess: (data) => {
          setAiReply(data);
          setShowReplyDialog(true);
        },
        onError: (error) => {
          console.error('AI Reply generation failed:', error);
        }
      }
    );
  };

  // Copy reply to clipboard
  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
      // You could add a toast notification here
    } catch (err) {
      console.error('Failed to copy text:', err);
    }
  };

  return (
    <div className={`group relative bg-white rounded-xl border p-5 transition-all duration-200 hover:shadow-md
      ${isHighRisk ? 'border-red-100 bg-red-50/10' : 'border-slate-100'}
    `}>
      {/* Header: Sender & Meta */}
      <div className="flex justify-between items-start mb-3">
        <div className="flex items-center gap-3">
          <div className={`w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold overflow-hidden
            ${isHighRisk ? 'bg-red-100 text-red-600' : 'bg-blue-100 text-blue-600'}
          `}>
            {item.Sender ? item.Sender[0].toUpperCase() : '?'}
          </div>
          <div>
            <h4 className="font-semibold text-slate-800 leading-tight">{item.Subject}</h4>
            <p className="text-xs text-slate-500 mt-0.5">{item.Sender} • {timeAgo}</p>
          </div>
        </div>
        
        {/* Risk Badge */}
        {item.Urgency && item.Urgency !== 'Low' && (
          <span className={`text-[10px] font-bold px-2 py-1 rounded-full uppercase tracking-wide flex items-center gap-1
            ${item.Urgency === 'High' ? 'bg-red-100 text-red-600' : 'bg-orange-100 text-orange-600'}
          `}>
            {item.Urgency === 'High' && <AlertTriangle className="w-3 h-3" />}
            {item.Urgency} {t('dashboard.risk')}
          </span>
        )}
      </div>

      {/* AI Summary Body */}
      <div className="ml-13 pl-0">
        <div className="bg-slate-50/80 rounded-lg p-3 text-sm text-slate-700 leading-relaxed border border-slate-100/50 relative">
            <div className="absolute top-3 left-0 w-1 h-full bg-blue-500 rounded-l-lg opacity-0"></div> 
            <span className="font-semibold text-blue-600/80 mr-1">{t('dashboard.aiSummary')}</span>
            {item.Summary || item.Snippet}
        </div>

        {/* Action Footer */}
        <div className="mt-4 flex items-center justify-between">
            <div className="flex gap-2">
                 <Button 
                    onClick={() => approveEmail(item.ID)}
                    className="h-8 bg-green-600 hover:bg-green-700 text-white gap-1.5 shadow-sm shadow-green-200 px-3"
                 >
                    <CheckCircle2 className="w-3.5 h-3.5" /> {t('dashboard.approve')}
                 </Button>
                
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button 
                      variant="outline" 
                      className="h-8 text-slate-600 border border-slate-200 hover:bg-slate-50 bg-transparent px-3"
                    >
                      <MoreVertical className="w-3.5 h-3.5" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent className="w-48">
                    <DropdownMenuItem onClick={() => snoozeEmail(item.ID, '4h')}>
                      <BellOff className="w-3.5 h-3.5 mr-2 text-slate-500" /> {t('dashboard.snooze4h')}
                    </DropdownMenuItem>
                    <DropdownMenuItem onClick={() => snoozeEmail(item.ID, 'tomorrow')}>
                      <BellOff className="w-3.5 h-3.5 mr-2 text-slate-500" /> {t('dashboard.snoozeTomorrow')}
                    </DropdownMenuItem>
                    <DropdownMenuItem onClick={() => snoozeEmail(item.ID, 'next_week')}>
                      <BellOff className="w-3.5 h-3.5 mr-2 text-slate-500" /> {t('dashboard.snoozeNextWeek')}
                    </DropdownMenuItem>
                    <DropdownMenuItem className="text-red-600 focus:bg-red-50 focus:text-red-700" onClick={() => dismissEmail(item.ID)}>
                      <Trash2 className="w-3.5 h-3.5 mr-2" /> {t('dashboard.dismiss')}
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>

                 <Button
                    onClick={handleReplyWithAI}
                    disabled={isGeneratingReply}
                    className="h-8 text-blue-600 border border-blue-200 hover:bg-blue-50 bg-transparent px-3 gap-1.5"
                 >
                    {isGeneratingReply ? (
                       <div className="w-3.5 h-3.5 border border-blue-600 border-t-transparent rounded-full animate-spin"></div>
                    ) : (
                       <Sparkles className="w-3.5 h-3.5" />
                    )}
                    {t('dashboard.replyWithAI')}
                 </Button>
            </div>
        </div>
      </div>

      {/* AI Reply Dialog */}
      {showReplyDialog && aiReply && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-xl shadow-xl max-w-2xl w-full max-h-[80vh] overflow-hidden">
            {/* Header */}
            <div className="border-b border-slate-200 p-4 flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Sparkles className="w-5 h-5 text-blue-600" />
                <h3 className="font-semibold text-slate-800">AI Generated Reply</h3>
                <span className="text-xs bg-blue-100 text-blue-700 px-2 py-1 rounded-full">
                  {Math.round(aiReply.confidence * 100)}% confidence
                </span>
              </div>
              <button
                onClick={() => setShowReplyDialog(false)}
                className="text-slate-400 hover:text-slate-600"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            {/* Reply Content */}
            <div className="p-4">
              <div className="bg-slate-50 rounded-lg p-4 text-sm text-slate-700 whitespace-pre-wrap max-h-96 overflow-y-auto border border-slate-200">
                {aiReply.reply}
              </div>
            </div>

            {/* Actions */}
            <div className="border-t border-slate-200 p-4 flex justify-end gap-2">
              <Button
                onClick={() => copyToClipboard(aiReply.reply)}
                variant="outline"
                className="h-8 text-slate-600 border border-slate-200 hover:bg-slate-50 bg-transparent px-3 gap-1.5"
              >
                <Copy className="w-3.5 h-3.5" />
                Copy
              </Button>
              <Button
                onClick={() => {
                  // Here you could open email client with the reply
                  window.location.href = `mailto:?body=${encodeURIComponent(aiReply.reply)}`;
                }}
                className="h-8 bg-blue-600 hover:bg-blue-700 text-white px-3"
              >
                Open in Email
              </Button>
              <Button
                onClick={() => setShowReplyDialog(false)}
                className="h-8 bg-slate-600 hover:bg-slate-700 text-white px-3"
              >
                Close
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
