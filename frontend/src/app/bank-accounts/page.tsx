'use client';

import React, { useState, useEffect } from 'react';
import { Plus, Edit, Trash2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { MainLayout } from '@/components/layout/main-layout';
import { useApi } from '@/components/providers/api-provider';
import { BankAccount } from '@/types/api';
import { formatCurrency } from '@/lib/utils-extended';
import { BankAccountForm } from '@/components/forms/bank-account-form';
import { ConfirmDialog } from '@/components/common/confirm-dialog';
import { toast } from 'sonner';

export default function BankAccountsPage() {
  const apiClient = useApi();
  const [bankAccounts, setBankAccounts] = useState<BankAccount[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingAccount, setEditingAccount] = useState<BankAccount | null>(null);
  const [deleteConfirm, setDeleteConfirm] = useState<{ open: boolean; account: BankAccount | null }>({
    open: false,
    account: null,
  });

  const loadBankAccounts = React.useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);
      const accounts = await apiClient.getBankAccounts();
      // 安全な初期化 - nullの場合は空配列にする
      setBankAccounts(Array.isArray(accounts) ? accounts : []);
    } catch (error) {
      const errorMessage = '銀行口座の取得に失敗しました';
      setError(errorMessage);
      toast.error(errorMessage);
      console.error('Failed to load bank accounts:', error);
      // エラーが発生した場合も空配列で安全に初期化
      setBankAccounts([]);
    } finally {
      setIsLoading(false);
    }
  }, [apiClient]);

  useEffect(() => {
    loadBankAccounts();
  }, [loadBankAccounts]);

  const handleCreateAccount = () => {
    setEditingAccount(null);
    setIsFormOpen(true);
  };

  const handleEditAccount = (account: BankAccount) => {
    setEditingAccount(account);
    setIsFormOpen(true);
  };

  const handleDeleteAccount = (account: BankAccount) => {
    setDeleteConfirm({ open: true, account });
  };

  const confirmDelete = async () => {
    if (!deleteConfirm.account) return;
    
    try {
      await apiClient.deleteBankAccount(deleteConfirm.account.id);
      toast.success('銀行口座を削除しました');
      loadBankAccounts();
    } catch (error) {
      toast.error('銀行口座の削除に失敗しました');
      console.error('Failed to delete bank account:', error);
    }
  };

  const handleFormSuccess = () => {
    setIsFormOpen(false);
    setEditingAccount(null);
    loadBankAccounts();
  };

  const totalBalance = Array.isArray(bankAccounts) ? bankAccounts.reduce((sum, account) => sum + account.balance, 0) : 0;

  if (isLoading) {
    return (
      <MainLayout>
        <div className="flex items-center justify-center h-64">
          <p>読み込み中...</p>
        </div>
      </MainLayout>
    );
  }

  if (error) {
    return (
      <MainLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-center">
            <p className="text-red-500 mb-4">{error}</p>
            <Button onClick={loadBankAccounts}>再試行</Button>
          </div>
        </div>
      </MainLayout>
    );
  }

  return (
    <MainLayout>
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold">銀行口座</h1>
            <p className="text-muted-foreground">
              銀行口座を管理し、残高を確認できます。
            </p>
          </div>
          <Button onClick={handleCreateAccount}>
            <Plus className="h-4 w-4 mr-2" />
            口座を追加
          </Button>
        </div>

        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <Card className="md:col-span-2 lg:col-span-3">
            <CardHeader>
              <CardTitle>総残高</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold">{formatCurrency(totalBalance)}</p>
            </CardContent>
          </Card>

          {Array.isArray(bankAccounts) && bankAccounts.length > 0 ? (
            bankAccounts.map((account) => (
              <Card key={account.id}>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-lg font-medium">{account.name}</CardTitle>
                  <div className="flex gap-2">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleEditAccount(account)}
                    >
                      <Edit className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleDeleteAccount(account)}
                    >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-bold">{formatCurrency(account.balance)}</p>
                <p className="text-xs text-muted-foreground mt-2">
                  作成日: {new Date(account.created_at).toLocaleDateString('ja-JP')}
                </p>
              </CardContent>
            </Card>
          ))
          ) : (
            <Card className="md:col-span-2 lg:col-span-3">
              <CardContent className="flex flex-col items-center justify-center py-12">
                <p className="text-muted-foreground mb-4">
                  まだ銀行口座が登録されていません。
                </p>
                <Button onClick={handleCreateAccount}>
                  <Plus className="h-4 w-4 mr-2" />
                  最初の口座を追加
                </Button>
              </CardContent>
            </Card>
          )}
        </div>
      </div>

      <BankAccountForm
        open={isFormOpen}
        onOpenChange={setIsFormOpen}
        account={editingAccount}
        onSuccess={handleFormSuccess}
      />

      <ConfirmDialog
        open={deleteConfirm.open}
        onOpenChange={(open) => setDeleteConfirm({ open, account: null })}
        title="銀行口座を削除"
        description={`「${deleteConfirm.account?.name}」を削除してもよろしいですか？この操作は取り消せません。`}
        confirmText="削除"
        onConfirm={confirmDelete}
        variant="destructive"
      />
    </MainLayout>
  );
}
