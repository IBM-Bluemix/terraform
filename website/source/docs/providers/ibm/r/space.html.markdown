---
layout: "ibm"
page_title: "IBM : Space"
sidebar_current: "docs-ibm-resource-space"
description: |-
  Manages IBM Space.
---

# ibm\_space

Create, update, or delete spaces for IBM Bluemix.

## Example Usage

```hcl
resource "ibm_space" "space" {
  name        = "myspace"
  org         = "myorg"
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

* `id` - The ID of the new space.
