import dynamic from 'next/dynamic';

const NetworkGraph = dynamic(() => import('@/components/insights/NetworkGraph'), { ssr: false });

export default function InsightsPage() {
    return (
        <div className="space-y-6">
            <h1 className="text-3xl font-bold tracking-tight">Insights</h1>
            <NetworkGraph />
        </div>
    );
}
