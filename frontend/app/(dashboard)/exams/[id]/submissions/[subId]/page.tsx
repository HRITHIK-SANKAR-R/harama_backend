'use client'

import { useEffect, useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import {
  getSubmission,
  getGrades,
  overrideGrade,
  triggerGrading,
  getFeedback,
} from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { useToast } from '@/hooks/use-toast'
import { ChevronLeft, ChevronRight, AlertCircle } from 'lucide-react'

export default function GradingReviewPage() {
  const params = useParams()
  const router = useRouter()
  const examId = params.id as string
  const submissionId = params.subId as string
  const { toast } = useToast()

  const [submission, setSubmission] = useState<any>(null)
  const [grades, setGrades] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [currentQuestionIdx, setCurrentQuestionIdx] = useState(0)
  const [overrideScores, setOverrideScores] = useState<Record<string, number>>({})
  const [overrideReasons, setOverrideReasons] = useState<Record<string, string>>({})
  const [overriding, setOverriding] = useState<Record<string, boolean>>({})
  const [grading, setGrading] = useState(false)

  useEffect(() => {
    fetchData()
  }, [submissionId])

  const fetchData = async () => {
    try {
      setLoading(true)
      const [subData, gradesData] = await Promise.all([
        getSubmission(submissionId),
        getGrades(submissionId).catch(() => []),
      ])

      setSubmission(subData)
      setGrades(gradesData || [])
    } catch (error: any) {
      toast({
        title: 'Error loading submission',
        description: error.message,
        variant: 'destructive',
      })
    } finally {
      setLoading(false)
    }
  }

  const handleTriggerGrading = async () => {
    setGrading(true)
    try {
      await triggerGrading(submissionId)
      toast({
        title: 'Grading started',
        description: 'This may take a few moments',
      })
      setTimeout(() => fetchData(), 2000)
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.message,
        variant: 'destructive',
      })
    } finally {
      setGrading(false)
    }
  }

  const handleOverrideGrade = async (grade: any) => {
    const newScore = overrideScores[grade.id]
    const reason = overrideReasons[grade.id]

    if (newScore === undefined) {
      toast({
        title: 'Error',
        description: 'Please enter a new score',
        variant: 'destructive',
      })
      return
    }

    if (!reason.trim()) {
      toast({
        title: 'Error',
        description: 'Please provide a reason for the override',
        variant: 'destructive',
      })
      return
    }

    setOverriding((prev) => ({ ...prev, [grade.id]: true }))

    try {
      await overrideGrade(submissionId, grade.question_id, newScore, reason)
      toast({
        title: 'Success',
        description: 'Grade overridden successfully',
      })
      fetchData()
      setOverrideScores((prev) => {
        const updated = { ...prev }
        delete updated[grade.id]
        return updated
      })
      setOverrideReasons((prev) => {
        const updated = { ...prev }
        delete updated[grade.id]
        return updated
      })
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.message,
        variant: 'destructive',
      })
    } finally {
      setOverriding((prev) => ({ ...prev, [grade.id]: false }))
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600" />
      </div>
    )
  }

  if (!submission) {
    return (
      <div className="p-6">
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <h3 className="text-lg font-medium text-white">Submission not found</h3>
          </CardContent>
        </Card>
      </div>
    )
  }

  const currentGrade = grades[currentQuestionIdx]

  return (
    <div className="space-y-6 p-6">
      {/* Header */}
      <div>
        <Button
          variant="ghost"
          onClick={() => router.back()}
          className="mb-4"
        >
          <ChevronLeft className="w-4 h-4 mr-2" />
          Back
        </Button>
        <h1 className="text-3xl font-bold text-white">
          Submission Review
        </h1>
        <div className="flex gap-4 mt-2 text-sm text-slate-400">
          <span>Student: {submission.student_id}</span>
          <span>•</span>
          <span>Status: {submission.processing_status}</span>
        </div>
      </div>

      {/* Status Card */}
      {submission.processing_status === 'pending' && (
        <Card className="bg-blue-950/30 border-blue-500/20">
          <CardContent className="pt-6 flex items-center justify-between">
            <div className="flex items-center gap-3">
              <AlertCircle className="w-5 h-5 text-blue-400" />
              <p className="text-sm text-blue-300">
                This submission hasn't been graded yet.
              </p>
            </div>
            <Button
              onClick={handleTriggerGrading}
              disabled={grading}
              size="sm"
            >
              {grading ? 'Starting...' : 'Trigger Grading'}
            </Button>
          </CardContent>
        </Card>
      )}

      {grades.length === 0 && submission.processing_status !== 'pending' ? (
        <Card className="bg-amber-950/30 border-amber-500/20">
          <CardContent className="pt-6 flex items-center justify-between">
            <div className="flex items-center gap-3">
              <AlertCircle className="w-5 h-5 text-amber-400" />
              <p className="text-sm text-amber-300">
                No grades available yet.
              </p>
            </div>
            <Button
              onClick={handleTriggerGrading}
              disabled={grading}
              size="sm"
            >
              {grading ? 'Starting...' : 'Trigger Grading'}
            </Button>
          </CardContent>
        </Card>
      ) : null}

      {grades.length > 0 && (
        <>
          {/* Answer & OCR */}
          {submission.answers?.[currentQuestionIdx] && (
            <Card>
              <CardHeader>
                <CardTitle>Student Answer</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {submission.ocr_results?.[currentQuestionIdx] && (
                  <div>
                    <p className="text-sm font-medium text-slate-300 mb-2">
                      OCR Confidence:{' '}
                      <Badge variant="outline">
                        {(
                          submission.ocr_results[currentQuestionIdx]
                            .confidence * 100
                        ).toFixed(0)}%
                      </Badge>
                    </p>
                    {submission.ocr_results[currentQuestionIdx].image_url && (
                      <div className="bg-slate-800/50 p-4 rounded-lg">
                        <img
                          src={submission.ocr_results[
                            currentQuestionIdx
                          ].image_url || "/placeholder.svg"}
                          alt="Student answer"
                          className="max-h-96 w-full object-contain"
                        />
                      </div>
                    )}
                  </div>
                )}

                <div>
                  <p className="text-sm font-medium text-slate-300 mb-2">
                    Extracted Text:
                  </p>
                  <p className="text-slate-300 bg-slate-800/50 p-3 rounded-lg">
                    {submission.answers[currentQuestionIdx]?.text ||
                      'No text extracted'}
                  </p>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Grading */}
          {currentGrade && (
            <Card>
              <CardHeader>
                <CardTitle>AI Grade</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-3 gap-4">
                  <div>
                    <p className="text-sm text-slate-400">Score</p>
                    <p className="text-2xl font-bold text-white">
                      {currentGrade.final_score.toFixed(1)}
                    </p>
                  </div>
                  <div>
                    <p className="text-sm text-slate-400">Confidence</p>
                    <p className="text-2xl font-bold text-white">
                      {(currentGrade.confidence * 100).toFixed(0)}%
                    </p>
                  </div>
                  <div>
                    <p className="text-sm text-slate-400">Status</p>
                    <Badge className="mt-1">{currentGrade.status}</Badge>
                  </div>
                </div>

                {currentGrade.reasoning && (
                  <div>
                    <p className="text-sm font-medium text-slate-300 mb-2">
                      AI Reasoning:
                    </p>
                    <p className="text-slate-300 bg-slate-800/50 p-3 rounded-lg">
                      {currentGrade.reasoning}
                    </p>
                  </div>
                )}

                {/* Override Option */}
                {currentGrade.status !== 'final' && (
                  <Dialog>
                    <DialogTrigger asChild>
                      <Button variant="outline" className="w-full bg-transparent">
                        Override Grade
                      </Button>
                    </DialogTrigger>
                    <DialogContent>
                      <DialogHeader>
                        <DialogTitle>Override Grade</DialogTitle>
                        <DialogDescription>
                          Provide a new score and reason for the override
                        </DialogDescription>
                      </DialogHeader>

                      <div className="space-y-4">
                        <div className="space-y-2">
                          <Label htmlFor="new-score">New Score</Label>
                          <Input
                            id="new-score"
                            type="number"
                            step="0.1"
                            value={overrideScores[currentGrade.id] || ''}
                            onChange={(e) =>
                              setOverrideScores((prev) => ({
                                ...prev,
                                [currentGrade.id]: parseFloat(
                                  e.target.value
                                ),
                              }))
                            }
                            placeholder="Enter new score"
                          />
                        </div>

                        <div className="space-y-2">
                          <Label htmlFor="reason">Reason</Label>
                          <Textarea
                            id="reason"
                            value={overrideReasons[currentGrade.id] || ''}
                            onChange={(e) =>
                              setOverrideReasons((prev) => ({
                                ...prev,
                                [currentGrade.id]: e.target.value,
                              }))
                            }
                            placeholder="Why are you overriding this grade?"
                            rows={3}
                          />
                        </div>

                        <Button
                          onClick={() => handleOverrideGrade(currentGrade)}
                          disabled={overriding[currentGrade.id]}
                          className="w-full"
                        >
                          {overriding[currentGrade.id]
                            ? 'Saving...'
                            : 'Confirm Override'}
                        </Button>
                      </div>
                    </DialogContent>
                  </Dialog>
                )}
              </CardContent>
            </Card>
          )}

          {/* Navigation */}
          <div className="flex gap-3">
            <Button
              onClick={() =>
                setCurrentQuestionIdx(Math.max(0, currentQuestionIdx - 1))
              }
              disabled={currentQuestionIdx === 0}
              variant="outline"
              className="flex-1"
            >
              <ChevronLeft className="w-4 h-4 mr-2" />
              Previous
            </Button>

            <div className="flex-1 flex items-center justify-center">
              <p className="text-sm text-slate-400">
                Question {currentQuestionIdx + 1} of {grades.length}
              </p>
            </div>

            <Button
              onClick={() =>
                setCurrentQuestionIdx(
                  Math.min(grades.length - 1, currentQuestionIdx + 1)
                )
              }
              disabled={currentQuestionIdx === grades.length - 1}
              variant="outline"
              className="flex-1"
            >
              Next
              <ChevronRight className="w-4 h-4 ml-2" />
            </Button>
          </div>
        </>
      )}
    </div>
  )
}
