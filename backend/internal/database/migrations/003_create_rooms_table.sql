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
    last_message_seq INTEGER DEFAULT 0,
    last_message_at TIMESTAMP, 
    is_active BOOLEAN DEFAULT true
);

CREATE INDEX idx_rooms_owner_id ON rooms(owner_id);
CREATE INDEX idx_rooms_room_sid ON rooms(room_sid);
