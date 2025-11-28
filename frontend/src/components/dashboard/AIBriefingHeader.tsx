import React from 'react';
import { Sparkles, Loader2 } from 'lucide-react';
import { useLanguage } from "@/lib/i18n/LanguageContext";
import { useUserProfile } from '@/hooks/useUserProfile';
import { useAuthStore } from '@/store/auth';

type ViewType = "executive" | "manager" | "dealmaker";

interface AIBriefingHeaderProps {
  currentView: ViewType;
  userName?: string;
}

export function AIBriefingHeader({ currentView, userName }: AIBriefingHeaderProps) {
  const { t } = useLanguage();
  const { data: userProfile, isLoading: profileLoading } = useUserProfile();
  const currentUser = useAuthStore(state => state.user);

  // Prioritize real user data from auth store, fallback to API profile, then provided userName
  const actualUserName = currentUser?.name || userProfile?.name || userName || "演示用户";

  // Mock data logic - in real app this comes from API
  const hour = new Date().getHours();
  const timeOfDay = hour < 12 ? 'greetingMorning' : hour < 18 ? 'greetingAfternoon' : 'greetingEvening';

  const getSummary = () => {
    switch (currentView) {
      case 'executive':
        return t('dashboard.executiveSummary')
          .replace('{riskCount}', '3')
          .replace('{decisionCount}', '5');
      case 'manager':
        return t('dashboard.managerSummary')
          .replace('{taskCount}', '12')
          .replace('{overdueCount}', '2');
      case 'dealmaker':
        return t('dashboard.dealmakerSummary')
          .replace('{opportunityCount}', '4')
          .replace('{followUpCount}', '7');
      default:
        return "";
    }
  };

  const showLoader = profileLoading && !currentUser?.name;

  return (
    <div className="relative overflow-hidden rounded-2xl bg-gradient-to-br from-blue-600 to-indigo-700 text-white p-8 shadow-lg mb-8">
      {/* Background Decorative Elements */}
      <div className="absolute top-0 right-0 -mt-10 -mr-10 w-40 h-40 bg-white/10 rounded-full blur-3xl"></div>
      <div className="absolute bottom-0 left-0 -mb-10 -ml-10 w-40 h-40 bg-blue-400/20 rounded-full blur-3xl"></div>

      <div className="relative z-10">
        <div className="flex items-center gap-2 mb-3 opacity-90">
          <Sparkles className="w-4 h-4 text-yellow-300" />
          <span className="text-xs font-bold uppercase tracking-wider">{t('dashboard.aiBriefing')}</span>
        </div>
        <h1 className="text-3xl md:text-4xl font-bold tracking-tight mb-3">
          {showLoader ? (
            <span className="flex items-center gap-2">
              {t(`dashboard.${timeOfDay}`)}, <Loader2 className="w-6 h-6 animate-spin" />
            </span>
          ) : (
            `${t(`dashboard.${timeOfDay}`)}, ${actualUserName}.`
          )}
        </h1>
        <p className="text-blue-50 text-lg md:text-xl leading-relaxed max-w-3xl font-medium">
          {getSummary()}
        </p>
      </div>
    </div>
  );
}
