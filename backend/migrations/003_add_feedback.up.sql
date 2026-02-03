CREATE TABLE feedback_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    question_id UUID NOT NULL,
    submission_id UUID NOT NULL,
    ai_score DECIMAL(5,2) NOT NULL,
    teacher_score DECIMAL(5,2) NOT NULL,
    delta DECIMAL(5,2) NOT NULL,
    ai_reasoning TEXT,
    teacher_reason TEXT,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_feedback_question ON feedback_events(question_id);
CREATE INDEX idx_feedback_submission ON feedback_events(submission_id);
