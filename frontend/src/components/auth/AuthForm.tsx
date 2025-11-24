'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link'; // Import Link
import { useLanguage } from '@/lib/i18n/LanguageContext';
import { useAuthStore } from '@/store/auth';
import { cn } from '@/lib/utils';
import { Loader2, Mail, Lock, User, AlertCircle, Eye, EyeOff } from 'lucide-react';
import { isAxiosError } from 'axios'; // Import isAxiosError

interface AuthInputProps {
  id: string;
  label: string;
  type: string;
  placeholder: string;
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  error?: string;
  icon: React.ElementType;
  showPasswordToggle?: boolean;
}

function AuthInput({ id, label, type, placeholder, value, onChange, error, icon: Icon, showPasswordToggle }: AuthInputProps) {
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
          className={cn(
            "w-full pl-10 pr-3 py-2 border rounded-lg text-sm transition-all duration-200",
            "focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500",
            error ? "border-red-500 ring-red-100" : "border-slate-200 bg-slate-50"
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

interface AuthFormProps {
  mode: 'login' | 'register';
  onModeChange: (mode: 'login' | 'register') => void;
}

export function AuthForm({ mode, onModeChange }: AuthFormProps) {
  const { t } = useLanguage();
  const router = useRouter();
  const loginUser = useAuthStore(state => state.login);
  const registerUser = useAuthStore(state => state.register);
  const user = useAuthStore(state => state.user);

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState(''); // Only for register
  const [fullName, setFullName] = useState(''); // Only for register
  const [isLoading, setIsLoading] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    // If already authenticated, redirect to dashboard
    if (user) {
      router.push('/dashboard');
    }
  }, [user, router]);

  const validateForm = () => {
    const newErrors: Record<string, string> = {};
    if (!email) newErrors.email = t('auth.errors.emailRequired');
    else if (!/^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/.test(email)) newErrors.email = t('auth.errors.emailInvalid');

    if (!password) newErrors.password = t('auth.errors.passwordLength'); // Simplified for now
    else if (password.length < 8) newErrors.password = t('auth.errors.passwordLength');

    if (mode === 'register') {
      if (!fullName) newErrors.fullName = t('auth.errors.nameRequired'); // Add to dict
      if (password !== confirmPassword) newErrors.confirmPassword = t('auth.errors.passwordMismatch');
    }
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validateForm()) return;

    setIsLoading(true);
    setErrors({});

    try {
      if (mode === 'login') {
        await loginUser(email, password);
      } else {
        await registerUser(fullName, email, password);
      }
    } catch (error: unknown) {
      console.error("Auth failed:", error);
      let errorMessage = "";
      if (isAxiosError(error) && error.response?.data?.error) {
        errorMessage = error.response.data.error;
      } else if (error instanceof Error) {
        errorMessage = error.message;
      } else {
        errorMessage = "An unknown error occurred.";
      }

      if (errorMessage.includes("invalid credentials") || errorMessage.includes("login failed")) {
        setErrors({ general: t('auth.errors.invalidCredentials') });
      } else if (errorMessage.includes("user already exists") || errorMessage.includes("duplicate key")) {
        setErrors({ email: t('auth.errors.emailTaken') });
      } else {
        setErrors({ general: t('auth.errors.registrationFailed') });
      }
    } finally {
      setIsLoading(false);
    }
  };

  const toggleMode = (newMode: 'login' | 'register') => {
    onModeChange(newMode);
    setErrors({}); // Clear errors on mode switch
  };

  // Dynamic title based on mode
  const title = mode === 'login' ? t('auth.welcomeBack') : t('auth.createAccount');
  const buttonText = mode === 'login' ? t('auth.continue') : t('auth.createAccount'); // Using createAccount for register button for now

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <h2 className="text-3xl font-bold text-slate-800 text-center mb-6">{title}</h2>

      {errors.general && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm flex items-center gap-2">
          <AlertCircle className="w-5 h-5 flex-shrink-0" />
          {errors.general}
        </div>
      )}

      {mode === 'register' && (
        <AuthInput
          id="fullName"
          label={t('auth.nameLabel')}
          type="text"
          placeholder={t('auth.namePlaceholder')}
          value={fullName}
          onChange={(e) => setFullName(e.target.value)}
          error={errors.fullName}
          icon={User}
        />
      )}

      <AuthInput
        id="email"
        label={t('auth.emailLabel')}
        type="email"
        placeholder={t('auth.emailPlaceholder')}
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        error={errors.email}
        icon={Mail}
      />

      <AuthInput
        id="password"
        label={t('auth.passwordLabel')}
        type="password"
        placeholder={t('auth.passwordPlaceholder')}
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        error={errors.password}
        icon={Lock}
        showPasswordToggle
      />

      {mode === 'register' && (
        <AuthInput
          id="confirmPassword"
          label={t('auth.confirmPasswordLabel')}
          type="password"
          placeholder={t('auth.confirmPasswordLabel')}
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          error={errors.confirmPassword}
          icon={Lock}
          showPasswordToggle
        />
      )}

      <button
        type="submit"
        disabled={isLoading}
        className={cn(
          "w-full flex items-center justify-center gap-2 px-5 py-3 rounded-lg text-white font-semibold shadow-md transition-all duration-200",
          "bg-blue-600 hover:bg-blue-700 disabled:opacity-60 disabled:cursor-not-allowed",
          isLoading && "bg-blue-500"
        )}
      >
        {isLoading && <Loader2 className="w-5 h-5 animate-spin" />}
        {isLoading ? (mode === 'login' ? t('auth.loggingIn') : t('auth.signingUp')) : buttonText}
      </button>

      {/* Forgot Password (Login Mode Only) */}
      {mode === 'login' && (
        <div className="text-center text-sm mt-4">
          <Link href="/auth/forgot-password" className="text-blue-600 hover:underline">
            {t('auth.forgotPassword')}
          </Link>
        </div>
      )}

      {/* Switch between Login/Register */}
      <div className="text-center text-sm text-slate-600 mt-6">
        {mode === 'login' ? t('auth.dontHaveAccount') : t('auth.alreadyHaveAccount')}{' '}
        <button
          type="button"
          onClick={() => toggleMode(mode === 'login' ? 'register' : 'login')}
          className="text-blue-600 hover:underline font-medium"
        >
          {mode === 'login' ? t('auth.createAccount') : t('auth.signInTitle')}
        </button>
      </div>

      {/* Terms of Service */}
      <p className="text-xs text-slate-500 text-center mt-4">
        {t('auth.terms')}
      </p>
    </form>
  );
}
