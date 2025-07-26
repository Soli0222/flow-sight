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
import { IncomeSource, BankAccount } from '@/types/api';
import { toast } from 'sonner';

const formSchema = z.object({
  name: z.string().min(1, '収入源名は必須です'),
  income_type: z.enum(['monthly_fixed', 'one_time'], { message: '収入タイプを選択してください' }),
  base_amount: z.string().min(1, '金額は必須です'),
  bank_account: z.string().min(1, '銀行口座を選択してください'),
  is_active: z.string(),
  payment_day: z.string().optional(),
  scheduled_date: z.string().optional(),
}).superRefine((data, ctx) => {
  if (data.income_type === 'monthly_fixed') {
    if (!data.payment_day || data.payment_day === '') {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: '月額固定収入の場合、支払日は必須です',
        path: ['payment_day'],
      });
    } else {
      const day = parseInt(data.payment_day);
      if (isNaN(day) || day < 1 || day > 31) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: '支払日は1〜31の数値で入力してください',
          path: ['payment_day'],
        });
      }
    }
  } else if (data.income_type === 'one_time') {
    if (!data.scheduled_date || data.scheduled_date === '') {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: '一時的収入の場合、予定日は必須です',
        path: ['scheduled_date'],
      });
    }
  }
});

type FormData = z.infer<typeof formSchema>;

interface IncomeFormProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  income?: IncomeSource | null;
  bankAccounts: BankAccount[];
  onSuccess: () => void;
}

export function IncomeForm({
  open,
  onOpenChange,
  income,
  bankAccounts,
  onSuccess,
}: IncomeFormProps) {
  const apiClient = useApi();
  const [isSubmitting, setIsSubmitting] = useState(false);

  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: '',
      income_type: 'monthly_fixed',
      base_amount: '',
      bank_account: '',
      is_active: 'true',
      payment_day: '25',
      scheduled_date: '',
    },
  });

  const incomeType = form.watch('income_type');

  useEffect(() => {
    if (income) {
      form.reset({
        name: income.name,
        income_type: income.income_type,
        base_amount: (income.base_amount / 100).toString(),
        bank_account: income.bank_account,
        is_active: income.is_active ? 'true' : 'false',
        payment_day: income.payment_day?.toString() || '25',
        scheduled_date: income.scheduled_date ? income.scheduled_date.split('T')[0] : '',
      });
    } else {
      form.reset({
        name: '',
        income_type: 'monthly_fixed',
        base_amount: '',
        bank_account: '',
        is_active: 'true',
        payment_day: '25',
        scheduled_date: '',
      });
    }
  }, [income, form]);

  const onSubmit = async (data: FormData) => {
    try {
      setIsSubmitting(true);
      
      const baseAmountInCents = Math.round(parseFloat(data.base_amount) * 100);
      
      const incomeData = {
        name: data.name,
        income_type: data.income_type,
        base_amount: baseAmountInCents,
        bank_account: data.bank_account,
        is_active: data.is_active === 'true',
        payment_day: data.income_type === 'monthly_fixed' && data.payment_day ? parseInt(data.payment_day) : undefined,
        scheduled_date: data.income_type === 'one_time' && data.scheduled_date ? data.scheduled_date : undefined,
      };

      if (income) {
        await apiClient.updateIncomeSource(income.id, incomeData);
        toast.success('収入源を更新しました');
      } else {
        await apiClient.createIncomeSource(incomeData);
        toast.success('収入源を作成しました');
      }

      onSuccess();
    } catch (error) {
      toast.error(income ? '収入源の更新に失敗しました' : '収入源の作成に失敗しました');
      console.error('Failed to save income source:', error);
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
            {income ? '収入源を編集' : '収入源を追加'}
          </DialogTitle>
        </DialogHeader>
        
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>収入源名</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="例: 給与"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="income_type"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>収入タイプ</FormLabel>
                  <Select onValueChange={field.onChange} defaultValue={field.value}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="収入タイプを選択" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value="monthly_fixed">月額固定</SelectItem>
                      <SelectItem value="one_time">一時的</SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="base_amount"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>金額 (円)</FormLabel>
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

            <FormField
              control={form.control}
              name="bank_account"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>振込先銀行口座</FormLabel>
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

            <FormField
              control={form.control}
              name="is_active"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>状態</FormLabel>
                  <Select onValueChange={field.onChange} defaultValue={field.value}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="状態を選択" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value="true">アクティブ</SelectItem>
                      <SelectItem value="false">非アクティブ</SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            {incomeType === 'one_time' && (
              <FormField
                control={form.control}
                name="scheduled_date"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>予定日</FormLabel>
                    <FormControl>
                      <Input
                        type="date"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            )}

            {incomeType === 'monthly_fixed' && (
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
                        placeholder="25"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            )}

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
                {isSubmitting ? '保存中...' : income ? '更新' : '作成'}
              </Button>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
