import { Dispatch, SetStateAction, useEffect, useState } from 'react';

export function useSessionState<T>(key: string, initialValue: T): [T, Dispatch<SetStateAction<T>>] {
  const [value, setValue] = useState<T>(() => {
    try {
      const raw = sessionStorage.getItem(key);
      if (raw == null) {
        return initialValue;
      }
      return JSON.parse(raw) as T;
    } catch {
      return initialValue;
    }
  });

  useEffect(() => {
    try {
      sessionStorage.setItem(key, JSON.stringify(value));
    } catch {
      // Ignore storage quota/private mode errors and keep in-memory state.
    }
  }, [key, value]);

  return [value, setValue];
}
