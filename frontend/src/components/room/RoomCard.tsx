'use client';

import { Room } from '@/types/room.types';
import { useRouter } from 'next/navigation';

interface RoomCardProps {
  room: Room;
}

export default function RoomCard({ room }: RoomCardProps) {
  const router = useRouter();

  return (
    <div 
      className="p-4 bg-white rounded-lg shadow-md cursor-pointer hover:shadow-lg transition-shadow"
      onClick={() => router.push(`/app/room/${room.id}`)}
    >
      <h3 className="text-lg font-bold">{room.room_name}</h3>
      <p className="text-gray-600">{room.description}</p>
      <div className="flex justify-between items-center mt-4">
        <span className="text-sm text-gray-500">{room.participant_count} participants</span>
        {room.is_owner && <span className="text-sm font-semibold text-blue-600">Owner</span>}
      </div>
    </div>
  );
}
