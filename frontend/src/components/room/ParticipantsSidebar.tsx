'use client';

import { useEffect, useState } from 'react';
import { participantService } from '@/services/api/participant.service';
import AddParticipantModal from './AddParticipantModal';

interface ParticipantsSidebarProps {
  roomId: string;
}

export default function ParticipantsSidebar({ roomId }: ParticipantsSidebarProps) {
  const [participants, setParticipants] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);

  useEffect(() => {
    loadParticipants();
  }, [roomId]);

  const loadParticipants = async () => {
    try {
      const data = await participantService.getParticipants(roomId);
      setParticipants(data);
    } catch (error) {
      console.error('Failed to load participants:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <aside className="w-80 p-4 bg-white rounded-lg shadow-md">
      <h2 className="text-lg font-bold mb-4">Participants ({participants.length})</h2>
      {loading ? (
        <div>Loading...</div>
      ) : (
        <ul>
          {participants.map(p => (
            <li key={p.id} className="flex items-center justify-between py-1">
              <span>{p.name}</span>
              <span className="text-sm text-gray-500">{p.role}</span>
            </li>
          ))}
        </ul>
      )}
      <button 
        className="w-full mt-4 px-4 py-2 bg-blue-600 text-white rounded"
        onClick={() => setIsModalOpen(true)}
      >
        + Add Participant
      </button>
      <AddParticipantModal 
        isOpen={isModalOpen} 
        onClose={() => setIsModalOpen(false)} 
        roomId={roomId} 
        onParticipantAdded={loadParticipants}
      />
    </aside>
  );
}
