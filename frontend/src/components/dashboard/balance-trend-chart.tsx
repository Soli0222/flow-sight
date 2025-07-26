'use client'

import { useState, useEffect } from 'react'
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useApi } from '@/components/providers/api-provider'
import { useTheme } from 'next-themes'
import { formatCurrency } from '@/lib/utils-extended'
import { toast } from 'sonner'

interface ChartDataPoint {
  date: string
  balance: number
  formattedDate: string
}

interface CustomTooltipProps {
  active?: boolean
  payload?: Array<{
    payload: ChartDataPoint
    value: number
  }>
  label?: string
}

export function BalanceTrendChart() {
  const apiClient = useApi()
  const { theme } = useTheme()
  const [chartData, setChartData] = useState<ChartDataPoint[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [chartColors, setChartColors] = useState({
    primary: '#3b82f6',
    muted: '#6b7280',
    border: '#e5e7eb'
  })

  // テーマに基づいてCSS変数から実際の色を取得
  useEffect(() => {
    const updateColors = () => {
      const root = document.documentElement
      const computedStyle = getComputedStyle(root)
      
      // CSS変数から計算された色を取得
      const primaryHsl = computedStyle.getPropertyValue('--primary').trim()
      
      // oklch値をhslに変換する関数
      const oklchToHsl = (oklchValue: string): string => {
        if (oklchValue.includes('oklch')) {
          // テーマに応じてフォールバック色を使用
          if (theme === 'dark') {
            return '#60a5fa' // blue-400相当
          } else {
            return '#3b82f6' // blue-500相当
          }
        }
        return oklchValue
      }
      
      setChartColors({
        primary: oklchToHsl(primaryHsl),
        muted: theme === 'dark' ? '#9ca3af' : '#6b7280',
        border: theme === 'dark' ? '#374151' : '#e5e7eb'
      })
    }
    
    updateColors()
    
    // テーマ変更時に色を更新
    const observer = new MutationObserver(updateColors)
    observer.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ['class']
    })
    
    return () => observer.disconnect()
  }, [theme])

  useEffect(() => {
    const loadChartData = async () => {
      try {
        setIsLoading(true)
        
        // 前後1年（24ヶ月）のキャッシュフロー予測を取得
        const projections = await apiClient.getCashflowProjection(24, false)
        
        // 月末残高データを作成（月次集計）
        const monthlyData = new Map<string, { balance: number; lastDate: string }>()
        
        projections.forEach((projection) => {
          const date = new Date(projection.date)
          const yearMonth = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
          
          // 各月の最後の日の残高を記録
          const existing = monthlyData.get(yearMonth)
          if (!existing || new Date(projection.date) > new Date(existing.lastDate)) {
            monthlyData.set(yearMonth, {
              balance: projection.balance,
              lastDate: projection.date
            })
          }
        })

        // チャート用のデータ形式に変換
        const chartPoints: ChartDataPoint[] = Array.from(monthlyData.entries())
          .map(([yearMonth, data]) => {
            const [year, month] = yearMonth.split('-')
            return {
              date: yearMonth,
              balance: data.balance,
              formattedDate: `${year}年${parseInt(month)}月`
            }
          })
          .sort((a, b) => a.date.localeCompare(b.date))

        setChartData(chartPoints)
      } catch (error) {
        console.error('Failed to load chart data:', error)
        toast.error('グラフデータの取得に失敗しました')
      } finally {
        setIsLoading(false)
      }
    }

    loadChartData()
  }, [apiClient])

  const CustomTooltip = ({ active, payload }: CustomTooltipProps) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload
      return (
        <div className="bg-background/95 backdrop-blur-sm border rounded-lg shadow-lg p-3 border-border">
          <p className="font-medium text-foreground">{data.formattedDate}</p>
          <p className={`text-sm font-medium ${data.balance >= 0 ? 'text-green-600' : 'text-red-600'}`}>
            残高: {formatCurrency(data.balance)}
          </p>
        </div>
      )
    }
    return null
  }

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>資金推移グラフ</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="h-80 flex items-center justify-center">
            <p className="text-muted-foreground">読み込み中...</p>
          </div>
        </CardContent>
      </Card>
    )
  }

  if (chartData.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>資金推移グラフ</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="h-80 flex items-center justify-center">
            <p className="text-muted-foreground">
              グラフデータがありません。銀行口座や収入源を設定してください。
            </p>
          </div>
        </CardContent>
      </Card>
    )
  }

  // 最小値と最大値を計算してY軸の範囲を適切に設定
  const minBalance = Math.min(...chartData.map(d => d.balance))
  const maxBalance = Math.max(...chartData.map(d => d.balance))
  const range = maxBalance - minBalance
  const padding = range * 0.1 // 10%のパディング
  
  return (
    <Card>
      <CardHeader>
        <CardTitle>資金推移グラフ</CardTitle>
        <p className="text-sm text-muted-foreground">
          前後1年間の月末残高推移
        </p>
      </CardHeader>
      <CardContent>
        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={chartData} margin={{ top: 20, right: 30, left: 20, bottom: 20 }}>
              <defs>
                <linearGradient id="balanceGradient" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stopColor={chartColors.primary} stopOpacity={0.8} />
                  <stop offset="50%" stopColor={chartColors.primary} stopOpacity={0.4} />
                  <stop offset="100%" stopColor={chartColors.primary} stopOpacity={0.05} />
                </linearGradient>
              </defs>
              <CartesianGrid 
                strokeDasharray="3 3" 
                stroke={chartColors.muted} 
                opacity={0.3}
              />
              <XAxis 
                dataKey="formattedDate"
                tick={{ 
                  fontSize: 12, 
                  fill: chartColors.muted,
                  fontWeight: 500
                }}
                axisLine={false}
                tickLine={false}
                interval="preserveStartEnd"
                angle={-45}
                textAnchor="end"
                height={80}
              />
              <YAxis 
                tick={{ 
                  fontSize: 12, 
                  fill: chartColors.muted,
                  fontWeight: 500
                }}
                axisLine={false}
                tickLine={false}
                tickFormatter={(value) => `¥${(value / 100).toLocaleString()}`}
                domain={[
                  Math.floor((minBalance - padding) / 10000) * 10000,
                  Math.ceil((maxBalance + padding) / 10000) * 10000
                ]}
              />
              <Tooltip content={<CustomTooltip />} />
              <Area 
                type="monotone" 
                dataKey="balance" 
                stroke={chartColors.primary}
                strokeWidth={3}
                fill="url(#balanceGradient)"
                dot={false}
                activeDot={{ 
                  r: 6, 
                  stroke: chartColors.primary, 
                  strokeWidth: 3,
                  fill: theme === 'dark' ? '#1f2937' : '#ffffff',
                  strokeOpacity: 1
                }}
              />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      </CardContent>
    </Card>
  )
}
