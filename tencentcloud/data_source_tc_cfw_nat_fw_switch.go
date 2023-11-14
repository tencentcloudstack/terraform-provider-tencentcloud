/*
Use this data source to query detailed information of cfw nat_fw_switch

Example Usage

```hcl
data "tencentcloud_cfw_nat_fw_switch" "nat_fw_switch" {
  search_value = ""
  status =
  vpc_id = ""
  nat_id = ""
  nat_ins_id = ""
  area = ""
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

func dataSourceTencentCloudCfwNatFwSwitch() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCfwNatFwSwitchRead,
		Schema: map[string]*schema.Schema{
			"search_value": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Search Value.",
			},

			"status": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Switch status ，1open，0close.",
			},

			"vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter the VPC to which the NAT firewall subnet switch belongs.",
			},

			"nat_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter the NAT gateway to which the NAT firewall subnet switch belongs.",
			},

			"nat_ins_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter the NAT firewall instance to which the NAT firewall subnet switch belongs.",
			},

			"area": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter the region to which the NAT firewall subnet switch belongs.",
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

func dataSourceTencentCloudCfwNatFwSwitchRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cfw_nat_fw_switch.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("search_value"); ok {
		paramMap["SearchValue"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("status"); v != nil {
		paramMap["Status"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		paramMap["VpcId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("nat_id"); ok {
		paramMap["NatId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("nat_ins_id"); ok {
		paramMap["NatInsId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("area"); ok {
		paramMap["Area"] = helper.String(v.(string))
	}

	service := CfwService{client: meta.(*TencentCloudClient).apiV3Conn}

	var data []*cfw.NatSwitchListData

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCfwNatFwSwitchByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		data = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(data))
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

			ids = append(ids, *natSwitchListData.SubnetId)
			tmpList = append(tmpList, natSwitchListDataMap)
		}

		_ = d.Set("data", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
