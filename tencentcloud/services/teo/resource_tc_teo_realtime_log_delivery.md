Provides a resource to create a teo teo_realtime_log_delivery

Example Usage

```hcl
resource "tencentcloud_teo_realtime_log_delivery" "teo_realtime_log_delivery" {
    area            = "overseas"
    delivery_status = "disabled"
    entity_list     = [
        "sid-2yvhjw98uaco",
    ]
    fields          = [
        "ServiceID",
        "ConnectTimeStamp",
        "DisconnetTimeStamp",
        "DisconnetReason",
        "ClientRealIP",
        "ClientRegion",
        "EdgeIP",
        "ForwardProtocol",
        "ForwardPort",
        "SentBytes",
        "ReceivedBytes",
        "LogTimeStamp",
    ]
    log_type        = "application"
    sample          = 0
    task_name       = "test"
    task_type       = "s3"
    zone_id         = "zone-2qtuhspy7cr6"

    log_format {
        field_delimiter  = ","
        format_type      = "json"
        record_delimiter = "\n"
        record_prefix    = "{"
        record_suffix    = "}"
    }

    s3 {
        access_id     = "xxxxxxxxxx"
        access_key    = "xxxxxxxxxx"
        bucket        = "test-1253833068"
        compress_type = "gzip"
        endpoint      = "https://test-1253833068.cos.ap-nanjing.myqcloud.com"
        region        = "ap-nanjing"
    }
}
```

Import

teo teo_realtime_log_delivery can be imported using the id, e.g.

```
terraform import tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery zoneId#taskId
```
