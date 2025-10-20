package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func MakeServerRequest(t *testing.T, call func() (*http.Response, error)) []byte {
	res, err := call()
	if err != nil {
		t.Fatalf("Failed to make request to test server: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	return body
}

func TestGetConfigs(t *testing.T) {
	app := BuildApplication()
	subject := httptest.NewServer(BuildServer(&app))
	defer subject.Close()

	err := app.ConfigDb.AddConfig(&Config{
		ConfigPath: ConfigPath{
			Service: "service1",
			Name:    "config1",
		},
		Type:         "string",
		DefaultValue: "value1",
	})
	if err != nil {
		t.Fatalf("Failed to add config to ConfigDb: %v", err)
	}

	body := MakeServerRequest(t, func() (*http.Response, error) {
		return http.Get(subject.URL + "/configs")
	})

	var response ListConfigsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	if len(response.Configs) != 1 {
		t.Errorf("Expected response to include 1 configs, but got %v", response.Configs)
	}
}

func TestAddConfig(t *testing.T) {
	app := BuildApplication()
	subject := httptest.NewServer(BuildServer(&app))
	defer subject.Close()

	req :=
		`{
			"config": {
				"service": "service1",
				"name": "config2",
				"type": "string",
				"defaultValue": "value2"
			}
		}`
	body := MakeServerRequest(t, func() (*http.Response, error) {
		return http.Post(
			subject.URL+"/configs",
			"application/json",
			strings.NewReader(req),
		)
	})

	var response PostConfigResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	if response.Message != "Success" {
		t.Errorf("Expected Success but got, but got %v", response.Message)
	}

	configs, err := app.ConfigDb.GetConfigs()
	if err != nil {
		t.Fatalf("Failed to get configs from ConfigDb: %v", err)
	}
	if len(configs) != 1 {
		t.Errorf("Expected 1 config in the ConfigDb, but got %v", configs)
	}
}

func TestGetConfig(t *testing.T) {
	app := BuildApplication()
	subject := httptest.NewServer(BuildServer(&app))
	defer subject.Close()

	configPath := ConfigPath{
		Service: "service1",
		Name:    "config1",
	}
	app.ConfigDb.AddConfig(&Config{
		ConfigPath:   configPath,
		Type:         "string",
		DefaultValue: "value1",
	})

	body := MakeServerRequest(t, func() (*http.Response, error) {
		return http.Get(subject.URL + "/configs/service1/config1")
	})

	var response GetConfigResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	if response.Config.Service != "service1" || response.Config.Name != "config1" {
		t.Errorf("Expected service1/config1, got %v/%v", response.Config.Service, response.Config.Name)
	}
	if response.Config.Type != "string" {
		t.Errorf("Expected type string, got %v", response.Config.Type)
	}
	if response.Config.DefaultValue != "value1" {
		t.Errorf("Expected defaultValue value1, got %v", response.Config.DefaultValue)
	}
}

func TestDeleteConfig(t *testing.T) {
	app := BuildApplication()
	subject := httptest.NewServer(BuildServer(&app))
	defer subject.Close()

	configPath := ConfigPath{
		Service: "service1",
		Name:    "config1",
	}
	app.ConfigDb.AddConfig(&Config{
		ConfigPath:   configPath,
		Type:         "string",
		DefaultValue: "value1",
	})

	// Delete the config
	req, err := http.NewRequest(http.MethodDelete, subject.URL+"/configs/service1/config1", nil)
	if err != nil {
		t.Fatalf("Failed to create DELETE request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send DELETE request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	// Confirm config is deleted
	_, err = app.ConfigDb.GetConfig(&configPath)
	if err == nil {
		t.Errorf("Expected error when getting deleted config, got nil")
	}
}

func TestGetOverrides(t *testing.T) {
	app := BuildApplication()
	subject := httptest.NewServer(BuildServer(&app))
	defer subject.Close()

	configPath := ConfigPath{
		Service: "service1",
		Name:    "config1",
	}
	app.ConfigDb.AddConfig(&Config{
		ConfigPath:   configPath,
		Type:         "string",
		DefaultValue: "value1",
	})
	// Add multiple overrides
	overrides := []Override{
		{OverrideKey: OverrideKey{EntityType: "user", EntityId: "1"}, Value: "A"},
		{OverrideKey: OverrideKey{EntityType: "user", EntityId: "2"}, Value: "B"},
	}
	for _, o := range overrides {
		app.ConfigDb.AddOverride(&configPath, &o)
	}

	url := subject.URL + "/configs/service1/config1/overrides"
	body := MakeServerRequest(t, func() (*http.Response, error) {
		return http.Get(url)
	})
	var response ListOverridesResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
}

func TestAddOverride(t *testing.T) {
	app := BuildApplication()
	subject := httptest.NewServer(BuildServer(&app))
	defer subject.Close()

	configPath := ConfigPath{
		Service: "service1",
		Name:    "config1",
	}
	app.ConfigDb.AddConfig(&Config{
		ConfigPath:   configPath,
		Type:         "string",
		DefaultValue: "value1",
	})
	reqBody := `{"override": {"entityType": "user", "entityId": "123", "value": "override1"}}`
	body := MakeServerRequest(t, func() (*http.Response, error) {
		return http.Post(
			subject.URL+"/configs/service1/config1/overrides",
			"application/json",
			strings.NewReader(reqBody),
		)
	})
	var response PostConfigOverrideResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	if response.Message != "Success" {
		t.Errorf("Expected Success but got %v", response.Message)
	}
	// Confirm override is added
	override, found, err := app.ConfigDb.GetOverride(&configPath, &OverrideKey{EntityType: "user", EntityId: "123"})
	if err != nil || !found {
		t.Fatalf("Expected override to be present, got error: %v", err)
	}
	if override.Value != "override1" {
		t.Errorf("Expected override value to be override1, got %v", override.Value)
	}
}

func TestGetOverride(t *testing.T) {
	app := BuildApplication()
	subject := httptest.NewServer(BuildServer(&app))
	defer subject.Close()

	configPath := ConfigPath{
		Service: "service1",
		Name:    "config1",
	}
	overrideKey := OverrideKey{EntityType: "user", EntityId: "123"}
	app.ConfigDb.AddConfig(&Config{
		ConfigPath:   configPath,
		Type:         "string",
		DefaultValue: "value1",
	})
	app.ConfigDb.AddOverride(&configPath, &Override{OverrideKey: overrideKey, Value: "override1"})

	body := MakeServerRequest(t, func() (*http.Response, error) {
		return http.Get(subject.URL + "/configs/service1/config1/overrides/user/123")
	})
	var response GetOverrideResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	if response.Override.Value != "override1" {
		t.Errorf("Expected override value to be override1, got %v", response.Override.Value)
	}
}

func TestDeleteOverride(t *testing.T) {
	app := BuildApplication()
	subject := httptest.NewServer(BuildServer(&app))
	defer subject.Close()

	configPath := ConfigPath{
		Service: "service1",
		Name:    "config1",
	}
	overrideKey := OverrideKey{EntityType: "user", EntityId: "123"}
	app.ConfigDb.AddConfig(&Config{
		ConfigPath:   configPath,
		Type:         "string",
		DefaultValue: "value1",
	})
	app.ConfigDb.AddOverride(&configPath, &Override{OverrideKey: overrideKey, Value: "override1"})

	// Delete the override
	req, err := http.NewRequest(http.MethodDelete, subject.URL+"/configs/service1/config1/overrides/user/123", nil)
	if err != nil {
		t.Fatalf("Failed to create DELETE request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send DELETE request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	// Confirm override is deleted
	_, found, err := app.ConfigDb.GetOverride(&configPath, &overrideKey)
	if found != false || err != nil {
		t.Errorf("Expected override to be deleted, but it still exists")
	}
}

func TestGetConfigDefaultValue(t *testing.T) {
	app := BuildApplication()
	subject := httptest.NewServer(BuildServer(&app))
	defer subject.Close()

	configPath := ConfigPath{
		Service: "service1",
		Name:    "config1",
	}
	app.ConfigDb.AddConfig(&Config{
		ConfigPath:   configPath,
		Type:         "string",
		DefaultValue: "value1",
	})
	override := Override{
		OverrideKey: OverrideKey{EntityType: "user", EntityId: "123"},
		Value:       "override1",
	}
	app.ConfigDb.AddOverride(&configPath, &override)

	reqBody := `{
		"attributes": {
		}
	}`

	body := MakeServerRequest(t, func() (*http.Response, error) {
		return http.Post(
			subject.URL+"/configs/service1/config1/value",
			"application/json",
			strings.NewReader(reqBody),
		)
	})

	var response GetConfigValueResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	config := response
	if config.Type != "string" {
		t.Errorf("Expected type to be string, but got %v", config.Type)
	}
	if config.Value != "value1" {
		t.Errorf("Expected value to be value1, but got %v", config.Value)
	}
}

func TestGetConfigValueOverrideBySingleEntity(t *testing.T) {
	app := BuildApplication()
	subject := httptest.NewServer(BuildServer(&app))
	defer subject.Close()

	configPath := ConfigPath{
		Service: "service1",
		Name:    "config1",
	}
	app.ConfigDb.AddConfig(&Config{
		ConfigPath:   configPath,
		Type:         "string",
		DefaultValue: "value1",
	})
	override := Override{
		OverrideKey: OverrideKey{EntityType: "user", EntityId: "123"},
		Value:       "override1",
	}
	app.ConfigDb.AddOverride(&configPath, &override)

	reqBody := `{"attributes": {"group": "456", "user": "123"}}`
	body := MakeServerRequest(t, func() (*http.Response, error) {
		return http.Post(
			subject.URL+"/configs/service1/config1/value",
			"application/json",
			strings.NewReader(reqBody),
		)
	})

	var response GetConfigValueResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	if response.Type != "string" {
		t.Errorf("Expected type to be string, but got %v", response.Type)
	}
	if response.Value != "override1" {
		t.Errorf("Expected value to be override1, but got %v", response.Value)
	}
}

func TestGetConfigValueOverrideByMultipleEntity(t *testing.T) {
	app := BuildApplication()
	subject := httptest.NewServer(BuildServer(&app))
	defer subject.Close()

	configPath := ConfigPath{
		Service: "service1",
		Name:    "config1",
	}
	app.ConfigDb.AddConfig(&Config{
		ConfigPath:   configPath,
		Type:         "string",
		DefaultValue: "value1",
	})
	// Add two overrides, only one should match
	app.ConfigDb.AddOverride(&configPath, &Override{
		OverrideKey: OverrideKey{EntityType: "user", EntityId: "123"},
		Value:       "override1",
	})
	app.ConfigDb.AddOverride(&configPath, &Override{
		OverrideKey: OverrideKey{EntityType: "group", EntityId: "456"},
		Value:       "override2",
	})

	reqBody := `{"attributes": {"user": "123", "group": "789"}}`
	body := MakeServerRequest(t, func() (*http.Response, error) {
		return http.Post(
			subject.URL+"/configs/service1/config1/value",
			"application/json",
			strings.NewReader(reqBody),
		)
	})

	var response GetConfigValueResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	if response.Type != "string" {
		t.Errorf("Expected type to be string, but got %v", response.Type)
	}
	if response.Value != "override1" {
		t.Errorf("Expected value to be override1, but got %v", response.Value)
	}
}
