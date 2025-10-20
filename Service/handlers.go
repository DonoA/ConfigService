package main

import (
	"encoding/json"
	"io"
	"maps"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type HttpResponse struct {
	Status int
	Data   []byte
}

func GetConfigPath(urlVars map[string]string) (*ConfigPath, error) {
	service, ok := urlVars["service"]
	if !ok {
		return nil, errors.New("missing service in url vars")
	}
	name, ok := urlVars["name"]
	if !ok {
		return nil, errors.New("missing name in url vars")
	}

	return &ConfigPath{
		Service: service,
		Name:    name,
	}, nil
}

func GetOverrideKey(urlVars map[string]string) (*OverrideKey, error) {
	entityType, ok := urlVars["entityType"]
	if !ok {
		return nil, errors.New("missing entityType in url vars")
	}
	entityId, ok := urlVars["entityId"]
	if !ok {
		return nil, errors.New("missing entityId in url vars")
	}
	return &OverrideKey{
		EntityType: entityType,
		EntityId:   entityId,
	}, nil
}

type Handlers struct {
	ConfigDb ConfigDb
}

func (h *Handlers) ListConfigs(r *http.Request) (*HttpResponse, error) {
	configs, err := h.ConfigDb.GetConfigs()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get configs from db")
	}
	response := ListConfigsResponse{
		Configs: configs,
	}
	body, err := json.Marshal(response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal response")
	}

	return &HttpResponse{
		Status: http.StatusOK,
		Data:   body,
	}, nil
}

func (h *Handlers) PostConfig(r *http.Request) (*HttpResponse, error) {
	var requestBody PostConfigRequest
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read request body")
	}
	err = json.Unmarshal(bodyBytes, &requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode request body")
	}
	if requestBody.Config.Name == "" ||
		requestBody.Config.Service == "" {
		return nil, errors.New("config name and service are required")
	}

	err = h.ConfigDb.AddConfig(&requestBody.Config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add config to db")
	}

	response := PostConfigResponse{
		Message: "Success",
	}
	respBytes, err := json.Marshal(response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal response")
	}
	return &HttpResponse{
		Status: http.StatusOK,
		Data:   respBytes,
	}, nil
}

func (h *Handlers) GetConfig(r *http.Request) (*HttpResponse, error) {
	urlVars := mux.Vars(r)
	configPath, err := GetConfigPath(urlVars)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config name from request")
	}

	config, err := h.ConfigDb.GetConfig(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config from db")
	}

	response := GetConfigResponse{
		Config: config,
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal response")
	}

	return &HttpResponse{
		Status: http.StatusOK,
		Data:   responseBytes,
	}, nil
}

func (h *Handlers) DeleteConfig(r *http.Request) (*HttpResponse, error) {
	urlVars := mux.Vars(r)
	configPath, err := GetConfigPath(urlVars)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config name from request")
	}

	err = h.ConfigDb.DeleteConfig(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to delete config from db")
	}

	response := DeleteConfigResponse{
		Message: "Success",
	}
	respBytes, err := json.Marshal(response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal response")
	}
	return &HttpResponse{
		Status: http.StatusOK,
		Data:   respBytes,
	}, nil
}

func (h *Handlers) ListOverrides(r *http.Request) (*HttpResponse, error) {
	urlVars := mux.Vars(r)
	configPath, err := GetConfigPath(urlVars)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config name from request")
	}
	overrides, err := h.ConfigDb.GetOverrides(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get overrides from db")
	}
	response := ListOverridesResponse{
		Overrides: overrides,
	}
	respBytes, err := json.Marshal(response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal response")
	}
	return &HttpResponse{
		Status: http.StatusOK,
		Data:   respBytes,
	}, nil
}

func (h *Handlers) PostOverride(r *http.Request) (*HttpResponse, error) {
	urlVars := mux.Vars(r)
	configPath, err := GetConfigPath(urlVars)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config name from request")
	}

	var requestBody PostConfigOverrideRequest
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read request body")
	}
	err = json.Unmarshal(bodyBytes, &requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode request body")
	}

	err = h.ConfigDb.AddOverride(configPath, &requestBody.Override)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add override to db")
	}

	response := PostConfigOverrideResponse{
		Message: "Success",
	}
	respBytes, err := json.Marshal(response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal response")
	}

	return &HttpResponse{
		Status: http.StatusOK,
		Data:   respBytes,
	}, nil
}

func (h *Handlers) GetOverride(r *http.Request) (*HttpResponse, error) {
	urlVars := mux.Vars(r)
	configPath, err := GetConfigPath(urlVars)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config name from request")
	}
	overrideKey, err := GetOverrideKey(urlVars)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get override key from request")
	}
	override, found, err := h.ConfigDb.GetOverride(configPath, overrideKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get override from db")
	}
	if !found {
		return nil, errors.New("override not found")
	}
	response := GetOverrideResponse{
		Override: override,
	}
	respBytes, err := json.Marshal(response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal response")
	}
	return &HttpResponse{
		Status: http.StatusOK,
		Data:   respBytes,
	}, nil
}

func (h *Handlers) DeleteOverride(r *http.Request) (*HttpResponse, error) {
	urlVars := mux.Vars(r)
	configPath, err := GetConfigPath(urlVars)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config name from request")
	}
	overrideKey, err := GetOverrideKey(urlVars)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get override key from request")
	}
	err = h.ConfigDb.DeleteOverride(configPath, overrideKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to delete override from db")
	}
	response := DeleteOverrideResponse{
		Message: "Success",
	}
	respBytes, err := json.Marshal(response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal response")
	}
	return &HttpResponse{
		Status: http.StatusOK,
		Data:   respBytes,
	}, nil
}

func (h *Handlers) GetConfigValue(r *http.Request) (*HttpResponse, error) {
	urlVars := mux.Vars(r)
	configPath, err := GetConfigPath(urlVars)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config name from request")
	}

	config, err := h.ConfigDb.GetConfig(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config from db")
	}
	configValue := config.DefaultValue

	var requestBody GetConfigValueRequest
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read request body")
	}
	err = json.Unmarshal(bodyBytes, &requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode request body")
	}

	entityAttributes := maps.All(requestBody.Attributes)
	for key, value := range entityAttributes {
		overrideKey := OverrideKey{
			EntityType: key,
			EntityId:   value,
		}

		override, found, err := h.ConfigDb.GetOverride(configPath, &overrideKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get override from db")
		}
		if found {
			configValue = override.Value
			break
		}
	}

	response := GetConfigValueResponse{
		Type:  config.Type,
		Value: configValue,
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal response")
	}

	return &HttpResponse{
		Status: http.StatusOK,
		Data:   responseBytes,
	}, nil
}
