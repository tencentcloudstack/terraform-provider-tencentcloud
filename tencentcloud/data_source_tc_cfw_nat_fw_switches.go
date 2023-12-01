/*
Use this data source to query detailed information of cfw nat_fw_switches

Example Usage

Query Nat instance'switch by instance id

```hcl
data "tencentcloud_cfw_nat_fw_switches" "example" {
  nat_ins_id = "cfwnat-18d2ba18"
}
```

Or filter by switch status

```hcl
data "tencentcloud_cfw_nat_fw_switches" "example" {
  nat_ins_id = "cfwnat-18d2ba18"
  status     = 1
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCfwNatFwSwitches() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCfwNatFwSwitchesRead,
		Schema: map[string]*schema.Schema{
			"nat_ins_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter the NAT firewall instance to which the NAT firewall subnet switch belongs.",
			},
			"status": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Switch status, 1 open; 0 close.",
			},
			"data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "NAT border firewall switch list data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet Id.",
						},
						"subnet_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet Name.",
						},
						"subnet_cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IPv4 CIDR.",
						},
						"route_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Route Id.",
						},
						"route_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Route Name.",
						},
						"cvm_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cvm Num.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Vpc Id.",
						},
						"vpc_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Vpc Name.",
						},
						"enable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Effective status.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Switch status.",
						},
						"nat_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "NAT gatway Id.",
						},
						"nat_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "NAT gatway name.",
						},
						"nat_ins_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "NAT firewall instance Id.",
						},
						"nat_ins_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "NAT firewall instance name.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"abnormal": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the switch is abnormal, 0: normal, 1: abnormal.",
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

func dataSourceTencentCloudCfwNatFwSwitchesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cfw_nat_fw_switches.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		ctx      = context.WithValue(context.TODO(), logIdKey, logId)
		service  = CfwService{client: meta.(*TencentCloudClient).apiV3Conn}
		data     []*cfw.NatSwitchListData
		natInsId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("nat_ins_id"); ok {
		paramMap["NatInsId"] = helper.String(v.(string))
		natInsId = v.(string)
	}

	if v, _ := d.GetOkExists("status"); v != nil {
		paramMap["Status"] = helper.IntInt64(v.(int))
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCfwNatFwSwitchesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		data = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(data))

	if data != nil {
		for _, natSwitchListData := range data {
			natSwitchListDataMap := map[string]interface{}{}

			if natSwitchListData.Id != nil {
				natSwitchListDataMap["id"] = natSwitchListData.Id
			}

			if natSwitchListData.SubnetId != nil {
				natSwitchListDataMap["subnet_id"] = natSwitchListData.SubnetId
			}

			if natSwitchListData.SubnetName != nil {
				natSwitchListDataMap["subnet_name"] = natSwitchListData.SubnetName
			}

			if natSwitchListData.SubnetCidr != nil {
				natSwitchListDataMap["subnet_cidr"] = natSwitchListData.SubnetCidr
			}

			if natSwitchListData.RouteId != nil {
				natSwitchListDataMap["route_id"] = natSwitchListData.RouteId
			}

			if natSwitchListData.RouteName != nil {
				natSwitchListDataMap["route_name"] = natSwitchListData.RouteName
			}

			if natSwitchListData.CvmNum != nil {
				natSwitchListDataMap["cvm_num"] = natSwitchListData.CvmNum
			}

			if natSwitchListData.VpcId != nil {
				natSwitchListDataMap["vpc_id"] = natSwitchListData.VpcId
			}

			if natSwitchListData.VpcName != nil {
				natSwitchListDataMap["vpc_name"] = natSwitchListData.VpcName
			}

			if natSwitchListData.Enable != nil {
				natSwitchListDataMap["enable"] = natSwitchListData.Enable
			}

			if natSwitchListData.Status != nil {
				natSwitchListDataMap["status"] = natSwitchListData.Status
			}

			if natSwitchListData.NatId != nil {
				natSwitchListDataMap["nat_id"] = natSwitchListData.NatId
			}

			if natSwitchListData.NatName != nil {
				natSwitchListDataMap["nat_name"] = natSwitchListData.NatName
			}

			if natSwitchListData.NatInsId != nil {
				natSwitchListDataMap["nat_ins_id"] = natSwitchListData.NatInsId
			}

			if natSwitchListData.NatInsName != nil {
				natSwitchListDataMap["nat_ins_name"] = natSwitchListData.NatInsName
			}

			if natSwitchListData.Region != nil {
				natSwitchListDataMap["region"] = natSwitchListData.Region
			}

			if natSwitchListData.Abnormal != nil {
				natSwitchListDataMap["abnormal"] = natSwitchListData.Abnormal
			}

			tmpList = append(tmpList, natSwitchListDataMap)
		}

		_ = d.Set("data", tmpList)
	}

	d.SetId(natInsId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
