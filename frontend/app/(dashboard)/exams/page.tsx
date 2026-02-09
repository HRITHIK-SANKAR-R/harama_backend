'use client'

import { useEffect, useState } from 'react'
import { listExams } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { useToast } from '@/hooks/use-toast'
import Link from 'next/link'
import { Plus, BookOpen, ChevronRight } from 'lucide-react'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

export default function ExamsPage() {
  const [exams, setExams] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [searchTerm, setSearchTerm] = useState('')
  const { toast } = useToast()

  useEffect(() => {
    fetchExams()
  }, [])

  const fetchExams = async () => {
    try {
      setLoading(true)
      const data = await listExams()
      setExams(data || [])
    } catch (error: any) {
      toast({
        title: 'Error loading exams',
        description: error.message,
        variant: 'destructive',
      })
    } finally {
      setLoading(false)
    }
  }

  const filteredExams = exams.filter((exam) =>
    exam.title.toLowerCase().includes(searchTerm.toLowerCase())
  )

  return (
    <div className="space-y-6 p-6 bg-gradient-to-br from-slate-950 to-slate-900 min-h-screen relative overflow-hidden">
      {/* Subtle background */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none -z-10">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-blue-600 rounded-full mix-blend-multiply filter blur-3xl opacity-10 animate-blob" />
        <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-purple-600 rounded-full mix-blend-multiply filter blur-3xl opacity-10 animate-blob animation-delay-2000" />
      </div>

      {/* Header */}
      <div className="flex items-center justify-between relative z-10">
        <div className="space-y-2">
          <h1 className="text-4xl font-bold text-white">Exams</h1>
          <p className="text-slate-400 text-sm">Manage all your assessments</p>
        </div>
        <Link href="/exams/new">
          <Button size="lg" className="bg-blue-600 hover:bg-blue-700 text-white font-semibold shadow-lg shadow-blue-500/20 transition-all duration-300 h-10">
            <Plus className="w-4 h-4 mr-2" />
            New Exam
          </Button>
        </Link>
      </div>

      {/* Search */}
      <div className="relative z-10">
        <input
          type="text"
          placeholder="Search exams..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="w-full px-4 py-2 pl-10 bg-slate-800 border border-slate-700 rounded-lg focus:ring-1 focus:ring-blue-500 focus:border-blue-500 outline-none text-white placeholder:text-slate-500 transition-all duration-300"
        />
        <svg className="absolute left-3 top-2.5 w-5 h-5 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>

      {/* Content */}
      {loading ? (
        <Card className="border border-slate-700 bg-slate-900/50 backdrop-blur-xl shadow-lg relative z-10">
          <CardContent className="flex flex-col items-center justify-center py-12">
            <div className="animate-spin rounded-full h-10 w-10 border-2 border-slate-700 border-t-blue-500" />
            <p className="text-slate-400 mt-3 text-sm">Loading exams...</p>
          </CardContent>
        </Card>
      ) : filteredExams.length === 0 ? (
        <Card className="border border-slate-700 bg-slate-900/50 backdrop-blur-xl shadow-lg relative z-10">
          <CardContent className="flex flex-col items-center justify-center py-12">
            <div className="w-12 h-12 bg-blue-500/10 rounded-lg flex items-center justify-center mb-3">
              <BookOpen className="w-6 h-6 text-blue-400" />
            </div>
            <h3 className="text-base font-semibold text-white">
              {exams.length === 0 ? 'No exams yet' : 'No exams found'}
            </h3>
            <p className="text-slate-400 text-sm mt-2 max-w-xs text-center">
              {exams.length === 0
                ? 'Create your first exam to start grading'
                : 'Try adjusting your search'}
            </p>
            {exams.length === 0 && (
              <Link href="/exams/new" className="mt-4">
                <Button className="bg-blue-600 hover:bg-blue-700 text-white font-semibold text-sm">
                  Create Exam
                </Button>
              </Link>
            )}
          </CardContent>
        </Card>
      ) : (
        <div className="rounded-lg border border-slate-700 shadow-lg overflow-hidden bg-slate-900/50 backdrop-blur-xl relative z-10">
          <Table>
            <TableHeader>
              <TableRow className="border-b border-slate-700">
                <TableHead className="text-slate-300 font-semibold text-sm">Title</TableHead>
                <TableHead className="text-slate-300 font-semibold text-sm">Subject</TableHead>
                <TableHead className="text-slate-300 font-semibold text-sm">Questions</TableHead>
                <TableHead className="text-slate-300 font-semibold text-sm">Created</TableHead>
                <TableHead className="text-right text-slate-300 font-semibold text-sm">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredExams.map((exam) => (
                <TableRow 
                  key={exam.id} 
                  className="border-b border-slate-700 hover:bg-slate-800/50 transition-colors duration-200"
                >
                  <TableCell className="font-medium text-white">{exam.title}</TableCell>
                  <TableCell className="text-slate-400">{exam.subject}</TableCell>
                  <TableCell>
                    <Badge className="bg-blue-600/20 text-blue-300 border border-blue-500/30 font-medium text-xs">
                      {exam.questions?.length || 0}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-slate-400 text-sm">
                    {new Date(exam.created_at).toLocaleDateString()}
                  </TableCell>
                  <TableCell className="text-right">
                    <Link href={`/exams/${exam.id}`}>
                      <Button variant="ghost" size="sm" className="text-blue-400 hover:text-blue-300 hover:bg-blue-500/10">
                        View
                        <ChevronRight className="w-4 h-4 ml-1" />
                      </Button>
                    </Link>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}
    </div>
  )
}
