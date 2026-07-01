export interface RouteListener {
  path: string;
  service: string;
}

export interface ListenersResponse {
  listeners?: RouteListener[];
}
