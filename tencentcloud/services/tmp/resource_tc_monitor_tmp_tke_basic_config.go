package tmp

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
)

const (
	SERVICE_MONITORS string = "service_monitors"
	POD_MONITORS     string = "pod_monitors"
	RAW_JOBS         string = "raw_jobs"
)

func ResourceTencentCloudMonitorTmpTkeBasicConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorTmpTkeBasicConfigCreate,
		Read:   resourceTencentCloudMonitorTmpTkeBasicConfigRead,
		Update: resourceTencentCloudMonitorTmpTkeBasicConfigUpdate,
		Delete: resourceTencentCloudMonitorTmpTkeBasicConfigDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of instance.",
			},

			"cluster_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Type of cluster.",
			},

			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of cluster.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name. The naming rule is: namespace/name. If you don&#39;t have any namespace, use the default namespace: kube-system, otherwise use the specified one.",
			},

			"metrics_name": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: tccommon.ValidateNotEmpty,
				},
				Description: "Configure the name of the metric to keep on.",
			},

			"config_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "config type, `service_monitors`, `pod_monitors`, `raw_jobs`.",
			},

			"config": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Full configuration in yaml format.",
			},
		},
	}
}

func resourceTencentCloudMonitorTmpTkeBasicConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_tke_basic_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		instanceId  string
		clusterType string
		clusterId   string
		name        string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}
	if v, ok := d.GetOk("cluster_type"); ok {
		clusterType = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	d.SetId(strings.Join([]string{instanceId, clusterType, clusterId, name}, tccommon.FILED_SP))

	return resourceTencentCloudMonitorTmpTkeBasicConfigUpdate(d, meta)
}

func resourceTencentCloudMonitorTmpTkeBasicConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_tke_basic_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	clusterType := idSplit[1]
	clusterId := idSplit[2]
	name := idSplit[3]

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	tmpTkeBasicConfig, err := service.DescribeTkeTmpBasicConfigById(ctx, clusterId, clusterType, instanceId)
	if err != nil {
		return err
	}

	if tmpTkeBasicConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorTmpTkeBasicConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("cluster_type", clusterType)
	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("name", name)

	configType, config, err := service.GetConfigType(name, tmpTkeBasicConfig)
	if err != nil {
		return err
	}
	_ = d.Set("config_type", configType)
	_ = d.Set("config", config.Config)

	return nil
}

func resourceTencentCloudMonitorTmpTkeBasicConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_tke_basic_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := monitor.NewModifyPrometheusConfigRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	clusterType := idSplit[1]
	clusterId := idSplit[2]
	name := idSplit[3]

	request.InstanceId = &instanceId
	request.ClusterType = &clusterType
	request.ClusterId = &clusterId

	if v, ok := d.GetOk("metrics_name"); ok {
		regexs := []string{}
		regexSet := v.(*schema.Set).List()
		for i := range regexSet {
			regex := regexSet[i].(string)
			regexs = append(regexs, regex)
		}

		service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		tmpTkeBasicConfig, err := service.DescribeTkeTmpBasicConfigById(ctx, clusterId, clusterType, instanceId)
		if err != nil {
			return err
		}
		configType, config, err := service.GetConfigType(name, tmpTkeBasicConfig)
		if err != nil {
			return err
		}

		serviceMonitors, podMonitors, rawMobs, err := configInit(configType, config, regexs)
		if err != nil {
			return err
		}

		if serviceMonitors != "" {
			prometheusConfig := []*monitor.PrometheusConfigItem{}
			prometheusConfig = append(prometheusConfig, &monitor.PrometheusConfigItem{
				Name:   &name,
				Config: &serviceMonitors,
			})
			request.ServiceMonitors = prometheusConfig
		}
		if podMonitors != "" {
			prometheusConfig := []*monitor.PrometheusConfigItem{}
			prometheusConfig = append(prometheusConfig, &monitor.PrometheusConfigItem{
				Name:   &name,
				Config: &podMonitors,
			})
			request.PodMonitors = prometheusConfig
		}
		if rawMobs != "" {
			prometheusConfig := []*monitor.PrometheusConfigItem{}
			prometheusConfig = append(prometheusConfig, &monitor.PrometheusConfigItem{
				Name:   &name,
				Config: &rawMobs,
			})
			request.RawJobs = prometheusConfig
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().ModifyPrometheusConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update monitor tmpTkeBasicConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMonitorTmpTkeBasicConfigRead(d, meta)
}

func resourceTencentCloudMonitorTmpTkeBasicConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_tke_basic_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func configInit(configType string, respParams *monitor.PrometheusConfigItem, regexs []string) (serviceMonitorConfig, podMonitorConfig, rawMobConfig string, errRet error) {
	config := svcmonitor.PrometheusConfig{
		Config: respParams.Config,
		Regex:  regexs,
	}
	switch configType {
	case SERVICE_MONITORS:
		serviceMonitor, err := config.UnmarshalToMap()
		if err != nil {
			errRet = err
			return
		}
		spec := serviceMonitor["spec"].(map[interface{}]interface{})["endpoints"].([]interface{})
		serviceMonitors, err := config.SetRegex(spec)
		serviceMonitor["spec"].(map[interface{}]interface{})["endpoints"] = serviceMonitors
		if err != nil {
			errRet = err
			return
		}
		serviceMonitorConfig, errRet = config.MarshalToYaml(&serviceMonitor)
		return
	case POD_MONITORS:
		serviceMonitor, err := config.UnmarshalToMap()
		if err != nil {
			errRet = err
			return
		}
		spec := serviceMonitor["spec"].(map[interface{}]interface{})["podMetricsEndpoints"].([]interface{})
		serviceMonitors, err := config.SetRegex(spec)
		serviceMonitor["spec"].(map[interface{}]interface{})["podMetricsEndpoints"] = serviceMonitors
		if err != nil {
			errRet = err
			return
		}
		podMonitorConfig, errRet = config.MarshalToYaml(&serviceMonitor)
		return
	case RAW_JOBS:
		rawMob, err := config.UnmarshalToMap()
		if err != nil {
			errRet = err
			return
		}
		configs := rawMob["scrape_configs"].([]interface{})
		rawMobConfigs, err := config.SetRegex(configs)
		rawMob["scrape_configs"] = rawMobConfigs
		if err != nil {
			errRet = err
			return
		}
		rawMobConfig, errRet = config.MarshalToYaml(&rawMob)
		return
	}
	return
}
