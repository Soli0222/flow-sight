'use client';

import React, { useState, useEffect } from 'react';
import { Plus, Edit, Trash2, Repeat } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { MainLayout } from '@/components/layout/main-layout';
import { useApi } from '@/components/providers/api-provider';
import { RecurringPayment, BankAccount } from '@/types/api';
import { formatCurrency } from '@/lib/utils-extended';
import { RecurringPaymentForm } from '@/components/forms/recurring-payment-form';
import { ConfirmDialog } from '@/components/common/confirm-dialog';
import { toast } from 'sonner';

export default function RecurringPaymentsPage() {
  const apiClient = useApi();
  const [recurringPayments, setRecurringPayments] = useState<RecurringPayment[]>([]);
  const [bankAccounts, setBankAccounts] = useState<BankAccount[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingPayment, setEditingPayment] = useState<RecurringPayment | null>(null);
  const [deleteConfirm, setDeleteConfirm] = useState<{ open: boolean; payment: RecurringPayment | null }>({
    open: false,
    payment: null,
  });

  const loadData = React.useCallback(async () => {
    try {
      setIsLoading(true);
      const [paymentsData, bankAccountsData] = await Promise.all([
        apiClient.getRecurringPayments(),
        apiClient.getBankAccounts(),
      ]);
      setRecurringPayments(paymentsData);
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

  const handleCreatePayment = () => {
    setEditingPayment(null);
    setIsFormOpen(true);
  };

  const handleEditPayment = (payment: RecurringPayment) => {
    setEditingPayment(payment);
    setIsFormOpen(true);
  };

  const handleDeletePayment = (payment: RecurringPayment) => {
    setDeleteConfirm({ open: true, payment });
  };

  const confirmDelete = async () => {
    if (!deleteConfirm.payment) return;
    
    try {
      await apiClient.deleteRecurringPayment(deleteConfirm.payment.id);
      toast.success('定期支払いを削除しました');
      loadData();
    } catch (error) {
      toast.error('定期支払いの削除に失敗しました');
      console.error('Failed to delete recurring payment:', error);
    }
  };

  const handleFormSuccess = () => {
    setIsFormOpen(false);
    setEditingPayment(null);
    loadData();
  };

  const getBankAccountName = (accountId: string) => {
    const account = bankAccounts.find(acc => acc.id === accountId);
    return account?.name || accountId;
  };

  const totalMonthlyPayments = recurringPayments
    .filter(payment => payment.is_active)
    .reduce((sum, payment) => sum + payment.amount, 0);

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
            <h1 className="text-3xl font-bold">定期支払い</h1>
            <p className="text-muted-foreground">
              月次の定期支払いやローンを管理できます。
            </p>
          </div>
          <Button onClick={handleCreatePayment}>
            <Plus className="h-4 w-4 mr-2" />
            定期支払いを追加
          </Button>
        </div>

        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <Card className="md:col-span-2 lg:col-span-3">
            <CardHeader>
              <CardTitle>月次支払い合計</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold">{formatCurrency(totalMonthlyPayments)}</p>
              <p className="text-sm text-muted-foreground mt-1">
                アクティブな定期支払いの合計
              </p>
            </CardContent>
          </Card>

          {recurringPayments.map((payment) => (
            <Card key={payment.id}>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <div className="flex items-center gap-2">
                  <Repeat className="h-5 w-5" />
                  <div>
                    <CardTitle className="text-lg font-medium">{payment.name}</CardTitle>
                    <Badge variant={payment.is_active ? "default" : "outline"} className="mt-1">
                      {payment.is_active ? 'アクティブ' : '非アクティブ'}
                    </Badge>
                  </div>
                </div>
                <div className="flex gap-2">
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleEditPayment(payment)}
                  >
                    <Edit className="h-4 w-4" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleDeletePayment(payment)}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <div>
                    <p className="text-sm text-muted-foreground">金額</p>
                    <p className="text-xl font-bold">{formatCurrency(payment.amount)}</p>
                  </div>
                  
                  <div>
                    <p className="text-sm text-muted-foreground">支払日</p>
                    <p className="font-medium">{payment.payment_day}日</p>
                  </div>
                  
                  <div>
                    <p className="text-sm text-muted-foreground">支払口座</p>
                    <p className="font-medium">{getBankAccountName(payment.bank_account)}</p>
                  </div>
                  
                  <div>
                    <p className="text-sm text-muted-foreground">開始年月</p>
                    <p className="font-medium">{payment.start_year_month}</p>
                  </div>
                  
                  {payment.total_payments && (
                    <div>
                      <p className="text-sm text-muted-foreground">残り回数</p>
                      <p className="font-medium">
                        {payment.remaining_payments || 0} / {payment.total_payments} 回
                      </p>
                    </div>
                  )}
                  
                  {payment.note && (
                    <div>
                      <p className="text-sm text-muted-foreground">備考</p>
                      <p className="text-sm">{payment.note}</p>
                    </div>
                  )}
                  
                  <p className="text-xs text-muted-foreground mt-2">
                    作成日: {new Date(payment.created_at).toLocaleDateString('ja-JP')}
                  </p>
                </div>
              </CardContent>
            </Card>
          ))}

          {recurringPayments.length === 0 && (
            <Card className="md:col-span-2 lg:col-span-3">
              <CardContent className="flex flex-col items-center justify-center py-12">
                <p className="text-muted-foreground mb-4">
                  まだ定期支払いが登録されていません。
                </p>
                <Button onClick={handleCreatePayment}>
                  <Plus className="h-4 w-4 mr-2" />
                  最初の定期支払いを追加
                </Button>
              </CardContent>
            </Card>
          )}
        </div>
      </div>

      <RecurringPaymentForm
        open={isFormOpen}
        onOpenChange={setIsFormOpen}
        payment={editingPayment}
        bankAccounts={bankAccounts}
        onSuccess={handleFormSuccess}
      />

      <ConfirmDialog
        open={deleteConfirm.open}
        onOpenChange={(open) => setDeleteConfirm({ open, payment: null })}
        title="定期支払いを削除"
        description={`「${deleteConfirm.payment?.name}」を削除してもよろしいですか？この操作は取り消せません。`}
        confirmText="削除"
        onConfirm={confirmDelete}
        variant="destructive"
      />
    </MainLayout>
  );
}
