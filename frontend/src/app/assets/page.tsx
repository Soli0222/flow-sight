'use client';

import React, { useState, useEffect } from 'react';
import { Plus, Edit, Trash2, CreditCard, Receipt } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { MainLayout } from '@/components/layout/main-layout';
import { useApi } from '@/components/providers/api-provider';
import { Asset, BankAccount } from '@/types/api';
import { AssetForm } from '@/components/forms/asset-form';
import { ConfirmDialog } from '@/components/common/confirm-dialog';
import { toast } from 'sonner';

export default function AssetsPage() {
  const apiClient = useApi();
  const [assets, setAssets] = useState<Asset[]>([]);
  const [bankAccounts, setBankAccounts] = useState<BankAccount[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingAsset, setEditingAsset] = useState<Asset | null>(null);
  const [deleteConfirm, setDeleteConfirm] = useState<{ open: boolean; asset: Asset | null }>({
    open: false,
    asset: null,
  });

  const loadData = React.useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);
      const [assetsData, bankAccountsData] = await Promise.all([
        apiClient.getAssets(),
        apiClient.getBankAccounts(),
      ]);
      // 安全な初期化 - nullの場合は空配列にする
      setAssets(Array.isArray(assetsData) ? assetsData : []);
      setBankAccounts(Array.isArray(bankAccountsData) ? bankAccountsData : []);
    } catch (error) {
      const errorMessage = 'データの取得に失敗しました';
      setError(errorMessage);
      toast.error(errorMessage);
      console.error('Failed to load data:', error);
      // エラーが発生した場合も空配列で安全に初期化
      setAssets([]);
      setBankAccounts([]);
    } finally {
      setIsLoading(false);
    }
  }, [apiClient]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const handleCreateAsset = () => {
    setEditingAsset(null);
    setIsFormOpen(true);
  };

  const handleEditAsset = (asset: Asset) => {
    setEditingAsset(asset);
    setIsFormOpen(true);
  };

  const handleDeleteAsset = (asset: Asset) => {
    setDeleteConfirm({ open: true, asset });
  };

  const confirmDelete = async () => {
    if (!deleteConfirm.asset) return;
    
    try {
      await apiClient.deleteAsset(deleteConfirm.asset.id);
      toast.success('資産を削除しました');
      loadData();
    } catch (error) {
      toast.error('資産の削除に失敗しました');
      console.error('Failed to delete asset:', error);
    }
  };

  const handleFormSuccess = () => {
    setIsFormOpen(false);
    setEditingAsset(null);
    loadData();
  };

  const getBankAccountName = (accountId: string) => {
    const account = bankAccounts.find(acc => acc.id === accountId);
    return account?.name || accountId;
  };

  const getAssetIcon = (assetType: string) => {
    return assetType === 'card' ? <CreditCard className="h-5 w-5" /> : <Receipt className="h-5 w-5" />;
  };

  const getAssetTypeText = (assetType: string) => {
    return assetType === 'card' ? 'クレジットカード' : 'ローン';
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
            <h1 className="text-3xl font-bold">資産管理</h1>
            <p className="text-muted-foreground">
              クレジットカードやローンを管理できます。
            </p>
          </div>
          <Button onClick={handleCreateAsset}>
            <Plus className="h-4 w-4 mr-2" />
            資産を追加
          </Button>
        </div>

        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {Array.isArray(assets) && assets.length > 0 ? (
            assets.map((asset) => (
              <Card key={asset.id}>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <div className="flex items-center gap-2">
                    {getAssetIcon(asset.asset_type)}
                    <div>
                      <CardTitle className="text-lg font-medium">{asset.name}</CardTitle>
                      <Badge variant="secondary" className="mt-1">
                        {getAssetTypeText(asset.asset_type)}
                      </Badge>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleEditAsset(asset)}
                    >
                      <Edit className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleDeleteAsset(asset)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <div>
                    <p className="text-sm text-muted-foreground">銀行口座</p>
                    <p className="font-medium">{getBankAccountName(asset.bank_account)}</p>
                  </div>
                  
                  {asset.closing_day && (
                    <div>
                      <p className="text-sm text-muted-foreground">締め日</p>
                      <p className="font-medium">{asset.closing_day}日</p>
                    </div>
                  )}
                  
                  <div>
                    <p className="text-sm text-muted-foreground">支払日</p>
                    <p className="font-medium">{asset.payment_day}日</p>
                  </div>
                  
                  <p className="text-xs text-muted-foreground mt-2">
                    作成日: {new Date(asset.created_at).toLocaleDateString('ja-JP')}
                  </p>
                </div>
              </CardContent>
            </Card>
          ))
          ) : (
            <Card className="md:col-span-2 lg:col-span-3">
              <CardContent className="flex flex-col items-center justify-center py-12">
                <p className="text-muted-foreground mb-4">
                  まだ資産が登録されていません。
                </p>
                <Button onClick={handleCreateAsset}>
                  <Plus className="h-4 w-4 mr-2" />
                  最初の資産を追加
                </Button>
              </CardContent>
            </Card>
          )}
        </div>
      </div>

      <AssetForm
        open={isFormOpen}
        onOpenChange={setIsFormOpen}
        asset={editingAsset}
        bankAccounts={bankAccounts}
        onSuccess={handleFormSuccess}
      />

      <ConfirmDialog
        open={deleteConfirm.open}
        onOpenChange={(open) => setDeleteConfirm({ open, asset: null })}
        title="資産を削除"
        description={`「${deleteConfirm.asset?.name}」を削除してもよろしいですか？この操作は取り消せません。`}
        confirmText="削除"
        onConfirm={confirmDelete}
        variant="destructive"
      />
    </MainLayout>
  );
}
