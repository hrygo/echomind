'use client';

import React, { useState, useEffect } from 'react';
import { RefreshCw, CheckCircle2, AlertCircle, Mail, Loader2 } from 'lucide-react';
import { useLanguage } from '@/lib/i18n/LanguageContext';
import { api } from '@/lib/api';
// import { cn } from '@/lib/utils'; // cn is not directly used in this file but by InputField in SmartMailboxForm
import { isAxiosError } from 'axios';
import { SmartMailboxForm } from '../onboarding/SmartMailboxForm'; // Re-use from onboarding
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '@/components/ui/Dialog'; // Assuming Radix UI Dialog
import { useToast } from '@/lib/hooks/useToast';
import { useConfirm } from '@/components/ui/ConfirmDialog';

interface AccountStatusResponse {
  has_account: boolean;
  email?: string;
  server_address?: string;
  server_port?: number;
  username?: string;
  is_connected?: boolean;
  last_sync_at?: string;
  error_message?: string;
}

export function ConnectionTab() {
  const { t } = useLanguage();
  const toast = useToast();
  const confirm = useConfirm();
  const [accountStatus, setAccountStatus] = useState<AccountStatusResponse | null>(null);
  const [isLoadingStatus, setIsLoadingStatus] = useState(true);
  const [isSyncing, setIsSyncing] = useState(false);
  const [showReconfigModal, setShowReconfigModal] = useState(false);

  const fetchAccountStatus = async () => {
    setIsLoadingStatus(true);
    try {
      const response = await api.get<AccountStatusResponse>('/settings/account');
      setAccountStatus(response.data);
    } catch (error) {
      console.error("Failed to fetch account status:", error);
      setAccountStatus({ has_account: false }); // Assume no account if error
    } finally {
      setIsLoadingStatus(false);
    }
  };

  useEffect(() => {
    fetchAccountStatus();
  }, []);

  const handleSyncNow = async () => {
    setIsSyncing(true);
    try {
      await api.post<{ message: string }>("/sync");
      toast.success(t('settings.connection.syncStarted'));
      fetchAccountStatus(); // Refresh status after sync
    } catch (error: unknown) {
      console.error("Sync failed:", error);
      if (isAxiosError(error) && error.response?.status === 400) {
        toast.error(error.response.data?.error || t('settings.connection.notConfigured'));
      } else {
        toast.error(t('settings.connection.syncFailed'));
      }
    } finally {
      setIsSyncing(false);
    }
  };

  return (
    <div className="space-y-8">
      {/* Mailbox Connection Section */}
      <div>
        <h3 className="text-xl font-bold text-slate-800">{t('settings.connection.title')}</h3>
        <p className="text-slate-700 text-sm mt-1">{t('settings.connectedMailboxDesc')}</p>
      </div>

      {isLoadingStatus ? (
        <div className="flex items-center gap-3 text-slate-500">
          <Loader2 className="w-5 h-5 animate-spin" />
          {t('common.loading')}
        </div>
      ) : accountStatus?.has_account ? (
        <div className="space-y-4 bg-slate-50 p-6 rounded-lg border border-slate-100">
          <div className="flex items-center gap-3">
            {accountStatus.is_connected ? (
              <CheckCircle2 className="w-6 h-6 text-green-600" />
            ) : (
              <AlertCircle className="w-6 h-6 text-red-600" />
            )}
            <div>
              <p className="text-lg font-semibold text-slate-800">
                {accountStatus.email} - {accountStatus.is_connected ? t('settings.connection.statusConnected') : t('settings.connection.statusDisconnected')}
              </p>
              {accountStatus.last_sync_at && (
                <p className="text-sm text-slate-600 mt-1">
                  {t('settings.connection.lastSynced').replace('{time}', new Date(accountStatus.last_sync_at).toLocaleString())}
                </p>
              )}
              {accountStatus.error_message && (
                <p className="text-sm text-red-500 mt-1">Error: {accountStatus.error_message}</p>
              )}
            </div>
          </div>

          <div className="flex gap-3">
            <button
              onClick={handleSyncNow}
              disabled={isSyncing}
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium shadow-md transition-colors disabled:opacity-60 disabled:cursor-not-allowed"
            >
              {isSyncing ? <Loader2 className="w-4 h-4 mr-2 animate-spin" /> : <RefreshCw className="w-4 h-4 mr-2" />}
              {isSyncing ? t('settings.connection.syncing') : t('settings.connection.syncNow')}
            </button>
            <button
              onClick={() => setShowReconfigModal(true)}
              className="px-4 py-2 bg-white border border-slate-200 rounded-lg text-sm font-medium text-slate-700 hover:bg-slate-50 transition-colors"
            >
              {t('settings.connection.reconfigure')}
            </button>
            <button
              onClick={() => {
                confirm(
                  t('settings.connection.disconnectConfirm'),
                  async () => {
                    try {
                      await api.delete('/settings/account');
                      toast.success(t('settings.connection.disconnectSuccess'));
                      fetchAccountStatus();
                    } catch (error) {
                      console.error("Failed to disconnect:", error);
                      toast.error(t('settings.connection.disconnectFailed'));
                    }
                  },
                  {
                    title: t('settings.connection.disconnectTitle'),
                    confirmText: t('common.confirm'),
                    cancelText: t('common.cancel')
                  }
                );
              }}
              className="px-4 py-2 bg-white border border-red-200 text-red-600 rounded-lg text-sm font-medium hover:bg-red-50 transition-colors"
            >
              {t('settings.connection.disconnect')}
            </button>
          </div>
        </div>
      ) : (
        <div className="text-center py-8 bg-slate-50 rounded-lg border border-dashed border-slate-200 text-slate-500">
          <Mail className="w-12 h-12 mx-auto mb-4" />
          <p className="text-lg font-semibold mb-2">{t('settings.connection.notConfigured')}</p>
          <p className="text-sm mb-4">{t('settings.connection.configureMailboxDesc')}</p>
          <button
            onClick={() => setShowReconfigModal(true)}
            className="px-6 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-base font-semibold shadow-md transition-colors"
          >
            {t('settings.connection.configureMailbox')}
          </button>
        </div>
      )}

      <Dialog open={showReconfigModal} onOpenChange={setShowReconfigModal}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>{t('settings.connection.configureMailbox')}</DialogTitle>
            <DialogDescription>
              {t('settings.connection.mailboxConfigDesc')}
            </DialogDescription>
          </DialogHeader>
          <SmartMailboxForm /> {/* Re-use the smart form from onboarding */}
        </DialogContent>
      </Dialog>
    </div>
  );
}
