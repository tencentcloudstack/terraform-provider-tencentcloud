/*
Use this data source to query dayu CC https policies

Example Usage

```hcl
data "tencentcloud_dayu_cc_https_policies" "name_test" {
  resource_type = tencentcloud_dayu_cc_https_policy.test_policy.resource_type
  resource_id = tencentcloud_dayu_cc_https_policy.test_policy.resource_id
  name = tencentcloud_dayu_cc_https_policy.test_policy.name
}
data "tencentcloud_dayu_cc_https_policies" "id_test" {
  resource_type = tencentcloud_dayu_cc_https_policy.test_policy.resource_type
  resource_id = tencentcloud_dayu_cc_https_policy.test_policy.resource_id
  policy_id = tencentcloud_dayu_cc_https_policy.test_policy.policy_id
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
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDayuCCHttpsPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDayuCCHttpsPoliciesRead,
		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(DAYU_RESOURCE_TYPE),
				Description:  "Type of the resource that the CC https policy works for, valid value is `bgpip`.",
			},
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the resource that the CC https policy works for.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the CC https policy to be queried.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 20),
				Description:  "Name of the CC https policy to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of CC https policies. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the resource that the CC self-define https policy works for.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the resource that the CC self-define https policy works for.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the CC self-define https policy.",
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Action mode.",
						},
						"switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicate the CC self-define https policy takes effect or not.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain that the CC self-define https policy works for.",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule id of the domain that the CC self-define https policy works for.",
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
							Description: "Rule list of the CC self-define https policy.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the CC self-define https policy.",
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the CC self-define https policy.",
						},
						"ip_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Ip of the CC self-define https policy.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDayuCCHttpsPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dayu_cc_https_policies.read")()

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
		listItem["switch"] = *policy.Switch > 0
		listItem["ip_list"] = helper.StringsInterfaces(policy.IpList)
		listItem["action"] = *policy.ExeMode
		listItem["rule_list"] = flattenCCRuleList(policy.RuleList)
		listItem["rule_id"] = *policy.RuleId
		listItem["domain"] = *policy.Domain

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
