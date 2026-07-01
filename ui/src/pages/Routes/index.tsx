import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { getRoutes } from '@/api/routes';
import type { RouteListener } from '@/types/route';
import { PageHeader } from '@/components/layout/PageHeader';
import { PageState } from '@/components/layout/PageState';
import { DataTable, type Column } from '@/components/tables/DataTable';

export const RoutesPage: React.FC = () => {
  const {
    data: routes,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['routes'],
    queryFn: getRoutes,
  });

  if (isLoading || error) {
    return (
      <div className="space-y-6">
        <PageHeader title="Routes" description="HTTP listener path routing rules" />
        <PageState isLoading={isLoading} error={error} />
      </div>
    );
  }

  const columns: Column<RouteListener>[] = [
    {
      header: 'Listener Path',
      accessorKey: 'path',
      cell: (item) => <span className="font-bold text-foreground">{item.path}</span>,
    },
    {
      header: 'Upstream Service',
      accessorKey: 'service',
      cell: (item) => (
        <span className="text-foreground/80 font-medium">{item.service || 'N/A'}</span>
      ),
    },
  ];

  return (
    <div className="space-y-6">
      <PageHeader title="Routes" description="HTTP listener path routing rules" />
      <DataTable
        columns={columns}
        data={routes ?? []}
        keyExtractor={(item, idx) => `${item.path}-${item.service}-${idx}`}
        emptyMessage="No listener routes configured."
      />
    </div>
  );
};
