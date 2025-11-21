"use client";

import { useState } from "react";
import { ExecutiveView } from "@/components/dashboard/ExecutiveView";
import { ManagerView } from "@/components/dashboard/ManagerView";
import { DealmakerView } from "@/components/dashboard/DealmakerView";
import { LayoutDashboard, Briefcase, Radar } from "lucide-react";

type ViewType = "executive" | "manager" | "dealmaker";

export default function DashboardHomePage() {
    const [currentView, setCurrentView] = useState<ViewType>("executive");

    return (
        <div className="min-h-[calc(100vh-theme(spacing.20)-theme(spacing.16))] flex flex-col space-y-6">
            {/* Header & View Switcher */}
            <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div>
                    <h1 className="text-3xl font-bold text-slate-800 tracking-tight">工作台</h1>
                    <p className="text-slate-500 text-sm mt-1">
                        {currentView === "executive" && "Morning Report - 全局视野，辅助决策"}
                        {currentView === "manager" && "Control Console - 任务闭环，高效执行"}
                        {currentView === "dealmaker" && "Radar - 洞察意图，激活关系"}
                    </p>
                </div>

                <div className="bg-slate-100 p-1 rounded-xl flex items-center self-start md:self-auto">
                    <button
                        onClick={() => setCurrentView("executive")}
                        className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 ${currentView === "executive"
                                ? "bg-white text-blue-600 shadow-sm"
                                : "text-slate-500 hover:text-slate-700 hover:bg-slate-200/50"
                            }`}
                    >
                        <LayoutDashboard className="w-4 h-4" />
                        高管视图
                    </button>
                    <button
                        onClick={() => setCurrentView("manager")}
                        className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 ${currentView === "manager"
                                ? "bg-white text-blue-600 shadow-sm"
                                : "text-slate-500 hover:text-slate-700 hover:bg-slate-200/50"
                            }`}
                    >
                        <Briefcase className="w-4 h-4" />
                        管理者视图
                    </button>
                    <button
                        onClick={() => setCurrentView("dealmaker")}
                        className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 ${currentView === "dealmaker"
                                ? "bg-white text-blue-600 shadow-sm"
                                : "text-slate-500 hover:text-slate-700 hover:bg-slate-200/50"
                            }`}
                    >
                        <Radar className="w-4 h-4" />
                        销售视图
                    </button>
                </div>
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
