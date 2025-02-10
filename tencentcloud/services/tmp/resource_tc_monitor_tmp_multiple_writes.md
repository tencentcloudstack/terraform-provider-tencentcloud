Provides a resource to create a monitor multiple writes

~> **NOTE:** For the same instance of prometheus, resource `tencentcloud_monitor_tmp_multiple_writes` and resource `tencentcloud_monitor_tmp_multiple_writes_list` cannot be used simultaneously. Resource `tencentcloud_monitor_tmp_multiple_writes` will been deprecated in version v1.81.166, Please use resource `tencentcloud_monitor_tmp_multiple_writes_list` instead.

~> **NOTE:** When using `<<EOT`, please pay attention to spaces, line breaks, indentation, etc.

~> **NOTE:** When importing, the unique id is separated by the first `#`.

Example Usage

```hcl
resource "tencentcloud_monitor_tmp_multiple_writes" "example" {
  instance_id = "prom-l9cl1ptk"

  remote_writes {
    url = "http://172.16.0.111:9090/api/v1/prom/write"
    url_relabel_config = trimspace(<<-EOT
        # 添加 label
        # - target_label: key
        #  replacement: value
        # 丢弃指标
        #- source_labels: [__name__]
        #  regex: kubelet_.+;
        #  action: drop
      EOT
    )
    headers {
      key   = "Key"
      value = "Value"
    }
  }
}
```

Import

monitor multiple writes can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_tmp_multiple_writes.example prom-l9cl1ptk#http://172.16.0.111:9090/api/v1/prom/write
```
