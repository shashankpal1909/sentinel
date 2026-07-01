import React, { useEffect, useState } from 'react';
import { type Theme, ThemeProviderContext } from './theme-context';

export function ThemeProvider({
  children,
  defaultTheme = 'system',
  storageKey = 'sentinel-theme',
}: {
  children: React.ReactNode;
  defaultTheme?: Theme;
  storageKey?: string;
}): React.ReactElement {
  const [theme, setThemeState] = useState<Theme>(() => {
    if (typeof window !== 'undefined') {
      return (localStorage.getItem(storageKey) as Theme) || defaultTheme;
    }
    return defaultTheme;
  });

  const [resolvedTheme, setResolvedTheme] = useState<'dark' | 'light'>('dark');

  useEffect(() => {
    const root = window.document.documentElement;

    const computeResolvedTheme = (t: Theme): 'dark' | 'light' => {
      if (t === 'system') {
        return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
      }
      return t;
    };

    const target = computeResolvedTheme(theme);
    setResolvedTheme(target);

    root.classList.remove('light', 'dark');
    root.classList.add(target);
    root.style.setProperty('color-scheme', target);
  }, [theme]);

  useEffect(() => {
    if (theme !== 'system') return;

    const media = window.matchMedia('(prefers-color-scheme: dark)');
    const listener = () => {
      const target = media.matches ? 'dark' : 'light';
      setResolvedTheme(target);
      const root = window.document.documentElement;
      root.classList.remove('light', 'dark');
      root.classList.add(target);
      root.style.setProperty('color-scheme', target);
    };

    media.addEventListener('change', listener);
    return () => media.removeEventListener('change', listener);
  }, [theme]);

  const setTheme = (newTheme: Theme) => {
    localStorage.setItem(storageKey, newTheme);
    setThemeState(newTheme);
  };

  return (
    <ThemeProviderContext.Provider value={{ theme, setTheme, resolvedTheme }}>
      {children}
    </ThemeProviderContext.Provider>
  );
}
