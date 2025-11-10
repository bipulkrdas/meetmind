import { create } from 'zustand';
import { Message } from '@/types/message.types';
import { Room } from '@/types/room.types';
import { Participant } from '@/types/participant.types';

interface MessageState {
  room: Room | null;
  messages: Message[];
  participants: Participant[];
  isLoading: boolean;
  error: string | null;
  setRoom: (room: Room) => void;
  setMessages: (messages: Message[]) => void;
  addMessage: (message: Message) => void;
  prependMessages: (messages: Message[]) => void;
  setParticipants: (participants: Participant[]) => void;
  addParticipant: (participant: Participant) => void;
  setLoading: (isLoading: boolean) => void;
  setError: (error: string | null) => void;
  reset: () => void;
}

export const useMessageStore = create<MessageState>((set) => ({
  room: null,
  messages: [],
  participants: [],
  isLoading: false,
  error: null,
  setRoom: (room) => set({ room }),
  setMessages: (messages) => set({ messages }),
  addMessage: (message) => set((state) => ({ messages: [...state.messages, message] })),
  prependMessages: (messages) => set((state) => ({ messages: [...messages, ...state.messages] })),
  setParticipants: (participants) => set({ participants }),
  addParticipant: (participant) => set((state) => ({ participants: [...state.participants, participant] })),
  setLoading: (isLoading) => set({ isLoading }),
  setError: (error) => set({ error }),
  reset: () => set({ room: null, messages: [], participants: [], isLoading: false, error: null }),
}));