'use client'

import { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import { getGradingTrends, getExam } from '@/lib/api'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { useToast } from '@/hooks/use-toast'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
  LineChart,
  Line,
} from 'recharts'

const COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444']

export default function AnalyticsPage() {
  const params = useParams()
  const examId = params.id as string
  const { toast } = useToast()

  const [exam, setExam] = useState<any>(null)
  const [trends, setTrends] = useState<any>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchData()
  }, [examId])

  const fetchData = async () => {
    try {
      setLoading(true)
      const [examData, trendsData] = await Promise.all([
        getExam(examId).catch(() => null),
        getGradingTrends(examId).catch(() => null),
      ])

      setExam(examData)
      setTrends(trendsData)
    } catch (error: any) {
      toast({
        title: 'Error loading analytics',
        description: error.message,
        variant: 'destructive',
      })
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600" />
      </div>
    )
  }

  // Sample data for demonstration
  const statusData = trends
    ? [
        {
          name: 'Auto Graded',
          value: trends.auto_graded || 0,
          color: COLORS[1],
        },
        {
          name: 'Needs Review',
          value: trends.needs_review || 0,
          color: COLORS[2],
        },
        {
          name: 'Overridden',
          value: trends.overridden || 0,
          color: COLORS[3],
        },
      ]
    : []

  const trendData = trends?.trends || []

  return (
    <div className="space-y-6 p-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-white">Analytics</h1>
        <p className="text-slate-400 mt-1">{exam?.title}</p>
      </div>

      {/* Summary Stats */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Total Submissions</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{trends?.total_submissions || 0}</div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Auto Graded</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              {trends?.auto_graded || 0}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Needs Review</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-amber-600">
              {trends?.needs_review || 0}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Avg Confidence</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-blue-600">
              {trends?.average_confidence
                ? `${(trends.average_confidence * 100).toFixed(0)}%`
                : 'N/A'}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Grading Status Distribution */}
      {statusData.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Grading Status Distribution</CardTitle>
            <CardDescription>
              Breakdown of submissions by status
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={statusData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, value }) => `${name}: ${value}`}
                  outerRadius={80}
                  fill="#8884d8"
                  dataKey="value"
                >
                  {statusData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      )}

      {/* Score Trend */}
      {trendData.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Score Trends Over Time</CardTitle>
            <CardDescription>Average score per submission date</CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <LineChart data={trendData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
                <XAxis dataKey="date" stroke="#94a3b8" />
                <YAxis stroke="#94a3b8" />
                <Tooltip contentStyle={{ backgroundColor: '#0f172a', border: '1px solid #334155', borderRadius: '8px', color: '#e2e8f0' }} />
                <Legend />
                <Line
                  type="monotone"
                  dataKey="avg_score"
                  stroke="#3b82f6"
                  dot={{ fill: '#3b82f6' }}
                  name="Avg Score"
                />
              </LineChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      )}

      {/* Summary */}
      <Card>
        <CardHeader>
          <CardTitle>Summary</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          {trends ? (
            <>
              <div className="flex justify-between">
                <span className="text-slate-400">Total Submissions</span>
                <span className="font-medium">{trends.total_submissions}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-slate-400">Average Score</span>
                <span className="font-medium">
                  {trends.average_score?.toFixed(2) || 'N/A'}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-slate-400">Average Confidence</span>
                <span className="font-medium">
                  {trends.average_confidence
                    ? `${(trends.average_confidence * 100).toFixed(0)}%`
                    : 'N/A'}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-slate-400">Status Breakdown</span>
                <div className="flex gap-2">
                  <Badge variant="outline">
                    {trends.auto_graded} Auto Graded
                  </Badge>
                  <Badge variant="outline">
                    {trends.needs_review} Needs Review
                  </Badge>
                </div>
              </div>
            </>
          ) : (
            <p className="text-slate-400">No analytics data available yet</p>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
