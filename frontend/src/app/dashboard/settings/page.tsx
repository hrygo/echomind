"use client";

import { useState } from "react";
import { User, Bell, Shield, CreditCard, Palette, Mail } from "lucide-react";
import { useLanguage } from "@/lib/i18n/LanguageContext";
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/Tabs'; // Assuming Radix UI Tabs
import { ProfileTab } from '@/components/settings/ProfileTab';
import { ConnectionTab } from '@/components/settings/ConnectionTab';
import { cn } from '@/lib/utils'; // Assuming cn utility exists

export default function SettingsPage() {
  const [activeTab, setActiveTab] = useState("profile"); // Default to profile tab
  const { t } = useLanguage();

  const tabs = [
    { id: "profile", label: t('settings.profile.title'), icon: User },
    { id: "connection", label: t('settings.connection.title'), icon: Mail },
    { id: "notifications", label: t('settings.notifications'), icon: Bell },
    { id: "security", label: t('settings.security'), icon: Shield },
    { id: "billing", label: t('settings.billing'), icon: CreditCard },
    { id: "appearance", label: t('settings.appearance'), icon: Palette },
  ];

  return (
    <div className="flex h-full bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden">
      {/* Settings Sidebar / Tabs List */}
      <Tabs value={activeTab} onValueChange={setActiveTab} orientation="vertical" className="flex w-full">
        <TabsList className="w-64 bg-slate-50 border-r border-slate-100 flex flex-col h-full justify-start p-3 space-y-1">
          <div className="p-3 border-b border-slate-100 mb-2">
            <h2 className="text-lg font-bold text-slate-800">{t('settings.title')}</h2>
          </div>
          {tabs.map((tab) => (
            <TabsTrigger
              key={tab.id}
              value={tab.id}
              className={cn(
                "w-full flex items-center justify-start gap-3 px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200",
                activeTab === tab.id
                  ? "bg-white text-blue-600 shadow-sm"
                  : "text-slate-700 hover:bg-slate-100 hover:text-slate-900"
              )}
            >
              <tab.icon className={cn("w-5 h-5 ", activeTab === tab.id ? "text-blue-600" : "text-slate-600")} />
              {tab.label}
            </TabsTrigger>
          ))}
        </TabsList>

        {/* Settings Content / Tabs Content */}
        <div className="flex-1 overflow-y-auto p-8">
          <div className="max-w-2xl mx-auto">
            <TabsContent value="profile">
              <ProfileTab />
            </TabsContent>
            <TabsContent value="connection">
              <ConnectionTab />
            </TabsContent>
            {/* Placeholder tabs for future development */}
            <TabsContent value="notifications">
              <div className="flex flex-col items-center justify-center h-[300px] text-slate-700 animate-in fade-in duration-300">
                <div className="p-4 bg-slate-100 rounded-full mb-4">
                  <Bell className="w-8 h-8 text-slate-600" />
                </div>
                <p>{t('settings.notifications')} {t('common.developing')}</p> {/* common.developing added to dict */}
              </div>
            </TabsContent>
            <TabsContent value="security">
              <div className="flex flex-col items-center justify-center h-[300px] text-slate-700 animate-in fade-in duration-300">
                <div className="p-4 bg-slate-100 rounded-full mb-4">
                  <Shield className="w-8 h-8 text-slate-600" />
                </div>
                <p>{t('settings.security')} {t('common.developing')}</p>
              </div>
            </TabsContent>
            <TabsContent value="billing">
              <div className="flex flex-col items-center justify-center h-[300px] text-slate-700 animate-in fade-in duration-300">
                <div className="p-4 bg-slate-100 rounded-full mb-4">
                  <CreditCard className="w-8 h-8 text-slate-600" />
                </div>
                <p>{t('settings.billing')} {t('common.developing')}</p>
              </div>
            </TabsContent>
            <TabsContent value="appearance">
              <div className="flex flex-col items-center justify-center h-[300px] text-slate-700 animate-in fade-in duration-300">
                <div className="p-4 bg-slate-100 rounded-full mb-4">
                  <Palette className="w-8 h-8 text-slate-600" />
                </div>
                <p>{t('settings.appearance')} {t('common.developing')}</p>
              </div>
            </TabsContent>
          </div>
        </div>
      </Tabs>
    </div>
  );
}
