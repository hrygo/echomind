import React, { useState } from 'react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { api } from '@/lib/api';
import { Loader2, Sparkles, Copy, Mail, CheckCircle } from 'lucide-react';
import { useLanguage } from '@/lib/i18n/LanguageContext';

interface AIDraftReplyModalProps {
  emailContent: string;
  isOpen: boolean;
  onClose: () => void;
}

export default function AIDraftReplyModal({ emailContent, isOpen, onClose }: AIDraftReplyModalProps) {
  const { t } = useLanguage();
  const [userPrompt, setUserPrompt] = useState(t('emailDetail.aiDraft.tonePrompts.professional'));
  const [draftReply, setDraftReply] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  const handleGenerateDraft = async () => {
    if (!emailContent.trim()) {
      setError(t('emailDetail.aiDraft.errors.emailRequired'));
      return;
    }
    if (!userPrompt.trim()) {
      setError(t('emailDetail.aiDraft.errors.promptRequired'));
      return;
    }

    setLoading(true);
    setError(null);
    setDraftReply('');
    try {
      const response = await api.post('/ai/draft', {
        emailContent: emailContent.trim(),
        userPrompt: userPrompt.trim()
      });
      setDraftReply(response.data.draft);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || t('emailDetail.aiDraft.errors.generateFailed'));
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error('Failed to copy text:', err);
    }
  };

  const openInEmailClient = () => {
    window.location.href = `mailto:?body=${encodeURIComponent(draftReply)}`;
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-4xl w-[calc(100vw-2rem)] sm:w-full max-h-[90vh] overflow-hidden flex flex-col">
        <DialogHeader className="pb-4 sm:pb-6 flex-shrink-0">
          <div className="flex items-center gap-3 sm:gap-4 pr-8">
            <div className="w-10 h-10 sm:w-12 sm:h-12 bg-blue-100 dark:bg-blue-900/30 rounded-xl flex items-center justify-center flex-shrink-0">
              <Sparkles className="w-5 h-5 sm:w-6 sm:h-6 text-blue-600 dark:text-blue-400" />
            </div>
            <div className="flex-1 min-w-0">
              <DialogTitle className="text-lg sm:text-xl font-bold text-foreground truncate">
                {t('emailDetail.aiDraft.title')}
              </DialogTitle>
              <DialogDescription className="text-sm sm:text-base text-muted-foreground mt-1 line-clamp-2">
                {t('emailDetail.aiDraft.description')}
              </DialogDescription>
            </div>
          </div>
        </DialogHeader>

        <div className="flex-1 overflow-y-auto overflow-x-hidden space-y-4 sm:space-y-6 pr-2">
          {/* Email Preview */}
          <div className="space-y-2 sm:space-y-3">
            <Label className="text-sm font-semibold text-foreground flex items-center gap-2">
              <Mail className="w-4 h-4 flex-shrink-0" />
              <span className="truncate">{t('emailDetail.aiDraft.originalEmail')}</span>
            </Label>
            <div className="p-3 sm:p-4 bg-muted/50 rounded-xl border max-h-32 sm:max-h-48 overflow-y-auto overflow-x-hidden">
              <p className="text-sm text-muted-foreground whitespace-pre-wrap leading-relaxed break-words">
                {emailContent.trim() || t('emailDetail.aiDraft.noContent')}
              </p>
            </div>
          </div>

          {/* Prompt Input */}
          <div className="space-y-2 sm:space-y-3">
            <Label htmlFor="userPrompt" className="text-sm font-semibold text-foreground truncate block">
              {t('emailDetail.aiDraft.instructions')}
            </Label>
            <Textarea
              id="userPrompt"
              value={userPrompt}
              onChange={(e) => setUserPrompt(e.target.value)}
              placeholder={t('emailDetail.aiDraft.instructionsPlaceholder')}
              className="min-h-[80px] sm:min-h-[100px] resize-none w-full"
            />
            <div className="flex flex-wrap gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setUserPrompt(t('emailDetail.aiDraft.tonePrompts.professional'))}
                className="text-xs h-8 px-3 whitespace-nowrap"
              >
                {t('emailDetail.aiDraft.toneButtons.professional')}
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setUserPrompt(t('emailDetail.aiDraft.tonePrompts.casual'))}
                className="text-xs h-8 px-3 whitespace-nowrap"
              >
                {t('emailDetail.aiDraft.toneButtons.casual')}
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setUserPrompt(t('emailDetail.aiDraft.tonePrompts.concise'))}
                className="text-xs h-8 px-3 whitespace-nowrap"
              >
                {t('emailDetail.aiDraft.toneButtons.concise')}
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setUserPrompt(t('emailDetail.aiDraft.tonePrompts.detailed'))}
                className="text-xs h-8 px-3 whitespace-nowrap"
              >
                {t('emailDetail.aiDraft.toneButtons.detailed')}
              </Button>
            </div>
          </div>

          {/* Generate Button */}
          <div className="flex justify-center">
            <Button
              onClick={handleGenerateDraft}
              disabled={loading}
              className="w-full sm:w-auto bg-blue-600 hover:bg-blue-700 text-white px-6 sm:px-8 py-3 rounded-xl font-medium shadow-lg shadow-blue-600/20 flex items-center justify-center gap-2 min-h-[44px] sm:min-h-[48px] whitespace-nowrap"
            >
              {loading ? (
                <>
                  <Loader2 className="w-4 h-4 animate-spin flex-shrink-0" />
                  <span className="truncate">{t('emailDetail.aiDraft.generating')}</span>
                </>
              ) : (
                <>
                  <Sparkles className="w-4 h-4 flex-shrink-0" />
                  <span className="truncate">{t('emailDetail.aiDraft.generateButton')}</span>
                </>
              )}
            </Button>
          </div>

          {/* Error Display */}
          {error && (
            <div className="p-3 sm:p-4 bg-destructive/10 border border-destructive/20 rounded-xl overflow-hidden">
              <p className="text-sm text-destructive break-words">{error}</p>
            </div>
          )}

          {/* Generated Draft */}
          {draftReply && (
            <div className="space-y-3 sm:space-y-4">
              <div className="flex items-center justify-between gap-2">
                <Label className="text-sm font-semibold text-foreground flex items-center gap-2 min-w-0">
                  <CheckCircle className="w-4 h-4 text-green-600 dark:text-green-400 flex-shrink-0" />
                  <span className="truncate">{t('emailDetail.aiDraft.generatedDraft')}</span>
                </Label>
                <div className="flex items-center gap-2 text-xs text-muted-foreground whitespace-nowrap flex-shrink-0">
                  {t('emailDetail.aiDraft.readyToUse')}
                </div>
              </div>

              <div className="p-3 sm:p-4 bg-background rounded-xl border shadow-sm overflow-hidden">
                <Textarea
                  value={draftReply}
                  readOnly
                  className="min-h-[150px] sm:min-h-[200px] border-0 resize-none focus-visible:ring-0 p-0 text-foreground leading-relaxed w-full"
                />
              </div>

              {/* Action Buttons */}
              <div className="flex flex-col sm:flex-row gap-2 sm:gap-3 justify-end">
                <Button
                  variant="outline"
                  onClick={() => copyToClipboard(draftReply)}
                  className="flex items-center justify-center gap-2 px-4 py-2 rounded-xl w-full sm:w-auto whitespace-nowrap"
                >
                  {copied ? (
                    <>
                      <CheckCircle className="w-4 h-4 text-green-600 dark:text-green-400 flex-shrink-0" />
                      <span className="truncate">{t('emailDetail.aiDraft.copied')}</span>
                    </>
                  ) : (
                    <>
                      <Copy className="w-4 h-4 flex-shrink-0" />
                      <span className="truncate">{t('emailDetail.aiDraft.copy')}</span>
                    </>
                  )}
                </Button>

                <Button
                  onClick={openInEmailClient}
                  className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-xl flex items-center justify-center gap-2 w-full sm:w-auto whitespace-nowrap"
                >
                  <Mail className="w-4 h-4 flex-shrink-0" />
                  <span className="hidden sm:inline truncate">{t('emailDetail.aiDraft.openInEmail')}</span>
                  <span className="sm:hidden truncate">{t('emailDetail.aiDraft.emailShort')}</span>
                </Button>
              </div>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
