-- 👤 Все пользователи (в том числе администраторы)
CREATE TABLE Users
(
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(10) CHECK (role IN ('admin', 'coach', 'medical', 'athlete')) NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 🎽 Справочник: Виды спорта
CREATE TABLE SportTypes
(
    sport_type_id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

-- 🏅 Справочник: Разряды
CREATE TABLE AthleteRanks
(
    rank_id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

-- 🧪 Справочник: Типы мед. тестов
CREATE TABLE MedicalTestTypes
(
    test_type_id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT
);

-- 🏋️ Справочник: Типы тренировок
CREATE TABLE TrainingSessionTypes
(
    session_type_id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

-- ⏱ Справочник: Форматы тренировок
CREATE TABLE TrainingFormatTypes
(
    format_type_id SERIAL PRIMARY KEY,
    name VARCHAR(20) UNIQUE NOT NULL
);

-- 🗓 Справочник: Периоды
CREATE TABLE PeriodTypes
(
    period_type_id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT
);

-- 🏢 Команды
CREATE TABLE Teams
(
    team_id SERIAL PRIMARY KEY,
    team_name VARCHAR(100) NOT NULL,
    sport_type_id INT REFERENCES SportTypes(sport_type_id)
);

-- 🎿 Справочник: Подстили (например, лыжи: классика/конёк)
CREATE TABLE SportSubtypes
(
    subtype_id SERIAL PRIMARY KEY,
    sport_type_id INT REFERENCES SportTypes(sport_type_id),
    name VARCHAR(50) NOT NULL
);

-- 👨‍🏫 Тренеры
CREATE TABLE Coaches
(
    coach_id INT PRIMARY KEY REFERENCES Users(user_id) ON DELETE CASCADE,
    phone VARCHAR(20),
    first_name VARCHAR(50) NOT NULL,
    middle_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    license_id VARCHAR(50),
    experience_level VARCHAR(20),
    sport_type_id INT REFERENCES SportTypes(sport_type_id)
);

-- 🧍‍♂️ Спортсмены
CREATE TABLE Athletes
(
    athlete_id INT PRIMARY KEY REFERENCES Users(user_id) ON DELETE CASCADE,
    phone VARCHAR(20),
    first_name VARCHAR(50) NOT NULL,
    middle_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    gender CHAR(1) CHECK (gender IN ('M', 'F')) NOT NULL,
    date_of_birth DATE NOT NULL,
    team_id INT REFERENCES Teams(team_id),
    sport_type_id INT REFERENCES SportTypes(sport_type_id),
    rank_id INT REFERENCES AthleteRanks(rank_id)
);

-- 🧍‍♂️ Периоды спортсмена
CREATE TABLE AthletePeriods
(
    period_id SERIAL PRIMARY KEY,
    athlete_id INT REFERENCES Athletes(athlete_id),
    period_type_id INT REFERENCES PeriodTypes(period_type_id),
    start_date DATE NOT NULL,
    end_date DATE,
    assigned_by INT REFERENCES Coaches(coach_id)
);

-- 👥 Связь тренер ↔ команда
CREATE TABLE CoachTeamLink
(
    coach_id INT REFERENCES Coaches(coach_id),
    team_id INT REFERENCES Teams(team_id),
    role_in_team VARCHAR(20) DEFAULT 'coach',
    PRIMARY KEY (coach_id, team_id)
);

-- 🏢 Спортшколы/ЦСП
CREATE TABLE Organizations
(
    organization_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- 👨‍⚕️ Врачи
CREATE TABLE MedicalStaff
(
    medical_staff_id INT PRIMARY KEY REFERENCES Users(user_id) ON DELETE CASCADE,
    first_name VARCHAR(50) NOT NULL,
    middle_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    specialization VARCHAR(100) NOT NULL,
    license_number VARCHAR(50) UNIQUE NOT NULL,
    verified BOOLEAN DEFAULT FALSE,
    organization_id INT REFERENCES Organizations(organization_id)
);

-- Список спортсменов, отправленных на осмотр/проаерку/пересдачу анализов/просто получивших рекомендации
CREATE TABLE MedicalAssignments
(
    assignment_id SERIAL PRIMARY KEY,
    athlete_id INT REFERENCES Athletes(athlete_id),
    assigned_by INT REFERENCES Coaches(coach_id),
    medical_staff_id INT REFERENCES MedicalStaff(medical_staff_id),
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    comment TEXT,
    status VARCHAR(20) DEFAULT 'assigned', -- assigned, reviewed, declined и т.п.
    reviewed_at TIMESTAMP,
    reviewed_by INT REFERENCES MedicalStaff(medical_staff_id)
);

-- 📊 Биометрия
CREATE TABLE BiometricData
(
    biometric_id SERIAL PRIMARY KEY,
    athlete_id INT REFERENCES Athletes(athlete_id),
    date DATE NOT NULL,
    time_of_day VARCHAR(10) CHECK (time_of_day IN ('morning', 'evening')) NOT NULL,
    time TIME NOT NULL,
    pulse INT,
    HRV INT,
    weight FLOAT
);

-- 🧪 Медицинские тесты
CREATE TABLE MedicalTests
(
    test_id SERIAL PRIMARY KEY,
    athlete_id INT REFERENCES Athletes(athlete_id),
    team_id INT REFERENCES Teams(team_id), -- 🔗 команда, за которую числился спортсмен на момент теста
    test_date DATE NOT NULL,
    test_type_id INT REFERENCES MedicalTestTypes(test_type_id),
    test_result TEXT NOT NULL
);

-- 📅 Тренировочный план
CREATE TABLE TrainingPlan
(
    plan_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_by INT REFERENCES Coaches(coach_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 📋 Пункты плана
CREATE TABLE TrainingPlanItem
(
    item_id SERIAL PRIMARY KEY,
    plan_id INT REFERENCES TrainingPlan(plan_id),
    session_date DATE NOT NULL,
    session_order INT DEFAULT 1,
    session_type_id INT REFERENCES TrainingSessionTypes(session_type_id),
    format_type_id INT REFERENCES TrainingFormatTypes(format_type_id),
    value VARCHAR(50),
    notes TEXT
);

-- 🧍‍♂️↔️📋 Назначения плана спортсмену
CREATE TABLE AssignedPlans
(
    assignment_id SERIAL PRIMARY KEY,
    athlete_id INT REFERENCES Athletes(athlete_id),
    plan_id INT REFERENCES TrainingPlan(plan_id),
    start_date DATE NOT NULL,
    assigned_by INT REFERENCES Coaches(coach_id)
);

-- 📄 Рекомендации врача
CREATE TABLE Recommendations
(
    recommendation_id SERIAL PRIMARY KEY,
    athlete_id INT REFERENCES Athletes(athlete_id),
    medical_staff_id INT REFERENCES MedicalStaff(medical_staff_id),
    recommendation_date DATE NOT NULL,
    recommendation_text TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'issued', -- issued, acknowledged, archived и т.д.
    reviewed_at TIMESTAMP,
    reviewed_by INT REFERENCES Coaches(coach_id)
);

-- ⬇️ Лог экспорта
CREATE TABLE DataExportLog
(
    export_id SERIAL PRIMARY KEY,
    athlete_id INT REFERENCES Athletes(athlete_id),
    export_date TIMESTAMP NOT NULL,
    export_format VARCHAR(10) CHECK (export_format IN ('CSV', 'JSON', 'XML')) NOT NULL,
    exported_by INT REFERENCES Users(user_id)
);

-- 📎 Загруженные файлы
CREATE TABLE MedicalFiles
(
    file_id SERIAL PRIMARY KEY,
    athlete_id INT REFERENCES Athletes(athlete_id),
    uploaded_by INT,
    upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    file_type VARCHAR(20),
    file_name VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    description TEXT
);

-- 📥 Импортированные данные
CREATE TABLE ImportedMedicalData
(
    import_id SERIAL PRIMARY KEY,
    athlete_id INT REFERENCES Athletes(athlete_id),
    source_file_id INT REFERENCES MedicalFiles(file_id),
    imported_by INT REFERENCES Users(user_id),
    imported_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    raw_data JSONB
);

-- 📈 Лог отчётов
CREATE TABLE ReportLog
(
    report_id SERIAL PRIMARY KEY,
    generated_by INT REFERENCES Users(user_id),
    team_id INT REFERENCES Teams(team_id),
    athlete_id INT REFERENCES Athletes(athlete_id),
    report_type VARCHAR(50),
    format VARCHAR(10),
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    report_path TEXT NOT NULL
);

-- 📨 Приглашения администраторов на регистрацию врачей и тренеров
CREATE TABLE AdminInvites
(
    invite_id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    role VARCHAR(10) CHECK (role IN ('coach', 'medical')) NOT NULL,
    full_name VARCHAR(150),
    license_number VARCHAR(50),
    organization_id INT REFERENCES Organizations(organization_id) ON DELETE CASCADE,
    invite_code VARCHAR(40) UNIQUE NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    created_by INT REFERENCES Users(user_id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 📬 Заявки спортсменов на вступление в команду
CREATE TABLE TeamJoinRequests
(
    request_id SERIAL PRIMARY KEY,
    athlete_id INT REFERENCES Athletes(athlete_id) ON DELETE CASCADE,
    team_id INT REFERENCES Teams(team_id) ON DELETE CASCADE,
    status VARCHAR(10) CHECK (status IN ('pending', 'approved', 'rejected')) DEFAULT 'pending',
    requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reviewed_by INT REFERENCES Coaches(coach_id) ON DELETE SET NULL,
    reviewed_at TIMESTAMP
);

CREATE TABLE CustomPlanItemOverrides
(
    override_id SERIAL PRIMARY KEY,
    item_id INT REFERENCES TrainingPlanItem(item_id) ON DELETE CASCADE,
    athlete_id INT REFERENCES Athletes(athlete_id) ON DELETE CASCADE,
    custom_value VARCHAR(50),
    custom_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
