import { apiClient } from '@/utils/api';

export interface CreatePostRequest {
  message: string;
}

export const postService = {
  async createPost(roomId: string, data: CreatePostRequest): Promise<void> {
    await apiClient.post(`/app/rooms/${roomId}/posts`, data);
  },

  async getPosts(roomId: string): Promise<any[]> {
    return apiClient.get(`/app/rooms/${roomId}/posts`);
  },

  async deletePost(roomId: string, postId: string): Promise<void> {
    await apiClient.delete(`/app/rooms/${roomId}/posts/${postId}`);
  },
};
