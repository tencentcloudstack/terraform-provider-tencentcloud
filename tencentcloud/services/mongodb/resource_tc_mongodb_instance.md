Provide a resource to create a Mongodb instance.

~> **NOTE:** If `availability_zone_list` needs to be changed, attention should be paid to cascading modifications of `available_zone` or `hidden_zone`.

~> **NOTE:** The `cpu` parameter takes effect only when the configuration is changed. Changing the `cpu` triggers the `ModifyDBInstanceSpec` API to adjust the CPU specification of the running MongoDB instance in-place. The supported CPU specifications can be obtained through the `DescribeSpecInfo` API.

Example Usage

```hcl
resource "tencentcloud_mongodb_instance" "example" {
  instance_name  = "tf-example"
  memory         = 4
  volume         = 100
  engine_version = "MONGO_40_WT"
  machine_type   = "HIO10G"
  available_zone = "ap-guangzhou-6"
  vpc_id         = "vpc-i5yyodl9"
  subnet_id      = "subnet-hhi88a58"
  project_id     = 0
  password       = "Password@2026"
}
```

Update the CPU specification of the MongoDB instance.

```hcl
resource "tencentcloud_mongodb_instance" "example" {
  instance_name  = "tf-example"
  memory         = 4
  volume         = 100
  engine_version = "MONGO_40_WT"
  machine_type   = "HIO10G"
  available_zone = "ap-guangzhou-6"
  vpc_id         = "vpc-i5yyodl9"
  subnet_id      = "subnet-hhi88a58"
  project_id     = 0
  password       = "Password@2026"
  cpu            = 2
}
```

Or

```hcl
resource "tencentcloud_mongodb_instance" "example" {
  instance_name  = "tf-example"
  memory         = 4
  volume         = 100
  engine_version = "MONGO_40_WT"
  machine_type   = "HIO10G"
  available_zone = "ap-guangzhou-6"
  availability_zone_list = [
    "ap-guangzhou-6",
    "ap-guangzhou-3",
    "ap-guangzhou-4",
  ]
  hidden_zone = "ap-guangzhou-4"
  vpc_id     = "vpc-i5yyodl9"
  subnet_id  = "subnet-hhi88a58"
  project_id = 0
  password   = "Password@2026"
}
```

Create instance with auto encryption

```hcl
resource "tencentcloud_mongodb_instance" "example" {
  instance_name         = "tf-example"
  cpu                   = 2
  memory                = 4
  volume                = 100
  engine_version        = "MONGO_80_WT"
  machine_type          = "GE.LD.T1"
  available_zone        = "ap-guangzhou-6"
  vpc_id                = "vpc-i5yyodl9"
  subnet_id             = "subnet-hhi88a58"
  project_id            = 0
  password              = "Password@2026"
  data_encryption       = "TDE"
  encryption_key_source = "auto"
}
```

Create instance with custom encryption

```hcl
resource "tencentcloud_mongodb_instance" "example" {
  instance_name         = "tf-example"
  cpu                   = 2
  memory                = 4
  volume                = 100
  engine_version        = "MONGO_80_WT"
  machine_type          = "GE.LD.T1"
  available_zone        = "ap-guangzhou-6"
  vpc_id                = "vpc-i5yyodl9"
  subnet_id             = "subnet-hhi88a58"
  project_id            = 0
  password              = "Password@2026"
  data_encryption       = "TDE"
  encryption_key_source = "manual"
  key_id                = "KMS-MONGODB"
  kms_region            = "ap-guangzhou"
}
```

Import

Mongodb instance can be imported using the id, e.g.

```
terraform import tencentcloud_mongodb_instance.example cmgo-41s6jwy4
```