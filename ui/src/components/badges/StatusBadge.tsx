import React from 'react';
import { cn } from '@/lib/utils';
import { Badge } from '@/components/ui/badge';

interface StatusBadgeProps {
  status: string;
  className?: string;
}

export const StatusBadge: React.FC<StatusBadgeProps> = ({ status, className }) => {
  const upper = status.toUpperCase();

  const isHealthy = upper === 'HEALTHY' || upper === 'UP' || upper === 'OK' || upper === 'ONLINE';
  const isUnhealthy =
    upper === 'UNHEALTHY' || upper === 'DOWN' || upper === 'OFFLINE' || upper === 'ERR';

  return (
    <Badge
      variant="outline"
      className={cn(
        'inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-md text-[11px] font-semibold border font-mono uppercase tracking-wider shadow-2xs select-none transition-colors duration-150',
        isHealthy && 'bg-success/10 text-success border-success/30 hover:bg-success/20',
        isUnhealthy && 'bg-error/10 text-error border-error/30 hover:bg-error/20',
        !isHealthy &&
          !isUnhealthy &&
          'bg-secondary/60 text-muted-foreground border-border hover:bg-secondary',
        className
      )}
    >
      <span
        className={cn(
          'w-1.5 h-1.5 rounded-full shrink-0',
          isHealthy && 'bg-success animate-pulse',
          isUnhealthy && 'bg-error',
          !isHealthy && !isUnhealthy && 'bg-muted-foreground'
        )}
      />
      {status}
    </Badge>
  );
};
