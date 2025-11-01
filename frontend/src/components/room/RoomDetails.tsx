'use client';

import { Room } from '@/types/room.types';
import { useRouter } from 'next/navigation';

interface RoomDetailsProps {
  room: Room;
}

export default function RoomDetails({ room }: RoomDetailsProps) {
  const router = useRouter();

  return (
    <div className="p-4 bg-white rounded-lg shadow-md">
      <h1 className="text-2xl font-bold">{room.room_name}</h1>
      <p className="text-gray-600">{room.description}</p>
      <div className="mt-4">
        <button 
          className="px-4 py-2 bg-green-600 text-white rounded"
          onClick={() => router.push(`/app/room/${room.id}/prep`)}
        >
          Join Meeting
        </button>
      </div>
    </div>
  );
}
