package tke

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func DataSourceTencentCloudTkeAddonValues() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTkeAddonValuesRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the TKE cluster, which can be obtained from the cluster list in the TKE console (https://console.cloud.tencent.com/tke2/cluster).",
			},

			"addon_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the addon, which can be obtained from the component management page in the cluster details page of the TKE console.",
			},

			"values": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Parameter list in JSON format. If the addon is installed, the parameters are rendered with the configured parameters and the default parameters of the chart; otherwise the value is empty.",
			},

			"default_values": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Parameter list supported by the addon in JSON format, whose value is the default value of the chart.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTkeAddonValuesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tke_addon_values.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	var clusterId string
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	var addonName string
	if v, ok := d.GetOk("addon_name"); ok {
		addonName = v.(string)
	}

	var resp *tke.DescribeAddonValuesResponseParams
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAddonValuesList(ctx, clusterId, addonName)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil {
			log.Printf("[DATASOURCE] read empty, skip SetId")
			return resource.NonRetryableError(fmt.Errorf("tencentcloud_tke_addon_values DescribeAddonValuesResponse is nil"))
		}

		resp = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	if resp != nil {
		if resp.Values != nil {
			_ = d.Set("values", *resp.Values)
		}

		if resp.DefaultValues != nil {
			_ = d.Set("default_values", *resp.DefaultValues)
		}
	}

	d.SetId(strings.Join([]string{clusterId, addonName}, tccommon.FILED_SP))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
