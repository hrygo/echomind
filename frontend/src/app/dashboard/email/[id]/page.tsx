"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { api } from "@/lib/api";
import { Button } from "@/components/ui/Button";
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
      } catch (err: any) {
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
    <div className="max-w-4xl mx-auto">
      <div className="mb-6">
        <Link href="/dashboard" className="text-blue-600 hover:underline">&larr; {t('emailDetail.backToInbox')}</Link>
      </div>

      {/* AI Insight Card */}
      {email.Summary && (
        <div className="bg-gradient-to-r from-indigo-50 to-blue-50 border border-indigo-100 rounded-lg p-6 mb-8 shadow-sm">
          <div className="flex justify-between items-start mb-4">
            <h2 className="text-lg font-semibold text-indigo-900 flex items-center gap-2">
              ✨ {t('emailDetail.aiInsights')}
            </h2>
            <div className="flex gap-2">
              {/* Sentiment Badge */}
              {email.Sentiment && (
                <span className={`text-xs font-medium px-2.5 py-0.5 rounded border 
                   ${email.Sentiment === 'Positive' ? 'bg-green-100 text-green-800 border-green-200' :
                    email.Sentiment === 'Negative' ? 'bg-red-100 text-red-800 border-red-200' :
                      'bg-gray-100 text-gray-800 border-gray-200'}`}>
                  {email.Sentiment}
                </span>
              )}
              {/* Urgency Badge */}
              {email.Urgency && (
                <span className={`text-xs font-medium px-2.5 py-0.5 rounded border 
                   ${email.Urgency === 'High' ? 'bg-orange-100 text-orange-800 border-orange-200' :
                    'bg-blue-100 text-blue-800 border-blue-200'}`}>
                  {email.Urgency} Urgency
                </span>
              )}
            </div>
          </div>
          <div className="text-indigo-800">
            <p className="font-medium mb-1">{t('emailDetail.summary')}:</p>
            <p className="mb-4">{email.Summary}</p>

            {email.ActionItems && email.ActionItems.length > 0 && (
              <div className="mt-4">
                <p className="font-medium mb-1">{t('emailDetail.actionItems')}:</p>
                <ul className="list-disc list-inside space-y-1">
                  {email.ActionItems.map((item, index) => (
                    <li key={index}>{item}</li>
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
                    className="flex items-center gap-1.5 px-3 py-1.5 bg-white hover:bg-indigo-50 text-indigo-600 border border-indigo-200 rounded-lg text-xs font-medium transition-colors shadow-sm"
                  >
                    {action.type === 'calendar_event' ? '📅' : '✅'}
                    {getActionLabel(action)}
                  </button>
                ))}
              </div>
            )}
          </div>
          <div className="mt-4 text-right">
            <Button onClick={() => setIsModalOpen(true)}>{t('emailDetail.generateDraft')}</Button>
          </div>
        </div>
      )}

      {/* Email Header */}
      <div className="bg-white shadow rounded-lg p-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-4">{email.Subject}</h1>
        <div className="flex justify-between items-center text-sm text-gray-500 mb-8 border-b pb-4">
          <div>
            <span className="font-medium text-gray-900">{t('emailDetail.from')}</span> {email.Sender}
          </div>
          <div>{new Date(email.Date).toLocaleString()}</div>
        </div>

        {/* Email Body */}
        <div className="prose max-w-none text-gray-800 whitespace-pre-wrap">
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
