import ProfileDetails from '@/components/profile/ProfileDetails';

export default function ProfilePage() {
  return (
    <div className="flex justify-center items-center min-h-screen">
      <div className="w-full max-w-md p-8 space-y-8 bg-white rounded-lg shadow-md">
        <h1 className="text-2xl font-bold text-center">Your Profile</h1>
        <ProfileDetails />
      </div>
    </div>
  );
}