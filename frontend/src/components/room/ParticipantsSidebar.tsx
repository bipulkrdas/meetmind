import { Fragment, useEffect, useState } from 'react';
import { Dialog, Transition } from '@headlessui/react';
import { participantService } from '@/services/api/participant.service';
import AddParticipantModal from './AddParticipantModal';

interface ParticipantsSidebarProps {
  roomId: string;
}

export default function ParticipantsSidebar({ roomId }: ParticipantsSidebarProps) {
  const [participants, setParticipants] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [isLinkModalOpen, setIsLinkModalOpen] = useState(false);
  const [meetingLink, setMeetingLink] = useState('');

  useEffect(() => {
    loadParticipants();
  }, [roomId]);

  const loadParticipants = async () => {
    setLoading(true);
    try {
      const data = await participantService.getParticipants(roomId);
      setParticipants(data);
    } catch (error) {
      console.error('Failed to load participants:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleInviteToJoin = async () => {
    try {
      await participantService.inviteParticipantsToJoinMeeting(roomId);
      alert('Invitations sent successfully!');
    } catch (error) {
      console.error('Failed to send invitations:', error);
      alert('Error sending invitations. Please try again.');
    }
  };

  const handleGenerateLink = async () => {
    try {
      const data = await participantService.generateMeetingUrl(roomId);
      setMeetingLink(data.meeting_url);
      setIsLinkModalOpen(true);
    } catch (error) {
      console.error('Failed to generate meeting link:', error);
      alert('Error generating meeting link. Please try again.');
    }
  };

  const copyToClipboard = () => {
    navigator.clipboard.writeText(meetingLink);
    alert('Link copied to clipboard!');
  };

  return (
    <aside className="w-80 p-4 bg-white rounded-lg shadow-md">
      <h2 className="text-lg font-bold mb-4">Participants ({participants.length})</h2>
      {loading ? (
        <div>Loading...</div>
      ) : (
        <ul>
          {participants.map((p) => (
            <li key={p.id} className="flex items-center justify-between py-1">
              <span>{p.name}</span>
              <span className="text-sm text-gray-500">{p.role}</span>
            </li>
          ))}
        </ul>
      )}
      <button
        className="w-full mt-4 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
        onClick={() => setIsAddModalOpen(true)}
      >
        + Add Participant
      </button>
      <button
        className="w-full mt-2 px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
        onClick={handleInviteToJoin}
      >
        Invite To Join
      </button>
      <button
        className="w-full mt-2 px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700"
        onClick={handleGenerateLink}
      >
        Generate Meeting Link
      </button>

      <AddParticipantModal
        isOpen={isAddModalOpen}
        onClose={() => setIsAddModalOpen(false)}
        roomId={roomId}
        onParticipantAdded={loadParticipants}
      />

      <Transition appear show={isLinkModalOpen} as={Fragment}>
        <Dialog as="div" className="relative z-50" onClose={() => setIsLinkModalOpen(false)}>
          <Transition.Child
            as={Fragment}
            enter="ease-out duration-300"
            enterFrom="opacity-0"
            enterTo="opacity-100"
            leave="ease-in duration-200"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
          >
            <div className="fixed inset-0 bg-black bg-opacity-50" />
          </Transition.Child>

          <div className="fixed inset-0 overflow-y-auto">
            <div className="flex min-h-full items-center justify-center p-4 text-center">
              <Transition.Child
                as={Fragment}
                enter="ease-out duration-300"
                enterFrom="opacity-0 scale-95"
                enterTo="opacity-100 scale-100"
                leave="ease-in duration-200"
                leaveFrom="opacity-100 scale-100"
                leaveTo="opacity-0 scale-95"
              >
                <Dialog.Panel className="w-full max-w-md transform overflow-hidden rounded-2xl bg-white dark:bg-gray-800 p-6 text-left align-middle shadow-xl transition-all">
                  <Dialog.Title
                    as="h3"
                    className="text-lg font-medium leading-6 text-gray-900 dark:text-white"
                  >
                    Generated Meeting Link
                  </Dialog.Title>
                  <div className="mt-2">
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                      Share this link for others to join.
                    </p>
                  </div>

                  <div className="mt-4 flex items-center space-x-2">
                    <input
                      type="text"
                      value={meetingLink}
                      readOnly
                      className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-gray-200"
                    />
                    <button
                      type="button"
                      onClick={copyToClipboard}
                      className="inline-flex justify-center rounded-md border border-transparent bg-blue-100 dark:bg-blue-800 px-4 py-2 text-sm font-medium text-blue-900 dark:text-blue-100 hover:bg-blue-200 dark:hover:bg-blue-700 focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2"
                    >
                      Copy
                    </button>
                  </div>

                  <div className="mt-4">
                    <button
                      type="button"
                      className="inline-flex justify-center rounded-md border border-transparent bg-gray-100 dark:bg-gray-600 px-4 py-2 text-sm font-medium text-gray-900 dark:text-gray-200 hover:bg-gray-200 dark:hover:bg-gray-500 focus:outline-none focus-visible:ring-2 focus-visible:ring-gray-500 focus-visible:ring-offset-2"
                      onClick={() => setIsLinkModalOpen(false)}
                    >
                      Close
                    </button>
                  </div>
                </Dialog.Panel>
              </Transition.Child>
            </div>
          </div>
        </Dialog>
      </Transition>
    </aside>
  );
}
