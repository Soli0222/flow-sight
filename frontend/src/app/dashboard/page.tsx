'use client'

import { useAuth } from '@/components/providers/auth-provider'
import { MainLayout } from '@/components/layout/main-layout'
import { useRouter } from 'next/navigation'
import { useEffect, useState, useCallback } from 'react'
import { useApi } from '@/components/providers/api-provider'
import { DashboardSummary } from '@/types/api'
import { formatCurrency, formatDate } from '@/lib/utils-extended'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { TrendingUp, TrendingDown, Wallet, CreditCard as CreditCardIcon } from 'lucide-react'
import { toast } from 'sonner'
import { BalanceTrendChart } from '@/components/dashboard/balance-trend-chart'

export default function DashboardPage() {
  const { user, isLoading: authLoading } = useAuth()
  const router = useRouter()
  const apiClient = useApi()
  
  const [dashboardData, setDashboardData] = useState<DashboardSummary | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  const loadDashboardData = useCallback(async () => {
    try {
      setIsLoading(true)
      const summary = await apiClient.getDashboardSummary()
      
      // データの安全性チェック
      if (summary && typeof summary === 'object') {
        // recent_activitiesが配列でない場合は空配列に設定
        if (!Array.isArray(summary.recent_activities)) {
          summary.recent_activities = []
        }
        setDashboardData(summary)
      } else {
        console.error('Invalid dashboard data received:', summary)
        toast.error('ダッシュボードデータの形式が正しくありません')
      }
    } catch (error) {
      console.error('Failed to load dashboard data:', error)
      toast.error('ダッシュボードデータの取得に失敗しました')
    } finally {
      setIsLoading(false)
    }
  }, [apiClient])

  useEffect(() => {
    if (!authLoading && !user) {
      router.push('/login')
    }
  }, [user, authLoading, router])

  useEffect(() => {
    if (user && !authLoading) {
      loadDashboardData()
    }
  }, [user, authLoading, loadDashboardData])

  if (authLoading || isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div>Loading...</div>
      </div>
    )
  }

  if (!user || !dashboardData) {
    return null
  }

  return (
    <MainLayout>
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold">ダッシュボード</h1>
          <p className="text-muted-foreground">
            {user.name}さん、Flow Sightへようこそ。金融管理の概要をここで確認できます。
          </p>
        </div>

        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">総残高</CardTitle>
              <Wallet className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {formatCurrency(dashboardData.total_balance)}
              </div>
              <p className="text-xs text-muted-foreground">
                全銀行口座の合計
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">今月の収入</CardTitle>
              <TrendingUp className="h-4 w-4 text-green-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-green-600">
                {formatCurrency(dashboardData.monthly_income)}
              </div>
              <p className="text-xs text-muted-foreground">
                当月の予定収入
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">今月の支出</CardTitle>
              <TrendingDown className="h-4 w-4 text-red-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-red-600">
                {formatCurrency(dashboardData.monthly_expense)}
              </div>
              <p className="text-xs text-muted-foreground">
                当月の予定支出
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">管理資産数</CardTitle>
              <CreditCardIcon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {dashboardData.total_assets}
              </div>
              <p className="text-xs text-muted-foreground">
                口座 + カード
              </p>
            </CardContent>
          </Card>
        </div>

        <BalanceTrendChart />

        <Card>
          <CardHeader>
            <CardTitle className="text-xl font-semibold">最近の活動</CardTitle>
          </CardHeader>
          <CardContent>
            {!dashboardData.recent_activities || dashboardData.recent_activities.length === 0 ? (
              <p className="text-muted-foreground">
                まだデータがありません。銀行口座や資産を追加して始めましょう。
              </p>
            ) : (
              <div className="space-y-4">
                {dashboardData.recent_activities.map((activity, index) => (
                  <div key={index} className="flex items-center justify-between p-4 border rounded-lg">
                    <div className="flex-1">
                      <div className="font-medium">
                        {formatDate(activity?.date || '')}
                      </div>
                      <div className="text-sm text-muted-foreground">
                        {(activity?.details || []).map((detail, detailIndex) => (
                          <div key={detailIndex}>• {detail?.description || ''}</div>
                        ))}
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="text-sm">
                        {activity.income > 0 && (
                          <span className="text-green-600">+{formatCurrency(activity.income)}</span>
                        )}
                        {activity.income > 0 && activity.expense > 0 && " | "}
                        {activity.expense > 0 && (
                          <span className="text-red-600">-{formatCurrency(activity.expense)}</span>
                        )}
                      </div>
                      <div className={`text-sm font-medium ${activity.balance >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                        残高: {formatCurrency(activity.balance)}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </MainLayout>
  )
}
