'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { createExam, addQuestion, setRubric } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { useToast } from '@/hooks/use-toast'
import { ChevronRight, Plus, Trash2 } from 'lucide-react'

type Step = 'basic' | 'questions' | 'review'

export default function CreateExamPage() {
  const router = useRouter()
  const { toast } = useToast()

  // Basic Info
  const [title, setTitle] = useState('')
  const [subject, setSubject] = useState('')
  const [description, setDescription] = useState('')

  // Questions
  const [questions, setQuestions] = useState<
    Array<{
      id: string
      question_text: string
      points: number
      answer_type: string
      number?: string
      group?: string
    }>
  >([])

  const [currentQuestion, setCurrentQuestion] = useState({
    question_text: '',
    points: 1,
    answer_type: 'short_answer',
    number: '',
    group: '',
  })

  // UI State
  const [step, setStep] = useState<Step>('basic')
  const [creating, setCreating] = useState(false)

  const addQuestionClick = () => {
    if (!currentQuestion.question_text.trim()) {
      toast({
        title: 'Error',
        description: 'Question text is required',
        variant: 'destructive',
      })
      return
    }

    setQuestions([
      ...questions,
      {
        id: `q-${Date.now()}`,
        ...currentQuestion,
      },
    ])

    setCurrentQuestion({
      question_text: '',
      points: 1,
      answer_type: 'short_answer',
      number: '',
      group: '',
    })
  }

  const removeQuestion = (id: string) => {
    setQuestions(questions.filter((q) => q.id !== id))
  }

  const handleCreateExam = async () => {
    if (!title.trim() || !subject.trim()) {
      toast({
        title: 'Error',
        description: 'Title and subject are required',
        variant: 'destructive',
      })
      return
    }

    setCreating(true)

    try {
      const examData = await createExam({
        title,
        subject,
        description,
        questions: questions.map((q) => ({
          question_text: q.question_text,
          points: q.points,
          answer_type: q.answer_type,
        })),
      })

      toast({
        title: 'Success',
        description: 'Exam created successfully',
      })

      router.push(`/exams/${examData.id}`)
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.message,
        variant: 'destructive',
      })
    } finally {
      setCreating(false)
    }
  }

  return (
    <div className="space-y-8 p-6 bg-gradient-to-br from-slate-950 to-slate-900 min-h-screen">
      <div className="max-w-4xl mx-auto space-y-8">
        {/* Header */}
        <div className="space-y-2">
          <h1 className="text-4xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-purple-600">Create New Exam</h1>
          <p className="text-slate-400">Set up your exam with questions and answer keys</p>
        </div>

        {/* Progress Steps */}
        <div className="flex gap-3">
          {(['basic', 'questions', 'review'] as const).map((s, i) => {
            const isComplete = (step === 'questions' && s === 'basic') || (step === 'review' && (s === 'basic' || s === 'questions'))
            const isActive = step === s
            return (
              <div key={s} className="flex items-center gap-3 flex-1">
                <button
                  onClick={() => setStep(s)}
                  className={`flex-shrink-0 w-10 h-10 rounded-full font-semibold transition-all ${
                    isActive
                      ? 'bg-gradient-to-r from-blue-600 to-blue-500 text-white shadow-lg scale-110'
                      : isComplete
                      ? 'bg-green-500 text-white'
                      : 'bg-slate-700 text-slate-300'
                  }`}
                >
                  {isComplete ? '✓' : i + 1}
                </button>
                <div className="hidden md:block">
                  <p className={`text-xs font-semibold uppercase tracking-wider ${isActive ? 'text-blue-400' : 'text-slate-400'}`}>{s}</p>
                </div>
                {i < 2 && (
                  <div className={`flex-1 h-1 rounded-full mx-1 ${isComplete || isActive ? 'bg-blue-600' : 'bg-slate-700'}`} />
                )}
              </div>
            )
          })}
        </div>

        {/* Step: Basic Info */}
        {step === 'basic' && (
          <Card className="border-0 shadow-lg dark:bg-slate-800/50 dark:border dark:border-slate-700">
            <CardHeader className="border-b border-slate-200 dark:border-slate-700">
              <CardTitle className="text-2xl">Basic Information</CardTitle>
              <CardDescription>Define your exam details</CardDescription>
            </CardHeader>
            <CardContent className="space-y-5 pt-6">
              <div className="space-y-2">
                <Label htmlFor="title">Exam Title *</Label>
                <Input
                  id="title"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder="e.g., Physics Midterm"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="subject">Subject *</Label>
                <Input
                  id="subject"
                  value={subject}
                  onChange={(e) => setSubject(e.target.value)}
                  placeholder="e.g., Physics"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="description">Description</Label>
                <Textarea
                  id="description"
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  placeholder="Add any additional notes about this exam"
                  rows={4}
                />
              </div>

              <Button onClick={() => setStep('questions')} className="w-full bg-gradient-to-r from-blue-600 to-blue-500 hover:from-blue-700 hover:to-blue-600 text-white font-medium h-11 rounded-lg mt-4">
                Next: Add Questions
                <ChevronRight className="w-4 h-4 ml-2" />
              </Button>
            </CardContent>
          </Card>
        )}

        {/* Step: Questions */}
        {step === 'questions' && (
          <Card className="border-0 shadow-lg bg-slate-800/50 border border-slate-700">
            <CardHeader className="border-b border-slate-700">
              <CardTitle className="text-2xl">Add Questions</CardTitle>
              <CardDescription>Create and manage exam questions</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6 pt-6">
              {/* Question List */}
              {questions.length > 0 && (
                <div className="space-y-3 pb-4 border-b">
                  {questions.map((q) => (
                    <div key={q.id} className="flex items-start justify-between p-3 bg-slate-800/50 rounded-lg">
                      <div className="flex-1">
                        <p className="font-medium text-white">{q.question_text}</p>
                        <div className="flex gap-3 mt-2 text-sm text-slate-400">
                          <span>{q.answer_type}</span>
                          <span>•</span>
                          <span>{q.points} marks</span>
                        </div>
                      </div>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => removeQuestion(q.id)}
                        className="text-red-600 hover:text-red-700"
                      >
                        <Trash2 className="w-4 h-4" />
                      </Button>
                    </div>
                  ))}
                </div>
              )}

              {/* Add Question Form */}
              <div className="space-y-4">
                <h3 className="font-medium">Add Question</h3>

                <div className="space-y-2">
                  <Label htmlFor="question">Question Text *</Label>
                  <Textarea
                    id="question"
                    value={currentQuestion.question_text}
                    onChange={(e) =>
                      setCurrentQuestion({
                        ...currentQuestion,
                        question_text: e.target.value,
                      })
                    }
                    placeholder="Enter the question text"
                    rows={3}
                  />
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="points">Points *</Label>
                    <Input
                      id="points"
                      type="number"
                      min="1"
                      value={currentQuestion.points}
                      onChange={(e) =>
                        setCurrentQuestion({
                          ...currentQuestion,
                          points: parseInt(e.target.value) || 1,
                        })
                      }
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="answerType">Answer Type *</Label>
                    <Select
                      value={currentQuestion.answer_type}
                      onValueChange={(value) =>
                        setCurrentQuestion({
                          ...currentQuestion,
                          answer_type: value,
                        })
                      }
                    >
                      <SelectTrigger id="answerType">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="short_answer">Short Answer</SelectItem>
                        <SelectItem value="essay">Essay</SelectItem>
                        <SelectItem value="mcq">Multiple Choice</SelectItem>
                        <SelectItem value="diagram">Diagram</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>

                <Button onClick={addQuestionClick} variant="outline" className="w-full bg-transparent">
                  <Plus className="w-4 h-4 mr-2" />
                  Add Question
                </Button>
              </div>

              <div className="flex gap-3">
                <Button
                  onClick={() => setStep('basic')}
                  variant="outline"
                  className="flex-1"
                >
                  Back
                </Button>
                <Button
                  onClick={() => setStep('review')}
                  disabled={questions.length === 0}
                  className="flex-1"
                >
                  Review & Create
                </Button>
              </div>
            </CardContent>
          </Card>
        )}

      {/* Step: Review */}
      {step === 'review' && (
        <Card className="border-0 shadow-lg bg-slate-800/50 border border-slate-700">
          <CardHeader className="border-b border-slate-700">
            <CardTitle className="text-2xl">Review & Create</CardTitle>
            <CardDescription>Verify all details before finalizing</CardDescription>
          </CardHeader>
          <CardContent className="space-y-6 pt-6">
              <div>
                <h3 className="font-medium text-white mb-2">Exam Details</h3>
                <dl className="space-y-2 text-sm">
                  <div className="flex justify-between">
                    <dt className="text-slate-400">Title:</dt>
                    <dd className="font-medium">{title}</dd>
                  </div>
                  <div className="flex justify-between">
                    <dt className="text-slate-400">Subject:</dt>
                    <dd className="font-medium">{subject}</dd>
                  </div>
                  <div className="flex justify-between">
                    <dt className="text-slate-400">Questions:</dt>
                    <dd className="font-medium">{questions.length}</dd>
                  </div>
                  <div className="flex justify-between">
                    <dt className="text-slate-400">Total Points:</dt>
                    <dd className="font-medium">
                      {questions.reduce((sum, q) => sum + q.points, 0)}
                    </dd>
                  </div>
                </dl>
              </div>

              <div>
                <h3 className="font-medium text-white mb-3">Questions</h3>
                <div className="space-y-2">
                  {questions.map((q, i) => (
                    <div key={q.id} className="text-sm p-2 bg-slate-800/50 rounded">
                      <p className="font-medium">
                        Q{i + 1}. {q.question_text}
                      </p>
                      <p className="text-slate-400 text-xs mt-1">
                        {q.answer_type} • {q.points} points
                      </p>
                    </div>
                  ))}
                </div>
              </div>

              <div className="flex gap-3">
                <Button
                  onClick={() => setStep('questions')}
                  variant="outline"
                  className="flex-1"
                >
                  Back
                </Button>
                <Button
                  onClick={handleCreateExam}
                  disabled={creating}
                  className="flex-1"
                >
                  {creating ? 'Creating...' : 'Create Exam'}
                </Button>
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  )
}
