import { http } from '@/lib/http';
import { getServices } from './services';
import type { BackendDetail } from '@/types/backend';

export async function getBackends(): Promise<BackendDetail[]> {
  try {
    const data = await http.get<{ backends?: BackendDetail[] }>('/backends');
    if (data && Array.isArray(data.backends)) {
      return data.backends;
    }
  } catch {
    // Fallback to deriving from /clusters if /backends endpoint returns 404
  }

  const clusters = await getServices();
  const backends: BackendDetail[] = [];

  for (const cluster of clusters) {
    const hc = cluster.health_check;
    for (const backend of cluster.backends) {
      backends.push({
        service: cluster.name,
        url: backend.url,
        state: backend.state,
        interval: hc?.interval ?? 'N/A',
        healthyThreshold: hc?.healthy_threshold ?? 0,
        unhealthyThreshold: hc?.unhealthy_threshold ?? 0,
      });
    }
  }

  return backends;
}
