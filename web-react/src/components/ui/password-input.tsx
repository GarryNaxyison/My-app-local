"use client";

import { useState, type InputHTMLAttributes } from "react";
import { Eye, EyeOff } from "lucide-react";
import { cn } from "@/lib/utils";

type PasswordInputFieldProps = InputHTMLAttributes<HTMLInputElement> & {
  label: string;
  invalid?: boolean;
  valid?: boolean;
  helperText?: string;
  errorText?: string;
};

export function PasswordInputField({
  label,
  invalid,
  valid,
  helperText,
  errorText,
  className,
  ...props
}: PasswordInputFieldProps) {
  const [visible, setVisible] = useState(false);
  return (
    <label className="password-input-v2">
      <span className={cn(valid && "is-valid", invalid && "is-invalid")}>{label}</span>
      <div className={cn("password-input-v2__control", invalid && "is-invalid")}>
        <input
          {...props}
          className={cn("password-input-v2__input", className)}
          type={visible ? "text" : "password"}
        />
        <button
          type="button"
          className="password-input-v2__toggle"
          onClick={() => setVisible((current) => !current)}
          aria-label={visible ? "Hide password" : "Show password"}
        >
          {visible ? <Eye size={16} /> : <EyeOff size={16} />}
        </button>
      </div>
      {invalid && errorText ? <small className="password-input-v2__error">{errorText}</small> : null}
      {!invalid && helperText ? <small className={cn("password-input-v2__helper", valid && "is-valid")}>{helperText}</small> : null}
    </label>
  );
}
