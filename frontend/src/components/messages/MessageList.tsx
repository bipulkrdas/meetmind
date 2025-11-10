'use client';

import { useEffect, useRef } from 'react';
import MessageItem from './MessageItem';
import { Message } from '@/types/message.types';

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
