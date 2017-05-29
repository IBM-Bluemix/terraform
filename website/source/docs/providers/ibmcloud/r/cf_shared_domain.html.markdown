---
layout: "ibmcloud"
page_title: "IBM Cloud: cf_shared_domain"
sidebar_current: "docs-ibmcloud-resource-cf-shared_domain"
description: |-
  Manages IBM Cloud Cloud Foundry shared domain.
---

# ibmcloud\_cf_shared_domain

Create, update, or delete CF shared domain on IBM Bluemix.

## Example Usage

```hcl
	
resource "ibmcloud_cf_shared_domain" "domain" {
		name              = "foo.com"
		router_group_guid = "3hG5jkjk4k34JH5666"
	}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the domain.
* `router_group_guid` - (Optional, string) The GUID of the router group.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the shared domain.

