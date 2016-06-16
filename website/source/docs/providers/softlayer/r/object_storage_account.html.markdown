---
layout: "softlayer"
page_title: "SoftLayer: Object Storage Account"
sidebar_current: "docs-softlayer-object-storage-account"
description: |-
  Provides Softlayer Object Storage Account
---

# softlayer_objectstorage_account

<div class="alert alert-block alert-info" style="padding-bottom:0"><strong>Note:</strong> For managing SoftLayer object storage **containers** and **objects**, please see the [Swift provider](/docs/providers/swift/index.html), since SoftLayer's object storage is an implementation of Swift object storage.</div>

Ensures there is an existing object storage account within your SoftLayer account. If there is an existing object storage, it will learn its account name and keep it as its ID for future usage. If there is no object storage account, it will order one for you and remember the account name. It is not meant to be used for managing the life cycle of an object storage account in SoftLayer (e.g. update, delete) at this time.

## Example Usage

```
resource "softlayer_objectstorage_account" "foo" {
}
```

## Argument Reference

No additional arguments needed.

## Computed Fields

* `id` - The object storage account name, which you can later use with [Swift resources](/docs/providers/swift/index.html).
