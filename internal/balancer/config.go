package balancer

type Config struct {
	Port string `json:"port"`
	HealthCheckInterval string `json:"healthCheckInterval"`
	Algorithm           string   `json:"algorithm"`
	Servers []string `json:"servers"`
}