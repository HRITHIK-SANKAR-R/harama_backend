'use client'

import { useState } from 'react'

export default function CreateExamPage() {
    const [title, setTitle] = useState('')
    const [subject, setSubject] = useState('')

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        const res = await fetch('/api/v1/exams', {
            method: 'POST',
            body: JSON.stringify({ title, subject }),
            headers: { 'Content-Type': 'application/json' }
        })
        if (res.ok) {
            alert('Exam created!')
        }
    }

    return (
        <div className="p-8">
            <h1 className="text-2xl font-bold mb-4">Create New Exam</h1>
            <form onSubmit={handleSubmit} className="space-y-4 max-w-md">
                <div>
                    <label className="block text-sm font-medium">Exam Title</label>
                    <input 
                        type="text" 
                        value={title} 
                        onChange={(e) => setTitle(e.target.value)}
                        className="mt-1 block w-full border rounded-md p-2"
                        placeholder="e.g. Midterm Physics"
                    />
                </div>
                <div>
                    <label className="block text-sm font-medium">Subject</label>
                    <input 
                        type="text" 
                        value={subject} 
                        onChange={(e) => setSubject(e.target.value)}
                        className="mt-1 block w-full border rounded-md p-2"
                        placeholder="e.g. Physics"
                    />
                </div>
                <button type="submit" className="bg-blue-600 text-white px-4 py-2 rounded-md">
                    Create Exam
                </button>
            </form>
        </div>
    )
}
