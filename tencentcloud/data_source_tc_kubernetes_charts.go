/*
Use this data source to query detailed information of kubernetes cluster addons.

Example Usage

```hcl
data "tencentcloud_kubernetes_charts" "name" {}
```
*/
package tencentcloud

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudKubernetesCharts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKubernetesChartsRead,
		Schema: map[string]*schema.Schema{
			"kind": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Kind of app chart. Available values: `log`, `scheduler`, `network`, `storage`, `monitor`, `dns`, `image`, `other`, `invisible`.",
			},
			"arch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Operation system app supported. Available values: `arm32`, `arm64`, `amd64`.",
			},
			"cluster_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cluster type. Available values: `tke`, `eks`.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"chart_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "App chart list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of chart.",
						},
						"label": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Label of chart.",
						},
						"latest_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Chart latest version.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudKubernetesChartsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_kubernetes_charts.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := TkeService{client: client}

	var (
		kind        = d.Get("kind").(string)
		arch        = d.Get("arch").(string)
		clusterType = d.Get("cluster_type").(string)
	)

	request := tke.NewGetTkeAppChartListRequest()
	if kind != "" {
		request.Kind = &kind
	}

	if arch != "" {
		request.Arch = &arch
	}

	if clusterType != "" {
		request.ClusterType = &clusterType
	}

	response, err := service.GetTkeAppChartList(ctx, request)
	if err != nil {
		return err
	}

	chartList := make([]interface{}, 0)

	for i := range response {
		item := response[i]
		chart := map[string]interface{}{
			"name":           item.Name,
			"latest_version": item.LatestVersion,
		}

		label := make(map[string]interface{})

		if err := json.Unmarshal([]byte(*item.Label), &label); err != nil {
			return err
		}

		chart["label"] = label

		chartList = append(chartList, chart)
	}

	err = d.Set("chart_list", chartList)

	if err != nil {
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		err := writeToFile(output.(string), chartList)
		if err != nil {
			return err
		}
	}

	ids := []string{kind, arch, clusterType}
	d.SetId("app_chart_" + helper.DataResourceIdsHash(ids))

	return nil
}
