package config

type RoutingConfig struct {
	RoutingURL string `env:"ROUTING_URL"`
}

type RouteBinding struct {
	RouteServiceURL string `json:"route_service_url"`
}
