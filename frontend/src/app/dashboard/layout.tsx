import { Sidebar } from "@/components/layout/Sidebar";
import { Header } from "@/components/layout/Header";
import AuthGuard from '@/components/auth/AuthGuard';
import { LanguageProvider } from "@/lib/i18n/LanguageContext";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <LanguageProvider>
      <AuthGuard>
        <div className="flex h-screen bg-slate-50">
          <Sidebar />
          <main className="flex-1 flex flex-col overflow-hidden relative">
            <Header />

            {/* Page Content */}
            <div className="flex-1 overflow-y-auto scroll-smooth">
              {/* Added a container to constrain width on large screens for better readability */}
              <div className="max-w-7xl mx-auto p-8 w-full">
                {children}
              </div>
            </div>
          </main>
        </div>
      </AuthGuard>
    </LanguageProvider>
  );
}