---
layout: "ibmcloud"
page_title: "IBM Cloud: ibmcloud_cf_service_plan"
sidebar_current: "docs-ibmcloud-datasource-cf-service-plan"
description: |-
  Get information about a Cloud Foundry service plan from IBM Bluemix.
---

# ibmcloud\_cf_service_plan

Import the details of an existing Cloud Foundry service plan from IBM Bluemix as a read-only data source. The fields of the data source can then be referenced by other resources within the same configuration by using interpolation syntax. 

## Example Usage

```hcl
data "ibmcloud_cf_service_plan" "service_plan" {
  service  = "cleardb"
  plan    = "spark"
}
```

## Argument Reference

The following arguments are supported:

* `service` - (Required, string) The name of the service offering. Use the `bluemix service offerings` command to find the name and the plan of the service that you require. Installing Bluemix cli can be found [here](https://console.ng.bluemix.net/docs/cli/reference/bluemix_cli/index.html#getting-started)
* `plan` - (Required, string) The name of the plan type supported by service. Use the `bluemix service offering` command to find the name and the plan of the service that you require.

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier of the service plan.  
