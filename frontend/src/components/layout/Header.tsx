'use client';

import Link from 'next/link';
import Image from 'next/image';
import DropdownMenu from './DropdownMenu';
import logo from '@/app/icon.png';

export default function Header() {
  return (
    <header className="flex items-center justify-between p-4 bg-white border-b">
      <Link href="/app/home" className="flex items-center gap-2">
        <div className="p-1 bg-gray-100 rounded-lg">
          <Image src={logo} alt="My Meetspace Logo" width={64} height={64} />
        </div>
        <span className="text-xl font-bold">My Meetspace</span>
      </Link>
      <div>
        <DropdownMenu />
      </div>
    </header>
  );
}
