---
layout: "ibmcloud"
page_title: "IBM Cloud: ibmcloud_cf_app"
sidebar_current: "docs-ibmcloud-datasource-cf-app"
description: |-
  Get information about an IBM Bluemix app.
---

# ibmcloud\_cf_app

Import the details of an existing IBM Bluemix app as a read-only data source. The fields of the data source can then be referenced by other resources within the same configuration by using interpolation syntax. 

## Example Usage

```hcl

	data "ibmcloud_cf_app" "testacc_ds_app" {
		name = "my-app"
		space_guid = "${ibmcloud_cf_app.app.space_guid}"
}


```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the application.
* `space_guid` - (Required, string) Define space guid to which application belongs. The values can be retrieved from data source `ibmcloud_cf_space`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the application.
* `memory` - Memory allocated to the application. In megabytes.
* `instances` - The number of instances of the application.
* `disk_quota` - The disk allocated to an instance of an application. In megabytes.
* `buildpack` - Buildpack used by the application.It can be a) Blank means autodetection; b) A Git Url pointing to a buildpack; c) Name of an installed buildpack.
* `environment_json` - Key/value pairs of all the environment variables. Does not include any system or service variables.
* `route_guid` - The route guids which are bound to the application.
* `service_instance_guid` - The service instance guids which are bound to the application.
* `package_state` - The state of the application package whether staged, pending etc.
* `state` - The state of the application.
