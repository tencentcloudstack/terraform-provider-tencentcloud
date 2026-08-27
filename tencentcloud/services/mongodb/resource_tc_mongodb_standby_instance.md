Provide a resource to create a Mongodb standby instance.

Example Usage

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

resource "tencentcloud_mongodb_standby_instance" "example" {
  instance_name          = "tf-example"
  cpu                    = 2
  memory                 = 4
  volume                 = 100
  available_zone         = "ap-guangzhou-7"
  vpc_id                 = "vpc-i5yyodl9"
  subnet_id              = "subnet-d4umunpy"
  project_id             = 0
  father_instance_id     = tencentcloud_mongodb_instance.example.id
  father_instance_region = "ap-guangzhou"
  data_encryption        = "TDE"
  encryption_key_source  = "manual"
  key_id                 = "KMS-MONGODB"
  kms_region             = "ap-guangzhou"

  tags = {
    createBy = "Terraform"
  }
}
```

Import

Mongodb standby instance can be imported using the id, e.g.

```
terraform import tencentcloud_mongodb_standby_instance.mongodb cmgo-41s6jwy4
```
