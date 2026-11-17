package tcmg

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMonitorGrafanaVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMonitorGrafanaVersionsRead,
		Schema: map[string]*schema.Schema{
			"versions": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Grafana available version list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alias": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Grafana version alias.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Grafana version number.",
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

func dataSourceTencentCloudMonitorGrafanaVersionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_grafana_versions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	var grafanaVersions []*monitor.GrafanaVersion
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMonitorGrafanaVersionsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result == nil {
			log.Printf("[DATASOURCE] read empty, skip SetId")
			return resource.NonRetryableError(fmt.Errorf("DescribeMonitorGrafanaVersionsByFilter return empty"))
		}
		if len(result) == 0 {
			log.Printf("[DATASOURCE] read empty, skip SetId")
			return resource.NonRetryableError(fmt.Errorf("DescribeMonitorGrafanaVersionsByFilter return empty"))
		}
		grafanaVersions = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(grafanaVersions))
	tmpList := make([]map[string]interface{}, 0, len(grafanaVersions))

	if grafanaVersions != nil {
		for _, grafanaVersion := range grafanaVersions {
			grafanaVersionMap := map[string]interface{}{}

			if grafanaVersion.Alias != nil {
				grafanaVersionMap["alias"] = grafanaVersion.Alias
			}

			if grafanaVersion.Version != nil {
				grafanaVersionMap["version"] = grafanaVersion.Version
			}

			if grafanaVersion.Version != nil {
				ids = append(ids, *grafanaVersion.Version)
			}
			tmpList = append(tmpList, grafanaVersionMap)
		}

		_ = d.Set("versions", tmpList)
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
