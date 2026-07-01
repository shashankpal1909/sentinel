import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { getServices } from '@/api/services';
import type { ServiceCluster } from '@/types/service';
import { PageHeader } from '@/components/layout/PageHeader';
import { PageState } from '@/components/layout/PageState';
import { DataTable, type Column } from '@/components/tables/DataTable';

export const ServicesPage: React.FC = () => {
  const {
    data: services,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['services'],
    queryFn: getServices,
  });

  if (isLoading || error) {
    return (
      <div className="space-y-6">
        <PageHeader title="Services" description="Configured upstream load balancing clusters" />
        <PageState isLoading={isLoading} error={error} />
      </div>
    );
  }

  const columns: Column<ServiceCluster>[] = [
    {
      header: 'Cluster Name',
      accessorKey: 'name',
      cell: (item) => <span className="font-bold text-foreground">{item.name}</span>,
    },
    {
      header: 'Strategy',
      accessorKey: 'strategy',
      cell: (item) => (
        <span className="px-2.5 py-0.5 rounded-[6px] bg-secondary text-secondary-foreground text-xs font-medium border border-border">
          {item.strategy || 'round-robin'}
        </span>
      ),
    },
    {
      header: 'Backend Count',
      cell: (item) => <span className="text-foreground">{item.backends?.length ?? 0}</span>,
    },
    {
      header: 'Health Check Path',
      cell: (item) => (
        <span className="text-secondary-foreground">{item.health_check?.path ?? 'None'}</span>
      ),
    },
  ];

  return (
    <div className="space-y-6">
      <PageHeader title="Services" description="Configured upstream load balancing clusters" />
      <DataTable
        columns={columns}
        data={services ?? []}
        keyExtractor={(item) => item.name}
        emptyMessage="No service clusters registered."
      />
    </div>
  );
};
