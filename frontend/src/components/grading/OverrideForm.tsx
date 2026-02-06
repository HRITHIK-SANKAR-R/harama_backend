'use client'

import { useState } from 'react'

interface OverrideFormProps {
    currentScore: number
    onSubmit: (score: number, reason: string) => void
}

export function OverrideForm({ currentScore, onSubmit }: OverrideFormProps) {
    const [score, setScore] = useState(currentScore.toString())
    const [reason, setReason] = useState('')

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        onSubmit(parseFloat(score), reason)
    }

    return (
        <div className="border rounded p-4 mt-4 bg-white shadow-sm">
            <h3 className="font-bold mb-4">Teacher Override</h3>
            <form onSubmit={handleSubmit} className="space-y-4">
                <div className="flex flex-col gap-2">
                    <label className="text-sm font-medium text-gray-700">New Score</label>
                    <input
                        type="number"
                        value={score}
                        onChange={(e) => setScore(e.target.value)}
                        className="border p-2 rounded focus:ring-2 focus:ring-blue-500 outline-none"
                        step="0.5"
                    />
                </div>
                <div className="flex flex-col gap-2">
                    <label className="text-sm font-medium text-gray-700">Reason for change</label>
                    <textarea
                        value={reason}
                        onChange={(e) => setReason(e.target.value)}
                        className="border p-2 rounded focus:ring-2 focus:ring-blue-500 outline-none h-20"
                        placeholder="Explain why you are overriding the AI score..."
                        required
                    />
                </div>
                <button 
                    type="submit"
                    className="w-full bg-blue-600 text-white px-4 py-2 rounded font-medium hover:bg-blue-700 transition-colors"
                >
                    Apply Override
                </button>
            </form>
        </div>
    )
}
