# FlightOwl API

By Ian Shakespeare

## Description

The API server for [FlightOwl](https://www.flightowl.app). Database management is handled by the `database` module. `auth` manages token generation and validation. All routes are defined in the `routes` module.

## Routes

| METHOD | ROUTE |
|---|---|
| POST | /register |
| POST | /login |
| GET | /user |
| GET | /flights/saved |
| POST | /flights |
| POST | /flights/check |
| POST | /flights/saved |

## Schema

```sql
CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    sex TEXT NOT NULL,
    date_joined TEXT NOT NULL,
    admin INTEGER DEFAULT 0 NOT NULL
);
CREATE TABLE IF NOT EXISTS flight_offers
(
    offer_id SERIAL PRIMARY KEY,
    date_saved TEXT NOT NULL,
    offer TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE
);
```

## Requirements

- Go v1.19 or later
- `.env` with `DATABASE_URL`, `API_KEY`, and `API_SECRET` fields
