'use client';

import { useEffect, useCallback } from 'react';
import { useParams } from 'next/navigation';
import { useMessageStore } from '@/stores/message.store';
import { messageService } from '@/services/api/message.service';
import { roomService } from '@/services/api/room.service';
import { participantService, AddParticipantRequest } from '@/services/api/participant.service';

export function useMessages() {
  const params = useParams();
  const roomId = params.roomId as string;
  
  const {
    room,
    messages,
    participants,
    isLoading,
    error,
    setRoom,
    setMessages,
    setParticipants,
    addParticipant: addParticipantToStore,
    setLoading,
    setError,
    reset,
    prependMessages,
  } = useMessageStore();

  useEffect(() => {
    const loadInitialData = async () => {
      if (!roomId) return;
      setLoading(true);
      setError(null);
      
      try {
        // Fetch data sequentially to avoid race conditions on the DB
        const roomData = await roomService.getRoomDetails(roomId);
        setRoom(roomData);

        const messageData = await messageService.fetchMessages(roomId, { limit: 20 });
        setMessages(messageData);

        const participantsData = await participantService.getParticipants(roomId);
        setParticipants(participantsData);

      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load room data');
      } finally {
        setLoading(false);
      }
    };

    loadInitialData();

    return () => {
      reset();
    };
  }, [roomId, setRoom, setMessages, setParticipants, setLoading, setError, reset]);

  const loadMoreMessages = useCallback(async () => {
    if (isLoading || messages.length === 0) return;

    const oldestMessageId = messages[0].id;
    setLoading(true);

    try {
      const olderMessages = await messageService.fetchMessages(roomId, {
        before: oldestMessageId,
        limit: 20,
      });
      
      prependMessages(olderMessages);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load more messages');
    } finally {
      setLoading(false);
    }
  }, [roomId, messages, isLoading, setLoading, setError, prependMessages]);

  const addParticipant = useCallback(async (participantData: AddParticipantRequest) => {
    if (!roomId) throw new Error("Room ID is not available");
    
    const newParticipant = await participantService.addParticipant(roomId, participantData);
    addParticipantToStore(newParticipant);
  }, [roomId, addParticipantToStore]);

  return {
    room,
    messages,
    participants,
    isLoading,
    error,
    loadMoreMessages,
    addParticipant,
  };
}