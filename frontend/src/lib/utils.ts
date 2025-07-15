import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export const formatTime = (timeString: string) => {
  // Convert "HH:mm:ss" to "HH:mm"
  return timeString.split(':').slice(0, 2).join(':');
};