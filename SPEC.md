# Client Portal - Project Specification

## Overview

Client portal for managing software projects, showing clients project progress, deliverables, and facilitating review workflows. Built for a software agency to replace scattered email/Slack communication with a unified project view.

## Problem Statement

Managing 7+ client projects with no central place for clients to see progress, review changes, or communicate. Currently using scattered emails, Slack, and Google Sheets for credentials.

## Target Users

- Agency owner/developers (admin)
- Clients viewing their project status

## Functional Requirements

### Auth & Users
- Users can login, logout, register, forgot password
- session based authentication 
- users can access authenticated content with their session 
- user/admin role 


### Projects
- Create project name, description, status
- Spec markdown field (agreed scope/requirements, rendered on frontend)
- Add clients and developers to projects
- Project members dev/client


### Environments & Services
- Project has multiple environments (Production, Staging, Development)
- Each environment has multiple services (Web, API, Queue, Worker, etc.)
- Services contain the actual URLs and BetterStack monitoring data
- Deliverables are deployed per environment (not per service)


### Deliverables
- Pre-planned feature slices/milestones
- Admin/developers post what they're building, then complete with documentation
- Completion includes images/screenshots and "how I completed it" writeup
- Deployed per environment tracking
- Clients can comment, admin/developer can respond
- Status: pending → in_progress → completed
- Fields: description, images, completion_description, status, completed_at, due_date 


### Reviews
- Ad-hoc change approval requests (not pre-planned like deliverables)
- Clients can accept or decline reviews
- Message thread for discussion/feedback
- Status: pending → approved/declined
- Lighter structure than deliverables

### Messages
- Client to dev/admin messaging
- live websocket, typing indicator
- Basically just a big group chat between devs/clients
- File upload (with size limits)
- Edit and delete messages


### Credentials
- Key, multiple value store for email accoutns and passwords
- only for developers to access 


### Notifications
- Simple external notifaction system plus notifaciton preferences in settings
- New develriabel, new review, etc


## Non-Functional Requirements

### Performance
- Page load < 2s
- API response < 200ms
- WebSocket message delivery < 500ms

### Security
- Argon2id password hashing
- Session-based auth (httpOnly cookies)
- HTTPS only
- Credentials encrypted at rest (AES-256)
- Role-based access control (admin/user + project member roles)

### Reliability
- Email delivery via Resend
- Messages persist when recipient offline
- Graceful WebSocket reconnect


## Entities
- users
- sessions
- projects
- project_meta (key-value JSONB for custom fields)
- project_members
- invites
- environments (id, project_id, name, order)
- services (id, environment_id, name, url, betterstack_monitor_id, status, uptime_percentage, response_time_ms, last_checked_at, ssl_expires_at)
- deliverables
- deliverable_environments (deliverable_id, environment_id, deployed_at)
- reviews
- messages
- credentials
- notification_preferences

## Architecture
- React + vite for frontend
- go for backend
- postgresql for db
- minio/r2/s3 for image storeage
- websockest
- resend for emails

## Auth System
- Argon2id for password hashing
- scs for session management

## API Endpoints
# Auth
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
POST   /api/v1/auth/forgot-password
POST   /api/v1/auth/reset-password

# Users
GET    /api/v1/users/me
GET    /api/v1/users?limit=20&cursor=abc123
PATCH  /api/v1/users/:id

# Projects
POST   /api/v1/projects
GET    /api/v1/projects?limit=20&cursor=abc123
GET    /api/v1/projects/:id
Patch  /api/v1/projects/:id
DELETE  /api/v1/projects/:id

# Project Meta (key-value pairs)
GET    /api/v1/projects/:id/meta
PUT    /api/v1/projects/:id/meta

# Project members
GET    /api/v1/projects/:id/members?limit=20&cursor=abc123
DELETE /api/v1/projects/:id/members/:userId

# Invites
POST  /api/v1/projects/:id/invites
GET   /api/v1/projects/:id/invites?limit=20&cursor=abc123
GET    /api/v1/invites/:token  # public validation
POST   /api/v1/invites/:token/accept     # accept
DELETE /api/v1/projects/:id/invites/:id

# Environments
GET    /api/v1/projects/:id/environments
POST   /api/v1/projects/:id/environments
PATCH  /api/v1/projects/:id/environments/:envId
DELETE /api/v1/projects/:id/environments/:envId

# Services
GET    /api/v1/projects/:id/environments/:envId/services
POST   /api/v1/projects/:id/environments/:envId/services
PATCH  /api/v1/projects/:id/environments/:envId/services/:svcId
DELETE /api/v1/projects/:id/environments/:envId/services/:svcId

# Deliverables
GET    /api/v1/projects/:id/deliverables
POST   /api/v1/projects/:id/deliverables
GET    /api/v1/projects/:id/deliverables/:delId
PATCH  /api/v1/projects/:id/deliverables/:delId
DELETE /api/v1/projects/:id/deliverables/:delId

# Deliverable Environments (mark deployed/undeployed to environment)
POST   /api/v1/projects/:id/deliverables/:delId/environments/:envId
DELETE /api/v1/projects/:id/deliverables/:delId/environments/:envId

# Reviews
GET    /api/v1/projects/:id/reviews
POST   /api/v1/projects/:id/reviews
GET    /api/v1/projects/:id/reviews/:reviewId
PATCH  /api/v1/projects/:id/reviews/:reviewId
DELETE /api/v1/projects/:id/reviews/:reviewId

# Messages (polymorphic - project chat & comments)
GET    /api/v1/projects/:id/messages?limit=20&cursor=abc123
POST   /api/v1/projects/:id/messages
PATCH  /api/v1/projects/:id/messages/:msgId
DELETE /api/v1/projects/:id/messages/:msgId
GET    /api/v1/projects/:id/deliverables/:delId/messages?limit=20&cursor=abc123
POST   /api/v1/projects/:id/deliverables/:delId/messages
PATCH  /api/v1/projects/:id/deliverables/:delId/messages/:msgId
DELETE /api/v1/projects/:id/deliverables/:delId/messages/:msgId
GET    /api/v1/projects/:id/reviews/:reviewId/messages?limit=20&cursor=abc123
POST   /api/v1/projects/:id/reviews/:reviewId/messages
PATCH  /api/v1/projects/:id/reviews/:reviewId/messages/:msgId
DELETE /api/v1/projects/:id/reviews/:reviewId/messages/:msgId

# Credentials (admin/dev only)
GET    /api/v1/projects/:id/credentials
POST   /api/v1/projects/:id/credentials
PATCH  /api/v1/projects/:id/credentials/:credId
DELETE /api/v1/projects/:id/credentials/:credId

# Notification Preferences
GET    /api/v1/users/me/notification-preferences
PATCH  /api/v1/users/me/notification-preferences

# WebSocket
WS     /api/v1/ws

# Webhooks
POST   /api/v1/webhooks/betterstack   # uptime status changes

## Third-Party Integrations
- Resend (transactional emails)
- BetterStack (uptime monitoring + webhooks)


## Screens
auth:
login
signup
forgot password
reset password

Admin:
Projects page


Client:
Projects page
view pr


## Feature Slices

### Slice 1: Auth + Users
- Session-based auth with Argon2id + SCS
- Login, logout, register
- Forgot password + reset password (Resend)
- GET /users/me, PATCH /users/:id
- Admin/user roles

**Done when:** Register, login, logout, reset password all work. Sessions persist.

### Slice 2: Projects + Members + Invites
- CRUD projects
- Project meta (JSONB)
- Invite tokens + accept flow
- Project members (developer/client roles)
- List projects filtered by membership

**Done when:** Admin creates project, sends invite link, user joins as dev/client.

### Slice 3: Environments + Services
- CRUD environments per project
- CRUD services per environment
- BetterStack webhook for status updates
- Display uptime, status, response time, SSL expiry

**Done when:** Add envs, add services, see live BetterStack data.

### Slice 4: Deliverables + Deployment Tracking
- CRUD deliverables
- Mark deployed to environments
- Comments (polymorphic messages)
- Status: pending → in_progress → completed

**Done when:** Post deliverable, mark deployed, client comments.

### Slice 5: Reviews
- CRUD review items (ad-hoc change approvals)
- Comments (polymorphic messages)
- Status: pending → approved/declined

**Done when:** Post review, client approves/declines with comments.

### Slice 6: Project Messages
- WebSocket real-time chat
- Typing indicators
- File uploads (S3/R2)
- Cursor pagination

**Done when:** Real-time chat works, files upload, history loads.

### Slice 7: Credentials
- AES-256 encrypted key-value store
- Developer-only access

**Done when:** Store/view credentials, clients blocked.

### Slice 8: Notifications
- Email via Resend (new deliverable, review, message)
- Notification preferences per user

**Done when:** Emails send, users can toggle off.

## MVP vs Future

MVP:
- Auth + password reset
- Projects + invites + members
- Environments + services + BetterStack
- Deliverables + deployment tracking
- Reviews
- Project messages (real-time)
- Credentials vault
- Email notifications

Future:
- File attachments on deliverables/feedback
- Activity log / audit trail
- Project archiving
- Dark mode
- Mobile app
- Slack/Discord notifications
- Time tracking
- Invoice generation
- Public status page

## Open Questions

(None remaining)
