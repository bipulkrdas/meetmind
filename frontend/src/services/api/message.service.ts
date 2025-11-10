import { Message } from '@/types/message.types';
import { apiClient, getBaseUrl } from '@/utils/api';

// Define the structure for a single transcript item
interface TranscriptItemContent {
  type: string;
  text: string;
}

interface TranscriptItem {
  timestamp: string;
  role: string;
  interrupted: boolean;
  content: TranscriptItemContent[];
  speaker_identity: string;
  speaker_name: string;
}

// Define the structure for the full transcript data
export interface FullTranscriptData { // Export this
  room_name: string;
  session_start: string;
  session_end: string;
  items: TranscriptItem[];
}

interface FetchMessagesOptions {
  limit?: number;
  before?: string;
  after?: string;
}

async function fetchMessages(
  roomId: string,
  options: FetchMessagesOptions = {}
): Promise<Message[]> {
  const { limit = 20, before, after } = options;

  const params = {
    limit,
    ...(before && { before }),
    ...(after && { after }),
  };

  return apiClient.get(`/app/rooms/${roomId}/messages`, { params });
}

async function sendMessage(
  roomId: string,
  data: { content: string; attachment_ids?: string[] }
): Promise<Message> {
  return apiClient.post(`/app/rooms/${roomId}/messages`, data);
}

async function fetchTranscript(
  roomId: string,
  messageId: string,
  s3KeyPath: string
): Promise<FullTranscriptData> {
  return apiClient.get(`/app/rooms/${roomId}/transcript/${messageId}/${s3KeyPath}`);
}

async function uploadFile(
  roomId: string,
  file: File,
  onProgress?: (progress: number) => void
): Promise<string> {
  return new Promise((resolve, reject) => {
    const formData = new FormData();
    formData.append('file', file);

    const xhr = new XMLHttpRequest();

    xhr.upload.addEventListener('progress', (e) => {
      if (e.lengthComputable && onProgress) {
        const progress = Math.round((e.loaded / e.total) * 100);
        onProgress(progress);
      }
    });

    xhr.addEventListener('load', () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        try {
          const response = JSON.parse(xhr.responseText);
          resolve(response.fileId);
        } catch (error) {
          reject(new Error('Invalid response from server'));
        }
      } else {
        reject(new Error(`Upload failed with status ${xhr.status}`));
      }
    });

    xhr.addEventListener('error', () => {
      reject(new Error('Network error during upload'));
    });

    xhr.addEventListener('abort', () => {
      reject(new Error('Upload aborted'));
    });

    xhr.open('POST', `${getBaseUrl()}/app/rooms/${roomId}/attachments`);
    xhr.setRequestHeader('Authorization', `Bearer ${getToken()}`);
    xhr.send(formData);
  });
}

export const messageService = {
    fetchMessages,
    sendMessage,
    uploadFile,
    fetchTranscript, // Export the new function
};

function getToken() {
 // throw new Error('Function not implemented.');
 return localStorage.getItem('auth_token');
}