---
layout: "ibmcloud"
page_title: "IBM Cloud: cf_space"
sidebar_current: "docs-ibmcloud-resource-cf-space"
description: |-
  Manages IBM Cloud Cloud Foundry space.
---

# ibmcloud\_cf_space

Create, update, or delete CF spaces for IBM Bluemix.

## Example Usage

```hcl

resource "ibmcloud_cf_space" "space" {
	name = "myspace"
	org = "myorg"
	space_quota = "myspacequota"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) A descriptive name used to identify a space.
* `org` - (Required, string) Name of the org this space belongs to.
* `space_quota` - (Optional, string) The name of the Space Quota Definition associated with the space.

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the new space.
