'use client';

import React, { useState, useEffect } from 'react';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import * as z from 'zod';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { useApi } from '@/components/providers/api-provider';
import { Asset, BankAccount } from '@/types/api';
import { toast } from 'sonner';

const formSchema = z.object({
  name: z.string().min(1, '資産名は必須です'),
  asset_type: z.enum(['card', 'loan'], { message: '資産タイプを選択してください' }),
  bank_account: z.string().min(1, '銀行口座を選択してください'),
  closing_day: z.string().optional(),
  payment_day: z.string().min(1, '支払日は必須です'),
});

type FormData = z.infer<typeof formSchema>;

interface AssetFormProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  asset?: Asset | null;
  bankAccounts: BankAccount[];
  onSuccess: () => void;
}

export function AssetForm({
  open,
  onOpenChange,
  asset,
  bankAccounts,
  onSuccess,
}: AssetFormProps) {
  const apiClient = useApi();
  const [isSubmitting, setIsSubmitting] = useState(false);

  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: '',
      asset_type: 'card',
      bank_account: '',
      closing_day: '',
      payment_day: '',
    },
  });

  const assetType = form.watch('asset_type');

  useEffect(() => {
    if (asset) {
      form.reset({
        name: asset.name,
        asset_type: asset.asset_type,
        bank_account: asset.bank_account,
        closing_day: asset.closing_day?.toString() || '',
        payment_day: asset.payment_day.toString(),
      });
    } else {
      form.reset({
        name: '',
        asset_type: 'card',
        bank_account: '',
        closing_day: '',
        payment_day: '',
      });
    }
  }, [asset, form]);

  const onSubmit = async (data: FormData) => {
    try {
      setIsSubmitting(true);
      
      const assetData = {
        name: data.name,
        asset_type: data.asset_type,
        bank_account: data.bank_account,
        closing_day: data.closing_day ? parseInt(data.closing_day) : undefined,
        payment_day: parseInt(data.payment_day),
      };

      if (asset) {
        await apiClient.updateAsset(asset.id, assetData);
        toast.success('資産を更新しました');
      } else {
        await apiClient.createAsset(assetData);
        toast.success('資産を作成しました');
      }

      onSuccess();
    } catch (error) {
      toast.error(asset ? '資産の更新に失敗しました' : '資産の作成に失敗しました');
      console.error('Failed to save asset:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen && !isSubmitting) {
      form.reset();
    }
    onOpenChange(newOpen);
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>
            {asset ? '資産を編集' : '資産を追加'}
          </DialogTitle>
        </DialogHeader>
        
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>資産名</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="例: 楽天カード"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="asset_type"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>資産タイプ</FormLabel>
                  <Select onValueChange={field.onChange} defaultValue={field.value}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="資産タイプを選択" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value="card">クレジットカード</SelectItem>
                      <SelectItem value="loan">ローン</SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="bank_account"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>銀行口座</FormLabel>
                  <Select onValueChange={field.onChange} defaultValue={field.value}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="銀行口座を選択" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {bankAccounts.map((account) => (
                        <SelectItem key={account.id} value={account.id}>
                          {account.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            {assetType === 'card' && (
              <FormField
                control={form.control}
                name="closing_day"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>締め日</FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        min="1"
                        max="31"
                        placeholder="例: 15"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            )}

            <FormField
              control={form.control}
              name="payment_day"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>支払日</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      min="1"
                      max="31"
                      placeholder="例: 27"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="flex justify-end gap-2 pt-4">
              <Button
                type="button"
                variant="outline"
                onClick={() => handleOpenChange(false)}
                disabled={isSubmitting}
              >
                キャンセル
              </Button>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? '保存中...' : asset ? '更新' : '作成'}
              </Button>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
