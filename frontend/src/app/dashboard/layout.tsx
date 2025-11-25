import { Sidebar } from "@/components/layout/Sidebar";
import { Header } from "@/components/layout/Header";
import { MobileSidebar } from "@/components/layout/MobileSidebar";
import AuthGuard from '@/components/auth/AuthGuard';

// Force dynamic rendering for all dashboard pages
export const dynamic = 'force-dynamic';


export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <AuthGuard>
      <div className="flex h-screen bg-slate-50">
        {/* Desktop Sidebar - Fixed */}
        <Sidebar className="hidden md:flex" />

        {/* Mobile Sidebar - Sheet */}
        <MobileSidebar />

        <main className="flex-1 flex flex-col relative w-full md:ml-64">
          <Header />

          {/* Page Content */}
          <div className="flex-1 overflow-y-auto scroll-smooth">
            {/* Added a container to constrain width on large screens for better readability */}
            <div className="max-w-7xl mx-auto p-6 md:p-10 w-full">
              {children}
            </div>
          </div>
        </main>
      </div>
    </AuthGuard>
  );
}
