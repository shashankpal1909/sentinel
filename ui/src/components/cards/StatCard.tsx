import React from 'react';
import type { LucideIcon } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';

interface StatCardProps {
  title: string;
  value: string | number;
  subtitle?: string;
  icon?: LucideIcon;
  className?: string;
  valueClassName?: string;
}

export const StatCard: React.FC<StatCardProps> = ({
  title,
  value,
  subtitle,
  icon: Icon,
  className,
  valueClassName,
}) => {
  return (
    <Card
      className={cn(
        'group transition-all duration-150 hover:border-foreground/20 hover:shadow-sm flex flex-col justify-between py-4',
        className
      )}
    >
      <CardHeader className="pb-2 flex flex-row items-center justify-between space-y-0">
        <CardTitle className="text-xs font-semibold uppercase tracking-wider text-muted-foreground group-hover:text-foreground transition-colors">
          {title}
        </CardTitle>
        {Icon && (
          <div className="p-2 rounded-lg bg-secondary text-muted-foreground group-hover:text-primary transition-colors duration-150">
            <Icon className="w-4 h-4 shrink-0" />
          </div>
        )}
      </CardHeader>
      <CardContent className="pt-0">
        <div
          className={cn(
            'text-2xl font-mono font-bold tracking-tight text-foreground',
            valueClassName
          )}
        >
          {value}
        </div>
        {subtitle && (
          <p className="mt-1.5 text-xs text-muted-foreground font-sans truncate leading-normal">
            {subtitle}
          </p>
        )}
      </CardContent>
    </Card>
  );
};
