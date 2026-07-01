import type { Backend } from './backend';

export interface HealthCheckConfig {
  path: string;
  interval: string;
  timeout?: string;
  healthy_threshold: number;
  unhealthy_threshold: number;
}

export interface ServiceCluster {
  name: string;
  strategy: string;
  health_check?: HealthCheckConfig;
  backends: Backend[];
}

export interface ClustersResponse {
  clusters?: ServiceCluster[];
}
