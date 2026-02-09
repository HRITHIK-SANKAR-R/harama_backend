'use client'

import { useEffect, useState } from 'react'
import { createClient } from '@/lib/supabase/client'
import { listExams, getGradingTrends } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { useToast } from '@/hooks/use-toast'
import Link from 'next/link'
import {
  BarChart3,
  BookOpen,
  CheckCircle2,
  FileText,
  Plus,
  TrendingUp,
} from 'lucide-react'

export default function DashboardPage() {
  const [exams, setExams] = useState<any[]>([])
  const [trends, setTrends] = useState<any>(null)
  const [loading, setLoading] = useState(true)
  const [user, setUser] = useState<any>(null)
  const { toast } = useToast()

  useEffect(() => {
    const fetchData = async () => {
      try {
        const supabase = createClient()
        const {
          data: { user },
        } = await supabase.auth.getUser()

        if (user) {
          setUser(user)
        }

        const [examsData, trendsData] = await Promise.all([
          listExams().catch(() => []),
          getGradingTrends().catch(() => null),
        ])

        setExams(examsData || [])
        setTrends(trendsData)
      } catch (error: any) {
        toast({
          title: 'Error loading dashboard',
          description: error.message,
          variant: 'destructive',
        })
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [toast])

  const stats = [
    {
      title: 'Total Exams',
      value: exams.length,
      icon: BookOpen,
      color: 'from-blue-500 to-blue-600',
      bgColor: 'bg-blue-500/10',
    },
    {
      title: 'Pending Reviews',
      value: trends?.needs_review || 0,
      icon: FileText,
      color: 'from-purple-500 to-purple-600',
      bgColor: 'bg-purple-500/10',
    },
    {
      title: 'Auto Graded',
      value: trends?.auto_graded || 0,
      icon: CheckCircle2,
      color: 'from-blue-600 to-blue-500',
      bgColor: 'bg-blue-600/10',
    },
    {
      title: 'Avg Confidence',
      value: trends?.average_confidence
        ? `${(trends.average_confidence * 100).toFixed(0)}%`
        : 'N/A',
      icon: TrendingUp,
      color: 'from-purple-600 to-purple-500',
      bgColor: 'bg-purple-600/10',
    },
  ]

  return (
    <div className="space-y-8 p-6 bg-gradient-to-br from-slate-950 to-slate-900 min-h-screen relative overflow-hidden">
      {/* Subtle background */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none -z-10">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-blue-600 rounded-full mix-blend-multiply filter blur-3xl opacity-10 animate-blob" />
        <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-purple-600 rounded-full mix-blend-multiply filter blur-3xl opacity-10 animate-blob animation-delay-2000" />
      </div>

      {/* Welcome Section */}
      <div className="space-y-2 relative z-10">
        <h1 className="text-4xl font-bold text-white">Dashboard</h1>
        <p className="text-slate-400">
          Welcome back, <span className="font-semibold text-white">{user?.user_metadata?.full_name || user?.email}</span>
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 relative z-10">
        {stats.map((stat, idx) => {
          const Icon = stat.icon
          return (
            <Card 
              key={stat.title} 
              className="border border-slate-700 bg-slate-900/50 backdrop-blur-xl shadow-lg hover:shadow-blue-500/20 transition-all duration-300 group overflow-hidden"
              style={{ animationDelay: `${idx * 100}ms` }}
            >
              <div className={`absolute top-0 right-0 w-28 h-28 bg-gradient-to-br ${stat.color} opacity-5 rounded-full -mr-10 -mt-10 group-hover:opacity-10 transition-all duration-300`} />
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-semibold text-slate-300">{stat.title}</CardTitle>
                <div className={`p-2 rounded-lg ${stat.bgColor}`}>
                  <Icon className={`w-4 h-4 text-transparent bg-clip-text bg-gradient-to-br ${stat.color}`} />
                </div>
              </CardHeader>
              <CardContent>
                <div className="text-3xl font-bold text-white">{stat.value}</div>
                <p className="text-xs text-slate-500 mt-2">Real-time</p>
              </CardContent>
            </Card>
          )
        })}
      </div>

      {/* Recent Exams Section */}
      <div className="space-y-4 relative z-10">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-2xl font-bold text-white">Recent Exams</h2>
            <p className="text-slate-400 text-sm mt-1">Manage and monitor your assessments</p>
          </div>
          <Link href="/exams/new">
            <Button className="bg-blue-600 hover:bg-blue-700 text-white font-semibold shadow-lg shadow-blue-500/20 transition-all duration-300 h-9">
              <Plus className="w-4 h-4 mr-2" />
              New Exam
            </Button>
          </Link>
        </div>

        {loading ? (
          <Card className="border border-slate-700 bg-slate-900/50 backdrop-blur-xl shadow-lg">
            <CardContent className="flex flex-col items-center justify-center py-12">
              <div className="animate-spin rounded-full h-10 w-10 border-2 border-slate-700 border-t-blue-500" />
              <p className="text-slate-400 mt-3 text-sm">Loading exams...</p>
            </CardContent>
          </Card>
        ) : exams.length === 0 ? (
          <Card className="border border-slate-700 bg-slate-900/50 backdrop-blur-xl shadow-lg">
            <CardContent className="flex flex-col items-center justify-center py-12">
              <div className="w-12 h-12 bg-blue-500/10 rounded-lg flex items-center justify-center mb-3">
                <BookOpen className="w-6 h-6 text-blue-400" />
              </div>
              <h3 className="text-base font-semibold text-white">No exams yet</h3>
              <p className="text-slate-400 text-sm mt-2 max-w-xs text-center">
                Create your first exam to start grading
              </p>
              <Link href="/exams/new" className="mt-4">
                <Button className="bg-blue-600 hover:bg-blue-700 text-white font-semibold text-sm">
                  Create Exam
                </Button>
              </Link>
            </CardContent>
          </Card>
        ) : (
          <div className="space-y-3">
            {exams.slice(0, 5).map((exam) => (
              <Card 
                key={exam.id} 
                className="border border-slate-700 bg-slate-900/50 backdrop-blur-xl hover:border-blue-500/50 hover:shadow-blue-500/10 hover:shadow-lg transition-all duration-300 group cursor-pointer"
              >
                <CardContent className="py-4">
                  <Link href={`/exams/${exam.id}`}>
                    <div className="flex items-center justify-between">
                      <div className="flex-1">
                        <h3 className="font-semibold text-white group-hover:text-blue-400 transition-colors">{exam.title}</h3>
                        <p className="text-sm text-slate-400 mt-1 flex gap-3">
                          <span>{exam.subject}</span>
                          <span>•</span>
                          <span>{exam.questions?.length || 0} questions</span>
                        </p>
                      </div>
                      <Badge className="ml-4 bg-blue-600/20 text-blue-300 border border-blue-500/30 font-medium text-xs">
                        {new Date(exam.created_at).toLocaleDateString()}
                      </Badge>
                    </div>
                  </Link>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
