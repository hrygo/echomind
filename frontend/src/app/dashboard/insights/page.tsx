import NetworkGraphClient from './NetworkGraphClient';
import InsightsHeader from './InsightsHeader'; // Import the new client component

export default function InsightsPage() {
    return (
        <section className="min-h-[calc(100vh-theme(spacing.20)-theme(spacing.16))] flex flex-col">
            <InsightsHeader />
            <div className="flex-1 bg-white rounded-2xl shadow-sm border border-slate-200/60 overflow-hidden relative">
                <NetworkGraphClient />
            </div>
        </section>
    );
}
