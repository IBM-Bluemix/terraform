---
layout: "ibmcloud"
page_title: "IBM Cloud: ibmcloud_cf_shared_domain"
sidebar_current: "docs-ibmcloud-datasource-cf-shared-domain"
description: |-
  Get information about an IBM Bluemix shared domain.
---

# ibmcloud\_cf_shared_domain

Import the details of an existing IBM Bluemix shared domain as a read-only data source. The fields of the data source can then be referenced by other resources within the same configuration by using interpolation syntax. 

## Example Usage

```hcl

data "ibmcloud_cf_shared_domain" "shared_domain" {
	name = "foo.com"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the shared domain.

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier of the shared domain.  
