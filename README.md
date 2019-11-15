# VK API Golang SDK Generator

## **! WARNING ! WORK IN PROGRESS !**

** Reviews appreciated! **

Tool to generate Golang (Go) SDK code for VK API implementation based on public JSON schema.

### Description

This tool uses official [VK API JSON schema](https://github.com/VKCOM/vk-api-schema) as input to generate VK API SDK code for Golang.

The result code is located here: [https://github.com/Burmuley/go-vkapi](https://github.com/Burmuley/go-vkapi).
See [README.md](https://github.com/Burmuley/go-vkapi/blob/master/README.md) in the target repository for details about SDK.

### Current status
What works:
* Golang types generation for all VK API data structures enlisted in [`objects.json`](https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/objects.json) schema; result code is located at [`objects`](https://github.com/Burmuley/go-vkapi/tree/master/objects) subdirectory
* Golang types generation for all VK API responses enlisted in [`responses.json`](https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/responses.json) schema; result code is located at [`responses`](https://github.com/Burmuley/go-vkapi/tree/master/responses) subdirectory
* Golang types generation for all VK API methods enlisted in [`metods.json`](https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/methods.json) schema; result code is located at [`repo root`](https://github.com/Burmuley/go-vkapi/tree/master)
* Include of static code (common interfaces and VK API interaction utils)
* Golang code formatting for generated code
* Documentation (i.e. description) is taken from JSON schema files, i.e. no documentation in JSON schema - no documentation in produced code

What TODO:
* Create output directory structure on the fly
* Add tests for the generator code
* Add code to generate tests for resulting SDK code

### Settings

The tool by default use these inputs:
* `objects` - [https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/objects.json](https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/objects.json) 
* `responses` - [https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/responses.json](https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/responses.json)
* `methods` - [https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/methods.json](https://raw.githubusercontent.com/VKCOM/vk-api-schema/master/methods.json)

Default output directory is `output` located in the repository root. 

These settings can be overridden by setting corresponding environment variables:
* to override `objects` - set `VK_API_SCHEMA_OBJECTS` environment variable  
* to override `responses` - set `VK_API_SCHEMA_RESPONSES` environment variable
* to override `methods` - set `VK_API_SCHEMA_METHODS`  environment variable
* to override `output` directory location - set `VK_API_SCHEMA_OUTPUT` environment variable

Variables values can be of two types:
1. HTTP URL to the target file (doesn't support any kind of authentication)
1. Local file path (format depends on OS, tested with unix-like ones)

### Automated builds

There is a CI/CD job triggered on each commit to `master` branch, which builds the tool and then commits to a remote repository, creating a merge request.

### Manual tool usage

Tool does not require any command line parameters.

```bash
$ go build -o go-vkapi-gen
$ ./go-vkapi-gen
```

### License
All the code (in this repository and produced by the tool) is licensed under [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) license. 