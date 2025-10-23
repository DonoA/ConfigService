ConfigService
===

A Terraform compatible, distributed config service and client libraries. Backed by redis, ConfigService is designed to support hundreds of client services, thousands of config options, and millions of unique overrides.

This repo contains 3 key elements:
1. The core service which serves requests for config status and updates
2. The Terraform provider which supports defining and updating config values in code
3. A basic web application which provides an overview of all config objects and their values

Config data is stored in a heirarchy with Service at the top level, followed by config name. This is intended to make config usage easy to locate within source code. Config values can have only 1 data type attached to them. These value types are: `bool`, `str`, `long`, and `float`

Config values are overriden based on the value of different entity ids. For example a single request might have a userId, groupId, and resourceId attached. Each of these entity types and their id value can be used to determine if any overrides should be applied to the resulting value of the config flag. If multiple entity ids have overrides attached to them, an override will be selected essentially at random among them.

