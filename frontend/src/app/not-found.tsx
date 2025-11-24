'use client';

import Link from 'next/link'

export default function NotFound() {
    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100 text-slate-800">
            <h2 className="text-4xl font-bold mb-4">404</h2>
            <p className="mb-8 text-lg text-slate-600">Page Not Found</p>
            <Link href="/dashboard" className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
                Return to Dashboard
            </Link>
        </div>
    )
}
