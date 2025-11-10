export interface Participant {
  id: string;
  room_id: string;
  participant_id?: string;
  user_id?: string;
  email: string;
  name: string;
  livekit_identity?: string;
  role: string;
  joined_at?: string;
  created_at: string;
  last_viewed_at?: string;
  last_read_seq_no: number;
  is_active: boolean;
}
