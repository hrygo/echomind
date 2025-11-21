import NetworkGraphClient from './NetworkGraphClient';

export default function InsightsPage() {
    return (
        <section className="h-[calc(100vh-6rem)] flex flex-col p-6">
            <div className="flex items-center justify-between mb-6">
                <h1 className="text-2xl font-bold text-slate-800">智能洞察</h1>
            </div>
            <div className="flex-1 bg-white rounded-2xl shadow-sm border border-slate-200/60 overflow-hidden flex flex-col">
                <div className="p-4 border-b border-slate-100 bg-slate-50/30">
                    <h2 className="text-sm font-semibold text-slate-600 uppercase tracking-wider">联系人关系图谱</h2>
                </div>
                <div className="flex-1 relative w-full min-h-0">
                    <div className="absolute inset-0">
                        <NetworkGraphClient />
                    </div>
                </div>
            </div>
        </section>
    );
}
