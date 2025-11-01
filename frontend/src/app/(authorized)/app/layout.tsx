import Header from '@/components/layout/Header';
import Sidebar from '@/components/layout/Sidebar';

export default function HomeLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex flex-col min-h-screen">
      <Header />
      <div className="flex flex-1">
        {/*
        <Sidebar />
        */}
        <main className="flex-1 p-8 bg-gray-100">
          {children}
        </main>
      </div>
    </div>
  );
}
