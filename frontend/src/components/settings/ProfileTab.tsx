'use client';

import React, { useEffect } from 'react';
import { Camera } from 'lucide-react';
import { useLanguage } from '@/lib/i18n/LanguageContext';
import { useAuthStore } from '@/store/auth';
import { api } from '@/lib/api';

export function ProfileTab() {
  const { t } = useLanguage();
  const user = useAuthStore(state => state.user);

  const [firstName, setFirstName] = React.useState('');
  const [lastName, setLastName] = React.useState('');

  useEffect(() => {
    if (user?.name) {
      const nameParts = user.name.split(' ');
      const newFirstName = nameParts[0] || '';
      const newLastName = nameParts.slice(1).join(' ') || '';

      if (firstName !== newFirstName) {
        setFirstName(newFirstName);
      }
      if (lastName !== newLastName) {
        setLastName(newLastName);
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [user?.name]); // Depend on user.name string to avoid object reference issues

  const updateUser = useAuthStore(state => state.updateUser);
  const [isLoading, setIsLoading] = React.useState(false);
  const [message, setMessage] = React.useState<{ type: 'success' | 'error'; text: string } | null>(null);

  // Auto-hide success messages after 3 seconds
  React.useEffect(() => {
    if (message?.type === 'success') {
      const timer = setTimeout(() => {
        setMessage(null);
      }, 3000);
      return () => clearTimeout(timer);
    }
  }, [message]);

  const handleSaveChanges = async () => {
    setIsLoading(true);
    setMessage(null);

    try {
      const fullName = `${firstName} ${lastName}`.trim();

      // API call to update profile
      await api.patch('/users/me', { name: fullName });

      // Update local store using the proper store method
      if (user) {
        updateUser({ name: fullName });
      }

      setMessage({ type: 'success', text: t('settings.profile.success') });
    } catch (error) {
      console.error("Failed to update profile:", error);
      setMessage({ type: 'error', text: t('settings.profile.error') });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="space-y-8">
      {/* Profile Section */}
      <div>
        <h3 className="text-xl font-bold text-slate-800">{t('settings.profile.title')}</h3>
        <p className="text-slate-700 text-sm mt-1">{t('settings.accountDesc')}</p>
      </div>

      <div className="flex items-center gap-6">
        <div className="w-20 h-20 rounded-full bg-slate-200 flex items-center justify-center text-2xl font-bold text-slate-700 relative group cursor-pointer overflow-hidden">
          {user?.name ? user.name[0].toUpperCase() : user?.email?.[0].toUpperCase()}
          <div className="absolute inset-0 bg-black/50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
            <Camera className="w-6 h-6 text-white" />
          </div>
        </div>
        <button
          onClick={() => alert("Avatar upload not implemented yet")} // Placeholder for avatar upload
          className="px-4 py-2 bg-white border border-slate-200 rounded-lg text-sm font-medium text-slate-700 hover:bg-slate-50 transition-colors"
        >
          {t('settings.profile.uploadAvatar')}
        </button>
      </div>

      <div className="grid grid-cols-2 gap-6">
        <div className="space-y-2">
          <label className="text-sm font-medium text-slate-700">{t('settings.firstName')}</label>
          <input
            type="text"
            className="w-full px-3 py-2 border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
            value={firstName}
            onChange={(e) => setFirstName(e.target.value)}
          />
        </div>
        <div className="space-y-2">
          <label className="text-sm font-medium text-slate-700">{t('settings.lastName')}</label>
          <input
            type="text"
            className="w-full px-3 py-2 border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
            value={lastName}
            onChange={(e) => setLastName(e.target.value)}
          />
        </div>
        <div className="col-span-2 space-y-2">
          <label className="text-sm font-medium text-slate-700">{t('settings.loginEmail')}</label>
          <input
            type="email"
            className="w-full px-3 py-2 border border-slate-200 rounded-lg bg-slate-50 text-slate-600 cursor-not-allowed"
            value={user?.email || ''}
            disabled
          />
          <p className="text-sm text-slate-700">{t('settings.loginEmailDesc')}</p>
        </div>
      </div>

      <div className="flex justify-end">
        <div className="space-y-3">
          {message && (
            <div
              className={`px-4 py-2 rounded-lg text-sm ${
                message.type === 'success'
                  ? 'bg-green-50 text-green-700 border border-green-200'
                  : 'bg-red-50 text-red-700 border border-red-200'
              }`}
            >
              {message.text}
            </div>
          )}
          <button
            onClick={handleSaveChanges}
            disabled={isLoading}
            className={`px-4 py-2 rounded-lg text-sm font-medium shadow-md transition-colors ${
              isLoading
                ? 'bg-gray-400 text-gray-200 cursor-not-allowed'
                : 'bg-blue-600 hover:bg-blue-700 text-white'
            }`}
          >
            {isLoading ? (
              <span className="flex items-center gap-2">
                <svg className="animate-spin h-4 w-4" viewBox="0 0 24 24">
                  <circle
                    className="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    strokeWidth="4"
                    fill="none"
                  />
                  <path
                    className="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  />
                </svg>
                {t('common.loading')}
              </span>
            ) : (
              t('settings.profile.saveChanges')
            )}
          </button>
        </div>
      </div>
    </div>
  );
}
