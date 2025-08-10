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
} from '@/types/api';

// 同じドメインを使用するため、環境変数が設定されていない場合は空文字（相対パス）を使用
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || '';

class ApiClient {
  private baseURL: string;

  constructor(baseURL: string = API_BASE_URL) {
    this.baseURL = `${baseURL}/api/v1`;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options.headers as Record<string, string>),
    };

    const config: RequestInit = {
      ...options,
      headers,
    };

    try {
      const response = await fetch(url, config);
      
      if (!response.ok) {
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

  async createCreditCard(creditCard: Omit<CreditCard, 'id' | 'created_at' | 'updated_at'>): Promise<CreditCard> {
    return this.request<CreditCard>('/credit-cards', {
      method: 'POST',
      body: JSON.stringify(creditCard),
    });
  }

  async updateCreditCard(id: string, creditCard: Omit<CreditCard, 'id' | 'created_at' | 'updated_at'>): Promise<CreditCard> {
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

  async createBankAccount(account: Omit<BankAccount, 'id' | 'created_at' | 'updated_at'>): Promise<BankAccount> {
    return this.request<BankAccount>('/bank-accounts', {
      method: 'POST',
      body: JSON.stringify(account),
    });
  }

  async updateBankAccount(id: string, account: Omit<BankAccount, 'id' | 'created_at' | 'updated_at'>): Promise<BankAccount> {
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

  async createIncomeSource(source: Omit<IncomeSource, 'id' | 'created_at' | 'updated_at'>): Promise<IncomeSource> {
    return this.request<IncomeSource>('/income-sources', {
      method: 'POST',
      body: JSON.stringify(source),
    });
  }

  async updateIncomeSource(id: string, source: Omit<IncomeSource, 'id' | 'created_at' | 'updated_at'>): Promise<IncomeSource> {
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

  async createRecurringPayment(payment: Omit<RecurringPayment, 'id' | 'created_at' | 'updated_at'>): Promise<RecurringPayment> {
    return this.request<RecurringPayment>('/recurring-payments', {
      method: 'POST',
      body: JSON.stringify(payment),
    });
  }

  async updateRecurringPayment(id: string, payment: Omit<RecurringPayment, 'id' | 'created_at' | 'updated_at'>): Promise<RecurringPayment> {
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

  // Cashflow Projection
  async getCashflowProjection(months: number, onlyChanges: boolean): Promise<CashflowProjection[]> {
    const params = new URLSearchParams({ months: months.toString(), onlyChanges: String(onlyChanges) });
    return this.request<CashflowProjection[]>(`/cashflow-projection?${params.toString()}`);
  }

  // Dashboard Summary
  async getDashboardSummary(): Promise<DashboardSummary> {
    return this.request<DashboardSummary>('/dashboard/summary');
  }

  // App Settings
  async getSettings(): Promise<AppSetting[]> {
    return this.request<AppSetting[]>('/settings');
  }

  async updateSettings(settings: UpdateSettingsRequest): Promise<{ success: boolean }> {
    return this.request<{ success: boolean }>(`/settings`, {
      method: 'PUT',
      body: JSON.stringify(settings),
    });
  }

  // Version
  async getVersion(): Promise<VersionInfo> {
    return this.request<VersionInfo>('/version');
  }
}

export default ApiClient;
