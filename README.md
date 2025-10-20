ConfigService
===

A Terraform compatible, distributed config service and client libraries. Backed by redis, ConfigService is designed to support hundreds of client services, thousands of config options, and millions of unique overrides.

This repo contains 3 key elements:
1. The core service which serves requests for config status and updates
2. The Terraform provider which supports defining and updating config values in code
3. A basic web application which provides an overview of all config objects and their values

Config data is stored in a heirarchy with Namespace at the top level, followed by service, followed by config name. This is intended to make config usage easy to locate within source code. Config values can have only 1 data type attached to them. These value types are: `bool`, `str`, `long`, and `float`

