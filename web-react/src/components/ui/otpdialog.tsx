"use client";

import { useEffect, useMemo, useState } from "react";
import { createPortal } from "react-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

type OTPDialogProps = {
  open: boolean;
  title: string;
  description: string;
  telegramUrl?: string;
  expiresAt?: string;
  verifying?: boolean;
  resending?: boolean;
  onClose: () => void;
  onVerify?: (code: string) => Promise<boolean | void> | boolean | void;
  onResend?: () => Promise<boolean | void> | boolean | void;
};

export default function OTPDialog({
  open,
  title,
  description,
  telegramUrl,
  expiresAt,
  verifying = false,
  resending = false,
  onClose,
  onVerify,
  onResend,
}: OTPDialogProps) {
  const [otp, setOtp] = useState(["", "", "", "", "", ""]);
  const [message, setMessage] = useState("");
  const [timeLeft, setTimeLeft] = useState(60);
  const canResend = timeLeft <= 0;
  const complete = useMemo(() => otp.every((digit) => digit !== ""), [otp]);

  useEffect(() => {
    if (!open) return;
    setOtp(["", "", "", "", "", ""]);
    setMessage("");
    setTimeLeft(60);
    window.setTimeout(() => document.getElementById("otp-0")?.focus(), 100);
  }, [open]);

  useEffect(() => {
    if (!open || timeLeft <= 0) return;
    const timer = window.setTimeout(() => setTimeLeft(timeLeft - 1), 1000);
    return () => window.clearTimeout(timer);
  }, [open, timeLeft]);

  if (!open) return null;

  const handleChange = (value: string, index: number) => {
    if (!/^\d?$/.test(value)) return;
    const updated = [...otp];
    updated[index] = value;
    setOtp(updated);
    if (value && index < otp.length - 1) document.getElementById(`otp-${index + 1}`)?.focus();
  };

  const handleVerify = async () => {
    if (!complete) {
      setMessage("Введите полный 6-значный код.");
      return;
    }
    const ok = await onVerify?.(otp.join(""));
    if (ok === false) {
      setMessage("Код не принят. Проверьте Telegram и попробуйте ещё раз.");
      return;
    }
    setMessage("Код принят. Можно продолжать.");
  };

  const handleResend = async () => {
    const ok = await onResend?.();
    if (ok === false) {
      setMessage("Не удалось отправить код повторно. Попробуйте позже.");
      return;
    }
    setMessage("Код отправлен повторно в Telegram.");
    setOtp(["", "", "", "", "", ""]);
    setTimeLeft(60);
    document.getElementById("otp-0")?.focus();
  };

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60).toString().padStart(2, "0");
    const secs = (seconds % 60).toString().padStart(2, "0");
    return `${mins}:${secs}`;
  };

  return createPortal(
    <div className="otp-backdrop-v2" role="dialog" aria-modal="true">
      <section className="otp-dialog-v2">
        <button className="otp-dialog-v2__close" type="button" onClick={onClose} aria-label="Close">
          x
        </button>
        <header className="otp-dialog-v2__header">
          <h2>{title}</h2>
          <p>{description}</p>
        </header>
        <p className="otp-dialog-v2__step">Шаг 1 из 1: подтвердите аккаунт</p>
        {telegramUrl ? (
          <a className="otp-dialog-v2__telegram" href={telegramUrl} target="_blank" rel="noreferrer">
            Открыть Telegram с кодом
          </a>
        ) : null}
        {expiresAt ? <p className="otp-dialog-v2__expires">Код действует до {expiresAt}</p> : null}
        <div className="otp-dialog-v2__inputs">
          {otp.map((digit, index) => (
            <Input
              key={index}
              id={`otp-${index}`}
              value={digit}
              onChange={(event) => handleChange(event.target.value, index)}
              onKeyDown={(event) => {
                if (event.key === "Backspace" && !digit && index > 0) document.getElementById(`otp-${index - 1}`)?.focus();
              }}
              className="h-14 w-14 rounded-md border border-muted-foreground text-center text-lg font-medium focus:border-primary focus:ring-1 focus:ring-primary"
              maxLength={1}
              inputMode="numeric"
            />
          ))}
        </div>
        {!canResend ? (
          <p className="otp-dialog-v2__timer">
            Повторная отправка через <strong>{formatTime(timeLeft)}</strong>
          </p>
        ) : null}
        <div className="otp-dialog-v2__actions">
          <Button className="w-full" onClick={() => void handleVerify()} disabled={verifying}>
            {verifying ? "Проверяем..." : "Подтвердить код"}
          </Button>
          <Button variant="outline" className="w-full justify-between" onClick={() => void handleResend()} disabled={!canResend || resending}>
            {resending ? "Отправляем..." : canResend ? "Отправить снова" : "Повторить отправку"}
            {!canResend ? <span className="text-xs text-muted-foreground">{formatTime(timeLeft)}</span> : null}
          </Button>
        </div>
        {message ? <p className="otp-dialog-v2__message">{message}</p> : null}
      </section>
    </div>,
    document.body,
  );
}
