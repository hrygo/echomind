import React, { useState } from 'react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '@/components/ui/Dialog';
import { Button } from '@/components/ui/Button';
import { Label } from '@/components/ui/Label';
import { Textarea } from '@/components/ui/Textarea';
import { api } from '@/lib/api';
import { Loader2, Sparkles, Copy, Mail, CheckCircle } from 'lucide-react';

interface AIDraftReplyModalProps {
  emailContent: string;
  isOpen: boolean;
  onClose: () => void;
}

export default function AIDraftReplyModal({ emailContent, isOpen, onClose }: AIDraftReplyModalProps) {
  const [userPrompt, setUserPrompt] = useState('Generate a professional email reply to this message.');
  const [draftReply, setDraftReply] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  const handleGenerateDraft = async () => {
    if (!emailContent.trim()) {
      setError('Email content is required.');
      return;
    }
    if (!userPrompt.trim()) {
      setError('User prompt is required.');
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
      setError(err.response?.data?.error || err.message || 'Failed to generate draft reply.');
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
      <DialogContent className="max-w-4xl w-full max-h-[95vh] mx-4 sm:mx-auto overflow-y-auto">
        <DialogHeader className="pb-4 sm:pb-6">
          <div className="flex items-center gap-3 sm:gap-4">
            <div className="w-10 h-10 sm:w-12 sm:h-12 bg-blue-100 dark:bg-blue-900/30 rounded-xl sm:rounded-xl flex items-center justify-center">
              <Sparkles className="w-5 h-5 sm:w-6 sm:h-6 text-blue-600 dark:text-blue-400" />
            </div>
            <div className="flex-1 min-w-0">
              <DialogTitle className="text-lg sm:text-xl font-bold text-foreground">
                AI Email Draft Generator
              </DialogTitle>
              <DialogDescription className="text-sm sm:text-base text-muted-foreground mt-1">
                Create a professional email reply powered by AI
              </DialogDescription>
            </div>
          </div>
        </DialogHeader>

        <div className="space-y-4 sm:space-y-6">
          {/* Email Preview */}
          <div className="space-y-2 sm:space-y-3">
            <Label className="text-sm font-semibold text-foreground flex items-center gap-2">
              <Mail className="w-4 h-4" />
              Original Email
            </Label>
            <div className="p-3 sm:p-4 bg-muted/50 rounded-xl border max-h-32 sm:max-h-48 overflow-y-auto">
              <p className="text-sm text-muted-foreground whitespace-pre-wrap leading-relaxed">
                {emailContent.trim() || 'No email content provided'}
              </p>
            </div>
          </div>

          {/* Prompt Input */}
          <div className="space-y-2 sm:space-y-3">
            <Label htmlFor="userPrompt" className="text-sm font-semibold text-foreground">
              Instructions for AI
            </Label>
            <Textarea
              id="userPrompt"
              value={userPrompt}
              onChange={(e) => setUserPrompt(e.target.value)}
              placeholder="Describe how you'd like to reply to this email..."
              className="min-h-[80px] sm:min-h-[100px] resize-none"
            />
            <div className="flex flex-wrap gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setUserPrompt('Generate a brief professional response.')}
                className="text-xs h-8 px-3"
              >
                Professional
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setUserPrompt('Generate a friendly and casual reply.')}
                className="text-xs h-8 px-3"
              >
                Casual
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setUserPrompt('Generate a concise and direct response.')}
                className="text-xs h-8 px-3"
              >
                Concise
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setUserPrompt('Generate a detailed and comprehensive reply.')}
                className="text-xs h-8 px-3"
              >
                Detailed
              </Button>
            </div>
          </div>

          {/* Generate Button */}
          <div className="flex justify-center">
            <Button
              onClick={handleGenerateDraft}
              disabled={loading}
              className="w-full sm:w-auto bg-blue-600 hover:bg-blue-700 text-white px-6 sm:px-8 py-3 rounded-xl font-medium shadow-lg shadow-blue-600/20 flex items-center gap-2 min-h-[44px] sm:min-h-[48px]"
            >
              {loading ? (
                <>
                  <Loader2 className="w-4 h-4 animate-spin" />
                  Generating AI Draft...
                </>
              ) : (
                <>
                  <Sparkles className="w-4 h-4" />
                  <span className="hidden sm:inline">Generate AI Draft</span>
                  <span className="sm:hidden">Generate Draft</span>
                </>
              )}
            </Button>
          </div>

          {/* Error Display */}
          {error && (
            <div className="p-3 sm:p-4 bg-destructive/10 border border-destructive/20 rounded-xl">
              <p className="text-sm text-destructive">{error}</p>
            </div>
          )}

          {/* Generated Draft */}
          {draftReply && (
            <div className="space-y-3 sm:space-y-4">
              <div className="flex items-center justify-between">
                <Label className="text-sm font-semibold text-foreground flex items-center gap-2">
                  <CheckCircle className="w-4 h-4 text-green-600 dark:text-green-400" />
                  AI Generated Draft
                </Label>
                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                  Ready to use
                </div>
              </div>

              <div className="p-3 sm:p-4 bg-background rounded-xl border shadow-sm">
                <Textarea
                  value={draftReply}
                  readOnly
                  className="min-h-[150px] sm:min-h-[200px] border-0 resize-none focus-visible:ring-0 p-0 text-foreground leading-relaxed"
                />
              </div>

              {/* Action Buttons */}
              <div className="flex flex-col sm:flex-row gap-2 sm:gap-3 justify-end">
                <Button
                  variant="outline"
                  onClick={() => copyToClipboard(draftReply)}
                  className="flex items-center justify-center gap-2 px-4 py-2 rounded-xl w-full sm:w-auto"
                >
                  {copied ? (
                    <>
                      <CheckCircle className="w-4 h-4 text-green-600 dark:text-green-400" />
                      Copied!
                    </>
                  ) : (
                    <>
                      <Copy className="w-4 h-4" />
                      Copy
                    </>
                  )}
                </Button>

                <Button
                  onClick={openInEmailClient}
                  className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-xl flex items-center justify-center gap-2 w-full sm:w-auto"
                >
                  <Mail className="w-4 h-4" />
                  <span className="hidden sm:inline">Open in Email</span>
                  <span className="sm:hidden">Email</span>
                </Button>
              </div>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
