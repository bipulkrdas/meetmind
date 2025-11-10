'use client';

import React, { useState } from 'react';
import { Message } from '@/types/message.types';
import { format } from 'date-fns';
import TranscriptModal from './TranscriptModal';
import { messageService, FullTranscriptData } from '@/services/api/message.service';

interface TranscriptDisplayProps {
  message: Message;
}

export default function TranscriptDisplay({ message }: TranscriptDisplayProps) {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [transcriptData, setTranscriptData] = useState<FullTranscriptData | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleViewTranscript = async () => {
    if (!message.extra_data?.transcript || !message.room_id || !message.id) {
      setError("Transcript data is missing.");
      return;
    }

    setIsLoading(true);
    setError(null);
    setIsModalOpen(true); // Open modal immediately to show loading state

    const { room_id, id } = message;
    const jsonKey = message.extra_data.transcript.s3_keys.json;

    try {
      const data: FullTranscriptData = await messageService.fetchTranscript(room_id, id, jsonKey);
      setTranscriptData(data);
    } catch (e: any) {
      console.error("Failed to fetch transcript:", e);
      setError(`Failed to load transcript: ${e.message}`);
      setTranscriptData(null); // Clear any previous data
    } finally {
      setIsLoading(false);
    }
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setTranscriptData(null); // Clear data when modal closes
    setError(null);
  };

  const sessionStart = message.extra_data?.transcript?.session_start
    ? new Date(message.extra_data.transcript.session_start)
    : null;
  const sessionEnd = message.extra_data?.transcript?.session_end
    ? new Date(message.extra_data.transcript.session_end)
    : null;

  return (
    <div className="p-4 bg-blue-50 dark:bg-blue-900 rounded-lg shadow-sm border border-blue-200 dark:border-blue-700">
      <h3 className="font-semibold text-blue-800 dark:text-blue-200">Meeting Transcript Available</h3>
      {sessionStart && sessionEnd && (
        <p className="text-sm text-blue-700 dark:text-blue-300 mt-1">
          Session: {format(sessionStart, 'MMM d, yyyy HH:mm')} - {format(sessionEnd, 'HH:mm')}
        </p>
      )}
      <button
        onClick={handleViewTranscript}
        className="mt-3 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        disabled={isLoading}
      >
        {isLoading ? 'Loading...' : 'View Transcript'}
      </button>

      {error && <p className="text-red-500 text-sm mt-2">{error}</p>}

      <TranscriptModal
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        transcript={transcriptData}
      />
    </div>
  );
}
