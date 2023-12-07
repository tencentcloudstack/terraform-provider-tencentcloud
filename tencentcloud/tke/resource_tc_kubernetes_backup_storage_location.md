Provide a resource to create tke backup storage location.

~> **NOTE:** To create this resource, you need to create a cos bucket with prefix "tke-backup" in advance.

Example Usage
```

	resource "tencentcloud_kubernetes_backup_storage_location" "example_backup" {
	  name            = "example-backup-1"
	  storage_region  = "ap-guangzhou" # region of you pre-created COS bucket
	  bucket          = "tke-backup-example-1" # bucket name of your pre-created COS bucket
	}

```

Import

tke backup storage location can be imported, e.g.

```
$ terraform import tencentcloud_kubernetes_backup_storage_location.test xxx
```