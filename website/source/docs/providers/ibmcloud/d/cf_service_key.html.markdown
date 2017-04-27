---
layout: "ibmcloud"
page_title: "IBM Cloud: ibmcloud_cf_service_key"
sidebar_current: "docs-ibmcloud-datasource-cf-service-key"
description: |-
  Get information about a Cloud Foundry service key from IBM Bluemix.
---

# ibmcloud\_cf_service_key

Import the details of an existing Cloud Foundry service key from IBM Bluemix as a read-only data source. The fields of the data source can then be referenced by other resources within the same configuration by using interpolation syntax. 

## Example Usage

```hcl
data "ibmcloud_cf_service_key" "serviceKeydata" {
  name  = "mycloudantdbKey"
  service_instance_name = "mycloudantdb"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the service key.
* `service_instance_name` - (Required) The name of the service instance that the service key is associated with.

## Attributes Reference

The following attributes are exported:

* `credentials` - The credentials associated with the key.  
