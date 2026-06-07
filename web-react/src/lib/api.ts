import type { SessionData } from "./types";

export class ApiError extends Error {
  status: number;
  payload: unknown;

  constructor(message: string, status: number, payload: unknown) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.payload = payload;
  }
}

type RequestOptions = Omit<RequestInit, "body"> & {
  body?: unknown;
};

const API_BASE = (import.meta.env.VITE_API_BASE || "").replace(/\/$/, "");

function errorMessage(payload: unknown, fallback: string) {
  if (payload && typeof payload === "object") {
    const record = payload as Record<string, unknown>;
    const error = record.error;
    if (typeof error === "string") return error;
    if (error && typeof error === "object") {
      const message = (error as Record<string, unknown>).message;
      if (typeof message === "string") return message;
    }
    if (typeof record.message === "string") return record.message;
  }
  return fallback;
}

export async function api<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const headers = new Headers(options.headers);
  let body: BodyInit | undefined;

  if (options.body !== undefined) {
    headers.set("Content-Type", "application/json");
    body = JSON.stringify(options.body);
  }

  const response = await fetch(API_BASE + path, {
    ...options,
    headers,
    body,
    credentials: "include",
  });

  const text = await response.text();
  let payload: unknown = null;
  if (text) {
    try {
      payload = JSON.parse(text);
    } catch {
      payload = text;
    }
  }

  if (!response.ok) {
    throw new ApiError(errorMessage(payload, response.statusText), response.status, payload);
  }

  return payload as T;
}

export async function apiForm<T>(path: string, formData: FormData): Promise<T> {
  const response = await fetch(API_BASE + path, {
    method: "POST",
    body: formData,
    credentials: "include",
  });
  const payload = await response.json().catch(() => null);
  if (!response.ok) {
    throw new ApiError(errorMessage(payload, response.statusText), response.status, payload);
  }
  return payload as T;
}

export async function apiBlob(path: string, body: unknown): Promise<Blob> {
  const response = await fetch(API_BASE + path, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
    credentials: "include",
  });
  if (!response.ok) {
    const payload = await response.json().catch(() => null);
    throw new ApiError(errorMessage(payload, response.statusText), response.status, payload);
  }
  return response.blob();
}

export function loadSession() {
  return api<SessionData>("/api/session");
}

export function logout() {
  return api<{ ok: boolean }>("/api/auth/logout", { method: "POST", body: {} });
}
