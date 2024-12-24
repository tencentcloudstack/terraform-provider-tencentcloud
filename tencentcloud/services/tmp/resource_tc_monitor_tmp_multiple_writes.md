Provides a resource to create a monitor monitor_tmp_multiple_writes

~> **NOTE:** When using `<<EOT`, please pay attention to spaces, line breaks, indentation, etc.

~> **NOTE:** When importing, the unique id is separated by the first `#`.

Example Usage

```hcl
resource "tencentcloud_monitor_tmp_multiple_writes" "monitor_tmp_multiple_writes" {
    instance_id = "prom-l9cl1ptk"

    remote_writes {
        label              = null
        max_block_size     = null
        url                = "http://172.16.0.111:9090/api/v1/prom/write"
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
    }
}
```

Import

monitor monitor_tmp_multiple_writes can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_tmp_multiple_writes.monitor_tmp_multiple_writes prom-l9cl1ptk#http://172.16.0.111:9090/api/v1/prom/write
```
