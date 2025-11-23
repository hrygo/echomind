"use client";

import { X, Folder, Hash, Users } from "lucide-react";
import { useState } from "react";
import { useContextStore } from "@/lib/store/contexts";
import { ContextInput } from '@/lib/api/contexts';
import { cn } from '@/lib/utils';
import { useLanguage } from '@/lib/i18n/LanguageContext';

interface CreateContextModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const COLORS = [
  { name: 'blue', bg: 'bg-blue-100', text: 'text-blue-700', border: 'border-blue-200' },
  { name: 'green', bg: 'bg-emerald-100', text: 'text-emerald-700', border: 'border-emerald-200' },
  { name: 'purple', bg: 'bg-purple-100', text: 'text-purple-700', border: 'border-purple-200' },
  { name: 'amber', bg: 'bg-amber-100', text: 'text-amber-700', border: 'border-amber-200' },
  { name: 'rose', bg: 'bg-rose-100', text: 'text-rose-700', border: 'border-rose-200' },
];

export function CreateContextModal({ isOpen, onClose }: CreateContextModalProps) {
  const { t } = useLanguage();
  const { addContext } = useContextStore();
  const [loading, setLoading] = useState(false);

  const [formData, setFormData] = useState<ContextInput>({
    name: '',
    color: 'blue',
    keywords: [],
    stakeholders: []
  });

  const [keywordInput, setKeywordInput] = useState('');
  const [stakeholderInput, setStakeholderInput] = useState('');

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name) return;

    setLoading(true);
    try {
      await addContext(formData);
      onClose();
      // Reset form
      setFormData({ name: '', color: 'blue', keywords: [], stakeholders: [] });
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  const addTag = (type: 'keywords' | 'stakeholders', value: string) => {
    if (!value.trim()) return;
    if (formData[type].includes(value.trim())) return;
    
    setFormData(prev => ({
      ...prev,
      [type]: [...prev[type], value.trim()]
    }));
  };

  const removeTag = (type: 'keywords' | 'stakeholders', index: number) => {
    setFormData(prev => ({
      ...prev,
      [type]: prev[type].filter((_, i) => i !== index)
    }));
  };

  const handleKeyDown = (e: React.KeyboardEvent, type: 'keywords' | 'stakeholders', value: string, setter: (v: string) => void) => {
    if (e.key === 'Enter' || e.key === ',') {
      e.preventDefault();
      addTag(type, value);
      setter('');
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/20 backdrop-blur-sm">
      <div className="bg-white rounded-2xl shadow-xl w-full max-w-md border border-slate-100 overflow-hidden animate-in fade-in zoom-in-95 duration-200">
        
        {/* Header */}
        <div className="px-6 py-4 border-b border-slate-100 flex items-center justify-between bg-slate-50/50">
          <h2 className="text-lg font-semibold text-slate-800 flex items-center gap-2">
            <Folder className="w-5 h-5 text-blue-600" />
            {t('createContextModal.title')}
          </h2>
          <button onClick={onClose} className="text-slate-400 hover:text-slate-600 transition-colors">
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-6 space-y-6">
          
          {/* Name & Color */}
          <div className="space-y-4">
            <div>
              <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5">{t('createContextModal.nameLabel')}</label>
              <input
                type="text"
                placeholder={t('createContextModal.namePlaceholder')}
                className="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                value={formData.name}
                onChange={e => setFormData(prev => ({ ...prev, name: e.target.value }))}
                autoFocus
              />
            </div>

            <div>
              <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-2">{t('createContextModal.colorLabel')}</label>
              <div className="flex gap-3">
                {COLORS.map(c => (
                  <button
                    key={c.name}
                    type="button"
                    onClick={() => setFormData(prev => ({ ...prev, color: c.name }))}
                    className={cn(
                      "w-8 h-8 rounded-full flex items-center justify-center transition-all ring-offset-2",
                      c.bg, c.border,
                      formData.color === c.name ? "ring-2 ring-blue-600 scale-110" : "hover:scale-105"
                    )}
                  >
                    {formData.color === c.name && <div className={`w-2.5 h-2.5 rounded-full ${c.text} bg-current`} />}
                  </button>
                ))}
              </div>
            </div>
          </div>

          <div className="h-px bg-slate-100" />

          {/* Rules Section */}
          <div className="space-y-4">
            <h3 className="text-sm font-medium text-slate-900">{t('createContextModal.rulesHeader')}</h3>
            
            {/* Keywords */}
            <div>
              <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5 flex items-center gap-1.5">
                <Hash className="w-3 h-3" /> {t('createContextModal.keywordsLabel')}
              </label>
              <div className="flex flex-wrap gap-2 mb-2">
                {formData.keywords.map((tag, i) => (
                  <span key={i} className="inline-flex items-center px-2 py-1 rounded bg-slate-100 border border-slate-200 text-xs font-medium text-slate-600">
                    {tag}
                    <button type="button" onClick={() => removeTag('keywords', i)} className="ml-1 text-slate-400 hover:text-red-500">
                      <X className="w-3 h-3" />
                    </button>
                  </span>
                ))}
              </div>
              <input
                type="text"
                placeholder={t('createContextModal.keywordsPlaceholder')}
                className="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm focus:outline-none focus:border-blue-500 transition-all"
                value={keywordInput}
                onChange={e => setKeywordInput(e.target.value)}
                onKeyDown={e => handleKeyDown(e, 'keywords', keywordInput, setKeywordInput)}
              />
            </div>

            {/* Stakeholders */}
            <div>
              <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1.5 flex items-center gap-1.5">
                <Users className="w-3 h-3" /> {t('createContextModal.stakeholdersLabel')}
              </label>
              <div className="flex flex-wrap gap-2 mb-2">
                {formData.stakeholders.map((tag, i) => (
                  <span key={i} className="inline-flex items-center px-2 py-1 rounded bg-slate-100 border border-slate-200 text-xs font-medium text-slate-600">
                    {tag}
                    <button type="button" onClick={() => removeTag('stakeholders', i)} className="ml-1 text-slate-400 hover:text-red-500">
                      <X className="w-3 h-3" />
                    </button>
                  </span>
                ))}
              </div>
              <input
                type="text"
                placeholder={t('createContextModal.stakeholdersPlaceholder')}
                className="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm focus:outline-none focus:border-blue-500 transition-all"
                value={stakeholderInput}
                onChange={e => setStakeholderInput(e.target.value)}
                onKeyDown={e => handleKeyDown(e, 'stakeholders', stakeholderInput, setStakeholderInput)}
              />
            </div>
          </div>

          {/* Footer */}
          <div className="pt-2 flex justify-end gap-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-sm font-medium text-slate-600 hover:bg-slate-100 rounded-lg transition-colors"
            >
              {t('createContextModal.cancel')}
            </button>
            <button
              type="submit"
              disabled={!formData.name || loading}
              className="px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:opacity-50 rounded-lg shadow-sm shadow-blue-200 transition-all"
            >
              {loading ? t('createContextModal.creating') : t('createContextModal.createContext')}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
