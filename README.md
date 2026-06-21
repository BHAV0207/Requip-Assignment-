# Requip User Management System

A comprehensive User Management system built with a Go backend (Gin + GORM + MySQL) and a React frontend (Vite + TypeScript + Tailwind CSS).

## Features
- **Full CRUD**: Create, Read, Update, and Soft Delete users.
- **Pagination**: Paginated user listing on both backend and frontend.
- **Validation**: Strict server-side and client-side validation for fields like Aadhaar, PAN, and Mobile.
- **Modern UI**: Sleek, responsive design using Tailwind CSS v4 and Lucide icons.
- **Unit Testing**: Comprehensive unit tests for backend Service and Handler layers using mocks.

## Tech Stack
- **Backend**: Go 1.25, Gin, GORM, MySQL 8.4.
- **Frontend**: React 19, TypeScript, Vite, Tailwind CSS v4, Axios.
- **Database**: MySQL (running in Docker).

## Setup Instructions

### Backend Setup
1. Navigate to the `Backend` directory.
2. Configure environmental variables in `.env` (you can copy `.env.example`).
3. Start the database using Docker:
   ```bash
   docker compose up -d
   ```
4. Run the backend server:
   ```bash
   go run cmd/server/main.go
   ```

### Frontend Setup
1. Navigate to the `Frontend` directory.
2. Install dependencies:
   ```bash
   npm install
   ```
3. Start the development server:
   ```bash
   npm run dev
   ```

## API Documentation

All API endpoints are prefixed with `/api/v1`. By default, the server runs on `http://localhost:8080`.

### User Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/users` | Create a new user record. |
| `GET` | `/users` | List users with pagination using `page` and `limit` query parameters. |
| `GET` | `/users/:id` | Retrieve detailed information for a specific user. |
| `PUT` | `/users/:id` | Update user details (Note: PAN and Aadhaar are immutable). |
| `DELETE` | `/users/:id` | Soft delete a user record. |

#### Create User Request Body
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "primaryMobile": "9876543210",
  "secondaryMobile": "9876543211",
  "aadhaar": "123456789012",
  "pan": "ABCDE1234F",
  "dateOfBirth": "1990-01-01",
  "placeOfBirth": "Mumbai",
  "currentAddress": "123 Street, Mumbai",
  "permanentAddress": "123 Street, Mumbai"
}
```

#### Running Tests
To run the backend unit tests and verify the logic:
```bash
cd Backend
go test ./internal/user/...
```


## Pain Points & Learnings

### Pain Points
- **MySQL Compatibility**: Migrating from Postgres required careful adjustment of UUID generation patterns, as MySQL handles UUIDs differently (using `char(36)` or similar instead of a native `uuid` type with defaults).
- **Port Conflicts**: Encountered issues with the default MySQL port (3306) being used by host services, requiring a shift to 33066 for the containerized database.
- **Tailwind CSS v4 Migration**: Transitioning to Tailwind v4 required a shift from the traditional `tailwind.config.js` to a CSS-first configuration, which has different build-time requirements and linting rules.
- **React 19 Linting**: Stricter `react-hooks/set-state-in-effect` rules in React 19 required refactoring components to use key-based remounting and `useCallback` to avoid cascading renders.

### Learnings
- **Interface-Driven Design**: Extracting interfaces for Repository and Service significantly improved the testability of the backend, allowing for clean mocking without complex setup.
- **Key-Based Remounting**: Using the `key` prop in React to reset complex forms proved to be a much more readable and robust pattern than syncing state via `useEffect`.
- **Phased Development**: Building the project in small, incremental chunks with frequent verification ensured that bugs were caught early and the system remained stable throughput the process.
