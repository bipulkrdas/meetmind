import NewPasswordForm from '@/components/auth/NewPasswordForm';
import { Suspense } from 'react';

export default function ResetPasswordFormPage() {
  return (
    <div className="flex justify-center items-center min-h-screen">
      <div className="w-full max-w-md p-8 space-y-8 bg-white rounded-lg shadow-md">
        <h1 className="text-2xl font-bold text-center">Set New Password</h1>
        <Suspense fallback={<div>Loading...</div>}>
          <NewPasswordForm />
        </Suspense>
      </div>
    </div>
  );
}
