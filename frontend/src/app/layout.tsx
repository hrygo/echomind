import type { Metadata } from "next";
import "./globals.css";
import { LanguageProvider } from "@/lib/i18n/LanguageContext";
import { QueryProvider } from "@/components/providers/QueryClientProvider";
import { ToastContainer } from "@/components/ui/ToastContainer";
import { ConfirmDialog } from "@/components/ui/ConfirmDialog";
import { ThemeProvider } from "@/components/theme/ThemeProviderNext";

export const metadata: Metadata = {
  title: "EchoMind",
  description: "Your external brain for email.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className="antialiased"
        suppressHydrationWarning
      >
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          <QueryProvider>
            <LanguageProvider>
              {children}
              <ToastContainer />
              <ConfirmDialog />
            </LanguageProvider>
          </QueryProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}