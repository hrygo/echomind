import type { Metadata } from "next";
import "./globals.css";
import { LanguageProvider } from "@/lib/i18n/LanguageContext";
import { ToastContainer } from '@/components/ui/Toast'; // Import ToastContainer

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
    <html lang="en">
      <body className="bg-gray-100 text-gray-900 antialiased">
        <LanguageProvider>
          {children}
          <ToastContainer />
        </LanguageProvider>
      </body>
    </html>
  );
}