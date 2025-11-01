export interface Room {
  id: string;
  room_name: string;
  room_sid: string;
  description?: string;
  owner_id: string;
  livekit_room_name: string;
  metadata: any | null;
  created_at: string;
  updated_at: string;
  is_active: boolean;
  participant_count: number;
  is_owner: boolean;
}

// Represents the nested room object from the backend
interface BackendRoom {
  id: string;
  room_name: string;
  room_sid: string;
  description?: string;
  owner_id: string;
  livekit_room_name: string;
  metadata: any | null;
  created_at: string;
  updated_at: string;
  is_active: boolean;
}

// Represents the actual structure of the API response for a single room
export interface RoomApiResponse {
  room: BackendRoom;
  participant_count: number;
  is_owner: boolean;
}
