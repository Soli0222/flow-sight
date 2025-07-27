import {
  formatCurrency,
  parseCurrency,
  formatDate,
  formatYearMonth,
  getCurrentYearMonth,
  getNextYearMonth,
  getPreviousYearMonth,
  validateRequired,
  validateAmount,
  validateYearMonth,
  validateDay,
} from '@/lib/utils-extended'

describe('utils-extended', () => {
  describe('formatCurrency', () => {
    it('should format cents to Japanese currency', () => {
      expect(formatCurrency(100000)).toBe('￥1,000')
      expect(formatCurrency(150)).toBe('￥2')
      expect(formatCurrency(0)).toBe('￥0')
      expect(formatCurrency(-50000)).toBe('-￥500')
    })

    it('should handle decimal places correctly', () => {
      expect(formatCurrency(99)).toBe('￥1')
      expect(formatCurrency(101)).toBe('￥1')
    })
  })

  describe('parseCurrency', () => {
    it('should parse formatted currency back to cents', () => {
      expect(parseCurrency('￥1,000')).toBe(100000)
      expect(parseCurrency('￥500')).toBe(50000)
      expect(parseCurrency('￥0')).toBe(0)
      expect(parseCurrency('-￥500')).toBe(-50000)
    })

    it('should handle different currency formats', () => {
      expect(parseCurrency('1000')).toBe(100000)
      expect(parseCurrency('1,000.50')).toBe(100050)
      expect(parseCurrency('$1,000')).toBe(100000)
    })

    it('should handle invalid input gracefully', () => {
      expect(parseCurrency('')).toBeNaN()
      expect(parseCurrency('abc')).toBeNaN()
    })
  })

  describe('formatDate', () => {
    it('should format ISO date strings to Japanese format', () => {
      expect(formatDate('2024-01-15T00:00:00Z')).toBe('2024/01/15')
      expect(formatDate('2023-12-31T00:00:00Z')).toBe('2023/12/31')
    })

    it('should handle different date formats', () => {
      expect(formatDate('2024-01-15')).toBe('2024/01/15')
    })
  })

  describe('formatYearMonth', () => {
    it('should format year-month string to Japanese format', () => {
      expect(formatYearMonth('2024-01')).toBe('2024年1月')
      expect(formatYearMonth('2024-12')).toBe('2024年12月')
    })
  })

  describe('getCurrentYearMonth', () => {
    it('should return current year-month in YYYY-MM format', () => {
      const result = getCurrentYearMonth()
      expect(result).toMatch(/^\d{4}-\d{2}$/)
      
      // 現在の日付と比較して確認
      const now = new Date()
      const expected = `${now.getFullYear()}-${(now.getMonth() + 1).toString().padStart(2, '0')}`
      expect(result).toBe(expected)
    })
  })

  describe('getNextYearMonth', () => {
    it('should get next month correctly', () => {
      expect(getNextYearMonth('2024-01')).toBe('2024-02')
      expect(getNextYearMonth('2024-11')).toBe('2024-12')
    })

    it('should handle year rollover', () => {
      expect(getNextYearMonth('2024-12')).toBe('2025-01')
    })
  })

  describe('getPreviousYearMonth', () => {
    it('should get previous month correctly', () => {
      expect(getPreviousYearMonth('2024-02')).toBe('2024-01')
      expect(getPreviousYearMonth('2024-12')).toBe('2024-11')
    })

    it('should handle year rollover', () => {
      expect(getPreviousYearMonth('2024-01')).toBe('2023-12')
    })
  })

  describe('validateRequired', () => {
    it('should validate string values', () => {
      expect(validateRequired('test')).toBe(true)
      expect(validateRequired('')).toBe(false)
      expect(validateRequired('   ')).toBe(false)
    })

    it('should validate number values', () => {
      expect(validateRequired(0)).toBe(true)
      expect(validateRequired(100)).toBe(true)
      expect(validateRequired(-1)).toBe(true)
    })
  })

  describe('validateAmount', () => {
    it('should validate positive amounts', () => {
      expect(validateAmount(0)).toBe(true)
      expect(validateAmount(100)).toBe(true)
      expect(validateAmount(0.01)).toBe(true)
    })

    it('should reject negative amounts', () => {
      expect(validateAmount(-1)).toBe(false)
      expect(validateAmount(-0.01)).toBe(false)
    })
  })

  describe('validateYearMonth', () => {
    it('should validate correct year-month format', () => {
      expect(validateYearMonth('2024-01')).toBe(true)
      expect(validateYearMonth('2024-12')).toBe(true)
      expect(validateYearMonth('1900-01')).toBe(true)
      expect(validateYearMonth('2100-12')).toBe(true)
    })

    it('should reject invalid formats', () => {
      expect(validateYearMonth('2024-1')).toBe(false)
      expect(validateYearMonth('24-01')).toBe(false)
      expect(validateYearMonth('2024/01')).toBe(false)
      expect(validateYearMonth('invalid')).toBe(false)
    })

    it('should reject invalid month values', () => {
      expect(validateYearMonth('2024-00')).toBe(false)
      expect(validateYearMonth('2024-13')).toBe(false)
    })

    it('should reject invalid year values', () => {
      expect(validateYearMonth('1899-01')).toBe(false)
      expect(validateYearMonth('2101-01')).toBe(false)
    })
  })

  describe('validateDay', () => {
    it('should validate valid day numbers', () => {
      expect(validateDay(1)).toBe(true)
      expect(validateDay(15)).toBe(true)
      expect(validateDay(31)).toBe(true)
    })

    it('should reject invalid day numbers', () => {
      expect(validateDay(0)).toBe(false)
      expect(validateDay(32)).toBe(false)
      expect(validateDay(-1)).toBe(false)
    })
  })
})
