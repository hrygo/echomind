'use server'

import { cookies } from 'next/headers'
import { redirect } from 'next/navigation'

interface AuthResponse {
  token: string
  user: {
    id: string
    email: string
    name: string
  }
}

interface ActionState {
  error?: string
  success?: boolean
  message?: string
}

export async function loginAction(
  prevState: ActionState | null,
  formData: FormData
): Promise<ActionState> {
  const email = formData.get('email') as string
  const password = formData.get('password') as string

  // Validation
  if (!email || !password) {
    return {
      error: 'Email and password are required',
      success: false,
    }
  }

  if (!email.includes('@')) {
    return {
      error: 'Invalid email format',
      success: false,
    }
  }

  try {
    // Call backend API
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/v1/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    })

    if (!response.ok) {
      const error = await response.json()
      return {
        error: error.message || 'Login failed',
        success: false,
      }
    }

    const data: AuthResponse = await response.json()

    // Set token in cookie
    const cookieStore = await cookies()
    cookieStore.set('auth-token', data.token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      maxAge: 60 * 60 * 24 * 7, // 7 days
    })

    // Redirect to dashboard
    redirect('/dashboard')
  } catch (error) {
    console.error('Login error:', error)
    return {
      error: error instanceof Error ? error.message : 'An unexpected error occurred',
      success: false,
    }
  }
}

export async function registerAction(
  prevState: ActionState | null,
  formData: FormData
): Promise<ActionState> {
  const email = formData.get('email') as string
  const password = formData.get('password') as string
  const name = formData.get('name') as string

  // Validation
  if (!email || !password || !name) {
    return {
      error: 'All fields are required',
      success: false,
    }
  }

  if (!email.includes('@')) {
    return {
      error: 'Invalid email format',
      success: false,
    }
  }

  if (password.length < 8) {
    return {
      error: 'Password must be at least 8 characters',
      success: false,
    }
  }

  try {
    // Call backend API
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/v1/auth/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password, name }),
    })

    if (!response.ok) {
      const error = await response.json()
      return {
        error: error.message || 'Registration failed',
        success: false,
      }
    }

    const data: AuthResponse = await response.json()

    // Set token in cookie
    const cookieStore = await cookies()
    cookieStore.set('auth-token', data.token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      maxAge: 60 * 60 * 24 * 7, // 7 days
    })

    // Redirect to onboarding
    redirect('/onboarding')
  } catch (error) {
    console.error('Registration error:', error)
    return {
      error: error instanceof Error ? error.message : 'An unexpected error occurred',
      success: false,
    }
  }
}

export async function logoutAction(): Promise<void> {
  const cookieStore = await cookies()
  cookieStore.delete('auth-token')
  redirect('/auth')
}
