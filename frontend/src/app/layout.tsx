import type { Metadata } from "next";

import "./globals.css";
import { Geist } from "next/font/google";
import { cn } from "@/lib/utils";
import { Providers } from "./providers";
import { ThemeToggle } from "@/components/theme-toggle";

const geist = Geist({subsets:['latin'],variable:'--font-sans'});

export const metadata: Metadata = {
  title: "My fitness",
  description: "desc",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="en" className={cn("font-sans", geist.variable)} suppressHydrationWarning
    >
      <body className="min-h-full flex flex-col">
        <Providers>
          <header className="flex items-center justify-between border-b px-6 py-4">
            <span className="font-semibold">My Fitness</span>
            <ThemeToggle />
          </header>
          <main className="flex-1">{children}</main>
        </Providers>
      </body>
    </html>
  );
}
