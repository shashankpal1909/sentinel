import React from 'react';
import { Activity, Sun, Moon, Monitor } from 'lucide-react';
import { useTheme } from '@/hooks/useTheme';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';

export const Header: React.FC = () => {
  const { theme, setTheme, resolvedTheme } = useTheme();

  const cycleTheme = () => {
    if (theme === 'dark') setTheme('light');
    else if (theme === 'light') setTheme('system');
    else setTheme('dark');
  };

  return (
    <header className="h-14 bg-card border-b border-border px-6 flex items-center justify-between shrink-0 transition-colors duration-150">
      <div className="flex items-center space-x-3">
        <span className="text-sm font-semibold tracking-tight text-foreground">
          Control Plane Telemetry
        </span>
      </div>

      <div className="flex items-center space-x-3">
        {/* Connection Status Badge */}
        <Badge
          variant="outline"
          className="flex items-center gap-1.5 bg-success/10 text-success border-success/30 px-3 py-1 rounded-lg text-xs font-mono tracking-wide shadow-2xs hover:bg-success/20"
          aria-label="Admin API Connected"
        >
          <Activity className="w-3.5 h-3.5 animate-pulse text-success shrink-0" />
          <span>Connected</span>
        </Badge>

        {/* Theme Toggle Button */}
        <Button
          onClick={cycleTheme}
          type="button"
          variant="outline"
          size="icon"
          aria-label={`Current theme: ${theme}. Click to switch theme.`}
          title={`Theme: ${theme.toUpperCase()} (Resolved: ${resolvedTheme})`}
          className="rounded-lg shadow-2xs transition-all duration-150 active:scale-[0.98]"
        >
          {theme === 'dark' && <Moon className="w-4 h-4 text-primary" />}
          {theme === 'light' && <Sun className="w-4 h-4 text-warning" />}
          {theme === 'system' && <Monitor className="w-4 h-4 text-muted-foreground" />}
        </Button>
      </div>
    </header>
  );
};
