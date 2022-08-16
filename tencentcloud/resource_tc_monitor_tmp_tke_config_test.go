package tencentcloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/pkg/errors"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_monitor_tmp_tke_config
	resource.AddTestSweepers("tencentcloud_monitor_tmp_tke_config", &resource.Sweeper{
		Name: "tencentcloud_monitor_tmp_tke_config",
		F: func(r string) error {
			logId := getLogId(contextNil)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn
			configId := packConfigId(defaultPrometheusId, defaultTkeClusterType, defaultTkeClusterId)

			service := TkeService{client}

			promConfigs, err := service.DescribeTkeTmpConfigById(logId, configId)

			if err != nil {
				return err
			}

			if promConfigs == nil {
				return fmt.Errorf("Prometheus config %s not exist", configId)
			}

			ServiceMonitors := transObj2StrNames(promConfigs.ServiceMonitors)
			PodMonitors := transObj2StrNames(promConfigs.PodMonitors)
			RawJobs := transObj2StrNames(promConfigs.RawJobs)
			service.DeleteTkeTmpConfigByName(logId, configId, ServiceMonitors, PodMonitors, RawJobs)

			return nil
		},
	})
}

func TestAccTencentCloudTmpTkeConfig_basic(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTmpTkeConfigDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccTmpTkeConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTmpTkeConfigExists("tencentcloud_monitor_tmp_tke_config.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_config.foo", "instance_id", defaultPrometheusId),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_config.foo", "cluster_type", defaultTkeClusterType),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_config.foo", "cluster_id", defaultTkeClusterId),
				),
			},
			{
				Config: testAccTmpTkeConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTmpTkeConfigExists("tencentcloud_monitor_tmp_tke_config.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_config.foo", "instance_id", defaultPrometheusId),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_config.foo", "cluster_type", defaultTkeClusterType),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_config.foo", "cluster_id", defaultTkeClusterId),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_tmp_tke_config.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTmpTkeConfigDestroy(configId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		client := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		service := TkeService{client}

		// rs, ok := s.RootModule().Resources[n]
		// if !ok {
		// 	return fmt.Errorf("not found: %s", n)
		// }

		// configId := rs.Primary.ID
		// if configId == "" {
		// 	return errors.New("no prometheus config ID is set")
		// }

		promConfigs, err := service.DescribeTkeTmpConfigById(logId, *configId)

		if err != nil {
			return err
		}

		if promConfigs == nil {
			return nil
		}

		if len(promConfigs.ServiceMonitors) != 0 || len(promConfigs.PodMonitors) != 0 || len(promConfigs.RawJobs) != 0 {
			return errors.New("promConfigs still exists")
		}

		return nil
	}
}

func testAccCheckTmpTkeConfigExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		client := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		service := TkeService{client}

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		configId := rs.Primary.ID
		if configId == "" {
			return errors.New("no prometheus config ID is set")
		}

		promConfigs, err := service.DescribeTkeTmpConfigById(logId, configId)
		if err != nil {
			return err
		}

		if promConfigs == nil || (len(promConfigs.ServiceMonitors) == 0 && len(promConfigs.PodMonitors) == 0 && len(promConfigs.RawJobs) == 0) {
			return fmt.Errorf("prometheus config not found: %s", rs.Primary.ID)
		}
		*id = configId

		return nil
	}
}

func packConfigId(instanceId string, clusterType string, clusterId string) (ids string) {
	ids = strings.Join([]string{instanceId, clusterType, clusterId}, FILED_SP)
	return
}

func transObj2StrNames(resList []*tke.PrometheusConfigItem) []*string {
	names := make([]*string, 0, len(resList))
	for _, res := range resList {
		if res.Name != nil {
			names = append(names, res.Name)
		}
	}
	return names
}

const testAccTmpTkeConfig_basic = `
resource "tencentcloud_monitor_tmp_tke_config" "foo" {
  instance_id  = "` + defaultPrometheusId + `"
  cluster_type = "` + defaultTkeClusterType + `"
  cluster_id   = "` + defaultTkeClusterId + `"
  raw_jobs {
    name   = "rawjob_test_001"
    config = "scrape_configs:\n- job_name: rawjob_test_001\n  honor_labels: true\n  honor_timestamps: false\n  scrape_interval: 20s\n  metrics_path: /metrics\n  scheme: http\n"
  }
}`
