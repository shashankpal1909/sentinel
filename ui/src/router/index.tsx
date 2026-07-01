import React from 'react';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { AppLayout } from '@/components/layout/AppLayout';
import { DashboardPage } from '@/pages/Dashboard';
import { RuntimePage } from '@/pages/Runtime';
import { ServicesPage } from '@/pages/Services';
import { RoutesPage } from '@/pages/Routes';
import { BackendsPage } from '@/pages/Backends';

const router = createBrowserRouter([
  {
    path: '/',
    element: <AppLayout />,
    children: [
      {
        index: true,
        element: <DashboardPage />,
      },
      {
        path: 'runtime',
        element: <RuntimePage />,
      },
      {
        path: 'services',
        element: <ServicesPage />,
      },
      {
        path: 'routes',
        element: <RoutesPage />,
      },
      {
        path: 'backends',
        element: <BackendsPage />,
      },
    ],
  },
]);

export const AppRouter: React.FC = () => {
  return <RouterProvider router={router} />;
};
