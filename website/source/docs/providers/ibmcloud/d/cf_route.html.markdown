---
layout: "ibmcloud"
page_title: "IBM Cloud: ibmcloud_cf_route"
sidebar_current: "docs-ibmcloud-datasource-cf-route"
description: |-
  Get information about an IBM Bluemix route.
---

# ibmcloud\_cf_route

Import the details of an existing IBM Bluemix route as a read-only data source. The fields of the data source can then be referenced by other resources within the same configuration by using interpolation syntax. 

## Example Usage

```hcl

	data "ibmcloud_cf_route" "route" {
		domain_guid       = "${data.ibmcloud_cf_shared_domain.domain.id}"
		space_guid        = "${data.ibmcloud_cf_space.spacedata.id}"
		host              = "somehost"
		path              = "/app"
	}


```

## Argument Reference

The following arguments are supported:

* `domain_guid` - (Required, string) The GUID of the associated domain. The values can be retrieved from data source `ibmcloud_cf_shared_domain`.
* `space_guid` - (Required, string) The GUID of the space where you want to create the route. The values can be retrieved from data source `ibmcloud_cf_space`.
* `host` - (Optional, string) The host portion of the route. Required for shared-domains.
* `port` - (Optional, int) The port of the route. Supported for domains of TCP router groups only.
* `path` - (Optional, string) The path for a route as raw text.Paths must be between 2 and 128 characters.Paths must start with a forward slash '/'.Paths must not contain a '?'.


## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier of the route.  
