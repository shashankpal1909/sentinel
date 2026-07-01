import React from 'react';
import { Play, RotateCcw, Check, AlertTriangle, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

interface ConfigToolbarProps {
  isDirty: boolean;
  isValidating: boolean;
  isApplying: boolean;
  isReloading: boolean;
  onValidate: () => void;
  onReload: () => void;
  onApply: () => void;
}

export const ConfigToolbar: React.FC<ConfigToolbarProps> = ({
  isDirty,
  isValidating,
  isApplying,
  isReloading,
  onValidate,
  onReload,
  onApply,
}) => {
  const isBusy = isValidating || isApplying || isReloading;

  return (
    <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 p-4 rounded-xl border border-border bg-card shadow-2xs">
      <div className="flex items-center gap-3">
        {isDirty ? (
          <Badge
            variant="outline"
            className="bg-warning/10 text-warning border-warning/30 flex items-center gap-1.5 px-2.5 py-1"
          >
            <AlertTriangle className="w-3.5 h-3.5" />
            <span>Unsaved Changes</span>
          </Badge>
        ) : (
          <Badge
            variant="outline"
            className="bg-success/10 text-success border-success/30 flex items-center gap-1.5 px-2.5 py-1"
          >
            <Check className="w-3.5 h-3.5" />
            <span>Synced</span>
          </Badge>
        )}
        <span className="text-xs text-muted-foreground hidden md:inline">
          {isDirty
            ? 'Editor differs from active gateway runtime'
            : 'Editor matches active runtime configuration'}
        </span>
      </div>

      <div className="flex items-center gap-2.5 w-full sm:w-auto justify-end">
        <Button
          variant="outline"
          size="sm"
          onClick={onValidate}
          disabled={isBusy}
          className="gap-1.5 font-mono text-xs"
        >
          {isValidating ? (
            <Loader2 className="w-3.5 h-3.5 animate-spin" />
          ) : (
            <Check className="w-3.5 h-3.5" />
          )}
          <span>Validate</span>
        </Button>

        <Button
          variant="outline"
          size="sm"
          onClick={onReload}
          disabled={isBusy}
          title="Reload active configuration from disk (Ctrl+R)"
          className="gap-1.5 font-mono text-xs"
        >
          {isReloading ? (
            <Loader2 className="w-3.5 h-3.5 animate-spin" />
          ) : (
            <RotateCcw className="w-3.5 h-3.5" />
          )}
          <span>{isReloading ? 'Reloading...' : 'Reload'}</span>
        </Button>

        <Button
          size="sm"
          onClick={onApply}
          disabled={!isDirty || isBusy}
          title="Apply configuration hot reload (Ctrl+S)"
          className="gap-1.5 font-mono text-xs"
        >
          {isApplying ? (
            <Loader2 className="w-3.5 h-3.5 animate-spin" />
          ) : (
            <Play className="w-3.5 h-3.5" />
          )}
          <span>{isApplying ? 'Applying...' : 'Apply'}</span>
        </Button>
      </div>
    </div>
  );
};
