import { useMutation } from '@tanstack/react-query';
import api from '@/lib/api';

// AI Reply request interface
export interface AIReplyRequest {
  emailId: string;
  tone?: 'professional' | 'casual' | 'friendly' | 'concise';
  context?: 'brief' | 'detailed' | 'urgent';
}

// AI Reply response interface
export interface AIReplyResponse {
  reply: string;
  confidence: number;
}

export const useAIReply = () => {
  return useMutation({
    mutationFn: async (requestData: AIReplyRequest): Promise<AIReplyResponse> => {
      const { data: responseData } = await api.post('/ai/reply', requestData);
      return responseData;
    },
    onSuccess: (responseData) => {
      console.log('AI Reply generated successfully:', {
        length: responseData.reply.length,
        confidence: responseData.confidence
      });
    },
    onError: (error) => {
      console.error('Failed to generate AI reply:', error);
    },
  });
};