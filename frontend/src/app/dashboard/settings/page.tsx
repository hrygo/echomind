"use client";

import { useState, useEffect } from "react";
import { User, Bell, Shield, CreditCard, Palette, Camera, RefreshCw, Mail, Eye, EyeOff } from "lucide-react";
import { useLanguage } from "@/lib/i18n/LanguageContext";
import apiClient from "@/lib/api";

export default function SettingsPage() {
  const [activeTab, setActiveTab] = useState("account");
  const { t } = useLanguage();
  const [isSyncing, setIsSyncing] = useState(false);
  const [lastSynced, setLastSynced] = useState<string | null>(null);
  const [showPassword, setShowPassword] = useState(false);

  // Mailbox configuration state
  const [mailboxConfig, setMailboxConfig] = useState({
    email: "user@example.com",
    password: "",
    imapServer: "imap.gmail.com",
    imapPort: "993",
    smtpServer: "smtp.gmail.com",
    smtpPort: "587"
  });

  useEffect(() => {
    // Fetch account status to get last synced time (mock implementation for now)
    // In a real app, this would be an API call like:
    // apiClient.get('/settings/account').then(res => setLastSynced(res.data.last_synced));
    setLastSynced("2 hours ago");
  }, []);

  const handleSync = async () => {
    setIsSyncing(true);
    try {
      await apiClient.post("/sync");
      setLastSynced("Just now");
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (error: any) {
      console.error("Sync failed:", error);
      if (error.response?.status === 400) {
        alert(error.response.data.error || "请先配置邮箱账户");
      } else {
        alert("同步失败，请稍后重试");
      }
    } finally {
      setIsSyncing(false);
    }
  };

  const handleSaveMailboxConfig = async () => {
    try {
      await apiClient.post('/settings/account', {
        email: mailboxConfig.email,
        server_address: mailboxConfig.imapServer,
        server_port: parseInt(mailboxConfig.imapPort),
        username: mailboxConfig.email, // Assuming username is email
        password: mailboxConfig.password
      });
      // console.log("Saving mailbox config:", mailboxConfig);
      alert("配置已保存");
    } catch (error) {
      console.error("Failed to save config:", error);
      alert("保存失败，请检查配置");
    }
  };

  const tabs = [
    { id: "account", label: t('settings.account'), icon: User },
    { id: "notifications", label: t('settings.notifications'), icon: Bell },
    { id: "security", label: t('settings.security'), icon: Shield },
    { id: "billing", label: t('settings.billing'), icon: CreditCard },
    { id: "appearance", label: t('settings.appearance'), icon: Palette },
  ];

  return (
    <div className="flex h-full bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden">
      {/* Settings Sidebar */}
      <div className="w-64 bg-slate-50 border-r border-slate-100 flex flex-col">
        <div className="p-6 border-b border-slate-100">
          <h2 className="text-lg font-bold text-slate-800">{t('settings.title')}</h2>
        </div>
        <nav className="flex-1 p-3 space-y-1">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`w-full flex items-center gap-3 px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200 ${activeTab === tab.id
                  ? "bg-white text-blue-600 shadow-sm"
                  : "text-slate-600 hover:bg-slate-100 hover:text-slate-900"
                }`}
            >
              <tab.icon className={`w-5 h-5 ${activeTab === tab.id ? "text-blue-600" : "text-slate-400"}`} />
              {tab.label}
            </button>
          ))}
        </nav>
      </div>

      {/* Settings Content */}
      <div className="flex-1 overflow-y-auto p-8">
        <div className="max-w-2xl mx-auto">
          {activeTab === "account" && (
            <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
              {/* Profile Section */}
              <div>
                <h3 className="text-xl font-bold text-slate-800">{t('settings.accountInfo')}</h3>
                <p className="text-slate-500 text-sm mt-1">{t('settings.accountDesc')}</p>
              </div>

              <div className="flex items-center gap-6">
                <div className="w-20 h-20 rounded-full bg-slate-200 flex items-center justify-center text-2xl font-bold text-slate-500 relative group cursor-pointer overflow-hidden">
                  U
                  <div className="absolute inset-0 bg-black/50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                    <Camera className="w-6 h-6 text-white" />
                  </div>
                </div>
                <button className="px-4 py-2 bg-white border border-slate-200 rounded-lg text-sm font-medium text-slate-700 hover:bg-slate-50 transition-colors">
                  {t('settings.changeAvatar')}
                </button>
              </div>

              <div className="grid grid-cols-2 gap-6">
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-700">{t('settings.firstName')}</label>
                  <input type="text" className="w-full px-3 py-2 border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all" defaultValue="User" />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-700">{t('settings.lastName')}</label>
                  <input type="text" className="w-full px-3 py-2 border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all" defaultValue="Name" />
                </div>
                <div className="col-span-2 space-y-2">
                  <label className="text-sm font-medium text-slate-700">{t('settings.loginEmail')}</label>
                  <input
                    type="email"
                    className="w-full px-3 py-2 border border-slate-200 rounded-lg bg-slate-50 text-slate-600 cursor-not-allowed"
                    defaultValue="user@example.com"
                    disabled
                  />
                  <p className="text-xs text-slate-400">{t('settings.loginEmailDesc')}</p>
                </div>
              </div>

              <hr className="border-slate-100" />

              {/* Connected Mailbox Section */}
              <div>
                <h3 className="text-xl font-bold text-slate-800 flex items-center gap-2">
                  <Mail className="w-5 h-5 text-blue-600" />
                  {t('settings.connectedMailbox')}
                </h3>
                <p className="text-slate-500 text-sm mt-1">{t('settings.connectedMailboxDesc')}</p>
              </div>

              <div className="bg-slate-50 rounded-xl p-6 border border-slate-100 space-y-6">
                {/* Sync Status */}
                <div className="flex items-center justify-between pb-4 border-b border-slate-200">
                  <div>
                    <p className="text-sm font-medium text-slate-500">{t('settings.syncStatus')}</p>
                    <div className="flex items-center gap-2 mt-1">
                      <div className="w-2 h-2 rounded-full bg-green-500"></div>
                      <span className="text-slate-800 font-medium">Active</span>
                      <span className="text-slate-400 text-sm ml-2">
                        {t('settings.lastSynced')}: {lastSynced}
                      </span>
                    </div>
                  </div>
                  <button
                    onClick={handleSync}
                    disabled={isSyncing}
                    className="flex items-center gap-2 px-4 py-2 bg-white border border-slate-200 rounded-lg text-sm font-medium text-slate-700 hover:bg-slate-50 hover:text-blue-600 transition-colors disabled:opacity-50"
                  >
                    <RefreshCw className={`w-4 h-4 ${isSyncing ? 'animate-spin' : ''}`} />
                    {isSyncing ? t('settings.syncing') : t('settings.syncNow')}
                  </button>
                </div>

                {/* Mailbox Configuration Form */}
                <div className="space-y-4">
                  <p className="text-sm font-medium text-slate-700">{t('settings.mailboxConfigDesc')}</p>

                  {/* Email and Password */}
                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <label className="text-sm font-medium text-slate-700">{t('settings.mailboxEmail')}</label>
                      <input
                        type="email"
                        className="w-full px-3 py-2 bg-white border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                        value={mailboxConfig.email}
                        onChange={(e) => setMailboxConfig({ ...mailboxConfig, email: e.target.value })}
                      />
                    </div>
                    <div className="space-y-2">
                      <label className="text-sm font-medium text-slate-700">{t('settings.mailboxPassword')}</label>
                      <div className="relative">
                        <input
                          type={showPassword ? "text" : "password"}
                          className="w-full px-3 py-2 pr-10 bg-white border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                          value={mailboxConfig.password}
                          onChange={(e) => setMailboxConfig({ ...mailboxConfig, password: e.target.value })}
                          placeholder="••••••••"
                        />
                        <button
                          type="button"
                          onClick={() => setShowPassword(!showPassword)}
                          className="absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
                        >
                          {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                        </button>
                      </div>
                    </div>
                  </div>

                  {/* IMAP Configuration */}
                  <div className="grid grid-cols-3 gap-4">
                    <div className="col-span-2 space-y-2">
                      <label className="text-sm font-medium text-slate-700">{t('settings.imapServer')}</label>
                      <input
                        type="text"
                        className="w-full px-3 py-2 bg-white border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                        value={mailboxConfig.imapServer}
                        onChange={(e) => setMailboxConfig({ ...mailboxConfig, imapServer: e.target.value })}
                      />
                    </div>
                    <div className="space-y-2">
                      <label className="text-sm font-medium text-slate-700">{t('settings.imapPort')}</label>
                      <input
                        type="text"
                        className="w-full px-3 py-2 bg-white border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                        value={mailboxConfig.imapPort}
                        onChange={(e) => setMailboxConfig({ ...mailboxConfig, imapPort: e.target.value })}
                      />
                    </div>
                  </div>

                  {/* SMTP Configuration */}
                  <div className="grid grid-cols-3 gap-4">
                    <div className="col-span-2 space-y-2">
                      <label className="text-sm font-medium text-slate-700">{t('settings.smtpServer')}</label>
                      <input
                        type="text"
                        className="w-full px-3 py-2 bg-white border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                        value={mailboxConfig.smtpServer}
                        onChange={(e) => setMailboxConfig({ ...mailboxConfig, smtpServer: e.target.value })}
                      />
                    </div>
                    <div className="space-y-2">
                      <label className="text-sm font-medium text-slate-700">{t('settings.smtpPort')}</label>
                      <input
                        type="text"
                        className="w-full px-3 py-2 bg-white border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                        value={mailboxConfig.smtpPort}
                        onChange={(e) => setMailboxConfig({ ...mailboxConfig, smtpPort: e.target.value })}
                      />
                    </div>
                  </div>

                  {/* Save Button */}
                  <div className="flex justify-end pt-2">
                    <button
                      onClick={handleSaveMailboxConfig}
                      className="px-6 py-2 bg-blue-600 text-white rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors"
                    >
                      {t('common.save')}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          )}

          {activeTab === "notifications" && (
            <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
              <div>
                <h3 className="text-xl font-bold text-slate-800">{t('settings.notifyPref')}</h3>
                <p className="text-slate-500 text-sm mt-1">{t('settings.notifyDesc')}</p>
              </div>

              <div className="space-y-4">
                {[t('settings.dailyDigestNotify'), t('settings.riskNotify'), t('settings.taskNotify')].map((item, i) => (
                  <div key={i} className="flex items-center justify-between p-4 border border-slate-100 rounded-xl">
                    <div>
                      <p className="font-medium text-slate-800">{item}</p>
                      <p className="text-xs text-slate-400">{t('settings.viaChannel')}</p>
                    </div>
                    <div className="w-11 h-6 bg-blue-600 rounded-full relative cursor-pointer">
                      <div className="absolute right-1 top-1 w-4 h-4 bg-white rounded-full shadow-sm"></div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Other tabs placeholders */}
          {['security', 'billing', 'appearance'].includes(activeTab) && (
            <div className="flex flex-col items-center justify-center h-full text-slate-400 animate-in fade-in duration-300">
              <div className="p-4 bg-slate-50 rounded-full mb-4">
                <Shield className="w-8 h-8" />
              </div>
              <p>该模块正在开发中...</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}