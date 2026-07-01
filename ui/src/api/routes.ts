import { http } from '@/lib/http';
import type { RouteListener, ListenersResponse } from '@/types/route';

export async function getRoutes(): Promise<RouteListener[]> {
  const response = await http.get<ListenersResponse>('/listeners');
  return response.listeners ?? [];
}
