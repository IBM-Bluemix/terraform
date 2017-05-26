---
layout: "ibmcloud"
page_title: "IBM Cloud: cf_route"
sidebar_current: "docs-ibmcloud-resource-cf-route"
description: |-
  Manages IBM Cloud Cloud Foundry route.
---

# ibmcloud\_cf_route

Create, update, or delete CF route on IBM Bluemix.

## Example Usage

```hcl
	
data "ibmcloud_cf_space" "spacedata" {
		space  = "space"
		org    = "someexample.com"
	}
		
data "ibmcloud_cf_shared_domain" "domain" {
		name        = "mybluemix.net"
	}
		
resource "ibmcloud_cf_route" "route" {
		domain_guid       = "${data.ibmcloud_cf_shared_domain.domain.id}"
		space_guid        = "${data.ibmcloud_cf_space.spacedata.id}"
		host              = "somehost172"
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

* `id` - The ID of the route.

