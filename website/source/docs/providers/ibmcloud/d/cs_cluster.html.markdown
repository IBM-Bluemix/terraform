---
layout: "ibmcloud"
page_title: "IBM Cloud: ibmcloud_cs_cluster"
sidebar_current: "docs-ibmcloud-datasource-cs-cluster"
description: |-
  Get information about a Kubernetes cluster on IBM Bluemix.
---

# ibmcloud\_cs_cluster


Import the details for a Kubernetes cluster on IBM Bluemix as a read-only data source. The fields of the data source can then be referenced by other resources within the same configuration using interpolation syntax. 


## Example Usage

```hcl
data "ibmcloud_cs_cluster" "cluster_foo" {
    cluster_name_id = "FOO"
    org_guid = "test"
	space_guid = "test_space"
	account_guid = "test_acc"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name_id` - (Required) Name or ID of the cluster.
* `org_guid` - (Required) The GUID for the Bluemix organization that the cluster is associated with. The value can be retrieved from the `ibmcloud_cf_org` data source.
* `space_guid` - (Required) The GUID for the Bluemix space that the cluster is associated with. The value can be retrieved from the `ibmcloud_cf_space` data source.
* `account_guid` - (Required) The GUID for the Bluemix account that the cluster is associated with. The value can be retrieved from the `ibmcloud_cf_account` data source.


## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier of the cluster.
* `worker_count` - Number of workers attached to the cluster.
* `workers` - IDs of the worker attached to the cluster.
