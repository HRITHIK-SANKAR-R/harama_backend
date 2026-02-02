CREATE TABLE escalations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    submission_id UUID NOT NULL REFERENCES submissions(id),
    question_id UUID NOT NULL REFERENCES questions(id),
    all_evaluations JSONB NOT NULL,
    variance DECIMAL(5,2) NOT NULL,
    escalated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    assigned_to UUID,
    status VARCHAR(50) NOT NULL DEFAULT 'pending'
);
