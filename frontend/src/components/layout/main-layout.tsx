"use client";

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { ReactNode } from 'react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import {
  Home,
  Settings,
  CreditCard,
  Landmark,
  DollarSign,
  RefreshCcw,
  BarChart3,
  CalendarDays,
  type LucideIcon,
} from 'lucide-react';

export function MainLayout({ children }: { children: ReactNode }) {
  const pathname = usePathname();

  const NavLink = ({ href, label, icon: Icon }: { href: string; label: string; icon: LucideIcon }) => {
    const active = pathname === href || (href !== '/' && pathname?.startsWith(href));
    return (
      <Link 
        href={href} 
        className={cn(
          'flex items-center gap-2 px-3 py-2 rounded-md text-sm hover:bg-accent', 
          active && 'bg-accent font-medium'
        )}
      >
        <Icon className="h-4 w-4" />
        <span>{label}</span>
      </Link>
    );
  };

  return (
    <div className="min-h-screen bg-background text-foreground">
      <header className="border-b">
        <div className="container mx-auto flex items-center justify-between h-14 px-4">
          <Link href="/" className="font-semibold text-lg">Flow Sight</Link>
          <nav className="flex items-center gap-2 flex-wrap">
            <NavLink href="/dashboard" label="ダッシュボード" icon={Home} />
            <NavLink href="/cashflow" label="キャッシュフロー" icon={BarChart3} />
            <NavLink href="/income" label="収入" icon={DollarSign} />
            <NavLink href="/recurring-payments" label="定期支出" icon={RefreshCcw} />
            <NavLink href="/bank-accounts" label="銀行口座" icon={Landmark} />
            <NavLink href="/credit-cards" label="クレジットカード" icon={CreditCard} />
            <NavLink href="/card-monthly-totals" label="カード月次集計" icon={CalendarDays} />
            <NavLink href="/settings" label="設定" icon={Settings} />
          </nav>
          <div className="flex items-center gap-2">
            <Button variant="outline" asChild>
              <a href="https://github.com/" target="_blank" rel="noreferrer">GitHub</a>
            </Button>
          </div>
        </div>
      </header>

      <main className="container mx-auto px-4 py-6">
        {children}
      </main>

      <footer className="border-t">
        <div className="container mx-auto px-4 py-4 text-sm text-muted-foreground">
          © {new Date().getFullYear()} Flow Sight
        </div>
      </footer>
    </div>
  );
}
