'use client';

import React, { useState, useEffect } from 'react';
import { Download, TrendingUp, TrendingDown, DollarSign } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { MainLayout } from '@/components/layout/main-layout';
import { useApi } from '@/components/providers/api-provider';
import { CashflowProjection } from '@/types/api';
import { formatCurrency, formatDate } from '@/lib/utils-extended';
import { toast } from 'sonner';

export default function CashflowPage() {
  const apiClient = useApi();
  const [projections, setProjections] = useState<CashflowProjection[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [projectionMonths, setProjectionMonths] = useState(6);

  const loadProjections = React.useCallback(async () => {
    try {
      setIsLoading(true);
      const data = await apiClient.getCashflowProjection(projectionMonths);
      setProjections(data);
    } catch (error) {
      toast.error('キャッシュフロー予測の取得に失敗しました');
      console.error('Failed to load cashflow projections:', error);
    } finally {
      setIsLoading(false);
    }
  }, [apiClient, projectionMonths]);

  useEffect(() => {
    loadProjections();
  }, [loadProjections]);

  const handleExportCSV = () => {
    if (projections.length === 0) {
      toast.error('エクスポートするデータがありません');
      return;
    }

    const csvContent = [
      ['日付', '収入', '支出', '残高'],
      ...projections.map(p => [
        p.date,
        (p.income / 100).toString(),
        (p.expense / 100).toString(),
        (p.balance / 100).toString(),
      ]),
    ]
      .map(row => row.join(','))
      .join('\n');

    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `cashflow_projection_${new Date().toISOString().split('T')[0]}.csv`;
    link.click();
    
    toast.success('CSVファイルをダウンロードしました');
  };

  const totalProjectedIncome = projections.reduce((sum, p) => sum + p.income, 0);
  const totalProjectedExpense = projections.reduce((sum, p) => sum + p.expense, 0);
  const netCashflow = totalProjectedIncome - totalProjectedExpense;

  if (isLoading) {
    return (
      <MainLayout>
        <div className="flex items-center justify-center h-64">
          <p>読み込み中...</p>
        </div>
      </MainLayout>
    );
  }

  return (
    <MainLayout>
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold">キャッシュフロー予測</h1>
            <p className="text-muted-foreground">
              将来の収支を予測し、資金計画を立てることができます。
            </p>
          </div>
          <div className="flex gap-2">
            <select
              value={projectionMonths}
              onChange={(e) => setProjectionMonths(Number(e.target.value))}
              className="px-3 py-2 border border-input rounded-md"
            >
              <option value={3}>3ヶ月</option>
              <option value={6}>6ヶ月</option>
              <option value={12}>12ヶ月</option>
            </select>
            <Button onClick={handleExportCSV}>
              <Download className="h-4 w-4 mr-2" />
              CSV出力
            </Button>
          </div>
        </div>

        <div className="grid gap-6 md:grid-cols-3">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">予測収入合計</CardTitle>
              <TrendingUp className="h-4 w-4 text-green-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-green-600">
                {formatCurrency(totalProjectedIncome)}
              </div>
              <p className="text-xs text-muted-foreground">
                {projectionMonths}ヶ月間の合計
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">予測支出合計</CardTitle>
              <TrendingDown className="h-4 w-4 text-red-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-red-600">
                {formatCurrency(totalProjectedExpense)}
              </div>
              <p className="text-xs text-muted-foreground">
                {projectionMonths}ヶ月間の合計
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">純キャッシュフロー</CardTitle>
              <DollarSign className={`h-4 w-4 ${netCashflow >= 0 ? 'text-green-600' : 'text-red-600'}`} />
            </CardHeader>
            <CardContent>
              <div className={`text-2xl font-bold ${netCashflow >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                {formatCurrency(netCashflow)}
              </div>
              <p className="text-xs text-muted-foreground">
                収入 - 支出
              </p>
            </CardContent>
          </Card>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>月別キャッシュフロー予測</CardTitle>
          </CardHeader>
          <CardContent>
            {projections.length === 0 ? (
              <p className="text-muted-foreground text-center py-8">
                キャッシュフロー予測データがありません。<br />
                銀行口座、収入源、資産を設定してください。
              </p>
            ) : (
              <div className="space-y-4">
                {projections.map((projection, index) => (
                  <div
                    key={index}
                    className="flex items-center justify-between p-4 border rounded-lg"
                  >
                    <div className="flex-1">
                      <div className="font-medium">
                        {formatDate(projection.date)}
                      </div>
                      <div className="text-sm text-muted-foreground">
                        収入: {formatCurrency(projection.income)} | 
                        支出: {formatCurrency(projection.expense)}
                      </div>
                      {projection.details.length > 0 && (
                        <div className="mt-2 space-y-1">
                          {projection.details.map((detail, detailIndex) => (
                            <div key={detailIndex} className="text-xs text-muted-foreground">
                              • {detail.description}: {formatCurrency(detail.amount)}
                            </div>
                          ))}
                        </div>
                      )}
                    </div>
                    <div className="text-right">
                      <div className={`text-lg font-bold ${projection.balance >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                        {formatCurrency(projection.balance)}
                      </div>
                      <div className="text-sm text-muted-foreground">残高</div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </MainLayout>
  );
}
