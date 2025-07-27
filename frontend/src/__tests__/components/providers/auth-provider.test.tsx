/**
 * @jest-environment jsdom
 */
import { render, screen, waitFor } from '@testing-library/react'
import { AuthProvider, useAuth } from '@/components/providers/auth-provider'
import { mockApiResponses } from '../../__mocks__/api-responses'
import Cookies from 'js-cookie'

// cookieのモック
jest.mock('js-cookie', () => ({
  get: jest.fn(),
  set: jest.fn(),
  remove: jest.fn(),
}))

// fetchのモック
global.fetch = jest.fn()

// localStorageのモック
const localStorageMock = {
  getItem: jest.fn(),
  setItem: jest.fn(),
  removeItem: jest.fn(),
  clear: jest.fn(),
}
Object.defineProperty(window, 'localStorage', {
  value: localStorageMock
})

// テスト用コンポーネント
const TestComponent = () => {
  const { user, token, isLoading } = useAuth()
  
  if (isLoading) {
    return <div>Loading...</div>
  }
  
  return (
    <div>
      <div data-testid="user-name">{user?.name || 'No user'}</div>
      <div data-testid="token">{token || 'No token'}</div>
    </div>
  )
}

describe('AuthProvider', () => {
  const mockFetch = fetch as jest.MockedFunction<typeof fetch>
  const mockCookiesGet = Cookies.get as jest.MockedFunction<typeof Cookies.get>
  const mockLocalStorageGet = localStorageMock.getItem as jest.MockedFunction<typeof localStorage.getItem>
  
  beforeEach(() => {
    jest.clearAllMocks()
    mockFetch.mockClear()
    ;(mockCookiesGet as jest.Mock).mockReturnValue(undefined)
    ;(mockLocalStorageGet as jest.Mock).mockReturnValue(null)
  })

  it('should provide initial state when no token or user data exists', async () => {
    // 明示的に空の状態を設定
    ;(mockCookiesGet as jest.Mock).mockReturnValue(undefined)
    ;(mockLocalStorageGet as jest.Mock).mockReturnValue(null)
    
    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )
    
    // Loading状態が終了するまで待機
    await waitFor(() => {
      expect(screen.queryByText('Loading...')).not.toBeInTheDocument()
    })

    expect(screen.getByTestId('user-name')).toHaveTextContent('No user')
    expect(screen.getByTestId('token')).toHaveTextContent('No token')
  })

  it('should handle successful user fetch', async () => {
    // localStorage とCookieのモック
    const mockToken = 'test-token'
    ;(mockCookiesGet as jest.Mock).mockReturnValue(mockToken)
    ;(mockLocalStorageGet as jest.Mock).mockReturnValue(null) // キャッシュされたユーザーデータなし
    
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => mockApiResponses.userInfo,
    } as Response)

    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )

    await waitFor(() => {
      expect(screen.queryByText('Loading...')).not.toBeInTheDocument()
    })

    expect(screen.getByTestId('user-name')).toHaveTextContent(mockApiResponses.userInfo.name)
    expect(screen.getByTestId('token')).toHaveTextContent(mockToken)
  })

  it('should handle failed user fetch', async () => {
    ;(mockCookiesGet as jest.Mock).mockReturnValue('invalid-token')
    ;(mockLocalStorageGet as jest.Mock).mockReturnValue(null) // localStorageにユーザーデータがない状態
    
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 401,
      text: async () => 'Unauthorized',
    } as Response)

    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )

    await waitFor(() => {
      expect(screen.queryByText('Loading...')).not.toBeInTheDocument()
    })

    expect(screen.getByTestId('user-name')).toHaveTextContent('No user')
    expect(screen.getByTestId('token')).toHaveTextContent('No token')
  })

  it('should handle no token case', async () => {
    ;(mockCookiesGet as jest.Mock).mockReturnValue(undefined)
    ;(mockLocalStorageGet as jest.Mock).mockReturnValue(null)

    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )

    await waitFor(() => {
      expect(screen.queryByText('Loading...')).not.toBeInTheDocument()
    })

    expect(screen.getByTestId('user-name')).toHaveTextContent('No user')
    expect(screen.getByTestId('token')).toHaveTextContent('No token')
    expect(mockFetch).not.toHaveBeenCalled()
  })
})
