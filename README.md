# FlightOwl API

By Ian Shakespeare

## Requirements

- Go v1.19 or later
- SQLite3
- `.env` with `apiKey` and `apiSecret` fields

## Schema

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    sex TEXT,
    date_joined TEXT NOT NULL -- YYYY-MM-DD HH:MM:SS.SSS
);
```
