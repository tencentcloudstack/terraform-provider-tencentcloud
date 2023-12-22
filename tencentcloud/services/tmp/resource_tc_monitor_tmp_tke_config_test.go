package tmp_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pkg/errors"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_monitor_tmp_tke_config
	resource.AddTestSweepers("tencentcloud_monitor_tmp_tke_config", &resource.Sweeper{
		Name: "tencentcloud_monitor_tmp_tke_config",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
			configId := packConfigId(tcacctest.DefaultPrometheusId, tcacctest.DefaultTkeClusterType, tcacctest.DefaultTkeClusterId)

			service := svcmonitor.NewMonitorService(client)

			promConfigs, err := service.DescribeTkeTmpConfigById(ctx, configId)

			if err != nil {
				return err
			}

			if promConfigs == nil {
				return fmt.Errorf("Prometheus config %s not exist", configId)
			}

			ServiceMonitors := transObj2StrNames(promConfigs.ServiceMonitors)
			PodMonitors := transObj2StrNames(promConfigs.PodMonitors)
			RawJobs := transObj2StrNames(promConfigs.RawJobs)
			err = service.DeleteTkeTmpConfigByName(ctx, configId, ServiceMonitors, PodMonitors, RawJobs)
			if err != nil {
				return err
			}

			return nil
		},
	})
}

func TestAccTencentCloudMonitorTmpTkeConfig_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTmpTkeConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTmpTkeConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTmpTkeConfigExists("tencentcloud_monitor_tmp_tke_config.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_config.basic", "instance_id", tcacctest.DefaultPrometheusId),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_config.basic", "cluster_type", tcacctest.DefaultTkeClusterType),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_config.basic", "cluster_id", tcacctest.DefaultTkeClusterId),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_config.basic", "raw_jobs.0.name", raw_job_name),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_config.basic", "service_monitors.0.name", service_monitors_name_fully),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_config.basic", "pod_monitors.0.name", pod_monitors_name_fully),
				),
			},
			{
				Config: testAccTmpTkeConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTmpTkeConfigExists("tencentcloud_monitor_tmp_tke_config.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_config.basic", "instance_id", tcacctest.DefaultPrometheusId),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_config.basic", "cluster_type", tcacctest.DefaultTkeClusterType),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_config.basic", "cluster_id", tcacctest.DefaultTkeClusterId),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_config.basic", "raw_jobs.0.config", "scrape_configs:\n- job_name: "+raw_job_name+"\n  scrape_interval: 20s\n  honor_labels: true\n"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_config.basic", "service_monitors.0.config", "apiVersion: monitoring.coreos.com/v1\nkind: ServiceMonitor\nmetadata:\n  name: "+service_monitors_name+"\n  namespace: kube-system\nspec:\n  endpoints:\n    - interval: 20s\n      port: 8080-8080-tcp\n      path: /metrics\n      relabelings:\n        - action: replace\n          sourceLabels:\n            - __meta_kubernetes_pod_label_app\n          targetLabel: application\n  namespaceSelector:\n    matchNames:\n      - test\n  selector:\n    matchLabels:\n      app: test"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_config.basic", "pod_monitors.0.config", "apiVersion: monitoring.coreos.com/v1\nkind: PodMonitor\nmetadata:\n  name: "+pod_monitors_name+"\n  namespace: kube-system\nspec:\n  podMetricsEndpoints:\n    - interval: 20s\n      port: metric-port\n      path: /metrics\n      relabelings:\n        - action: replace\n          sourceLabels:\n            - instance\n          regex: (.*)\n          targetLabel: instance\n          replacement: xxxxxx\n  namespaceSelector:\n    matchNames:\n      - test\n  selector:\n    matchLabels:\n      k8s-app: test"),
				),
			},
		},
	})
}

func testAccCheckTmpTkeConfigDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_tmp_tke_config" {
			continue
		}

		promConfigs, err := service.DescribeTkeTmpConfigById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if promConfigs == nil {
			return nil
		}

		for _, config := range promConfigs.ServiceMonitors {
			if *config.Name == service_monitors_name_fully {
				return errors.New("promConfigs service_monitors still exists")
			}
		}

		for _, config := range promConfigs.PodMonitors {
			if *config.Name == pod_monitors_name_fully {
				return errors.New("promConfigs pod_monitors still exists")
			}
		}

		for _, config := range promConfigs.RawJobs {
			if *config.Name == raw_job_name {
				return errors.New("promConfigs raw_jobs still exists")
			}
		}
	}
	return nil
}

func testAccCheckTmpTkeConfigExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("instance id is not set")
		}

		tkeService := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		promConfigs, err := tkeService.DescribeTkeTmpConfigById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if promConfigs == nil || (len(promConfigs.ServiceMonitors) == 0 && len(promConfigs.PodMonitors) == 0 && len(promConfigs.RawJobs) == 0) {
			return fmt.Errorf("prometheus config not found: %s", rs.Primary.ID)
		}

		return nil
	}
}

func packConfigId(instanceId string, clusterType string, clusterId string) (ids string) {
	ids = strings.Join([]string{instanceId, clusterType, clusterId}, tccommon.FILED_SP)
	return
}

func transObj2StrNames(resList []*monitor.PrometheusConfigItem) []*string {
	names := make([]*string, 0, len(resList))
	for _, res := range resList {
		if res.Name != nil {
			names = append(names, res.Name)
		}
	}
	return names
}

const (
	raw_job_name                = "raw_jobs_001"
	pod_monitors_name           = "pod-monitor-001"
	service_monitors_name       = "service-monitor-001"
	pod_monitors_name_fully     = "kube-system/pod-monitor-001"
	service_monitors_name_fully = "kube-system/service-monitor-001"
)

const testAccTmpTkeConfigVar = `
variable "prometheus_id" {
  default = "` + tcacctest.DefaultPrometheusId + `"
}

variable "tke_cluster_type" {
  default = "` + tcacctest.DefaultTkeClusterType + `"
}

variable "tke_cluster_id" {
  default = "` + tcacctest.DefaultTkeClusterId + `"
}

variable "pod_monitors_name_fully" {
  default = "` + pod_monitors_name_fully + `"
}
  
variable "service_monitors_name_fully" {
  default = "` + service_monitors_name_fully + `"
}

variable "raw_job_name" {
  default = "` + raw_job_name + `"
}

variable "pod_monitors_name" {
  default = "` + pod_monitors_name + `"
}

variable "service_monitors_name" {
  default = "` + service_monitors_name + `"
}
`

const testAccTmpTkeConfig_basic = testAccTmpTkeConfigVar + `
resource "tencentcloud_monitor_tmp_tke_config" "basic" {
  instance_id  = var.prometheus_id
  cluster_type = var.tke_cluster_type
  cluster_id   = var.tke_cluster_id
  service_monitors {
    name   = var.service_monitors_name_fully
    config = "apiVersion: monitoring.coreos.com/v1\nkind: ServiceMonitor\nmetadata:\n  name: ` + service_monitors_name + `\n  namespace: kube-system\nspec:\n  endpoints:\n    - interval: 115s\n      port: 8080-8080-tcp\n      path: /metrics\n      relabelings:\n        - action: replace\n          sourceLabels:\n            - __meta_kubernetes_pod_label_app\n          targetLabel: application\n  namespaceSelector:\n    matchNames:\n      - test\n  selector:\n    matchLabels:\n      app: test"
  }

  pod_monitors {
    name   = var.pod_monitors_name_fully
    config = "apiVersion: monitoring.coreos.com/v1\nkind: PodMonitor\nmetadata:\n  name: ` + pod_monitors_name + `\n  namespace: kube-system\nspec:\n  podMetricsEndpoints:\n    - interval: 15s\n      port: metric-port\n      path: /metrics\n      relabelings:\n        - action: replace\n          sourceLabels:\n            - instance\n          regex: (.*)\n          targetLabel: instance\n          replacement: xxxxxx\n  namespaceSelector:\n    matchNames:\n      - test\n  selector:\n    matchLabels:\n      k8s-app: test"
  }

  raw_jobs {
    name   = var.raw_job_name
    config = "scrape_configs:\n- job_name: ` + raw_job_name + `\n  honor_labels: true\n"
  }
}`

const testAccTmpTkeConfig_update = testAccTmpTkeConfigVar + `
resource "tencentcloud_monitor_tmp_tke_config" "basic" {
  instance_id  = var.prometheus_id
  cluster_type = var.tke_cluster_type
  cluster_id   = var.tke_cluster_id
  service_monitors {
    name   = var.service_monitors_name_fully
    config = "apiVersion: monitoring.coreos.com/v1\nkind: ServiceMonitor\nmetadata:\n  name: ` + service_monitors_name + `\n  namespace: kube-system\nspec:\n  endpoints:\n    - interval: 20s\n      port: 8080-8080-tcp\n      path: /metrics\n      relabelings:\n        - action: replace\n          sourceLabels:\n            - __meta_kubernetes_pod_label_app\n          targetLabel: application\n  namespaceSelector:\n    matchNames:\n      - test\n  selector:\n    matchLabels:\n      app: test"
  }

  pod_monitors {
    name   = var.pod_monitors_name_fully
    config = "apiVersion: monitoring.coreos.com/v1\nkind: PodMonitor\nmetadata:\n  name: ` + pod_monitors_name + `\n  namespace: kube-system\nspec:\n  podMetricsEndpoints:\n    - interval: 20s\n      port: metric-port\n      path: /metrics\n      relabelings:\n        - action: replace\n          sourceLabels:\n            - instance\n          regex: (.*)\n          targetLabel: instance\n          replacement: xxxxxx\n  namespaceSelector:\n    matchNames:\n      - test\n  selector:\n    matchLabels:\n      k8s-app: test"
  }

  raw_jobs {
    name   = var.raw_job_name
    config = "scrape_configs:\n- job_name: ` + raw_job_name + `\n  scrape_interval: 20s\n  honor_labels: true\n"
  }
}`
