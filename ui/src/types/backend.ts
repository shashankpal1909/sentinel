export interface Backend {
  url: string;
  state: string;
  healthy?: boolean;
}

export interface BackendDetail {
  service: string;
  url: string;
  state: string;
  healthy?: boolean;
  interval: string;
  healthyThreshold: number;
  unhealthyThreshold: number;
}
