'use client';

import React, { useState, useEffect } from 'react';
import { Plus, Edit, Trash2, DollarSign, Calendar } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { MainLayout } from '@/components/layout/main-layout';
import { useApi } from '@/components/providers/api-provider';
import { IncomeSource, BankAccount } from '@/types/api';
import { formatCurrency } from '@/lib/utils-extended';
import { IncomeForm } from '@/components/forms/income-form';
import { ConfirmDialog } from '@/components/common/confirm-dialog';
import { toast } from 'sonner';

export default function IncomePage() {
  const apiClient = useApi();
  const [incomeSources, setIncomeSources] = useState<IncomeSource[]>([]);
  const [bankAccounts, setBankAccounts] = useState<BankAccount[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingIncome, setEditingIncome] = useState<IncomeSource | null>(null);
  const [deleteConfirm, setDeleteConfirm] = useState<{ open: boolean; income: IncomeSource | null }>({
    open: false,
    income: null,
  });

  const loadData = React.useCallback(async () => {
    try {
      setIsLoading(true);
      const [incomeData, bankAccountsData] = await Promise.all([
        apiClient.getIncomeSources(),
        apiClient.getBankAccounts(),
      ]);
      setIncomeSources(incomeData);
      setBankAccounts(bankAccountsData);
    } catch (error) {
      toast.error('データの取得に失敗しました');
      console.error('Failed to load data:', error);
    } finally {
      setIsLoading(false);
    }
  }, [apiClient]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const handleCreateIncome = () => {
    setEditingIncome(null);
    setIsFormOpen(true);
  };

  const handleEditIncome = (income: IncomeSource) => {
    setEditingIncome(income);
    setIsFormOpen(true);
  };

  const handleDeleteIncome = (income: IncomeSource) => {
    setDeleteConfirm({ open: true, income });
  };

  const confirmDelete = async () => {
    if (!deleteConfirm.income) return;
    
    try {
      await apiClient.deleteIncomeSource(deleteConfirm.income.id);
      toast.success('収入源を削除しました');
      loadData();
    } catch (error) {
      toast.error('収入源の削除に失敗しました');
      console.error('Failed to delete income source:', error);
    }
  };

  const handleFormSuccess = () => {
    setIsFormOpen(false);
    setEditingIncome(null);
    loadData();
  };

  const getBankAccountName = (accountId: string) => {
    const account = bankAccounts.find(acc => acc.id === accountId);
    return account?.name || accountId;
  };

  const getIncomeTypeText = (incomeType: string) => {
    return incomeType === 'monthly_fixed' ? '月額固定' : '一時的';
  };

  const getIncomeIcon = (incomeType: string) => {
    return incomeType === 'monthly_fixed' ? <DollarSign className="h-5 w-5" /> : <Calendar className="h-5 w-5" />;
  };

  const totalMonthlyIncome = incomeSources
    .filter(income => income.income_type === 'monthly_fixed' && income.is_active)
    .reduce((sum, income) => sum + income.base_amount, 0);

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
            <h1 className="text-3xl font-bold">収入管理</h1>
            <p className="text-muted-foreground">
              収入源を管理し、月額収入を確認できます。
            </p>
          </div>
          <Button onClick={handleCreateIncome}>
            <Plus className="h-4 w-4 mr-2" />
            収入源を追加
          </Button>
        </div>

        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <Card className="md:col-span-2 lg:col-span-3">
            <CardHeader>
              <CardTitle>月額収入合計</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold">{formatCurrency(totalMonthlyIncome)}</p>
              <p className="text-sm text-muted-foreground mt-1">
                アクティブな月額固定収入の合計
              </p>
            </CardContent>
          </Card>

          {incomeSources.map((income) => (
            <Card key={income.id}>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <div className="flex items-center gap-2">
                  {getIncomeIcon(income.income_type)}
                  <div>
                    <CardTitle className="text-lg font-medium">{income.name}</CardTitle>
                    <div className="flex gap-2 mt-1">
                      <Badge variant="secondary">
                        {getIncomeTypeText(income.income_type)}
                      </Badge>
                      <Badge variant={income.is_active ? "default" : "outline"}>
                        {income.is_active ? 'アクティブ' : '非アクティブ'}
                      </Badge>
                    </div>
                  </div>
                </div>
                <div className="flex gap-2">
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleEditIncome(income)}
                  >
                    <Edit className="h-4 w-4" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleDeleteIncome(income)}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <div>
                    <p className="text-sm text-muted-foreground">金額</p>
                    <p className="text-xl font-bold">{formatCurrency(income.base_amount)}</p>
                  </div>
                  
                  <div>
                    <p className="text-sm text-muted-foreground">振込先</p>
                    <p className="font-medium">{getBankAccountName(income.bank_account)}</p>
                  </div>
                  
                  {income.scheduled_year_month && (
                    <div>
                      <p className="text-sm text-muted-foreground">予定年月</p>
                      <p className="font-medium">{income.scheduled_year_month}</p>
                    </div>
                  )}
                  
                  <p className="text-xs text-muted-foreground mt-2">
                    作成日: {new Date(income.created_at).toLocaleDateString('ja-JP')}
                  </p>
                </div>
              </CardContent>
            </Card>
          ))}

          {incomeSources.length === 0 && (
            <Card className="md:col-span-2 lg:col-span-3">
              <CardContent className="flex flex-col items-center justify-center py-12">
                <p className="text-muted-foreground mb-4">
                  まだ収入源が登録されていません。
                </p>
                <Button onClick={handleCreateIncome}>
                  <Plus className="h-4 w-4 mr-2" />
                  最初の収入源を追加
                </Button>
              </CardContent>
            </Card>
          )}
        </div>
      </div>

      <IncomeForm
        open={isFormOpen}
        onOpenChange={setIsFormOpen}
        income={editingIncome}
        bankAccounts={bankAccounts}
        onSuccess={handleFormSuccess}
      />

      <ConfirmDialog
        open={deleteConfirm.open}
        onOpenChange={(open) => setDeleteConfirm({ open, income: null })}
        title="収入源を削除"
        description={`「${deleteConfirm.income?.name}」を削除してもよろしいですか？この操作は取り消せません。`}
        confirmText="削除"
        onConfirm={confirmDelete}
        variant="destructive"
      />
    </MainLayout>
  );
}
