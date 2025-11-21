"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useSearchParams, useRouter } from 'next/navigation';
import apiClient from "@/lib/api";

// Define Email interface based on API response
interface Email {
  ID: string; // Changed from number to string (UUID)
  Subject: string;
  Sender: string;
  Date: string;
  Summary: string;
  Sentiment: string;
  Urgency: string;
  Category: string; // New field from AI analysis
  ActionItems: string[]; // New field from AI analysis
}

export default function DashboardPage() {
  const [emails, setEmails] = useState<Email[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const searchParams = useSearchParams();
  const router = useRouter();
  const categoryFilter = searchParams.get('category');

  useEffect(() => {
    async function fetchEmails() {
      try {
        setLoading(true);
        setError(null);
        const url = categoryFilter ? `/emails?category=${categoryFilter}` : "/emails";
        const response = await apiClient.get<Email[]>(url);
        setEmails(response.data);
      } catch (err: unknown) {
        console.error("Error fetching emails:", err);
        if (err instanceof Error) {
          setError(err.message);
        } else {
          setError("Failed to load emails.");
        }
      } finally {
        setLoading(false);
      }
    }

    fetchEmails();
  }, [categoryFilter]); // Re-fetch when categoryFilter changes

  const handleSync = async () => {
    try {
      const response = await apiClient.post<{ message: string }>("/sync");
      alert(response.data.message);
      // Optionally refetch or poll if sync is quick, or rely on worker for updates
    } catch (err: any) {
      console.error("Sync error:", err);
      if (err.response && err.response.status === 400) {
        if (confirm(err.response.data.error + "\n\nGo to Settings now?")) {
            router.push('/dashboard/settings');
        }
      } else if (err instanceof Error) {
        alert(`Sync failed: ${err.message}`);
      } else {
        alert("Sync failed: An unknown error occurred.");
      }
    }
  };

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-semibold text-gray-900">Inbox</h1>
        <button
          onClick={handleSync}
          className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 text-sm font-medium"
        >
          Sync Now
        </button>
      </div>

      <div className="bg-white shadow rounded-lg overflow-hidden">
        {loading ? (
          <div className="p-8 text-center text-gray-500">Loading emails...</div>
        ) : error ? (
          <div className="p-8 text-center text-red-500">Error: {error}</div>
        ) : emails.length === 0 ? (
          <div className="p-8 text-center text-gray-500">No emails found. Try syncing!</div>
        ) : (
          <ul className="divide-y divide-gray-200">
            {emails.map((email) => (
              <li key={email.ID} className="hover:bg-gray-50 transition">
                <Link href={`/dashboard/email/${email.ID}`} className="block p-6">
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm font-medium text-gray-900">{email.Sender}</span>
                    <span className="text-xs text-gray-500">{new Date(email.Date).toLocaleString()}</span>
                  </div>
                  <h3 className="text-lg font-medium text-gray-900 mb-1">{email.Subject}</h3>
                  <div className="flex items-center gap-2 mt-2">
                    <p className="text-sm text-gray-600 line-clamp-2 flex-1">{email.Summary || "Analysis pending..."}</p>
                    {email.Sentiment && (
                      <span className={`text-xs font-medium px-2.5 py-0.5 rounded border shrink-0 ${email.Sentiment === 'Positive' ? 'bg-green-100 text-green-800 border-green-200' : email.Sentiment === 'Negative' ? 'bg-red-100 text-red-800 border-red-200' : 'bg-gray-100 text-gray-800 border-gray-200'}`}>
                        {email.Sentiment}
                      </span>
                    )}
                    {email.Urgency && (
                      <span className={`text-xs font-medium px-2.5 py-0.5 rounded border shrink-0 ${email.Urgency === 'High' ? 'bg-orange-100 text-orange-800 border-orange-200' : email.Urgency === 'Medium' ? 'bg-yellow-100 text-yellow-800 border-yellow-200' : 'bg-blue-100 text-blue-800 border-blue-200'}`}>
                        {email.Urgency} Urgency
                      </span>
                    )}
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