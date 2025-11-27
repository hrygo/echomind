"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { api } from "@/lib/api";
import { Button } from "@/components/ui/button";
import AIDraftReplyModal from "@/components/email/AIDraftReplyModal";
import { useLanguage } from "@/lib/i18n/LanguageContext";
import { createTask } from "@/lib/api/tasks"; // Import the new task API
import { useTaskStore } from "@/store/task"; // Import the task store

interface SmartAction {
  type: string;
  label: string;
  data: Record<string, string>;
}

interface EmailDetail {
  ID: string; // Changed from number to string (UUID)
  Subject: string;
  Sender: string;
  Date: string;
  BodyText: string;
  Summary: string;
  Sentiment: string;
  Urgency: string;
  Category: string; // New field from AI analysis
  ActionItems: string[]; // New field from AI analysis
  SmartActions?: SmartAction[];
}

export default function EmailDetailPage() {
  const params = useParams();
  const { id } = params;
  const [email, setEmail] = useState<EmailDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const { t } = useLanguage();
  const { addTask } = useTaskStore(); // Use addTask from the store

  useEffect(() => {
    if (!id) return;

    async function fetchEmail() {
      try {
        // Ensure `id` is treated as a string for UUID
        const response = await api.get<EmailDetail>(`/emails/${id}`);
        setEmail(response.data);
      } catch (err: unknown) {
        console.error("Error fetching email:", err);
        if (err instanceof Error) {
          setError(err.message);
        } else {
          setError("Failed to load email.");
        }
      } finally {
        setLoading(false);
      }
    }

    fetchEmail();
  }, [id]);

  const getActionLabel = (action: SmartAction) => {
    if (action.type === 'calendar_event') return t('emailDetail.addToCalendar');
    if (action.type === 'create_task') return t('emailDetail.createTask');
    return action.label;
  };

  // 翻译情感值
  const getSentimentLabel = (sentiment: string) => {
    const key = sentiment.toLowerCase() as 'positive' | 'neutral' | 'negative';
    return t(`emailDetail.sentiment.${key}`);
  };

  // 翻译紧急程度
  const getUrgencyLabel = (urgency: string) => {
    const key = urgency.toLowerCase() as 'high' | 'medium' | 'low';
    return t(`emailDetail.urgency.${key}`);
  };

  const handleSmartAction = async (action: SmartAction) => {
    if (action.type === 'create_task') {
      try {
        const newTask = await createTask({
          title: action.data.title || t('emailDetail.createTaskDefault'), // Fallback title
          description: `From email: ${email?.Subject || ''}`,
          source_email_id: email?.ID,
          due_date: action.data.deadline, // Use 'deadline' for due_date
        });
        addTask(newTask); // Add task to Zustand store
        alert(t('emailDetail.taskCreatedSuccess')); // Use i18n for success message
      } catch (err: unknown) {
        console.error("Failed to create task:", err);
        alert(t('emailDetail.taskCreatedError')); // Use i18n for error message
      }
    } else if (action.type === 'calendar_event') {
      // For calendar events, still use alert for now as per plan
      alert(`Executing action: ${getActionLabel(action)}\nData: ${JSON.stringify(action.data, null, 2)}`);
    } else {
      alert(`Unknown action type: ${action.type}`);
    }
  };

  if (loading) return <div className="p-8">{t('common.loading')}</div>;
  if (error) return <div className="p-8 text-red-500">{t('common.error')}: {error}</div>;
  if (!email) return <div className="p-8">{t('common.noResults')}</div>;

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
      <div className="mb-6">
        <Link href="/dashboard" className="text-blue-600 hover:underline break-words">&larr; {t('emailDetail.backToInbox')}</Link>
      </div>

      {/* AI Insight Card */}
      {email.Summary && (
        <div className="bg-gradient-to-r from-indigo-50 to-blue-50 border border-indigo-100 rounded-lg p-4 sm:p-6 mb-8 shadow-sm overflow-hidden">
          <div className="flex flex-col sm:flex-row justify-between items-start gap-3 sm:gap-0 mb-4">
            <h2 className="text-lg font-semibold text-indigo-900 flex items-center gap-2 min-w-0">
              <span className="flex-shrink-0">✨</span>
              <span className="truncate">{t('emailDetail.aiInsights')}</span>
            </h2>
            <div className="flex gap-2 flex-wrap">
              {/* Sentiment Badge */}
              {email.Sentiment && (
                <span className={`text-xs font-medium px-2.5 py-0.5 rounded border whitespace-nowrap
                   ${email.Sentiment === 'Positive' ? 'bg-green-100 text-green-800 border-green-200' :
                    email.Sentiment === 'Negative' ? 'bg-red-100 text-red-800 border-red-200' :
                      'bg-gray-100 text-gray-800 border-gray-200'}`}>
                  {getSentimentLabel(email.Sentiment)}
                </span>
              )}
              {/* Urgency Badge */}
              {email.Urgency && (
                <span className={`text-xs font-medium px-2.5 py-0.5 rounded border whitespace-nowrap
                   ${email.Urgency === 'High' ? 'bg-orange-100 text-orange-800 border-orange-200' :
                    'bg-blue-100 text-blue-800 border-blue-200'}`}>
                  {getUrgencyLabel(email.Urgency)}
                </span>
              )}
            </div>
          </div>
          <div className="text-indigo-800 overflow-hidden">
            <p className="font-medium mb-1">{t('emailDetail.summary')}:</p>
            <p className="mb-4 break-words">{email.Summary}</p>

            {email.ActionItems && email.ActionItems.length > 0 && (
              <div className="mt-4 overflow-hidden">
                <p className="font-medium mb-1">{t('emailDetail.actionItems')}:</p>
                <ul className="list-disc list-inside space-y-1">
                  {email.ActionItems.map((item, index) => (
                    <li key={index} className="break-words">{item}</li>
                  ))}
                </ul>
              </div>
            )}

            {/* Smart Actions */}
            {email.SmartActions && email.SmartActions.length > 0 && (
              <div className="mt-4 border-t border-indigo-100 pt-4 flex flex-wrap gap-2">
                {email.SmartActions.map((action, idx) => (
                  <button
                    key={idx}
                    onClick={() => handleSmartAction(action)}
                    className="flex items-center gap-1.5 px-3 py-1.5 bg-white hover:bg-indigo-50 text-indigo-600 border border-indigo-200 rounded-lg text-xs font-medium transition-colors shadow-sm whitespace-nowrap"
                  >
                    <span className="flex-shrink-0">{action.type === 'calendar_event' ? '📅' : '✅'}</span>
                    <span className="truncate">{getActionLabel(action)}</span>
                  </button>
                ))}
              </div>
            )}
          </div>
          <div className="mt-4 text-right">
            <Button onClick={() => setIsModalOpen(true)} className="w-full sm:w-auto">{t('emailDetail.generateDraft')}</Button>
          </div>
        </div>
      )}

      {/* Email Header */}
      <div className="bg-white shadow rounded-lg p-4 sm:p-8 overflow-hidden">
        <h1 className="text-xl sm:text-2xl font-bold text-gray-900 mb-4 break-words">{email.Subject}</h1>
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-2 text-sm text-gray-500 mb-8 border-b pb-4">
          <div className="min-w-0 w-full sm:w-auto">
            <span className="font-medium text-gray-900">{t('emailDetail.from')}</span> 
            <span className="break-all">{email.Sender}</span>
          </div>
          <div className="whitespace-nowrap flex-shrink-0">{new Date(email.Date).toLocaleString()}</div>
        </div>

        {/* Email Body */}
        <div className="prose max-w-none text-gray-800 whitespace-pre-wrap break-words overflow-wrap-anywhere">
          {email.BodyText}
        </div>
      </div>
      <AIDraftReplyModal
        emailContent={email.BodyText}
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
      />
    </div>
  );
}
