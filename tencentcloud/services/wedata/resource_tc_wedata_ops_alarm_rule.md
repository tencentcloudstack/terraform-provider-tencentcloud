Provides a resource to create a wedata ops alarm rule

Example Usage

```hcl
resource "tencentcloud_wedata_ops_alarm_rule" "wedata_ops_alarm_rule" {
    alarm_level         = 1
    alarm_rule_name     = "tf_test"
    alarm_types         = [
        "failure",
    ]
    description         = "ccc"
    monitor_object_ids  = [
        "20230906105118824",
    ]
    monitor_object_type = 1
    project_id          = "1859317240494305280"

    alarm_groups {
        alarm_escalation_interval      = 15
        alarm_escalation_recipient_ids = []
        alarm_recipient_ids            = [
            "100029411056",
        ]
        alarm_recipient_type           = 1
        alarm_ways                     = [
            "1",
        ]

        notification_fatigue {
            notify_count    = 1
            notify_interval = 5

            quiet_intervals {
                days_of_week = [
                    6,
                    7,
                ]
                end_time     = "21:00:00"
                start_time   = "10:00:00"
            }
        }
    }

    alarm_rule_detail {
        data_backfill_or_rerun_trigger = 1
        trigger                        = 2
    }
}
```


Import

wedata ops alarm rule can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule projectId#askId
```