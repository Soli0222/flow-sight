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
import { CreditCard, BankAccount } from '@/types/api';
import { toast } from 'sonner';

const formSchema = z.object({
  name: z.string().min(1, 'クレジットカード名は必須です'),
  bank_account: z.string().min(1, '銀行口座を選択してください'),
  closing_day: z.string().optional(),
  payment_day: z.string().min(1, '支払日は必須です'),
});

type FormData = z.infer<typeof formSchema>;

interface CreditCardFormProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  creditCard?: CreditCard | null;
  bankAccounts: BankAccount[];
  onSuccess: () => void;
}

export function CreditCardForm({
  open,
  onOpenChange,
  creditCard,
  bankAccounts,
  onSuccess,
}: CreditCardFormProps) {
  const apiClient = useApi();
  const [isSubmitting, setIsSubmitting] = useState(false);

  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: '',
      bank_account: '',
      closing_day: '',
      payment_day: '',
    },
  });

  useEffect(() => {
    if (creditCard) {
      form.reset({
        name: creditCard.name,
        bank_account: creditCard.bank_account,
        closing_day: creditCard.closing_day?.toString() || '',
        payment_day: creditCard.payment_day.toString(),
      });
    } else {
      form.reset({
        name: '',
        bank_account: '',
        closing_day: '',
        payment_day: '',
      });
    }
  }, [creditCard, form]);

  const onSubmit = async (data: FormData) => {
    try {
      setIsSubmitting(true);
      
      const creditCardData = {
        name: data.name,
        bank_account: data.bank_account,
        closing_day: data.closing_day ? parseInt(data.closing_day) : undefined,
        payment_day: parseInt(data.payment_day),
      };

      if (creditCard) {
        await apiClient.updateCreditCard(creditCard.id, creditCardData);
        toast.success('クレジットカードを更新しました');
      } else {
        await apiClient.createCreditCard(creditCardData);
        toast.success('クレジットカードを作成しました');
      }

      onSuccess();
    } catch (error) {
      toast.error(creditCard ? 'クレジットカードの更新に失敗しました' : 'クレジットカードの作成に失敗しました');
      console.error('Failed to save credit card:', error);
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
            {creditCard ? 'クレジットカードを編集' : 'クレジットカードを追加'}
          </DialogTitle>
        </DialogHeader>
        
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>クレジットカード名</FormLabel>
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
                {isSubmitting ? '保存中...' : creditCard ? '更新' : '作成'}
              </Button>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
