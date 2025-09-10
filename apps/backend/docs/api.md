# API Documentation

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

All protected endpoints require a JWT token in the Authorization header:

```http
Authorization: Bearer <jwt_token>
```

## Response Format

All API responses follow this format:

```json
{
  "success": true,
  "data": {...},
  "message": "Success message",
  "error": null
}
```

Error responses:

```json
{
  "success": false,
  "data": null,
  "message": "Error description",
  "error": {
    "code": "ERROR_CODE",
    "details": "Additional error details"
  }
}
```

## Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict
- `422` - Unprocessable Entity
- `500` - Internal Server Error

---

## Authentication Endpoints

### Register User

```http
POST /api/v1/users
```

**Request Body:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "secure_password123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "username": "john_doe",
    "email": "john@example.com",
    "role": "user"
  }
}
```

### Login

```http
POST /api/v1/auth/login
```

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "secure_password123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "username": "john_doe",
      "email": "john@example.com",
      "role": "user"
    },
    "token": "jwt_token_here"
  }
}
```

---

## User Endpoints

### Get User Profile

```http
GET /api/v1/users/profile
```

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "username": "john_doe",
    "email": "john@example.com",
    "avatar_url": "https://example.com/avatar.jpg",
    "role": "user",
    "is_verified": true,
    "created_at": "2024-01-01T00:00:00Z",
    "profile": {
      "about": "Anime enthusiast",
      "location": "Tokyo, Japan",
      "website": "https://johndoe.com",
      "birthdate": "1990-01-01T00:00:00Z"
    }
  }
}
```

### Update Profile

```http
PUT /api/v1/users/profile
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "about": "Updated bio",
  "location": "New location",
  "website": "https://newsite.com"
}
```

### Follow User

```http
POST /api/v1/users/follow
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "user_id": "target_user_uuid"
}
```

### Get Followers

```http
GET /api/v1/users/{user_id}/followers?limit=20&offset=0
```

### Get Following

```http
GET /api/v1/users/{user_id}/following?limit=20&offset=0
```

---

## Content Endpoints

### List Content

```http
GET /api/v1/content?type=anime&status=ongoing&limit=20&offset=0
```

**Query Parameters:**
- `type` - Content type: `anime`, `manga`, `movie`, `music`
- `status` - Content status: `ongoing`, `completed`, `upcoming`
- `genre` - Filter by genre
- `search` - Search query
- `sort` - Sort by: `popularity`, `rating`, `release_date`
- `limit` - Number of results (max 100)
- `offset` - Pagination offset

**Response:**
```json
{
  "success": true,
  "data": {
    "content": [
      {
        "id": "uuid",
        "title": "Attack on Titan",
        "description": "Epic anime series...",
        "type": "anime",
        "status": "completed",
        "genres": ["action", "drama"],
        "rating": 9.2,
        "release_date": "2013-04-07",
        "cover_url": "https://example.com/cover.jpg"
      }
    ],
    "total": 150,
    "limit": 20,
    "offset": 0
  }
}
```

### Get Content Details

```http
GET /api/v1/content/{content_id}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "title": "Attack on Titan",
    "description": "Detailed description...",
    "type": "anime",
    "status": "completed",
    "genres": ["action", "drama"],
    "rating": 9.2,
    "release_date": "2013-04-07",
    "cover_url": "https://example.com/cover.jpg",
    "metadata": {
      "episodes": 87,
      "duration": 24,
      "studio": "WIT Studio"
    },
    "stats": {
      "views": 1000000,
      "favorites": 50000,
      "reviews": 1250
    }
  }
}
```

### Add to Favorites

```http
POST /api/v1/content/{content_id}/favorite
```

**Headers:** `Authorization: Bearer <token>`

### Remove from Favorites

```http
DELETE /api/v1/content/{content_id}/favorite
```

**Headers:** `Authorization: Bearer <token>`

### Get User Favorites

```http
GET /api/v1/users/favorites?type=anime&limit=20&offset=0
```

**Headers:** `Authorization: Bearer <token>`

---

## Watch Party Endpoints

### Create Watch Party

```http
POST /api/v1/parties
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "content_id": "content_uuid",
  "name": "My Anime Night",
  "description": "Watching Attack on Titan together",
  "max_participants": 10,
  "is_public": true
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "party_uuid",
    "host_id": "user_uuid",
    "content_id": "content_uuid",
    "name": "My Anime Night",
    "description": "Watching Attack on Titan together",
    "status": "waiting",
    "max_participants": 10,
    "current_participants": 1,
    "is_public": true,
    "created_at": "2024-01-01T00:00:00Z",
    "join_code": "ABC123"
  }
}
```

### Join Watch Party

```http
POST /api/v1/parties/{party_id}/join
```

**Headers:** `Authorization: Bearer <token>`

**Alternative with join code:**
```http
POST /api/v1/parties/join
```

**Request Body:**
```json
{
  "join_code": "ABC123"
}
```

### Get Party Details

```http
GET /api/v1/parties/{party_id}
```

**Headers:** `Authorization: Bearer <token>`

### Update Party State

```http
POST /api/v1/parties/{party_id}/state
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "action": "play",
  "timestamp": 1234567,
  "episode": 5
}
```

**Actions:**
- `play` - Start playback
- `pause` - Pause playback
- `seek` - Seek to timestamp
- `next_episode` - Go to next episode
- `prev_episode` - Go to previous episode

### Leave Party

```http
POST /api/v1/parties/{party_id}/leave
```

**Headers:** `Authorization: Bearer <token>`

### Send Chat Message

```http
POST /api/v1/parties/{party_id}/messages
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "message": "This scene is amazing!",
  "timestamp": 1234567
}
```

### Get Chat Messages

```http
GET /api/v1/parties/{party_id}/messages?limit=50&offset=0
```

**Headers:** `Authorization: Bearer <token>`

---

## History Endpoints

### Update Watch Progress

```http
POST /api/v1/content/{content_id}/progress
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "episode": 5,
  "timestamp": 1234567,
  "completed": false
}
```

### Get Watch History

```http
GET /api/v1/users/history/watch?limit=20&offset=0
```

**Headers:** `Authorization: Bearer <token>`

### Get Read History

```http
GET /api/v1/users/history/read?limit=20&offset=0
```

**Headers:** `Authorization: Bearer <token>`

---

## Social Endpoints

### Create Discussion

```http
POST /api/v1/discussions
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "title": "What do you think about the latest episode?",
  "content": "Discussion content...",
  "content_id": "content_uuid"
}
```

### Get Discussions

```http
GET /api/v1/discussions?content_id=uuid&limit=20&offset=0
```

### Add Comment

```http
POST /api/v1/content/{content_id}/comments
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "comment": "Great episode!",
  "parent_id": "parent_comment_uuid" // Optional for replies
}
```

### Add Review

```http
POST /api/v1/content/{content_id}/reviews
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "rating": 9,
  "title": "Amazing series!",
  "content": "Detailed review...",
  "spoiler": false
}
```

---

## Device Management

### Register Device

```http
POST /api/v1/devices
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "device_name": "John's iPhone",
  "platform": "mobile"
}
```

### Get User Devices

```http
GET /api/v1/users/devices
```

**Headers:** `Authorization: Bearer <token>`

### Transfer Content

```http
POST /api/v1/devices/transfer
```

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "content_id": "content_uuid",
  "target_device_id": "device_uuid",
  "timestamp": 1234567
}
```

---

## Notifications

### Get Notifications

```http
GET /api/v1/notifications?unread_only=true&limit=20&offset=0
```

**Headers:** `Authorization: Bearer <token>`

### Mark as Read

```http
PUT /api/v1/notifications/{notification_id}/read
```

**Headers:** `Authorization: Bearer <token>`

### Mark All as Read

```http
PUT /api/v1/notifications/read
```

**Headers:** `Authorization: Bearer <token>`

---

## Health Check

### Health Status

```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T00:00:00Z",
  "version": "1.0.0",
  "services": {
    "database": "healthy",
    "external_apis": "healthy"
  }
}
```

---

## Rate Limiting

API requests are rate limited:

- **Authenticated users**: 1000 requests per hour
- **Unauthenticated users**: 100 requests per hour

Rate limit headers are included in responses:
- `X-RateLimit-Limit`: Request limit
- `X-RateLimit-Remaining`: Remaining requests
- `X-RateLimit-Reset`: Reset timestamp

---

## WebSocket Events

For real-time features like watch parties and notifications:

### Connection

```javascript
const ws = new WebSocket('ws://localhost:8080/ws?token=jwt_token');
```

### Watch Party Events

```json
{
  "type": "party_state_change",
  "data": {
    "party_id": "uuid",
    "action": "play",
    "timestamp": 1234567,
    "user_id": "uuid"
  }
}
```

### Notification Events

```json
{
  "type": "notification",
  "data": {
    "id": "uuid",
    "type": "follow",
    "title": "New follower",
    "message": "John started following you"
  }
}
```
