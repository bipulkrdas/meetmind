export const livekitService = {
  getLiveKitToken(): string | null {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('livekit_token');
    }
    return null;
  },
};
