import React from 'react'

export const AIReasoningPanel = ({ reasoning, confidence }: { reasoning: string, confidence: number }) => (
    <div>
        <div>Confidence: {confidence}</div>
        <div>{reasoning}</div>
    </div>
)
