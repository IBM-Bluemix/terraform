---
layout: "ibmcloud"
page_title: "IBM Cloud: cf_service_instance"
sidebar_current: "docs-ibmcloud-resource-cf-service-instance"
description: |-
  Manages IBM Cloud Cloud Foundry service instance.
---

# ibmcloud\_cf_service_instance

Crate, update, or delete CF service instances on IBM Bluemix.

## Example Usage

```hcl
data "ibmcloud_cf_space" "spacedata" {
  space  = "prod"
  org    = "somexample.com"
}

resource "ibmcloud_cf_service_instance" "service_instance" {
  name              = "test"
  space_guid        = "${data.ibmcloud_cf_space.spacedata.id}"
  service           = "cloudantNoSQLDB"
  plan              = "Lite"
  tags              = ["cluster-service", "cluster-bind"]
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) A descriptive name used to identify the service instance.
* `space_guid` - (Required, string) The GUID of the space where you want to create the service. The values can be retrieved from data source `ibmcloud_cf_space`.
* `service` - (Required, string) The name of the service offering. Use the `bluemix service offerings` command to find the name and the plan of the service that you require. Installing Bluemix cli can be found [here](https://console.ng.bluemix.net/docs/cli/reference/bluemix_cli/index.html#getting-started)
* `plan` - (Required, string) The name of the plan type supported by service. Use the `bluemix service offerings` command to find the name and the plan of the service that you require.
* `metadata` - (Optional, map) Valid JSON object containing service-specific configuration parameters.
* `tags` - (Optional, list) User-provided tags.
* `parameters` - (Optional, map) Arbitrary parameters to pass along to the service broker. Must be a JSON object.

## Attributes Reference6

The following attributes are exported:

* `id` - The ID of the new service instance.
* `credentials` - The credentials associated with the service instance.
* `service_plan_guid` - The plan of the service offering used by this service instance 

