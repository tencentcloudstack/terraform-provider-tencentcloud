Provides a resource to create a CLB target group attachment is bound to the load balancing listener or forwarding rule.

Example Usage

```hcl
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-rule-basic"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id        = tencentcloud_clb_instance.clb_basic.id
  port          = 1
  protocol      = "HTTP"
  listener_name = "listener_basic"
}

resource "tencentcloud_clb_listener_rule" "rule_basic" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  listener_id         = tencentcloud_clb_listener.listener_basic.listener_id
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}

resource "tencentcloud_clb_target_group" "test"{
    target_group_name = "test-target-keep-1"
}

resource "tencentcloud_clb_target_group_attachment" "group" {
    clb_id          = tencentcloud_clb_instance.clb_basic.id
    listener_id     = tencentcloud_clb_listener.listener_basic.listener_id
    rule_id         = tencentcloud_clb_listener_rule.rule_basic.rule_id
    target_group_id = tencentcloud_clb_target_group.test.id
}
```

Import

CLB target group attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_target_group_attachment.group lbtg-odareyb2#lbl-bicjmx3i#lb-cv0iz74c#loc-ac6uk7b6
```