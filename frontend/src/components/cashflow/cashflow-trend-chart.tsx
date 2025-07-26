'use client'

import { useMemo, useEffect, useState } from 'react'
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useTheme } from 'next-themes'
import { formatCurrency } from '@/lib/utils-extended'
import { CashflowProjection } from '@/types/api'

interface ChartDataPoint {
  date: string
  balance: number
  income: number
  expense: number
  formattedDate: string
  originalDate: string
}

interface CustomTooltipProps {
  active?: boolean
  payload?: Array<{
    payload: ChartDataPoint
    value: number
    dataKey: string
  }>
  label?: string
}

interface CashflowTrendChartProps {
  projections: CashflowProjection[]
  projectionMonths: number
}

export function CashflowTrendChart({ projections, projectionMonths }: CashflowTrendChartProps) {
  const { theme } = useTheme()
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

  const chartData = useMemo(() => {
    if (projections.length === 0) return []

    // 月次集計用のマップ
    const monthlyData = new Map<string, {
      balance: number
      income: number
      expense: number
      lastDate: string
    }>()

    projections.forEach((projection) => {
      const date = new Date(projection.date)
      const yearMonth = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
      
      const existing = monthlyData.get(yearMonth)
      
      if (!existing) {
        // 新しい月の場合
        monthlyData.set(yearMonth, {
          balance: projection.balance,
          income: projection.income,
          expense: projection.expense,
          lastDate: projection.date
        })
      } else {
        // 既存の月の場合、最後の日の残高を使用し、収入と支出は累積
        if (new Date(projection.date) > new Date(existing.lastDate)) {
          monthlyData.set(yearMonth, {
            balance: projection.balance, // 最後の日の残高
            income: existing.income + projection.income, // 累積収入
            expense: existing.expense + projection.expense, // 累積支出
            lastDate: projection.date
          })
        } else {
          monthlyData.set(yearMonth, {
            ...existing,
            income: existing.income + projection.income,
            expense: existing.expense + projection.expense
          })
        }
      }
    })

    // チャート用のデータ形式に変換
    const chartPoints: ChartDataPoint[] = Array.from(monthlyData.entries())
      .map(([yearMonth, data]) => {
        const [year, month] = yearMonth.split('-')
        return {
          date: yearMonth,
          balance: data.balance,
          income: data.income,
          expense: data.expense,
          formattedDate: `${year}年${parseInt(month)}月`,
          originalDate: data.lastDate
        }
      })
      .sort((a, b) => a.date.localeCompare(b.date))

    return chartPoints
  }, [projections])

  const CustomTooltip = ({ active, payload }: CustomTooltipProps) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload
      return (
        <div className="bg-background/95 backdrop-blur-sm border rounded-lg shadow-lg p-3 border-border">
          <p className="font-medium text-foreground mb-2">{data.formattedDate}</p>
          <div className="space-y-1">
            <p className="text-sm text-green-600">
              収入: {formatCurrency(data.income)}
            </p>
            <p className="text-sm text-red-600">
              支出: {formatCurrency(data.expense)}
            </p>
            <p className={`text-sm font-medium ${data.balance >= 0 ? 'text-green-600' : 'text-red-600'}`}>
              残高: {formatCurrency(data.balance)}
            </p>
          </div>
        </div>
      )
    }
    return null
  }

  if (chartData.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>キャッシュフロー推移グラフ</CardTitle>
          <p className="text-sm text-muted-foreground">
            {projectionMonths >= 12 
              ? `当月から${Math.floor(projectionMonths / 12)}年${projectionMonths % 12 > 0 ? `${projectionMonths % 12}ヶ月` : ''}間の推移`
              : `当月から${projectionMonths}ヶ月間の推移`
            }
          </p>
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
  const padding = range * 0.1 || 10000 // 最小パディング
  
  return (
    <Card>
      <CardHeader>
        <CardTitle>キャッシュフロー推移グラフ</CardTitle>
        <p className="text-sm text-muted-foreground">
          {projectionMonths >= 12 
            ? `当月から${Math.floor(projectionMonths / 12)}年${projectionMonths % 12 > 0 ? `${projectionMonths % 12}ヶ月` : ''}間の月末残高推移`
            : `当月から${projectionMonths}ヶ月間の月末残高推移`
          }
        </p>
      </CardHeader>
      <CardContent>
        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={chartData} margin={{ top: 20, right: 30, left: 20, bottom: 60 }}>
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
