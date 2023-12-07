Provides a resource to create a dlc data_engine

Example Usage

```hcl
resource "tencentcloud_dlc_data_engine" "data_engine" {
  engine_type = "spark"
  data_engine_name = "testSpark"
  cluster_type = "spark_cu"
  mode = 1
  auto_resume = false
  size = 16
  pay_mode = 0
  min_clusters = 1
  max_clusters = 1
  default_data_engine = false
  cidr_block = "10.255.0.0/16"
  message = "test spark1"
  time_span = 1
  time_unit = "h"
  auto_suspend = false
  crontab_resume_suspend = 0
  engine_exec_type = "BATCH"
}
```

Import

dlc data_engine can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_data_engine.data_engine data_engine_id
```