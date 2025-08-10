'use client'
import React, { createContext, useContext } from 'react'

interface AuthContextType {
  user: null
  token: null
  login: () => void
  logout: () => void
  isLoading: false
}

const AuthContext = createContext<AuthContextType>({
  user: null,
  token: null,
  login: () => {},
  logout: () => {},
  isLoading: false,
})

export function AuthProvider({ children }: { children: React.ReactNode }) {
  return <>{children}</>
}

export function useAuth() {
  return useContext(AuthContext)
}
