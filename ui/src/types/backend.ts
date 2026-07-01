export interface Backend {
  url: string;
  state: string;
}

export interface BackendDetail {
  service: string;
  url: string;
  state: string;
  interval: string;
  healthyThreshold: number;
  unhealthyThreshold: number;
}
