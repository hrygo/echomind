import NetworkGraphClient from './NetworkGraphClient';

export default function InsightsPage() {
    return (
        <section className="min-h-screen p-8">
            <h1 className="text-3xl font-bold mb-6">智能洞察</h1>
            <div className="bg-card p-6 rounded-lg shadow-md">
                <h2 className="text-xl font-semibold mb-4">联系人关系图谱</h2>
                <div className="h-[600px] w-full">
                    <NetworkGraphClient />
                </div>
            </div>
        </section>
    );
}
