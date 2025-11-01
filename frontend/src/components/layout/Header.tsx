'use client';

import Link from 'next/link';
import DropdownMenu from './DropdownMenu';

export default function Header() {
  return (
    <header className="flex items-center justify-between p-4 bg-white border-b">
      <div>
        <Link href="/app/home">
          <span className="text-xl font-bold">My Rooms</span>
        </Link>
      </div>
      <div>
        <DropdownMenu />
      </div>
    </header>
  );
}
