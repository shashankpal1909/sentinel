export interface RuntimeInfo {
  version: number;
  loaded_at: string;
  services: number;
  routes: number;
  healthy_backends?: number;
  unhealthy_backends?: number;
}
