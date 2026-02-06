import { GradingView } from '@/components/grading/GradingView'

export default function GradingPage({ params }: { params: { id: string } }) {
    return (
        <div className="p-8">
            <h1 className="text-3xl font-bold mb-6">Review Submission</h1>
            <GradingView submissionId={params.id} />
        </div>
    )
}
