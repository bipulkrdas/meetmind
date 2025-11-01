import { apiClient } from '@/utils/api';
import { User } from '@/types/user.types';

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

export const authService = {
  async signUp(data: SignUpRequest): Promise<void> {
    await apiClient.post('/auth/signup', data);
  },

  async signIn(data: SignInRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/auth/signin', data);
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

  async getMe(): Promise<User> {
    return apiClient.get<User>('/auth/me');
  },
};
