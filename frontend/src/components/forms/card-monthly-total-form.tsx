'use client';

import React, { useState } from 'react';
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
import { CreditCard, CardMonthlyTotal } from '@/types/api';
import { toast } from 'sonner';

const formSchema = z.object({
  credit_card_id: z.string().min(1, 'カードを選択してください'),
  year_month: z.string().min(1, '対象年月は必須です').regex(/^\d{4}-\d{2}$/, '年月は YYYY-MM 形式で入力してください'),
  total_amount: z.number().min(0, '金額は0以上である必要があります'),
  is_confirmed: z.boolean(),
});

type FormData = z.infer<typeof formSchema>;

interface CardMonthlyTotalFormProps {
  total?: CardMonthlyTotal | null;
  creditCards: CreditCard[];
  onSuccess: () => void;
  onCancel: () => void;
}

export function CardMonthlyTotalForm({ 
  total, 
  creditCards, 
  onSuccess, 
  onCancel 
}: CardMonthlyTotalFormProps) {
  const apiClient = useApi();
  const [isLoading, setIsLoading] = useState(false);

  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      credit_card_id: total?.credit_card_id || '',
      year_month: total?.year_month || '',
      total_amount: total ? total.total_amount / 100 : 0,
      is_confirmed: total?.is_confirmed || false,
    },
  });

  const onSubmit = async (data: FormData) => {
    try {
      setIsLoading(true);
      
      const totalData = {
        credit_card_id: data.credit_card_id,
        year_month: data.year_month,
        total_amount: Math.round(data.total_amount * 100), // Convert to cents
        is_confirmed: data.is_confirmed,
      };

      if (total) {
        await apiClient.updateCardMonthlyTotal(total.id, totalData);
        toast.success('月次利用額を更新しました');
      } else {
        await apiClient.createCardMonthlyTotal(totalData);
        toast.success('月次利用額を登録しました');
      }

      onSuccess();
    } catch (error) {
      toast.error(total ? '月次利用額の更新に失敗しました' : '月次利用額の登録に失敗しました');
      console.error('Failed to save card monthly total:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Dialog open={true} onOpenChange={(open) => !open && onCancel()}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>
            {total ? '月次利用額を編集' : '月次利用額を追加'}
          </DialogTitle>
        </DialogHeader>
        
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="credit_card_id"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>クレジットカード</FormLabel>
                  <Select 
                    onValueChange={field.onChange} 
                    defaultValue={field.value}
                    disabled={isLoading}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="クレジットカードを選択" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {creditCards.map((creditCard) => (
                        <SelectItem key={creditCard.id} value={creditCard.id}>
                          {creditCard.name}
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
              name="year_month"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>対象年月</FormLabel>
                  <FormControl>
                    <Input 
                      {...field} 
                      type="month"
                      disabled={isLoading}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="total_amount"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>利用額</FormLabel>
                  <FormControl>
                    <Input 
                      {...field} 
                      type="number"
                      min="0"
                      step="1"
                      placeholder="0"
                      onChange={(e) => field.onChange(Number(e.target.value))}
                      disabled={isLoading}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="is_confirmed"
              render={({ field }) => (
                <FormItem className="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm">
                  <div className="space-y-0.5">
                    <FormLabel className="text-base">確定済み</FormLabel>
                    <div className="text-sm text-muted-foreground">
                      この月次利用額を確定としてマークします
                    </div>
                  </div>
                  <FormControl>
                    <input
                      type="checkbox"
                      checked={field.value}
                      onChange={field.onChange}
                      disabled={isLoading}
                      className="h-4 w-4"
                    />
                  </FormControl>
                </FormItem>
              )}
            />

            <div className="flex gap-2 pt-4">
              <Button 
                type="button" 
                variant="outline" 
                className="flex-1"
                onClick={onCancel}
                disabled={isLoading}
              >
                キャンセル
              </Button>
              <Button 
                type="submit" 
                className="flex-1"
                disabled={isLoading}
              >
                {isLoading ? '保存中...' : total ? '更新' : '作成'}
              </Button>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
