DROP TABLE IF EXISTS BiometricData;

CREATE TABLE BiometricData
(
    biometric_id SERIAL PRIMARY KEY,
    athlete_id INT REFERENCES Athletes(athlete_id),
    date timestamp NOT NULL,
    morning_pulse INT,
    evening_pulse INT,
    HRV INT,
    weight FLOAT
);