'use client';

import { useEffect, useState } from 'react';
import { roomService } from '@/services/api/room.service';
import { Room } from '@/types/room.types';
import RoomCard from './RoomCard';
import { useRouter } from 'next/navigation';

export default function RoomList() {
  const router = useRouter();
  const [rooms, setRooms] = useState<Room[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadRooms();
  }, []);

  const loadRooms = async () => {
    try {
      const data = await roomService.getUserRooms();
      setRooms(data);
    } catch (error) {
      console.error('Failed to load rooms:', error);
    } finally {
      setLoading(false);
    }
  };

  const filteredRooms = searchQuery
    ? rooms.filter(room =>
        room.room_name.toLowerCase().includes(searchQuery.toLowerCase())
      )
    : rooms;

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-4">
        <input
          type="text"
          placeholder="Search rooms..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="flex-1 px-4 py-2 border rounded"
        />
        <button
          onClick={() => router.push('/app/home/create-room')}
          className="px-4 py-2 bg-blue-600 text-white rounded"
        >
          + Add Room
        </button>
      </div>

      {loading ? (
        <div>Loading rooms...</div>
      ) : (
        <div className="grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
          {filteredRooms.map(room => (
            <RoomCard key={room.id} room={room} />
          ))}
        </div>
      )}
    </div>
  );
}
