# LiveKit Video Conference Application - System Design & Specification

# Business Requirement
for backend golang server use the end to end layered architecture, handler, services, and repository for database etc. for email , make it interface based so we are flexible to use multiple providers like sendGrid, and mailjet etc. for now use these two email providers implementation. livekit API call can be in the service layers. design modular so that tehre are separation of concerns.This is a live video conference application using LiveKit (https://livekit.io) where I will be using livekit cloud that will give me a API_KEY, project_url which we can put in the .env file., Add environment variable .env file for the backend go lang service as part of the specifciation. Front end will contain directory structures like src/app/, src/services for api call, src/components for components, src/utils for util functions Business requirements: 1.) The front end will allow users to register and sign-up for the service. for now keep it simple just by retyping the password twice and entering email address, name and passwords. In future we can incorporate other fields like phone number, single sign on option too. for now keep it simple with username, email and password. 2.) Allow sign in to the service. At the backend it will validate password authentication and create jwt, right now keep simple jwt auth token, NO NEED for refresh token because that will make UI little complex. For now keep simple. 3. Allow password reset if user forgets password, there will be a link on the sign in page to reset password, which will have email field that user will enter. At server end with 'reset_password' end point, this will set token and send the user email address with the link to reset password from the email. Create the nextjs route and page for the reset password form when user clicks on the link sent by email. 3. when user logs-in, they will follow the 'auth' route pages in the nextjs UI routes and user will have ‘home’ page which will have a header with a drop down menu with items such as ‘profile settings’, ‘log out’ etc. under the header is the main section which will display the list of rooms, a search box at the top of the list of rom, a “+” button to add a new room.  3-a.  when user logs in, at the server side, the server service layer will also call the "token" end point of LiveKit cloud to get the token. for the logged in user. On the side bar there will be items like 'my rooms', profile settings, my invites. rooms, invites etc will have their own database tables. the rooms table will have fields like user_id, room_id, room_name(from UI form), another table that is required is room_participants table which will have columns such as room_id, email address, participant_id (a unique identifier). Focus on using LiveKit documentation on room management and participant management to add field to the room and participant_management to get a more sense on what other fields could be required . 4. For room management and adding participant top a room, on the UI when user clicks a room from the list of rooms in the home, the room details page will be displayed. Room details page will have a header section and a main section. 5.  Main section will have two section, like a messaging app window like slack type of app, with a list of posts ( which will be populated from a ‘posts’ table with fields such as id, creator_id, ‘message’, etc. in the first section, and at the bottom section will be a text box to enter a post by the user.  6. Header section will display room name, a button indicating “show members”,  room description,  button indicating “Add Members”,  and button for “Join Meeting" 7. Clicking on Members will open a right side sidebar that will display the list of participants, and also a button with “Add Button”.  The right side bar can be displayed or hidden based on the “Members” button click to toggle display and show members. On the sidebar,  there will be a section or table displaying participants details (name field, email field etc), at top right of the participants list section, a button to "Add participants" which will open a modal like component where user will enter the email address and name for the participants to add multiple participants through the modal.  8. At the backend api (in auth routes), /room, /<room_id>/participant will add participant. to add the participants. 9. when a participant is added from front end to backend, the participant will receive an URL in their email with the right structure, may me room id and participant id, so when the participant clicks on the link, they can join the LiveKit Room using the link to participate in the vidoe and audio conference session.  10. “Join Meeting” at the header will start a live kit video conference session using the React Component that live kit reacts SDK provides. The live video meeting will be happening through a next’s route page, when user clicks the “Join Meeting”, it will take to a nextjs page route called /prep_page, which will be similar to GoogleMeet screen where the user can perform some set up like background selection etc (I do not know how to do background selections in reacts, nextjs - please see if you have design and library to do available for this.), then this route page will have the “Join” button to finally join the meeting room, which will also be delivered from the nextjs route page /room. 11. At the end of the meeting the user will come back to the room home page which will display the header and the list of Posts.


# System Design and Specification

## 1. Technology Stack

### Frontend
- **Framework**: Next.js 14+ with App Router
- **UI Library**: React 18+
- **API Communication**: Fetch API
- **LiveKit SDK**: @livekit/components-react, livekit-client
- **Styling**: Tailwind CSS

### Backend
- **Language**: Go 1.21+
- **Database**: PostgreSQL 15+
- **Architecture**: Layered (Handler → Service → Repository)
- **Authentication**: JWT (simple, no refresh tokens)
- **Email**: Interface-based (SendGrid, Mailjet)
- **LiveKit**: LiveKit Server SDK for Go

---

## 2. Directory Structure

### 2.1 Frontend Structure (Next.js)

```
frontend/
├── src/
│   ├── app/
│   │   ├── layout.tsx
│   │   ├── page.tsx
│   │   ├── auth/
│   │   │   ├── signup/
│   │   │   │   └── page.tsx
│   │   │   ├── signin/
│   │   │   │   └── page.tsx
│   │   │   ├── reset-password/
│   │   │   │   └── page.tsx
│   │   │   └── reset-password-form/
│   │   │       └── page.tsx
│   │   ├── home/
│   │   │   ├── page.tsx
│   │   │   └── layout.tsx
│   │   ├── room/
│   │   │   └── [roomId]/
│   │   │       ├── page.tsx
│   │   │       ├── prep/
│   │   │       │   └── page.tsx
│   │   │       └── meeting/
│   │   │           └── page.tsx
│   │   └── profile/
│   │       └── page.tsx
│   ├── components/
│   │   ├── auth/
│   │   │   ├── SignUpForm.tsx
│   │   │   ├── SignInForm.tsx
│   │   │   └── ResetPasswordForm.tsx
│   │   ├── layout/
│   │   │   ├── Header.tsx
│   │   │   ├── Sidebar.tsx
│   │   │   └── DropdownMenu.tsx
│   │   ├── room/
│   │   │   ├── RoomList.tsx
│   │   │   ├── RoomCard.tsx
│   │   │   ├── RoomDetails.tsx
│   │   │   ├── RoomHeader.tsx
│   │   │   ├── ParticipantsSidebar.tsx
│   │   │   ├── AddParticipantModal.tsx
│   │   │   └── PostsList.tsx
│   │   ├── meeting/
│   │   │   ├── VideoConference.tsx
│   │   │   ├── PrepRoom.tsx
│   │   │   └── BackgroundSelector.tsx
│   │   └── common/
│   │       ├── Button.tsx
│   │       ├── Input.tsx
│   │       └── Modal.tsx
│   ├── services/
│   │   ├── api/
│   │   │   ├── auth.service.ts
│   │   │   ├── room.service.ts
│   │   │   ├── participant.service.ts
│   │   │   ├── post.service.ts
│   │   │   └── livekit.service.ts
│   │   └── storage/
│   │       └── auth.storage.ts
│   ├── utils/
│   │   ├── api.ts
│   │   ├── constants.ts
│   │   └── validation.ts
│   ├── types/
│   │   ├── auth.types.ts
│   │   ├── room.types.ts
│   │   └── user.types.ts
│   └── middleware.ts
├── public/
├── .env.local
├── next.config.js
├── package.json
└── tsconfig.json
```

### 2.2 Backend Structure (Go)

```
backend/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── room_handler.go
│   │   ├── participant_handler.go
│   │   └── post_handler.go
│   ├── service/
│   │   ├── auth_service.go
│   │   ├── room_service.go
│   │   ├── participant_service.go
│   │   ├── post_service.go
│   │   ├── livekit_service.go
│   │   └── email/
│   │       ├── email_service.go
│   │       ├── sendgrid_provider.go
│   │       └── mailjet_provider.go
│   ├── repository/
│   │   ├── user_repository.go
│   │   ├── room_repository.go
│   │   ├── participant_repository.go
│   │   └── post_repository.go
│   ├── model/
│   │   ├── user.go
│   │   ├── room.go
│   │   ├── participant.go
│   │   ├── post.go
│   │   └── invite.go
│   ├── middleware/
│   │   ├── auth_middleware.go
│   │   └── cors_middleware.go
│   ├── database/
│   │   ├── postgres.go
│   │   └── migrations/
│   │       ├── 001_create_users_table.sql
│   │       ├── 002_create_rooms_table.sql
│   │       ├── 003_create_participants_table.sql
│   │       ├── 004_create_posts_table.sql
│   │       └── 005_create_invites_table.sql
│   └── utils/
│       ├── jwt.go
│       ├── password.go
│       └── validator.go
├── pkg/
│   └── logger/
│       └── logger.go
├── .env
├── go.mod
├── go.sum
└── Makefile
```

---

## 3. Database Schema

### 3.1 Users Table

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP,
    is_active BOOLEAN DEFAULT true
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
```

### 3.2 Password Reset Tokens Table

```sql
CREATE TABLE password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_reset_tokens_token ON password_reset_tokens(token);
CREATE INDEX idx_reset_tokens_user_id ON password_reset_tokens(user_id);
```

### 3.3 Rooms Table

```sql
CREATE TABLE rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_name VARCHAR(255) NOT NULL,
    room_sid VARCHAR(255), -- LiveKit room SID
    description TEXT,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    livekit_room_name VARCHAR(255),
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true
);

CREATE INDEX idx_rooms_owner_id ON rooms(owner_id);
CREATE INDEX idx_rooms_room_sid ON rooms(room_sid);
```

### 3.4 Room Participants Table

```sql
CREATE TABLE room_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    participant_id UUID, -- Can be NULL for external participants
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    livekit_identity VARCHAR(255), -- LiveKit participant identity
    role VARCHAR(50) DEFAULT 'participant', -- owner, moderator, participant
    joined_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    UNIQUE(room_id, email)
);

CREATE INDEX idx_participants_room_id ON room_participants(room_id);
CREATE INDEX idx_participants_user_id ON room_participants(user_id);
CREATE INDEX idx_participants_email ON room_participants(email);
```

### 3.5 Posts Table

```sql
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    creator_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT false
);

CREATE INDEX idx_posts_room_id ON posts(room_id);
CREATE INDEX idx_posts_creator_id ON posts(creator_id);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
```

### 3.6 Invites Table

```sql
CREATE TABLE invites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    inviter_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    invitee_email VARCHAR(255) NOT NULL,
    invitee_name VARCHAR(255),
    token VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(50) DEFAULT 'pending', -- pending, accepted, expired
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    accepted_at TIMESTAMP
);

CREATE INDEX idx_invites_room_id ON invites(room_id);
CREATE INDEX idx_invites_token ON invites(token);
CREATE INDEX idx_invites_email ON invites(invitee_email);
```

---

## 4. Backend Go Implementation

### 4.1 Environment Variables (.env)

```env
# Server Configuration
PORT=8080
ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=livekit_app
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRY_HOURS=24

# LiveKit Configuration
LIVEKIT_API_KEY=your-livekit-api-key
LIVEKIT_API_SECRET=your-livekit-api-secret
LIVEKIT_URL=wss://your-project.livekit.cloud

# Email Configuration
EMAIL_PROVIDER=sendgrid # or mailjet
SENDGRID_API_KEY=your-sendgrid-api-key
SENDGRID_FROM_EMAIL=noreply@yourapp.com
SENDGRID_FROM_NAME=LiveKit App
MAILJET_API_KEY=your-mailjet-api-key
MAILJET_SECRET_KEY=your-mailjet-secret-key
MAILJET_FROM_EMAIL=noreply@yourapp.com
MAILJET_FROM_NAME=LiveKit App

# Frontend URL
FRONTEND_URL=http://localhost:3000

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000
```

### 4.2 Data Models

#### model/user.go

```go
package model

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID           uuid.UUID  `json:"id" db:"id"`
    Username     string     `json:"username" db:"username"`
    Email        string     `json:"email" db:"email"`
    PasswordHash string     `json:"-" db:"password_hash"`
    Name         string     `json:"name" db:"name"`
    CreatedAt    time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
    LastLogin    *time.Time `json:"last_login" db:"last_login"`
    IsActive     bool       `json:"is_active" db:"is_active"`
}

type UserSignUpRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Name     string `json:"name" validate:"required,min=2"`
    Password string `json:"password" validate:"required,min=8"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type UserSignInRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
    Token        string    `json:"token"`
    User         User      `json:"user"`
    LiveKitToken string    `json:"livekit_token"`
    ExpiresAt    time.Time `json:"expires_at"`
}

type PasswordResetRequest struct {
    Email string `json:"email" validate:"required,email"`
}

type PasswordResetConfirm struct {
    Token       string `json:"token" validate:"required"`
    NewPassword string `json:"new_password" validate:"required,min=8"`
}
```

#### model/room.go

```go
package model

import (
    "time"
    "github.com/google/uuid"
)

type Room struct {
    ID              uuid.UUID              `json:"id" db:"id"`
    RoomName        string                 `json:"room_name" db:"room_name"`
    RoomSID         *string                `json:"room_sid" db:"room_sid"`
    Description     *string                `json:"description" db:"description"`
    OwnerID         uuid.UUID              `json:"owner_id" db:"owner_id"`
    LiveKitRoomName *string                `json:"livekit_room_name" db:"livekit_room_name"`
    Metadata        map[string]interface{} `json:"metadata" db:"metadata"`
    CreatedAt       time.Time              `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time              `json:"updated_at" db:"updated_at"`
    IsActive        bool                   `json:"is_active" db:"is_active"`
}

type CreateRoomRequest struct {
    RoomName    string  `json:"room_name" validate:"required,min=3,max=100"`
    Description *string `json:"description"`
}

type RoomResponse struct {
    Room             Room                `json:"room"`
    ParticipantCount int                 `json:"participant_count"`
    IsOwner          bool                `json:"is_owner"`
}
```

#### model/participant.go

```go
package model

import (
    "time"
    "github.com/google/uuid"
)

type RoomParticipant struct {
    ID               uuid.UUID  `json:"id" db:"id"`
    RoomID           uuid.UUID  `json:"room_id" db:"room_id"`
    ParticipantID    *uuid.UUID `json:"participant_id" db:"participant_id"`
    UserID           *uuid.UUID `json:"user_id" db:"user_id"`
    Email            string     `json:"email" db:"email"`
    Name             string     `json:"name" db:"name"`
    LiveKitIdentity  *string    `json:"livekit_identity" db:"livekit_identity"`
    Role             string     `json:"role" db:"role"`
    JoinedAt         *time.Time `json:"joined_at" db:"joined_at"`
    CreatedAt        time.Time  `json:"created_at" db:"created_at"`
    IsActive         bool       `json:"is_active" db:"is_active"`
}

type AddParticipantRequest struct {
    Email string `json:"email" validate:"required,email"`
    Name  string `json:"name" validate:"required,min=2"`
}

type ParticipantInviteResponse struct {
    ParticipantID uuid.UUID `json:"participant_id"`
    InviteToken   string    `json:"invite_token"`
    InviteURL     string    `json:"invite_url"`
}
```

#### model/post.go

```go
package model

import (
    "time"
    "github.com/google/uuid"
)

type Post struct {
    ID        uuid.UUID `json:"id" db:"id"`
    RoomID    uuid.UUID `json:"room_id" db:"room_id"`
    CreatorID uuid.UUID `json:"creator_id" db:"creator_id"`
    Message   string    `json:"message" db:"message"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
    IsDeleted bool      `json:"is_deleted" db:"is_deleted"`
}

type CreatePostRequest struct {
    Message string `json:"message" validate:"required,min=1,max=5000"`
}

type PostWithCreator struct {
    Post
    CreatorName string `json:"creator_name"`
}
```

### 4.3 Email Service Interface

#### service/email/email_service.go

```go
package email

import (
    "context"
)

// EmailProvider is the interface that all email providers must implement
type EmailProvider interface {
    SendEmail(ctx context.Context, to, subject, htmlContent, textContent string) error
    SendTemplateEmail(ctx context.Context, to string, templateID string, data map[string]interface{}) error
}

// EmailService handles email operations
type EmailService struct {
    provider EmailProvider
    fromEmail string
    fromName  string
}

func NewEmailService(provider EmailProvider, fromEmail, fromName string) *EmailService {
    return &EmailService{
        provider:  provider,
        fromEmail: fromEmail,
        fromName:  fromName,
    }
}

func (s *EmailService) SendPasswordResetEmail(ctx context.Context, to, resetToken, resetURL string) error {
    subject := "Password Reset Request"
    htmlContent := `
        <html>
        <body>
            <h2>Password Reset Request</h2>
            <p>You have requested to reset your password. Click the link below to reset:</p>
            <a href="` + resetURL + `">Reset Password</a>
            <p>This link will expire in 1 hour.</p>
            <p>If you did not request this, please ignore this email.</p>
        </body>
        </html>
    `
    textContent := "Password reset link: " + resetURL
    
    return s.provider.SendEmail(ctx, to, subject, htmlContent, textContent)
}

func (s *EmailService) SendRoomInviteEmail(ctx context.Context, to, roomName, inviteURL string) error {
    subject := "You've been invited to join a meeting room"
    htmlContent := `
        <html>
        <body>
            <h2>Meeting Room Invitation</h2>
            <p>You have been invited to join the room: <strong>` + roomName + `</strong></p>
            <p>Click the link below to join:</p>
            <a href="` + inviteURL + `">Join Room</a>
        </body>
        </html>
    `
    textContent := "Join room " + roomName + ": " + inviteURL
    
    return s.provider.SendEmail(ctx, to, subject, htmlContent, textContent)
}
```

#### service/email/sendgrid_provider.go

```go
package email

import (
    "context"
    "github.com/sendgrid/sendgrid-go"
    "github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridProvider struct {
    apiKey    string
    fromEmail string
    fromName  string
}

func NewSendGridProvider(apiKey, fromEmail, fromName string) *SendGridProvider {
    return &SendGridProvider{
        apiKey:    apiKey,
        fromEmail: fromEmail,
        fromName:  fromName,
    }
}

func (p *SendGridProvider) SendEmail(ctx context.Context, to, subject, htmlContent, textContent string) error {
    from := mail.NewEmail(p.fromName, p.fromEmail)
    toEmail := mail.NewEmail("", to)
    message := mail.NewSingleEmail(from, subject, toEmail, textContent, htmlContent)
    
    client := sendgrid.NewSendClient(p.apiKey)
    _, err := client.Send(message)
    return err
}

func (p *SendGridProvider) SendTemplateEmail(ctx context.Context, to string, templateID string, data map[string]interface{}) error {
    // Implement template-based sending if needed
    return nil
}
```

#### service/email/mailjet_provider.go

```go
package email

import (
    "context"
    mailjet "github.com/mailjet/mailjet-apiv3-go"
)

type MailjetProvider struct {
    client    *mailjet.Client
    fromEmail string
    fromName  string
}

func NewMailjetProvider(apiKey, secretKey, fromEmail, fromName string) *MailjetProvider {
    client := mailjet.NewMailjetClient(apiKey, secretKey)
    return &MailjetProvider{
        client:    client,
        fromEmail: fromEmail,
        fromName:  fromName,
    }
}

func (p *MailjetProvider) SendEmail(ctx context.Context, to, subject, htmlContent, textContent string) error {
    messagesInfo := []mailjet.InfoMessagesV31{
        {
            From: &mailjet.RecipientV31{
                Email: p.fromEmail,
                Name:  p.fromName,
            },
            To: &mailjet.RecipientsV31{
                mailjet.RecipientV31{
                    Email: to,
                },
            },
            Subject:  subject,
            TextPart: textContent,
            HTMLPart: htmlContent,
        },
    }
    
    messages := mailjet.MessagesV31{Info: messagesInfo}
    _, err := p.client.SendMailV31(&messages)
    return err
}

func (p *MailjetProvider) SendTemplateEmail(ctx context.Context, to string, templateID string, data map[string]interface{}) error {
    // Implement template-based sending if needed
    return nil
}
```

### 4.4 LiveKit Service

#### service/livekit_service.go

```go
package service

import (
    "context"
    "time"
    
    "github.com/livekit/protocol/auth"
    lksdk "github.com/livekit/server-sdk-go"
)

type LiveKitService struct {
    apiKey    string
    apiSecret string
    url       string
}

func NewLiveKitService(apiKey, apiSecret, url string) *LiveKitService {
    return &LiveKitService{
        apiKey:    apiKey,
        apiSecret: apiSecret,
        url:       url,
    }
}

// GenerateToken creates a LiveKit access token for a participant
func (s *LiveKitService) GenerateToken(identity, roomName string, canPublish, canSubscribe bool) (string, error) {
    at := auth.NewAccessToken(s.apiKey, s.apiSecret)
    grant := &auth.VideoGrant{
        RoomJoin:     true,
        Room:         roomName,
        CanPublish:   &canPublish,
        CanSubscribe: &canSubscribe,
    }
    at.AddGrant(grant).
        SetIdentity(identity).
        SetValidFor(24 * time.Hour)
    
    return at.ToJWT()
}

// CreateRoom creates a new room in LiveKit
func (s *LiveKitService) CreateRoom(ctx context.Context, roomName string) (*lksdk.Room, error) {
    roomClient := lksdk.NewRoomServiceClient(s.url, s.apiKey, s.apiSecret)
    
    room, err := roomClient.CreateRoom(ctx, &lksdk.CreateRoomRequest{
        Name:         roomName,
        EmptyTimeout: 300, // 5 minutes
        MaxParticipants: 50,
    })
    
    return room, err
}

// ListParticipants returns all participants in a room
func (s *LiveKitService) ListParticipants(ctx context.Context, roomName string) ([]*lksdk.ParticipantInfo, error) {
    roomClient := lksdk.NewRoomServiceClient(s.url, s.apiKey, s.apiSecret)
    
    participants, err := roomClient.ListParticipants(ctx, &lksdk.ListParticipantsRequest{
        Room: roomName,
    })
    
    if err != nil {
        return nil, err
    }
    
    return participants, nil
}

// DeleteRoom removes a room from LiveKit
func (s *LiveKitService) DeleteRoom(ctx context.Context, roomName string) error {
    roomClient := lksdk.NewRoomServiceClient(s.url, s.apiKey, s.apiSecret)
    
    _, err := roomClient.DeleteRoom(ctx, &lksdk.DeleteRoomRequest{
        Room: roomName,
    })
    
    return err
}
```

### 4.5 Repository Layer (Pseudo Code)

#### repository/user_repository.go

```go
package repository

import (
    "context"
    "database/sql"
    "github.com/google/uuid"
    "yourapp/internal/model"
)

type UserRepository interface {
    Create(ctx context.Context, user *model.User) error
    GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
    GetByEmail(ctx context.Context, email string) (*model.User, error)
    GetByUsername(ctx context.Context, username string) (*model.User, error)
    Update(ctx context.Context, user *model.User) error
    UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
}

type userRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
    query := `
        INSERT INTO users (username, email, password_hash, name)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at, is_active
    `
    // Execute query and scan result
    return nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
    query := `
        SELECT id, username, email, password_hash, name, created_at, updated_at, last_login, is_active
        FROM users
        WHERE email = $1 AND is_active = true
    `
    // Execute query and scan result
    return nil, nil
}

// Other methods...
```

#### repository/room_repository.go

```go
package repository

import (
    "context"
    "database/sql"
    "github.com/google/uuid"
    "yourapp/internal/model"
)

type RoomRepository interface {
    Create(ctx context.Context, room *model.Room) error
    GetByID(ctx context.Context, id uuid.UUID) (*model.Room, error)
    GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*model.Room, error)
    Update(ctx context.Context, room *model.Room) error
    Delete(ctx context.Context, id uuid.UUID) error
    GetRoomsByUser(ctx context.Context, userID uuid.UUID) ([]*model.Room, error)
}

type roomRepository struct {
    db *sql.DB
}

func NewRoomRepository(db *sql.DB) RoomRepository {
    return &roomRepository{db: db}
}

func (r *roomRepository) Create(ctx context.Context, room *model.Room) error {
    query := `
        INSERT INTO rooms (room_name, description, owner_id, livekit_room_name, room_sid)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at, is_active
    `
    // Execute query
    return nil
}

func (r *roomRepository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*model.Room, error) {
    query := `
        SELECT id, room_name, room_sid, description, owner_id, livekit_room_name, 
               metadata, created_at, updated_at, is_active
        FROM rooms
        WHERE owner_id = $1 AND is_active = true
        ORDER BY created_at DESC
    `
    // Execute query and scan results
    return nil, nil
}

// Other methods...
```

### 4.6 Service Layer (Pseudo Code)

#### service/auth_service.go

```go
package service

import (
    "context"
    "errors"
    "time"
    "github.com/google/uuid"
    "yourapp/internal/model"
    "yourapp/internal/repository"
    "yourapp/internal/utils"
    "yourapp/internal/service/email"
)

type AuthService struct {
    userRepo       repository.UserRepository
    resetTokenRepo repository.PasswordResetTokenRepository
    emailService   *email.EmailService
    livekitService *LiveKitService
    jwtSecret      string
    frontendURL    string
}

func NewAuthService(
    userRepo repository.UserRepository,
    resetTokenRepo repository.PasswordResetTokenRepository,
    emailService *email.EmailService,
    livekitService *LiveKitService,
    jwtSecret, frontendURL string,
) *AuthService {
    return &AuthService{
        userRepo:       userRepo,
        resetTokenRepo: resetTokenRepo,
        emailService:   emailService,
        livekitService: livekitService,
        jwtSecret:      jwtSecret,
        frontendURL:    frontendURL,
    }
}

func (s *AuthService) SignUp(ctx context.Context, req *model.UserSignUpRequest) error {
    // 1. Validate request
    if req.Password != req.ConfirmPassword {
        return errors.New("passwords do not match")
    }
    
    // 2. Check if user exists
    existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
    if existingUser != nil {
        return errors.New("email already exists")
    }
    
    // 3. Hash password
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return err
    }
    
    // 4. Create user
    user := &model.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: hashedPassword,
        Name:         req.Name,
    }
    
    return s.userRepo.Create(ctx, user)
}

func (s *AuthService) SignIn(ctx context.Context, req *model.UserSignInRequest) (*model.AuthResponse, error) {
    // 1. Get user by email
    user, err := s.userRepo.GetByEmail(ctx, req.Email)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }
    
    // 2. Verify password
    if !utils.CheckPassword(req.Password, user.PasswordHash) {
        return nil, errors.New("invalid credentials")
    }
    
    // 3. Generate JWT token
    token, expiresAt, err := utils.GenerateJWT(user.ID.String(), user.Email, s.jwtSecret)
    if err != nil {
        return nil, err
    }
    
    // 4. Generate LiveKit token
    livekitToken, err := s.livekitService.GenerateToken(
        user.ID.String(),
        "default", // Default room or user-specific
        true,      // Can publish
        true,      // Can subscribe
    )
    if err != nil {
        return nil, err
    }
    
    // 5. Update last login
    s.userRepo.UpdateLastLogin(ctx, user.ID)
    
    return &model.AuthResponse{
        Token:        token,
        User:         *user,
        LiveKitToken: livekitToken,
        ExpiresAt:    expiresAt,
    }, nil
}

func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
    // 1. Get user by email
    user, err := s.userRepo.GetByEmail(ctx, email)
    if err != nil {
        // Don't reveal if user exists
        return nil
    }
    
    // 2. Generate reset token
    resetToken := uuid.New().String()
    expiresAt := time.Now().Add(1 * time.Hour)
    
    // 3. Save reset token
    err = s.resetTokenRepo.Create(ctx, user.ID, resetToken, expiresAt)
    if err != nil {
        return err
    }
    
    // 4. Send reset email
    resetURL := s.frontendURL + "/auth/reset-password-form?token=" + resetToken
    return s.emailService.SendPasswordResetEmail(ctx, user.Email, resetToken, resetURL)
}

func (s *AuthService) ResetPassword(ctx context.Context, req *model.PasswordResetConfirm) error {
    // 1. Verify token
    resetToken, err := s.resetTokenRepo.GetByToken(ctx, req.Token)
    if err != nil {
        return errors.New("invalid or expired token")
    }
    
    if resetToken.Used || time.Now().After(resetToken.ExpiresAt) {
        return errors.New("invalid or expired token")
    }
    
    // 2. Get user
    user, err := s.userRepo.GetByID(ctx, resetToken.UserID)
    if err != nil {
        return err
    }
    
    // 3. Hash new password
    hashedPassword, err := utils.HashPassword(req.NewPassword)
    if err != nil {
        return err
    }
    
    // 4. Update user password
    user.PasswordHash = hashedPassword
    err = s.userRepo.Update(ctx, user)
    if err != nil {
        return err
    }
    
    // 5. Mark token as used
    return s.resetTokenRepo.MarkAsUsed(ctx, resetToken.ID)
}
```

#### service/room_service.go

```go
package service

import (
    "context"
    "errors"
    "github.com/google/uuid"
    "yourapp/internal/model"
    "yourapp/internal/repository"
)

type RoomService struct {
    roomRepo        repository.RoomRepository
    participantRepo repository.ParticipantRepository
    livekitService  *LiveKitService
}

func NewRoomService(
    roomRepo repository.RoomRepository,
    participantRepo repository.ParticipantRepository,
    livekitService *LiveKitService,
) *RoomService {
    return &RoomService{
        roomRepo:        roomRepo,
        participantRepo: participantRepo,
        livekitService:  livekitService,
    }
}

func (s *RoomService) CreateRoom(ctx context.Context, userID uuid.UUID, req *model.CreateRoomRequest) (*model.Room, error) {
    // 1. Generate unique LiveKit room name
    livekitRoomName := "room_" + uuid.New().String()
    
    // 2. Create room in LiveKit
    lkRoom, err := s.livekitService.CreateRoom(ctx, livekitRoomName)
    if err != nil {
        return nil, err
    }
    
    // 3. Create room in database
    room := &model.Room{
        RoomName:        req.RoomName,
        Description:     req.Description,
        OwnerID:         userID,
        LiveKitRoomName: &livekitRoomName,
        RoomSID:         &lkRoom.Sid,
    }
    
    err = s.roomRepo.Create(ctx, room)
    if err != nil {
        // Cleanup LiveKit room if DB creation fails
        s.livekitService.DeleteRoom(ctx, livekitRoomName)
        return nil, err
    }
    
    // 4. Add owner as participant
    err = s.participantRepo.Create(ctx, &model.RoomParticipant{
        RoomID: room.ID,
        UserID: &userID,
        Email:  "", // Get from user context
        Name:   "", // Get from user context
        Role:   "owner",
    })
    
    return room, nil
}

func (s *RoomService) GetUserRooms(ctx context.Context, userID uuid.UUID) ([]*model.RoomResponse, error) {
    // 1. Get rooms owned by user
    ownedRooms, err := s.roomRepo.GetByOwnerID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // 2. Get rooms where user is participant
    participantRooms, err := s.roomRepo.GetRoomsByUser(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // 3. Combine and format response
    roomMap := make(map[uuid.UUID]*model.RoomResponse)
    
    for _, room := range ownedRooms {
        count, _ := s.participantRepo.CountByRoomID(ctx, room.ID)
        roomMap[room.ID] = &model.RoomResponse{
            Room:             *room,
            ParticipantCount: count,
            IsOwner:          true,
        }
    }
    
    for _, room := range participantRooms {
        if _, exists := roomMap[room.ID]; !exists {
            count, _ := s.participantRepo.CountByRoomID(ctx, room.ID)
            roomMap[room.ID] = &model.RoomResponse{
                Room:             *room,
                ParticipantCount: count,
                IsOwner:          false,
            }
        }
    }
    
    // Convert map to slice
    rooms := make([]*model.RoomResponse, 0, len(roomMap))
    for _, room := range roomMap {
        rooms = append(rooms, room)
    }
    
    return rooms, nil
}

func (s *RoomService) GetRoomDetails(ctx context.Context, roomID, userID uuid.UUID) (*model.RoomResponse, error) {
    // 1. Get room
    room, err := s.roomRepo.GetByID(ctx, roomID)
    if err != nil {
        return nil, err
    }
    
    // 2. Check if user has access
    hasAccess, err := s.participantRepo.UserHasAccess(ctx, roomID, userID)
    if err != nil || !hasAccess {
        return nil, errors.New("access denied")
    }
    
    // 3. Get participant count
    count, err := s.participantRepo.CountByRoomID(ctx, roomID)
    if err != nil {
        return nil, err
    }
    
    return &model.RoomResponse{
        Room:             *room,
        ParticipantCount: count,
        IsOwner:          room.OwnerID == userID,
    }, nil
}

func (s *RoomService) DeleteRoom(ctx context.Context, roomID, userID uuid.UUID) error {
    // 1. Get room
    room, err := s.roomRepo.GetByID(ctx, roomID)
    if err != nil {
        return err
    }
    
    // 2. Check ownership
    if room.OwnerID != userID {
        return errors.New("only room owner can delete")
    }
    
    // 3. Delete from LiveKit
    if room.LiveKitRoomName != nil {
        s.livekitService.DeleteRoom(ctx, *room.LiveKitRoomName)
    }
    
    // 4. Delete from database (cascade will handle participants)
    return s.roomRepo.Delete(ctx, roomID)
}
```

#### service/participant_service.go

```go
package service

import (
    "context"
    "github.com/google/uuid"
    "yourapp/internal/model"
    "yourapp/internal/repository"
    "yourapp/internal/service/email"
)

type ParticipantService struct {
    participantRepo repository.ParticipantRepository
    roomRepo        repository.RoomRepository
    inviteRepo      repository.InviteRepository
    emailService    *email.EmailService
    livekitService  *LiveKitService
    frontendURL     string
}

func NewParticipantService(
    participantRepo repository.ParticipantRepository,
    roomRepo repository.RoomRepository,
    inviteRepo repository.InviteRepository,
    emailService *email.EmailService,
    livekitService *LiveKitService,
    frontendURL string,
) *ParticipantService {
    return &ParticipantService{
        participantRepo: participantRepo,
        roomRepo:        roomRepo,
        inviteRepo:      inviteRepo,
        emailService:    emailService,
        livekitService:  livekitService,
        frontendURL:     frontendURL,
    }
}

func (s *ParticipantService) AddParticipant(
    ctx context.Context,
    roomID uuid.UUID,
    inviterID uuid.UUID,
    req *model.AddParticipantRequest,
) (*model.ParticipantInviteResponse, error) {
    // 1. Get room
    room, err := s.roomRepo.GetByID(ctx, roomID)
    if err != nil {
        return nil, err
    }
    
    // 2. Check if participant already exists
    existing, _ := s.participantRepo.GetByRoomAndEmail(ctx, roomID, req.Email)
    if existing != nil {
        return nil, errors.New("participant already added")
    }
    
    // 3. Create participant
    participantID := uuid.New()
    participant := &model.RoomParticipant{
        ParticipantID: &participantID,
        RoomID:        roomID,
        Email:         req.Email,
        Name:          req.Name,
        Role:          "participant",
    }
    
    err = s.participantRepo.Create(ctx, participant)
    if err != nil {
        return nil, err
    }
    
    // 4. Create invite token
    inviteToken := uuid.New().String()
    expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days
    
    err = s.inviteRepo.Create(ctx, roomID, inviterID, req.Email, req.Name, inviteToken, expiresAt)
    if err != nil {
        return nil, err
    }
    
    // 5. Send invitation email
    inviteURL := s.frontendURL + "/room/" + roomID.String() + "/join?token=" + inviteToken
    err = s.emailService.SendRoomInviteEmail(ctx, req.Email, room.RoomName, inviteURL)
    if err != nil {
        // Log error but don't fail
    }
    
    return &model.ParticipantInviteResponse{
        ParticipantID: participantID,
        InviteToken:   inviteToken,
        InviteURL:     inviteURL,
    }, nil
}

func (s *ParticipantService) GetRoomParticipants(ctx context.Context, roomID uuid.UUID) ([]*model.RoomParticipant, error) {
    return s.participantRepo.GetByRoomID(ctx, roomID)
}

func (s *ParticipantService) RemoveParticipant(ctx context.Context, roomID, participantID uuid.UUID) error {
    return s.participantRepo.Delete(ctx, participantID)
}

func (s *ParticipantService) GenerateParticipantToken(
    ctx context.Context,
    roomID uuid.UUID,
    inviteToken string,
) (string, error) {
    // 1. Verify invite token
    invite, err := s.inviteRepo.GetByToken(ctx, inviteToken)
    if err != nil {
        return "", errors.New("invalid invite")
    }
    
    if invite.RoomID != roomID || time.Now().After(invite.ExpiresAt) {
        return "", errors.New("invalid or expired invite")
    }
    
    // 2. Get room
    room, err := s.roomRepo.GetByID(ctx, roomID)
    if err != nil {
        return "", err
    }
    
    // 3. Generate LiveKit token
    identity := invite.InviteeEmail
    token, err := s.livekitService.GenerateToken(
        identity,
        *room.LiveKitRoomName,
        true, // Can publish
        true, // Can subscribe
    )
    
    if err != nil {
        return "", err
    }
    
    // 4. Mark invite as accepted
    s.inviteRepo.MarkAsAccepted(ctx, invite.ID)
    
    return token, nil
}
```

### 4.7 Handler Layer (Pseudo Code)

#### handler/auth_handler.go

```go
package handler

import (
    "encoding/json"
    "net/http"
    "yourapp/internal/model"
    "yourapp/internal/service"
)

type AuthHandler struct {
    authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
    var req model.UserSignUpRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request")
        return
    }
    
    err := h.authService.SignUp(r.Context(), &req)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, err.Error())
        return
    }
    
    respondWithJSON(w, http.StatusCreated, map[string]string{
        "message": "User created successfully",
    })
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
    var req model.UserSignInRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request")
        return
    }
    
    authResp, err := h.authService.SignIn(r.Context(), &req)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, err.Error())
        return
    }
    
    respondWithJSON(w, http.StatusOK, authResp)
}

func (h *AuthHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
    var req model.PasswordResetRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request")
        return
    }
    
    err := h.authService.RequestPasswordReset(r.Context(), req.Email)
    if err != nil {
        // Always return success to prevent email enumeration
    }
    
    respondWithJSON(w, http.StatusOK, map[string]string{
        "message": "If the email exists, a reset link has been sent",
    })
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
    var req model.PasswordResetConfirm
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request")
        return
    }
    
    err := h.authService.ResetPassword(r.Context(), &req)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, err.Error())
        return
    }
    
    respondWithJSON(w, http.StatusOK, map[string]string{
        "message": "Password reset successfully",
    })
}
```

#### handler/room_handler.go

```go
package handler

import (
    "encoding/json"
    "net/http"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "yourapp/internal/model"
    "yourapp/internal/service"
)

type RoomHandler struct {
    roomService *service.RoomService
}

func NewRoomHandler(roomService *service.RoomService) *RoomHandler {
    return &RoomHandler{roomService: roomService}
}

func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
    userID := getUserIDFromContext(r.Context())
    
    var req model.CreateRoomRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request")
        return
    }
    
    room, err := h.roomService.CreateRoom(r.Context(), userID, &req)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    respondWithJSON(w, http.StatusCreated, room)
}

func (h *RoomHandler) GetUserRooms(w http.ResponseWriter, r *http.Request) {
    userID := getUserIDFromContext(r.Context())
    
    rooms, err := h.roomService.GetUserRooms(r.Context(), userID)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    respondWithJSON(w, http.StatusOK, rooms)
}

func (h *RoomHandler) GetRoomDetails(w http.ResponseWriter, r *http.Request) {
    userID := getUserIDFromContext(r.Context())
    roomID, err := uuid.Parse(mux.Vars(r)["roomId"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid room ID")
        return
    }
    
    room, err := h.roomService.GetRoomDetails(r.Context(), roomID, userID)
    if err != nil {
        respondWithError(w, http.StatusForbidden, err.Error())
        return
    }
    
    respondWithJSON(w, http.StatusOK, room)
}

func (h *RoomHandler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
    userID := getUserIDFromContext(r.Context())
    roomID, err := uuid.Parse(mux.Vars(r)["roomId"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid room ID")
        return
    }
    
    err = h.roomService.DeleteRoom(r.Context(), roomID, userID)
    if err != nil {
        respondWithError(w, http.StatusForbidden, err.Error())
        return
    }
    
    respondWithJSON(w, http.StatusOK, map[string]string{
        "message": "Room deleted successfully",
    })
}
```

### 4.8 Main Application Setup

#### cmd/api/main.go

```go
package main

import (
    "database/sql"
    "log"
    "net/http"
    "os"
    
    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
    
    "yourapp/internal/config"
    "yourapp/internal/handler"
    "yourapp/internal/middleware"
    "yourapp/internal/repository"
    "yourapp/internal/service"
    "yourapp/internal/service/email"
)

func main() {
    // Load configuration
    cfg := config.Load()
    
    // Initialize database
    db, err := sql.Open("postgres", cfg.DatabaseURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()
    
    // Initialize repositories
    userRepo := repository.NewUserRepository(db)
    roomRepo := repository.NewRoomRepository(db)
    participantRepo := repository.NewParticipantRepository(db)
    postRepo := repository.NewPostRepository(db)
    resetTokenRepo := repository.NewPasswordResetTokenRepository(db)
    inviteRepo := repository.NewInviteRepository(db)
    
    // Initialize email provider
    var emailProvider email.EmailProvider
    if cfg.EmailProvider == "sendgrid" {
        emailProvider = email.NewSendGridProvider(
            cfg.SendGridAPIKey,
            cfg.SendGridFromEmail,
            cfg.SendGridFromName,
        )
    } else {
        emailProvider = email.NewMailjetProvider(
            cfg.MailjetAPIKey,
            cfg.MailjetSecretKey,
            cfg.MailjetFromEmail,
            cfg.MailjetFromName,
        )
    }
    
    emailService := email.NewEmailService(
        emailProvider,
        cfg.FromEmail,
        cfg.FromName,
    )
    
    // Initialize LiveKit service
    livekitService := service.NewLiveKitService(
        cfg.LiveKitAPIKey,
        cfg.LiveKitAPISecret,
        cfg.LiveKitURL,
    )
    
    // Initialize services
    authService := service.NewAuthService(
        userRepo,
        resetTokenRepo,
        emailService,
        livekitService,
        cfg.JWTSecret,
        cfg.FrontendURL,
    )
    
    roomService := service.NewRoomService(
        roomRepo,
        participantRepo,
        livekitService,
    )
    
    participantService := service.NewParticipantService(
        participantRepo,
        roomRepo,
        inviteRepo,
        emailService,
        livekitService,
        cfg.FrontendURL,
    )
    
    postService := service.NewPostService(postRepo, roomRepo)
    
    // Initialize handlers
    authHandler := handler.NewAuthHandler(authService)
    roomHandler := handler.NewRoomHandler(roomService)
    participantHandler := handler.NewParticipantHandler(participantService)
    postHandler := handler.NewPostHandler(postService)
    
    // Setup router
    r := mux.NewRouter()
    
    // Middleware
    r.Use(middleware.CORSMiddleware(cfg.CORSAllowedOrigins))
    r.Use(middleware.LoggingMiddleware)
    
    // Public routes
    r.HandleFunc("/api/auth/signup", authHandler.SignUp).Methods("POST")
    r.HandleFunc("/api/auth/signin", authHandler.SignIn).Methods("POST")
    r.HandleFunc("/api/auth/reset-password", authHandler.RequestPasswordReset).Methods("POST")
    r.HandleFunc("/api/auth/reset-password/confirm", authHandler.ResetPassword).Methods("POST")
    
    // Protected routes
    api := r.PathPrefix("/api").Subrouter()
    api.Use(middleware.AuthMiddleware(cfg.JWTSecret))
    
    // Room routes
    api.HandleFunc("/rooms", roomHandler.CreateRoom).Methods("POST")
    api.HandleFunc("/rooms", roomHandler.GetUserRooms).Methods("GET")
    api.HandleFunc("/rooms/{roomId}", roomHandler.GetRoomDetails).Methods("GET")
    api.HandleFunc("/rooms/{roomId}", roomHandler.DeleteRoom).Methods("DELETE")
    
    // Participant routes
    api.HandleFunc("/rooms/{roomId}/participants", participantHandler.AddParticipant).Methods("POST")
    api.HandleFunc("/rooms/{roomId}/participants", participantHandler.GetParticipants).Methods("GET")
    api.HandleFunc("/rooms/{roomId}/participants/{participantId}", participantHandler.RemoveParticipant).Methods("DELETE")
    api.HandleFunc("/rooms/{roomId}/join", participantHandler.JoinRoom).Methods("POST")
    
    // Post routes
    api.HandleFunc("/rooms/{roomId}/posts", postHandler.CreatePost).Methods("POST")
    api.HandleFunc("/rooms/{roomId}/posts", postHandler.GetPosts).Methods("GET")
    api.HandleFunc("/rooms/{roomId}/posts/{postId}", postHandler.DeletePost).Methods("DELETE")
    
    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}
```

---

## 5. Frontend Implementation

### 5.1 API Service Layer

#### src/services/api/auth.service.ts

```typescript
import { apiClient } from '@/utils/api';

export interface SignUpRequest {
  username: string;
  email: string;
  name: string;
  password: string;
  confirm_password: string;
}

export interface SignInRequest {
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: User;
  livekit_token: string;
  expires_at: string;
}

export interface User {
  id: string;
  username: string;
  email: string;
  name: string;
  created_at: string;
}

export const authService = {
  async signUp(data: SignUpRequest): Promise<void> {
    await apiClient.post('/auth/signup', data);
  },

  async signIn(data: SignInRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/auth/signin', data);
    // Store token
    localStorage.setItem('auth_token', response.token);
    localStorage.setItem('livekit_token', response.livekit_token);
    return response;
  },

  async requestPasswordReset(email: string): Promise<void> {
    await apiClient.post('/auth/reset-password', { email });
  },

  async resetPassword(token: string, newPassword: string): Promise<void> {
    await apiClient.post('/auth/reset-password/confirm', {
      token,
      new_password: newPassword,
    });
  },

  logout() {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('livekit_token');
  },

  getToken(): string | null {
    return localStorage.getItem('auth_token');
  },

  getLiveKitToken(): string | null {
    return localStorage.getItem('livekit_token');
  },
};
```

#### src/services/api/room.service.ts

```typescript
import { apiClient } from '@/utils/api';

export interface Room {
  id: string;
  room_name: string;
  description?: string;
  owner_id: string;
  created_at: string;
  participant_count: number;
  is_owner: boolean;
}

export interface CreateRoomRequest {
  room_name: string;
  description?: string;
}

export const roomService = {
  async createRoom(data: CreateRoomRequest): Promise<Room> {
    return apiClient.post<Room>('/rooms', data);
  },

  async getUserRooms(): Promise<Room[]> {
    return apiClient.get<Room[]>('/rooms');
  },

  async getRoomDetails(roomId: string): Promise<Room> {
    return apiClient.get<Room>(`/rooms/${roomId}`);
  },

  async deleteRoom(roomId: string): Promise<void> {
    await apiClient.delete(`/rooms/${roomId}`);
  },
};
```

#### src/utils/api.ts

```typescript
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

class APIClient {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const token = localStorage.getItem('auth_token');
    
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${this.baseURL}${endpoint}`, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({}));
      throw new Error(error.message || 'Request failed');
    }

    return response.json();
  }

  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'GET' });
  }

  async post<T>(endpoint: string, data: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async put<T>(endpoint: string, data: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE' });
  }
}

export const apiClient = new APIClient(API_BASE_URL);
```

### 5.2 Key Components (Pseudo Code)

#### src/components/auth/SignUpForm.tsx

```typescript
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { authService } from '@/services/api/auth.service';

export default function SignUpForm() {
  const router = useRouter();
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    name: '',
    password: '',
    confirm_password: '',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    if (formData.password !== formData.confirm_password) {
      setError('Passwords do not match');
      setLoading(false);
      return;
    }

    try {
      await authService.signUp(formData);
      router.push('/auth/signin?registered=true');
    } catch (err: any) {
      setError(err.message || 'Sign up failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <input
        type="text"
        placeholder="Username"
        value={formData.username}
        onChange={(e) => setFormData({ ...formData, username: e.target.value })}
        required
      />
      <input
        type="email"
        placeholder="Email"
        value={formData.email}
        onChange={(e) => setFormData({ ...formData, email: e.target.value })}
        required
      />
      <input
        type="text"
        placeholder="Full Name"
        value={formData.name}
        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
        required
      />
      <input
        type="password"
        placeholder="Password"
        value={formData.password}
        onChange={(e) => setFormData({ ...formData, password: e.target.value })}
        required
      />
      <input
        type="password"
        placeholder="Confirm Password"
        value={formData.confirm_password}
        onChange={(e) => setFormData({ ...formData, confirm_password: e.target.value })}
        required
      />
      {error && <div className="text-red-500">{error}</div>}
      <button type="submit" disabled={loading}>
        {loading ? 'Signing up...' : 'Sign Up'}
      </button>
    </form>
  );
}
```

#### src/components/room/RoomList.tsx

```typescript
'use client';

import { useEffect, useState } from 'react';
import { roomService, Room } from '@/services/api/room.service';
import RoomCard from './RoomCard';

export default function RoomList() {
  const [rooms, setRooms] = useState<Room[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadRooms();
  }, []);

  const loadRooms = async () => {
    try {
      const data = await roomService.getUserRooms();
      setRooms(data);
    } catch (error) {
      console.error('Failed to load rooms:', error);
    } finally {
      setLoading(false);
    }
  };

  const filteredRooms = rooms.filter(room =>
    room.room_name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-4">
        <input
          type="text"
          placeholder="Search rooms..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="flex-1 px-4 py-2 border rounded"
        />
        <button
          onClick={() => router.push('/home/create-room')}
          className="px-4 py-2 bg-blue-600 text-white rounded"
        >
          + Add Room
        </button>
      </div>

      {loading ? (
        <div>Loading rooms...</div>
      ) : (
        <div className="grid gap-4">
          {filteredRooms.map(room => (
            <RoomCard key={room.id} room={room} />
          ))}
        </div>
      )}
    </div>
  );
}
```

#### src/components/room/RoomDetails.tsx

```typescript
'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import RoomHeader from './RoomHeader';
import PostsList from './PostsList';
import ParticipantsSidebar from './ParticipantsSidebar';

interface Props {
  roomId: string;
}

export default function RoomDetails({ roomId }: Props) {
  const router = useRouter();
  const [showParticipants, setShowParticipants] = useState(false);
  const [room, setRoom] = useState(null);

  useEffect(() => {
    loadRoomDetails();
  }, [roomId]);

  const loadRoomDetails = async () => {
    try {
      const data = await roomService.getRoomDetails(roomId);
      setRoom(data);
    } catch (error) {
      console.error('Failed to load room:', error);
    }
  };

  const handleJoinMeeting = () => {
    router.push(`/room/${roomId}/prep`);
  };

  return (
    <div className="flex h-screen">
      <div className="flex-1 flex flex-col">
        <RoomHeader
          room={room}
          onShowMembers={() => setShowParticipants(!showParticipants)}
          onJoinMeeting={handleJoinMeeting}
        />
        <div className="flex-1 overflow-auto">
          <PostsList roomId={roomId} />
        </div>
      </div>

      {showParticipants && (
        <ParticipantsSidebar
          roomId={roomId}
          onClose={() => setShowParticipants(false)}
        />
      )}
    </div>
  );
}
```

#### src/components/room/ParticipantsSidebar.tsx

```typescript
'use client';

import { useState, useEffect } from 'react';
import { participantService } from '@/services/api/participant.service';
import AddParticipantModal from './AddParticipantModal';

interface Props {
  roomId: string;
  onClose: () => void;
}

export default function ParticipantsSidebar({ roomId, onClose }: Props) {
  const [participants, setParticipants] = useState([]);
  const [showAddModal, setShowAddModal] = useState(false);

  useEffect(() => {
    loadParticipants();
  }, [roomId]);

  const loadParticipants = async () => {
    try {
      const data = await participantService.getParticipants(roomId);
      setParticipants(data);
    } catch (error) {
      console.error('Failed to load participants:', error);
    }
  };

  const handleAddParticipant = async (email: string, name: string) => {
    try {
      await participantService.addParticipant(roomId, { email, name });
      loadParticipants();
      setShowAddModal(false);
    } catch (error) {
      console.error('Failed to add participant:', error);
    }
  };

  return (
    <div className="w-80 bg-white border-l p-4">
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-lg font-semibold">Participants</h3>
        <button onClick={onClose}>×</button>
      </div>

      <button
        onClick={() => setShowAddModal(true)}
        className="w-full mb-4 px-4 py-2 bg-blue-600 text-white rounded"
      >
        + Add Participant
      </button>

      <div className="space-y-2">
        {participants.map(participant => (
          <div key={participant.id} className="p-3 border rounded">
            <div className="font-medium">{participant.name}</div>
            <div className="text-sm text-gray-600">{participant.email}</div>
            <div className="text-xs text-gray-500">{participant.role}</div>
          </div>
        ))}
      </div>

      {showAddModal && (
        <AddParticipantModal
          onAdd={handleAddParticipant}
          onClose={() => setShowAddModal(false)}
        />
      )}
    </div>
  );
}
```

#### src/components/meeting/PrepRoom.tsx

```typescript
'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useLocalParticipant, useMediaDevices } from '@livekit/components-react';
import BackgroundSelector from './BackgroundSelector';

interface Props {
  roomId: string;
}

export default function PrepRoom({ roomId }: Props) {
  const router = useRouter();
  const [videoEnabled, setVideoEnabled] = useState(true);
  const [audioEnabled, setAudioEnabled] = useState(true);
  const [selectedBackground, setSelectedBackground] = useState('none');
  const { devices } = useMediaDevices();

  const handleJoin = () => {
    router.push(`/room/${roomId}/meeting`);
  };

  return (
    <div className="h-screen flex flex-col items-center justify-center bg-gray-900">
      <div className="max-w-4xl w-full p-8">
        <h2 className="text-white text-2xl mb-8">Ready to join?</h2>

        <div className="bg-black rounded-lg mb-8 aspect-video relative">
          {/* Video preview */}
          <video
            autoPlay
            muted
            className="w-full h-full object-cover rounded-lg"
          />
          
          <div className="absolute bottom-4 left-1/2 transform -translate-x-1/2 flex gap-4">
            <button
              onClick={() => setVideoEnabled(!videoEnabled)}
              className="px-4 py-2 bg-gray-700 text-white rounded"
            >
              {videoEnabled ? '📹' : '📹❌'}
            </button>
            <button
              onClick={() => setAudioEnabled(!audioEnabled)}
              className="px-4 py-2 bg-gray-700 text-white rounded"
            >
              {audioEnabled ? '🎤' : '🎤❌'}
            </button>
          </div>
        </div>

        <BackgroundSelector
          onSelect={setSelectedBackground}
          selected={selectedBackground}
        />

        <div className="mt-8 flex justify-center gap-4">
          <button
            onClick={() => router.back()}
            className="px-6 py-3 bg-gray-700 text-white rounded"
          >
            Cancel
          </button>
          <button
            onClick={handleJoin}
            className="px-6 py-3 bg-blue-600 text-white rounded"
          >
            Join Now
          </button>
        </div>
      </div>
    </div>
  );
}
```

#### src/components/meeting/VideoConference.tsx

```typescript
'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import {
  LiveKitRoom,
  VideoConference as LKVideoConference,
  GridLayout,
  ParticipantTile,
  RoomAudioRenderer,
  ControlBar,
  useTracks,
} from '@livekit/components-react';
import '@livekit/components-styles';
import { Track } from 'livekit-client';

interface Props {
  roomId: string;
  token: string;
}

export default function VideoConference({ roomId, token }: Props) {
  const router = useRouter();
  const serverUrl = process.env.NEXT_PUBLIC_LIVEKIT_URL!;

  const handleDisconnect = () => {
    router.push(`/room/${roomId}`);
  };

  return (
    <div className="h-screen">
      <LiveKitRoom
        video={true}
        audio={true}
        token={token}
        serverUrl={serverUrl}
        connect={true}
        onDisconnected={handleDisconnect}
      >
        <LKVideoConference />
        <RoomAudioRenderer />
      </LiveKitRoom>
    </div>
  );
}

// Alternative custom implementation
export function CustomVideoConference({ roomId, token }: Props) {
  const router = useRouter();
  const serverUrl = process.env.NEXT_PUBLIC_LIVEKIT_URL!;
  const tracks = useTracks(
    [
      { source: Track.Source.Camera, withPlaceholder: true },
      { source: Track.Source.ScreenShare, withPlaceholder: false },
    ],
    { onlySubscribed: false },
  );

  return (
    <div className="h-screen flex flex-col bg-gray-900">
      <LiveKitRoom
        video={true}
        audio={true}
        token={token}
        serverUrl={serverUrl}
        connect={true}
        onDisconnected={() => router.push(`/room/${roomId}`)}
      >
        <div className="flex-1 p-4">
          <GridLayout tracks={tracks}>
            <ParticipantTile />
          </GridLayout>
        </div>
        
        <div className="p-4">
          <ControlBar />
        </div>
        
        <RoomAudioRenderer />
      </LiveKitRoom>
    </div>
  );
}
```

#### src/components/meeting/BackgroundSelector.tsx

```typescript
'use client';

import { useState } from 'react';

interface Props {
  onSelect: (background: string) => void;
  selected: string;
}

export default function BackgroundSelector({ onSelect, selected }: Props) {
  const backgrounds = [
    { id: 'none', label: 'No Background', preview: '/bg-none.jpg' },
    { id: 'blur', label: 'Blur', preview: '/bg-blur.jpg' },
    { id: 'office', label: 'Office', preview: '/bg-office.jpg' },
    { id: 'beach', label: 'Beach', preview: '/bg-beach.jpg' },
  ];

  return (
    <div className="space-y-4">
      <h3 className="text-white text-lg">Background Effects</h3>
      <div className="grid grid-cols-4 gap-4">
        {backgrounds.map(bg => (
          <button
            key={bg.id}
            onClick={() => onSelect(bg.id)}
            className={`relative rounded-lg overflow-hidden border-2 ${
              selected === bg.id ? 'border-blue-500' : 'border-gray-600'
            }`}
          >
            <div className="aspect-video bg-gray-800">
              {/* Background preview */}
            </div>
            <div className="p-2 text-sm text-white">{bg.label}</div>
          </button>
        ))}
      </div>
      
      <p className="text-sm text-gray-400">
        Note: Background effects require browser support for background blur/replacement.
        Consider using libraries like @mediapipe/selfie_segmentation or @tensorflow/tfjs
        for custom background effects.
      </p>
    </div>
  );
}
```

### 5.3 Page Routes

#### src/app/home/page.tsx

```typescript
'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Header from '@/components/layout/Header';
import Sidebar from '@/components/layout/Sidebar';
import RoomList from '@/components/room/RoomList';
import { authService } from '@/services/api/auth.service';

export default function HomePage() {
  const router = useRouter();

  useEffect(() => {
    const token = authService.getToken();
    if (!token) {
      router.push('/auth/signin');
    }
  }, []);

  return (
    <div className="flex h-screen">
      <Sidebar />
      <div className="flex-1 flex flex-col">
        <Header />
        <main className="flex-1 overflow-auto p-8">
          <RoomList />
        </main>
      </div>
    </div>
  );
}
```

#### src/app/room/[roomId]/page.tsx

```typescript
import RoomDetails from '@/components/room/RoomDetails';

export default function RoomPage({ params }: { params: { roomId: string } }) {
  return <RoomDetails roomId={params.roomId} />;
}
```

#### src/app/room/[roomId]/prep/page.tsx

```typescript
import PrepRoom from '@/components/meeting/PrepRoom';

export default function PrepPage({ params }: { params: { roomId: string } }) {
  return <PrepRoom roomId={params.roomId} />;
}
```

#### src/app/room/[roomId]/meeting/page.tsx

```typescript
'use client';

import { useEffect, useState } from 'react';
import { useSearchParams } from 'next/navigation';
import VideoConference from '@/components/meeting/VideoConference';
import { authService } from '@/services/api/auth.service';

export default function MeetingPage({ params }: { params: { roomId: string } }) {
  const searchParams = useSearchParams();
  const [token, setToken] = useState('');

  useEffect(() => {
    // Get token from query params or generate new one
    const inviteToken = searchParams.get('token');
    
    if (inviteToken) {
      // Participant joining via invite
      generateParticipantToken(inviteToken);
    } else {
      // User joining their own room
      const lkToken = authService.getLiveKitToken();
      setToken(lkToken || '');
    }
  }, []);

  const generateParticipantToken = async (inviteToken: string) => {
    try {
      const response = await participantService.generateToken(
        params.roomId,
        inviteToken
      );
      setToken(response.token);
    } catch (error) {
      console.error('Failed to generate token:', error);
    }
  };

  if (!token) {
    return <div>Loading...</div>;
  }

  return <VideoConference roomId={params.roomId} token={token} />;
}
```

#### src/app/auth/reset-password-form/page.tsx

```typescript
'use client';

import { useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { authService } from '@/services/api/auth.service';

export default function ResetPasswordFormPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const token = searchParams.get('token');
  
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (password !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    if (!token) {
      setError('Invalid reset link');
      return;
    }

    setLoading(true);
    try {
      await authService.resetPassword(token, password);
      router.push('/auth/signin?reset=success');
    } catch (err: any) {
      setError(err.message || 'Password reset failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="max-w-md w-full p-8 bg-white rounded-lg shadow">
        <h2 className="text-2xl font-bold mb-6">Reset Your Password</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <input
            type="password"
            placeholder="New Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            minLength={8}
          />
          <input
            type="password"
            placeholder="Confirm Password"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            required
            minLength={8}
          />
          {error && <div className="text-red-500">{error}</div>}
          <button type="submit" disabled={loading}>
            {loading ? 'Resetting...' : 'Reset Password'}
          </button>
        </form>
      </div>
    </div>
  );
}
```

---

## 6. API Endpoints Summary

### 6.1 Authentication Routes

```
POST   /api/auth/signup                 - Create new user account
POST   /api/auth/signin                 - Sign in and get JWT + LiveKit token
POST   /api/auth/reset-password         - Request password reset email
POST   /api/auth/reset-password/confirm - Confirm password reset with token
```

### 6.2 Room Routes (Protected)

```
POST   /api/rooms                       - Create new room
GET    /api/rooms                       - Get all user's rooms
GET    /api/rooms/{roomId}              - Get room details
DELETE /api/rooms/{roomId}              - Delete room
```

### 6.3 Participant Routes (Protected)

```
POST   /api/rooms/{roomId}/participants              - Add participant to room
GET    /api/rooms/{roomId}/participants              - Get all participants
DELETE /api/rooms/{roomId}/participants/{participantId} - Remove participant
POST   /api/rooms/{roomId}/join                      - Join room with invite token
```

### 6.4 Post Routes (Protected)

```
POST   /api/rooms/{roomId}/posts        - Create post in room
GET    /api/rooms/{roomId}/posts        - Get all posts in room
DELETE /api/rooms/{roomId}/posts/{postId} - Delete post
```

---

## 7. Environment Configuration

### 7.1 Backend (.env)

```env
# Already provided above in section 4.1
```

### 7.2 Frontend (.env.local)

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api
NEXT_PUBLIC_LIVEKIT_URL=wss://your-project.livekit.cloud
NEXT_PUBLIC_APP_URL=http://localhost:3000
```

---

## 8. Key Implementation Notes

### 8.1 LiveKit Integration

1. **Token Generation**: LiveKit tokens are generated server-side with appropriate permissions
2. **Room Management**: Rooms are created in both LiveKit and PostgreSQL
3. **Participant Identity**: Use email or user ID as LiveKit identity
4. **Token Validity**: Tokens expire after 24 hours by default

### 8.2 Security Considerations

1. **JWT Authentication**: Simple JWT without refresh tokens
2. **Password Hashing**: Use bcrypt with appropriate cost factor
3. **SQL Injection**: Use parameterized queries
4. **CORS**: Configure allowed origins properly
5. **Rate Limiting**: Consider implementing rate limiting for API endpoints
6. **Input Validation**: Validate all user inputs on both frontend and backend

### 8.3 Email Provider Flexibility

The email service uses an interface-based design allowing easy switching between providers:
```go
// Switch provider in main.go
if cfg.EmailProvider == "sendgrid" {
    emailProvider = email.NewSendGridProvider(...)
} else {
    emailProvider = email.NewMailjetProvider(...)
}
```

### 8.4 Database Migrations

Use a migration tool like `golang-migrate` or `goose`:

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path internal/database/migrations -database "postgres://..." up
```

### 8.5 Background Effects Library Recommendations

For implementing background effects in the PrepRoom:
- **@mediapipe/selfie_segmentation**: Google's solution for background segmentation
- **@tensorflow/tfjs**: For custom ML-based background effects
- **@livekit/track-processors**: LiveKit's built-in track processors

---

## 9. Development Workflow

### 9.1 Backend Development

```bash
# Install dependencies
go mod download

# Run migrations
make migrate-up

# Start server
go run cmd/api/main.go

# Or with hot reload
air
```

### 9.2 Frontend Development

```bash
# Install dependencies
npm install

# Run development server
npm run dev

# Build for production
npm run build
```

---

## 10. Future Enhancements

1. **Refresh Token Implementation**: Add refresh token rotation for better security
2. **Single Sign-On (SSO)**: Integrate OAuth providers (Google, GitHub, etc.)
3. **Phone Number Verification**: Add SMS verification
4. **Recording**: Implement meeting recording using LiveKit Egress
5. **Screen Sharing**: Enhanced screen sharing controls
6. **Chat**: Real-time chat during meetings using WebSocket
7. **Webhooks**: LiveKit webhook handling for events
8. **Analytics**: Track meeting duration, participant activity
9. **Notifications**: Real-time notifications for invites, mentions
10. **Mobile App**: React Native mobile application

---

## 11. Testing Strategy

### 11.1 Backend Testing

```go
// Example test structure
func TestAuthService_SignUp(t *testing.T) {
    // Setup
    mockRepo := &MockUserRepository{}
    service := NewAuthService(mockRepo, ...)
    
    // Test
    err := service.SignUp(ctx, &SignUpRequest{...})
    
    // Assert
    assert.NoError(t, err)
}
```

### 11.2 Frontend Testing

```typescript
// Example test with React Testing Library
import { render, screen } from '@testing-library/react';
import SignUpForm from '@/components/auth/SignUpForm';

test('renders sign up form', () => {
  render(<SignUpForm />);
  expect(screen.getByPlaceholderText(/username/i)).toBeInTheDocument();
});
```

---

## 12. Deployment Considerations

### 12.1 Backend Deployment
- Deploy to cloud platforms (AWS, GCP, Azure, DigitalOcean)
- Use Docker for containerization
- Set up CI/CD pipeline
- Configure environment variables securely
- Use managed PostgreSQL service

### 12.2 Frontend Deployment
- Deploy to Vercel (optimized for Next.js)
- Alternative: Netlify, AWS Amplify
- Configure environment variables
- Set up custom domain
- Enable HTTPS

### 12.3 Database
- Use managed PostgreSQL (AWS RDS, GCP Cloud SQL, etc.)
- Set up automated backups
- Configure read replicas for scaling
- Monitor query performance

---

This specification provides a complete foundation for building your LiveKit video conference application with clear separation of concerns, modular architecture, and scalability in mind.