import type { Metadata } from "next";
import "@/styles/globals.css";

export const metadata: Metadata = {
  title: "HARaMA - AI Exam Grading",
  description: "Human-Assisted Reasoning & Automated Marking Assistant",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="antialiased bg-gray-50 min-h-screen">
        <nav className="bg-white border-b border-gray-200 px-8 py-4 flex justify-between items-center sticky top-0 z-50">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 bg-blue-600 rounded flex items-center justify-center text-white font-bold">H</div>
            <span className="text-xl font-bold text-gray-800 tracking-tight">HARaMA</span>
          </div>
          <div className="flex gap-6 text-sm font-medium text-gray-600">
            <a href="/" className="hover:text-blue-600">Home</a>
            <a href="/dashboard" className="hover:text-blue-600">Dashboard</a>
            <a href="/exams" className="hover:text-blue-600">Exams</a>
          </div>
          <div>
            <div className="w-8 h-8 rounded-full bg-gray-200 border border-gray-300"></div>
          </div>
        </nav>
        <main>{children}</main>
      </body>
    </html>
  );
}
