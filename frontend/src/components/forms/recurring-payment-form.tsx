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
import { Textarea } from '@/components/ui/textarea';
import { Button } from '@/components/ui/button';
import { useApi } from '@/components/providers/api-provider';
import { RecurringPayment, BankAccount } from '@/types/api';
import { toast } from 'sonner';

const formSchema = z.object({
  name: z.string().min(1, '支払い名は必須です'),
  amount: z.string().min(1, '金額は必須です'),
  payment_day: z.string().min(1, '支払日は必須です'),
  bank_account: z.string().min(1, '銀行口座を選択してください'),
  start_year_month: z.string().min(1, '開始年月は必須です'),
  is_active: z.string(),
  note: z.string().optional(),
});

type FormData = z.infer<typeof formSchema>;

interface RecurringPaymentFormProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  payment?: RecurringPayment | null;
  bankAccounts: BankAccount[];
  onSuccess: () => void;
}

export function RecurringPaymentForm({
  open,
  onOpenChange,
  payment,
  bankAccounts,
  onSuccess,
}: RecurringPaymentFormProps) {
  const apiClient = useApi();
  const [isSubmitting, setIsSubmitting] = useState(false);

  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: '',
      amount: '',
      payment_day: '',
      bank_account: '',
      start_year_month: '',
      is_active: 'true',
      note: '',
    },
  });

  useEffect(() => {
    if (payment) {
      form.reset({
        name: payment.name,
        amount: (payment.amount / 100).toString(),
        payment_day: payment.payment_day.toString(),
        bank_account: payment.bank_account,
        start_year_month: payment.start_year_month,
        is_active: payment.is_active ? 'true' : 'false',
        note: payment.note || '',
      });
    } else {
      form.reset({
        name: '',
        amount: '',
        payment_day: '',
        bank_account: '',
        start_year_month: '',
        is_active: 'true',
        note: '',
      });
    }
  }, [payment, form]);

  const onSubmit = async (data: FormData) => {
    try {
      setIsSubmitting(true);
      
      const amountInCents = Math.round(parseFloat(data.amount) * 100);
      
      const paymentData = {
        name: data.name,
        amount: amountInCents,
        payment_day: parseInt(data.payment_day),
        bank_account: data.bank_account,
        start_year_month: data.start_year_month,
        is_active: data.is_active === 'true',
        note: data.note || undefined,
      };

      if (payment) {
        await apiClient.updateRecurringPayment(payment.id, paymentData);
        toast.success('定期支払いを更新しました');
      } else {
        await apiClient.createRecurringPayment(paymentData);
        toast.success('定期支払いを作成しました');
      }

      onSuccess();
    } catch (error) {
      toast.error(payment ? '定期支払いの更新に失敗しました' : '定期支払いの作成に失敗しました');
      console.error('Failed to save recurring payment:', error);
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
      <DialogContent className="sm:max-w-[500px] max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>
            {payment ? '定期支払いを編集' : '定期支払いを追加'}
          </DialogTitle>
        </DialogHeader>
        
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>支払い名</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="例: 電気代"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="amount"
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

            <FormField
              control={form.control}
              name="bank_account"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>支払口座</FormLabel>
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
              name="start_year_month"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>開始年月</FormLabel>
                  <FormControl>
                    <Input
                      type="month"
                      {...field}
                    />
                  </FormControl>
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

            <FormField
              control={form.control}
              name="note"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>備考</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="メモや備考を入力..."
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
                {isSubmitting ? '保存中...' : payment ? '更新' : '作成'}
              </Button>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
