import RoomList from '@/components/room/RoomList';

export default function HomePage() {
  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">My Rooms</h1>
      <RoomList />
    </div>
  );
}
