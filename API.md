# VocabBunny Core API

> Base URL: `/api/v1`

All protected endpoints require a `Bearer` JWT token in the `Authorization` header. Tokens carry a **scope** (`app` or `bo`) and the server validates scope per route group.

---

## Table of Contents

- [Health](#health)
- [Authentication](#authentication)
- [Users (BO)](#users-bo)
- [Roles (BO)](#roles-bo)
- [Auth Identities (BO)](#auth-identities-bo)
- [Content — Order Nos (BO)](#content--order-nos-bo)
- [Content — Sections (BO)](#content--sections-bo)
- [Content — Lessons (BO)](#content--lessons-bo)
- [Content — Units (BO)](#content--units-bo)
- [Content — Question Sets (BO)](#content--question-sets-bo)
- [Content — Questions (BO)](#content--questions-bo)
- [Content — Question Choices (BO)](#content--question-choices-bo)
- [Content — Tags (BO)](#content--tags-bo)
- [Content — Media Assets (BO)](#content--media-assets-bo)
- [WebSocket](#websocket)
- [Common Structures](#common-structures)
- [Auth & Permissions Reference](#auth--permissions-reference)

---

## Health

### `GET /health`

No authentication required.

**Response `200`**
```json
{
  "status": "ok",
  "env": "production",
  "server_time": "2026-05-07T10:00:00Z",
  "timezone": "UTC",
  "hostname": "api-server-1"
}
```

---

## Authentication

Public endpoints — no token required.

### `POST /api/v1/app/auth/login/password`

Login as an **app user** (learner).

**Request**
```json
{
  "email_or_username": "rabbit@example.com",
  "password": "s3cret"
}
```

**Response `200`**
```json
{
  "access_token": "eyJ...",
  "refresh_token": "eyJ...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_expires_in": 2592000,
  "user": { "...": "UserResponse" }
}
```

---

### `POST /api/v1/app/auth/refresh`

Issue a new **app access token** from a valid app refresh token.

**Request**
```json
{
  "refresh_token": "eyJ..."
}
```

**Response `200`**
```json
{
  "access_token": "eyJ...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "user": { "...": "UserResponse" }
}
```

**Notes**

- Returns `401` when the refresh token is invalid or expired.
- Returns `403` when the refresh token scope does not match the route.

---

### `POST /api/v1/bo/auth/login/password`

Login as a **back-office user** (admin / content team). Returns a token scoped to `bo`.

**Request / Response** — same shape as app login above.

---

### `POST /api/v1/bo/auth/refresh`

Issue a new **back-office access token** from a valid `bo` refresh token.

**Request / Response** — same shape as app refresh above.

---

## Users (BO)

> Scope: `bo` | Base path: `/api/v1/bo/users`

| Method | Path | Permission / Role | Description |
|--------|------|-------------------|-------------|
| `GET` | `/bo/users` | `USER_READ` | List users with filters |
| `POST` | `/bo/users` | `ADMIN`, `CONTENT_ADMIN` | Create a user |
| `GET` | `/bo/users/by-email` | `USER_READ` | Look up user by email |
| `GET` | `/bo/users/by-username` | `USER_READ` | Look up user by username |
| `GET` | `/bo/users/:id` | `USER_READ` | Get user by ID |
| `PUT` | `/bo/users/:id` | `ADMIN`, `CONTENT_ADMIN` | Update user |
| `DELETE` | `/bo/users/:id` | `USER_BAN` | Delete user |

### Query Parameters — `GET /bo/users`

| Param | Type | Description |
|-------|------|-------------|
| `page` | int | Page number |
| `limit` | int | Items per page |
| `sort_by` | string | Field to sort by |
| `sort_order` | string | `ASC` or `DESC` |
| `search` | string | Search term |
| `role_id` | uuid | Filter by role |
| `status` | string | `INACTIVE` \| `ACTIVE` \| `BANNED` \| `DELETED` |
| `include_auth` | bool | Include auth identities in response |

### `GET /bo/users` — Response Example

```json
{
  "items": [],
  "paging": {
    "page": 1,
    "limit": 20,
    "total": 0,
    "total_pages": 0
  },
  "query": {
    "include_auth": false
  }
}
```

### `POST /bo/users` — Request Body

```json
{
  "email": "user@example.com",
  "username": "rabbituser",
  "display_name": "Rabbit User",
  "avatar_id": "uuid-optional",
  "status": "ACTIVE",
  "role_ids": ["uuid-role-1"]
}
```

### `PUT /bo/users/:id` — Request Body

All fields optional. Send only what should change.

```json
{
  "display_name": "New Name",
  "status": "BANNED"
}
```

### Delete Response `200`

```json
{ "id": "uuid", "status": "deleted" }
```

---

## Roles (BO)

> Scope: `bo` | Requires role: `ADMIN` | Base path: `/api/v1/bo/roles`

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/bo/roles` | List roles |
| `POST` | `/bo/roles` | Create role |
| `GET` | `/bo/roles/:id` | Get role by ID |
| `PUT` | `/bo/roles/:id` | Update role |
| `DELETE` | `/bo/roles/:id` | Delete role |

### `POST /bo/roles` — Request Body

```json
{
  "name": "CONTENT_ADMIN",
  "description": "Can manage all content",
  "permissions": [
    { "permission_code": "CONTENT_READ" },
    { "permission_code": "CONTENT_WRITE" },
    { "permission_code": "CONTENT_PUBLISH" }
  ]
}
```

Available `name` values: `ADMIN` | `CONTENT_ADMIN` | `MODERATOR` | `USER`

Available `permission_code` values: `CONTENT_READ` | `CONTENT_WRITE` | `CONTENT_PUBLISH` | `USER_READ` | `USER_BAN` | `ANALYTICS_READ` | `SYSTEM_CONFIG`

### Query Parameters — `GET /bo/roles`

| Param | Type | Description |
|-------|------|-------------|
| `search` | string | Search by name |
| `permission` | string | Filter by permission code |
| `page`, `limit`, `sort_by`, `sort_order` | — | Pagination & sorting |

---

## Auth Identities (BO)

> Scope: `bo` | Requires role: `ADMIN` | Base path: `/api/v1/bo/auth-identities`

An **auth identity** links a user to a login provider (password, Google, Apple).

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/bo/auth-identities` | List identities |
| `POST` | `/bo/auth-identities` | Create identity |
| `GET` | `/bo/auth-identities/:id` | Get identity by ID |
| `PUT` | `/bo/auth-identities/:id` | Update identity |
| `DELETE` | `/bo/auth-identities/:id` | Delete identity |

### `POST /bo/auth-identities` — Request Body

```json
{
  "user_id": "uuid",
  "provider": "PASSWORD",
  "password": "initial-password"
}
```

`provider` values: `PASSWORD` | `GOOGLE` | `APPLE`

### Query Parameters — `GET /bo/auth-identities`

| Param | Type | Description |
|-------|------|-------------|
| `user_id` | uuid | Filter by user |
| `provider` | string | Filter by provider |
| `page`, `limit`, `sort_by`, `sort_order` | — | Pagination & sorting |

### `GET /bo/auth-identities` — Response Example

```json
{
  "items": [],
  "paging": {
    "page": 1,
    "limit": 20,
    "total": 0,
    "total_pages": 0
  },
  "query": {}
}
```

---

## Content — Order Nos (BO)

> Scope: `bo` | Base path: `/api/v1/bo/content/order-nos`

Use this endpoint to fetch the current highest `order_no` for each content resource that supports manual ordering.

| Method | Path | Permission | Description |
|--------|------|------------|-------------|
| `GET` | `/bo/content/order-nos/last` | `CONTENT_READ` | Get the latest `order_no` grouped by content type |

### `GET /bo/content/order-nos/last` — Response Example

```json
{
  "sections": 12,
  "lessons": 34,
  "units": 20,
  "question_sets": 15,
  "questions": 87
}
```

---

## Content — Sections (BO)

> Scope: `bo` | Base path: `/api/v1/bo/content/sections`

Sections are the top level of the content hierarchy: **Section → Lesson → Unit → Question Set → Question**.

| Method | Path | Permission | Description |
|--------|------|------------|-------------|
| `GET` | `/bo/content/sections` | `CONTENT_READ` | List sections |
| `POST` | `/bo/content/sections` | `CONTENT_WRITE` | Create section |
| `GET` | `/bo/content/sections/:id` | `CONTENT_READ` | Get section by ID |
| `PUT` | `/bo/content/sections/:id` | `CONTENT_WRITE` | Update section |
| `DELETE` | `/bo/content/sections/:id` | `CONTENT_WRITE` | Delete section |

### `POST /bo/content/sections` — Request Body

```json
{
  "slug": "beginner-vocab",
  "title": "Beginner Vocabulary",
  "description": "Start here if you are new.",
  "order_no": 1,
  "is_published": false
}
```

### Query Parameters — `GET /bo/content/sections`

| Param | Type | Description |
|-------|------|-------------|
| `search` | string | Search in title / slug |
| `is_published` | bool | Filter by published status |
| `page`, `limit`, `sort_by`, `sort_order` | — | Pagination & sorting |

---

## Content — Lessons (BO)

> Scope: `bo` | Base path: `/api/v1/bo/content/lessons`

| Method | Path | Permission | Description |
|--------|------|------------|-------------|
| `GET` | `/bo/content/lessons` | `CONTENT_READ` | List lessons |
| `POST` | `/bo/content/lessons` | `CONTENT_WRITE` | Create lesson |
| `GET` | `/bo/content/lessons/:id` | `CONTENT_READ` | Get lesson by ID |
| `PUT` | `/bo/content/lessons/:id` | `CONTENT_WRITE` | Update lesson |
| `DELETE` | `/bo/content/lessons/:id` | `CONTENT_WRITE` | Delete lesson |

### `POST /bo/content/lessons` — Request Body

```json
{
  "section_id": "uuid",
  "slug": "colors-and-shapes",
  "title": "Colors & Shapes",
  "description": "Learn basic color and shape vocabulary.",
  "order_no": 1,
  "is_published": false
}
```

### Additional Query Parameter — `GET /bo/content/lessons`

| Param | Type | Description |
|-------|------|-------------|
| `section_id` | uuid | Filter by parent section |
| `is_published` | bool | Filter by published status |

---

## Content — Units (BO)

> Scope: `bo` | Base path: `/api/v1/bo/content/units`

| Method | Path | Permission | Description |
|--------|------|------------|-------------|
| `GET` | `/bo/content/units` | `CONTENT_READ` | List units |
| `POST` | `/bo/content/units` | `CONTENT_WRITE` | Create unit |
| `GET` | `/bo/content/units/:id` | `CONTENT_READ` | Get unit by ID |
| `PUT` | `/bo/content/units/:id` | `CONTENT_WRITE` | Update unit |
| `DELETE` | `/bo/content/units/:id` | `CONTENT_WRITE` | Delete unit |

### `POST /bo/content/units` — Request Body

```json
{
  "lesson_id": "uuid",
  "slug": "red-things",
  "title": "Red Things",
  "order_no": 1,
  "is_published": false
}
```

---

## Content — Question Sets (BO)

> Scope: `bo` | Base path: `/api/v1/bo/content/question-sets`

A question set groups questions within a unit and tracks a version number for curriculum iteration.

| Method | Path | Permission | Description |
|--------|------|------------|-------------|
| `GET` | `/bo/content/question-sets` | `CONTENT_READ` | List question sets |
| `POST` | `/bo/content/question-sets` | `CONTENT_WRITE` | Create question set |
| `GET` | `/bo/content/question-sets/:id` | `CONTENT_READ` | Get question set by ID |
| `PUT` | `/bo/content/question-sets/:id` | `CONTENT_WRITE` | Update question set |
| `DELETE` | `/bo/content/question-sets/:id` | `CONTENT_WRITE` | Delete question set |

### `POST /bo/content/question-sets` — Request Body

```json
{
  "unit_id": "uuid",
  "slug": "red-things-v1",
  "title": "Red Things — Practice Set",
  "estimated_seconds": 300,
  "order_no": 1,
  "version": 1,
  "is_published": false
}
```

### Additional Query Parameters — `GET /bo/content/question-sets`

| Param | Type | Description |
|-------|------|-------------|
| `unit_id` | uuid | Filter by parent unit |
| `version` | int | Filter by version |
| `is_published` | bool | Filter by published status |

---

## Content — Questions (BO)

> Scope: `bo` | Base path: `/api/v1/bo/content/questions`

| Method | Path | Permission | Description |
|--------|------|------------|-------------|
| `GET` | `/bo/content/questions` | `CONTENT_READ` | List questions |
| `POST` | `/bo/content/questions` | `CONTENT_WRITE` | Create question |
| `GET` | `/bo/content/questions/:id` | `CONTENT_READ` | Get question by ID |
| `PUT` | `/bo/content/questions/:id` | `CONTENT_WRITE` | Update question |
| `DELETE` | `/bo/content/questions/:id` | `CONTENT_WRITE` | Delete question |

### `POST /bo/content/questions` — Request Body

```json
{
  "question_set_id": "uuid",
  "type": "MULTIPLE_CHOICE",
  "question_text": "What color is the apple?",
  "explanation": "Apples are typically red.",
  "image_url": "https://cdn.example.com/apple.png",
  "difficulty": 1,
  "order_no": 1,
  "is_active": true,
  "choices": [
    { "choice_text": "Red", "choice_order": 1, "is_correct": true },
    { "choice_text": "Blue", "choice_order": 2, "is_correct": false }
  ],
  "tag_ids": ["uuid-tag-1"]
}
```

### Additional Query Parameters — `GET /bo/content/questions`

| Param | Type | Default | Description |
|-------|------|---------|-------------|
| `question_set_id` | uuid | — | Filter by question set |
| `type` | string | — | Filter by question type |
| `is_active` | bool | — | Filter by active status |
| `include_choices` | bool | `true` | Include choices in response |
| `include_tags` | bool | `true` | Include tags in response |

---

## Content — Question Choices (BO)

> Scope: `bo` | Base path: `/api/v1/bo/content/question-choices`

Manage individual answer choices independently of their parent question.

| Method | Path | Permission | Description |
|--------|------|------------|-------------|
| `GET` | `/bo/content/question-choices` | `CONTENT_READ` | List choices |
| `POST` | `/bo/content/question-choices` | `CONTENT_WRITE` | Create choice |
| `GET` | `/bo/content/question-choices/:id` | `CONTENT_READ` | Get choice by ID |
| `PUT` | `/bo/content/question-choices/:id` | `CONTENT_WRITE` | Update choice |
| `DELETE` | `/bo/content/question-choices/:id` | `CONTENT_WRITE` | Delete choice |

### `POST /bo/content/question-choices` — Request Body

```json
{
  "question_id": "uuid",
  "choice_text": "Red",
  "choice_order": 1,
  "is_correct": true
}
```

---

## Content — Tags (BO)

> Scope: `bo` | Base path: `/api/v1/bo/content/tags`

Tags are attached to questions for filtering and curriculum organization.

| Method | Path | Permission | Description |
|--------|------|------------|-------------|
| `GET` | `/bo/content/tags` | `CONTENT_READ` | List tags |
| `POST` | `/bo/content/tags` | `CONTENT_WRITE` | Create tag |
| `GET` | `/bo/content/tags/:id` | `CONTENT_READ` | Get tag by ID |
| `PUT` | `/bo/content/tags/:id` | `CONTENT_WRITE` | Update tag |
| `DELETE` | `/bo/content/tags/:id` | `CONTENT_WRITE` | Delete tag |

### `POST /bo/content/tags` — Request Body

```json
{ "name": "animals" }
```

---

## Content — Media Assets (BO)

> Scope: `bo` | Base path: `/api/v1/bo/content/media-assets`

Media assets support both **external** storage (S3, R2, GCS, local paths) and **database** storage (binary data).

| Method | Path | Permission | Description |
|--------|------|------------|-------------|
| `GET` | `/bo/content/media-assets` | `CONTENT_READ` | List media assets |
| `POST` | `/bo/content/media-assets` | `CONTENT_WRITE` | Create media asset |
| `GET` | `/bo/content/media-assets/:id` | `CONTENT_READ` | Get asset by ID |
| `PUT` | `/bo/content/media-assets/:id` | `CONTENT_WRITE` | Update asset |
| `DELETE` | `/bo/content/media-assets/:id` | `CONTENT_WRITE` | Delete asset |

### `POST /bo/content/media-assets` — Request Body

```json
{
  "asset_type": "IMAGE",
  "purpose": "QUESTION_IMAGE",
  "storage_mode": "EXTERNAL",
  "storage_provider": "S3",
  "bucket": "vocabunny-assets",
  "object_key": "questions/apple.png",
  "url": "https://cdn.example.com/apple.png",
  "mime_type": "image/png",
  "file_size_bytes": 48200,
  "is_public": true
}
```

**`asset_type`** — `IMAGE` | `VIDEO` | `DOCUMENT`

**`purpose`** — `AVATAR` | `QUESTION_IMAGE` | `BADGE_ICON` | `TROPHY_ICON` | `BANNER`

**`storage_mode`** — `EXTERNAL` | `DATABASE`

**`storage_provider`** — `S3` | `R2` | `GCS` | `LOCAL` (required when `storage_mode` is `EXTERNAL`)

### Query Parameters — `GET /bo/content/media-assets`

| Param | Type | Description |
|-------|------|-------------|
| `owner_user_id` | uuid | Filter by owning user |
| `asset_type` | string | Filter by asset type |
| `purpose` | string | Filter by purpose |
| `storage_mode` | string | Filter by storage mode |
| `storage_provider` | string | Filter by provider |
| `is_public` | bool | Filter by visibility |

---

## WebSocket

The server initialises a WebSocket manager (`app.Websocket`) but **no public WebSocket routes are registered yet**. Expect future endpoints here for real-time game state and leaderboards.

---

## Common Structures

### Paginated List Response

```json
{
  "items": [],
  "paging": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  },
  "query": {}
}
```

`total_pages` is calculated from `total` and `limit`. When `total` is `0`, `total_pages` is also `0`.

### Delete Response

```json
{ "id": "uuid", "status": "deleted" }
```

### Error Response

Standard HTTP status codes. Body shape varies by error type (validation errors include a `details` array).

```json
{
  "message": "validation failed",
  "details": [
    { "field": "email", "error": "must be a valid email address" }
  ]
}
```

---

## Auth & Permissions Reference

### Token Scopes

| Scope | Route group | Description |
|-------|-------------|-------------|
| `app` | `/api/v1/app/...` | End-user (learner) access |
| `bo` | `/api/v1/bo/...` | Back-office admin access |

### Roles

| Role | Description |
|------|-------------|
| `ADMIN` | Full system access |
| `CONTENT_ADMIN` | Content creation and management |
| `MODERATOR` | Moderation actions |
| `USER` | Basic learner access |

### Permission Codes

| Code | Grants |
|------|--------|
| `CONTENT_READ` | Read all content resources |
| `CONTENT_WRITE` | Create and update content |
| `CONTENT_PUBLISH` | Publish / unpublish content |
| `USER_READ` | Read user profiles |
| `USER_BAN` | Ban or delete users |
| `ANALYTICS_READ` | Access analytics data |
| `SYSTEM_CONFIG` | Modify system configuration |
