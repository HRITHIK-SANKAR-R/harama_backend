import React from 'react'

interface AIReasoningPanelProps {
    reasoning: string
    confidence: number
}

export const AIReasoningPanel = ({ reasoning, confidence }: AIReasoningPanelProps) => {
    const confidenceColor = confidence > 0.8 ? 'text-green-600' : confidence > 0.6 ? 'text-yellow-600' : 'text-red-600'
    const confidenceBg = confidence > 0.8 ? 'bg-green-100' : confidence > 0.6 ? 'bg-yellow-100' : 'bg-red-100'

    return (
        <div className="bg-white rounded-lg shadow-md border border-gray-200 overflow-hidden">
            <div className="bg-blue-600 px-4 py-2 flex justify-between items-center">
                <h3 className="text-white font-semibold">Gemini 3 Reasoning</h3>
                <span className={`${confidenceBg} ${confidenceColor} px-2 py-1 rounded-full text-xs font-bold`}>
                    {(confidence * 100).toFixed(0)}% Confidence
                </span>
            </div>
            <div className="p-4">
                <div className="prose prose-sm max-w-none text-gray-700 italic">
                    "{reasoning}"
                </div>
                <div className="mt-4 flex gap-2">
                    <span className="bg-blue-50 text-blue-700 px-2 py-1 rounded text-[10px] uppercase font-bold tracking-wider">Multimodal</span>
                    <span className={`${confidence > 0.15 ? 'bg-orange-50 text-orange-700' : 'bg-green-50 text-green-700'} px-2 py-1 rounded text-[10px] uppercase font-bold tracking-wider`}>
                        {confidence < 0.7 ? 'Escalated' : 'Verified'}
                    </span>
                </div>
            </div>
        </div>
    )
}
