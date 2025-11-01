import { apiClient } from '@/utils/api';
import { Room, RoomApiResponse } from '@/types/room.types';

export interface CreateRoomRequest {
  room_name: string;
  description?: string;
}

export const roomService = {
  async createRoom(data: CreateRoomRequest): Promise<Room> {
    // Assuming createRoom returns a flattened Room object, adjust if necessary
    return apiClient.post<Room>('/app/rooms', data);
  },

  async getUserRooms(): Promise<Room[]> {
    const response = await apiClient.get<RoomApiResponse[]>('/app/rooms');
    return response.map(res => ({
      ...res.room,
      participant_count: res.participant_count,
      is_owner: res.is_owner,
    }));
  },

  async getRoomDetails(roomId: string): Promise<Room> {
    const res = await apiClient.get<RoomApiResponse>(`/app/rooms/${roomId}`);
    return {
      ...res.room,
      participant_count: res.participant_count,
      is_owner: res.is_owner,
    };
  },

  async deleteRoom(roomId: string): Promise<void> {
    await apiClient.delete(`/app/rooms/${roomId}`);
  },
};
