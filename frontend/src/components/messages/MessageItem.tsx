'use client';

import { Message, MessageType } from '@/types/message.types';
import { formatDistanceToNow } from 'date-fns';
import { FileText } from 'lucide-react';
import TranscriptDisplay from './TranscriptDisplay';

interface MessageItemProps {
  message: Message;
}

export default function MessageItem({ message }: MessageItemProps) {
  return (
    <div className="flex gap-3 group hover:bg-gray-100 dark:hover:bg-gray-800 px-4 py-2 rounded-lg transition-colors">
      {/* Avatar */}
      <div className="flex-shrink-0">
        <div className="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white font-semibold">
          {message.message_type === "meeting_transcript"? "A" : message.username.charAt(0).toUpperCase()}
        </div>
      </div>

      {/* Message Content */}
      <div className="flex-1 min-w-0">
        <div className="flex items-baseline gap-2">
          <span className="font-semibold text-gray-900 dark:text-gray-100">
            {message.message_type === "meeting_transcript" ? "AI Agent" : message.username}
          </span>
          <span className="text-xs text-gray-500 dark:text-gray-400">
            {formatDistanceToNow(new Date(message.created_at), { addSuffix: true })}
          </span>
          {message.edited && (
            <span className="text-xs text-gray-400 dark:text-gray-500">(edited)</span>
          )}
        </div>
        
        {message.message_type === "meeting_transcript" ? (
          <TranscriptDisplay message={message} />
        ) : (
          <>
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
                    onClick={() => window.open(attachment.storage_url, '_blank')}
                  >
                    {attachment.file_type.startsWith('image/') ? (
                      <img
                        src={attachment.storage_url}
                        alt={attachment.file_name}
                        className="w-full max-w-sm rounded-lg"
                      />
                    ) : (
                      <>
                        <FileText className="w-5 h-5 text-blue-600 dark:text-blue-400 flex-shrink-0" />
                        <div className="flex-1 min-w-0">
                          <p className="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">
                            {attachment.file_name}
                          </p>
                          <p className="text-xs text-gray-500 dark:text-gray-400">
                            {(attachment.file_size / 1024).toFixed(1)} KB
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
          </>
        )}
      </div>
    </div>
  );
}
