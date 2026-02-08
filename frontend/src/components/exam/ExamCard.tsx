'use client';

import { Exam } from '@/types'
import Link from 'next/link'
import { Calendar, BookOpen, ChevronRight } from 'lucide-react'

interface ExamCardProps {
    exam: Exam
}

export const ExamCard = ({ exam }: { exam: Exam }) => {
    return (
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 hover:shadow-md transition-shadow">
            <div className="flex justify-between items-start mb-4">
                <div className="p-2 bg-blue-50 rounded-lg">
                    <BookOpen className="w-6 h-6 text-blue-600" />
                </div>
                <span className="text-xs font-medium px-2 py-1 bg-gray-100 text-gray-600 rounded-full">
                    {exam.subject}
                </span>
            </div>
            
            <h3 className="text-xl font-bold text-gray-900 mb-2">{exam.title}</h3>
            
            <div className="flex items-center text-sm text-gray-500 mb-6">
                <Calendar className="w-4 h-4 mr-2" />
                {new Date(exam.createdAt).toLocaleDateString()}
            </div>
            
            <div className="flex items-center justify-between pt-4 border-t border-gray-50">
                <span className="text-sm font-medium text-gray-700">
                    {exam.questions?.length || 0} Questions
                </span>
                <Link 
                    href={`/exams/${exam.id}`}
                    className="flex items-center text-sm font-bold text-blue-600 hover:text-blue-700"
                >
                    View Details
                    <ChevronRight className="w-4 h-4 ml-1" />
                </Link>
            </div>
        </div>
    )
}
