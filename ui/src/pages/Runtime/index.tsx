import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { getRuntime } from '@/api/runtime';
import { PageHeader } from '@/components/layout/PageHeader';
import { PageState } from '@/components/layout/PageState';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';

export const RuntimePage: React.FC = () => {
  const {
    data: runtime,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['runtime'],
    queryFn: getRuntime,
  });

  if (isLoading || error) {
    return (
      <div className="space-y-6">
        <PageHeader title="Runtime State" description="Current gateway process metadata" />
        <PageState isLoading={isLoading} error={error} />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <PageHeader title="Runtime State" description="Current gateway process metadata" />

      <Card className="max-w-2xl">
        <CardHeader className="border-b border-border pb-3.5">
          <CardTitle className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
            Runtime Snapshot
          </CardTitle>
        </CardHeader>
        <CardContent className="pt-6">
          <dl className="grid grid-cols-1 sm:grid-cols-2 gap-y-6 gap-x-8 font-mono">
            <div>
              <dt className="text-xs text-muted-foreground font-sans">Runtime Version</dt>
              <dd className="mt-1.5 text-lg font-bold text-foreground">
                {runtime ? `v${runtime.version}` : 'N/A'}
              </dd>
            </div>
            <div>
              <dt className="text-xs text-muted-foreground font-sans">Loaded Timestamp</dt>
              <dd className="mt-1.5 text-xs font-medium text-foreground">
                {runtime?.loaded_at || 'Never'}
              </dd>
            </div>
            <div>
              <dt className="text-xs text-muted-foreground font-sans">Registered Services</dt>
              <dd className="mt-1.5 text-lg font-bold text-foreground">{runtime?.services ?? 0}</dd>
            </div>
            <div>
              <dt className="text-xs text-muted-foreground font-sans">Configured Routes</dt>
              <dd className="mt-1.5 text-lg font-bold text-foreground">{runtime?.routes ?? 0}</dd>
            </div>
          </dl>
        </CardContent>
      </Card>
    </div>
  );
};
