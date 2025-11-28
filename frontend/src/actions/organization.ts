'use server'

import { cookies } from 'next/headers'
import { revalidatePath } from 'next/cache'

interface ActionState {
  error?: string
  success?: boolean
  message?: string
  data?: unknown
}

export async function createOrganizationAction(
  prevState: ActionState | null,
  formData: FormData
): Promise<ActionState> {
  const name = formData.get('name') as string
  const description = formData.get('description') as string

  // Validation
  if (!name) {
    return {
      error: 'Organization name is required',
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

    // Call backend API
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/v1/organizations`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ name, description }),
    })

    if (!response.ok) {
      const error = await response.json()
      return {
        error: error.message || 'Failed to create organization',
        success: false,
      }
    }

    const data = await response.json()

    // Revalidate the dashboard
    revalidatePath('/dashboard')

    return {
      success: true,
      message: 'Organization created successfully',
      data,
    }
  } catch (error) {
    console.error('Create organization error:', error)
    return {
      error: error instanceof Error ? error.message : 'An unexpected error occurred',
      success: false,
    }
  }
}

export async function updateOrganizationAction(
  prevState: ActionState | null,
  formData: FormData
): Promise<ActionState> {
  const id = formData.get('id') as string
  const name = formData.get('name') as string
  const description = formData.get('description') as string

  // Validation
  if (!id) {
    return {
      error: 'Organization ID is required',
      success: false,
    }
  }

  if (!name) {
    return {
      error: 'Organization name is required',
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

    // Call backend API
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/v1/organizations/${id}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ name, description }),
    })

    if (!response.ok) {
      const error = await response.json()
      return {
        error: error.message || 'Failed to update organization',
        success: false,
      }
    }

    const data = await response.json()

    // Revalidate the dashboard
    revalidatePath('/dashboard')

    return {
      success: true,
      message: 'Organization updated successfully',
      data,
    }
  } catch (error) {
    console.error('Update organization error:', error)
    return {
      error: error instanceof Error ? error.message : 'An unexpected error occurred',
      success: false,
    }
  }
}

export async function inviteMemberAction(
  prevState: ActionState | null,
  formData: FormData
): Promise<ActionState> {
  const organizationId = formData.get('organizationId') as string
  const email = formData.get('email') as string
  const role = formData.get('role') as string

  // Validation
  if (!organizationId || !email || !role) {
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

  try {
    const cookieStore = await cookies()
    const token = cookieStore.get('auth-token')?.value

    if (!token) {
      return {
        error: 'Not authenticated',
        success: false,
      }
    }

    // Call backend API
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/v1/organizations/${organizationId}/invite`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, role }),
    })

    if (!response.ok) {
      const error = await response.json()
      return {
        error: error.message || 'Failed to invite member',
        success: false,
      }
    }

    return {
      success: true,
      message: 'Invitation sent successfully',
    }
  } catch (error) {
    console.error('Invite member error:', error)
    return {
      error: error instanceof Error ? error.message : 'An unexpected error occurred',
      success: false,
    }
  }
}
