"use client";

import { useState } from "react";
import { User, Bell, Shield, CreditCard, Palette, Mail } from "lucide-react";
import { useLanguage } from "@/lib/i18n/LanguageContext";
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'; // Assuming Radix UI Tabs
import { ProfileTab } from '@/components/settings/ProfileTab';
import { ConnectionTab } from '@/components/settings/ConnectionTab';
import { AppearanceTab } from '@/components/settings/AppearanceTab';
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
    <div className="flex h-full bg-card rounded-2xl shadow-sm border border-border overflow-hidden">
      {/* Settings Sidebar / Tabs List */}
      <Tabs value={activeTab} onValueChange={setActiveTab} orientation="vertical" className="flex w-full">
        <TabsList className="w-64 bg-muted border-r border-border flex flex-col h-full justify-start p-4 space-y-1">
          <div className="px-4 pb-3 border-b border-border">
            <h2 className="text-lg font-bold text-foreground">{t('settings.title')}</h2>
          </div>
          {tabs.map((tab) => (
            <TabsTrigger
              key={tab.id}
              value={tab.id}
              className={cn(
                "w-full flex items-center justify-start gap-3 px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200",
                activeTab === tab.id
                  ? "bg-card text-primary shadow-sm"
                  : "text-muted-foreground hover:bg-accent hover:text-foreground"
              )}
            >
              <tab.icon className={cn("w-5 h-5 ", activeTab === tab.id ? "text-primary" : "text-muted-foreground")} />
              {tab.label}
            </TabsTrigger>
          ))}
        </TabsList>

        {/* Settings Content / Tabs Content */}
        <div className="flex-1 overflow-y-auto p-8 md:p-10">
          <div className="max-w-5xl mx-auto">
            <TabsContent value="profile">
              <ProfileTab />
            </TabsContent>
            <TabsContent value="connection">
              <ConnectionTab />
            </TabsContent>
            <TabsContent value="appearance">
              <AppearanceTab />
            </TabsContent>
            {/* Placeholder tabs for future development */}
            <TabsContent value="notifications">
              <div className="flex flex-col items-center justify-center min-h-[400px] text-foreground animate-in fade-in duration-300">
                <div className="p-4 bg-muted rounded-full mb-4">
                  <Bell className="w-8 h-8 text-muted-foreground" />
                </div>
                <p className="text-lg font-medium mb-2">{t('settings.notifications')}</p>
                <p className="text-sm text-muted-foreground">{t('common.developing')}</p>
              </div>
            </TabsContent>
            <TabsContent value="security">
              <div className="flex flex-col items-center justify-center min-h-[400px] text-foreground animate-in fade-in duration-300">
                <div className="p-4 bg-muted rounded-full mb-4">
                  <Shield className="w-8 h-8 text-muted-foreground" />
                </div>
                <p className="text-lg font-medium mb-2">{t('settings.security')}</p>
                <p className="text-sm text-muted-foreground">{t('common.developing')}</p>
              </div>
            </TabsContent>
            <TabsContent value="billing">
              <div className="flex flex-col items-center justify-center min-h-[400px] text-foreground animate-in fade-in duration-300">
                <div className="p-4 bg-muted rounded-full mb-4">
                  <CreditCard className="w-8 h-8 text-muted-foreground" />
                </div>
                <p className="text-lg font-medium mb-2">{t('settings.billing')}</p>
                <p className="text-sm text-muted-foreground">{t('common.developing')}</p>
              </div>
            </TabsContent>
          </div>
        </div>
      </Tabs>
    </div>
  );
}
