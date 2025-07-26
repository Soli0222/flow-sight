'use client';

import React, { useState, useEffect } from 'react';
import { Plus, Edit, Trash2, Calendar, Check, X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { MainLayout } from '@/components/layout/main-layout';
import { useApi } from '@/components/providers/api-provider';
import { CreditCard, CardMonthlyTotal } from '@/types/api';
import { CardMonthlyTotalForm } from '@/components/forms/card-monthly-total-form';
import { ConfirmDialog } from '@/components/common/confirm-dialog';
import { toast } from 'sonner';
import { formatCurrency } from '@/lib/utils-extended';

export default function CardMonthlyTotalsPage() {
  const apiClient = useApi();
  const [creditCards, setCreditCards] = useState<CreditCard[]>([]);
  const [cardMonthlyTotals, setCardMonthlyTotals] = useState<CardMonthlyTotal[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingTotal, setEditingTotal] = useState<CardMonthlyTotal | null>(null);
  const [deleteConfirm, setDeleteConfirm] = useState<{ open: boolean; total: CardMonthlyTotal | null }>({
    open: false,
    total: null,
  });

  const loadData = React.useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);
      
      // まずクレジットカードデータを取得
      const creditCardsData = await apiClient.getCreditCards();
      const creditCardsList = Array.isArray(creditCardsData) ? creditCardsData : [];
      setCreditCards(creditCardsList);

      // クレジットカードがある場合のみ月次利用額データを取得
      if (creditCardsList.length > 0) {
        const allTotals: CardMonthlyTotal[] = [];
        for (const creditCard of creditCardsList) {
          try {
            const totals = await apiClient.getCardMonthlyTotals(creditCard.id);
            if (Array.isArray(totals)) {
              allTotals.push(...totals);
            }
          } catch (error) {
            console.warn(`Failed to load totals for credit card ${creditCard.id}:`, error);
          }
        }
        setCardMonthlyTotals(allTotals);
      } else {
        setCardMonthlyTotals([]);
      }
    } catch (error) {
      const errorMessage = 'データの取得に失敗しました';
      setError(errorMessage);
      toast.error(errorMessage);
      console.error('Failed to load data:', error);
      setCreditCards([]);
      setCardMonthlyTotals([]);
    } finally {
      setIsLoading(false);
    }
  }, [apiClient]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const handleCreateTotal = () => {
    setEditingTotal(null);
    setIsFormOpen(true);
  };

  const handleEditTotal = (total: CardMonthlyTotal) => {
    setEditingTotal(total);
    setIsFormOpen(true);
  };

  const handleDeleteTotal = (total: CardMonthlyTotal) => {
    setDeleteConfirm({ open: true, total });
  };

  const confirmDelete = async () => {
    if (!deleteConfirm.total) return;

    try {
      await apiClient.deleteCardMonthlyTotal(deleteConfirm.total.id);
      toast.success('月次利用額を削除しました');
      loadData();
    } catch (error) {
      toast.error('月次利用額の削除に失敗しました');
      console.error('Failed to delete card monthly total:', error);
    } finally {
      setDeleteConfirm({ open: false, total: null });
    }
  };

  const handleFormSuccess = () => {
    setIsFormOpen(false);
    setEditingTotal(null);
    loadData();
  };

  const toggleConfirmation = async (total: CardMonthlyTotal) => {
    try {
      await apiClient.updateCardMonthlyTotal(total.id, {
        ...total,
        is_confirmed: !total.is_confirmed,
      });
      toast.success(`確認状態を${!total.is_confirmed ? '確認済み' : '未確認'}に変更しました`);
      loadData();
    } catch (error) {
      toast.error('確認状態の更新に失敗しました');
      console.error('Failed to update confirmation status:', error);
    }
  };

  const getAssetName = (creditCardId: string) => {
    const creditCard = creditCards.find(c => c.id === creditCardId);
    return creditCard?.name || creditCardId;
  };  const sortedTotals = cardMonthlyTotals.sort((a, b) => 
    b.year_month.localeCompare(a.year_month)
  );

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
            <h1 className="text-3xl font-bold">カード月次利用額</h1>
            <p className="text-muted-foreground">
              クレジットカードの月次利用額を管理できます。
            </p>
          </div>
          <Button 
            onClick={handleCreateTotal} 
            disabled={creditCards.length === 0}
          >
            <Plus className="h-4 w-4 mr-2" />
            月次利用額を追加
          </Button>
        </div>

        {error && (
          <Card className="border-destructive">
            <CardContent className="pt-6">
              <p className="text-destructive">{error}</p>
            </CardContent>
          </Card>
        )}

        {creditCards.length === 0 && (
          <Card>
            <CardContent className="pt-6">
              <div className="text-center py-8">
                <Calendar className="h-12 w-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-semibold mb-2">カード資産が登録されていません</h3>
                <p className="text-muted-foreground mb-4">
                  カード月次利用額を管理するには、まずクレジットカードをクレジットカード管理ページで登録する必要があります。
                </p>
                <Button asChild>
                  <a href="/credit-cards">クレジットカード管理ページへ</a>
                </Button>
              </div>
            </CardContent>
          </Card>
        )}

        {creditCards.length > 0 && (
          <>
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
              {sortedTotals.map((total) => (
                <Card key={total.id}>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <div className="flex items-center gap-2">
                      <Calendar className="h-5 w-5" />
                      <div>
                        <CardTitle className="text-lg font-medium">
                          {getAssetName(total.credit_card_id)}
                        </CardTitle>
                        <p className="text-sm text-muted-foreground">
                          {total.year_month}
                        </p>
                      </div>
                    </div>
                    <div className="flex gap-2">
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => handleEditTotal(total)}
                      >
                        <Edit className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => handleDeleteTotal(total)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      <div>
                        <p className="text-2xl font-bold">
                          {formatCurrency(total.total_amount)}
                        </p>
                      </div>
                      
                      <div className="flex items-center justify-between">
                        <div>
                          <p className="text-sm text-muted-foreground">確認状態</p>
                          <Badge 
                            variant={total.is_confirmed ? "default" : "secondary"}
                            className="mt-1"
                          >
                            {total.is_confirmed ? '確認済み' : '未確認'}
                          </Badge>
                        </div>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => toggleConfirmation(total)}
                        >
                          {total.is_confirmed ? (
                            <X className="h-4 w-4 mr-1" />
                          ) : (
                            <Check className="h-4 w-4 mr-1" />
                          )}
                          {total.is_confirmed ? '未確認にする' : '確認済みにする'}
                        </Button>
                      </div>
                      
                      <p className="text-xs text-muted-foreground">
                        作成日: {new Date(total.created_at).toLocaleDateString('ja-JP')}
                      </p>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>

            {sortedTotals.length === 0 && (
              <Card>
                <CardContent className="pt-6">
                  <div className="text-center py-8">
                    <Calendar className="h-12 w-12 mx-auto text-muted-foreground mb-4" />
                    <h3 className="text-lg font-semibold mb-2">月次利用額が登録されていません</h3>
                    <p className="text-muted-foreground mb-4">
                      上の「月次利用額を追加」ボタンから最初の月次利用額を追加してください。
                    </p>
                  </div>
                </CardContent>
              </Card>
            )}
          </>
        )}
      </div>

      {/* フォームダイアログ */}
      {isFormOpen && (
        <CardMonthlyTotalForm
          total={editingTotal}
          creditCards={creditCards}
          onSuccess={handleFormSuccess}
          onCancel={() => setIsFormOpen(false)}
        />
      )}

      {/* 削除確認ダイアログ */}
      <ConfirmDialog
        open={deleteConfirm.open}
        onOpenChange={(open) => !open && setDeleteConfirm({ open: false, total: null })}
        title="月次利用額の削除"
        description={`${deleteConfirm.total ? getAssetName(deleteConfirm.total.credit_card_id) : ''}の${deleteConfirm.total?.year_month}の月次利用額を削除しますか？この操作は取り消せません。`}
        onConfirm={confirmDelete}
      />
    </MainLayout>
  );
}
