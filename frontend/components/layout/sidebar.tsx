'use client'

import { useState } from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import {
  LayoutDashboard,
  BookOpen,
  Settings,
  Menu,
  X,
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'

const navItems = [
  {
    label: 'Dashboard',
    href: '/dashboard',
    icon: LayoutDashboard,
  },
  {
    label: 'Exams',
    href: '/exams',
    icon: BookOpen,
  },
  {
    label: 'Settings',
    href: '/settings',
    icon: Settings,
  },
]

export function Sidebar() {
  const pathname = usePathname()
  const [open, setOpen] = useState(true)

  return (
    <>
      {/* Mobile toggle */}
      <Button
        variant="ghost"
        size="icon"
        className="fixed top-4 left-4 z-40 md:hidden"
        onClick={() => setOpen(!open)}
      >
        {open ? <X /> : <Menu />}
      </Button>

      {/* Sidebar */}
      <aside
        className={cn(
          'bg-slate-950 border-r border-slate-800 transition-all duration-300 fixed md:relative h-full z-30',
          open ? 'w-64' : 'w-0 md:w-64 overflow-hidden'
        )}
      >
        <div className="p-6 space-y-8 h-full flex flex-col">
          {/* Logo */}
          <Link href="/dashboard" className="flex items-center gap-3 font-bold text-xl">
            <div className="w-9 h-9 bg-gradient-to-br from-blue-500 to-blue-700 rounded-lg flex items-center justify-center text-white text-sm font-bold shadow-lg shadow-blue-500/20">
              H
            </div>
            <span className="hidden sm:inline text-white tracking-tight">HARaMA</span>
          </Link>

          {/* Navigation */}
          <nav className="space-y-1.5 flex-1">
            {navItems.map((item) => {
              const Icon = item.icon
              const isActive = pathname === item.href || pathname.startsWith(item.href + '/')

              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={cn(
                    'flex items-center gap-3 px-4 py-2.5 rounded-lg transition-all duration-200',
                    isActive
                      ? 'bg-blue-600/15 text-blue-400 font-medium border border-blue-500/20'
                      : 'text-slate-400 hover:bg-slate-800/80 hover:text-slate-200'
                  )}
                  onClick={() => setOpen(false)}
                >
                  <Icon className="w-5 h-5" />
                  <span className="text-sm">{item.label}</span>
                </Link>
              )
            })}
          </nav>

          {/* Bottom branding */}
          <div className="pt-4 border-t border-slate-800">
            <p className="text-xs text-slate-600">AI-Powered Grading</p>
          </div>
        </div>
      </aside>

      {/* Mobile overlay */}
      {open && (
        <div
          className="fixed inset-0 bg-black/50 z-20 md:hidden"
          onClick={() => setOpen(false)}
        />
      )}
    </>
  )
}
