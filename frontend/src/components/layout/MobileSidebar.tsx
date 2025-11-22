"use client";

import { useUIStore } from "@/store/ui";
import { Sheet, SheetContent, SheetHeader, SheetTitle } from "@/components/ui/Sheet";
import { Sidebar } from "./Sidebar";
import { useLanguage } from "@/lib/i18n/LanguageContext";

export function MobileSidebar() {
  const { isMobileSidebarOpen, closeMobileSidebar } = useUIStore();
  const { t } = useLanguage();

  return (
    <Sheet open={isMobileSidebarOpen} onOpenChange={closeMobileSidebar}>
      <SheetContent side="left" className="p-0 w-72 bg-white">
        <SheetHeader className="absolute top-0 w-full p-4 border-b border-gray-200 bg-white z-10">
          <SheetTitle className="sr-only">{t('common.appName')} Navigation</SheetTitle>
        </SheetHeader>
        <Sidebar className="w-full h-full border-none shadow-none relative pt-20" />
      </SheetContent>
    </Sheet>
  );
}
