package gaap

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudGaapProxiesStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapProxiesStatusRead,
		Schema: map[string]*schema.Schema{
			"proxy_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of Proxy IDs.",
			},

			"instance_status_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Proxy status list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "proxy instance ID.",
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
							Description: "proxy status.Among them:\n" +
								"- RUNNING indicates running;\n" +
								"- CREATING indicates being created;\n" +
								"- DESTROYING indicates being destroyed;\n" +
								"- OPENING indicates being opened;\n" +
								"- CLOSING indicates being closed;\n" +
								"- Closed indicates that it has been closed;\n" +
								"- ADJUSTING represents a configuration change in progress;\n" +
								"- ISOLATING indicates being isolated;\n" +
								"- ISOLATED indicates that it has been isolated;\n" +
								"- MOVING indicates that migration is in progress.",
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

func dataSourceTencentCloudGaapProxiesStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_gaap_proxies_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		paramMap["InstanceIds"] = helper.InterfacesStringsPoint(instanceIdsSet)
	}

	if v, ok := d.GetOk("proxy_ids"); ok {
		proxyIdsSet := v.(*schema.Set).List()
		paramMap["ProxyIds"] = helper.InterfacesStringsPoint(proxyIdsSet)
	}

	service := GaapService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var instanceStatusSet []*gaap.ProxyStatus

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeGaapProxiesStatusByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instanceStatusSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceStatusSet))
	tmpList := make([]map[string]interface{}, 0, len(instanceStatusSet))

	if instanceStatusSet != nil {
		for _, proxyStatus := range instanceStatusSet {
			proxyStatusMap := map[string]interface{}{}

			if proxyStatus.InstanceId != nil {
				proxyStatusMap["instance_id"] = proxyStatus.InstanceId
			}

			if proxyStatus.Status != nil {
				proxyStatusMap["status"] = proxyStatus.Status
			}

			ids = append(ids, *proxyStatus.InstanceId)
			tmpList = append(tmpList, proxyStatusMap)
		}

		_ = d.Set("instance_status_set", tmpList)
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
