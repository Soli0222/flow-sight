// API response types based on backend Swagger documentation

export interface CreditCard {
  id: string;
  user_id: string;
  name: string;
  bank_account: string;
  closing_day?: number; // Closing day of the month
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
  credit_card_id: string;
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
  payment_day?: number; // For monthly_fixed income (1-31)
  scheduled_date?: string; // For one_time income (ISO date string)
  scheduled_year_month?: string; // For one-time income (backward compatibility)
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
  total_payments?: number; // For loans, undefined means infinite payments
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

export interface DashboardSummary {
  total_balance: number;
  monthly_income: number;
  monthly_expense: number;
  total_assets: number;
  recent_activities: CashflowProjection[];
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

export interface VersionInfo {
  version: string;
}

export interface UserInfo {
  id: string;
  email: string;
  name: string;
  picture: string;
  google_id: string;
  created_at: string;
  updated_at: string;
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
