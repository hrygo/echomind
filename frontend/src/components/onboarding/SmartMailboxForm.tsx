'use client';

import React, { useState, useEffect, useMemo } from 'react';
import { Mail, Lock, Server, Globe, Key, Wifi, AlertCircle, Eye, EyeOff, CheckCircle2, Loader2 } from 'lucide-react';
import { useLanguage } from '@/lib/i18n/LanguageContext';
import { useOnboardingStore } from '@/store';
import { cn } from '@/lib/utils';
import { detectProvider } from '@/lib/constants/mail_providers';
import { api } from '@/lib/api';
import { isAxiosError } from 'axios';

interface InputFieldProps {
  id: string;
  label: string;
  type: string;
  placeholder: string;
  value: string | number;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  error?: string;
  icon: React.ElementType;
  disabled?: boolean;
  showPasswordToggle?: boolean;
}

function InputField({ id, label, type, placeholder, value, onChange, error, icon: Icon, disabled, showPasswordToggle }: InputFieldProps) {
  const [showPassword, setShowPassword] = useState(false);
  const inputType = showPasswordToggle && showPassword ? 'text' : type;

  return (
    <div className="space-y-1">
      <label htmlFor={id} className="block text-sm font-medium text-slate-700">
        {label}
      </label>
      <div className="relative">
        <span className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400">
          <Icon className="w-5 h-5" />
        </span>
        <input
          id={id}
          type={inputType}
          placeholder={placeholder}
          value={value}
          onChange={onChange}
          disabled={disabled}
          className={cn(
            "w-full pl-10 pr-3 py-2 border rounded-lg text-sm transition-all duration-200",
            "focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500",
            error ? "border-red-500 ring-red-100" : "border-slate-200 bg-slate-50",
            disabled && "cursor-not-allowed bg-slate-100"
          )}
        />
        {showPasswordToggle && (
          <span 
            className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 cursor-pointer hover:text-slate-600"
            onClick={() => setShowPassword(prev => !prev)}
          >
            {showPassword ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
          </span>
        )}
      </div>
      {error && <p className="text-sm text-red-500 mt-1">{error}</p>}
    </div>
  );
}

export function SmartMailboxForm() {
  const { t } = useLanguage();
  const { mailbox, setMailboxConfig, setStep } = useOnboardingStore();

  const [showAdvanced, setShowAdvanced] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [connectionStatus, setConnectionStatus] = useState<'idle' | 'detecting' | 'connecting' | 'success' | 'failed'>('idle');
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [generalError, setGeneralError] = useState<string | null>(null);

  const detectedProvider = useMemo(() => {
    return detectProvider(mailbox.email);
  }, [mailbox.email]);

  useEffect(() => {
    // Auto-fill advanced settings if provider is detected and not already manually set
    if (detectedProvider && !showAdvanced) {
      setMailboxConfig({
        imapServer: detectedProvider.imap.host,
        imapPort: detectedProvider.imap.port,
        smtpServer: detectedProvider.smtp.host,
        smtpPort: detectedProvider.smtp.port,
        providerConfig: detectedProvider,
      });
    } else if (!detectedProvider && !showAdvanced) {
        // Clear auto-filled if email changes to unknown domain and not in advanced mode
        setMailboxConfig({
            imapServer: '',
            imapPort: 0,
            smtpServer: '',
            smtpPort: 0,
            providerConfig: null,
        });
    }
  }, [detectedProvider, setMailboxConfig, showAdvanced]);

  const validateForm = () => {
    const newErrors: Record<string, string> = {};
    setGeneralError(null);

    if (!mailbox.email) newErrors.email = t('auth.errors.emailRequired');
    else if (!/^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/.test(mailbox.email)) newErrors.email = t('auth.errors.emailInvalid');

    if (!mailbox.password) newErrors.password = t('onboarding.step2.errors.passwordRequired'); // Need to add this to dict

    if (showAdvanced || !detectedProvider) {
      if (!mailbox.imapServer) newErrors.imapServer = t('onboarding.step2.errors.imapHostRequired');
      if (!mailbox.imapPort) newErrors.imapPort = t('onboarding.step2.errors.imapPortRequired');
      else if (isNaN(mailbox.imapPort) || mailbox.imapPort <= 0 || mailbox.imapPort > 65535) newErrors.imapPort = t('onboarding.step2.errors.invalidPort');
      
      if (!mailbox.smtpServer) newErrors.smtpServer = t('onboarding.step2.errors.smtpHostRequired');
      if (!mailbox.smtpPort) newErrors.smtpPort = t('onboarding.step2.errors.smtpPortRequired');
      else if (isNaN(mailbox.smtpPort) || mailbox.smtpPort <= 0 || mailbox.smtpPort > 65535) newErrors.smtpPort = t('onboarding.step2.errors.invalidPort');
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validateForm()) return;

    setIsLoading(true);
    setConnectionStatus('connecting');
    setErrors({});
    setGeneralError(null);

    try {
      // Backend will attempt to connect and save
      await api.post('/settings/account', {
        email: mailbox.email,
        username: mailbox.email, // Assuming email is username
        password: mailbox.password,
        server_address: mailbox.imapServer,
        server_port: mailbox.imapPort,
        // The backend `EmailAccountInput` only has one server/port. This needs clarification.
        // For now, we'll send IMAP details, assuming backend handles SMTP separately or doesn't need it at this stage.
        // TODO: Backend `EmailAccountInput` needs IMAP and SMTP fields separately.
      });
      setConnectionStatus('success');
      setStep(3); // Move to next step
    } catch (error: unknown) {
      console.error("Mailbox connection failed:", error);
      setConnectionStatus('failed');
      if (isAxiosError(error) && error.response) {
        const backendError = error.response.data?.error || t('onboarding.step2.errors.unknownConnection');
        setGeneralError(backendError);
        // More granular error mapping can be done here based on backend error codes
      } else {
        setGeneralError(t('onboarding.step2.errors.networkError'));
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen-auth p-4">
      <h1 className="text-4xl font-extrabold text-slate-900 mb-4 text-center">
        {t('onboarding.step2.title')}
      </h1>
      <p className="text-lg text-slate-600 mb-10 text-center max-w-xl">
        {t('onboarding.step2.subtitle')}
      </p>

      <form onSubmit={handleSubmit} className="w-full max-w-md space-y-6">
        {generalError && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm flex items-center gap-2">
            <AlertCircle className="w-5 h-5 flex-shrink-0" />
            {generalError}
          </div>
        )}

        <InputField
          id="email"
          label={t('onboarding.step2.emailLabel')}
          type="email"
          placeholder={t('auth.emailPlaceholder')}
          value={mailbox.email}
          onChange={(e) => setMailboxConfig({ email: e.target.value })}
          error={errors.email}
          icon={Mail}
        />

        {detectedProvider?.requiresAppPassword && (
          <p className="text-sm text-slate-500 mt-2 flex items-center gap-2">
            <Key className="w-4 h-4" />
            {t('onboarding.step2.passwordHint').replace('{provider}', detectedProvider.name)}
            {detectedProvider.helpLink && (
              <a href={detectedProvider.helpLink} target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline">
                ({t('common.howTo')}) {/* Need to add common.howTo to dict */}
              </a>
            )}
          </p>
        )}

        <InputField
          id="password"
          label={t('onboarding.step2.passwordLabel')}
          type="password"
          placeholder={t('auth.passwordPlaceholder')}
          value={mailbox.password}
          onChange={(e) => setMailboxConfig({ password: e.target.value })}
          error={errors.password}
          icon={Lock}
          showPasswordToggle
        />

        <div className="space-y-4">
          <button
            type="button"
            onClick={() => setShowAdvanced(prev => !prev)}
            className="text-blue-600 hover:underline text-sm flex items-center gap-2"
          >
            <Server className="w-4 h-4" /> {t('onboarding.step2.advancedSettings')}
          </button>

          {(showAdvanced || !detectedProvider) && (
            <div className="space-y-4 bg-slate-50 p-4 rounded-lg border border-slate-100">
              <InputField
                id="imapHost"
                label={t('onboarding.step2.imapHost')}
                type="text"
                placeholder="e.g. imap.example.com"
                value={mailbox.imapServer}
                onChange={(e) => setMailboxConfig({ imapServer: e.target.value })}
                error={errors.imapServer}
                icon={Wifi}
              />
              <InputField
                id="imapPort"
                label={t('onboarding.step2.imapPort')}
                type="number"
                placeholder="e.g. 993"
                value={mailbox.imapPort === 0 ? '' : mailbox.imapPort} // Display empty for 0
                onChange={(e) => setMailboxConfig({ imapPort: parseInt(e.target.value) || 0 })}
                error={errors.imapPort}
                icon={Globe}
              />
              <InputField
                id="smtpHost"
                label={t('onboarding.step2.smtpHost')}
                type="text"
                placeholder="e.g. smtp.example.com"
                value={mailbox.smtpServer}
                onChange={(e) => setMailboxConfig({ smtpServer: e.target.value })}
                error={errors.smtpServer}
                icon={Wifi}
              />
              <InputField
                id="smtpPort"
                label={t('onboarding.step2.smtpPort')}
                type="number"
                placeholder="e.g. 587"
                value={mailbox.smtpPort === 0 ? '' : mailbox.smtpPort} // Display empty for 0
                onChange={(e) => setMailboxConfig({ smtpPort: parseInt(e.target.value) || 0 })}
                error={errors.smtpPort}
                icon={Globe}
              />
            </div>
          )}
        </div>

        <button
          type="submit"
          disabled={isLoading}
          className={cn(
            "w-full flex items-center justify-center gap-2 px-5 py-3 rounded-lg text-white font-semibold shadow-md transition-all duration-200",
            "bg-blue-600 hover:bg-blue-700 disabled:opacity-60 disabled:cursor-not-allowed"
          )}
        >
          {isLoading ? <Loader2 className="w-5 h-5 animate-spin" /> : <CheckCircle2 className="w-5 h-5" />}
          {isLoading ? t('onboarding.step2.connecting') : t('onboarding.step2.connect')}
        </button>

        {connectionStatus === 'success' && (
          <p className="text-center text-green-600 text-sm mt-4 flex items-center justify-center gap-2">
            <CheckCircle2 className="w-5 h-5" /> {t('onboarding.step2.success')}
          </p>
        )}
      </form>
    </div>
  );
}
