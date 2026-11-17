Provides a resource to manage MongoDB instance elastic CPU scale-up.

~> **NOTE:** This resource is used to manually enable elastic CPU scale-up for MongoDB instances. When the resource is created, `ScaleUpDBInstanceCpu` is called to scale up the CPU. When the resource is destroyed, `ScaleDownDBInstanceCpu` is called to disable the elastic CPU scale-up and shrink back.

~> **NOTE:** Suggest using instances with `machine_type` values of `HIO10G` to operate CPU elastic scaling, and if CPU elastic scaling is required, it is not recommended to actually configure the instance CPU value.

Example Usage

```hcl
resource "tencentcloud_mongodb_instance" "example" {
  instance_name  = "tf-example"
  memory         = 4
  volume         = 100
  engine_version = "MONGO_80_WT"
  machine_type   = "HIO10G"
  available_zone = "ap-guangzhou-6"
  vpc_id         = "vpc-i5yyodl9"
  subnet_id      = "subnet-hhi88a58"
  project_id     = 0
  password       = "Password@2026"
}

resource "tencentcloud_mongodb_scale_up_db_instance_cpu" "example" {
  instance_id = tencentcloud_mongodb_instance.example.id
  extra_cpu   = 2
}
```

Import

MongoDB instance elastic CPU scale-up can be imported using the instance id, e.g.

```
terraform import tencentcloud_mongodb_scale_up_db_instance_cpu.example cmgo-p8vnipr5
```
