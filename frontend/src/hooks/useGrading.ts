export const useGrading = (submissionId: string) => {
    return {
        grading: {
            score: 85,
            confidence: 0.9,
            aiReasoning: "Good job",
            answer: {
                text: "Sample Answer",
                image: "/sample.png"
            }
        },
        loading: false,
        applyOverride: (score: number) => console.log(score)
    }
}
