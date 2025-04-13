# Status Page Application

A simple status page application for monitoring and displaying service status and incidents, similar to StatusPage or Cachet.

## Features

- Authentication and user management
- Organization-based multi-tenancy
- Service status management
- Incident and maintenance tracking
- Real-time updates via WebSockets
- Public status pages for each organization

## Tech Stack

- **Backend**:
  - Go with Gin web framework
  - GORM for database ORM
  - PostgreSQL for persistent storage
  - WebSockets for real-time updates
  - JWT for authentication

- **Frontend**:
  - React
  - React Router for navigation
  - Axios for API requests
  - WebSocket for real-time updates
  - Minimal UI with inline styling (no extra UI libraries)

## Getting Started

### Prerequisites

- Go (1.16+)
- PostgreSQL
- Node.js (14+)
- npm or yarn

### Setup Instructions

#### Backend Setup

1. Navigate to the backend directory:
   ```
   cd backend
   ```

2. Install Go dependencies:
   ```
   go mod tidy
   ```

3. Create a `.env` file based on `.env.example`:
   ```
   cp .env.example .env
   ```

4. Set up your PostgreSQL database and update the credentials in the `.env` file.

5. Run the backend server:
   ```
   go run main.go
   ```

#### Frontend Setup

1. Navigate to the frontend directory:
   ```
   cd frontend
   ```

2. Install dependencies:
   ```
   npm install
   ```

3. Run the development server:
   ```
   npm start
   ```

4. The application will be available at `http://localhost:3000`

## API Endpoints

### Authentication

- `POST /api/auth/signup` - Register a new user and organization
- `POST /api/auth/login` - User login

### Services

- `GET /api/services` - Get all services for the user's organization
- `GET /api/services/:id` - Get a specific service
- `POST /api/services` - Create a new service
- `PUT /api/services/:id` - Update a service
- `DELETE /api/services/:id` - Delete a service

### Incidents

- `GET /api/incidents` - Get all incidents for the user's organization
- `GET /api/incidents/:id` - Get a specific incident
- `POST /api/incidents` - Create a new incident
- `PUT /api/incidents/:id` - Update an incident
- `DELETE /api/incidents/:id` - Delete an incident
- `POST /api/incidents/:id/updates` - Add an update to an incident

### Public Status Pages

- `GET /api/public/:orgId/services` - Get services for the public status page
- `GET /api/public/:orgId/incidents` - Get active incidents for the public status page

### WebSockets

- `GET /api/ws/:orgId` - WebSocket connection for real-time updates

## Important Notes

- Look for the `--->>here<<---` comments in the codebase, which highlight areas that might need your attention for:
  - Database configuration and connections
  - Authentication settings
  - WebSocket implementation
  - Other important parts of the application

## Project Structure

### Backend

```
backend/
├── api/            # API handlers
├── config/         # Configuration files
├── db/             # Database connection and migrations
├── middleware/     # Middleware (auth, logging, etc.)
├── models/         # Data models
├── services/       # Business logic services
├── utils/          # Utility functions
├── .env.example    # Environment variables example
├── go.mod          # Go modules
├── go.sum          # Go modules checksum
└── main.go         # Entry point
```

### Frontend

```
frontend/
├── public/         # Static files
└── src/
    ├── components/ # Reusable components
    ├── contexts/   # React contexts
    ├── hooks/      # Custom hooks
    ├── pages/      # Page components
    ├── services/   # API services
    ├── utils/      # Utility functions
    ├── App.js      # Main component
    └── index.js    # Entry point
```

## License

This project is open source and available under the [MIT License](LICENSE). 