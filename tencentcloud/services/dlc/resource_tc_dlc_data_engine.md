Provides a resource to create a DLC data engine

Example Usage

```hcl
resource "tencentcloud_dlc_data_engine" "example" {
  engine_type        = "spark"
  data_engine_name   = "tf-example"
  cluster_type       = "spark_cu"
  mode               = 1
  auto_resume        = false
  size               = 16
  min_clusters       = 1
  max_clusters       = 1
  cidr_block         = "10.255.0.0/16"
  message            = "DLC data engine demo."
  image_version_name = "Standard-S 1.1"
  engine_exec_type   = "BATCH"
  engine_generation  = "Native"
  session_resource_template {
    driver_size          = "medium"
    executor_max_numbers = 7
    executor_nums        = 1
    executor_size        = "medium"
  }
}
```

Import

DLC data engine can be imported using the dataEngineName#dataEngineId, e.g.

```
terraform import tencentcloud_dlc_data_engine.example tf-example#DataEngine-d3gk8r5h
```
