package cfw

import (
	"context"
	"strconv"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"
)

func DataSourceTencentCloudCfwEdgeFwSwitches() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCfwEdgeFwSwitchesRead,
		Schema: map[string]*schema.Schema{
			"data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Ip switch list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "public ip.",
						},
						"public_ip_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Public IP type.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Id.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Name.",
						},
						"intranet_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Intranet Ip.",
						},
						"asset_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Asset Type.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "region.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "status.",
						},
						"switch_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "switch mode.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "vpc id.",
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

func dataSourceTencentCloudCfwEdgeFwSwitchesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cfw_edge_fw_switches.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		data    []*cfw.EdgeIpInfo
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCfwEdgeFwSwitchesByFilter(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}

		data = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(data))

	if data != nil {
		for _, edgeFwSwitchData := range data {
			edgeFwSwitchDataMap := map[string]interface{}{}

			if edgeFwSwitchData.PublicIp != nil {
				edgeFwSwitchDataMap["public_ip"] = edgeFwSwitchData.PublicIp
			}

			if edgeFwSwitchData.PublicIpType != nil {
				edgeFwSwitchDataMap["public_ip_type"] = edgeFwSwitchData.PublicIpType
			}

			if edgeFwSwitchData.InstanceId != nil {
				edgeFwSwitchDataMap["instance_id"] = edgeFwSwitchData.InstanceId
			}

			if edgeFwSwitchData.InstanceName != nil {
				edgeFwSwitchDataMap["instance_name"] = edgeFwSwitchData.InstanceName
			}

			if edgeFwSwitchData.IntranetIp != nil {
				edgeFwSwitchDataMap["intranet_ip"] = edgeFwSwitchData.IntranetIp
			}

			if edgeFwSwitchData.AssetType != nil {
				edgeFwSwitchDataMap["asset_type"] = edgeFwSwitchData.AssetType
			}

			if edgeFwSwitchData.Region != nil {
				edgeFwSwitchDataMap["region"] = edgeFwSwitchData.Region
			}

			if edgeFwSwitchData.Status != nil {
				edgeFwSwitchDataMap["status"] = edgeFwSwitchData.Status
			}

			if edgeFwSwitchData.SwitchMode != nil {
				edgeFwSwitchDataMap["switch_mode"] = edgeFwSwitchData.SwitchMode
			}

			if edgeFwSwitchData.VpcId != nil {
				edgeFwSwitchDataMap["vpc_id"] = edgeFwSwitchData.VpcId
			}

			tmpList = append(tmpList, edgeFwSwitchDataMap)
		}

		_ = d.Set("data", tmpList)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
