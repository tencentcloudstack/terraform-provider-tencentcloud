Provides a resource for an AS (Auto scaling) lifecycle hook.

Example Usage

Create a basic LifecycleHook

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "as"
}

data "tencentcloud_images" "image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "TencentOS Server 3.2 (Final)"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-example"
  cidr_block        = "10.0.0.0/16"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
}

resource "tencentcloud_as_scaling_config" "example" {
  configuration_name = "tf-example"
  image_id           = data.tencentcloud_images.image.images.0.image_id
  instance_types     = ["SA1.SMALL1", "SA2.SMALL1", "SA2.SMALL2", "SA2.SMALL4"]
  instance_name_settings {
    instance_name = "test-ins-name"
  }
}

resource "tencentcloud_as_scaling_group" "example" {
  scaling_group_name = "tf-example"
  configuration_id   = tencentcloud_as_scaling_config.example.id
  max_size           = 1
  min_size           = 0
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_ids         = [tencentcloud_subnet.subnet.id]
}

resource "tencentcloud_as_lifecycle_hook" "example" {
  scaling_group_id         = tencentcloud_as_scaling_group.example.id
  lifecycle_hook_name      = "tf-as-lifecycle-hook"
  lifecycle_transition     = "INSTANCE_LAUNCHING"
  default_result           = "CONTINUE"
  heartbeat_timeout        = 500
  notification_metadata    = "tf test"
}
```

If `notification_target_type` is `CMQ_QUEUE`

```hcl
resource "tencentcloud_as_lifecycle_hook" "example" {
  scaling_group_id         = tencentcloud_as_scaling_group.example.id
  lifecycle_hook_name      = "tf-as-lifecycle-hook"
  lifecycle_transition     = "INSTANCE_LAUNCHING"
  default_result           = "CONTINUE"
  heartbeat_timeout        = 500
  notification_metadata    = "tf test"
  notification_target_type = "CMQ_QUEUE"
  notification_queue_name  = "lifcyclehook"
}
```

Or `notification_target_type` is `CMQ_TOPIC`

```hcl
resource "tencentcloud_as_lifecycle_hook" "example" {
  scaling_group_id         = tencentcloud_as_scaling_group.example.id
  lifecycle_hook_name      = "tf-as-lifecycle-hook"
  lifecycle_transition     = "INSTANCE_LAUNCHING"
  default_result           = "CONTINUE"
  heartbeat_timeout        = 500
  notification_metadata    = "tf test"
  notification_target_type = "CMQ_TOPIC"
  notification_topic_name  = "lifcyclehook"
}
```