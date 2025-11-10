'use client';

import { useMessages } from '@/hooks/useMessages';
import MessageList from '@/components/messages/MessageList';
import MessageInput from '@/components/messages/MessageInput';
import RoomDetails from '@/components/room/RoomDetails';
import ParticipantsSidebar from '@/components/room/ParticipantsSidebar';

export default function RoomPage() {
  const { room, messages, participants, isLoading, error, loadMoreMessages, addParticipant } = useMessages();

  return (
    <div className="flex h-[calc(100vh-4rem)]">
      <main className="flex-1 flex flex-col bg-gray-50 dark:bg-gray-900">
        {room && <RoomDetails room={room} />}

        <div className="flex-1 overflow-hidden">
          {error && (
            <div className="p-4 mx-4 mt-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
              <p className="text-red-800 dark:text-red-200">{error}</p>
            </div>
          )}
          
          {isLoading && messages.length === 0 ? (
            <div className="flex items-center justify-center h-full">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600" />
            </div>
          ) : (
            <MessageList messages={messages} onLoadMore={loadMoreMessages} />
          )}
        </div>

        <MessageInput />
      </main>
      <aside className="w-80 border-l border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-4">
        {room && (
          <ParticipantsSidebar
            roomId={room.id}
            participants={participants}
            isLoading={isLoading && participants.length === 0}
            onAddParticipant={addParticipant}
          />
        )}
      </aside>
    </div>
  );
}
