import Link from 'next/link';

export function Sidebar() {
  return (
    <aside className="w-64 bg-white border-r border-gray-200 min-h-screen flex flex-col">
      <div className="p-6 border-b border-gray-200">
        <h1 className="text-2xl font-bold text-blue-600">EchoMind</h1>
      </div>
      <nav className="flex-1 p-4 space-y-1">
        <NavItem href="/dashboard" label="Inbox" active />
        <NavItem href="/dashboard/tasks" label="Tasks" />
        <NavItem href="/dashboard/insights" label="Insights" />
        <NavItem href="/dashboard/settings" />
      </nav>
      <div className="p-4 border-t border-gray-200 mt-auto">
        <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Categories</h3>
        <div className="space-y-1">
          <NavItem href="/dashboard?category=Work" label="Work" />
          <NavItem href="/dashboard?category=Personal" label="Personal" />
          <NavItem href="/dashboard?category=Newsletter" label="Newsletter" />
          <NavItem href="/dashboard?category=Notification" label="Notification" />
          <NavItem href="/dashboard?category=Spam" label="Spam" />
        </div>
      </div>
    </aside>
  );
}

function NavItem({ href, label, active = false }: { href: string; label: string; active?: boolean }) {
  return (
    <Link
      href={href}
      className={`flex items-center px-4 py-2 text-sm font-medium rounded-md ${
        active
          ? 'bg-blue-50 text-blue-700'
          : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
      }`}
    >
      {label}
    </Link>
  );
}
