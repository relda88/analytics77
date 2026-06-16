package sconfiguration

type ServiceConfiguration struct {
	pathConfiguration string
}

func NewServiceConfiguration() (*ServiceConfiguration, error) {
	return &ServiceConfiguration{},
		nil
}

func (s *ServiceConfiguration) Get() map[string]any {
	return map[string]any{}
}
