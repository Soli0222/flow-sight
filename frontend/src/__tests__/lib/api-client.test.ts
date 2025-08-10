/**
 * @jest-environment jsdom
 */
import ApiClient from '@/lib/api-client'
import { mockApiResponses } from '../__mocks__/api-responses'

// fetchのモック
global.fetch = jest.fn()

describe('ApiClient', () => {
  let apiClient: ApiClient
  const mockFetch = fetch as jest.MockedFunction<typeof fetch>

  beforeEach(() => {
    jest.clearAllMocks()
    apiClient = new ApiClient()
  })

  afterEach(() => {
    jest.resetAllMocks()
  })

  describe('request method', () => {
    it('should make a successful GET request', async () => {
      const mockData = { test: 'data' }
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => mockData,
      } as Response)

      const result = await apiClient['request']('/test')
      
      expect(fetch).toHaveBeenCalledWith('/api/v1/test', {
        headers: {
          'Content-Type': 'application/json',
        },
      })
      expect(result).toEqual(mockData)
    })

    it('should handle 204 No Content responses', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 204,
      } as Response)

      const result = await apiClient['request']('/test')
      expect(result).toBeUndefined()
    })

    it('should handle other HTTP errors', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
      } as Response)

      await expect(apiClient['request']('/test')).rejects.toThrow('HTTP error! status: 500')
    })
  })

  describe('Credit Cards API', () => {
    it('should get all credit cards', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => mockApiResponses.creditCards,
      } as Response)

      const result = await apiClient.getCreditCards()
      
      expect(fetch).toHaveBeenCalledWith('/api/v1/credit-cards', expect.any(Object))
      expect(result).toEqual(mockApiResponses.creditCards)
    })

    it('should get a single credit card', async () => {
      const creditCard = mockApiResponses.creditCards[0]
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => creditCard,
      } as Response)

      const result = await apiClient.getCreditCard('1')
      
      expect(fetch).toHaveBeenCalledWith('/api/v1/credit-cards/1', expect.any(Object))
      expect(result).toEqual(creditCard)
    })

    it('should create a credit card', async () => {
      const newCreditCard = {
        name: 'New Card',
        bank_account: 'bank1',
        payment_day: 10,
      }
      const createdCard = { ...newCreditCard, id: '2', created_at: '2024-01-01', updated_at: '2024-01-01' }
      
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 201,
        json: async () => createdCard,
      } as Response)

      const result = await apiClient.createCreditCard(newCreditCard)
      
      expect(fetch).toHaveBeenCalledWith('/api/v1/credit-cards', {
        method: 'POST',
        body: JSON.stringify(newCreditCard),
        headers: {
          'Content-Type': 'application/json',
        },
      })
      expect(result).toEqual(createdCard)
    })

    it('should delete a credit card', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 204,
      } as Response)

      await apiClient.deleteCreditCard('1')
      
      expect(fetch).toHaveBeenCalledWith('/api/v1/credit-cards/1', {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
        },
      })
    })
  })

  describe('Bank Accounts API', () => {
    it('should get all bank accounts', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => mockApiResponses.bankAccounts,
      } as Response)

      const result = await apiClient.getBankAccounts()
      
      expect(fetch).toHaveBeenCalledWith('/api/v1/bank-accounts', expect.any(Object))
      expect(result).toEqual(mockApiResponses.bankAccounts)
    })
  })
})
