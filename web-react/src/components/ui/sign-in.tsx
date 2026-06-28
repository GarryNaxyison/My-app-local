import React, { useState } from "react";
import { Eye, EyeOff } from "lucide-react";
import { AuthGenerativeScene } from "./auth-generative-scene";

const GoogleIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 48 48" aria-hidden="true">
    <path fill="#FFC107" d="M43.611 20.083H42V20H24v8h11.303c-1.649 4.657-6.08 8-11.303 8-6.627 0-12-5.373-12-12s12-5.373 12-12c3.059 0 5.842 1.154 7.961 3.039l5.657-5.657C34.046 6.053 29.268 4 24 4 12.955 4 4 12.955 4 24s8.955 20 20 20 20-8.955 20-20c0-2.641-.21-5.236-.611-7.743z" />
    <path fill="#FF3D00" d="M6.306 14.691l6.571 4.819C14.655 15.108 18.961 12 24 12c3.059 0 5.842 1.154 7.961 3.039l5.657-5.657C34.046 6.053 29.268 4 24 4 16.318 4 9.656 8.337 6.306 14.691z" />
    <path fill="#4CAF50" d="M24 44c5.166 0 9.86-1.977 13.409-5.192l-6.19-5.238C29.211 35.091 26.715 36 24 36c-5.202 0-9.619-3.317-11.283-7.946l-6.522 5.025C9.505 39.556 16.227 44 24 44z" />
    <path fill="#1976D2" d="M43.611 20.083H42V20H24v8h11.303c-.792 2.237-2.231 4.166-4.087 5.571l6.19 5.238C42.022 35.026 44 30.038 44 24c0-2.641-.21-5.236-.611-7.743z" />
  </svg>
);

export interface Testimonial {
  avatarSrc: string;
  name: string;
  handle: string;
  text: string;
}

export type SignInMode = "login" | "register";

interface SignInPageProps {
  mode?: SignInMode;
  title?: React.ReactNode;
  description?: React.ReactNode;
  heroImageSrc?: string;
  testimonials?: Testimonial[];
  status?: React.ReactNode;
  captchaSlot?: React.ReactNode;
  extraFields?: React.ReactNode;
  submitting?: boolean;
  labels?: Partial<{
    login: string;
    loginPlaceholder: string;
    password: string;
    passwordPlaceholder: string;
    passwordHint: string;
    passwordConfirm: string;
    passwordConfirmPlaceholder: string;
    passwordConfirmHint: string;
    remember: string;
    resetPassword: string;
    submit: string;
    divider: string;
    social: string;
    createPrompt: string;
    createAction: string;
    loginPrompt: string;
    loginAction: string;
    referral: string;
    referralPlaceholder: string;
    hidePassword: string;
    showPassword: string;
  }>;
  socialIcon?: React.ReactNode;
  onSignIn?: (event: React.FormEvent<HTMLFormElement>) => void;
  onGoogleSignIn?: () => void;
  onResetPassword?: () => void;
  onCreateAccount?: () => void;
  onLoginMode?: () => void;
}

const GlassInputWrapper = ({ children }: { children: React.ReactNode }) => (
  <div className="rounded-2xl border border-[color:var(--line)] bg-[color-mix(in_srgb,var(--surface-strong)_64%,transparent)] backdrop-blur-sm transition-colors focus-within:border-violet-400/70 focus-within:bg-violet-500/10">
    {children}
  </div>
);

export const SignInPage: React.FC<SignInPageProps> = ({
  mode = "login",
  title = <span className="font-light tracking-tighter text-[color:var(--text)]">Welcome</span>,
  description = "Access your account and continue your journey with us",
  heroImageSrc,
  testimonials = [],
  status,
  captchaSlot,
  extraFields,
  submitting = false,
  labels = {},
  socialIcon,
  onSignIn,
  onGoogleSignIn,
  onResetPassword,
  onCreateAccount,
  onLoginMode,
}) => {
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const isRegister = mode === "register";
  const hasSocialSignIn = Boolean(onGoogleSignIn);
  void heroImageSrc;
  void testimonials;

  return (
    <div className="sign-in-page-v2 bg-[color:var(--bg)] text-[color:var(--text)]">
      <AuthGenerativeScene className="sign-in-page-v2__background-scene" />
      <section className="sign-in-page-v2__form-section flex items-center justify-center">
        <div className="sign-in-page-v2__form-card w-full max-w-md">
          <div className="flex flex-col gap-6">
            <h1 className="sign-in-page-v2__title animate-element animate-delay-100 font-semibold">{title}</h1>
            <p className="animate-element animate-delay-200 text-[color:var(--muted)]">{description}</p>
            {status ? <div className="animate-element animate-delay-200 auth-status-v2">{status}</div> : null}

            <form className="space-y-5" onSubmit={onSignIn}>
              <div className="animate-element animate-delay-300">
                <label className="text-sm font-medium text-[color:var(--muted)]">{labels.login || "Email Address"}</label>
                <GlassInputWrapper>
                  <input name="login" type="text" autoComplete="username" placeholder={labels.loginPlaceholder || "Enter your email address"} className="w-full rounded-2xl border-0 bg-transparent p-4 text-sm focus:outline-none focus:ring-0" />
                </GlassInputWrapper>
              </div>

              <div className="animate-element animate-delay-400">
                <label className="text-sm font-medium text-[color:var(--muted)]">{labels.password || "Password"}</label>
                <GlassInputWrapper>
                  <div className="relative">
                    <input name="password" type={showPassword ? "text" : "password"} autoComplete={isRegister ? "new-password" : "current-password"} placeholder={labels.passwordPlaceholder || "Enter your password"} className="w-full rounded-2xl border-0 bg-transparent p-4 pr-12 text-sm focus:outline-none focus:ring-0" />
                    <button type="button" onClick={() => setShowPassword(!showPassword)} className="absolute inset-y-0 right-3 flex items-center" aria-label={showPassword ? labels.hidePassword || "Hide password" : labels.showPassword || "Show password"}>
                      {showPassword ? <EyeOff className="h-5 w-5 text-[color:var(--muted)] transition-colors hover:text-[color:var(--text)]" /> : <Eye className="h-5 w-5 text-[color:var(--muted)] transition-colors hover:text-[color:var(--text)]" />}
                    </button>
                  </div>
                </GlassInputWrapper>
                {isRegister && labels.passwordHint ? <small className="mt-2 block text-xs font-medium text-[color:var(--muted)]">{labels.passwordHint}</small> : null}
              </div>

              {isRegister ? (
                <div className="animate-element animate-delay-500">
                  <label className="text-sm font-medium text-[color:var(--muted)]">{labels.passwordConfirm || "Repeat password"}</label>
                  <GlassInputWrapper>
                    <div className="relative">
                      <input name="password_confirm" type={showConfirmPassword ? "text" : "password"} autoComplete="new-password" placeholder={labels.passwordConfirmPlaceholder || "Repeat password"} className="w-full rounded-2xl border-0 bg-transparent p-4 pr-12 text-sm focus:outline-none focus:ring-0" />
                      <button type="button" onClick={() => setShowConfirmPassword(!showConfirmPassword)} className="absolute inset-y-0 right-3 flex items-center" aria-label={showConfirmPassword ? labels.hidePassword || "Hide password" : labels.showPassword || "Show password"}>
                        {showConfirmPassword ? <EyeOff className="h-5 w-5 text-[color:var(--muted)] transition-colors hover:text-[color:var(--text)]" /> : <Eye className="h-5 w-5 text-[color:var(--muted)] transition-colors hover:text-[color:var(--text)]" />}
                      </button>
                    </div>
                  </GlassInputWrapper>
                  {labels.passwordConfirmHint ? <small className="mt-2 block text-xs font-medium text-[color:var(--muted)]">{labels.passwordConfirmHint}</small> : null}
                </div>
              ) : null}

              {isRegister ? (
                <div className="animate-element animate-delay-500">
                  <label className="text-sm font-medium text-[color:var(--muted)]">{labels.referral || "Referral code"}</label>
                  <GlassInputWrapper>
                    <input name="referral_code" type="text" autoComplete="off" placeholder={labels.referralPlaceholder || "Optional"} className="w-full rounded-2xl border-0 bg-transparent p-4 text-sm focus:outline-none focus:ring-0" />
                  </GlassInputWrapper>
                </div>
              ) : null}

              {extraFields}
              {captchaSlot}

              {!isRegister ? (
                <div className="animate-element animate-delay-500 flex items-center justify-between gap-3 text-sm">
                  <label className="flex cursor-pointer items-center gap-3">
                    <input type="checkbox" name="rememberMe" className="custom-checkbox" />
                    <span className="text-[color:color-mix(in_srgb,var(--text)_90%,transparent)]">{labels.remember || "Keep me signed in"}</span>
                  </label>
                  <a href="#" onClick={(event) => { event.preventDefault(); onResetPassword?.(); }} className="text-violet-400 transition-colors hover:underline">{labels.resetPassword || "Reset password"}</a>
                </div>
              ) : null}

              <button type="submit" disabled={submitting} className="animate-element animate-delay-600 w-full rounded-2xl bg-[color:var(--primary)] py-4 font-medium text-white transition-colors hover:bg-[color-mix(in_srgb,var(--primary)_90%,black)] disabled:cursor-not-allowed disabled:opacity-60">
                {submitting ? "..." : labels.submit || "Sign In"}
              </button>
            </form>

            {hasSocialSignIn ? (
              <>
                <div className="animate-element animate-delay-700 relative flex items-center justify-center">
                  <span className="w-full border-t border-[color:var(--line)]"></span>
                  <span className="absolute bg-[color:var(--bg)] px-4 text-sm text-[color:var(--muted)]">{labels.divider || "Or continue with"}</span>
                </div>

                <button type="button" onClick={onGoogleSignIn} className="auth-social-button-v2 animate-element animate-delay-800 flex w-full items-center justify-center gap-3 rounded-2xl border border-[color:var(--line)] py-3 transition-colors hover:bg-[color:var(--surface-muted)]">
                  {socialIcon || <GoogleIcon />}
                  {labels.social || "Continue with Google"}
                </button>
              </>
            ) : null}

            <p className="animate-element animate-delay-900 text-center text-sm text-[color:var(--muted)]">
              {isRegister ? labels.loginPrompt || "Already have an account?" : labels.createPrompt || "New to our platform?"}{" "}
              <a href="#" onClick={(event) => { event.preventDefault(); isRegister ? onLoginMode?.() : onCreateAccount?.(); }} className="text-violet-400 transition-colors hover:underline">
                {isRegister ? labels.loginAction || "Sign In" : labels.createAction || "Create Account"}
              </a>
            </p>
          </div>
        </div>
      </section>

    </div>
  );
};

export default SignInPage;
