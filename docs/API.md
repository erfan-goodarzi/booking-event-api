# Event Booking API Documentation

API collection for the Event Booking service.

## Base URL

```
http://localhost:8080
```

Set the `baseUrl` variable to your server address. After logging in, the `token` variable is automatically populated via the Login or Refresh Token response scripts.

---

## Variables

| Variable | Default Value | Description |
|----------|--------------|-------------|
| `baseUrl` | `http://localhost:8080` | Base URL of the API server |
| `token` | *(auto-populated)* | JWT bearer token, set automatically after login |

---

## Authentication

Protected endpoints require a **Bearer Token** in the `Authorization` header:

```
Authorization: Bearer {{token}}
```

The token is automatically set after a successful **Login** or **Refresh Token** request.

---

## Endpoints

### Health

#### Health Check

```
GET {{baseUrl}}/health
```

Checks whether the API server is running and healthy.

**Authentication:** None required

**Request Headers:** None

**Response:** `200 OK` — server is healthy

---

### Auth

#### Signup

```
POST {{baseUrl}}/auth/signup
```

Registers a new user account.

**Authentication:** None required

**Request Headers:**

| Header | Value |
|--------|-------|
| `Content-Type` | `application/json` |

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "password123",
  "username": "johndoe"
}
```

| Field | Type | Description |
|-------|------|-------------|
| `email` | string | User's email address |
| `password` | string | User's password |
| `username` | string | Desired username |

---

#### Login

```
POST {{baseUrl}}/auth/login
```

Authenticates a user and returns a JWT token. The token is automatically saved to the `token` collection variable.

**Authentication:** None required

**Request Headers:**

| Header | Value |
|--------|-------|
| `Content-Type` | `application/json` |

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

| Field | Type | Description |
|-------|------|-------------|
| `email` | string | User's registered email |
| `password` | string | User's password |

**Response:** Returns a `token` field which is automatically stored in `{{token}}`.

---

#### Refresh Token

```
POST {{baseUrl}}/auth/refresh
```

Refreshes the current JWT token. The new token is automatically saved to the `token` collection variable.

**Authentication:** Bearer Token (`{{token}}`)

**Response:** Returns a new `token` which is automatically stored in `{{token}}`.

---

#### Logout

```
POST {{baseUrl}}/auth/logout
```

Logs out the current user and invalidates the session.

**Authentication:** Bearer Token (`{{token}}`)

---

### Events

#### Get All Events

```
GET {{baseUrl}}/events
```

Retrieves a list of all available events.

**Authentication:** None required

**Response:** `200 OK` — array of event objects

---

#### Get Event by ID

```
GET {{baseUrl}}/events/:id
```

Retrieves details for a specific event.

**Authentication:** None required

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | The unique ID of the event |

**Response:** `200 OK` — event object

---

#### Create Event

```
POST {{baseUrl}}/events
```

Creates a new event. Requires authentication.

**Authentication:** Bearer Token (`{{token}}`)

**Request Headers:**

| Header | Value |
|--------|-------|
| `Content-Type` | `application/json` |

**Request Body:**

```json
{
  "title": "Tech Conference 2026",
  "description": "Annual technology conference",
  "location": "San Francisco, CA",
  "dateTime": "2026-09-15T09:00:00Z",
  "duration": 480
}
```

| Field | Type | Description |
|-------|------|-------------|
| `title` | string | Title of the event |
| `description` | string | Description of the event |
| `location` | string | Physical or virtual location |
| `dateTime` | string (ISO 8601) | Start date and time of the event |
| `duration` | integer | Duration of the event in minutes |

**Response:** `201 Created` — created event object

---

#### Update Event

```
PUT {{baseUrl}}/events/:id
```

Updates an existing event. Requires authentication.

**Authentication:** Bearer Token (`{{token}}`)

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | The unique ID of the event to update |

**Request Headers:**

| Header | Value |
|--------|-------|
| `Content-Type` | `application/json` |

**Request Body:**

```json
{
  "title": "Updated Tech Conference 2026",
  "description": "Updated description",
  "location": "New York, NY",
  "dateTime": "2026-09-20T09:00:00Z",
  "duration": 360
}
```

| Field | Type | Description |
|-------|------|-------------|
| `title` | string | Updated title of the event |
| `description` | string | Updated description |
| `location` | string | Updated location |
| `dateTime` | string (ISO 8601) | Updated start date and time |
| `duration` | integer | Updated duration in minutes |

**Response:** `200 OK` — updated event object

---

#### Delete Event

```
DELETE {{baseUrl}}/events/:id
```

Deletes an existing event. Requires authentication.

**Authentication:** Bearer Token (`{{token}}`)

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | The unique ID of the event to delete |

**Response:** `200 OK` or `204 No Content`

---

### Tickets

#### Create Ticket

```
POST {{baseUrl}}/events/:id/tickets
```

Creates a ticket for a specific event. Requires authentication.

**Authentication:** Bearer Token (`{{token}}`)

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | The unique ID of the event |

**Request Headers:**

| Header | Value |
|--------|-------|
| `Content-Type` | `application/json` |

**Request Body:**

```json
{
  "type": "general",
  "price": 49.99,
  "quantity": 100
}
```

| Field | Type | Description |
|-------|------|-------------|
| `type` | string | Ticket type (e.g. `general`, `vip`) |
| `price` | number | Price per ticket |
| `quantity` | integer | Number of tickets available |

**Response:** `201 Created` — created ticket object

---

#### Register for Ticket

```
POST {{baseUrl}}/tickets/:id/register
```

Register the current user for a specific ticket. Requires authentication.

**Authentication:** Bearer Token (`{{token}}`)

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | The unique ID of the ticket |

**Request Headers:**

| Header | Value |
|--------|-------|
| `Content-Type` | `application/json` |

**Request Body:**

```json
{
  "user_id": "abc123"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `user_id` | string | Yes | The ID of the user to register |

**Responses:**

| Status | Description |
|--------|-------------|
| `200 OK` | Registration successful |
| `401 Unauthorized` | Missing or invalid token |
| `404 Not Found` | Ticket not found |

---

## Endpoint Summary

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/health` | Health check | No |
| `POST` | `/auth/signup` | Register a new user | No |
| `POST` | `/auth/login` | Login and get token | No |
| `POST` | `/auth/refresh` | Refresh JWT token | Yes |
| `POST` | `/auth/logout` | Logout | Yes |
| `GET` | `/events` | Get all events | No |
| `GET` | `/events/:id` | Get event by ID | No |
| `POST` | `/events` | Create a new event | Yes |
| `PUT` | `/events/:id` | Update an event | Yes |
| `DELETE` | `/events/:id` | Delete an event | Yes |
| `POST` | `/events/:id/tickets` | Create a ticket for an event | Yes |
| `POST` | `/tickets/:id/register` | Register for a ticket | Yes |