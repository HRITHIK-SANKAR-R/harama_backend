/**
 * Centralized environment variable validation.
 * Import this wherever env vars are needed to get type-safe,
 * validated values with clear error messages on misconfiguration.
 */

function missing(name: string) {
  throw new Error(
    `Missing required environment variable: ${name}. ` +
      `Please add it to your .env.local file. See .env.example for reference.`
  )
}

const NEXT_PUBLIC_SUPABASE_URL = process.env.NEXT_PUBLIC_SUPABASE_URL || missing('NEXT_PUBLIC_SUPABASE_URL')
const NEXT_PUBLIC_SUPABASE_ANON_KEY =
  process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY || missing('NEXT_PUBLIC_SUPABASE_ANON_KEY')
const NEXT_PUBLIC_API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export const env = {
  /** Supabase project URL (e.g. https://xxx.supabase.co) */
  NEXT_PUBLIC_SUPABASE_URL,
  /** Supabase anonymous/public key */
  NEXT_PUBLIC_SUPABASE_ANON_KEY,
  /** Backend API base URL */
  NEXT_PUBLIC_API_URL,
} as const
