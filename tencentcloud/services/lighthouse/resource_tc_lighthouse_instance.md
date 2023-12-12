Provides a resource to create a lighthouse instance.

Example Usage

```hcl
resource "tencentcloud_lighthouse_firewall_template" "firewall_template" {
  template_name="empty-template"
}

resource "tencentcloud_lighthouse_instance" "lighthouse" {
  bundle_id    = "bundle2022_gen_01"
  blueprint_id = "lhbp-f1lkcd41"

  period     = 1
  renew_flag = "NOTIFY_AND_AUTO_RENEW"

  instance_name = "hello world"
  zone          = "ap-guangzhou-3"

  containers {
    container_image = "ccr.ccs.tencentyun.com/qcloud/nginx"
    container_name = "nginx"
    envs {
      key = "key"
      value = "value"
    }
    envs {
      key = "key2"
      value = "value2"
    }
    publish_ports {
      host_port = 80
      container_port = 80
      ip = "127.0.0.1"
      protocol = "tcp"
    }
    publish_ports {
      host_port = 36000
      container_port = 36000
      ip = "127.0.0.1"
      protocol = "tcp"
    }
    volumes {
      container_path = "/data"
      host_path = "/tmp"
    }
    volumes {
      container_path = "/var"
      host_path = "/tmp"
    }
    command = "ls -l"
  }

  containers {
    container_image = "ccr.ccs.tencentyun.com/qcloud/resty"
    container_name = "resty"
    envs {
      key = "key2"
      value = "value2"
    }
    publish_ports {
      host_port = 80
      container_port = 80
      ip = "127.0.0.1"
      protocol = "udp"
    }

    volumes {
      container_path = "/var"
      host_path = "/tmp"
    }
    command = "echo \"hello\""
  }
  firewall_template_id = tencentcloud_lighthouse_firewall_template.firewall_template.id
}
```