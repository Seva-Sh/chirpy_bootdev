# Chirpy

Chirpy is a local-server based posts uploading application that allows multiple user access functionality.

## Setup

### Prerequisites
 - **Go 1.22+**
 - **PostgreSQL**

### Configuration
Create a `.env` file in the root directory with the following variables:
- `DB_URL`: `postgres://user:password@localhost:5432/database_name`
- `JWT_SECRET`: `your_secure_random_string`
- `POLKA_KEY`: `your_polka_api_key`
- `PLATFORM`: `dev`

## Getting Started
1. Clone the repository: `git clone https://github.com/Seva-Sh/chirpy_bootdev`
2. Install dependencies: `go mod download`
3. Run database migrations: `cd sql/schema && goose postgres "postgres://postgres:2991@localhost:5432/chirpy?sslmode=disable" up`
4. Build the application: `go build -o out`
5. Start the server: `./out`

## Features

Currently supported functionality:
 - User authentication and JWT handling
 - Creating and retrieving "chirps" (posts)
 - Local server configuration via Go

### API Endpoints

- `GET /api/healthz` - Check if the server is alive 
- `GET /admin/metrics` - Check the count of requests
- `GET /api/chirps` - Get chirps (posts). Allows optional query parameters :
    - `author_id`: Returns the chirp written by the selected author
    - `sort`: Sort chirps in ascending (`asc`, default) or descending (`desc`) order
- `GET /api/chirps/{chirpID}` - Get specified chirp via ID
- `POST /admin/reset` - Reset the count of requests
- `POST /api/users` - Create a new user and upload to databse. Requires email and password as parameters
- `POST /api/chirps` - Create a new chirp and upload to database. Requires string as a chirp, cannot be more than 140 characters
- `POST /api/login` - Authenticate a user and receive a JWT. Requires password and email as parameters
- `POST /api/refresh` - Update JWT expiration by 1 hour. Requires token string as parameter
- `POST /api/revoke` - Revoke the refresh token
- `PUT /api/users` - Update user info with new email or password. Requires password and email as parameter
- `DELETE /api/chirps/{chirpID}` - Delete chirp via ID
- `POST /api/polka/webhooks` - Upgrades user to Chirpy Red