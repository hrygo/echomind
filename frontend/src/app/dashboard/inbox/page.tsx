"use client";

import { useEffect, useState, useCallback } from "react";
import Link from "next/link";
import { useSearchParams, useRouter } from 'next/navigation';
import { api } from "@/lib/api";
import { Button } from "@/components/ui/button"; // Assuming a global Button component with good styling
import { isAxiosError } from 'axios';

import { useLanguage } from "@/lib/i18n/LanguageContext";
import { useToast } from "@/lib/hooks/useToast";
import { useConfirm } from "@/components/ui/confirm-dialog";

interface Email {
  ID: string;
  Subject: string;
  Sender: string;
  Date: string;
  Summary: string;
  Sentiment: string;
  Urgency: string;
  Category: string;
  ActionItems: string[];
}

export default function DashboardPage() {
  const [emails, setEmails] = useState<Email[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [syncLoading, setSyncLoading] = useState(false);
  const searchParams = useSearchParams();
  const router = useRouter();
  const { t } = useLanguage();
  const toast = useToast();
  const confirm = useConfirm();
  const categoryFilter = searchParams.get('category');
  const folderFilter = searchParams.get('folder');
  const smartFilter = searchParams.get('filter'); // 'smart' for Smart Inbox

  const getTitle = () => {
    if (smartFilter === 'smart') return t('sidebar.smartInbox');
    if (categoryFilter === 'Work') return t('sidebar.work');
    if (categoryFilter === 'Personal') return t('sidebar.personal');
    if (categoryFilter === 'Newsletter') return t('sidebar.newsletter');
    if (categoryFilter === 'Notification') return t('sidebar.notification');
    if (categoryFilter === 'Spam') return t('sidebar.spam');
    if (categoryFilter === 'Finance') return t('sidebar.finance');
    if (folderFilter === 'sent') return t('sidebar.sent');
    if (folderFilter === 'drafts') return t('sidebar.drafts');
    if (folderFilter === 'trash') return t('sidebar.trash');
    return t('sidebar.inbox');
  };

  const fetchEmails = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      let url = "/emails";
      if (categoryFilter) {
        url += `?category=${categoryFilter}`;
      } else if (folderFilter) {
        url += `?folder=${folderFilter}`;
      } else if (smartFilter === 'smart') {
        // For smart inbox, we might want a different API or a filter for existing emails
        url += `?filter=smart`;
      }
      const response = await api.get<Email[]>(url);
      setEmails(response.data);
    } catch (err: unknown) {
      console.error("Error fetching emails:", err);
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("Failed to load emails.");
      }
    }
    finally {
      setLoading(false);
    }
  }, [categoryFilter, folderFilter, smartFilter]);

  useEffect(() => {
    fetchEmails();
  }, [fetchEmails]);

  const handleSync = async () => {
    if (syncLoading) return;

    setSyncLoading(true);
    try {
      await api.post<{ message: string }>("/sync");
      toast.success(t('inbox.syncStarted'));
      fetchEmails();
    } catch (error: unknown) {
      console.error("Sync error:", error);
      if (isAxiosError(error) && error.response?.status === 400) {
        confirm(
          t('inbox.noAccountConfigured'),
          () => {
            router.push('/dashboard/settings');
          },
          {
            title: t('inbox.configureAccount'),
            confirmText: t('common.goToSettings'),
            cancelText: t('common.cancel')
          }
        );
      } else if (error instanceof Error) {
        toast.error(`${t('inbox.syncFailed')}: ${error.message}`);
      } else {
        toast.error(t('inbox.syncFailedUnknown'));
      }
    } finally {
      setSyncLoading(false);
    }
  };

  const getSentimentLabel = (sentiment: string) => {
    const map: Record<string, string> = {
      'Positive': '积极',
      'Negative': '消极',
      'Neutral': '中性'
    };
    return map[sentiment] || sentiment;
  };

  const getUrgencyLabel = (urgency: string) => {
    const map: Record<string, string> = {
      'High': '高优先级',
      'Medium': '中优先级',
      'Low': '低优先级'
    };
    return map[urgency] || urgency;
  };

  return (
    <div className="min-h-[calc(100vh-theme(spacing.20)-theme(spacing.16))] flex flex-col">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-slate-800 tracking-tight">{getTitle()}</h1>
        <Button
          onClick={handleSync}
          disabled={syncLoading}
          className="inline-flex items-center gap-2 px-5 py-2.5 rounded-full bg-blue-600 hover:bg-blue-700 text-white font-semibold shadow-md hover:shadow-lg transition-all duration-200 disabled:opacity-60 disabled:cursor-not-allowed"
        >
          {syncLoading ? (
            <svg className="animate-spin h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          ) : (
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="lucide lucide-rotate-cw"><path d="M21 12a9 9 0 0 0-9-9V3a10 10 0 0 1 10 10Z"/><path d="M21 21v-3.5L16 18"/><path d="M3 12a9 9 0 0 0 9 9v0a10 10 0 0 1-10-10Z"/><path d="M3 3v3.5L8 6"/></svg>
          )}
          {syncLoading ? "同步中..." : "立即同步"}
        </Button>
      </div>

      <div className="flex-1 bg-white shadow-lg rounded-2xl overflow-hidden border border-slate-100">
        {loading ? (
          <div className="p-8 text-center text-slate-500 text-lg flex items-center justify-center h-full">
            <svg className="animate-spin h-6 w-6 mr-3 text-blue-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            正在加载邮件...
          </div>
        ) : error ? (
          <div className="p-8 text-center text-red-500 text-lg flex items-center justify-center h-full">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="lucide lucide-alert-triangle mr-2 text-red-500"><path d="M10.29 3.86L1.86 18.14a2 2 0 0 0 1.74 3.09h16.8a2 2 0 0 0 1.74-3.09L13.71 3.86a2 2 0 0 0-3.42 0Z"/><path d="M12 9v4"/><path d="M12 17h.01"/></svg>
            错误: {error}
          </div>
        ) : emails.length === 0 ? (
          <div className="p-8 text-center text-slate-500 text-lg flex items-center justify-center h-full">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="lucide lucide-inbox mr-2 text-slate-400"><polyline points="22 12 16 12 14 15 10 15 8 12 2 12"/><path d="M5.45 5.11L2 12v6a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-6l-3.45-6.89A2 2 0 0 0 16.76 4H7.24a2 2 0 0 0-1.79 1.11z"/></svg>
            暂无邮件，请尝试同步！
          </div>
        ) : (
          <ul className="divide-y divide-slate-100">
            {emails.map((email) => (
              <li key={email.ID} className="group relative transition-all duration-200">
                <Link href={`/dashboard/email/${email.ID}`} className="block p-6 hover:bg-slate-50 active:bg-slate-100">
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm font-semibold text-slate-800 group-hover:text-blue-600 transition-colors mr-4 truncate flex-1">
                      {email.Sender}
                    </span>
                    <span className="text-xs text-slate-500 font-medium whitespace-nowrap">
                      {new Date(email.Date).toLocaleString('zh-CN', { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })}
                    </span>
                  </div>
                  <h3 className="text-lg font-bold text-slate-900 mb-2 leading-tight">
                    {email.Subject}
                  </h3>
                  <div className="flex items-start gap-3 mt-3">
                    <p className="text-sm text-slate-600 line-clamp-2 flex-1">
                      {email.Summary || "正在分析..."}
                    </p>
                    <div className="flex flex-shrink-0 flex-col items-end gap-1.5">
                      {email.Sentiment && (
                        <span className={`text-xs font-semibold px-3 py-1 rounded-full transition-colors
                          ${email.Sentiment === 'Positive' ? 'bg-green-50 text-green-700 border border-green-200'
                            : email.Sentiment === 'Negative' ? 'bg-red-50 text-red-700 border border-red-200'
                            : 'bg-slate-50 text-slate-600 border border-slate-200'}`}>
                          {getSentimentLabel(email.Sentiment)}
                        </span>
                      )}
                      {email.Urgency && (
                        <span className={`text-xs font-semibold px-3 py-1 rounded-full transition-colors
                          ${email.Urgency === 'High' ? 'bg-orange-50 text-orange-700 border border-orange-200'
                            : email.Urgency === 'Medium' ? 'bg-yellow-50 text-yellow-700 border border-yellow-200'
                            : 'bg-blue-50 text-blue-700 border border-blue-200'}`}>
                          {getUrgencyLabel(email.Urgency)}
                        </span>
                      )}
                    </div>
                  </div>
                </Link>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}