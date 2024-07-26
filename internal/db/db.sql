DROP TABLE IF EXISTS limit_law CASCADE;
CREATE TYPE balance_identifiaction_level AS ENUM ('anonymous', 'simplified', 'full');
CREATE TABLE limit_law (
    id SERIAL PRIMARY KEY,
    identifiaction_level balance_identifiaction_level,
    balance_min INT NOT NULL,
    balance_max INT NOT NULL
);

DROP TABLE IF EXISTS balance CASCADE;
CREATE TABLE balance (
    id SERIAL PRIMARY KEY,
    balance INT NOT NULL,
    hold INT NOT NULL,identification_level INTEGER,
    FOREIGN KEY (identification_level) REFERENCES limit_law (id),
    created_at DATE NOT NULL,
    updated_at DATE NOT NULL DEFAULT CURRENT_DATE
);


INSERT INTO limit_law (identifiaction_level, balance_min, balance_max) VALUES ('anonymous', 0, 15000);
INSERT INTO limit_law (identifiaction_level, balance_min, balance_max) VALUES ('simplified', 0, 60000);
INSERT INTO limit_law (identifiaction_level, balance_min, balance_max) VALUES ('full', 0, 600000);