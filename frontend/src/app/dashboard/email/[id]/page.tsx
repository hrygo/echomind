"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import apiClient from "@/lib/api";

interface EmailDetail {
  ID: string; // Changed from number to string (UUID)
  Subject: string;
  Sender: string;
  Date: string;
  BodyText: string;
  Summary: string;
  Sentiment: string;
  Urgency: string;
}

export default function EmailDetailPage() {
  const params = useParams();
  const { id } = params;
  const [email, setEmail] = useState<EmailDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;

    async function fetchEmail() {
      try {
        // Ensure `id` is treated as a string for UUID
        const response = await apiClient.get<EmailDetail>(`/emails/${id}`);
        setEmail(response.data);
      } catch (err: any) {
        console.error("Error fetching email:", err);
        setError(err.message || "Failed to load email.");
      } finally {
        setLoading(false);
      }
    }

    fetchEmail();
  }, [id]);

  if (loading) return <div className="p-8">Loading...</div>;
  if (error) return <div className="p-8 text-red-500">Error: {error}</div>;
  if (!email) return <div className="p-8">Email not found</div>;

  return (
    <div className="max-w-4xl mx-auto">
      <div className="mb-6">
        <Link href="/dashboard" className="text-blue-600 hover:underline">&larr; Back to Inbox</Link>
      </div>

      {/* AI Insight Card */}
      {email.Summary && (
        <div className="bg-gradient-to-r from-indigo-50 to-blue-50 border border-indigo-100 rounded-lg p-6 mb-8 shadow-sm">
          <div className="flex justify-between items-start mb-4">
            <h2 className="text-lg font-semibold text-indigo-900 flex items-center gap-2">
              ✨ AI Insights
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
            <p className="font-medium mb-1">Summary:</p>
            <p>{email.Summary}</p>
          </div>
        </div>
      )}

      {/* Email Header */}
      <div className="bg-white shadow rounded-lg p-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-4">{email.Subject}</h1>
        <div className="flex justify-between items-center text-sm text-gray-500 mb-8 border-b pb-4">
          <div>
            <span className="font-medium text-gray-900">From:</span> {email.Sender}
          </div>
          <div>{new Date(email.Date).toLocaleString()}</div>
        </div>

        {/* Email Body */}
        <div className="prose max-w-none text-gray-800 whitespace-pre-wrap">
          {email.BodyText}
        </div>
      </div>
    </div>
  );
}
