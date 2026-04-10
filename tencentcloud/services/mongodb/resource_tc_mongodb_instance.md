Provide a resource to create a Mongodb instance.

~> **NOTE:** If `availability_zone_list` needs to be changed, attention should be paid to cascading modifications of `available_zone` or `hidden_zone`.

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
  password       = "Password@123"
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
  password   = "Password@123"
}
```

Import

Mongodb instance can be imported using the id, e.g.

```
terraform import tencentcloud_mongodb_instance.example cmgo-41s6jwy4
```