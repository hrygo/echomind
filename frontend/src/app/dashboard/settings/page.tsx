"use client";

import { useState, useEffect } from "react";
import { useRouter } from 'next/navigation';
import apiClient from "@/lib/api";
import { Input } from "@/components/ui/Input";
import { Button } from "@/components/ui/Button";
import { Label } from "@/components/ui/Label";
import { Card } from "@/components/ui/Card";

interface EmailAccountStatus {
  has_account: boolean;
  email?: string;
  server_address?: string;
  server_port?: number;
  username?: string;
  is_connected?: boolean;
  last_sync_at?: string;
  error_message?: string;
}

export default function SettingsPage() {
  const [accountStatus, setAccountStatus] = useState<EmailAccountStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [formLoading, setFormLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [formData, setFormData] = useState({
    email: "",
    server_address: "",
    server_port: 993,
    username: "",
    password: "",
  });
  const router = useRouter();

  useEffect(() => {
    fetchAccountStatus();
  }, []);

  const fetchAccountStatus = async () => {
    try {
      setLoading(true);
      const response = await apiClient.get<EmailAccountStatus>("/settings/account");
      setAccountStatus(response.data);
      if (response.data.has_account && response.data.email) {
        setFormData(prev => ({
            ...prev,
            email: response.data.email || "",
            server_address: response.data.server_address || "",
            server_port: response.data.server_port || 993,
            username: response.data.username || "",
            password: "", // Never pre-fill password
        }));
      }
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || "Failed to fetch account status");
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: name === "server_port" ? parseInt(value) || 0 : value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormLoading(true);
    setError(null);
    try {
      await apiClient.post("/settings/account", formData);
      alert("Account connected successfully!");
      router.refresh(); // Refresh page data
      fetchAccountStatus(); // Re-fetch status to update UI
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || "Failed to connect account.");
    } finally {
      setFormLoading(false);
    }
  };

  if (loading) {
    return <div className="p-8 text-center">Loading settings...</div>;
  }

  return (
    <div className="max-w-2xl mx-auto py-8">
      <h1 className="text-3xl font-bold mb-8 text-gray-900">Settings</h1>

      <Card className="p-6">
        <h2 className="text-2xl font-semibold mb-6 text-gray-800">Email Account Connection</h2>

        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-6" role="alert">
            <strong className="font-bold">Error!</strong>
            <span className="block sm:inline"> {error}</span>
          </div>
        )}

        {accountStatus?.has_account ? (
          <div>
            <p className="mb-4 text-gray-700">Your email account is currently connected.</p>
            <div className="space-y-3 mb-6">
              <p><strong>Email:</strong> {accountStatus.email}</p>
              <p><strong>Server:</strong> {accountStatus.server_address}:{accountStatus.server_port}</p>
              <p><strong>Username:</strong> {accountStatus.username}</p>
              <p><strong>Status:</strong> 
                <span className={`font-medium ${accountStatus.is_connected ? 'text-green-600' : 'text-red-600'}`}>
                  {accountStatus.is_connected ? 'Connected' : 'Disconnected'}
                </span>
              </p>
              {accountStatus.last_sync_at && (
                <p><strong>Last Synced:</strong> {new Date(accountStatus.last_sync_at).toLocaleString()}</p>
              )}
              {accountStatus.error_message && (
                <p className="text-red-500"><strong>Last Error:</strong> {accountStatus.error_message}</p>
              )}
            </div>
            {/* Future: Add Disconnect Button */}
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <Label htmlFor="email">Email Address</Label>
              <Input
                id="email"
                name="email"
                type="email"
                value={formData.email}
                onChange={handleChange}
                required
                placeholder="your.email@example.com"
              />
            </div>
            <div>
              <Label htmlFor="server_address">IMAP Server Address</Label>
              <Input
                id="server_address"
                name="server_address"
                type="text"
                value={formData.server_address}
                onChange={handleChange}
                required
                placeholder="imap.example.com"
              />
            </div>
            <div>
              <Label htmlFor="server_port">IMAP Server Port</Label>
              <Input
                id="server_port"
                name="server_port"
                type="number"
                value={formData.server_port}
                onChange={handleChange}
                required
                placeholder="993"
              />
            </div>
            <div>
              <Label htmlFor="username">IMAP Username</Label>
              <Input
                id="username"
                name="username"
                type="text"
                value={formData.username}
                onChange={handleChange}
                required
                placeholder="your.username or email"
              />
            </div>
            <div>
              <Label htmlFor="password">IMAP App Password</Label>
              <Input
                id="password"
                name="password"
                type="password"
                value={formData.password}
                onChange={handleChange}
                required
                placeholder="Your app-specific password"
              />
              <p className="text-sm text-gray-500 mt-1">
                For Gmail/Outlook, this is often an app-specific password, not your main account password.
              </p>
            </div>
            <Button type="submit" className="w-full" disabled={formLoading}>
              {formLoading ? "Connecting..." : "Connect Account"}
            </Button>
          </form>
        )}
      </Card>
    </div>
  );
}