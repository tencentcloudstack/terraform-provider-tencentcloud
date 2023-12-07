Provides a resource to create a clickhouse instance.

Example Usage

```hcl
resource "tencentcloud_clickhouse_instance" "cdwch_instance" {
  zone="ap-guangzhou-6"
  ha_flag=true
  vpc_id="vpc-xxxxxx"
  subnet_id="subnet-xxxxxx"
  product_version="21.8.12.29"
  data_spec {
    spec_name="SCH6"
    count=2
    disk_size=300
  }
  common_spec {
    spec_name="SCH6"
    count=3
    disk_size=300
  }
  charge_type="POSTPAID_BY_HOUR"
  instance_name="tf-test-clickhouse"
}
```

PREPAID instance

```hcl
resource "tencentcloud_clickhouse_instance" "cdwch_instance_prepaid" {
  zone="ap-guangzhou-6"
  ha_flag=true
  vpc_id="vpc-xxxxxx"
  subnet_id="subnet-xxxxxx"
  product_version="21.8.12.29"
  data_spec {
    spec_name="SCH6"
    count=2
    disk_size=300
  }
  common_spec {
    spec_name="SCH6"
    count=3
    disk_size=300
  }
  charge_type="PREPAID"
  renew_flag=1
  time_span=1
  instance_name="tf-test-clickhouse-prepaid"
}
```

Import

Clickhouse instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_clickhouse_instance.foo cdwch-xxxxxx
```