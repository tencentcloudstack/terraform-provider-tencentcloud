Use this data source to query monitor data. for complex queries, use (https://github.com/tencentyun/tencentcloud-exporter)

Example Usage

```hcl
data "tencentcloud_instances" "instances" {
}

#cvm
data "tencentcloud_monitor_data" "cvm_monitor_data" {
  namespace   = "QCE/CVM"
  metric_name = "CPUUsage"
  dimensions {
    name  = "InstanceId"
    value = data.tencentcloud_instances.instances.instance_list[0].instance_id
  }
  period     = 300
  start_time = "2020-04-28T18:45:00+08:00"
  end_time   = "2020-04-28T19:00:00+08:00"
}

#cos
data "tencentcloud_monitor_data" "cos_monitor_data" {
  namespace   = "QCE/COS"
  metric_name = "InternetTraffic"
  dimensions {
    name  = "appid"
    value = "1258798060"
  }
  dimensions {
    name  = "bucket"
    value = "test-1258798060"
  }

  period     = 300
  start_time = "2020-04-28T18:30:00+08:00"
  end_time   = "2020-04-28T19:00:00+08:00"
}
```