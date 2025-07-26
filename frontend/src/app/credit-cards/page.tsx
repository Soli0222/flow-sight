'use client';

import React, { useState, useEffect } from 'react';
import { Plus, Edit, Trash2, CreditCard } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { MainLayout } from '@/components/layout/main-layout';
import { useApi } from '@/components/providers/api-provider';
import { CreditCard as CreditCardType, BankAccount } from '@/types/api';
import { CreditCardForm } from '@/components/forms/credit-card-form';
import { ConfirmDialog } from '@/components/common/confirm-dialog';
import { toast } from 'sonner';

export default function CreditCardsPage() {
  const apiClient = useApi();
  const [creditCards, setCreditCards] = useState<CreditCardType[]>([]);
  const [bankAccounts, setBankAccounts] = useState<BankAccount[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingCreditCard, setEditingCreditCard] = useState<CreditCardType | null>(null);
  const [deleteConfirm, setDeleteConfirm] = useState<{ open: boolean; creditCard: CreditCardType | null }>({
    open: false,
    creditCard: null,
  });

  const loadData = React.useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);
      const [creditCardsData, bankAccountsData] = await Promise.all([
        apiClient.getCreditCards(),
        apiClient.getBankAccounts(),
      ]);
      setCreditCards(Array.isArray(creditCardsData) ? creditCardsData : []);
      setBankAccounts(Array.isArray(bankAccountsData) ? bankAccountsData : []);
    } catch (error) {
      const errorMessage = 'データの取得に失敗しました';
      setError(errorMessage);
      toast.error(errorMessage);
      console.error('Failed to load data:', error);
      // エラーが発生した場合も空配列で安全に初期化
      setCreditCards([]);
      setBankAccounts([]);
    } finally {
      setIsLoading(false);
    }
  }, [apiClient]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const handleCreateCreditCard = () => {
    setEditingCreditCard(null);
    setIsFormOpen(true);
  };

  const handleEditCreditCard = (creditCard: CreditCardType) => {
    setEditingCreditCard(creditCard);
    setIsFormOpen(true);
  };

  const handleDeleteCreditCard = (creditCard: CreditCardType) => {
    setDeleteConfirm({ open: true, creditCard });
  };

  const confirmDelete = async () => {
    if (!deleteConfirm.creditCard) return;
    
    try {
      await apiClient.deleteCreditCard(deleteConfirm.creditCard.id);
      toast.success('クレジットカードを削除しました');
      loadData();
    } catch (error) {
      toast.error('クレジットカードの削除に失敗しました');
      console.error('Failed to delete credit card:', error);
    }
  };

  const handleFormSuccess = () => {
    setIsFormOpen(false);
    setEditingCreditCard(null);
    loadData();
  };

  const getBankAccountName = (accountId: string) => {
    const account = bankAccounts.find(acc => acc.id === accountId);
    return account?.name || accountId;
  };

  const getAssetIcon = () => {
    return <CreditCard className="h-5 w-5" />;
  };

  const getAssetTypeText = () => {
    return 'クレジットカード';
  };

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
            <Button onClick={loadData}>再試行</Button>
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
            <h1 className="text-3xl font-bold">クレジットカード管理</h1>
            <p className="text-muted-foreground">
              クレジットカードを管理できます。
            </p>
          </div>
          <Button onClick={handleCreateCreditCard}>
            <Plus className="h-4 w-4 mr-2" />
            クレジットカードを追加
          </Button>
        </div>

        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {Array.isArray(creditCards) && creditCards.length > 0 ? (
            creditCards.map((creditCard) => (
              <Card key={creditCard.id}>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <div className="flex items-center gap-2">
                    {getAssetIcon()}
                    <div>
                      <CardTitle className="text-lg font-medium">{creditCard.name}</CardTitle>
                      <Badge variant="secondary" className="mt-1">
                        {getAssetTypeText()}
                      </Badge>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleEditCreditCard(creditCard)}
                    >
                      <Edit className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleDeleteCreditCard(creditCard)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <div>
                    <p className="text-sm text-muted-foreground">銀行口座</p>
                    <p className="font-medium">{getBankAccountName(creditCard.bank_account)}</p>
                  </div>
                  
                  {creditCard.closing_day && (
                    <div>
                      <p className="text-sm text-muted-foreground">締め日</p>
                      <p className="font-medium">{creditCard.closing_day}日</p>
                    </div>
                  )}
                  
                  <div>
                    <p className="text-sm text-muted-foreground">支払日</p>
                    <p className="font-medium">{creditCard.payment_day}日</p>
                  </div>
                  
                  <p className="text-xs text-muted-foreground mt-2">
                    作成日: {new Date(creditCard.created_at).toLocaleDateString('ja-JP')}
                  </p>
                </div>
              </CardContent>
            </Card>
          ))
          ) : (
            <Card className="md:col-span-2 lg:col-span-3">
              <CardContent className="flex flex-col items-center justify-center py-12">
                <p className="text-muted-foreground mb-4">
                  まだクレジットカードが登録されていません。
                </p>
                <Button onClick={handleCreateCreditCard}>
                  <Plus className="h-4 w-4 mr-2" />
                  最初のクレジットカードを追加
                </Button>
              </CardContent>
            </Card>
          )}
        </div>
      </div>

      <CreditCardForm
        open={isFormOpen}
        onOpenChange={setIsFormOpen}
        creditCard={editingCreditCard}
        bankAccounts={bankAccounts}
        onSuccess={handleFormSuccess}
      />

      <ConfirmDialog
        open={deleteConfirm.open}
        onOpenChange={(open) => setDeleteConfirm({ open, creditCard: null })}
        title="クレジットカードを削除"
        description={`「${deleteConfirm.creditCard?.name}」を削除してもよろしいですか？この操作は取り消せません。`}
        confirmText="削除"
        onConfirm={confirmDelete}
        variant="destructive"
      />
    </MainLayout>
  );
}
