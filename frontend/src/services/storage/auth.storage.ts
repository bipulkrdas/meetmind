export const authStorage = {
  getToken(): string | null {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('auth_token');
    }
    return null;
  },

  getLiveKitToken(): string | null {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('livekit_token');
    }
    return null;
  },

  setTokens(authToken: string, livekitToken: string) {
    if (typeof window !== 'undefined') {
      localStorage.setItem('auth_token', authToken);
      localStorage.setItem('livekit_token', livekitToken);
    }
  },

  clearTokens() {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('auth_token');
      localStorage.removeItem('livekit_token');
    }
  },
};
