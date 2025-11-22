import React, { useState } from 'react';
import * as Dialog from '@radix-ui/react-dialog';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Label } from '@/components/ui/Label';
import { api } from '@/lib/api';

interface AIDraftReplyModalProps {
  emailContent: string;
  isOpen: boolean;
  onClose: () => void;
}

export default function AIDraftReplyModal({ emailContent, isOpen, onClose }: AIDraftReplyModalProps) {
  const [userPrompt, setUserPrompt] = useState('');
  const [draftReply, setDraftReply] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleGenerateDraft = async () => {
    setLoading(true);
    setError(null);
    setDraftReply('');
    try {
      const response = await api.post('/ai/draft', { emailContent, userPrompt });
      setDraftReply(response.data.draft);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Failed to generate draft reply.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog.Root open={isOpen} onOpenChange={onClose}>
      <Dialog.Portal>
        <Dialog.Overlay className="bg-black/50 data-[state=open]:animate-overlayShow fixed inset-0" />
        <Dialog.Content className="data-[state=open]:animate-contentShow fixed top-[50%] left-[50%] max-h-[85vh] w-[90vw] max-w-[500px] translate-x-[-50%] translate-y-[-50%] rounded-[6px] bg-white p-[25px] shadow-[hsl(206_22%_7%_/_35%)_0px_10px_38px_-10px,_hsl(206_22%_7%_/_20%)_0px_10px_20px_-15px] focus:outline-none">
          <Dialog.Title className="text-mauve12 m-0 text-[17px] font-semibold">Generate AI Draft Reply</Dialog.Title>
          <Dialog.Description className="text-mauve11 mt-4 mb-5 text-[15px] leading-normal">
            Provide a prompt to generate a draft reply for this email.
          </Dialog.Description>
          <fieldset className="mb-[15px] flex items-center gap-5">
            <Label htmlFor="userPrompt" className="text-grass11 w-[90px] text-right text-[15px]">
              Your Prompt
            </Label>
            <Input
              id="userPrompt"
              value={userPrompt}
              onChange={(e) => setUserPrompt(e.target.value)}
              className="text-violet11 shadow-violet7 focus:shadow-violet8 inline-flex h-[35px] w-full flex-1 items-center justify-center rounded-[4px] px-[10px] text-[15px] leading-none shadow-[0_0_0_1px] outline-none focus:shadow-[0_0_0_2px]"
            />
          </fieldset>
          <div className="mt-[25px] flex justify-end">
            <Button onClick={handleGenerateDraft} disabled={loading}>
              {loading ? 'Generating...' : 'Generate Draft'}
            </Button>
          </div>

          {error && <div className="text-red-500 mt-4">Error: {error}</div>}

          {draftReply && (
            <div className="mt-8 p-4 bg-gray-50 rounded-md border border-gray-200">
              <h3 className="text-md font-semibold mb-2">Draft Reply:</h3>
              <p className="text-gray-700 whitespace-pre-wrap">{draftReply}</p>
            </div>
          )}

          <Dialog.Close asChild>
            <Button className="text-violet11 hover:bg-violet4 focus:shadow-violet7 absolute top-[10px] right-[10px] inline-flex h-[25px] w-[25px] appearance-none items-center justify-center rounded-full focus:shadow-[0_0_0_2px] focus:outline-none" aria-label="Close">
              X
            </Button>
          </Dialog.Close>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
