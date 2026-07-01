import React from 'react';
import { cn } from '@/lib/utils';

interface PageHeaderProps {
  title: string;
  description?: string;
  className?: string;
  children?: React.ReactNode;
}

export const PageHeader: React.FC<PageHeaderProps> = ({
  title,
  description,
  className,
  children,
}) => {
  return (
    <div
      className={cn(
        'flex flex-col sm:flex-row sm:items-center sm:justify-between pb-5 mb-6 border-b border-border transition-colors duration-150',
        className
      )}
    >
      <div>
        <h1 className="text-2xl font-bold tracking-tight text-foreground">{title}</h1>
        {description && (
          <p className="text-sm text-muted-foreground mt-1 font-normal leading-relaxed">
            {description}
          </p>
        )}
      </div>
      {children && <div className="mt-4 sm:mt-0 flex items-center gap-2">{children}</div>}
    </div>
  );
};
