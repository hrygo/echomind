import { useState, useEffect } from 'react';

const STORAGE_KEY = 'echomind_search_history';
const MAX_HISTORY_ITEMS = 10;

export function useSearchHistory() {
  const [history, setHistory] = useState<string[]>([]);

  useEffect(() => {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored) {
      try {
        const parsed = JSON.parse(stored);
        setTimeout(() => setHistory(parsed), 0);
      } catch (e) {
        console.error('Failed to parse search history', e);
      }
    }
  }, []);

  const addToHistory = (query: string) => {
    if (!query.trim()) return;
    const trimmed = query.trim();
    
    setHistory((prev) => {
      const newHistory = [trimmed, ...prev.filter(item => item !== trimmed)].slice(0, MAX_HISTORY_ITEMS);
      localStorage.setItem(STORAGE_KEY, JSON.stringify(newHistory));
      return newHistory;
    });
  };

  const removeFromHistory = (query: string) => {
    setHistory((prev) => {
      const newHistory = prev.filter(item => item !== query);
      localStorage.setItem(STORAGE_KEY, JSON.stringify(newHistory));
      return newHistory;
    });
  };

  const clearHistory = () => {
    setHistory([]);
    localStorage.removeItem(STORAGE_KEY);
  };

  return {
    history,
    addToHistory,
    removeFromHistory,
    clearHistory
  };
}
