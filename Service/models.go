package main

type Config struct {
	ConfigPath
	Type         string `json:"type"`
	DefaultValue string `json:"defaultValue"`
}

type ConfigPath struct {
	Service string `json:"service"`
	Name    string `json:"name"`
}

type Override struct {
	OverrideKey
	Value string `json:"value"`
}

type OverrideKey struct {
	EntityType string `json:"entityType"`
	EntityId   string `json:"entityId"`
}

type SimpleResponse struct {
	Message string `json:"message"`
}

type ListConfigsResponse struct {
	Configs []Config `json:"configs"`
}

type PostConfigRequest struct {
	Config Config `json:"config"`
}

type PostConfigResponse = SimpleResponse

type GetConfigResponse struct {
	Config Config `json:"config"`
}

type PostConfigOverrideRequest struct {
	Override Override `json:"override"`
}

type PostConfigOverrideResponse = SimpleResponse

type DeleteConfigResponse = SimpleResponse

type ListOverridesResponse struct {
	Overrides []Override `json:"overrides"`
}

type GetOverrideResponse struct {
	Override Override `json:"override"`
}

type DeleteOverrideResponse = SimpleResponse

type GetConfigValueRequest struct {
	Attributes map[string]string `json:"attributes"`
}

type GetConfigValueResponse struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
