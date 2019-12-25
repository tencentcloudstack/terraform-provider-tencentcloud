resource "tencentcloud_vpc" "vpc" {
  name       = "tf-as-vpc"
  cidr_block = "10.2.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "tf-as-subnet"
  cidr_block        = "10.2.11.0/24"
  availability_zone = var.availability_zone
}

resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "tf-as-configuration"
  image_id           = "img-9qabwvbn"
  instance_types     = [var.instance_type]
  project_id         = 0
  system_disk_type   = "CLOUD_PREMIUM"
  system_disk_size   = "50"

  data_disk {
    disk_type = "CLOUD_PREMIUM"
    disk_size = 50
  }

  internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
  internet_max_bandwidth_out = 10
  public_ip_assigned         = true
  password                   = "test123#"
  enhanced_security_service  = false
  enhanced_monitor_service   = false
  user_data                  = "test"

  instance_tags = {
    tag = "as"
  }
}

resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name   = "tf-as-scaling-group"
  configuration_id     = tencentcloud_as_scaling_config.launch_configuration.id
  max_size             = var.max_size
  min_size             = var.min_size
  vpc_id               = tencentcloud_vpc.vpc.id
  subnet_ids           = [tencentcloud_subnet.subnet.id]
  project_id           = 0
  default_cooldown     = 400
  desired_capacity     = var.desired_capacity
  termination_policies = ["NEWEST_INSTANCE"]
  retry_policy         = "INCREMENTAL_INTERVALS"

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_as_scaling_policy" "scaling_policy" {
  scaling_group_id    = tencentcloud_as_scaling_group.scaling_group.id
  policy_name         = "tf-as-scaling-policy"
  adjustment_type     = "EXACT_CAPACITY"
  adjustment_value    = 2
  comparison_operator = "GREATER_THAN"
  metric_name         = "CPU_UTILIZATION"
  threshold           = 80
  period              = 300
  continuous_time     = 10
  statistic           = "AVERAGE"
  cooldown            = 360
}

resource "tencentcloud_as_schedule" "schedule" {
  scaling_group_id     = tencentcloud_as_scaling_group.scaling_group.id
  schedule_action_name = "tf-as-schedule"
  max_size             = var.max_size
  min_size             = var.min_size
  desired_capacity     = var.desired_capacity
  start_time           = "2020-01-01T00:00:00+08:00"
  end_time             = "2020-12-01T00:00:00+08:00"
  recurrence           = "0 0 */1 * *"
}

resource "tencentcloud_as_lifecycle_hook" "lifecycle_hook" {
  scaling_group_id         = tencentcloud_as_scaling_group.scaling_group.id
  lifecycle_hook_name      = "tf-as-lifecycle"
  lifecycle_transition     = "INSTANCE_TERMINATING"
  default_result           = "ABANDON"
  heartbeat_timeout        = 300
  notification_metadata    = "tf lifecycle test"
  notification_target_type = "CMQ_QUEUE"
  notification_queue_name  = "notification"
}

resource "tencentcloud_as_notification" "notification" {
  scaling_group_id            = tencentcloud_as_scaling_group.scaling_group.id
  notification_types          = ["SCALE_OUT_FAILED"]
  notification_user_group_ids = ["76955"]
}

data "tencentcloud_as_scaling_groups" "scaling_groups_tags" {
  tags = tencentcloud_as_scaling_group.scaling_group.tags
}
