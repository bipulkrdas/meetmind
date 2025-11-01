'use client';
import { useState, useEffect } from 'react';
import { useRouter, useParams, useSearchParams } from 'next/navigation';
import { PreJoin, LocalUserChoices } from '@livekit/components-react';
import { participantService } from '@/services/api/participant.service';

export default function PrepRoom() {
  const router = useRouter();
  const params = useParams();
  const searchParams = useSearchParams();
  const roomId = params.roomId as string;
 

  const [liveKitToken, setLivekitToken] = useState('')

  useEffect(() => {
    if (roomId) {
      const fetchToken = async () => {
        try {
          const response = await participantService.joinRoomInternal(roomId);
          setLivekitToken(response.livekit_token);
        } catch (error) {
          console.error("Failed to fetch livekit token", error);
          // Handle error, maybe redirect to an error page or show a message
        }
      };
      fetchToken();
    }
  }, [roomId]);

  const handleJoin = (choices: LocalUserChoices) => {
    // In a real app, you might want to save these choices to localStorage
    // so the user's name and device choices are remembered for next time.
    if (liveKitToken) {
      router.push(`/app/room/${roomId}/meeting?livekit_token=${liveKitToken}`);
    } else {
        // Handle case where livekit token is not yet available
        console.error("Livekit token not available");
    }
  };

  return (
    <div className="p-8 bg-white rounded-lg shadow-md" data-lk-theme="default">
      <h1 className="text-2xl font-bold mb-4">Ready to join?</h1>
      <PreJoin onSubmit={handleJoin} />
    </div>
  );
}
