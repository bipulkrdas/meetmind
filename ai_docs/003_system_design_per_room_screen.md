# Messages Component System Design Document

## 1. Overview

This document outlines the architecture for a real-time messaging interface in Next.js, similar to Slack/Teams rooms. The system uses WebSocket for real-time updates and Zustand for lightweight state management.

---

## 2. Technology Stack

### Core Technologies
- **Framework**: Next.js 14+ (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: Zustand (lightweight, performant, modern alternative to Redux)
- **Real-time Communication**: WebSocket (connection established at Home screen level)

### Why Zustand?
- **Performance**: No unnecessary re-renders, selector-based subscriptions
- **Simplicity**: No boilerplate, minimal API surface
- **Modern**: Hooks-based, TypeScript-first
- **Bundle Size**: ~1KB (vs Redux ~20KB)
- **Learning Curve**: Extremely gentle compared to Redux

---

## 3. Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                    Home Screen (/)                       │
│  - Establishes WebSocket connection                     │
│  - Displays room list                                │
│  - Maintains connection context                         │
└──────────────────────┬──────────────────────────────────┘
                       │
                       │ Navigate to room
                       ▼
┌─────────────────────────────────────────────────────────┐
│         Messages Screen (/:roomId/messages)          │
│                                                          │
│  ┌────────────────────────────────────────────────┐   │
│  │  useEffect Hook (Component Mount)               │   │
│  │  1. Fetch last 20 messages (REST API)          │   │
│  │  2. Join room room (WebSocket send)         │   │
│  │  3. Subscribe to WebSocket messages            │   │
│  └────────────────────────────────────────────────┘   │
│                                                          │
│  ┌────────────────────────────────────────────────┐   │
│  │  Zustand Store (Channel-Specific)              │   │
│  │  - Messages array                               │   │
│  │  - Loading states                               │   │
│  │  - User typing indicators                       │   │
│  │  - Actions: addMessage, updateMessage, etc.    │   │
│  └────────────────────────────────────────────────┘   │
│                                                          │
│  ┌────────────────────────────────────────────────┐   │
│  │  UI Components                                  │   │
│  │  - MessageList (virtualized)                   │   │
│  │  - MessageItem                                  │   │
│  │  - MessageInput                                 │   │
│  │  - TypingIndicator                              │   │
│  └────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

---

## 4. File Structure

```
app/myrooms
├── [roomId]/
│   └── messages/
│       ├── page.tsx                    # Main Messages route component
│       ├── components/
│       │   ├── MessageList.tsx         # Virtualized message list
│       │   ├── MessageItem.tsx         # Individual message bubble
│       │   ├── MessageInput.tsx        # Input field with send button
│       │   ├── TypingIndicator.tsx     # Shows who's typing
│       │   └── MessageHeader.tsx       # Channel info header
│       ├── hooks/
│       │   ├── useMessages.ts          # Main hook for message logic
│       │   └── useMessageWebSocket.ts  # WebSocket subscription hook
│       └── store/
│           └── messageStore.ts         # Zustand store definition
│
├
│
├── contexts/
│   └── WebSocketContext.tsx            # Global WebSocket provider
│
├── lib/
│   ├── api/
│   │   └── messages.ts                 # API client for REST endpoints
│   └── websocket/
│       └── client.ts                   # WebSocket client wrapper
│
└── types/
    ├── message.ts                      # Message type definitions
    └── websocket.ts                    # WebSocket event types
```

---

## 5. Type Definitions

### types/message.ts
```typescript
export interface Message {
  id: string;
  roomId: string;
  userId: string;
  username: string;
  content: string;
  created_at: Date;
  updated_at: Date
  edited?: boolean;
  attachments?: Attachment[];
  is_deleted?: boolean;
}

export interface Reaction {
  emoji: string;
  userIds: string[];
  count: number;
}

export interface Attachment {
  id: string;
  name: string;
  url: string;
  type: string;
  size: number;
}

export interface TypingUser {
  userId: string;
  username: string;
}
```

### types/websocket.ts
```typescript
export enum WebSocketEventType {
  JOIN_ROOM = 'join_room',
  LEAVE_ROOM = 'leave_room',
  NEW_MESSAGE = 'new_message',
  MESSAGE_UPDATED = 'message_updated',
  MESSAGE_DELETED = 'message_deleted',
  USER_TYPING = 'user_typing',
  USER_STOPPED_TYPING = 'user_stopped_typing',
}

export interface WebSocketEvent {
  type: WebSocketEventType;
  payload: any;
  roomId: string;
  timestamp: number;
}

export interface JoinChannelPayload {
  roomId: string;
  userId: string;
}

export interface NewMessagePayload {
  message: Message;
}

export interface TypingPayload {
  userId: string;
  username: string;
  roomId: string;
}
```

---

## 6. State Management with Zustand

### app/[roomId]/messages/store/messageStore.ts

```typescript
import { create } from 'zustand';
import { Message, TypingUser } from '@/types/message';

interface MessageState {
  // State
  messages: Message[];
  isLoading: boolean;
  error: string | null;
  typingUsers: TypingUser[];
  hasMore: boolean;
  
  // Actions
  setMessages: (messages: Message[]) => void;
  addMessage: (message: Message) => void;
  updateMessage: (messageId: string, updates: Partial<Message>) => void;
  deleteMessage: (messageId: string) => void;
  prependMessages: (messages: Message[]) => void;
  setLoading: (isLoading: boolean) => void;
  setError: (error: string | null) => void;
  addTypingUser: (user: TypingUser) => void;
  removeTypingUser: (userId: string) => void;
  reset: () => void;
}

export const useMessageStore = create<MessageState>((set) => ({
  // Initial state
  messages: [],
  isLoading: false,
  error: null,
  typingUsers: [],
  hasMore: true,

  // Actions
  setMessages: (messages) => set({ messages }),
  
  addMessage: (message) =>
    set((state) => ({
      messages: [...state.messages, message],
    })),
  
  updateMessage: (messageId, updates) =>
    set((state) => ({
      messages: state.messages.map((msg) =>
        msg.id === messageId ? { ...msg, ...updates } : msg
      ),
    })),
  
  deleteMessage: (messageId) =>
    set((state) => ({
      messages: state.messages.filter((msg) => msg.id !== messageId),
    })),
  
  prependMessages: (messages) =>
    set((state) => ({
      messages: [...messages, ...state.messages],
      hasMore: messages.length > 0,
    })),
  
  setLoading: (isLoading) => set({ isLoading }),
  
  setError: (error) => set({ error }),
  
  addTypingUser: (user) =>
    set((state) => {
      if (state.typingUsers.some((u) => u.userId === user.userId)) {
        return state;
      }
      return { typingUsers: [...state.typingUsers, user] };
    }),
  
  removeTypingUser: (userId) =>
    set((state) => ({
      typingUsers: state.typingUsers.filter((u) => u.userId !== userId),
    })),
  
  reset: () =>
    set({
      messages: [],
      isLoading: false,
      error: null,
      typingUsers: [],
      hasMore: true,
    }),
}));
```

---

## 7. WebSocket Context (Global)

### app/contexts/WebSocketContext.tsx

```typescript
'use client';

import React, { createContext, useContext, useEffect, useRef, useState } from 'react';

interface WebSocketContextType {
  socket: WebSocket | null;
  isConnected: boolean;
  send: (data: any) => void;
  subscribe: (roomId: string, handler: (event: any) => void) => () => void;
}

const WebSocketContext = createContext<WebSocketContextType | null>(null);

export function WebSocketProvider({ children }: { children: React.ReactNode }) {
  const socketRef = useRef<WebSocket | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const handlersRef = useRef<Map<string, Set<(event: any) => void>>>(new Map());

  useEffect(() => {
    // Initialize WebSocket connection
    const ws = new WebSocket(process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:3001');
    
    ws.onopen = () => {
      console.log('WebSocket connected');
      setIsConnected(true);
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      const { roomId } = data;
      
      // Notify all subscribers for this room
      const handlers = handlersRef.current.get(roomId);
      if (handlers) {
        handlers.forEach((handler) => handler(data));
      }
    };

    ws.onclose = () => {
      console.log('WebSocket disconnected');
      setIsConnected(false);
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    socketRef.current = ws;

    return () => {
      ws.close();
    };
  }, []);

  const send = (data: any) => {
    if (socketRef.current?.readyState === WebSocket.OPEN) {
      socketRef.current.send(JSON.stringify(data));
    }
  };

  const subscribe = (roomId: string, handler: (event: any) => void) => {
    if (!handlersRef.current.has(roomId)) {
      handlersRef.current.set(roomId, new Set());
    }
    handlersRef.current.get(roomId)!.add(handler);

    // Return unsubscribe function
    return () => {
      const handlers = handlersRef.current.get(roomId);
      if (handlers) {
        handlers.delete(handler);
        if (handlers.size === 0) {
          handlersRef.current.delete(roomId);
        }
      }
    };
  };

  return (
    <WebSocketContext.Provider value={{ socket: socketRef.current, isConnected, send, subscribe }}>
      {children}
    </WebSocketContext.Provider>
  );
}

export const useWebSocket = () => {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error('useWebSocket must be used within WebSocketProvider');
  }
  return context;
};
```

---

## 8. Custom Hooks

### app/[roomId]/messages/hooks/useMessages.ts

```typescript
import { useEffect, useCallback } from 'react';
import { useParams } from 'next/navigation';
import { useMessageStore } from '../store/messageStore';
import { fetchMessages } from '@/lib/api/messages';
import { useMessageWebSocket } from './useMessageWebSocket';

export function useMessages() {
  const params = useParams();
  const roomId = params.roomId as string;
  
  const {
    messages,
    isLoading,
    error,
    setMessages,
    setLoading,
    setError,
    reset,
  } = useMessageStore();

  // Setup WebSocket subscription
  useMessageWebSocket(roomId);

  // Fetch initial messages on mount
  useEffect(() => {
    const loadMessages = async () => {
      setLoading(true);
      setError(null);
      
      try {
        const data = await fetchMessages(roomId, { limit: 20 });
        setMessages(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load messages');
      } finally {
        setLoading(false);
      }
    };

    loadMessages();

    // Cleanup: reset store when leaving room
    return () => {
      reset();
    };
  }, [roomId, setMessages, setLoading, setError, reset]);

  const loadMoreMessages = useCallback(async () => {
    if (isLoading || messages.length === 0) return;

    const oldestMessageId = messages[0].id;
    setLoading(true);

    try {
      const olderMessages = await fetchMessages(roomId, {
        before: oldestMessageId,
        limit: 20,
      });
      
      useMessageStore.getState().prependMessages(olderMessages);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load more messages');
    } finally {
      setLoading(false);
    }
  }, [roomId, messages, isLoading, setLoading, setError]);

  return {
    messages,
    isLoading,
    error,
    loadMoreMessages,
  };
}
```

### app/[roomId]/messages/hooks/useMessageWebSocket.ts

```typescript
import { useEffect } from 'react';
import { useWebSocket } from '@/app/contexts/WebSocketContext';
import { useMessageStore } from '../store/messageStore';
import { WebSocketEventType } from '@/types/websocket';

export function useMessageWebSocket(roomId: string) {
  const { send, subscribe, isConnected } = useWebSocket();
  const { addMessage, updateMessage, deleteMessage, addTypingUser, removeTypingUser } = useMessageStore();

  useEffect(() => {
    if (!isConnected) return;

    // Join room room
    send({
      type: WebSocketEventType.JOIN_ROOM,
      payload: { roomId },
    });

    // Subscribe to room events
    const unsubscribe = subscribe(roomId, (event) => {
      switch (event.type) {
        case WebSocketEventType.NEW_MESSAGE:
          addMessage(event.payload.message);
          break;
        
        case WebSocketEventType.MESSAGE_UPDATED:
          updateMessage(event.payload.messageId, event.payload.updates);
          break;
        
        case WebSocketEventType.MESSAGE_DELETED:
          deleteMessage(event.payload.messageId);
          break;
        
        case WebSocketEventType.USER_TYPING:
          addTypingUser({
            userId: event.payload.userId,
            username: event.payload.username,
          });
          break;
        
        case WebSocketEventType.USER_STOPPED_TYPING:
          removeTypingUser(event.payload.userId);
          break;
      }
    });

    // Cleanup: leave room room
    return () => {
      send({
        type: WebSocketEventType.LEAVE_ROOM,
        payload: { roomId },
      });
      unsubscribe();
    };
  }, [roomId, isConnected, send, subscribe, addMessage, updateMessage, deleteMessage, addTypingUser, removeTypingUser]);
}
```

---

## 9. Main Messages Page Component

### app/[roomId]/messages/page.tsx

```typescript
'use client';

import { useMessages } from './hooks/useMessages';
import MessageList from './components/MessageList';
import MessageInput from './components/MessageInput';
import MessageHeader from './components/MessageHeader';
import TypingIndicator from './components/TypingIndicator';
import { useMessageStore } from './store/messageStore';

export default function MessagesPage() {
  const { messages, isLoading, error, loadMoreMessages } = useMessages();
  const typingUsers = useMessageStore((state) => state.typingUsers);

  return (
    <div className="flex flex-col h-screen bg-gray-50 dark:bg-gray-900">
      {/* Header */}
      <MessageHeader />

      {/* Messages List */}
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

      {/* Typing Indicator */}
      {typingUsers.length > 0 && <TypingIndicator users={typingUsers} />}

      {/* Message Input */}
      <MessageInput />
    </div>
  );
}
```

---

## 10. UI Components

### MessageList Component (with Virtualization)

```typescript
// myrooms/[roomId]/messages/components/MessageList.tsx
'use client';

import { useEffect, useRef } from 'react';
import MessageItem from './MessageItem';
import { Message } from '@/types/message';

interface MessageListProps {
  messages: Message[];
  onLoadMore: () => void;
}

export default function MessageList({ messages, onLoadMore }: MessageListProps) {
  const scrollRef = useRef<HTMLDivElement>(null);
  const bottomRef = useRef<HTMLDivElement>(null);

  // Auto-scroll to bottom on new messages
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  // Infinite scroll: load more when scrolled to top
  const handleScroll = () => {
    if (scrollRef.current) {
      const { scrollTop } = scrollRef.current;
      if (scrollTop === 0) {
        onLoadMore();
      }
    }
  };

  return (
    <div
      ref={scrollRef}
      onScroll={handleScroll}
      className="h-full overflow-y-auto px-4 py-6 space-y-4"
    >
      {messages.map((message) => (
        <MessageItem key={message.id} message={message} />
      ))}
      <div ref={bottomRef} />
    </div>
  );
}
```

### MessageItem Component

```typescript
// app/myrooms/[roomId]/messages/components/MessageItem.tsx
import { Message } from '@/types/message';
import { formatDistanceToNow } from 'date-fns';
import { FileText } from 'lucide-react';

interface MessageItemProps {
  message: Message;
}

export default function MessageItem({ message }: MessageItemProps) {
  return (
    <div className="flex gap-3 group hover:bg-gray-100 dark:hover:bg-gray-800 px-4 py-2 rounded-lg transition-colors">
      {/* Avatar */}
      <div className="flex-shrink-0">
        <div className="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white font-semibold">
          {message.username.charAt(0).toUpperCase()}
        </div>
      </div>

      {/* Message Content */}
      <div className="flex-1 min-w-0">
        <div className="flex items-baseline gap-2">
          <span className="font-semibold text-gray-900 dark:text-gray-100">
            {message.username}
          </span>
          <span className="text-xs text-gray-500 dark:text-gray-400">
            {formatDistanceToNow(new Date(message.timestamp), { addSuffix: true })}
          </span>
          {message.edited && (
            <span className="text-xs text-gray-400 dark:text-gray-500">(edited)</span>
          )}
        </div>
        <p className="text-gray-800 dark:text-gray-200 mt-1 break-words">
          {message.content}
        </p>
        
        {/* Attachments */}
        {message.attachments && message.attachments.length > 0 && (
          <div className="mt-2 space-y-2">
            {message.attachments.map((attachment) => (
              <div
                key={attachment.id}
                className="flex items-center gap-2 p-2 bg-gray-100 dark:bg-gray-700 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors cursor-pointer"
                onClick={() => window.open(attachment.url, '_blank')}
              >
                {attachment.type.startsWith('image/') ? (
                  <img
                    src={attachment.url}
                    alt={attachment.name}
                    className="w-full max-w-sm rounded-lg"
                  />
                ) : (
                  <>
                    <FileText className="w-5 h-5 text-blue-600 dark:text-blue-400 flex-shrink-0" />
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">
                        {attachment.name}
                      </p>
                      <p className="text-xs text-gray-500 dark:text-gray-400">
                        {(attachment.size / 1024).toFixed(1)} KB
                      </p>
                    </div>
                  </>
                )}
              </div>
            ))}
          </div>
        )}
        
        {/* Reactions */}
        {message.reactions && message.reactions.length > 0 && (
          <div className="flex gap-2 mt-2">
            {message.reactions.map((reaction, idx) => (
              <button
                key={idx}
                className="flex items-center gap-1 px-2 py-1 bg-gray-200 dark:bg-gray-700 rounded-full text-sm hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
              >
                <span>{reaction.emoji}</span>
                <span className="text-xs text-gray-600 dark:text-gray-300">
                  {reaction.count}
                </span>
              </button>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
```

### MessageInput Component (Enhanced with File Upload)

```typescript
// app/myrooms/[roomId]/messages/components/MessageInput.tsx
'use client';

import { useState, useRef } from 'react';
import { useParams } from 'next/navigation';
import { useWebSocket } from '@/app/contexts/WebSocketContext';
import { sendMessage, uploadFiles } from '@/lib/api/messages';
import { WebSocketEventType } from '@/types/websocket';
import { Paperclip, X, Send, Image, FileText, Film } from 'lucide-react';

interface AttachedFile {
  id: string;
  file: File;
  preview?: string;
  uploading: boolean;
  progress: number;
  error?: string;
}

const MAX_FILE_SIZE = 10 * 1024 * 1024; // 10MB
const ALLOWED_FILE_TYPES = {
  images: ['image/jpeg', 'image/png', 'image/gif', 'image/webp'],
  videos: ['video/mp4', 'video/webm', 'video/quicktime'],
  documents: [
    'application/pdf',
    'application/msword',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    'application/vnd.ms-excel',
    'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    'text/plain',
  ],
};

const ALL_ALLOWED_TYPES = [
  ...ALLOWED_FILE_TYPES.images,
  ...ALLOWED_FILE_TYPES.videos,
  ...ALLOWED_FILE_TYPES.documents,
];

export default function MessageInput() {
  const [content, setContent] = useState('');
  const [isSending, setIsSending] = useState(false);
  const [attachedFiles, setAttachedFiles] = useState<AttachedFile[]>([]);
  const [isDragging, setIsDragging] = useState(false);
  const params = useParams();
  const roomId = params.roomId as string;
  const { send } = useWebSocket();
  const typingTimeoutRef = useRef<NodeJS.Timeout>();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const dragCounterRef = useRef(0);

  const handleTyping = () => {
    send({
      type: WebSocketEventType.USER_TYPING,
      payload: { roomId },
    });

    if (typingTimeoutRef.current) {
      clearTimeout(typingTimeoutRef.current);
    }

    typingTimeoutRef.current = setTimeout(() => {
      send({
        type: WebSocketEventType.USER_STOPPED_TYPING,
        payload: { roomId },
      });
    }, 3000);
  };

  const validateFile = (file: File): string | null => {
    if (file.size > MAX_FILE_SIZE) {
      return `File "${file.name}" exceeds 10MB limit`;
    }

    if (!ALL_ALLOWED_TYPES.includes(file.type)) {
      return `File type "${file.type}" is not supported`;
    }

    return null;
  };

  const generateFilePreview = (file: File): Promise<string | undefined> => {
    return new Promise((resolve) => {
      if (file.type.startsWith('image/')) {
        const reader = new FileReader();
        reader.onload = (e) => resolve(e.target?.result as string);
        reader.onerror = () => resolve(undefined);
        reader.readAsDataURL(file);
      } else {
        resolve(undefined);
      }
    });
  };

  const handleFileSelect = async (files: FileList | null) => {
    if (!files || files.length === 0) return;

    const newFiles: AttachedFile[] = [];

    for (let i = 0; i < files.length; i++) {
      const file = files[i];
      const error = validateFile(file);

      if (error) {
        alert(error);
        continue;
      }

      const preview = await generateFilePreview(file);

      newFiles.push({
        id: `${Date.now()}-${i}`,
        file,
        preview,
        uploading: false,
        progress: 0,
      });
    }

    setAttachedFiles((prev) => [...prev, ...newFiles]);
  };

  const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    handleFileSelect(e.target.files);
    // Reset input value to allow selecting the same file again
    e.target.value = '';
  };

  const removeFile = (fileId: string) => {
    setAttachedFiles((prev) => prev.filter((f) => f.id !== fileId));
  };

  const handleDragEnter = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    dragCounterRef.current++;
    if (e.dataTransfer.items && e.dataTransfer.items.length > 0) {
      setIsDragging(true);
    }
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    dragCounterRef.current--;
    if (dragCounterRef.current === 0) {
      setIsDragging(false);
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
  };

  const handleDrop = async (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
    dragCounterRef.current = 0;

    const files = e.dataTransfer.files;
    await handleFileSelect(files);
  };

  const uploadAttachedFiles = async (): Promise<string[]> => {
    const uploadedFileIds: string[] = [];

    for (const attachedFile of attachedFiles) {
      try {
        setAttachedFiles((prev) =>
          prev.map((f) =>
            f.id === attachedFile.id ? { ...f, uploading: true, progress: 0 } : f
          )
        );

        // Simulate upload progress (replace with actual upload logic)
        const fileId = await uploadFiles(roomId, attachedFile.file, (progress) => {
          setAttachedFiles((prev) =>
            prev.map((f) =>
              f.id === attachedFile.id ? { ...f, progress } : f
            )
          );
        });

        uploadedFileIds.push(fileId);

        setAttachedFiles((prev) =>
          prev.map((f) =>
            f.id === attachedFile.id ? { ...f, uploading: false, progress: 100 } : f
          )
        );
      } catch (error) {
        setAttachedFiles((prev) =>
          prev.map((f) =>
            f.id === attachedFile.id
              ? { ...f, uploading: false, error: 'Upload failed' }
              : f
          )
        );
        throw error;
      }
    }

    return uploadedFileIds;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if ((!content.trim() && attachedFiles.length === 0) || isSending) return;

    setIsSending(true);

    try {
      let attachmentIds: string[] = [];

      // Upload files first if any
      if (attachedFiles.length > 0) {
        attachmentIds = await uploadAttachedFiles();
      }

      // Send message with attachments
      await sendMessage(roomId, {
        content: content.trim(),
        attachmentIds,
      });

      setContent('');
      setAttachedFiles([]);

      send({
        type: WebSocketEventType.USER_STOPPED_TYPING,
        payload: { roomId },
      });
    } catch (error) {
      console.error('Failed to send message:', error);
      alert('Failed to send message. Please try again.');
    } finally {
      setIsSending(false);
    }
  };

  const getFileIcon = (file: File) => {
    if (file.type.startsWith('image/')) return <Image className="w-5 h-5" />;
    if (file.type.startsWith('video/')) return <Film className="w-5 h-5" />;
    return <FileText className="w-5 h-5" />;
  };

  return (
    <div className="border-t border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
      {/* File Attachments Preview */}
      {attachedFiles.length > 0 && (
        <div className="px-4 pt-4 pb-2">
          <div className="flex flex-wrap gap-2">
            {attachedFiles.map((attachedFile) => (
              <div
                key={attachedFile.id}
                className="relative group bg-gray-100 dark:bg-gray-700 rounded-lg p-2 flex items-center gap-2 max-w-xs"
              >
                {/* File Preview/Icon */}
                {attachedFile.preview ? (
                  <img
                    src={attachedFile.preview}
                    alt={attachedFile.file.name}
                    className="w-12 h-12 object-cover rounded"
                  />
                ) : (
                  <div className="w-12 h-12 flex items-center justify-center bg-gray-200 dark:bg-gray-600 rounded text-gray-600 dark:text-gray-300">
                    {getFileIcon(attachedFile.file)}
                  </div>
                )}

                {/* File Info */}
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">
                    {attachedFile.file.name}
                  </p>
                  <p className="text-xs text-gray-500 dark:text-gray-400">
                    {(attachedFile.file.size / 1024).toFixed(1)} KB
                  </p>

                  {/* Upload Progress */}
                  {attachedFile.uploading && (
                    <div className="mt-1">
                      <div className="w-full bg-gray-200 dark:bg-gray-600 rounded-full h-1">
                        <div
                          className="bg-blue-600 h-1 rounded-full transition-all duration-300"
                          style={{ width: `${attachedFile.progress}%` }}
                        />
                      </div>
                    </div>
                  )}

                  {/* Error Message */}
                  {attachedFile.error && (
                    <p className="text-xs text-red-600 dark:text-red-400 mt-1">
                      {attachedFile.error}
                    </p>
                  )}
                </div>

                {/* Remove Button */}
                {!attachedFile.uploading && (
                  <button
                    onClick={() => removeFile(attachedFile.id)}
                    className="absolute -top-2 -right-2 bg-red-500 hover:bg-red-600 text-white rounded-full p-1 opacity-0 group-hover:opacity-100 transition-opacity"
                  >
                    <X className="w-3 h-3" />
                  </button>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Input Area */}
      <div className="p-4">
        <form
          onSubmit={handleSubmit}
          onDragEnter={handleDragEnter}
          onDragLeave={handleDragLeave}
          onDragOver={handleDragOver}
          onDrop={handleDrop}
          className={`relative flex gap-2 ${
            isDragging
              ? 'ring-2 ring-blue-500 ring-offset-2 dark:ring-offset-gray-800 rounded-lg'
              : ''
          }`}
        >
          {/* Drag and Drop Overlay */}
          {isDragging && (
            <div className="absolute inset-0 bg-blue-500/10 dark:bg-blue-500/20 border-2 border-dashed border-blue-500 rounded-lg flex items-center justify-center z-10 pointer-events-none">
              <div className="text-center">
                <Paperclip className="w-12 h-12 mx-auto mb-2 text-blue-600 dark:text-blue-400" />
                <p className="text-sm font-medium text-blue-600 dark:text-blue-400">
                  Drop files to attach
                </p>
              </div>
            </div>
          )}

          {/* File Input (Hidden) */}
          <input
            ref={fileInputRef}
            type="file"
            multiple
            accept={ALL_ALLOWED_TYPES.join(',')}
            onChange={handleFileInputChange}
            className="hidden"
          />

          {/* Attach File Button */}
          <button
            type="button"
            onClick={() => fileInputRef.current?.click()}
            disabled={isSending}
            className="flex-shrink-0 p-2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            title="Attach files"
          >
            <Paperclip className="w-5 h-5" />
          </button>

          {/* Text Input */}
          <input
            type="text"
            value={content}
            onChange={(e) => {
              setContent(e.target.value);
              handleTyping();
            }}
            placeholder="Type a message..."
            className="flex-1 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 placeholder-gray-500 dark:placeholder-gray-400"
            disabled={isSending}
          />

          {/* Send Button */}
          <button
            type="submit"
            disabled={(!content.trim() && attachedFiles.length === 0) || isSending}
            className="flex-shrink-0 px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 dark:disabled:bg-gray-600 text-white rounded-lg font-medium transition-colors disabled:cursor-not-allowed flex items-center gap-2"
          >
            {isSending ? (
              <>
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white" />
                Sending...
              </>
            ) : (
              <>
                <Send className="w-4 h-4" />
                Send
              </>
            )}
          </button>
        </form>

        {/* File Type Info */}
        <p className="text-xs text-gray-500 dark:text-gray-400 mt-2">
          Supported: Images, Videos, Documents (PDF, Word, Excel) • Max 10MB per file
        </p>
      </div>
    </div>
  );
}
```

---

## 11. API Layer

### lib/api/messages.ts

```typescript
import { Message } from '@/types/message';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3000/api';

interface FetchMessagesOptions {
  limit?: number;
  before?: string;
  after?: string;
}

export async function fetchMessages(
  roomId: string,
  options: FetchMessagesOptions = {}
): Promise<Message[]> {
  const { limit = 20, before, after } = options;
  
  const params = new URLSearchParams({
    limit: limit.toString(),
    ...(before && { before }),
    ...(after && { after }),
  });

  const response = await fetch(
    `${API_BASE}/rooms/${roomId}/messages?${params}`,
    {
      headers: {
        'Content-Type': 'application/json',
        // Add auth headers
        Authorization: `Bearer ${getAuthToken()}`,
      },
    }
  );

  if (!response.ok) {
    throw new Error('Failed to fetch messages');
  }

  return response.json();
}

export async function sendMessage(
  roomId: string,
  data: { content: string; attachmentIds?: string[] }
): Promise<Message> {
  const response = await fetch(
    `${API_BASE}/rooms/${roomId}/messages`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${getAuthToken()}`,
      },
      body: JSON.stringify(data),
    }
  );

  if (!response.ok) {
    throw new Error('Failed to send message');
  }

  return response.json();
}

export async function uploadFiles(
  roomId: string,
  file: File,
  onProgress?: (progress: number) => void
): Promise<string> {
  return new Promise((resolve, reject) => {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('roomId', roomId);

    const xhr = new XMLHttpRequest();

    // Track upload progress
    xhr.upload.addEventListener('progress', (e) => {
      if (e.lengthComputable && onProgress) {
        const progress = Math.round((e.loaded / e.total) * 100);
        onProgress(progress);
      }
    });

    xhr.addEventListener('load', () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        try {
          const response = JSON.parse(xhr.responseText);
          resolve(response.fileId);
        } catch (error) {
          reject(new Error('Invalid response from server'));
        }
      } else {
        reject(new Error(`Upload failed with status ${xhr.status}`));
      }
    });

    xhr.addEventListener('error', () => {
      reject(new Error('Network error during upload'));
    });

    xhr.addEventListener('abort', () => {
      reject(new Error('Upload aborted'));
    });

    xhr.open('POST', `${API_BASE}/rooms/${roomId}/attachments`);
    xhr.setRequestHeader('Authorization', `Bearer ${getAuthToken()}`);
    xhr.send(formData);
  });
}

function getAuthToken(): string {
  // Implement your auth token retrieval logic
  return localStorage.getItem('auth_token') || '';
}
```

---

## 12. Responsive Design & Dark Mode

### Tailwind Configuration

All components use Tailwind's dark mode utilities:
- `dark:bg-gray-900` for backgrounds
- `dark:text-gray-100` for text
- `dark:border-gray-700` for borders

### Mobile Responsiveness
- Use `flex-col` layouts for vertical stacking
- `h-screen` ensures full viewport height
- `overflow-hidden` and `overflow-y-auto` for proper scrolling
- Touch-friendly tap targets (minimum 44x44px)

---

## 13. Performance Optimizations

### 1. Zustand Store
- Selector-based subscriptions prevent unnecessary re-renders
- Only components using specific state slices re-render

### 2. Message Virtualization
- For rooms with 1000+ messages, implement `react-window` or `react-virtual`
- Only render visible messages in viewport

### 3. WebSocket Event Batching
- Batch multiple typing events to reduce server load
- Debounce typing indicators

### 4. Lazy Loading
- Load older messages on demand (infinite scroll)
- Prefetch next batch when user scrolls near top

### 5. Memoization
- Use `React.memo()` for MessageItem components
- Prevent re-renders of unchanged messages

---

## 14. Error Handling

### Network Errors
- Display user-friendly error messages
- Implement retry logic for failed API calls
- Show offline indicator when WebSocket disconnects

### State Errors
- Reset store on critical errors
- Log errors to monitoring service (Sentry, LogRocket)

---

## 15. Testing Strategy

### Unit Tests
- Test Zustand store actions and state updates
- Test utility functions and date formatters

### Integration Tests
- Test WebSocket message handling
- Test API calls with mock responses

### E2E Tests
- Test complete user flow: login → select room → send message
- Test real-time message updates
- Test pagination and infinite scroll

---

## 16. Deployment Considerations

### Environment Variables
```
NEXT_PUBLIC_WS_URL=wss://api.example.com/ws
NEXT_PUBLIC_API_URL=https://api.example.com
```

### Build Optimization
- Enable Next.js bundle analyzer
- Tree-shake unused Tailwind classes
- Minimize Zustand store bundle

### Monitoring
- Track WebSocket connection health
- Monitor message delivery latency
- Track render performance with React DevTools Profiler

---

## 17. Future Enhancements

1. **Message Threading**: Nested replies to messages
2. **Rich Text Editor**: Markdown support, @mentions, emojis
3. **File Uploads**: Image/document attachments
4. **Message Search**: Full-text search across rooms
5. **Unread Indicators**: Track last read position
6. **Message Reactions**: Expandable emoji picker
7. **Voice/Video Calls**: WebRTC integration
8. **Presence Indicators**: Online/offline status

---

## 18. Security Considerations

- Validate all messages on server-side
- Sanitize user input to prevent XSS attacks
- Implement rate limiting for message sending
- Use HTTPS/WSS for all connections
- Implement proper authentication/authorization
- Encrypt sensitive data in transit and at rest

---

## Conclusion

This architecture provides a scalable, performant, and maintainable solution for a real-time messaging interface. Zustand offers excellent performance with minimal boilerplate, while the component-based architecture ensures code reusability and testability. The WebSocket integration enables seamless real-time updates, and the responsive design ensures a great user experience across all devices.