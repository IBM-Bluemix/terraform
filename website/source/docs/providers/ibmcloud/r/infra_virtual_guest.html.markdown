---
layout: "ibmcloud"
page_title: "IBM Cloud: infra_virtual_guest"
sidebar_current: "docs-ibmcloud-resource-infra-virtual-guest"
description: |-
  Manages IBM Cloud infrastructure virtual guests.
---

# ibmcloud\_infra_virtual_guest

Provides a resource for virtual guests. This allows virtual guests to be created, updated, and deleted. 

For additional details, see the [Bluemix Infrastructure (SoftLayer) API docs](http://sldn.softlayer.com/reference/services/SoftLayer_Virtual_Guest).

## Example Usage

In the following example, you can create a virtual guest using a Debian image.

```hcl
resource "ibmcloud_infra_virtual_guest" "twc_terraform_sample" {
    hostname = "twc-terraform-sample-name"
    domain = "bar.example.com"
    os_reference_code = "DEBIAN_7_64"
    datacenter = "wdc01"
    network_speed = 10
    hourly_billing = true
    private_network_only = false
    cores = 1
    memory = 1024
    disks = [25, 10, 20]
    user_metadata = "{\"value\":\"newvalue\"}"
    dedicated_acct_host_only = true
    local_disk = false
    public_vlan_id = 1391277
    private_vlan_id = 7721931
}
```

In the following example, you can create a virtual guest using a block device template.

```hcl
resource "ibmcloud_infra_virtual_guest" "terraform-sample-BDTGroup" {
   hostname = "terraform-sample-blockDeviceTemplateGroup"
   domain = "bar.example.com"
   datacenter = "ams01"
   public_network_speed = 10
   hourly_billing = false
   cores = 1
   memory = 1024
   local_disk = false
   image_id = 12345
   tags = [
     "collectd",
     "mesos-master"
   ]
   public_subnet = "50.97.46.160/28"
   private_subnet = "10.56.109.128/26"
}
```

## Argument Reference

The following arguments are supported:

*   `hostname` - (Optional) Hostname for the computing instance.
*   `domain` - (Required)  Domain for the computing instance.
*   `cores` - (Required) The number of CPU cores to allocate.
*   `memory` - (Required) The amount of memory to allocate, expressed in megabytes.
*   `datacenter` -  (Required) Specify which data center the instance is to be provisioned in.
*   `hourly_billing` - (Optional) Specify the billing type for the instance. When set to `true`, the computing instance is billed on hourly usage, otherwise it is billed on a monthly basis. Default value: `true`.
*   `local_disk`- (Optional) Specify the disk type for the instance. When set to `true`, the disks for the computing instance are provisioned on the host that it runs, otherwise SAN disks are provisioned. Default value: `true`.
*   `dedicated_acct_host_only` - (Optional) Specify whether or not the instance must only run on hosts with instances from the same account. Default value: `false`.
*   `os_reference_code` - (Optional) An operating system reference code that is used to provision the computing instance. [Get a complete list of the OS reference codes available](https://api.softlayer.com/rest/v3/SoftLayer_Virtual_Guest_Block_Device_Template_Group/getVhdImportSoftwareDescriptions.json?objectMask=referenceCode) (use your API key as the password to log in). 

    **NOTE**: Conflicts with`image_id`.
*   `image_id` - (Optional) The image template ID to be used to provision the computing instance. Note this is not the global identifier (UUID), but the image template group ID that should point to a valid global identifier. You can get the image template ID in the SoftLayer Customer Portal by navigating to **Devices > Manage > Images**. Click the desired image and take note of the ID number in the browser URL location. 

    **NOTE**: Conflicts with `os_reference_code`. If you don't know the ID(s) for your image templates, [you can reference them by name](../d/infra_image_template.html).
*   `network_speed` - (Optional) Specify the connection speed (in Mbps) for the instance's network components. Default value: `100`.
*   `private_network_only` - (Optional) When set to `true`, a compute instance only has access to the private network. Default value: `false`.
*   `public_vlan_id` - (Optional) Public VLAN ID for the public network interface of the instance. Accepted values are in the [VLAN doc](https://control.softlayer.com/network/vlans). Click the desired VLAN and note the ID in the resulting URL. You can also [refer to a VLAN by name using a data source](../d/infra_vlan.html).
* `private_vlan_id` - (Optional) Private VLAN ID for the private network interface of the instance. Accepted values are in the [VLAN doc](https://control.softlayer.com/network/vlans). Click the desired VLAN and note the ID in the resulting URL. You can also [refer to a VLAN by name using a data source](../d/infra_vlan.md).
* `public_subnet` - (Optional) Public subnet for the public network interface of the instance. Accepted values are primary public networks and can be found in the [subnets doc](https://control.softlayer.com/network/subnets).
* `private_subnet` - (Optional) Private subnet for the private network interface of the instance. Accepted values are primary private networks and can be found in the  [subnets doc](https://control.softlayer.com/network/subnets).
* `disks` - (Optional, array) Numeric disk sizes in GBs. Block device and disk image settings for the computing instance. Defaults to the smallest available capacity for the primary disk are used. If an image template is specified, the disk capacity is provided by the template.
* `user_metadata` - (Optional) Arbitrary data to be made available to the computing instance.
* `ssh_key_ids` - (Optional) An array of numbers. SSH key IDs to install on the computing instance upon provisioning.
    **NOTE**: If you don't know the ID(s) for your SSH keys, [you can reference your SSH keys by their labels](../d/infra_ssh_key.html).
* `post_install_script_uri` - (Optional)  As defined in the [Bluemix Infrastructure (SoftLayer) API](https://sldn.softlayer.com/reference/datatypes/SoftLayer_Virtual_Guest_SupplementalCreateObjectOptions).
* `tags` - (Optional, array of strings) Set tags on the virtual guest. Permitted characters include: A-Z, 0-9, whitespace, _ (underscore), - (hyphen), . (period), and : (colon). All other characters are removed.
* `ipv6_enabled` - (Optional) Provides a primary public IPv6 address. Default value: `false`.
*  `secondary_ip_count` - (Optional) Provides secondary public IPv4 addresses. Accepted values are `4` and `8`. 
*  `wait_time_minutes` - (Optional) The duration, expressed in minutes, to wait for the virtual guest to become available before declaring it as created. It is also the same amount of time waited for no active transactions before proceeding with an update or deletion. Default value: `90`.


## Attributes Reference

The following attributes are exported:

* `id` - ID of the virtual guest.
* `ipv4_address` - Public IPv4 address of the virtual guest.
* `ip_address_id_private` - Unique ID for the private IPv4 address assigned to the virtual guest.
* `ipv4_address_private` - Private IPv4 address of the virtual guest.
* `ip_address_id` - Unique ID for the public IPv4 address assigned to the virtual guest.
* `ipv6_address` - Public IPv6 address of the virtual guest. It is provided when `ipv6_enabled` is set to `true`.
* `ipv6_address_id` - Unique ID for the public IPv6 address assigned to the virtual_guest. It is provided when `ipv6_enabled` is set to `true`.
* `public_ipv6_subnet` - Public IPv6 subnet. It is provided when `ipv6_enabled` is set to `true`.
* `secondary_ip_addresses` - Public secondary IPv4 addresses of the virtual guest.
