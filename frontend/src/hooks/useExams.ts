'use client';

import { useState, useEffect } from 'react';
import { Exam } from '@/types';
import { api } from '@/lib/api';

export function useExams() {
    const [exams, setExams] = useState<Exam[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchExams = async () => {
        try {
            setLoading(true);
            const data = await api.listExams();
            setExams(data);
            setError(null);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred while fetching exams');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchExams();
    }, []);

    return { 
        exams, 
        loading, 
        error, 
        refresh: fetchExams 
    };
}
