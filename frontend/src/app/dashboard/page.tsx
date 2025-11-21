"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
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
}

export default function DashboardPage() {
  const [emails, setEmails] = useState<Email[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchEmails() {
      try {
        const response = await apiClient.get<Email[]>("/emails");
        setEmails(response.data);
      } catch (err: any) {
        console.error("Error fetching emails:", err);
        setError(err.message || "Failed to load emails.");
      } finally {
        setLoading(false);
      }
    }

    fetchEmails();
  }, []);

  const handleSync = async () => {
    try {
      const response = await apiClient.post<{ message: string }>("/sync");
      alert(response.data.message);
      // Optionally refetch or poll if sync is quick, or rely on worker for updates
    } catch (err: any) {
      console.error("Sync error:", err);
      alert(`Sync failed: ${err.message}`);
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
                    {email.Urgency === "High" && (
                      <span className="bg-red-100 text-red-800 text-xs font-medium px-2.5 py-0.5 rounded border border-red-200 shrink-0">
                        High Priority
                      </span>
                    )}
                    {email.Sentiment === "Negative" && (
                      <span className="bg-orange-100 text-orange-800 text-xs font-medium px-2.5 py-0.5 rounded border border-orange-200 shrink-0">
                        Negative
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