"use client";

import Link from 'next/link';
import { usePathname, useSearchParams } from 'next/navigation';
import {
  Inbox,
  CheckSquare,
  BarChart2,
  Settings,
  Briefcase,
  User,
  Newspaper,
  Bell,
  AlertOctagon,
  Zap,
  LayoutDashboard
} from 'lucide-react';

import { useLanguage } from "@/lib/i18n/LanguageContext";

export function Sidebar() {
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const currentCategory = searchParams.get('category');
  const { t } = useLanguage();

  const isActive = (href: string, category?: string) => {
    if (category) {
      return currentCategory === category && pathname === '/dashboard/inbox';
    }
    if (href === '/dashboard' && pathname === '/dashboard' && !currentCategory) {
      return true;
    }
    if (href !== '/dashboard' && pathname.startsWith(href)) {
      return true;
    }
    return false;
  };

  return (
    <aside className="w-72 bg-white border-r border-slate-100 min-h-screen flex flex-col shadow-[1px_0_20px_0_rgba(0,0,0,0.02)] z-20">
      {/* Logo Section */}
      <div className="h-20 flex items-center px-6 border-b border-slate-50">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 bg-gradient-to-br from-blue-600 to-indigo-600 rounded-xl flex items-center justify-center shadow-lg shadow-blue-200">
            <Zap className="text-white w-6 h-6 fill-current" />
          </div>
          <div>
            <h1 className="text-xl font-bold text-slate-800 tracking-tight leading-tight">{t('common.appName')}</h1>
            <p className="text-[10px] font-medium text-slate-400 uppercase tracking-wider">{t('common.appSlogan')}</p>
          </div>
        </div>
      </div>

      <div className="flex-1 flex flex-col overflow-y-auto py-6">
        {/* Main Navigation */}
        <nav className="space-y-1 px-3">
          <SectionLabel>{t('sidebar.mainMenu')}</SectionLabel>
          <NavItem
            href="/dashboard"
            label={t('sidebar.dashboard')}
            icon={LayoutDashboard}
            active={isActive('/dashboard') && pathname === '/dashboard'}
          />
          <NavItem
            href="/dashboard/inbox"
            label={t('sidebar.inbox')}
            icon={Inbox}
            active={isActive('/dashboard/inbox') || (pathname === '/dashboard' && !!currentCategory)}
          />
          <NavItem
            href="/dashboard/tasks"
            label={t('sidebar.tasks')}
            icon={CheckSquare}
            active={isActive('/dashboard/tasks')}
          />
          <NavItem
            href="/dashboard/insights"
            label={t('sidebar.insights')}
            icon={BarChart2}
            active={isActive('/dashboard/insights')}
          />
        </nav>

        {/* Spacer */}
        <div className="h-8"></div>

        {/* Smart Folders */}
        <div className="space-y-1 px-3">
          <SectionLabel>{t('sidebar.smartFolders')}</SectionLabel>
          <NavItem
            href="/dashboard/inbox?category=Work"
            label={t('sidebar.work')}
            icon={Briefcase}
            active={isActive('/dashboard/inbox', 'Work')}
          />
          <NavItem
            href="/dashboard/inbox?category=Personal"
            label={t('sidebar.personal')}
            icon={User}
            active={isActive('/dashboard/inbox', 'Personal')}
          />
          <NavItem
            href="/dashboard/inbox?category=Newsletter"
            label={t('sidebar.newsletter')}
            icon={Newspaper}
            active={isActive('/dashboard/inbox', 'Newsletter')}
          />
          <NavItem
            href="/dashboard/inbox?category=Notification"
            label={t('sidebar.notification')}
            icon={Bell}
            active={isActive('/dashboard/inbox', 'Notification')}
          />
          <NavItem
            href="/dashboard/inbox?category=Spam"
            label={t('sidebar.spam')}
            icon={AlertOctagon}
            active={isActive('/dashboard/inbox', 'Spam')}
          />
        </div>
      </div>
    </aside>
  );
}

function SectionLabel({ children }: { children: React.ReactNode }) {
  return (
    <div className="px-4 py-2 text-[11px] font-bold text-slate-400 uppercase tracking-widest">
      {children}
    </div>
  );
}

function NavItem({
  href,
  label,
  icon: Icon,
  active = false
}: {
  href: string;
  label: string;
  icon: any;
  active?: boolean
}) {
  return (
    <Link
      href={href}
      className={`flex items-center px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200 group ${active
        ? 'bg-blue-50 text-blue-700 shadow-sm'
        : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'
        }`}
    >
      <Icon
        className={`w-5 h-5 mr-3 transition-colors ${active ? 'text-blue-600' : 'text-slate-400 group-hover:text-slate-600'
          }`}
        strokeWidth={active ? 2.5 : 2}
      />
      {label}
      {active && (
        <div className="ml-auto w-1.5 h-1.5 bg-blue-600 rounded-full" />
      )}
    </Link>
  );
}

