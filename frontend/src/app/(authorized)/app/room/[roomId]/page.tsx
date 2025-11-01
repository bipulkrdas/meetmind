'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import { roomService } from '@/services/api/room.service';
import { Room } from '@/types/room.types';
import RoomDetails from '@/components/room/RoomDetails';
import ParticipantsSidebar from '@/components/room/ParticipantsSidebar';
import PostsList from '@/components/room/PostsList';

export default function RoomPage() {
  const params = useParams();
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
    <div className="flex gap-4">
      <div className="flex-1">
        <RoomDetails room={room} />
        <PostsList roomId={roomId} />
      </div>
      <ParticipantsSidebar roomId={roomId} />
    </div>
  );
}
