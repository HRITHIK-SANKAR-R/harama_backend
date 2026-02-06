import Link from 'next/link'

export default function DashboardPage() {
    const submissions = [
        { id: '8f14f066-cd9a-4c67-bde5-4a09876eae12', student: 'Student 123', status: 'Needs Review', variance: '18%' },
        { id: '2', student: 'Student 456', status: 'Completed', variance: '2%' },
    ]

    return (
        <div className="p-8">
            <h1 className="text-3xl font-bold mb-6">Teacher Dashboard</h1>
            
            <div className="grid gap-4">
                <h2 className="text-xl font-semibold">Pending Review</h2>
                {submissions.map(sub => (
                    <div key={sub.id} className="border p-4 rounded flex justify-between items-center bg-white shadow-sm">
                        <div>
                            <p className="font-medium">{sub.student}</p>
                            <p className="text-sm text-gray-500">Status: <span className={sub.status === 'Needs Review' ? 'text-red-500' : 'text-green-500'}>{sub.status}</span></p>
                        </div>
                        <div className="text-right">
                            <p className="text-sm font-mono">Variance: {sub.variance}</p>
                            <Link href={`/grading/${sub.id}`} className="text-blue-600 hover:underline">
                                Grade Submission {'->'}
                            </Link>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    )
}
