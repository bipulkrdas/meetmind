export type MessageType = "user_message" | "meeting_transcript" | "participant_joined";

export interface S3Keys {
  json: string;
  text: string;
}

export interface HTTPSUrls {
  json_https_url: string;
  text_https_url: string;
}

export interface TranscriptData {
  bucket: string;
  region: string;
  s3_keys: S3Keys;
  https_urls: HTTPSUrls;
  session_start: string; // ISO 8601 string
  session_end: string;   // ISO 8601 string
}

export interface ExtraData {
  transcript?: TranscriptData;
  // Add other types of extra data here as needed
}

export interface Message {
  id: string;
  room_id: string;
  user_id: string;
  username: string;
  seq_no: number;
  content: string;
  message_type: MessageType; // Added
  extra_data?: ExtraData;    // Added
  created_at: string;
  updated_at: string;
  edited?: boolean;
  attachments?: Attachment[];
  deleted_at?: string;
  reactions?: Reaction[];
}

export interface Reaction {
  emoji: string;
  userIds: string[];
  count: number;
}

export interface Attachment {
  id: string;
  file_name: string;
  storage_url: string;
  file_type: string;
  file_size: number;
}