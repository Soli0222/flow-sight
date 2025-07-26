export const formatCurrency = (amountInCents: number): string => {
  const amount = amountInCents / 100;
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(amount);
};

export const parseCurrency = (formattedAmount: string): number => {
  // Remove currency symbols and convert to number
  const cleanAmount = formattedAmount.replace(/[^0-9.-]/g, '');
  const amount = parseFloat(cleanAmount);
  return Math.round(amount * 100); // Convert to cents
};

export const formatDate = (dateString: string): string => {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  }).format(date);
};

export const formatYearMonth = (yearMonth: string): string => {
  const [year, month] = yearMonth.split('-');
  const date = new Date(parseInt(year), parseInt(month) - 1);
  return new Intl.DateTimeFormat('ja-JP', {
    year: 'numeric',
    month: 'long',
  }).format(date);
};

export const getCurrentYearMonth = (): string => {
  const now = new Date();
  const year = now.getFullYear();
  const month = (now.getMonth() + 1).toString().padStart(2, '0');
  return `${year}-${month}`;
};

export const getNextYearMonth = (yearMonth: string): string => {
  const [year, month] = yearMonth.split('-').map(Number);
  const date = new Date(year, month - 1); // month is 0-indexed
  date.setMonth(date.getMonth() + 1);
  const newYear = date.getFullYear();
  const newMonth = (date.getMonth() + 1).toString().padStart(2, '0');
  return `${newYear}-${newMonth}`;
};

export const getPreviousYearMonth = (yearMonth: string): string => {
  const [year, month] = yearMonth.split('-').map(Number);
  const date = new Date(year, month - 1); // month is 0-indexed
  date.setMonth(date.getMonth() - 1);
  const newYear = date.getFullYear();
  const newMonth = (date.getMonth() + 1).toString().padStart(2, '0');
  return `${newYear}-${newMonth}`;
};

export const validateRequired = (value: string | number): boolean => {
  if (typeof value === 'string') {
    return value.trim().length > 0;
  }
  return value !== null && value !== undefined;
};

export const validateAmount = (amount: number): boolean => {
  return amount >= 0;
};

export const validateYearMonth = (yearMonth: string): boolean => {
  const regex = /^\d{4}-\d{2}$/;
  if (!regex.test(yearMonth)) return false;
  
  const [year, month] = yearMonth.split('-').map(Number);
  return year >= 1900 && year <= 2100 && month >= 1 && month <= 12;
};

export const validateDay = (day: number): boolean => {
  return day >= 1 && day <= 31;
};
