'use client';
import { useState, useEffect } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { PreJoin, LocalUserChoices } from '@livekit/components-react';
import { participantService } from '@/services/api/participant.service';

export default function PrepRoom() {
  const router = useRouter();
  const params = useParams();
  const roomId = params.roomId as string;

  const [liveKitToken, setLivekitToken] = useState('');

  useEffect(() => {
    if (roomId) {
      const fetchToken = async () => {
        try {
          const response = await participantService.joinRoomInternal(roomId);
          setLivekitToken(response.livekit_token);
        } catch (error) {
          console.error('Failed to fetch livekit token', error);
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
      router.push(`/app/uroom/${roomId}/meeting?livekit_token=${liveKitToken}`);
    } else {
      // Handle case where livekit token is not yet available
      console.error('Livekit token not available');
    }
  };

  const handleCancel = () => {
    router.push(`/app/uroom/${roomId}`);
  };

  return (
    <div className="p-8 bg-white rounded-lg shadow-md" data-lk-theme="default">
      <div className="mb-4">
        <button onClick={handleCancel} className="text-gray-600 hover:text-gray-800">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
        </button>
      </div>
      <h1 className="text-2xl font-bold mb-4 text-center">Ready to join?</h1>
      <PreJoin onSubmit={handleJoin} />
      <div className="mt-4 flex justify-center">
        <button
          onClick={handleCancel}
          className="bg-gray-200 text-gray-800 px-4 py-2 rounded-md hover:bg-gray-300"
        >
          Cancel
        </button>
      </div>
    </div>
  );
}
