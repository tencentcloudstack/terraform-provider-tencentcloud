package tke_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctke "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tke"

	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_eks_ci
	resource.AddTestSweepers("tencentcloud_eks_ci", &resource.Sweeper{
		Name: "tencentcloud_eks_ci",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
			service := svctke.NewEksService(client)

			cis, err := service.DescribeEksContainerInstancesByFilter(ctx, nil, 100, 0)

			if err != nil {
				return err
			}

			var ids []*string
			for i := range cis {
				ci := cis[i]
				name := *ci.EksCiName
				if tcacctest.IsResourcePersist(name, nil) {
					continue
				}
				ids = append(ids, ci.EksCiId)
			}
			request := tke.NewDeleteEKSContainerInstancesRequest()
			request.EksCiIds = ids
			err = service.DeleteEksContainerInstance(ctx, request)

			if err != nil {
				return err
			}

			return nil
		},
	})
}

func TestAccTencentCloudEKSContainerInstance_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckEksCiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEksCi,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("tencentcloud_eks_container_instance.foo"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "name", "foo"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "cpu", "2"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "memory", "4"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "cpu_type", "intel"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "restart_policy", "Always"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "security_groups.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "container.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "container.0.name", "nginx"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "container.0.image", "nginx"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "container.0.liveness_probe.0.init_delay_seconds", "1"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "container.0.liveness_probe.0.timeout_seconds", "3"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "container.0.liveness_probe.0.period_seconds", "10"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "container.0.liveness_probe.0.success_threshold", "1"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "container.0.liveness_probe.0.failure_threshold", "3"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "container.0.liveness_probe.0.http_get_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "container.0.liveness_probe.0.http_get_port", "80"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "container.0.liveness_probe.0.http_get_scheme", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "container.0.readiness_probe.0.init_delay_seconds", "1"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "container.0.readiness_probe.0.timeout_seconds", "3"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "container.0.readiness_probe.0.period_seconds", "10"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "container.0.readiness_probe.0.success_threshold", "1"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "container.0.readiness_probe.0.failure_threshold", "3"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "container.0.readiness_probe.0.tcp_socket_port", "81"),
					resource.TestCheckResourceAttr("tencentcloud_eks_container_instance.foo", "init_container.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_eks_container_instance.foo", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_eks_container_instance.foo", "subnet_id"),
				),
			},
		},
	})
}

func testAccCheckEksCiDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	eksService := svctke.NewEksService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_eks_container_instance" {
			continue
		}
		_, has, err := eksService.DescribeEksContainerInstanceById(ctx, rs.Primary.ID)

		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				_, has, err = eksService.DescribeEksContainerInstanceById(ctx, rs.Primary.ID)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
		}

		if err != nil {
			return err
		}

		if has {
			return fmt.Errorf("eks container instance still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

const testAccEksCi = tcacctest.DefaultVpcVariable + `
data "tencentcloud_security_groups" "group" {
  name = "default"
}

resource "tencentcloud_eks_container_instance" "foo" {
  name = "foo"
  vpc_id = var.vpc_id
  subnet_id = var.subnet_id
  cpu = 2
  cpu_type = "intel"
  restart_policy = "Always"
  memory = 4
  security_groups = [data.tencentcloud_security_groups.group.security_groups[0].security_group_id]
  container {
    name = "nginx"
    image = "nginx"
    liveness_probe {
      init_delay_seconds = 1
      timeout_seconds = 3
      period_seconds = 10
      success_threshold = 1
      failure_threshold = 3
      http_get_path = "/"
      http_get_port = 80
      http_get_scheme = "HTTP"
    }
    readiness_probe {
      init_delay_seconds = 1
      timeout_seconds = 3
      period_seconds = 10
      success_threshold = 1
      failure_threshold = 3
      tcp_socket_port = 81
    }
  }
  init_container {
    name = "alpine"
    image = "alpine"
  }
}`
