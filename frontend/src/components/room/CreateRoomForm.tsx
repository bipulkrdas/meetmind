'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { roomService } from '@/services/api/room.service';

export default function CreateRoomForm() {
  const router = useRouter();
  const [formData, setFormData] = useState({
    room_name: '',
    description: '',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const newRoom = await roomService.createRoom(formData);
      router.push(`/app/room/${newRoom.id}`);
    } catch (err: any) {
      setError(err.message || 'Failed to create room');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <input
        type="text"
        name="room_name"
        placeholder="Room Name"
        value={formData.room_name}
        onChange={handleChange}
        required
        className="w-full px-4 py-2 border rounded"
      />
      <textarea
        name="description"
        placeholder="Description"
        value={formData.description}
        onChange={handleChange}
        className="w-full px-4 py-2 border rounded"
      />
      {error && <div className="text-red-500">{error}</div>}
      <button type="submit" disabled={loading} className="w-full px-4 py-2 bg-blue-600 text-white rounded disabled:bg-gray-400">
        {loading ? 'Creating room...' : 'Create Room'}
      </button>
    </form>
  );
}
