import { http } from '@/lib/http';
import type { RuntimeInfo } from '@/types/runtime';

export async function getRuntime(): Promise<RuntimeInfo> {
  return http.get<RuntimeInfo>('/runtime');
}
