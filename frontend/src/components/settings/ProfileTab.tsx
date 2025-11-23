'use client';

import React from 'react';
import { Camera } from 'lucide-react';
import { useLanguage } from '@/lib/i18n/LanguageContext';
import { useAuthStore } from '@/store/auth';

export function ProfileTab() {
  const { t } = useLanguage();
  const user = useAuthStore(state => state.user);

  // Dummy state for profile fields (in a real app, these would be managed with more robust state/forms)
  const [firstName, setFirstName] = React.useState(user?.name?.split(' ')[0] || '');
  const [lastName, setLastName] = React.useState(user?.name?.split(' ')[1] || '');

  const handleSaveChanges = () => {
    // TODO: Implement API call to update user profile
    alert('Saving changes...'); // Placeholder
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
          {user?.name ? user.name[0].toUpperCase() : user?.email[0].toUpperCase()}
          <div className="absolute inset-0 bg-black/50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
            <Camera className="w-6 h-6 text-white" />
          </div>
        </div>
        <button 
          onClick={handleSaveChanges} // Placeholder for avatar upload
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
