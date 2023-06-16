CREATE TABLE units (
    id INTEGER NOT NULL PRIMARY KEY,
    name TEXT UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE nutrients (
    id INTEGER NOT NULL PRIMARY KEY,
    name TEXT UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE foods (
    id INTEGER NOT NULL PRIMARY KEY,
    name TEXT UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE foods_nutrients (
    food_id INTEGER NOT NULL,
    nutrient_id INTEGER NOT NULL,
    FOREIGN KEY(food_id) REFERENCES foods(id)
    FOREIGN KEY(nutrient_id) REFERENCES nutrients(id)
)
