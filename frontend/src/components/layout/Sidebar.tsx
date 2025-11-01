'use client';

import Link from 'next/link';

export default function Sidebar() {
  return (
    <aside className="w-64 p-4 bg-gray-50 border-r">
      <nav>
        <ul>
          <li>
            <Link href="/home">
              <span className="text-gray-700 hover:text-blue-600">Home</span>
            </Link>
          </li>
        </ul>
      </nav>
    </aside>
  );
}
