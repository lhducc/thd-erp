import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export const formatTime = (timeString: string) => {
    if (!timeString) return ''
  // Convert "HH:mm:ss" to "HH:mm"
  return timeString.split(':').slice(0, 2).join(':');
};

export function formatDate(isoDate: string): string {
  const date = new Date(isoDate);
  const day = String(date.getUTCDate()).padStart(2, '0');
  const month = String(date.getUTCMonth() + 1).padStart(2, '0');
  const year = date.getUTCFullYear();
  return `${day}/${month}/${year}`;
}

export const getTimeFromTimestamp = (timestamp: string): string => {
  const date = new Date(timestamp);
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  return `${hours}:${minutes}`;
};

export const toRFC3339 = (date: Date) => {
  const pad = (n: number) => n.toString().padStart(2, '0');
  const year = date.getFullYear();
  const month = pad(date.getMonth() + 1);
  const day = pad(date.getDate());
  const hours = pad(date.getHours());
  const minutes = pad(date.getMinutes());
  const seconds = pad(date.getSeconds());

  const tzOffset = -date.getTimezoneOffset();
  const sign = tzOffset >= 0 ? '+' : '-';
  const tzHours = pad(Math.floor(Math.abs(tzOffset) / 60));
  const tzMinutes = pad(Math.abs(tzOffset) % 60);

  return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}${sign}${tzHours}:${tzMinutes}`;
};

export function unwrap<T extends Record<string, any>>(obj: T): T[keyof T] {
    const key = Object.keys(obj)[0] as keyof T;
    return obj[key];
}

export function toVietnamISOString(date: Date = new Date()): string {
    const offsetMs = 7 * 60 * 60 * 1000; // UTC+7
    const localDate = new Date(date.getTime() + offsetMs);
    return localDate.toISOString().replace("Z", "+07:00");
}

