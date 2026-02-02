'use client'

import { useState, useEffect } from 'react'
import { AnswerDisplay } from './AnswerDisplay'
import { AIReasoningPanel } from './AIReasoningPanel'
import { OverrideForm } from './OverrideForm'
import { useGrading } from '@/hooks/useGrading'
import { LoadingSpinner } from '@/components/ui/LoadingSpinner'

export function GradingView({ submissionId }: { submissionId: string }) {
    const [overrideScore, setOverrideScore] = useState<number | null>(null)
    const { grading, loading, applyOverride } = useGrading(submissionId)
    
    if (loading) return <LoadingSpinner />
    
    return (
        <div className="grid grid-cols-2 gap-4">
            <AnswerDisplay answer={grading.answer} />
            <AIReasoningPanel 
                reasoning={grading.aiReasoning}
                confidence={grading.confidence}
            />
            <OverrideForm 
                currentScore={grading.score}
                onSubmit={(score) => applyOverride(score)}
            />
        </div>
    )
}
