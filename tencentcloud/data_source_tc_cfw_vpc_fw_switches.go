/*
Use this data source to query detailed information of cfw vpc_fw_switches

Example Usage

```hcl
data "tencentcloud_cfw_vpc_fw_switches" "example" {
  vpc_ins_id = "cfwg-c8c2de41"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"
)

func dataSourceTencentCloudCfwVpcFwSwitches() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCfwVpcFwSwitchesRead,
		Schema: map[string]*schema.Schema{
			"vpc_ins_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Firewall instance id.",
			},
			"switch_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Switch list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Firewall switch ID.",
						},
						"switch_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Firewall switch name.",
						},
						"switch_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "switch mode.",
						},
						"enable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Switch status 0: off, 1: on.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Switch status 0: normal, 1: switching.",
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

func dataSourceTencentCloudCfwVpcFwSwitchesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cfw_vpc_fw_switches.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = CfwService{client: meta.(*TencentCloudClient).apiV3Conn}
		switchList []*cfw.FwGroupSwitchShow
		vpcInsId   string
	)

	if v, ok := d.GetOk("vpc_ins_id"); ok {
		vpcInsId = v.(string)
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCfwVpcFwSwitchesByFilter(ctx, vpcInsId)
		if e != nil {
			return retryError(e)
		}

		switchList = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(switchList))

	if switchList != nil {
		for _, fwGroupSwitcheshow := range switchList {
			fwGroupSwitcheshowMap := map[string]interface{}{}

			if fwGroupSwitcheshow.SwitchId != nil {
				fwGroupSwitcheshowMap["switch_id"] = fwGroupSwitcheshow.SwitchId
			}

			if fwGroupSwitcheshow.SwitchName != nil {
				fwGroupSwitcheshowMap["switch_name"] = fwGroupSwitcheshow.SwitchName
			}

			if fwGroupSwitcheshow.SwitchMode != nil {
				fwGroupSwitcheshowMap["switch_mode"] = fwGroupSwitcheshow.SwitchMode
			}

			if fwGroupSwitcheshow.Enable != nil {
				fwGroupSwitcheshowMap["enable"] = fwGroupSwitcheshow.Enable
			}

			if fwGroupSwitcheshow.Status != nil {
				fwGroupSwitcheshowMap["status"] = fwGroupSwitcheshow.Status
			}

			tmpList = append(tmpList, fwGroupSwitcheshowMap)
		}

		_ = d.Set("switch_list", tmpList)
	}

	d.SetId(vpcInsId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
