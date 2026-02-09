'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { createClient } from '@/lib/supabase/client'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useToast } from '@/hooks/use-toast'
import { LogOut } from 'lucide-react'

export default function SettingsPage() {
  const router = useRouter()
  const { toast } = useToast()
  const supabase = createClient()
  const [loading, setLoading] = useState(false)

  const handleLogout = async () => {
    setLoading(true)
    try {
      await supabase.auth.signOut()
      toast({
        title: 'Logged out',
        description: 'You have been logged out successfully',
      })
      router.push('/login')
      router.refresh()
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to log out',
        variant: 'destructive',
      })
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-6 p-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-white">Settings</h1>
        <p className="text-slate-400 mt-1">Manage your account and preferences</p>
      </div>

      {/* API Configuration */}
      <Card>
        <CardHeader>
          <CardTitle>API Configuration</CardTitle>
          <CardDescription>
            Configure the API endpoint for HARaMA backend
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="api-url">API Base URL</Label>
            <Input
              id="api-url"
              value={process.env.NEXT_PUBLIC_API_URL || ''}
              readOnly
              className="bg-slate-800/50"
            />
            <p className="text-xs text-slate-400 mt-1">
              This is configured in your environment variables. Update{' '}
              <code className="bg-slate-800 px-2 py-1 rounded text-blue-400">
                NEXT_PUBLIC_API_URL
              </code>{' '}
              in your project settings.
            </p>
          </div>
        </CardContent>
      </Card>

      {/* Security */}
      <Card>
        <CardHeader>
          <CardTitle>Security</CardTitle>
          <CardDescription>Manage your authentication</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <p className="text-sm text-slate-400">
            Your account is secured with Supabase Auth using JWT tokens. All API
            requests are authenticated automatically.
          </p>
          <Button
            onClick={handleLogout}
            disabled={loading}
            variant="destructive"
            className="w-full"
          >
            <LogOut className="w-4 h-4 mr-2" />
            {loading ? 'Logging out...' : 'Log Out'}
          </Button>
        </CardContent>
      </Card>

      {/* About */}
      <Card className="bg-blue-950/30 border-blue-500/20">
        <CardHeader>
          <CardTitle className="text-blue-300">About HARaMA</CardTitle>
        </CardHeader>
        <CardContent className="text-sm text-blue-400 space-y-2">
          <p>
            HARaMA (Handwritten Analysis and Mark Allocation) is an AI-powered
            exam grading system that uses machine learning to automatically
            grade student responses.
          </p>
          <p>
            Version: 1.0.0 • Built with Next.js, Supabase, and Gemini Vision
            API
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
