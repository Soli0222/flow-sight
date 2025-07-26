import { MainLayout } from '@/components/layout/main-layout';

export default function Home() {
  return (
    <MainLayout>
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold">ダッシュボード</h1>
          <p className="text-muted-foreground">
            Flow Sightへようこそ。金融管理の概要をここで確認できます。
          </p>
        </div>

        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
          <div className="rounded-lg border bg-card p-6">
            <h3 className="text-sm font-medium text-muted-foreground">総残高</h3>
            <p className="text-2xl font-bold">¥0</p>
          </div>
          <div className="rounded-lg border bg-card p-6">
            <h3 className="text-sm font-medium text-muted-foreground">今月の収入</h3>
            <p className="text-2xl font-bold">¥0</p>
          </div>
          <div className="rounded-lg border bg-card p-6">
            <h3 className="text-sm font-medium text-muted-foreground">今月の支出</h3>
            <p className="text-2xl font-bold">¥0</p>
          </div>
          <div className="rounded-lg border bg-card p-6">
            <h3 className="text-sm font-medium text-muted-foreground">資産数</h3>
            <p className="text-2xl font-bold">0</p>
          </div>
        </div>

        <div className="rounded-lg border bg-card p-6">
          <h2 className="text-xl font-semibold mb-4">最近の活動</h2>
          <p className="text-muted-foreground">
            まだデータがありません。銀行口座や資産を追加して始めましょう。
          </p>
        </div>
      </div>
    </MainLayout>
  );
}
