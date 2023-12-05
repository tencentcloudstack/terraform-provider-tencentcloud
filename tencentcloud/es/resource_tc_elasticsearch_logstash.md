Provides a resource to create a elasticsearch logstash

Example Usage

```hcl
resource "tencentcloud_elasticsearch_logstash" "logstash" {
  instance_name    = "logstash-test"
  zone             = "ap-guangzhou-6"
  logstash_version = "7.14.2"
  vpc_id           = "vpc-4owdpnwr"
  subnet_id        = "subnet-4o0zd840"
  node_num         = 1
  charge_type      = "POSTPAID_BY_HOUR"
  node_type        = "LOGSTASH.SA2.MEDIUM4"
  disk_type        = "CLOUD_SSD"
  disk_size        = 20
  license_type     = "xpack"
  operation_duration {
    periods    = [1, 2, 3, 4, 5, 6, 0]
    time_start = "02:00"
    time_end   = "06:00"
    time_zone  = "UTC+8"
  }
  tags = {
    tagKey = "tagValue"
  }
}
```

Import

elasticsearch logstash can be imported using the id, e.g.

```
terraform import tencentcloud_elasticsearch_logstash.logstash logstash_id
```