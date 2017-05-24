---
layout: "ibmcloud"
page_title: "IBM Cloud: ibmcloud_cf_domainse"
sidebar_current: "docs-ibmcloud-datasource-cf-domain"
description: |-
  Get information about an IBM Bluemix domain.
---

# ibmcloud\_cf_domain

Import the details of an existing IBM Bluemix domain as a read-only data source. The fields of the data source can then be referenced by other resources within the same configuration by using interpolation syntax. 

## Example Usage

```hcl

// Shared Domain

data "ibmcloud_cf_domain" "testacc_domain" {
	name = mybluemix.net"
	domain_type = "shared"
}

// Private Domain

data "ibmcloud_cf_domain" "testacc_domain" {
	name = "privatedomain.net"
	domain_type = "private"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the domain.
* `domain_type` - (Optional) The type of the domain.Accepted values are 'shared' or 'private'

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier of the domain.  
