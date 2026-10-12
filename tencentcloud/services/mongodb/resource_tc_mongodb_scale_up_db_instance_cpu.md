Provides a resource to manage MongoDB instance elastic CPU scale-up.

~> **NOTE:** This resource is used to manually enable elastic CPU scale-up for MongoDB instances. When the resource is created, `ScaleUpDBInstanceCpu` is called to scale up the CPU. When the resource is destroyed, `ScaleDownDBInstanceCpu` is called to disable the elastic CPU scale-up and shrink back.

Example Usage

```hcl
resource "tencentcloud_mongodb_scale_up_db_instance_cpu" "example" {
  instance_id = "cmgo-p8vnipr5"
  extra_cpu   = 4
}
```

Import

MongoDB instance elastic CPU scale-up can be imported using the instance id, e.g.

```
terraform import tencentcloud_mongodb_scale_up_db_instance_cpu.example cmgo-p8vnipr5
```