Provides a resource to create a CBS storage.

-> **NOTE:** When creating an encrypted disk, if `kms_key_id` is not entered, the product side will generate a key by default.

-> **NOTE:** When using CBS encrypted disk, it is necessary to add `CVM_QcsRole` role and `QcloudKMSAccessForCVMRole` strategy to the account.

Example Usage

Create a standard CBS storage

```hcl
resource "tencentcloud_cbs_storage" "example" {
  storage_name      = "tf-example"
  storage_type      = "CLOUD_SSD"
  storage_size      = 100
  availability_zone = "ap-guangzhou-3"
  project_id        = 0
  encrypt           = false

  tags = {
    createBy = "Terraform"
  }
}
```

Create an encrypted CBS storage with customize kms_key_id

```hcl
resource "tencentcloud_cbs_storage" "example" {
  storage_name      = "tf-example"
  storage_type      = "CLOUD_SSD"
  storage_size      = 100
  availability_zone = "ap-guangzhou-3"
  project_id        = 0
  kms_key_id        = "2e860789-7ef0-11ef-8d1c-5254001955d1"
  encrypt           = true

  tags = {
    createBy = "Terraform"
  }
}
```

Create an encrypted CBS storage with default generated kms_key_id

```hcl
resource "tencentcloud_cbs_storage" "example" {
  storage_name      = "tf-example"
  storage_type      = "CLOUD_SSD"
  storage_size      = 100
  availability_zone = "ap-guangzhou-3"
  project_id        = 0
  encrypt           = true

  tags = {
    createBy = "Terraform"
  }
}
```

Create an encrypted CBS storage with encrypt_type

```hcl
resource "tencentcloud_cbs_storage" "example" {
  storage_name      = "tf-example"
  storage_type      = "CLOUD_SSD"
  storage_size      = 100
  availability_zone = "ap-guangzhou-3"
  project_id        = 0
  encrypt           = true
  encrypt_type      = "ENCRYPT_V2"

  tags = {
    createBy = "Terraform"
  }
}
```

Create a dedicated cluster CBS storage

```hcl
resource "tencentcloud_cbs_storage" "example" {
  storage_name         = "tf-example"
  storage_type         = "CLOUD_SSD"
  storage_size         = 100
  availability_zone    = "ap-guangzhou-4"
  dedicated_cluster_id = "cluster-262n63e8"
  charge_type          = "DEDICATED_CLUSTER_PAID"
  project_id           = 0
  encrypt              = false

  tags = {
    createBy = "Terraform"
  }
}
```

Import

CBS storage can be imported using the id, e.g.

```
$ terraform import tencentcloud_cbs_storage.example disk-41s6jwy4
```
