resource "tencentcloud_security_group" "foo" {
  name = "example"
}

resource "tencentcloud_vpc" "foo" {
  name       = "ci-test-eni-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "foo" {
  availability_zone = var.availability_zone
  name              = "ci-test-eni-subnet"
  vpc_id            = tencentcloud_vpc.foo.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_instance" "foo" {
  instance_name              = "example"
  availability_zone          = var.availability_zone
  image_id                   = "img-9qabwvbn"
  instance_type              = "S2.SMALL1"
  system_disk_type           = "CLOUD_PREMIUM"
  internet_max_bandwidth_out = 0
  vpc_id                     = tencentcloud_vpc.foo.id
  subnet_id                  = tencentcloud_subnet.foo.id
}

resource "tencentcloud_clb_instance" "example" {
  clb_name                  = "example"
  network_type              = var.network_type
  project_id                = 0
  vpc_id                    = tencentcloud_vpc.foo.id
  target_region_info_region = "ap-guangzhou"
  target_region_info_vpc_id = tencentcloud_vpc.foo.id
  security_groups           = [tencentcloud_security_group.foo.id]
}

resource "tencentcloud_clb_listener" "listener_tcp" {
  clb_id                     = tencentcloud_clb_instance.example.id
  listener_name              = "listener_tcp"
  port                       = 22
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
  health_check_port          = 200
  health_check_type          = "HTTP"
  health_check_http_code     = 2
  health_check_http_version  = "HTTP/1.0"
  health_check_http_method   = "GET"
}

resource "tencentcloud_clb_attachment" "attachment_tcp" {
  clb_id      = tencentcloud_clb_instance.example.id
  listener_id = tencentcloud_clb_listener.listener_tcp.listener_id

  targets {
    instance_id = tencentcloud_instance.foo.id
    port        = 22
    weight      = 10
  }
}

resource "tencentcloud_clb_listener" "listener_https" {
  clb_id               = tencentcloud_clb_instance.example.id
  listener_name        = "listener_https"
  port                 = 443
  protocol             = "HTTPS"
  certificate_ssl_mode = "UNIDIRECTIONAL"
  certificate_id       = "f8k7ke6a"
}

resource "tencentcloud_clb_listener_rule" "rule_https" {
  clb_id              = tencentcloud_clb_instance.example.id
  listener_id         = tencentcloud_clb_listener.listener_https.listener_id
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}

resource "tencentcloud_clb_attachment" "attachment_https" {
  clb_id      = tencentcloud_clb_instance.example.id
  listener_id = tencentcloud_clb_listener.listener_https.listener_id
  rule_id     = tencentcloud_clb_listener_rule.rule_https.rule_id

  targets {
    instance_id = tencentcloud_instance.foo.id
    port        = 443
    weight      = 10
  }
}

resource "tencentcloud_clb_listener" "listener_http_src" {
  clb_id        = tencentcloud_clb_instance.example.id
  port          = 8080
  protocol      = "HTTP"
  listener_name = "listener_http_src"
}

resource "tencentcloud_clb_listener_rule" "rule_http_src" {
  clb_id              = tencentcloud_clb_instance.example.id
  listener_id         = tencentcloud_clb_listener.listener_http_src.listener_id
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}

resource "tencentcloud_clb_listener" "listener_http_dst" {
  clb_id        = tencentcloud_clb_instance.example.id
  port          = 80
  protocol      = "HTTP"
  listener_name = "listener_http_dst"
}

resource "tencentcloud_clb_listener_rule" "rule_http_dst" {
  clb_id              = tencentcloud_clb_instance.example.id
  listener_id         = tencentcloud_clb_listener.listener_http_dst.listener_id
  domain              = "abcd.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}

resource "tencentcloud_clb_redirection" "redirection_http" {
  clb_id             = tencentcloud_clb_instance.example.id
  source_listener_id = tencentcloud_clb_listener.listener_http_src.listener_id
  target_listener_id = tencentcloud_clb_listener.listener_http_dst.listener_id
  source_rule_id     = tencentcloud_clb_listener_rule.rule_http_src.rule_id
  target_rule_id     = tencentcloud_clb_listener_rule.rule_http_dst.rule_id
}

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

data "tencentcloud_images" "my_favorite_image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}

data "tencentcloud_instance_types" "my_favorite_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S3"]
  }

  cpu_core_count = 1
  memory_size    = 1
}

data "tencentcloud_availability_zones" "default" {
}

resource "tencentcloud_vpc" "app" {
  cidr_block = "10.0.0.0/16"
  name       = "awesome_app_vpc"
}

resource "tencentcloud_subnet" "app" {
  vpc_id            = tencentcloud_vpc.app.id
  availability_zone = data.tencentcloud_availability_zones.default.zones.0.name
  name              = "awesome_app_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_instance" "my_awesome_app" {
  instance_name              = "awesome_app"
  availability_zone          = data.tencentcloud_availability_zones.default.zones.0.name
  image_id                   = data.tencentcloud_images.my_favorite_image.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.my_favorite_instance_types.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  hostname                   = "user"
  project_id                 = 0
  vpc_id                     = tencentcloud_vpc.app.id
  subnet_id                  = tencentcloud_subnet.app.id
  internet_max_bandwidth_out = 20

  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 50
    encrypt        = false
  }

  tags = {
    tagKey = "tagValue"
  }
}

data "tencentcloud_instances" "foo" {
  instance_id = tencentcloud_instance.my_awesome_app.id
}

resource "tencentcloud_clb_target_group" "test_instance_attachment"{
  target_group_name = "test"
  vpc_id            = tencentcloud_vpc.app.id
}

resource "tencentcloud_clb_target_group_instance_attachment" "test"{
  target_group_id = tencentcloud_clb_target_group.test_instance_attachment.id
  bind_ip         = data.tencentcloud_instances.foo.instance_list[0].private_ip
  port            = 88
  weight          = 3
}

data "tencentcloud_clb_target_groups" "target_group_info_id" {
  target_group_id = tencentcloud_clb_target_group.test.id
}

data "tencentcloud_clb_instances" "instances" {
  clb_id = tencentcloud_clb_instance.example.id
}

data "tencentcloud_clb_listeners" "listeners" {
  clb_id      = tencentcloud_clb_instance.example.id
  listener_id = tencentcloud_clb_listener.listener_tcp.listener_id
}

data "tencentcloud_clb_listener_rules" "rules" {
  clb_id      = tencentcloud_clb_instance.example.id
  listener_id = tencentcloud_clb_listener.listener_https.listener_id
  domain      = tencentcloud_clb_listener_rule.rule_https.domain
  url         = tencentcloud_clb_listener_rule.rule_https.url
}

data "tencentcloud_clb_attachments" "attachments" {
  clb_id      = tencentcloud_clb_instance.example.id
  listener_id = tencentcloud_clb_listener.listener_https.listener_id
  rule_id     = tencentcloud_clb_attachment.attachment_https.rule_id
}

data "tencentcloud_clb_redirections" "redirections" {
  clb_id             = tencentcloud_clb_instance.example.id
  source_listener_id = tencentcloud_clb_redirection.redirection_http.source_listener_id
  source_rule_id     = tencentcloud_clb_redirection.redirection_http.source_rule_id
}
