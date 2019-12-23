/*
Use this data source to query detailed information of CBS snapshot policies.

Example Usage

```hcl
data "tencentcloud_cbs_snapshot_policies" "policies" {
  snapshot_policy_id = "snap-f3io7adt"
  snapshot_policy_name = "test"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCbsSnapshotPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCbsSnapshotPoliciesRead,

		Schema: map[string]*schema.Schema{
			"snapshot_policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the snapshot policy to be queried.",
			},
			"snapshot_policy_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the snapshot policy to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"snapshot_policy_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of snapshot policy. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snapshot_policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the snapshot policy.",
						},
						"snapshot_policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the snapshot policy.",
						},
						"repeat_weekdays": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Computed:    true,
							Description: "Trigger days of periodic snapshot.",
						},
						"repeat_hours": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Computed:    true,
							Description: "Trigger hours of periodic snapshot.",
						},
						"retention_days": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Retention days of the snapshot.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the snapshot policy.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the snapshot policy.",
						},
						"attached_storage_ids": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Storage ids that the snapshot policy attached.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCbsSnapshotPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cbs_snapshot_policies.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	var policyId string
	var policyName string
	if v, ok := d.GetOk("snapshot_policy_id"); ok {
		policyId = v.(string)
	}
	if v, ok := d.GetOk("snapshot_policy_name"); ok {
		policyName = v.(string)
	}
	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var policies []*cbs.AutoSnapshotPolicy
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		policies, errRet = cbsService.DescribeSnapshotPolicy(ctx, policyId, policyName)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read cbs snapshot policies failed, reason:%s\n ", logId, err.Error())
		return err
	}

	ids := make([]string, 0, len(policies))
	policyList := make([]map[string]interface{}, 0, len(policies))
	for _, policy := range policies {
		mapping := map[string]interface{}{
			"snapshot_policy_id":   policy.AutoSnapshotPolicyId,
			"snapshot_policy_name": policy.AutoSnapshotPolicyName,
			"retention_days":       policy.RetentionDays,
			"status":               policy.AutoSnapshotPolicyState,
			"create_time":          policy.CreateTime,
			"attached_storage_ids": helper.StringsInterfaces(policy.DiskIdSet),
		}
		if len(policy.Policy) < 1 {
			continue
		}
		mapping["repeat_weekdays"] = helper.Uint64sInterfaces(policy.Policy[0].DayOfWeek)
		mapping["repeat_hours"] = helper.Uint64sInterfaces(policy.Policy[0].Hour)
		policyList = append(policyList, mapping)
		ids = append(ids, *policy.AutoSnapshotPolicyId)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if err = d.Set("snapshot_policy_list", policyList); err != nil {
		log.Printf("[CRITAL]%s provider set snapshot policy list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err = writeToFile(output.(string), policyList); err != nil {
			return err
		}
	}
	return nil
}
