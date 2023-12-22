package lighthouse_test

import (
	"context"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	svclighthouse "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/lighthouse"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_lighthouse_instance
	resource.AddTestSweepers("tencentcloud_lighthouse_instance", &resource.Sweeper{
		Name: "tencentcloud_lighthouse_instance",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)

			request := lighthouse.NewDescribeInstancesRequest()
			response, err := cli.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().DescribeInstances(request)
			if err != nil {
				return err
			}
			instances := response.Response.InstanceSet
			service := svclighthouse.NewLightHouseService(cli.(tccommon.ProviderMeta).GetAPIV3Conn())

			for _, instance := range instances {
				name := *instance.InstanceName
				created, err := time.Parse("2006-01-02 15:04:05", *instance.CreatedTime)
				if err != nil {
					continue
				}
				if tcacctest.IsResourcePersist(name, &created) {
					continue
				}
				if innerErr := service.DeleteLighthouseInstanceById(ctx, *instance.InstanceId); innerErr != nil {
					continue
				}
			}
			return nil
		},
	})
}

func TestAccTencentCloudLighthouseInstanceResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_lighthouse_instance.instance", "instance_name", "terraform"),
				),
			},
			{
				Config: testAccLighthouseInstance_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_lighthouse_instance.instance", "renew_flag", "NOTIFY_AND_MANUAL_RENEW"),
				),
			},
		},
	})
}

const testAccLighthouseInstance = `
data "tencentcloud_lighthouse_bundle" "bundle" {
}

resource "tencentcloud_lighthouse_firewall_template" "firewall_template" {
  template_name="empty-template"
}

resource "tencentcloud_lighthouse_instance" "instance" {
  bundle_id    = data.tencentcloud_lighthouse_bundle.bundle.bundle_set.0.bundle_id
  blueprint_id = "lhbp-f1lkcd41"

  period     = 1
  renew_flag = "NOTIFY_AND_AUTO_RENEW"

  instance_name = "terraform"
  zone          = "ap-guangzhou-3"
  isolate_data_disk = true
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
`

const testAccLighthouseInstance_update = `
data "tencentcloud_lighthouse_bundle" "bundle" {
}

resource "tencentcloud_lighthouse_firewall_template" "firewall_template" {
  template_name="empty-template"
}

resource "tencentcloud_lighthouse_instance" "instance" {
  bundle_id    = data.tencentcloud_lighthouse_bundle.bundle.bundle_set.1.bundle_id
  blueprint_id = "lhbp-f1lkcd41"

  period     = 1
  renew_flag = "NOTIFY_AND_MANUAL_RENEW"

  instance_name = "terraform"
  zone          = "ap-guangzhou-3"
  isolate_data_disk = true

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
`
