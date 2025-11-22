import NetworkGraphClient from './NetworkGraphClient';
import InsightsHeader from './InsightsHeader'; // Import the new client component

export default function InsightsPage() {
    return (
        <section className="h-[calc(100vh-6rem)] flex flex-col p-6">
            <InsightsHeader /> {/* Use the client component for translated headers */}
            <div className="flex-1 bg-white rounded-2xl shadow-sm border border-slate-200/60 overflow-hidden flex flex-col">
                
                <div className="flex-1 relative w-full min-h-0">
                    <div className="absolute inset-0">
                        <NetworkGraphClient />
                    </div>
                </div>
            </div>
        </section>
    );
}
