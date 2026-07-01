import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { Cpu, Clock, Server, Route, CheckCircle2, AlertCircle } from 'lucide-react';
import { getRuntime } from '@/api/runtime';
import { getBackends } from '@/api/backends';
import { PageHeader } from '@/components/layout/PageHeader';
import { PageState } from '@/components/layout/PageState';
import { StatCard } from '@/components/cards/StatCard';

export const DashboardPage: React.FC = () => {
  const runtimeQuery = useQuery({
    queryKey: ['runtime'],
    queryFn: getRuntime,
  });

  const backendsQuery = useQuery({
    queryKey: ['backends'],
    queryFn: getBackends,
  });

  const isLoading = runtimeQuery.isLoading || backendsQuery.isLoading;
  const error = runtimeQuery.error || backendsQuery.error;

  if (isLoading || error) {
    return (
      <div className="space-y-6">
        <PageHeader title="Overview" description="System telemetry and infrastructure snapshot" />
        <PageState isLoading={isLoading} error={error} />
      </div>
    );
  }

  const runtime = runtimeQuery.data;
  const backends = backendsQuery.data ?? [];

  const healthyBackends = backends.filter((b) => {
    const s = b.state.toUpperCase();
    return s === 'HEALTHY' || s === 'UP' || s === 'OK' || s === 'ONLINE';
  }).length;

  const unhealthyBackends = backends.length - healthyBackends;

  return (
    <div className="space-y-6">
      <PageHeader title="Overview" description="System telemetry and infrastructure snapshot" />

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
        <StatCard
          title="Runtime Version"
          value={runtime ? `v${runtime.version}` : 'N/A'}
          subtitle="Active configuration epoch"
          icon={Cpu}
        />
        <StatCard
          title="Loaded At"
          value={runtime?.loaded_at ? new Date(runtime.loaded_at).toLocaleTimeString() : 'N/A'}
          subtitle={runtime?.loaded_at || 'Configuration timestamp'}
          icon={Clock}
        />
        <StatCard
          title="Active Services"
          value={runtime?.services ?? 0}
          subtitle="Registered clusters"
          icon={Server}
        />
        <StatCard
          title="Active Routes"
          value={runtime?.routes ?? 0}
          subtitle="Configured listeners"
          icon={Route}
        />
        <StatCard
          title="Healthy Backends"
          value={healthyBackends}
          subtitle="Passing active probes"
          icon={CheckCircle2}
          valueClassName="text-success"
        />
        <StatCard
          title="Unhealthy Backends"
          value={unhealthyBackends}
          subtitle="Failing or degraded probes"
          icon={AlertCircle}
          valueClassName={unhealthyBackends > 0 ? 'text-error' : 'text-muted-foreground'}
        />
      </div>
    </div>
  );
};
