'use client';

import React from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { usePathname, useRouter } from 'next/navigation';
import { 
  LayoutDashboard, 
  CreditCard, 
  Landmark, 
  TrendingUp, 
  Repeat,
  BarChart3,
  Settings,
  Menu,
  X,
  LogOut,
  User,
  Calendar
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { ThemeToggle } from '@/components/providers/theme-toggle';
import { useAuth } from '@/components/providers/auth-provider';
import { cn } from '@/lib/utils';
import { useState } from 'react';

const navigation = [
  { name: 'ダッシュボード', href: '/dashboard', icon: LayoutDashboard },
  { name: '銀行口座', href: '/bank-accounts', icon: Landmark },
  { name: 'クレジットカード管理', href: '/credit-cards', icon: CreditCard },
  { name: 'カード月次利用額', href: '/card-monthly-totals', icon: Calendar },
  { name: '定期支払い', href: '/recurring-payments', icon: Repeat },
  { name: '収入管理', href: '/income', icon: TrendingUp },
  { name: 'キャッシュフロー', href: '/cashflow', icon: BarChart3 },
  { name: '設定', href: '/settings', icon: Settings },
];

interface MainLayoutProps {
  children: React.ReactNode;
}

export function MainLayout({ children }: MainLayoutProps) {
  const pathname = usePathname();
  const router = useRouter();
  const { user, logout } = useAuth();
  const [sidebarOpen, setSidebarOpen] = useState(false);

  const handleLogout = () => {
    logout();
    router.push('/login');
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Mobile sidebar */}
      <div className={cn(
        "fixed inset-0 z-40 lg:hidden",
        sidebarOpen ? "block" : "hidden"
      )}>
        <div className="fixed inset-0 bg-black bg-opacity-25" onClick={() => setSidebarOpen(false)} />
        <nav className="fixed top-0 left-0 bottom-0 flex flex-col w-64 bg-card border-r">
          <div className="flex items-center justify-between p-4">
            <h1 className="text-xl font-bold">Flow Sight</h1>
            <Button variant="ghost" size="icon" onClick={() => setSidebarOpen(false)}>
              <X className="h-5 w-5" />
            </Button>
          </div>
          <div className="flex-1 px-2 pb-4">
            {navigation.map((item) => {
              const Icon = item.icon;
              const isActive = pathname === item.href;
              return (
                <Link
                  key={item.name}
                  href={item.href}
                  className={cn(
                    "flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                    isActive
                      ? "bg-primary text-primary-foreground"
                      : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
                  )}
                  onClick={() => setSidebarOpen(false)}
                >
                  <Icon className="h-4 w-4" />
                  {item.name}
                </Link>
              );
            })}
          </div>
        </nav>
      </div>

      {/* Desktop sidebar */}
      <nav className="hidden lg:fixed lg:inset-y-0 lg:z-40 lg:flex lg:w-64 lg:flex-col">
        <div className="flex flex-col flex-1 min-h-screen bg-card border-r">
          <div className="flex items-center h-16 px-4">
            <h1 className="text-xl font-bold">Flow Sight</h1>
          </div>
          <div className="flex-1 px-2 pb-4">
            {navigation.map((item) => {
              const Icon = item.icon;
              const isActive = pathname === item.href;
              return (
                <Link
                  key={item.name}
                  href={item.href}
                  className={cn(
                    "flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                    isActive
                      ? "bg-primary text-primary-foreground"
                      : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
                  )}
                >
                  <Icon className="h-4 w-4" />
                  {item.name}
                </Link>
              );
            })}
          </div>
        </div>
      </nav>

      {/* Main content */}
      <div className="lg:pl-64">
        {/* Top header */}
        <header className="flex items-center justify-between h-16 px-4 border-b bg-card">
          <Button
            variant="ghost"
            size="icon"
            className="lg:hidden"
            onClick={() => setSidebarOpen(true)}
          >
            <Menu className="h-5 w-5" />
          </Button>
          <div className="flex items-center gap-4 ml-auto">
            {user && (
              <div className="flex items-center gap-2">
                {user.picture ? (
                  <Image
                    src={user.picture}
                    alt={user.name}
                    width={32}
                    height={32}
                    className="rounded-full"
                  />
                ) : (
                  <User className="h-8 w-8" />
                )}
                <span className="text-sm">{user.name}</span>
              </div>
            )}
            <ThemeToggle />
            <Button
              variant="ghost"
              size="icon"
              onClick={handleLogout}
              title="ログアウト"
            >
              <LogOut className="h-4 w-4" />
            </Button>
          </div>
        </header>

        {/* Page content */}
        <main className="p-6 min-h-screen bg-background">
          {children}
        </main>
      </div>
    </div>
  );
}
