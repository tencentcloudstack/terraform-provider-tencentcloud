/*
Use this data source to query dayu CC http policies

Example Usage

```hcl
data "tencentcloud_dayu_cc_http_policies" "id_test" {
  resource_type = tencentcloud_dayu_cc_http_policy.test_policy.resource_type
  resource_id   = tencentcloud_dayu_cc_http_policy.test_policy.resource_id
  policy_id     = tencentcloud_dayu_cc_http_policy.test_policy.policy_id
}
data "tencentcloud_dayu_cc_http_policies" "name_test" {
  resource_type = tencentcloud_dayu_cc_http_policy.test_policy.resource_type
  resource_id   = tencentcloud_dayu_cc_http_policy.test_policy.resource_id
  name          = tencentcloud_dayu_cc_http_policy.test_policy.name
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDayuCCHttpPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDayuCCHttpPoliciesRead,
		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(DAYU_RESOURCE_TYPE),
				Description:  "Type of the resource that the CC http policy works for, valid values are `bgpip`, `bgp`, `bgp-multip` and `net`.",
			},
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the resource that the CC http policy works for.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the CC http policy to be queried.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 20),
				Description:  "Name of the CC http policy to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of CC http policies. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the resource that the CC self-define http policy works for.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the resource that the CC self-define http policy works for.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the CC self-define http policy.",
						},
						"smode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Match mode.",
						},
						"frequency": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Max frequency per minute.",
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Action mode.",
						},
						"switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicate the CC self-define http policy takes effect or not.",
						},
						"rule_list": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"skey": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Key of the rule.",
									},
									"operator": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Operator of the rule.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule value.",
									},
								},
							},
							Description: "Rule list of the CC self-define http policy.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the CC self-define http policy.",
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the CC self-define http policy.",
						},
						"ip_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "IP of the CC self-define http policy.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDayuCCHttpPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dayu_cc_http_policies.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DayuService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	resourceType := d.Get("resource_type").(string)
	resourceId := d.Get("resource_id").(string)
	policyId := d.Get("policy_id").(string)
	name := d.Get("name").(string)

	ccPolicies := make([]*dayu.CCPolicy, 0)
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, _, err := service.DescribeCCSelfdefinePolicies(ctx, resourceType, resourceId, name, policyId)
		if err != nil {
			return retryError(err)
		}
		ccPolicies = result
		return nil
	})

	if err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(ccPolicies))
	ids := make([]string, 0, len(ccPolicies))

	listItem := make(map[string]interface{})
	for _, policy := range ccPolicies {
		listItem["name"] = *policy.Name
		listItem["create_time"] = *policy.CreateTime
		listItem["policy_id"] = *policy.SetId
		listItem["smode"] = *policy.Smode
		listItem["switch"] = *policy.Switch > 0
		listItem["ip_list"] = helper.StringsInterfaces(policy.IpList)
		if *policy.Smode == "matching" {
			listItem["action"] = *policy.ExeMode
			listItem["rule_list"] = flattenCCRuleList(policy.RuleList)
		} else {
			listItem["frequency"] = int(*policy.Frequency)
		}

		list = append(list, listItem)
		ids = append(ids, *policy.SetId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return writeToFile(output.(string), list)
	}
	return nil

}
