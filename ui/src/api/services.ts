import { http } from '@/lib/http';
import type { ServiceCluster, ClustersResponse } from '@/types/service';

export async function getServices(): Promise<ServiceCluster[]> {
  const response = await http.get<ClustersResponse>('/clusters');
  return response.clusters ?? [];
}
