import { CreditCard, BankAccount, DashboardSummary, CashflowProjection } from '@/types/api'

export const mockCreditCard: CreditCard = {
  id: '1',
  name: 'テストカード',
  bank_account: 'bank1',
  closing_day: 25,
  payment_day: 10,
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
}

export const mockBankAccount: BankAccount = {
  id: 'bank1',
  name: 'テスト銀行',
  balance: 100000,
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
}

export const mockCashflowProjection: CashflowProjection = {
  date: '2024-01-15',
  income: 50000,
  expense: 30000,
  balance: 120000,
  details: [
    {
      type: 'income',
      description: 'テスト収入',
      amount: 50000,
    },
  ],
}

export const mockDashboardSummary: DashboardSummary = {
  total_balance: 100000,
  monthly_income: 250000,
  monthly_expense: 150000,
  total_assets: 500000,
  recent_activities: [mockCashflowProjection],
}

export const mockApiResponses = {
  creditCards: [mockCreditCard],
  bankAccounts: [mockBankAccount],
  dashboardSummary: mockDashboardSummary,
}
