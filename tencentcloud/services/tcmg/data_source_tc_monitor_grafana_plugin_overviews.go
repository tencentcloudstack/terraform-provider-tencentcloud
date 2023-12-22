package tcmg

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMonitorGrafanaPluginOverviews() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMonitorGrafanaPluginOverviewsRead,
		Schema: map[string]*schema.Schema{
			"plugin_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Plugin set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plugin_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Grafana plugin ID.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Grafana plugin version.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMonitorGrafanaPluginOverviewsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_grafana_plugin_overviews.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	var pluginSet []*monitor.GrafanaPlugin
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMonitorGrafanaPluginOverviewsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		pluginSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(pluginSet))
	tmpList := make([]map[string]interface{}, 0, len(pluginSet))

	if pluginSet != nil {
		for _, grafanaPlugin := range pluginSet {
			grafanaPluginMap := map[string]interface{}{}

			if grafanaPlugin.PluginId != nil {
				grafanaPluginMap["plugin_id"] = grafanaPlugin.PluginId
			}

			if grafanaPlugin.Version != nil {
				grafanaPluginMap["version"] = grafanaPlugin.Version
			}

			ids = append(ids, *grafanaPlugin.PluginId)
			tmpList = append(tmpList, grafanaPluginMap)
		}

		_ = d.Set("plugin_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
