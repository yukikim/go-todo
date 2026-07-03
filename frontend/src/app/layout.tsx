import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Go Todo",
  description: "Todo dashboard for the Go API",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja">
      <body>{children}</body>
    </html>
  );
}
