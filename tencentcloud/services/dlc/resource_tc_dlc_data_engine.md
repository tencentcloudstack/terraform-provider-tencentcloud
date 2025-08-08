Provides a resource to create a DLC data engine

Example Usage

```hcl
resource "tencentcloud_dlc_data_engine" "example" {
  engine_type            = "spark"
  data_engine_name       = "tf-example"
  cluster_type           = "spark_cu"
  mode                   = 1
  auto_resume            = false
  size                   = 16
  min_clusters           = 1
  max_clusters           = 1
  cidr_block             = "10.255.0.0/16"
  message                = "demo"
  pay_mode               = 0
  time_unit              = "h"
  time_span              = 3600
  auto_renew             = 2
  engine_exec_type       = "SQL"
  auto_suspend           = false
  default_data_engine    = false
  crontab_resume_suspend = 0
  engine_generation      = "Native"
}
```

Import

DLC data engine can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_data_engine.example tf-example#DataEngine-d3gk8r5h
```
