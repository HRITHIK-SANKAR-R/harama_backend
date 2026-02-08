ALTER TABLE grades ADD CONSTRAINT grades_submission_question_key UNIQUE (submission_id, question_id);
