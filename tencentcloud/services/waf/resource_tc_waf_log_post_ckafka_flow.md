Provides a resource to create a WAF log post ckafka flow

Example Usage

If vip_type is 1

```hcl
resource "tencentcloud_waf_log_post_ckafka_flow" "example" {
  ckafka_region = "ap-guangzhou"
  ckafka_id     = "ckafka-k9m5vwar"
  brokers       = "ckafka-k9m5vwar.ap-guangzhou.ckafka.tencentcloudmq.com:50000"
  compression   = "snappy"
  vip_type      = 1
  log_type      = 2
  topic         = "tf-test"
  kafka_version = "2.8.1"
}
```

If vip_type is 2

```hcl
resource "tencentcloud_waf_log_post_ckafka_flow" "example" {
  ckafka_region = "ap-guangzhou"
  ckafka_id     = "ckafka-k9m5vwar"
  brokers       = "11.135.14.110:18737"
  compression   = "snappy"
  vip_type      = 2
  log_type      = 1
  topic         = "tf-example"
  kafka_version = "2.8.1"
}
```

Import

WAF log post ckafka flow can be imported using the id, e.g.

```
# If log_type is 1
terraform import tencentcloud_waf_log_post_ckafka_flow.example 100536#1

# If log_type is 2
terraform import tencentcloud_waf_log_post_ckafka_flow.example 100541#2
```
