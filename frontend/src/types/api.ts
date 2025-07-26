// API response types based on backend Swagger documentation

export interface Asset {
  id: string;
  user_id: string;
  name: string;
  asset_type: 'card' | 'loan';
  bank_account: string;
  closing_day?: number; // For credit cards
  payment_day: number;
  created_at: string;
  updated_at: string;
}

export interface BankAccount {
  id: string;
  user_id: string;
  name: string;
  balance: number; // Amount in cents
  created_at: string;
  updated_at: string;
}

export interface CardMonthlyTotal {
  id: string;
  asset_id: string;
  year_month: string; // Format: "2024-01"
  total_amount: number; // Amount in cents
  is_confirmed: boolean;
  created_at: string;
  updated_at: string;
}

export interface IncomeSource {
  id: string;
  user_id: string;
  name: string;
  income_type: 'monthly_fixed' | 'one_time';
  base_amount: number; // Amount in cents
  bank_account: string;
  is_active: boolean;
  scheduled_year_month?: string; // For one-time income
  created_at: string;
  updated_at: string;
}

export interface MonthlyIncomeRecord {
  id: string;
  income_source_id: string;
  year_month: string; // Format: "2024-01"
  actual_amount: number; // Amount in cents
  is_confirmed: boolean;
  note?: string;
  created_at: string;
  updated_at: string;
}

export interface RecurringPayment {
  id: string;
  user_id: string;
  name: string;
  amount: number; // Amount in cents
  payment_day: number;
  bank_account: string;
  start_year_month: string; // Format: "2024-01"
  total_payments?: number; // For loans
  remaining_payments?: number;
  is_active: boolean;
  note?: string;
  created_at: string;
  updated_at: string;
}

export interface CashflowProjectionDetail {
  type: 'income' | 'recurring_payment' | 'card_payment';
  description: string;
  amount: number;
}

export interface CashflowProjection {
  date: string;
  income: number;
  expense: number;
  balance: number;
  details: CashflowProjectionDetail[];
}

export interface AppSetting {
  id: string;
  user_id: string;
  key: string;
  value: string;
  created_at: string;
  updated_at: string;
}

export interface UpdateSettingsRequest {
  settings: Record<string, string>;
}

// API Error Response
export interface ApiError {
  message: string;
  code?: string;
}

// Common API response wrapper
export interface ApiResponse<T> {
  data: T;
  error?: ApiError;
}
