DROP TABLE IF EXISTS BiometricData;

CREATE TABLE IF NOT EXISTS BiometricData
(
    biometric_id SERIAL PRIMARY KEY,
    athlete_id INT REFERENCES Athletes(athlete_id),
    date timestamp NOT NULL,
    morning_pulse INT,
    evening_pulse INT,
    HRV INT,
    weight FLOAT
);


CREATE TABLE IF NOT EXISTS TrainingPlan
(
    plan_id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    team_id INT NOT NULL REFERENCES Teams(team_id),
    description TEXT,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft', -- 'active', 'completed'
    period_type_id INT REFERENCES PeriodTypes(period_type_id),
    created_by INT REFERENCES Coaches(coach_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS TrainingPlanItem
(
    plan_item_id SERIAL PRIMARY KEY,
    plan_id INT NOT NULL REFERENCES TrainingPlan(plan_id) ON DELETE CASCADE,
    workout_date DATE NOT NULL,
    time TIME, -- Можно NULL если не указано
    session_type_id INT REFERENCES TrainingSessionTypes(session_type_id),
    format_type_id INT REFERENCES TrainingFormatTypes(format_type_id),
    value VARCHAR(50),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS PersonalizedTrainingPlan
(
    id SERIAL PRIMARY KEY,
    athlete_id INT NOT NULL REFERENCES Athletes(athlete_id),
    base_plan_id INT REFERENCES TrainingPlan(plan_id), -- От какого шаблона скопировано
    title VARCHAR(100) NOT NULL,
    description TEXT,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    created_by INT REFERENCES Coaches(coach_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS PersonalizedTrainingPlanItem
(
    id SERIAL PRIMARY KEY,
    plan_id INT NOT NULL REFERENCES PersonalizedTrainingPlan(id) ON DELETE CASCADE,
    workout_date DATE NOT NULL,
    time TIME,
    session_type_id INT REFERENCES TrainingSessionTypes(session_type_id),
    format_type_id INT REFERENCES TrainingFormatTypes(format_type_id),
    value VARCHAR(50),
    description TEXT
);

INSERT INTO TrainingFormatTypes (name) VALUES
                                           ('время'),
                                           ('дистанция'),
                                           ('интервалы'),
                                           ('круги'),
                                           ('пульс'),
                                           ('RPE'),
                                           ('нагрузка'),
                                           ('нет'),
                                           ('другое');

INSERT INTO TrainingSessionTypes (name) VALUES
                                            ('зарядка'),
                                            ('кросс'),
                                            ('велосипед'),
                                            ('плавание'),
                                            ('гребля'),
                                            ('лыжи'),
                                            ('ходьба'),
                                            ('силовая'),
                                            ('отдых'),
                                            ('контрольный старт'),
                                            ('сбор'),
                                            ('растяжка'),
                                            ('заминка'),
                                            ('разминка'),
                                            ('триатлон'),
                                            ('другое');