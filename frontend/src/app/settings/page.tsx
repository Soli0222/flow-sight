'use client';

import React, { useState, useEffect } from 'react';
import { Save, User, Info } from 'lucide-react';
import Image from 'next/image';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { MainLayout } from '@/components/layout/main-layout';
import { useApi } from '@/components/providers/api-provider';
import { toast } from 'sonner';
import { VersionInfo, UserInfo } from '@/types/api';
import { FRONTEND_VERSION } from '@/lib/version';

export default function SettingsPage() {
  const apiClient = useApi();
  const [settings, setSettings] = useState<Record<string, string>>({
    minimum_monthly_expense: '0',
    notification_enabled: 'true',
    theme: 'light',
  });
  const [versionInfo, setVersionInfo] = useState<VersionInfo | null>(null);
  const [userInfo, setUserInfo] = useState<UserInfo | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);

  const loadSettings = React.useCallback(async () => {
    try {
      setIsLoading(true);
      
      // Load settings
      const settingsData = await apiClient.getSettings();
      
      // データの安全性チェック
      if (!Array.isArray(settingsData)) {
        console.error('Invalid settings data received:', settingsData);
        toast.error('設定データの形式が正しくありません');
        return;
      }
      
      const settingsMap = settingsData.reduce((acc, setting) => {
        if (setting && setting.key && setting.value !== undefined) {
          acc[setting.key] = setting.value;
        }
        return acc;
      }, {} as Record<string, string>);
      
      setSettings({
        minimum_monthly_expense: settingsMap.minimum_monthly_expense 
          ? String(Math.round(Number(settingsMap.minimum_monthly_expense) / 100))
          : '0', // Convert cents to yen for display
        notification_enabled: settingsMap.notification_enabled || 'true',
        theme: settingsMap.theme || 'light',
      });

      // Load version info
      try {
        const version = await apiClient.getVersion();
        setVersionInfo(version);
      } catch (error) {
        console.error('Failed to load version info:', error);
      }

      // Load user info
      try {
        const user = await apiClient.getCurrentUser();
        setUserInfo(user);
      } catch (error) {
        console.error('Failed to load user info:', error);
      }
    } catch (error) {
      toast.error('設定の取得に失敗しました');
      console.error('Failed to load settings:', error);
    } finally {
      setIsLoading(false);
    }
  }, [apiClient]);

  useEffect(() => {
    loadSettings();
  }, [loadSettings]);

  const handleSave = async () => {
    try {
      setIsSaving(true);
      // Convert yen to cents before saving
      const settingsToSave = {
        ...settings,
        minimum_monthly_expense: String(Number(settings.minimum_monthly_expense) * 100)
      };
      await apiClient.updateSettings({ settings: settingsToSave });
      toast.success('設定を保存しました');
    } catch (error) {
      toast.error('設定の保存に失敗しました');
      console.error('Failed to save settings:', error);
    } finally {
      setIsSaving(false);
    }
  };

  const updateSetting = (key: string, value: string) => {
    setSettings(prev => ({ ...prev, [key]: value }));
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

  return (
    <MainLayout>
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold">設定</h1>
          <p className="text-muted-foreground">
            アプリケーションの設定を管理できます。
          </p>
        </div>

        <div className="grid gap-6 md:grid-cols-1 lg:grid-cols-2">
          <Card>
            <CardHeader>
              <CardTitle>金融設定</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="minimum_monthly_expense">最低月支出 (円)</Label>
                <Input
                  id="minimum_monthly_expense"
                  type="number"
                  value={settings.minimum_monthly_expense}
                  onChange={(e) => updateSetting('minimum_monthly_expense', e.target.value)}
                  placeholder="0"
                />
                <p className="text-sm text-muted-foreground">
                  キャッシュフロー予測で使用される最低月支出額
                </p>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>アプリケーション設定</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="notification_enabled">通知</Label>
                <select
                  id="notification_enabled"
                  className="w-full p-2 border border-input rounded-md"
                  value={settings.notification_enabled}
                  onChange={(e) => updateSetting('notification_enabled', e.target.value)}
                >
                  <option value="true">有効</option>
                  <option value="false">無効</option>
                </select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="theme">テーマ</Label>
                <select
                  id="theme"
                  className="w-full p-2 border border-input rounded-md"
                  value={settings.theme}
                  onChange={(e) => updateSetting('theme', e.target.value)}
                >
                  <option value="light">ライト</option>
                  <option value="dark">ダーク</option>
                  <option value="system">システム</option>
                </select>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Info className="h-5 w-5" />
                アプリケーション情報
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label>フロントエンドバージョン</Label>
                <Input
                  value={FRONTEND_VERSION}
                  readOnly
                  className="bg-muted"
                />
              </div>
              {versionInfo && (
                <div className="space-y-2">
                  <Label>バックエンドバージョン</Label>
                  <Input
                    value={versionInfo.version}
                    readOnly
                    className="bg-muted"
                  />
                </div>
              )}
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <User className="h-5 w-5" />
                アカウント情報
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {userInfo && (
                <>
                  <div className="flex items-center gap-4">
                    <Image
                      src={userInfo.picture || '/default-avatar.png'}
                      alt={userInfo.name || 'User'}
                      width={64}
                      height={64}
                      className="rounded-full"
                    />
                    <div>
                      <h3 className="font-semibold">{userInfo.name || 'Unknown'}</h3>
                      <p className="text-sm text-muted-foreground">{userInfo.email || 'No email'}</p>
                    </div>
                  </div>
                  <div className="space-y-2">
                    <Label>ユーザーID</Label>
                    <Input
                      value={userInfo.id || ''}
                      readOnly
                      className="bg-muted font-mono text-xs"
                    />
                  </div>
                  <div className="space-y-2">
                    <Label>Google ID</Label>
                    <Input
                      value={userInfo.google_id || ''}
                      readOnly
                      className="bg-muted font-mono text-xs"
                    />
                  </div>
                  <div className="space-y-2">
                    <Label>登録日時</Label>
                    <Input
                      value={userInfo.created_at ? new Date(userInfo.created_at).toLocaleString('ja-JP') : ''}
                      readOnly
                      className="bg-muted"
                    />
                  </div>
                </>
              )}
            </CardContent>
          </Card>
        </div>

        <div className="flex justify-end">
          <Button onClick={handleSave} disabled={isSaving}>
            <Save className="h-4 w-4 mr-2" />
            {isSaving ? '保存中...' : '設定を保存'}
          </Button>
        </div>
      </div>
    </MainLayout>
  );
}
