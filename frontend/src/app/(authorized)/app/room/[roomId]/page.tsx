'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { roomService } from '@/services/api/room.service';
import { Room } from '@/types/room.types';
import RoomDetails from '@/components/room/RoomDetails';
import ParticipantsSidebar from '@/components/room/ParticipantsSidebar';
import PostsList from '@/components/room/PostsList';

export default function RoomPage() {
  const params = useParams();
  const router = useRouter();
  const roomId = params.roomId as string;
  const [room, setRoom] = useState<Room | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (roomId) {
      loadRoomDetails();
    }
  }, [roomId]);

  const loadRoomDetails = async () => {
    try {
      const data = await roomService.getRoomDetails(roomId);
      setRoom(data);
    } catch (error) {
      console.error('Failed to load room details:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div>Loading room...</div>;
  }

  if (!room) {
    return <div>Room not found</div>;
  }

  return (
    <div>
      <div className="mb-4">
        <button
          onClick={() => router.push('/app/home')}
          className="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-500"
        >
          <span className="sr-only">Back to home</span>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-6 w-6 text-gray-700 dark:text-gray-300"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            strokeWidth={2}
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M10 19l-7-7m0 0l7-7m-7 7h18"
            />
          </svg>
        </button>
      </div>
      <div className="flex gap-4">
        <div className="flex-1">
          <RoomDetails room={room} />
          <PostsList roomId={roomId} />
        </div>
       {/* <ParticipantsSidebar roomId={roomId} /> */}
      </div>
    </div>
  );
}
