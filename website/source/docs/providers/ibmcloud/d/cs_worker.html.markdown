---
layout: "ibmcloud"
page_title: "IBM Cloud: ibmcloud_cs_worker"
sidebar_current: "docs-ibmcloud-datasource-cs-worker"
description: |-
  Get information about a worker node that is attached to a Kubernetes cluster on IBM Bluemix.
---

# ibmcloud\_cs_worker


Import details about the worker node of a Kubernetes cluster as a read-only data source. The fields of the data source can then be referenced by other resources within the same configuration using interpolation syntax. 


## Example Usage

```hcl
data "ibmcloud_cs_cluster" "cluster_foo" {
    worker_id = "dev-mex10-pa70c4414695c041518603bfd0cd6e333a-w1"
    org_guid = "test"
	space_guid = "test_space"
	account_guid = "test_acc"
}
```

## Argument Reference

The following arguments are supported:

* `worker_id` - (Required) ID of the worker node attached to the cluster.
* `org_guid` - (Required) The GUID for the Bluemix organization that the cluster is associated with. The value can be retrieved from the `ibmcloud_cf_org` data source.
* `space_guid` - (Required) The GUID for the Bluemix space that the cluster is associated with. The value can be retrieved from the `ibmcloud_cf_space` data source.
* `account_guid` - (Required) The GUID for the Bluemix account that the cluster is associated with. The value can be retrieved from the `ibmcloud_cf_account` data source.


## Attributes Reference

The following attributes are exported:

* `state` - The unique identifier of the cluster.
* `status` - Number of workers nodes attached to the cluster.
* `private_vlan` - The private VLAN of the worker node.
* `public_vlan` -  The public VLAN of the worker node.
* `private_ip` - The private IP of the worker node.
* `public_ip` -  The public IP of the worker node.
