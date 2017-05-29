---
layout: "ibmcloud"
page_title: "IBM Cloud: cs_cluster_bind_service"
sidebar_current: "docs-ibmcloud-resource-cs-cluster-bind-service"
description: |-
  Manages IBM Cloud infrastructure cluster.
---

# ibmcloud\_cs_cluster_bind_service

Bind an IBM Cloud service to a Kubernetes namespace. With this resource, you can attach an existing service to an existing Kubernetes cluster. 

## Example Usage

In the following example, you can bind a service to a cluster.

```hcl
resource "ibmcloud_cs_cluster_service_bind" "bind_service" {
  cluster_name_id          = "cluster_name"
  service_instance_space_guid = "space_guid"
  service_instance_name_id = "service_name"
  namespace_id 			   = "default"
  org_guid = "test"
  space_guid = "test_space"
  account_guid = "test_account"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name_id` - (Required) Name or ID of the cluster.
* `service_instance_space_guid` - (Required) The space GUID the service instance is associated with.
* `service_instance_name_id` - (Required) The name or ID of the service that you want to bind to the cluster.
* `namespace_id` - (Required) The Kubernetes namespace.
* `org_guid` - (Required) The GUID for the Bluemix organization that the cluster is associated with. The values can be retrieved from data source `ibmcloud_cf_org`.
* `space_guid` - (Required) The GUID for the Bluemix space that the cluster is associated with. The values can be retrieved from data source `ibmcloud_cf_space`.
* `account_guid` - (Optional) The GUID for the Bluemix account that the cluster is associated with. The values can be retrieved from data source `ibmcloud_cf_account`.
    
## Attributes Reference

The following attributes are exported:

* `service_instance_name_id` - The name or ID of the service that is bound to the cluster.
* `namespace_id` -  The Kubernetes namespace.
* `space_guid` - The Bluemix space GUID. 
* `secret_name` - The secret name.
