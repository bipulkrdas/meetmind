import { apiClient } from '@/utils/api';
import { Participant } from '@/types/participant.types';

export interface AddParticipantRequest {
  email: string;
  name: string;
}

export const participantService = {
  async addParticipant(roomId: string, data: AddParticipantRequest): Promise<Participant> {
    return apiClient.post(`/app/rooms/${roomId}/participants`, data);
  },

  async getParticipants(roomId:string): Promise<Participant[]> {
    return apiClient.get(`/app/rooms/${roomId}/participants`);
  },

  async removeParticipant(roomId: string, participantId: string): Promise<void> {
    await apiClient.delete(`/app/rooms/${roomId}/participants/${participantId}`);
  },

  async joinRoomExternal(roomId: string, token: string): Promise<{ livekit_token: string }> {
    return apiClient.post(`/rooms/${roomId}/join_external?token=${token}`, {});
  },

  async joinRoomInternal(roomId: string): Promise<{ livekit_token: string }> {
    return apiClient.post(`/app/rooms/${roomId}/join_internal`, {});
  },

  async inviteParticipantsToJoinMeeting(roomId: string): Promise<void> {
    await apiClient.post(`/app/rooms/${roomId}/invite_participants_to_join_meeting`, {});
  },

  async generateMeetingUrl(roomId: string): Promise<{ meeting_url: string }> {
    return apiClient.post(`/app/rooms/${roomId}/generate_meeting_url`, {});
  },
};
