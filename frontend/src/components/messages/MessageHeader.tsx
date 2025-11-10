'use client';

import { Room } from '@/types/room.types';

interface MessageHeaderProps {
  room: Room | null;
}

export default function MessageHeader({ room }: MessageHeaderProps) {
  return (
    <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
      <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">
        {room?.room_name || 'Loading...'}
      </h2>
    </div>
  );
}