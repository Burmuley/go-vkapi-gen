# VK API Golang SDK Generator

**!WARNING! Work in Progress!**

Tool to generate Golang (Go) SDK code for VK API implementation based on public JSON schema.

### Description

This tool uses official [VK API JSON schema](https://github.com/VKCOM/vk-api-schema) as input ro generate SDK code for Golang.

The result code is located here: [https://gitlab.com/Burmuley/go-vkapi](https://gitlab.com/Burmuley/go-vkapi).
See [README.md](https://gitlab.com/Burmuley/go-vkapi/blob/master/README.md) in the target repository for details about SDK.

### Settings

The tool by default use these sources:
* `objects` - [https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/objects.json](https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/objects.json) 
* `responses` - [https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/responses.json](https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/responses.json)
* `methods` - [https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/methods.json](https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/methods.json)

These settings can be overridden by setting corresponding environment variables:
* to override `objects` - VK_API_SCHEMA_OBJECTS  
* to override `responses` - VK_API_SCHEMA_RESPONSES
* to override `methods` - VK_API_SCHEMA_METHODS

### Usage
```bash
$ go build -o go-vkapi-gen
$ ./go-vkapi-gen
```
