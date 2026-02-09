import { createClient } from './supabase/client'
import { env } from '@/lib/env'

const API_BASE = env.NEXT_PUBLIC_API_URL

async function getAuthHeaders() {
  const supabase = createClient()
  const {
    data: { session },
  } = await supabase.auth.getSession()

  if (!session?.access_token) {
    throw new Error('Not authenticated')
  }

  return {
    Authorization: `Bearer ${session.access_token}`,
    'Content-Type': 'application/json',
  }
}

async function getAuthHeadersFormData() {
  const supabase = createClient()
  const {
    data: { session },
  } = await supabase.auth.getSession()

  if (!session?.access_token) {
    throw new Error('Not authenticated')
  }

  return {
    Authorization: `Bearer ${session.access_token}`,
  }
}

export interface CreateExamRequest {
  title: string
  subject: string
  description?: string
  questions?: Array<{
    question_text: string
    points: number
    answer_type: string
  }>
}

export interface Question {
  id: string
  question_text: string
  points: number
  answer_type: string
  question_number?: string
  question_group?: string
  rubric?: Record<string, any>
}

export interface Exam {
  id: string
  title: string
  subject: string
  description?: string
  questions?: Question[]
  created_at: string
}

export interface Submission {
  id: string
  exam_id: string
  student_id: string
  processing_status: string
  ocr_results?: Array<{
    page_number: number
    raw_text: string
    confidence: number
    image_url: string
  }>
  answers?: Array<{
    question_id: string
    text: string
  }>
  uploaded_at: string
}

export interface Grade {
  id: string
  submission_id: string
  question_id: string
  final_score: number
  max_score: number
  ai_score?: number
  override_score?: number
  confidence: number
  reasoning: string
  criteria_met?: string[]
  mistakes_found?: string[]
  ai_evaluator_id?: string
  status: string
  created_at: string
  updated_at: string
}

// Exam Endpoints
export async function createExam(data: CreateExamRequest) {
  const headers = await getAuthHeaders()
  const response = await fetch(`${API_BASE}/api/v1/exams`, {
    method: 'POST',
    headers,
    body: JSON.stringify(data),
  })

  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Failed to create exam')
  }

  return response.json()
}

export async function listExams(): Promise<Exam[]> {
  const headers = await getAuthHeaders()
  const response = await fetch(`${API_BASE}/api/v1/exams`, {
    headers,
  })

  if (!response.ok) {
    throw new Error('Failed to list exams')
  }

  return response.json()
}

export async function getExam(id: string): Promise<Exam> {
  const headers = await getAuthHeaders()
  const response = await fetch(`${API_BASE}/api/v1/exams/${id}`, {
    headers,
  })

  if (!response.ok) {
    throw new Error('Failed to get exam')
  }

  return response.json()
}

export async function addQuestion(
  examId: string,
  data: {
    question_text: string
    points: number
    answer_type: string
    question_number?: string
    question_group?: string
  }
) {
  const headers = await getAuthHeaders()
  const response = await fetch(
    `${API_BASE}/api/v1/exams/${examId}/questions`,
    {
      method: 'POST',
      headers,
      body: JSON.stringify(data),
    }
  )

  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Failed to add question')
  }

  return response.json()
}

export async function setRubric(
  questionId: string,
  data: {
    full_credit_criteria?: Array<{ description: string; points: number }>
    partial_credit_rules?: any[]
    common_mistakes?: any[]
    key_concepts?: string[]
  }
) {
  const headers = await getAuthHeaders()
  const response = await fetch(
    `${API_BASE}/api/v1/questions/${questionId}/rubric`,
    {
      method: 'PUT',
      headers,
      body: JSON.stringify(data),
    }
  )

  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Failed to set rubric')
  }

  return response.json()
}

// Submission Endpoints
export async function uploadSubmission(
  examId: string,
  file: File,
  studentId: string
): Promise<Submission> {
  const headers = await getAuthHeadersFormData()

  const formData = new FormData()
  formData.append('file', file)
  formData.append('student_id', studentId)

  const response = await fetch(
    `${API_BASE}/api/v1/exams/${examId}/submissions`,
    {
      method: 'POST',
      headers,
      body: formData,
    }
  )

  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Failed to upload submission')
  }

  return response.json()
}

export async function getSubmission(id: string): Promise<Submission> {
  const headers = await getAuthHeaders()
  const response = await fetch(`${API_BASE}/api/v1/submissions/${id}`, {
    headers,
  })

  if (!response.ok) {
    throw new Error('Failed to get submission')
  }

  return response.json()
}

export async function triggerGrading(submissionId: string) {
  const headers = await getAuthHeaders()
  const response = await fetch(
    `${API_BASE}/api/v1/submissions/${submissionId}/trigger-grading`,
    {
      method: 'POST',
      headers,
    }
  )

  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Failed to trigger grading')
  }

  return response.json()
}

// Grade Endpoints
export async function getGrades(submissionId: string): Promise<Grade[]> {
  const headers = await getAuthHeaders()
  const response = await fetch(
    `${API_BASE}/api/v1/submissions/${submissionId}/grades`,
    {
      headers,
    }
  )

  if (!response.ok) {
    throw new Error('Failed to get grades')
  }

  return response.json()
}

export async function overrideGrade(
  submissionId: string,
  questionId: string,
  newScore: number,
  reason: string
) {
  const headers = await getAuthHeaders()
  const response = await fetch(
    `${API_BASE}/api/v1/submissions/${submissionId}/questions/${questionId}/override`,
    {
      method: 'POST',
      headers,
      body: JSON.stringify({ new_score: newScore, reason }),
    }
  )

  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Failed to override grade')
  }

  return response.json()
}

export async function getFeedback(
  submissionId: string,
  questionId: string
) {
  const headers = await getAuthHeaders()
  const response = await fetch(
    `${API_BASE}/api/v1/submissions/${submissionId}/questions/${questionId}/feedback`,
    {
      headers,
    }
  )

  if (!response.ok) {
    throw new Error('Failed to get feedback')
  }

  return response.json()
}

// Analytics Endpoints
export async function getGradingTrends(examId?: string) {
  const headers = await getAuthHeaders()
  const url = examId
    ? `${API_BASE}/api/v1/analytics/grading-trends?exam_id=${examId}`
    : `${API_BASE}/api/v1/analytics/grading-trends`
  const response = await fetch(url, {
    headers,
  })

  if (!response.ok) {
    throw new Error('Failed to get grading trends')
  }

  return response.json()
}

export async function getQuestionAnalysis(questionId: string) {
  const headers = await getAuthHeaders()
  const response = await fetch(
    `${API_BASE}/api/v1/questions/${questionId}/analysis`,
    {
      headers,
    }
  )

  if (!response.ok) {
    throw new Error('Failed to get question analysis')
  }

  return response.json()
}

export async function exportGrades(examId: string, format: string = 'csv') {
  const headers = await getAuthHeaders()
  const response = await fetch(`${API_BASE}/api/v1/exams/${examId}/export`, {
    method: 'POST',
    headers,
    body: JSON.stringify({ format }),
  })

  if (!response.ok) {
    throw new Error('Failed to export grades')
  }

  return response.blob()
}

export async function adaptRubric(questionId: string) {
  const headers = await getAuthHeaders()
  const response = await fetch(
    `${API_BASE}/api/v1/questions/${questionId}/adapt-rubric`,
    {
      method: 'POST',
      headers,
    }
  )

  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Failed to adapt rubric')
  }

  return response.json()
}

export async function getAuditLogs(entityId: string, entityType: string = 'submission') {
  const headers = await getAuthHeaders()
  const response = await fetch(
    `${API_BASE}/api/v1/audit/${entityId}?type=${entityType}`,
    { headers }
  )

  if (!response.ok) {
    throw new Error('Failed to get audit logs')
  }

  return response.json()
}

export async function uploadBatchSubmissions(
  examId: string,
  files: File[],
  studentMapping: Record<string, string>
) {
  const headers = await getAuthHeadersFormData()

  const formData = new FormData()
  files.forEach((file) => formData.append('files[]', file))
  formData.append('student_mapping', JSON.stringify(studentMapping))

  const response = await fetch(
    `${API_BASE}/api/v1/exams/${examId}/submissions/batch`,
    {
      method: 'POST',
      headers,
      body: formData,
    }
  )

  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Failed to upload batch')
  }

  return response.json()
}

export async function getBatchStatus(batchId: string) {
  const headers = await getAuthHeaders()
  const response = await fetch(`${API_BASE}/api/v1/batches/${batchId}/status`, {
    headers,
  })

  if (!response.ok) {
    throw new Error('Failed to get batch status')
  }

  return response.json()
}
