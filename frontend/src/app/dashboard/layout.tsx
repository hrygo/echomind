import { Sidebar } from "@/components/layout/Sidebar";
import AuthGuard from '@/components/auth/AuthGuard';

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <AuthGuard>
      <div className="flex h-screen bg-gray-100">
        <Sidebar />
        <main className="flex-1 overflow-y-auto p-8 text-gray-900">
          {children}
        </main>
      </div>
    </AuthGuard>
  );
}