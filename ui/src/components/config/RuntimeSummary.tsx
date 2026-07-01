import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { Cpu, Server, Route, Clock, CheckCircle2, AlertCircle } from 'lucide-react';
import { getRuntime } from '@/api/runtime';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';

export const RuntimeSummary: React.FC = () => {
  const runtimeQuery = useQuery({
    queryKey: ['runtime'],
    queryFn: getRuntime,
  });

  const isLoading = runtimeQuery.isLoading;
  const runtime = runtimeQuery.data;

  const healthyBackends = runtime?.healthy_backends ?? 0;
  const unhealthyBackends = runtime?.unhealthy_backends ?? 0;

  return (
    <Card className="w-full">
      <CardHeader className="border-b border-border pb-3">
        <CardTitle className="text-xs font-semibold uppercase tracking-wider text-muted-foreground flex items-center gap-2">
          <Cpu className="w-3.5 h-3.5" />
          Runtime Information
        </CardTitle>
      </CardHeader>
      <CardContent className="pt-4 space-y-4 font-mono text-xs">
        {isLoading ? (
          <div className="space-y-3">
            <Skeleton className="h-4 w-3/4" />
            <Skeleton className="h-4 w-1/2" />
            <Skeleton className="h-4 w-2/3" />
            <Skeleton className="h-4 w-full" />
          </div>
        ) : (
          <div className="divide-y divide-border/60">
            <div className="flex items-center justify-between py-2.5 first:pt-0">
              <span className="text-muted-foreground font-sans">Version</span>
              <span className="font-bold text-foreground">
                {runtime ? `v${runtime.version}` : 'N/A'}
              </span>
            </div>
            <div className="flex items-center justify-between py-2.5">
              <span className="text-muted-foreground font-sans flex items-center gap-1.5">
                <Clock className="w-3.5 h-3.5 text-muted-foreground/70" />
                Loaded At
              </span>
              <span className="text-[11px] font-medium text-foreground">
                {runtime?.loaded_at ? new Date(runtime.loaded_at).toLocaleTimeString() : 'N/A'}
              </span>
            </div>
            <div className="flex items-center justify-between py-2.5">
              <span className="text-muted-foreground font-sans flex items-center gap-1.5">
                <Server className="w-3.5 h-3.5 text-muted-foreground/70" />
                Services
              </span>
              <span className="font-semibold text-foreground">{runtime?.services ?? 0}</span>
            </div>
            <div className="flex items-center justify-between py-2.5">
              <span className="text-muted-foreground font-sans flex items-center gap-1.5">
                <Route className="w-3.5 h-3.5 text-muted-foreground/70" />
                Routes
              </span>
              <span className="font-semibold text-foreground">{runtime?.routes ?? 0}</span>
            </div>
            <div className="flex items-center justify-between py-2.5">
              <span className="text-muted-foreground font-sans flex items-center gap-1.5">
                <CheckCircle2 className="w-3.5 h-3.5 text-success" />
                Healthy Backends
              </span>
              <span className="font-semibold text-success">{healthyBackends}</span>
            </div>
            <div className="flex items-center justify-between py-2.5 last:pb-0">
              <span className="text-muted-foreground font-sans flex items-center gap-1.5">
                <AlertCircle className="w-3.5 h-3.5 text-error" />
                Unhealthy Backends
              </span>
              <span
                className={`font-semibold ${unhealthyBackends > 0 ? 'text-error' : 'text-muted-foreground'}`}
              >
                {unhealthyBackends}
              </span>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
};
