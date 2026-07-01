import * as yaml from 'js-yaml';
import { http } from '@/lib/http';

export interface ApplyConfigResponse {
  status: string;
  message: string;
}

export interface ReloadConfigResponse {
  status: string;
  message: string;
}

export async function getConfig(): Promise<string> {
  const res = await http.get<unknown>('/config_dump');
  if (typeof res === 'string') {
    try {
      const parsed = JSON.parse(res);
      return yaml.dump(parsed, { indent: 2, lineWidth: -1 });
    } catch {
      return res;
    }
  }
  return yaml.dump(res, { indent: 2, lineWidth: -1 });
}

export async function applyConfig(yamlString: string): Promise<ApplyConfigResponse> {
  return http.postRaw<ApplyConfigResponse>('/config', yamlString, 'application/x-yaml');
}

export async function reloadConfig(): Promise<ReloadConfigResponse> {
  return http.post<ReloadConfigResponse>('/reload');
}
