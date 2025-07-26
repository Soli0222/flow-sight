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
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { useApi } from '@/components/providers/api-provider';
import { BankAccount } from '@/types/api';
import { toast } from 'sonner';

const formSchema = z.object({
  name: z.string().min(1, '口座名は必須です'),
  balance: z.string().min(1, '残高は必須です'),
});

type FormData = z.infer<typeof formSchema>;

interface BankAccountFormProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  account?: BankAccount | null;
  onSuccess: () => void;
}

export function BankAccountForm({
  open,
  onOpenChange,
  account,
  onSuccess,
}: BankAccountFormProps) {
  const apiClient = useApi();
  const [isSubmitting, setIsSubmitting] = useState(false);

  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: '',
      balance: '0',
    },
  });

  useEffect(() => {
    if (account) {
      form.reset({
        name: account.name,
        balance: (account.balance / 100).toString(),
      });
    } else {
      form.reset({
        name: '',
        balance: '0',
      });
    }
  }, [account, form]);

  const onSubmit = async (data: FormData) => {
    try {
      setIsSubmitting(true);
      
      const balanceInCents = Math.round(parseFloat(data.balance) * 100);
      
      const accountData = {
        name: data.name,
        balance: balanceInCents,
      };

      if (account) {
        await apiClient.updateBankAccount(account.id, accountData);
        toast.success('銀行口座を更新しました');
      } else {
        await apiClient.createBankAccount(accountData);
        toast.success('銀行口座を作成しました');
      }

      onSuccess();
    } catch (error) {
      toast.error(account ? '銀行口座の更新に失敗しました' : '銀行口座の作成に失敗しました');
      console.error('Failed to save bank account:', error);
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
            {account ? '銀行口座を編集' : '銀行口座を追加'}
          </DialogTitle>
        </DialogHeader>
        
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>口座名</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="例: 三菱UFJ銀行 普通預金"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="balance"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>残高 (円)</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      step="0.01"
                      placeholder="0"
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
                {isSubmitting ? '保存中...' : account ? '更新' : '作成'}
              </Button>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
