import SignInForm from '@/components/auth/SignInForm';
import Link from 'next/link';

export default function SignInPage() {
  return (
    <div className="flex justify-center items-center min-h-screen">
      <div className="w-full max-w-md p-8 space-y-8 bg-white rounded-lg shadow-md">
        <h1 className="text-2xl font-bold text-center">Sign In</h1>
        <SignInForm />
        <div className="text-center">
          <Link href="/auth/reset-password">
            <span className="text-sm text-blue-600 hover:underline">Forgot Password?</span>
          </Link>
        </div>
      </div>
    </div>
  );
}
