'use client'

import { useEffect, Suspense, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { useAuth } from '@/components/providers/auth-provider'

function AuthCallbackContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { login } = useAuth()
  const [processed, setProcessed] = useState(false)

  useEffect(() => {
    if (processed) return // 既に処理済みの場合は何もしない

    const handleCallback = async () => {
      setProcessed(true) // 処理開始をマーク
      
      const error = searchParams.get('error')

      if (error) {
        console.error('OAuth error:', error)
        router.push('/login?error=oauth_failed')
        return
      }

      // まず、バックエンドからリダイレクトされた場合のトークンとユーザー情報をチェック
      const token = searchParams.get('token')
      const userStr = searchParams.get('user')
      
      if (token && userStr) {
        try {
          const user = JSON.parse(decodeURIComponent(userStr))
          console.log('Logging in with token and user:', { token: token.substring(0, 20) + '...', user })
          login(token, user)
          router.push('/dashboard')
          return
        } catch (parseError) {
          console.error('Failed to parse user data:', parseError)
          router.push('/login?error=parse_failed')
          return
        }
      }

      // トークンがない場合は、Googleからの直接コールバック（codeパラメータ）をチェック
      const code = searchParams.get('code')
      if (!code) {
        console.log('No code or token found, redirecting to login')
        router.push('/login?error=no_code')
        return
      }

      try {
        // 上記で処理できない場合は、直接バックエンドのコールバックURLにリダイレクト
        console.log('Redirecting to backend callback with code:', code.substring(0, 20) + '...')
        const apiUrl = process.env.NEXT_PUBLIC_API_URL || ''
        window.location.href = `${apiUrl}/api/v1/auth/google/callback?code=${code}`
      } catch (error) {
        console.error('Failed to complete OAuth flow:', error)
        router.push('/login?error=network_error')
      }
    }

    // searchParamsが存在する場合のみ実行
    if (searchParams.toString()) {
      handleCallback()
    }
  }, [searchParams, router, login, processed]) // loginを依存配列に含める

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto"></div>
        <p className="mt-4">ログイン処理中...</p>
      </div>
    </div>
  )
}

export default function AuthCallbackPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto"></div>
          <p className="mt-4">読み込み中...</p>
        </div>
      </div>
    }>
      <AuthCallbackContent />
    </Suspense>
  )
}
