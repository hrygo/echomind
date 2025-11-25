'use client';

import * as React from 'react';
import { Moon, Sun, Monitor } from 'lucide-react';
import { useTheme } from './ThemeProvider';

export function ThemeToggle() {
  const { setTheme, theme } = useTheme();

  const toggleTheme = () => {
    if (theme === 'light') {
      setTheme('dark');
    } else if (theme === 'dark') {
      setTheme('system');
    } else {
      setTheme('light');
    }
  };

  const getIcon = () => {
    if (theme === 'light') {
      return <Sun className="h-[1.2rem] w-[1.2rem]" />;
    } else if (theme === 'dark') {
      return <Moon className="h-[1.2rem] w-[1.2rem]" />;
    } else {
      return <Monitor className="h-[1.2rem] w-[1.2rem]" />;
    }
  };

  return (
    <button
      onClick={toggleTheme}
      className="inline-flex items-center justify-center rounded-full text-sm font-medium transition-all duration-200 hover:bg-slate-100 dark:hover:bg-slate-800 p-2 text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-100"
      title={`当前主题: ${theme === 'light' ? '明亮' : theme === 'dark' ? '黑暗' : '跟随系统'}。点击切换: ${theme === 'light' ? '黑暗' : theme === 'dark' ? '跟随系统' : '明亮'}`}
    >
      {getIcon()}
      <span className="sr-only">Toggle theme</span>
    </button>
  );
}