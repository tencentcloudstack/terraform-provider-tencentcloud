/*
Use this data source to query detailed information of rum whitelist

Example Usage

```hcl
data "tencentcloud_rum_whitelist" "whitelist" {
  instance_id = "rum-pasZKEI3RLgakj"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudRumWhitelist() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRumWhitelistRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance ID, such as taw-123.",
			},

			"whitelist_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "While list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remarks.",
						},
						"whitelist_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "uin: business identifier.",
						},
						"aid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Business identifier.",
						},
						"ttl": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End time.",
						},
						"wid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Auto-Increment allowlist ID.",
						},
						"create_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creator ID.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
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

func dataSourceTencentCloudRumWhitelistRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_rum_whitelist.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	rumService := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	var whitelistSet []*rum.Whitelist
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := rumService.DescribeRumWhitelistByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		whitelistSet = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Rum whitelistSet failed, reason:%+v", logId, err)
		return err
	}

	whitelistSetList := []interface{}{}
	ids := make([]string, 0, len(whitelistSet))
	if whitelistSet != nil {
		for _, whitelistSet := range whitelistSet {
			ids = append(ids, *whitelistSet.ID)

			whitelistSetMap := map[string]interface{}{}
			if whitelistSet.Remark != nil {
				whitelistSetMap["remark"] = whitelistSet.Remark
			}
			if whitelistSet.WhitelistUin != nil {
				whitelistSetMap["whitelist_uin"] = whitelistSet.WhitelistUin
			}
			if whitelistSet.Aid != nil {
				whitelistSetMap["aid"] = whitelistSet.Aid
			}
			if whitelistSet.Ttl != nil {
				whitelistSetMap["ttl"] = whitelistSet.Ttl
			}
			if whitelistSet.ID != nil {
				whitelistSetMap["wid"] = whitelistSet.ID
			}
			if whitelistSet.CreateUser != nil {
				whitelistSetMap["create_user"] = whitelistSet.CreateUser
			}
			if whitelistSet.CreateTime != nil {
				whitelistSetMap["create_time"] = whitelistSet.CreateTime
			}

			whitelistSetList = append(whitelistSetList, whitelistSetMap)
		}
		_ = d.Set("whitelist_set", whitelistSetList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), whitelistSetList); e != nil {
			return e
		}
	}

	return nil
}
