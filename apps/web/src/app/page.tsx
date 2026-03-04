"use client";

import { useState, useEffect } from "react";

interface HealthStatus {
  status: string;
  timestamp: string;
}

export default function HomePage() {
  const [health, setHealth] = useState<HealthStatus | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch("/api/health")
      .then((res) => res.json())
      .then((data) => {
        setHealth({
          status: data.status || "OK",
          timestamp: new Date().toISOString(),
        });
        setLoading(false);
      })
      .catch(() => {
        setHealth({ status: "Error", timestamp: new Date().toISOString() });
        setLoading(false);
      });
  }, []);

  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-900 via-blue-900 to-indigo-900">
      <div className="container mx-auto px-4 py-16">
        <div className="text-center mb-16">
          <h1 className="text-6xl font-bold text-white mb-4">🐱 Neko-Sync</h1>
          <p className="text-xl text-gray-300 mb-8">
            Anime, manga, and content synchronization platform
          </p>

          {loading ? (
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-white"></div>
          ) : (
            <div
              className={`inline-block px-4 py-2 rounded-full text-sm font-medium ${
                health?.status === "OK"
                  ? "bg-green-500 text-white"
                  : "bg-red-500 text-white"
              }`}
            >
              Backend: {health?.status}
            </div>
          )}
        </div>

        <div className="grid md:grid-cols-3 gap-8 mb-16">
          <div className="bg-white/10 backdrop-blur-lg rounded-lg p-6 text-white">
            <h3 className="text-xl font-semibold mb-4">📺 Watch Parties</h3>
            <p className="text-gray-300">
              Synchronized viewing experience with friends. Watch anime together
              in real-time.
            </p>
          </div>

          <div className="bg-white/10 backdrop-blur-lg rounded-lg p-6 text-white">
            <h3 className="text-xl font-semibold mb-4">📱 Cross-Device Sync</h3>
            <p className="text-gray-300">
              Continue watching seamlessly across all your devices. Never lose
              your progress.
            </p>
          </div>

          <div className="bg-white/10 backdrop-blur-lg rounded-lg p-6 text-white">
            <h3 className="text-xl font-semibold mb-4">🤝 Social Features</h3>
            <p className="text-gray-300">
              Follow friends, share recommendations, and discover new content
              together.
            </p>
          </div>
        </div>

        <div className="text-center">
          <div className="space-y-4">
            <h2 className="text-2xl font-semibold text-white mb-8">
              API Status
            </h2>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 max-w-2xl mx-auto">
              <div className="bg-white/10 backdrop-blur-lg rounded-lg p-4">
                <div className="text-green-400 font-semibold">✓ Health</div>
                <div className="text-sm text-gray-300">/api/health</div>
              </div>
              <div className="bg-white/10 backdrop-blur-lg rounded-lg p-4">
                <div className="text-yellow-400 font-semibold">○ Users</div>
                <div className="text-sm text-gray-300">/api/v1/users</div>
              </div>
              <div className="bg-white/10 backdrop-blur-lg rounded-lg p-4">
                <div className="text-yellow-400 font-semibold">○ Content</div>
                <div className="text-sm text-gray-300">/api/v1/content</div>
              </div>
              <div className="bg-white/10 backdrop-blur-lg rounded-lg p-4">
                <div className="text-yellow-400 font-semibold">○ Parties</div>
                <div className="text-sm text-gray-300">/api/v1/parties</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
