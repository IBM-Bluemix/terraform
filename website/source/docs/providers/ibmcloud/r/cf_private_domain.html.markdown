---
layout: "ibmcloud"
page_title: "IBM Cloud: cf_private_domain"
sidebar_current: "docs-ibmcloud-resource-cf-private_domain"
description: |-
  Manages IBM Cloud Cloud Foundry private domain.
---

# ibmcloud\_cf_private_domain

Create, update, or delete CF private domain on IBM Bluemix.

## Example Usage

```hcl
	
data "ibmcloud_cf_org" "orgdata" {
		org    = "someexample.com"
	}
		
resource "ibmcloud_cf_private_domain" "domain" {
		name              = "foo.com"
		org_guid          = "${data.ibmcloud_cf_org.orgdata.id}"
	}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the domain.
* `org_guid` - (Required, string) The GUID of the organization that owns the domain. The values can be retrieved from data source `ibmcloud_cf_org`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the private domain.

