import React from "react";
import "./globals.css";

export const metadata = {
  title: "Neko-Sync",
  description: "Anime, manga, and content synchronization platform",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
