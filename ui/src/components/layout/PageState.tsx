import React from 'react';
import { AlertCircle } from 'lucide-react';
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert';
import { Skeleton } from '@/components/ui/skeleton';

interface PageStateProps {
  isLoading?: boolean;
  error?: Error | null;
  loadingMessage?: string;
}

export const PageState: React.FC<PageStateProps> = ({
  isLoading,
  error,
  loadingMessage = 'Synchronizing gateway state...',
}) => {
  if (isLoading) {
    return (
      <div className="space-y-4" role="status" aria-label={loadingMessage}>
        <div className="bg-card border border-border rounded-xl p-4 flex items-center gap-3 text-xs font-mono text-muted-foreground shadow-2xs">
          <div className="w-2 h-2 rounded-full bg-primary animate-ping shrink-0" />
          <span>{loadingMessage}</span>
        </div>

        {/* Shimmer skeleton grid */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className="h-28 rounded-xl bg-card border border-border/80 p-5 space-y-3 shadow-2xs"
            >
              <Skeleton className="h-3 w-24 rounded" />
              <Skeleton className="h-6 w-16 rounded" />
              <Skeleton className="h-2.5 w-32 rounded" />
            </div>
          ))}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <Alert variant="destructive" className="border-destructive/40 bg-destructive/10">
        <AlertCircle className="size-4" />
        <AlertTitle>State Synchronization Failure</AlertTitle>
        <AlertDescription>
          {error.message || 'An unexpected error occurred communicating with the control plane.'}
        </AlertDescription>
      </Alert>
    );
  }

  return null;
};
