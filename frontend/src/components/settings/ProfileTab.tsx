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
      setFirstName(nameParts[0] || '');
      // Handle names with more than two parts correctly
      setLastName(nameParts.slice(1).join(' ') || '');
    }
  }, [user]); // This effect runs when the user object is loaded or changes

  const handleSaveChanges = async () => {
    try {
      const fullName = `${firstName} ${lastName}`.trim();
      await api.patch('/users/me', { name: fullName });
      // Update local store
      useAuthStore.setState(state => ({
        user: state.user ? { ...state.user, name: fullName } : null
      }));
      alert(t('settings.profile.success')); // Need to add key
    } catch (error) {
      console.error("Failed to update profile:", error);
      alert(t('settings.profile.error')); // Need to add key
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
        <button
          onClick={handleSaveChanges}
          className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium shadow-md transition-colors"
        >
          {t('settings.profile.saveChanges')}
        </button>
      </div>
    </div>
  );
}
