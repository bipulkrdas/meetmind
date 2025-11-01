'use client';

import { useState } from 'react';
import { participantService } from '@/services/api/participant.service';
import Modal from '@/components/common/Modal';

interface AddParticipantModalProps {
  isOpen: boolean;
  onClose: () => void;
  roomId: string;
  onParticipantAdded: () => void;
}

export default function AddParticipantModal({ isOpen, onClose, roomId, onParticipantAdded }: AddParticipantModalProps) {
  const [formData, setFormData] = useState({ email: '', name: '' });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await participantService.addParticipant(roomId, formData);
      onParticipantAdded();
      onClose();
    } catch (err: any) {
      setError(err.message || 'Failed to add participant');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <h2 className="text-lg font-bold mb-4">Add Participant</h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <input
          type="email"
          name="email"
          placeholder="Email"
          value={formData.email}
          onChange={handleChange}
          required
          className="w-full px-4 py-2 border rounded"
        />
        <input
          type="text"
          name="name"
          placeholder="Name"
          value={formData.name}
          onChange={handleChange}
          required
          className="w-full px-4 py-2 border rounded"
        />
        {error && <div className="text-red-500">{error}</div>}
        <div className="flex justify-end gap-4">
          <button type="button" onClick={onClose} className="px-4 py-2 bg-gray-200 rounded">Cancel</button>
          <button type="submit" disabled={loading} className="px-4 py-2 bg-blue-600 text-white rounded disabled:bg-gray-400">
            {loading ? 'Adding...' : 'Add'}
          </button>
        </div>
      </form>
    </Modal>
  );
}
