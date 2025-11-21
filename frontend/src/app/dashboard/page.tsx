"use client";

import { useState } from "react";
import { ExecutiveView } from "@/components/dashboard/ExecutiveView";
import { ManagerView } from "@/components/dashboard/ManagerView";
import { DealmakerView } from "@/components/dashboard/DealmakerView";
import { AIBriefingHeader } from "@/components/dashboard/AIBriefingHeader";
import { LayoutDashboard, Briefcase, Radar } from "lucide-react";
import { useLanguage } from "@/lib/i18n/LanguageContext";

type ViewType = "executive" | "manager" | "dealmaker";

export default function DashboardHomePage() {
    const [currentView, setCurrentView] = useState<ViewType>("executive");
    const { t } = useLanguage();

    return (
        <div className="min-h-[calc(100vh-theme(spacing.20)-theme(spacing.16))] flex flex-col">
            
            {/* AI Header */}
            <AIBriefingHeader currentView={currentView} />

            {/* View Switcher Tabs */}
            <div className="flex items-center gap-1 mb-8 bg-slate-100 p-1.5 rounded-xl w-fit border border-slate-200/50">
                <ViewTab 
                    active={currentView === "executive"} 
                    onClick={() => setCurrentView("executive")}
                    icon={LayoutDashboard}
                    label={t('dashboard.executiveView')}
                />
                <ViewTab 
                    active={currentView === "manager"} 
                    onClick={() => setCurrentView("manager")}
                    icon={Briefcase}
                    label={t('dashboard.managerView')}
                />
                <ViewTab 
                    active={currentView === "dealmaker"} 
                    onClick={() => setCurrentView("dealmaker")}
                    icon={Radar}
                    label={t('dashboard.dealmakerView')}
                />
            </div>

            {/* View Content */}
            <div className="flex-1 animate-in fade-in slide-in-from-bottom-4 duration-300">
                {currentView === "executive" && <ExecutiveView />}
                {currentView === "manager" && <ManagerView />}
                {currentView === "dealmaker" && <DealmakerView />}
            </div>
        </div>
    );
}

function ViewTab({ active, onClick, icon: Icon, label }: { active: boolean, onClick: () => void, icon: React.ElementType, label: string }) {
    return (
        <button
            onClick={onClick}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-semibold transition-all duration-200 ${
                active
                    ? "bg-white text-blue-600 shadow-sm"
                    : "text-slate-500 hover:text-slate-700 hover:bg-slate-200/50"
            }`}
        >
            <Icon className={`w-4 h-4 ${active ? "text-blue-600" : "text-slate-400"}`} />
            {label}
        </button>
    );
}