'use client'

import React from "react"

import { useState, useCallback } from 'react'
import { useParams } from 'next/navigation'
import { uploadSubmission, uploadBatchSubmissions } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useToast } from '@/hooks/use-toast'
import { Upload, File, Trash2 } from 'lucide-react'

export default function SubmissionUploadPage() {
  const params = useParams()
  const examId = params.id as string
  const { toast } = useToast()

  const [files, setFiles] = useState<Array<{ file: File; studentId: string }>>(
    []
  )
  const [uploadMode, setUploadMode] = useState<'single' | 'batch'>('batch')
  const [studentId, setStudentId] = useState('')
  const [uploading, setUploading] = useState(false)

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
  }, [])

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()

    const droppedFiles = Array.from(e.dataTransfer.files)
    handleFiles(droppedFiles)
  }, [])

  const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      const selectedFiles = Array.from(e.target.files)
      handleFiles(selectedFiles)
    }
  }

  const handleFiles = (newFiles: File[]) => {
    const validFiles = newFiles.filter((file) => {
      const isValid =
        file.type === 'application/pdf' ||
        file.type.startsWith('image/')

      if (!isValid) {
        toast({
          title: 'Invalid file type',
          description: `${file.name} is not a PDF or image file`,
          variant: 'destructive',
        })
      }

      return isValid
    })

    if (uploadMode === 'single') {
      if (validFiles.length > 0) {
        const file = validFiles[0]
        setFiles([{ file, studentId }])
      }
    } else {
      setFiles((prev) => [
        ...prev,
        ...validFiles.map((file) => ({
          file,
          studentId: extractStudentId(file.name),
        })),
      ])
    }
  }

  const extractStudentId = (filename: string) => {
    const withoutExtension = filename.replace(/\.[^/.]+$/, '')
    return withoutExtension
  }

  const updateStudentId = (fileIndex: number, newStudentId: string) => {
    setFiles((prev) => {
      const updated = [...prev]
      updated[fileIndex] = {
        ...updated[fileIndex],
        studentId: newStudentId,
      }
      return updated
    })
  }

  const removeFile = (fileIndex: number) => {
    setFiles((prev) => prev.filter((_, i) => i !== fileIndex))
  }

  const handleUpload = async () => {
    if (files.length === 0) {
      toast({
        title: 'Error',
        description: 'Please select at least one file',
        variant: 'destructive',
      })
      return
    }

    // Check if all files have student IDs
    if (files.some((f) => !f.studentId.trim())) {
      toast({
        title: 'Error',
        description: 'All files must have a student ID',
        variant: 'destructive',
      })
      return
    }

    setUploading(true)

    try {
      if (uploadMode === 'single') {
        const { file, studentId } = files[0]
        await uploadSubmission(examId, file, studentId)
        toast({
          title: 'Success',
          description: 'Submission uploaded successfully',
        })
      } else {
        const studentMapping: Record<string, string> = {}
        files.forEach(({ file, studentId }) => {
          studentMapping[file.name] = studentId
        })

        await uploadBatchSubmissions(
          examId,
          files.map((f) => f.file),
          studentMapping
        )

        toast({
          title: 'Success',
          description: `${files.length} submissions uploaded successfully`,
        })
      }

      setFiles([])
    } catch (error: any) {
      toast({
        title: 'Upload failed',
        description: error.message,
        variant: 'destructive',
      })
    } finally {
      setUploading(false)
    }
  }

  return (
    <div className="space-y-6 p-6">
      <div>
        <h1 className="text-3xl font-bold text-white">Upload Submissions</h1>
        <p className="text-slate-400 mt-1">Upload student answer sheets for grading</p>
      </div>

      {/* Mode Selection */}
      <Card>
        <CardHeader>
          <CardTitle>Upload Mode</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex gap-4">
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="radio"
                value="single"
                checked={uploadMode === 'single'}
                onChange={(e) => {
                  setUploadMode(e.target.value as 'single' | 'batch')
                  setFiles([])
                }}
              />
              <span>Single Upload</span>
            </label>
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="radio"
                value="batch"
                checked={uploadMode === 'batch'}
                onChange={(e) => {
                  setUploadMode(e.target.value as 'single' | 'batch')
                  setFiles([])
                }}
              />
              <span>Batch Upload</span>
            </label>
          </div>
        </CardContent>
      </Card>

      {/* Upload Zone */}
      <Card>
        <CardHeader>
          <CardTitle>Select Files</CardTitle>
          <CardDescription>
            Upload PDF or image files (JPG, PNG)
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div
            onDragOver={handleDragOver}
            onDrop={handleDrop}
            className="border-2 border-dashed border-slate-600 rounded-lg p-8 text-center hover:border-blue-500 transition-colors bg-slate-900/30"
          >
            <Upload className="w-12 h-12 text-slate-500 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-white">
              Drag files here
            </h3>
            <p className="text-slate-400 mt-1">or</p>
            <Label className="mt-4 inline-block">
              <Button variant="outline" asChild>
                <span>Browse Files</span>
              </Button>
              <input
                type="file"
                multiple={uploadMode === 'batch'}
                accept="application/pdf,image/*"
                onChange={handleFileInput}
                className="hidden"
              />
            </Label>
          </div>

          {/* File List */}
          {files.length > 0 && (
            <div className="space-y-3">
              <h3 className="font-medium text-white">
                Selected Files ({files.length})
              </h3>
              {files.map((item, i) => (
                <div
                  key={i}
                  className="flex items-center gap-3 p-3 border border-slate-700 rounded-lg bg-slate-900/30"
                >
                  <File className="w-5 h-5 text-slate-500 flex-shrink-0" />
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-white truncate">
                      {item.file.name}
                    </p>
                    <p className="text-xs text-slate-500">
                      {(item.file.size / 1024 / 1024).toFixed(2)} MB
                    </p>
                  </div>

                  {uploadMode === 'batch' ? (
                    <div className="flex items-center gap-2 flex-shrink-0">
                      <Input
                        value={item.studentId}
                        onChange={(e) =>
                          updateStudentId(i, e.target.value)
                        }
                        placeholder="Student ID"
                        className="w-32 text-sm"
                      />
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => removeFile(i)}
                        className="text-red-600"
                      >
                        <Trash2 className="w-4 h-4" />
                      </Button>
                    </div>
                  ) : (
                    <div className="flex items-center gap-2 flex-shrink-0">
                      <Input
                        value={studentId}
                        onChange={(e) => setStudentId(e.target.value)}
                        placeholder="Student ID"
                        className="w-32 text-sm"
                      />
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}

          {/* Upload Button */}
          {files.length > 0 && (
            <Button
              onClick={handleUpload}
              disabled={uploading}
              className="w-full"
              size="lg"
            >
              {uploading
                ? 'Uploading...'
                : `Upload ${files.length} File${files.length !== 1 ? 's' : ''}`}
            </Button>
          )}
        </CardContent>
      </Card>

      {/* Info */}
      <Card className="bg-blue-950/30 border-blue-500/20">
        <CardContent className="pt-6">
          <h4 className="font-medium text-blue-300 mb-2">Supported Formats</h4>
          <ul className="text-sm text-blue-400 space-y-1">
            <li>• PDF files (.pdf)</li>
            <li>• Image files (.jpg, .png, .jpeg)</li>
            <li>• Maximum file size: 50 MB per file</li>
          </ul>
        </CardContent>
      </Card>
    </div>
  )
}
