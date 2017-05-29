---
layout: "ibmcloud"
page_title: "IBM Cloud: ibmcloud_cf_private_domain"
sidebar_current: "docs-ibmcloud-datasource-cf-private-domain"
description: |-
  Get information about an IBM Bluemix private domain.
---

# ibmcloud\_cf_private_domain

Import the details of an existing IBM Bluemix private domain as a read-only data source. The fields of the data source can then be referenced by other resources within the same configuration by using interpolation syntax. 

## Example Usage

```hcl

data "ibmcloud_cf_private_domain" "private_domain" {
	name = "foo.com"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the private domain.

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier of the private domain.  
