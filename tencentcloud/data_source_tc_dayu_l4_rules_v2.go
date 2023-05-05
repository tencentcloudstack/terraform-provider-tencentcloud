/*
Use this data source to query dayu new layer 4 rules

Example Usage

```hcl
data "tencentcloud_dayu_l4_rules_v2" "tencentcloud_dayu_l4_rules_v2" {
    business = "bgpip"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDayuL4RulesV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDayuL4RulesReadV2,
		Schema: map[string]*schema.Schema{
			"business": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(DAYU_RESOURCE_TYPE),
				Description:  "Type of the resource that the layer 4 rule works for, valid values are `bgpip`, `bgp`, `bgp-multip` and `net`.",
			},
			"virtual_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Virtual port of resource.",
			},
			"ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Ip of the resource.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of layer 4 rules. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol of the rule.",
						},
						"source_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The source port of the layer 4 rule.",
						},
						"virtual_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The virtual port of the layer 4 rule.",
						},
						"keeptime": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The keeptime of the layer 4 rule.",
						},
						"source_list": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source IP or domain.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Weight of the source.",
									},
								},
							},
							Description: "Source list of the rule.",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the 4 layer rule.",
						},
						"lb_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "LB type of the rule, `1` for weight cycling and `2` for IP hash.",
						},
						"keep_enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "session hold switch.",
						},
						"source_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Source type, `1` for source of host, `2` for source of IP.",
						},
						"rule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the rule.",
						},
						"remove_switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Remove the watermark state.",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule modification time.",
						},
						"region": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Corresponding regional information.",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Bind the resource IP information.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Bind the resource ID information.",
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

func dataSourceTencentCloudDayuL4RulesReadV2(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dayu_l4_rules_v2.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DayuService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	business := d.Get("business").(string)
	extendParams := make(map[string]interface{})
	if v, ok := d.GetOk("ip"); ok {
		extendParams["ip"] = v.(string)
	}
	if v, ok := d.GetOk("virtual_port"); ok {
		extendParams["virtual_port"] = v.(int)
	}

	rules := make([]*dayu.NewL4RuleEntry, 0)
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, err := service.DescribeNewL4Rules(ctx, business, extendParams)

		if err != nil {
			return retryError(err)
		}
		rules = result
		return nil
	})

	if err != nil {
		return err
	}
	list := make([]map[string]interface{}, 0)
	for _, rule := range rules {
		tmpRule := make(map[string]interface{})
		tmpRule["protocol"] = *rule.Protocol
		tmpRule["source_port"] = *rule.SourcePort
		tmpRule["virtual_port"] = *rule.VirtualPort
		tmpRule["keeptime"] = *rule.KeepEnable
		tmpSourceList := make([]map[string]interface{}, 0)
		for _, source := range rule.SourceList {
			tmpSource := make(map[string]interface{})
			tmpSource["source"] = *source.Source
			tmpSource["weight"] = *source.Weight
			tmpSourceList = append(tmpSourceList, tmpSource)
		}
		tmpRule["source_list"] = tmpSourceList
		tmpRule["rule_id"] = *rule.RuleId
		tmpRule["lb_type"] = *rule.LbType
		tmpRule["keep_enable"] = *rule.KeepEnable == 1
		tmpRule["source_type"] = *rule.SourceType
		tmpRule["rule_name"] = *rule.RuleName
		tmpRule["remove_switch"] = *rule.RemoveSwitch == 1
		tmpRule["modify_time"] = *rule.ModifyTime
		tmpRule["region"] = *rule.Region
		tmpRule["ip"] = *rule.Ip
		tmpRule["id"] = *rule.Id
		list = append(list, tmpRule)
	}
	ids := make([]string, 0, len(list))
	for _, listItem := range list {
		ids = append(ids, listItem["rule_id"].(string))
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("list", list); e != nil {
		log.Printf("[CRITAL]%s provider set rules fail, reason:%s\n", logId, e.Error())
		return e
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return writeToFile(output.(string), list)
	}
	return nil

}
