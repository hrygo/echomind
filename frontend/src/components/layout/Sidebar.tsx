"use client";

import Link from 'next/link';
import { usePathname, useSearchParams } from 'next/navigation';
import {
  Inbox,
  CheckSquare,
  LayoutDashboard,
  Sparkles,
  Share2,
  Send,
  FileEdit,
  Trash2,
  Briefcase,
  User,
  Newspaper,
  CreditCard,
  Bell,
  Zap
} from 'lucide-react';

import { OrgSwitcher } from "./OrgSwitcher";
import { useLanguage } from "@/lib/i18n/LanguageContext";

export function Sidebar() {
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const currentCategory = searchParams.get('category');
  const currentFilter = searchParams.get('filter');
  const currentFolder = searchParams.get('folder');
  const { t } = useLanguage();

  // Active state logic
  const isActive = (href: string, params?: { category?: string, filter?: string, folder?: string }) => {
    // Strict match for Dashboard
    if (href === '/dashboard' && pathname === '/dashboard' && !currentCategory && !currentFilter && !currentFolder) {
      return true;
    }

    // Match path
    if (pathname !== href.split('?')[0]) {
      return false;
    }

    // Match parameters if provided
    if (params) {
      if (params.category && currentCategory !== params.category) return false;
      if (params.filter && currentFilter !== params.filter) return false;
      if (params.folder && currentFolder !== params.folder) return false;
      
      // If params are required but url has different ones (exclusive check)
      if (!params.category && currentCategory) return false;
      if (!params.filter && currentFilter) return false;
      if (!params.folder && currentFolder) return false;
      
      return true;
    }

    // Default match for simple paths (ensure no extra params strictly)
    return !currentCategory && !currentFilter && !currentFolder;
  };

  return (
    <aside className="w-64 bg-white border-r border-slate-100 min-h-screen flex flex-col shadow-[1px_0_20px_0_rgba(0,0,0,0.02)] z-40 fixed left-0 top-0 h-full">
      {/* Logo Section */}
      <div className="h-20 flex items-center px-6 border-b border-slate-50 shrink-0">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 bg-gradient-to-br from-blue-600 to-indigo-600 rounded-xl flex items-center justify-center shadow-lg shadow-blue-200/50">
            <Zap className="text-white w-6 h-6 fill-current" />
          </div>
          <div>
            <h1 className="text-xl font-bold text-slate-800 tracking-tight leading-none">{t('common.appName')}</h1>
            <p className="text-[11px] font-medium text-slate-400 mt-1 tracking-wide">{t('common.appSlogan')}</p>
          </div>
        </div>
      </div>

      {/* Org Switcher - Hidden for Personal Intelligence Phase
      <OrgSwitcher />
      */}

      {/* Scrollable Navigation Area */}
      <div className="flex-1 flex flex-col overflow-y-auto py-4 custom-scrollbar">
        
        {/* SECTION A: INTELLIGENCE */}
        <nav className="space-y-0.5 px-2 mb-6">
          <SectionLabel>{t('sidebar.intelligence')}</SectionLabel>
          <NavItem
            href="/dashboard"
            label={t('sidebar.dashboard')}
            icon={LayoutDashboard}
            active={isActive('/dashboard')}
          />
          <NavItem
            href="/dashboard/inbox?filter=smart"
            label={t('sidebar.smartInbox')}
            icon={Sparkles}
            active={isActive('/dashboard/inbox', { filter: 'smart' })}
            iconColor="text-amber-500"
          />
          <NavItem
            href="/dashboard/tasks"
            label={t('sidebar.actionItems')}
            icon={CheckSquare}
            active={isActive('/dashboard/tasks')}
          />
          <NavItem
            href="/dashboard/insights"
            label={t('sidebar.network')}
            icon={Share2}
            active={isActive('/dashboard/insights')}
          />
        </nav>

        {/* SECTION B: MAILBOX */}
        <nav className="space-y-0.5 px-2 mb-6">
          <SectionLabel>{t('sidebar.mailbox')}</SectionLabel>
          <NavItem
            href="/dashboard/inbox"
            label={t('sidebar.inbox')}
            icon={Inbox}
            active={isActive('/dashboard/inbox')}
          />
          <NavItem
            href="/dashboard/inbox?folder=sent"
            label={t('sidebar.sent')}
            icon={Send}
            active={isActive('/dashboard/inbox', { folder: 'sent' })}
          />
          <NavItem
            href="/dashboard/inbox?folder=drafts"
            label={t('sidebar.drafts')}
            icon={FileEdit}
            active={isActive('/dashboard/inbox', { folder: 'drafts' })}
          />
          <NavItem
            href="/dashboard/inbox?folder=trash"
            label={t('sidebar.trash')}
            icon={Trash2}
            active={isActive('/dashboard/inbox', { folder: 'trash' })}
          />
        </nav>

        {/* SECTION C: LABELS */}
        <nav className="space-y-0.5 px-2">
          <SectionLabel>{t('sidebar.labels')}</SectionLabel>
          <NavItem
            href="/dashboard/inbox?category=Work"
            label={t('sidebar.work')}
            icon={Briefcase}
            active={isActive('/dashboard/inbox', { category: 'Work' })}
          />
          <NavItem
            href="/dashboard/inbox?category=Personal"
            label={t('sidebar.personal')}
            icon={User}
            active={isActive('/dashboard/inbox', { category: 'Personal' })}
          />
          <NavItem
            href="/dashboard/inbox?category=Newsletter"
            label={t('sidebar.newsletter')}
            icon={Newspaper}
            active={isActive('/dashboard/inbox', { category: 'Newsletter' })}
          />
           <NavItem
            href="/dashboard/inbox?category=Finance"
            label={t('sidebar.finance')}
            icon={CreditCard}
            active={isActive('/dashboard/inbox', { category: 'Finance' })}
          />
          <NavItem
            href="/dashboard/inbox?category=Notification"
            label={t('sidebar.notification')}
            icon={Bell}
            active={isActive('/dashboard/inbox', { category: 'Notification' })}
          />
        </nav>
      </div>

      {/* Footer: User Profile & Settings removed as it is in Header */}
    </aside>
  );
}

function SectionLabel({ children }: { children: React.ReactNode }) {
  return (
    <div className="px-3 py-1.5 mt-2 text-[10px] font-bold text-slate-400/80 uppercase tracking-widest">
      {children}
    </div>
  );
}

function NavItem({
  href,
  label,
  icon: Icon,
  active = false,
  iconColor
}: {
  href: string;
  label: string;
  icon: React.ElementType;
  active?: boolean;
  iconColor?: string;
}) {
  return (
    <Link
      href={href}
      className={`flex items-center px-3 py-2 text-[13.5px] font-medium rounded-lg transition-all duration-200 group relative ${
        active
          ? 'bg-blue-50/80 text-blue-700'
          : 'text-slate-600 hover:bg-slate-100/80 hover:text-slate-900'
      }`}
    >
      <Icon
        className={`w-[18px] h-[18px] mr-3 transition-colors ${
          active ? 'text-blue-600' : (iconColor || 'text-slate-400 group-hover:text-slate-600')
        }`}
        strokeWidth={active ? 2.5 : 2}
      />
      <span className="truncate">{label}</span>
      {active && (
        <div className="absolute right-2 w-1.5 h-1.5 bg-blue-600 rounded-full" />
      )}
    </Link>
  );
}