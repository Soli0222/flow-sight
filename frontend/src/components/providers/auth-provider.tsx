'use client'

import React, { createContext, useContext, useEffect, useState, useCallback } from 'react'
import Cookies from 'js-cookie'

interface User {
  id: string
  email: string
  name: string
  picture: string
}

interface AuthContextType {
  user: User | null
  token: string | null
  login: (token: string, user: User) => void
  logout: () => void
  isLoading: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  const fetchUserInfo = useCallback(async (authToken: string, retryCount = 0) => {
    try {
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || ''
      const response = await fetch(`${apiUrl}/api/v1/auth/me`, {
        headers: {
          'Authorization': `Bearer ${authToken}`,
        },
      })

      if (response.ok) {
        const userData = await response.json()
        setUser(userData)
        localStorage.setItem('user_data', JSON.stringify(userData))
      } else {
        const errorText = await response.text()
        console.error('Failed to fetch user info:', response.status, errorText)
        
        // 初回ログイン時のデータベース一貫性問題に対してリトライを実装
        if (response.status === 500 && retryCount < 3 && errorText.includes('no rows in result set')) {
          console.log(`Retrying to fetch user info (attempt ${retryCount + 1}/3)...`)
          setTimeout(() => {
            fetchUserInfo(authToken, retryCount + 1)
          }, 1000 * (retryCount + 1)) // 1秒、2秒、3秒の間隔でリトライ
          return
        }
        
        // Token is invalid, remove it
        Cookies.remove('auth_token')
        localStorage.removeItem('user_data')
        setToken(null)
      }
    } catch (error) {
      console.error('Failed to fetch user info:', error)
      
      // ネットワークエラーの場合もリトライを試行
      if (retryCount < 3) {
        console.log(`Retrying to fetch user info due to network error (attempt ${retryCount + 1}/3)...`)
        setTimeout(() => {
          fetchUserInfo(authToken, retryCount + 1)
        }, 1000 * (retryCount + 1))
        return
      }
      
      Cookies.remove('auth_token')
      localStorage.removeItem('user_data')
      setToken(null)
    } finally {
      if (retryCount === 0) { // 最初の試行でのみ setIsLoading(false) を呼ぶ
        setIsLoading(false)
      }
    }
  }, [])

  useEffect(() => {
    const initializeAuth = async () => {
      // Check for existing token and user data on app start
      const savedToken = Cookies.get('auth_token')
      const savedUser = localStorage.getItem('user_data')
      
      if (savedToken && savedUser) {
        try {
          const userData = JSON.parse(savedUser)
          setToken(savedToken)
          setUser(userData)
          setIsLoading(false)
        } catch (error) {
          console.error('Failed to parse saved user data:', error)
          // If user data is corrupted, fetch from server
          setToken(savedToken)
          await fetchUserInfo(savedToken)
        }
      } else if (savedToken) {
        setToken(savedToken)
        // Fetch user info if we have token but no user data
        await fetchUserInfo(savedToken)
      } else {
        setIsLoading(false)
      }
    }
    
    initializeAuth()
  }, [fetchUserInfo])

  const login = useCallback((authToken: string, userData: User) => {
    setToken(authToken)
    setUser(userData)
    Cookies.set('auth_token', authToken, { expires: 1/24 }) // 1 hour (1/24 days)
    localStorage.setItem('user_data', JSON.stringify(userData))
  }, [])

  const logout = useCallback(() => {
    setToken(null)
    setUser(null)
    Cookies.remove('auth_token')
    localStorage.removeItem('user_data')
  }, [])

  return (
    <AuthContext.Provider value={{ user, token, login, logout, isLoading }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
