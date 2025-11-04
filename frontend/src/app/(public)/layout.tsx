"use client";
import { Video } from "lucide-react";
import Link from "next/link";
import { useState } from "react";

const Header = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  return (
    <header className="bg-white shadow-md">
      <div className="container mx-auto px-4 py-4 flex justify-between items-center">
        <div className="text-2xl font-bold">
          <div className="flex items-center space-x-3">
                <div className="bg-blue-600 p-2 rounded-lg">
                  <Video className="w-6 h-6 text-white" />
                </div>
                <span className="text-xl font-bold text-gray-900 dark:text-gray-900"> <Link href="/">Meetspace</Link></span>
              </div>

        </div>
        <nav className="hidden md:flex items-center space-x-4">
          <Link href="/how-it-works" className="text-gray-600 hover:text-gray-900">
            How it works
          </Link>
          <Link href="/contact" className="text-gray-600 hover:text-gray-900">
            Contact us
          </Link>
        </nav>
        <div className="hidden md:flex items-center space-x-4">
          <Link href="/auth/signin" className="text-gray-600 hover:text-gray-900">
            SignIn
          </Link>
          <Link href="/auth/signup" className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">
            New User, Signup Here
          </Link>
        </div>
        <div className="md:hidden flex items-center">
          <button onClick={() => setIsMenuOpen(!isMenuOpen)} className="text-gray-600 hover:text-gray-900 focus:outline-none">
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d={isMenuOpen ? "M6 18L18 6M6 6l12 12" : "M4 6h16M4 12h16M4 18h16"}></path>
            </svg>
          </button>
        </div>
      </div>
      {isMenuOpen && (
        <div className="md:hidden">
          <nav className="px-4 pt-2 pb-4 space-y-2">
            <Link href="/how-it-works" className="block text-gray-600 hover:text-gray-900">
              How it works
            </Link>
            <Link href="/contact" className="block text-gray-600 hover:text-gray-900">
              Contact us
            </Link>
            <Link href="/auth/signin" className="block text-gray-600 hover:text-gray-900">
              SignIn
            </Link>
            <Link href="/auth/signup" className="block bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">
              New User, Signup Here
            </Link>
          </nav>
        </div>
      )}
    </header>
  );
};

const Footer = () => {
  return (
    <footer className="bg-gray-800 text-white py-8">
      <div className="container mx-auto px-4">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
          <div>
            <h3 className="font-bold mb-4">Company</h3>
            <ul className="space-y-2">
              <li><Link href="/contact" className="hover:underline">Contact Us</Link></li>
              <li><Link href="/support" className="hover:underline">Support</Link></li>
            </ul>
          </div>
          <div>
            <h3 className="font-bold mb-4">How it works</h3>
            <ul className="space-y-2">
              <li><Link href="/how-it-works" className="hover:underline">How it works</Link></li>
            </ul>
          </div>
          <div>
            <h3 className="font-bold mb-4">Legal</h3>
            <ul className="space-y-2">
              <li><Link href="/terms-of-service" className="hover:underline">Terms of Service</Link></li>
              <li><Link href="/privacy-policy" className="hover:underline">Privacy Policy</Link></li>
            </ul>
          </div>
          <div>
            <h3 className="font-bold mb-4">Resources</h3>
            <ul className="space-y-2">
              <li><Link href="/blog" className="hover:underline">Blog</Link></li>
              <li><Link href="/white-papers" className="hover:underline">White Papers</Link></li>
            </ul>
          </div>
        </div>
        <div className="mt-8 pt-8 border-t border-gray-700 text-center text-sm">
          <p>&copy; {new Date().getFullYear()} AIMeeting. All rights reserved.</p>
        </div>
      </div>
    </footer>
  );
};


export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div className="flex flex-col min-h-screen">
      <Header/>
      <main className="flex-grow">{children}</main>
      <Footer/>
    </div>
  );
}