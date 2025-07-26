import {
  CreditCard,
  BankAccount,
  CardMonthlyTotal,
  IncomeSource,
  MonthlyIncomeRecord,
  RecurringPayment,
  CashflowProjection,
  DashboardSummary,
  AppSetting,
  UpdateSettingsRequest,
  VersionInfo,
  UserInfo,
} from '@/types/api';
import Cookies from 'js-cookie';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:4000';

class ApiClient {
  private baseURL: string;

  constructor(baseURL: string = API_BASE_URL) {
    this.baseURL = `${baseURL}/api/v1`;
  }

  private getAuthHeaders(): Record<string, string> {
    const token = Cookies.get('auth_token');
    return token ? { 'Authorization': `Bearer ${token}` } : {};
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...this.getAuthHeaders(),
      ...options.headers as Record<string, string>,
    };

    const config: RequestInit = {
      ...options,
      headers,
    };

    try {
      const response = await fetch(url, config);
      
      if (!response.ok) {
        if (response.status === 401) {
          // Unauthorized - redirect to login
          Cookies.remove('auth_token');
          window.location.href = '/login';
          return Promise.reject(new Error('Unauthorized'));
        }
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      // Handle No Content responses (204)
      if (response.status === 204) {
        return undefined as T;
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error('API request failed:', error);
      throw error;
    }
  }

  // Credit Cards API
  async getCreditCards(): Promise<CreditCard[]> {
    return this.request<CreditCard[]>('/credit-cards');
  }

  async getCreditCard(id: string): Promise<CreditCard> {
    return this.request<CreditCard>(`/credit-cards/${id}`);
  }

  async createCreditCard(creditCard: Omit<CreditCard, 'id' | 'created_at' | 'updated_at' | 'user_id'>): Promise<CreditCard> {
    return this.request<CreditCard>('/credit-cards', {
      method: 'POST',
      body: JSON.stringify(creditCard),
    });
  }

  async updateCreditCard(id: string, creditCard: Omit<CreditCard, 'id' | 'created_at' | 'updated_at' | 'user_id'>): Promise<CreditCard> {
    return this.request<CreditCard>(`/credit-cards/${id}`, {
      method: 'PUT',
      body: JSON.stringify(creditCard),
    });
  }

  async deleteCreditCard(id: string): Promise<void> {
    await this.request<void>(`/credit-cards/${id}`, {
      method: 'DELETE',
    });
  }

  // Bank Accounts API
  async getBankAccounts(): Promise<BankAccount[]> {
    return this.request<BankAccount[]>('/bank-accounts');
  }

  async getBankAccount(id: string): Promise<BankAccount> {
    return this.request<BankAccount>(`/bank-accounts/${id}`);
  }

  async createBankAccount(account: Omit<BankAccount, 'id' | 'created_at' | 'updated_at' | 'user_id'>): Promise<BankAccount> {
    return this.request<BankAccount>('/bank-accounts', {
      method: 'POST',
      body: JSON.stringify(account),
    });
  }

  async updateBankAccount(id: string, account: Omit<BankAccount, 'id' | 'created_at' | 'updated_at' | 'user_id'>): Promise<BankAccount> {
    return this.request<BankAccount>(`/bank-accounts/${id}`, {
      method: 'PUT',
      body: JSON.stringify(account),
    });
  }

  async deleteBankAccount(id: string): Promise<void> {
    await this.request<void>(`/bank-accounts/${id}`, {
      method: 'DELETE',
    });
  }

  // Card Monthly Totals API
  async getCardMonthlyTotals(creditCardId: string): Promise<CardMonthlyTotal[]> {
    return this.request<CardMonthlyTotal[]>(`/card-monthly-totals?credit_card_id=${creditCardId}`);
  }

  async getCardMonthlyTotal(id: string): Promise<CardMonthlyTotal> {
    return this.request<CardMonthlyTotal>(`/card-monthly-totals/${id}`);
  }

  async createCardMonthlyTotal(total: Omit<CardMonthlyTotal, 'id' | 'created_at' | 'updated_at'>): Promise<CardMonthlyTotal> {
    return this.request<CardMonthlyTotal>('/card-monthly-totals', {
      method: 'POST',
      body: JSON.stringify(total),
    });
  }

  async updateCardMonthlyTotal(id: string, total: Omit<CardMonthlyTotal, 'id' | 'created_at' | 'updated_at'>): Promise<CardMonthlyTotal> {
    return this.request<CardMonthlyTotal>(`/card-monthly-totals/${id}`, {
      method: 'PUT',
      body: JSON.stringify(total),
    });
  }

  async deleteCardMonthlyTotal(id: string): Promise<void> {
    await this.request<void>(`/card-monthly-totals/${id}`, {
      method: 'DELETE',
    });
  }

  // Income Sources API
  async getIncomeSources(): Promise<IncomeSource[]> {
    return this.request<IncomeSource[]>('/income-sources');
  }

  async getIncomeSource(id: string): Promise<IncomeSource> {
    return this.request<IncomeSource>(`/income-sources/${id}`);
  }

  async createIncomeSource(source: Omit<IncomeSource, 'id' | 'created_at' | 'updated_at' | 'user_id'>): Promise<IncomeSource> {
    return this.request<IncomeSource>('/income-sources', {
      method: 'POST',
      body: JSON.stringify(source),
    });
  }

  async updateIncomeSource(id: string, source: Omit<IncomeSource, 'id' | 'created_at' | 'updated_at' | 'user_id'>): Promise<IncomeSource> {
    return this.request<IncomeSource>(`/income-sources/${id}`, {
      method: 'PUT',
      body: JSON.stringify(source),
    });
  }

  async deleteIncomeSource(id: string): Promise<void> {
    await this.request<void>(`/income-sources/${id}`, {
      method: 'DELETE',
    });
  }

  // Monthly Income Records API
  async getMonthlyIncomeRecords(incomeSourceId: string): Promise<MonthlyIncomeRecord[]> {
    return this.request<MonthlyIncomeRecord[]>(`/monthly-income-records?income_source_id=${incomeSourceId}`);
  }

  async getMonthlyIncomeRecord(id: string): Promise<MonthlyIncomeRecord> {
    return this.request<MonthlyIncomeRecord>(`/monthly-income-records/${id}`);
  }

  async createMonthlyIncomeRecord(record: Omit<MonthlyIncomeRecord, 'id' | 'created_at' | 'updated_at'>): Promise<MonthlyIncomeRecord> {
    return this.request<MonthlyIncomeRecord>('/monthly-income-records', {
      method: 'POST',
      body: JSON.stringify(record),
    });
  }

  async updateMonthlyIncomeRecord(id: string, record: Omit<MonthlyIncomeRecord, 'id' | 'created_at' | 'updated_at'>): Promise<MonthlyIncomeRecord> {
    return this.request<MonthlyIncomeRecord>(`/monthly-income-records/${id}`, {
      method: 'PUT',
      body: JSON.stringify(record),
    });
  }

  async deleteMonthlyIncomeRecord(id: string): Promise<void> {
    await this.request<void>(`/monthly-income-records/${id}`, {
      method: 'DELETE',
    });
  }

  // Recurring Payments API
  async getRecurringPayments(): Promise<RecurringPayment[]> {
    return this.request<RecurringPayment[]>('/recurring-payments');
  }

  async getRecurringPayment(id: string): Promise<RecurringPayment> {
    return this.request<RecurringPayment>(`/recurring-payments/${id}`);
  }

  async createRecurringPayment(payment: Omit<RecurringPayment, 'id' | 'created_at' | 'updated_at' | 'user_id'>): Promise<RecurringPayment> {
    return this.request<RecurringPayment>('/recurring-payments', {
      method: 'POST',
      body: JSON.stringify(payment),
    });
  }

  async updateRecurringPayment(id: string, payment: Omit<RecurringPayment, 'id' | 'created_at' | 'updated_at' | 'user_id'>): Promise<RecurringPayment> {
    return this.request<RecurringPayment>(`/recurring-payments/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payment),
    });
  }

  async deleteRecurringPayment(id: string): Promise<void> {
    await this.request<void>(`/recurring-payments/${id}`, {
      method: 'DELETE',
    });
  }

  // Cashflow Projection API
  async getCashflowProjection(months: number = 36, onlyChanges: boolean = true): Promise<CashflowProjection[]> {
    return this.request<CashflowProjection[]>(`/cashflow-projection?months=${months}&onlyChanges=${onlyChanges}`);
  }

  // Dashboard API
  async getDashboardSummary(): Promise<DashboardSummary> {
    return this.request<DashboardSummary>('/dashboard/summary');
  }

  // Settings API
  async getSettings(): Promise<AppSetting[]> {
    return this.request<AppSetting[]>('/settings');
  }

  async updateSettings(settings: UpdateSettingsRequest): Promise<Record<string, string>> {
    return this.request<Record<string, string>>('/settings', {
      method: 'PUT',
      body: JSON.stringify(settings),
    });
  }

  // Version API
  async getVersion(): Promise<VersionInfo> {
    return this.request<VersionInfo>('/version');
  }

  // User API
  async getCurrentUser(): Promise<UserInfo> {
    return this.request<UserInfo>('/auth/me');
  }
}

export const apiClient = new ApiClient();
export default ApiClient;
