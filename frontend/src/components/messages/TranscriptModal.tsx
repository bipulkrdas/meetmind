'use client';

import React from 'react';
import { X } from 'lucide-react';
import { format } from 'date-fns';
import { FullTranscriptData } from '@/services/api/message.service'; // Import from message.service.ts

interface TranscriptModalProps {
  isOpen: boolean;
  onClose: () => void;
  transcript: FullTranscriptData | null;
}

export default function TranscriptModal({ isOpen, onClose, transcript }: TranscriptModalProps) {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 p-4">
      <div className="relative w-full max-w-3xl max-h-[90vh] bg-white dark:bg-gray-800 rounded-lg shadow-xl flex flex-col">
        <div className="flex justify-between items-center p-4 border-b border-gray-200 dark:border-gray-700">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">Meeting Transcript</h2>
          <button onClick={onClose} className="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
            <X className="w-6 h-6" />
          </button>
        </div>
        <div className="p-4 overflow-y-auto flex-1">
          {transcript ? (
            <div>
              <p className="text-sm text-gray-600 dark:text-gray-300 mb-2">
                <span className="font-medium">Room:</span> {transcript.room_name}
              </p>
              <p className="text-sm text-gray-600 dark:text-gray-300 mb-4">
                <span className="font-medium">Session:</span> {format(new Date(transcript.session_start), 'MMM d, yyyy HH:mm')} - {format(new Date(transcript.session_end), 'HH:mm')}
              </p>
              <div className="space-y-4">
                {transcript.items.map((item, index) => (
                  <div key={index} className="flex items-start gap-3">
                    <div className="flex-shrink-0 w-16 text-xs text-gray-500 dark:text-gray-400">
                      {format(new Date(item.timestamp), 'HH:mm:ss')}
                    </div>
                    <div className="flex-1">
                      <p className="font-semibold text-gray-900 dark:text-gray-100">
                        {item.speaker_name || item.speaker_identity}:
                      </p>
                      {item.content.map((c, cIdx) => (
                        <p key={cIdx} className="text-gray-800 dark:text-gray-200">
                          {c.text}
                        </p>
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          ) : (
            <p className="text-gray-500 dark:text-gray-400">Loading transcript...</p>
          )}
        </div>
      </div>
    </div>
  );
}
