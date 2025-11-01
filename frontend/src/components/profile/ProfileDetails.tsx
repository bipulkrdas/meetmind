'use client';

import { useEffect, useState } from 'react';
import { authService, User } from '@/services/api/auth.service';

export default function ProfileDetails() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const userData = await authService.getMe();
        setUser(userData);
      } catch (err: any) {
        setError(err.message || 'Failed to fetch user data');
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
  }, []);

  if (loading) {
    return <div>Loading profile...</div>;
  }

  if (error) {
    return <div className="text-red-500">Error: {error}</div>;
  }

  if (!user) {
    return <div>No user data found.</div>;
  }

  return (
    <div className="space-y-4">
      <div>
        <label className="block text-sm font-medium text-gray-700">Username</label>
        <p className="mt-1 p-2 border rounded-md bg-gray-50">{user.username}</p>
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-700">Email</label>
        <p className="mt-1 p-2 border rounded-md bg-gray-50">{user.email}</p>
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-700">Full Name</label>
        <p className="mt-1 p-2 border rounded-md bg-gray-50">{user.name}</p>
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-700">Member Since</label>
        <p className="mt-1 p-2 border rounded-md bg-gray-50">{new Date(user.created_at).toLocaleDateString()}</p>
      </div>
    </div>
  );
}