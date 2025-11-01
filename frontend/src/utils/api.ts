// const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

/**
 * Constructs the full API URL for a given endpoint.
 * In development, it returns a relative path (e.g., '/api/dashboard/balance') to be handled by the Next.js proxy.
 * In production, it constructs the full absolute URL using the NEXT_PUBLIC_API_URL environment variable.
 *
 * @param endpoint The API endpoint path, starting with a slash (e.g., '/dashboard/balance').
 * @returns The complete URL for the API call.
 */
export const getBaseUrl = (): string => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL;

  // Production: Use the full URL from environment variables.
  // The NEXT_PUBLIC_API_URL should be the base URL of the backend (e.g., https://api.example.com)
  if (process.env.NODE_ENV === 'production' && apiUrl) {
    return `${apiUrl}/api`;
  }

  // Development: Use a relative path for the Next.js proxy.
  return `/api`;
};

class APIClient {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const token = localStorage.getItem('auth_token');
    
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${this.baseURL}${endpoint}`, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({}));
      throw new Error(error.message || 'Request failed');
    }

    return response.json();
  }

  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'GET' });
  }

  async post<T>(endpoint: string, data: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async put<T>(endpoint: string, data: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE' });
  }
}

export const apiClient = new APIClient(getBaseUrl());
