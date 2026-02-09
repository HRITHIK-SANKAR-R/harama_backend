'use client'

import { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import { getExam } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { useToast } from '@/hooks/use-toast'
import Link from 'next/link'
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion'
import { ChevronRight, BookOpen, Upload } from 'lucide-react'

export default function ExamDetailsPage() {
  const params = useParams()
  const examId = params.id as string
  const [exam, setExam] = useState<any>(null)
  const [loading, setLoading] = useState(true)
  const { toast } = useToast()

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
            <BookOpen className="w-12 h-12 text-slate-500 mb-4" />
            <h3 className="text-lg font-medium text-white">Exam not found</h3>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="space-y-6 p-6">
      {/* Header */}
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-3xl font-bold text-white">{exam.title}</h1>
          <div className="flex gap-4 mt-2 text-sm text-slate-400">
            <span>Subject: {exam.subject}</span>
            <span>•</span>
            <span>Created: {new Date(exam.created_at).toLocaleDateString()}</span>
          </div>
        </div>
        <Link href={`/exams/${examId}/upload`}>
          <Button>
            <Upload className="w-4 h-4 mr-2" />
            Upload Submissions
          </Button>
        </Link>
      </div>

      {/* Description */}
      {exam.description && (
        <Card>
          <CardContent className="pt-6">
            <p className="text-slate-300">{exam.description}</p>
          </CardContent>
        </Card>
      )}

      {/* Questions */}
      <Card>
        <CardHeader>
          <CardTitle>Questions</CardTitle>
          <CardDescription>
            {exam.questions?.length || 0} questions • {exam.questions?.reduce((sum: number, q: any) => sum + q.points, 0) || 0} total points
          </CardDescription>
        </CardHeader>
        <CardContent>
          {exam.questions?.length === 0 ? (
            <p className="text-slate-400">No questions added yet</p>
          ) : (
            <Accordion type="single" collapsible className="w-full">
              {exam.questions?.map((question: any, i: number) => (
                <AccordionItem key={question.id} value={question.id}>
                  <AccordionTrigger className="hover:no-underline">
                    <div className="flex items-start justify-between flex-1 text-left">
                      <div>
                        <h3 className="font-medium text-white">
                          Q{i + 1}. {question.question_text}
                        </h3>
                        <div className="flex gap-2 mt-2">
                          <Badge variant="outline">{question.answer_type}</Badge>
                          <Badge variant="secondary">{question.points} points</Badge>
                        </div>
                      </div>
                    </div>
                  </AccordionTrigger>
                  <AccordionContent>
                    <div className="space-y-4 pt-4">
                      {question.rubric && (
                        <div>
                          <h4 className="font-medium text-white mb-2">Rubric</h4>
                          {question.rubric.full_credit_criteria && (
                            <div>
                              <p className="text-sm font-medium text-slate-300 mb-2">
                                Full Credit Criteria:
                              </p>
                              <ul className="text-sm space-y-1 ml-4">
                                {question.rubric.full_credit_criteria.map(
                                  (criteria: any, j: number) => (
                                    <li key={j} className="text-slate-400">
                                      • {criteria.description} ({criteria.points} points)
                                    </li>
                                  )
                                )}
                              </ul>
                            </div>
                          )}
                          {question.rubric.key_concepts && (
                            <div className="mt-3">
                              <p className="text-sm font-medium text-slate-300 mb-2">
                                Key Concepts:
                              </p>
                              <div className="flex flex-wrap gap-2">
                                {question.rubric.key_concepts.map((concept: string, j: number) => (
                                  <Badge key={j} variant="outline">
                                    {concept}
                                  </Badge>
                                ))}
                              </div>
                            </div>
                          )}
                        </div>
                      )}
                      {!question.rubric && (
                        <p className="text-sm text-slate-400">
                          No rubric added yet. Add one from the submission grading page.
                        </p>
                      )}
                    </div>
                  </AccordionContent>
                </AccordionItem>
              ))}
            </Accordion>
          )}
        </CardContent>
      </Card>

      {/* Actions */}
      <Card>
        <CardHeader>
          <CardTitle>Actions</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 gap-3">
            <Link href={`/exams/${examId}/upload`}>
              <Button variant="outline" className="w-full bg-transparent">
                <Upload className="w-4 h-4 mr-2" />
                Upload Submissions
              </Button>
            </Link>
            <Link href={`/exams/${examId}/results`}>
              <Button variant="outline" className="w-full bg-transparent">
                View Results
                <ChevronRight className="w-4 h-4 ml-2" />
              </Button>
            </Link>
            <Link href={`/exams/${examId}/analytics`}>
              <Button variant="outline" className="w-full bg-transparent">
                Analytics
                <ChevronRight className="w-4 h-4 ml-2" />
              </Button>
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
