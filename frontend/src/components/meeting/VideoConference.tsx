'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter, useSearchParams } from 'next/navigation';
import { LiveKitRoom, VideoConference as LiveKitVideoConference } from '@livekit/components-react';

export default function VideoConference() {
  const router = useRouter();
  const params = useParams();
  const searchParams = useSearchParams();
  const roomId = params.roomId as string;
  const [livekitToken, setLivekitToken] = useState('');

  useEffect(() => {
    const token = searchParams.get('livekit_token');
    if (token) {
      setLivekitToken(token);
    } else {
      // Handle case where token is not available
      router.push(`/app/uroom/${roomId}/prep`);
    }
  }, [roomId, router, searchParams]);

  return (
    <LiveKitRoom
      token={livekitToken}
      serverUrl={process.env.NEXT_PUBLIC_LIVEKIT_URL!}
      connect={true}
      video={true}
      audio={true}
      onDisconnected={() => router.push(`/app/uroom/${roomId}`)}
      data-lk-theme="default"
    >
      <LiveKitVideoConference />
    </LiveKitRoom>
  );
}
