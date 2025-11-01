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
