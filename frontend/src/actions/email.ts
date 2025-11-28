'use server'

import { cookies } from 'next/headers'
import { revalidatePath } from 'next/cache'

interface ActionState {
  error?: string
  success?: boolean
  message?: string
}

export async function syncEmailsAction(
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  _prevState: ActionState | null
): Promise<ActionState> {
  try {
    const cookieStore = await cookies()
    const token = cookieStore.get('auth-token')?.value

    if (!token) {
      return {
        error: 'Not authenticated',
        success: false,
      }
    }

    // Call backend sync API
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/v1/sync`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    })

    if (!response.ok) {
      const error = await response.json()
      return {
        error: error.message || 'Sync failed',
        success: false,
      }
    }

    // Revalidate the inbox page
    revalidatePath('/dashboard/inbox')

    return {
      success: true,
      message: 'Email sync started successfully',
    }
  } catch (error) {
    console.error('Sync error:', error)
    return {
      error: error instanceof Error ? error.message : 'An unexpected error occurred',
      success: false,
    }
  }
}

export async function deleteEmailAction(
  prevState: ActionState | null,
  formData: FormData
): Promise<ActionState> {
  const emailId = formData.get('emailId') as string

  if (!emailId) {
    return {
      error: 'Email ID is required',
      success: false,
    }
  }

  try {
    const cookieStore = await cookies()
    const token = cookieStore.get('auth-token')?.value

    if (!token) {
      return {
        error: 'Not authenticated',
        success: false,
      }
    }

    // Call backend delete API
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/v1/emails/${emailId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })

    if (!response.ok) {
      const error = await response.json()
      return {
        error: error.message || 'Delete failed',
        success: false,
      }
    }

    // Revalidate the inbox page
    revalidatePath('/dashboard/inbox')

    return {
      success: true,
      message: 'Email deleted successfully',
    }
  } catch (error) {
    console.error('Delete error:', error)
    return {
      error: error instanceof Error ? error.message : 'An unexpected error occurred',
      success: false,
    }
  }
}

export async function archiveEmailAction(
  prevState: ActionState | null,
  formData: FormData
): Promise<ActionState> {
  const emailId = formData.get('emailId') as string

  if (!emailId) {
    return {
      error: 'Email ID is required',
      success: false,
    }
  }

  try {
    const cookieStore = await cookies()
    const token = cookieStore.get('auth-token')?.value

    if (!token) {
      return {
        error: 'Not authenticated',
        success: false,
      }
    }

    // Call backend archive API
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/v1/emails/${emailId}/archive`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })

    if (!response.ok) {
      const error = await response.json()
      return {
        error: error.message || 'Archive failed',
        success: false,
      }
    }

    // Revalidate the inbox page
    revalidatePath('/dashboard/inbox')

    return {
      success: true,
      message: 'Email archived successfully',
    }
  } catch (error) {
    console.error('Archive error:', error)
    return {
      error: error instanceof Error ? error.message : 'An unexpected error occurred',
      success: false,
    }
  }
}
