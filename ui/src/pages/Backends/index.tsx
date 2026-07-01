import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { getBackends } from '@/api/backends';
import type { BackendDetail } from '@/types/backend';
import { PageHeader } from '@/components/layout/PageHeader';
import { PageState } from '@/components/layout/PageState';
import { DataTable, type Column } from '@/components/tables/DataTable';
import { StatusBadge } from '@/components/badges/StatusBadge';

export const BackendsPage: React.FC = () => {
  const {
    data: backends,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['backends'],
    queryFn: getBackends,
  });

  if (isLoading || error) {
    return (
      <div className="space-y-6">
        <PageHeader title="Backends" description="Upstream server pool health and thresholds" />
        <PageState isLoading={isLoading} error={error} />
      </div>
    );
  }

  const columns: Column<BackendDetail>[] = [
    {
      header: 'Service Cluster',
      accessorKey: 'service',
      cell: (item) => <span className="font-semibold text-foreground">{item.service}</span>,
    },
    {
      header: 'Endpoint URL',
      accessorKey: 'url',
      cell: (item) => <span className="text-foreground/90 font-mono">{item.url}</span>,
    },
    {
      header: 'Health Status',
      accessorKey: 'state',
      cell: (item) => <StatusBadge status={item.state} />,
    },
    {
      header: 'Probe Interval',
      accessorKey: 'interval',
    },
    {
      header: 'Healthy Threshold',
      accessorKey: 'healthyThreshold',
    },
    {
      header: 'Unhealthy Threshold',
      accessorKey: 'unhealthyThreshold',
    },
  ];

  return (
    <div className="space-y-6">
      <PageHeader title="Backends" description="Upstream server pool health and thresholds" />
      <DataTable
        columns={columns}
        data={backends ?? []}
        keyExtractor={(item, idx) => `${item.service}-${item.url}-${idx}`}
        emptyMessage="No backend endpoints active."
      />
    </div>
  );
};
