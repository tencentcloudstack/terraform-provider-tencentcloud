/*
Use this data source to query detailed information of waf user_domains

Example Usage

```hcl
data "tencentcloud_waf_user_domains" "user_domains" {}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudWafUserDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafUserDomainsRead,
		Schema: map[string]*schema.Schema{
			"users_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Domain infos.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"appid": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "User appid.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name.",
						},
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain unique id.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance unique id.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name.",
						},
						"edition": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance type, sparta-waf represents SAAS WAF, clb-waf represents CLB WAF.",
						},
						"level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance level infoNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"write_config": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Switch for accessing log fieldsNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"cls": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CLS switch 1: write, 0: do not writeNote: This field may return null, indicating that a valid value cannot be obtained.",
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

func dataSourceTencentCloudWafUserDomainsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_waf_user_domains.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId     = getLogId(contextNil)
		ctx       = context.WithValue(context.TODO(), logIdKey, logId)
		service   = WafService{client: meta.(*TencentCloudClient).apiV3Conn}
		usersInfo []*waf.UserDomainInfo
	)

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafUserDomainsByFilter(ctx)
		if e != nil {
			return retryError(e)
		}

		usersInfo = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(usersInfo))
	tmpList := make([]map[string]interface{}, 0, len(usersInfo))

	if usersInfo != nil {
		for _, userDomainInfo := range usersInfo {
			userDomainInfoMap := map[string]interface{}{}

			if userDomainInfo.Appid != nil {
				userDomainInfoMap["appid"] = userDomainInfo.Appid
			}

			if userDomainInfo.Domain != nil {
				userDomainInfoMap["domain"] = userDomainInfo.Domain
			}

			if userDomainInfo.DomainId != nil {
				userDomainInfoMap["domain_id"] = userDomainInfo.DomainId
			}

			if userDomainInfo.InstanceId != nil {
				userDomainInfoMap["instance_id"] = userDomainInfo.InstanceId
			}

			if userDomainInfo.InstanceName != nil {
				userDomainInfoMap["instance_name"] = userDomainInfo.InstanceName
			}

			if userDomainInfo.Edition != nil {
				userDomainInfoMap["edition"] = userDomainInfo.Edition
			}

			if userDomainInfo.Level != nil {
				userDomainInfoMap["level"] = userDomainInfo.Level
			}

			if userDomainInfo.WriteConfig != nil {
				userDomainInfoMap["write_config"] = userDomainInfo.WriteConfig
			}

			if userDomainInfo.Cls != nil {
				userDomainInfoMap["cls"] = userDomainInfo.Cls
			}

			ids = append(ids, *userDomainInfo.DomainId)
			tmpList = append(tmpList, userDomainInfoMap)
		}

		_ = d.Set("users_info", tmpList)
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
