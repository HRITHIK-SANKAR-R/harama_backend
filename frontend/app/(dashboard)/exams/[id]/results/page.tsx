'use client'

import { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import { getExam, exportGrades } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { useToast } from '@/hooks/use-toast'
import Link from 'next/link'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Download, FileText } from 'lucide-react'

export default function ResultsPage() {
  const params = useParams()
  const examId = params.id as string
  const { toast } = useToast()

  const [exam, setExam] = useState<any>(null)
  const [loading, setLoading] = useState(true)
  const [exporting, setExporting] = useState(false)

  useEffect(() => {
    fetchExam()
  }, [examId])

  const fetchExam = async () => {
    try {
      setLoading(true)
      const data = await getExam(examId)
      setExam(data)
    } catch (error: any) {
      toast({
        title: 'Error loading exam',
        description: error.message,
        variant: 'destructive',
      })
    } finally {
      setLoading(false)
    }
  }

  const handleExport = async (format: string) => {
    setExporting(true)
    try {
      const blob = await exportGrades(examId, format)
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `grades.${format === 'csv' ? 'csv' : format === 'xlsx' ? 'xlsx' : 'pdf'}`
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)

      toast({
        title: 'Success',
        description: `Grades exported as ${format.toUpperCase()}`,
      })
    } catch (error: any) {
      toast({
        title: 'Export failed',
        description: error.message,
        variant: 'destructive',
      })
    } finally {
      setExporting(false)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600" />
      </div>
    )
  }

  if (!exam) {
    return (
      <div className="p-6">
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <h3 className="text-lg font-medium text-white">Exam not found</h3>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="space-y-6 p-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-white">Results</h1>
          <p className="text-slate-400 mt-1">{exam.title}</p>
        </div>
        <div className="flex gap-2">
          <Button
            onClick={() => handleExport('csv')}
            disabled={exporting}
            variant="outline"
          >
            <Download className="w-4 h-4 mr-2" />
            CSV
          </Button>
          <Button
            onClick={() => handleExport('xlsx')}
            disabled={exporting}
            variant="outline"
          >
            <Download className="w-4 h-4 mr-2" />
            Excel
          </Button>
        </div>
      </div>

      {/* Summary */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Total Questions</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{exam.questions?.length || 0}</div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Total Points</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {exam.questions?.reduce((sum: number, q: any) => sum + q.points, 0) || 0}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Questions</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-1">
              {exam.questions?.map((q: any, i: number) => (
                <div key={q.id} className="text-xs text-slate-400">
                  Q{i + 1}: {q.points} points
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Questions Breakdown */}
      <Card>
        <CardHeader>
          <CardTitle>Questions Details</CardTitle>
          <CardDescription>Points and grading information per question</CardDescription>
        </CardHeader>
        <CardContent>
          {exam.questions?.length === 0 ? (
            <p className="text-slate-400">No questions in this exam</p>
          ) : (
            <div className="rounded-lg border border-slate-700 overflow-hidden">
              <Table>
                <TableHeader>
                  <TableRow className="bg-slate-800/50">
                    <TableHead>Question</TableHead>
                    <TableHead>Type</TableHead>
                    <TableHead>Points</TableHead>
                    <TableHead>Group</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {exam.questions?.map((q: any, i: number) => (
                    <TableRow key={q.id} className="hover:bg-slate-800/30 border-slate-700">
                      <TableCell className="font-medium">
                        {q.question_text.substring(0, 50)}...
                      </TableCell>
                      <TableCell>
                        <Badge variant="outline">{q.answer_type}</Badge>
                      </TableCell>
                      <TableCell>{q.points}</TableCell>
                      <TableCell>
                        {q.question_group ? (
                          <Badge variant="secondary">{q.question_group}</Badge>
                        ) : (
                          <span className="text-slate-500">-</span>
                        )}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Instructions */}
      <Card className="bg-blue-950/30 border-blue-500/20">
        <CardContent className="pt-6">
          <h4 className="font-medium text-blue-300 mb-3 flex items-center gap-2">
            <FileText className="w-4 h-4" />
            How to view student results
          </h4>
          <ol className="text-sm text-blue-400 space-y-2 list-decimal list-inside">
            <li>Go to the exam details page</li>
            <li>Upload student answer sheets</li>
            <li>Trigger grading for each submission</li>
            <li>Review AI grades and override if needed</li>
            <li>Export results when ready</li>
          </ol>
        </CardContent>
      </Card>
    </div>
  )
}
