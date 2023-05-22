/*
Provides a resource to create a vpc snapshot_policy

Example Usage

```hcl
resource "tencentcloud_vpc_snapshot_policy" "snapshot_policy" {
  snapshot_policy_name = "terraform-test"
  backup_type          = "time"
  cos_bucket           = "cos-lock-1308919341"
  cos_region           = "ap-guangzhou"
  create_new_cos       = false
  keep_time            = 2

  backup_policies {
    backup_day  = "monday"
    backup_time = "00:00:00"
  }
  backup_policies {
    backup_day  = "tuesday"
    backup_time = "02:03:03"
  }
  backup_policies {
    backup_day  = "wednesday"
    backup_time = "04:13:23"
  }
}
```

Import

vpc snapshot_policy can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_snapshot_policy.snapshot_policy snapshot_policy_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudVpcSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcSnapshotPolicyCreate,
		Read:   resourceTencentCloudVpcSnapshotPolicyRead,
		Update: resourceTencentCloudVpcSnapshotPolicyUpdate,
		Delete: resourceTencentCloudVpcSnapshotPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"snapshot_policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Snapshot policy name.",
			},
			"backup_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Backup strategy type, `operate`: operate backup, `time`: schedule backup.",
			},
			"keep_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The retention time supports 1 to 365 days.",
			},
			"create_new_cos": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to create a new cos bucket, the default is False.Note: This field may return null, indicating that no valid value can be obtained.",
			},
			"cos_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The region where the cos bucket is located.",
			},
			"cos_bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "cos bucket.",
			},
			"backup_policies": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Time backup strategy. Note: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_day": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Backup cycle time, the value can be monday, tuesday, wednesday, thursday, friday, saturday, sunday.",
						},
						"backup_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Backup time point, format:HH:mm:ss.",
						},
					},
				},
			},
			"enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enabled state, True-enabled, False-disabled, the default is True.",
			},
			"snapshot_policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Snapshot policy Id.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.Note: This field may return null, indicating that no valid value can be obtained.",
			},
		},
	}
}

func resourceTencentCloudVpcSnapshotPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_snapshot_policy.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request          = vpc.NewCreateSnapshotPoliciesRequest()
		response         = vpc.NewCreateSnapshotPoliciesResponse()
		snapshotPolicyId string
	)
	snapshotPolicy := vpc.SnapshotPolicy{}

	if v, ok := d.GetOk("snapshot_policy_name"); ok {
		snapshotPolicy.SnapshotPolicyName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("backup_type"); ok {
		snapshotPolicy.BackupType = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("keep_time"); ok {
		snapshotPolicy.KeepTime = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOkExists("create_new_cos"); ok {
		snapshotPolicy.CreateNewCos = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOk("cos_region"); ok {
		snapshotPolicy.CosRegion = helper.String(v.(string))
	}
	if v, ok := d.GetOk("cos_bucket"); ok {
		snapshotPolicy.CosBucket = helper.String(v.(string))
	}
	if v, ok := d.GetOk("snapshot_policy_id"); ok {
		snapshotPolicy.SnapshotPolicyId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("backup_policies"); ok {
		for _, item := range v.([]interface{}) {
			backupPoliciesMap := item.(map[string]interface{})
			backupPolicy := vpc.BackupPolicy{}
			if v, ok := backupPoliciesMap["backup_day"]; ok {
				backupPolicy.BackupDay = helper.String(v.(string))
			}
			if v, ok := backupPoliciesMap["backup_time"]; ok {
				backupPolicy.BackupTime = helper.String(v.(string))
			}
			snapshotPolicy.BackupPolicies = append(snapshotPolicy.BackupPolicies, &backupPolicy)
		}
	}
	if v, ok := d.GetOkExists("enable"); ok {
		snapshotPolicy.Enable = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("create_time"); ok {
		snapshotPolicy.CreateTime = helper.String(v.(string))
	}
	request.SnapshotPolicies = append(request.SnapshotPolicies, &snapshotPolicy)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().CreateSnapshotPolicies(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc snapshotPolicies failed, reason:%+v", logId, err)
		return err
	}

	snapshotPolicies := response.Response.SnapshotPolicies
	if len(snapshotPolicies) < 1 {
		return fmt.Errorf("create vpc snapshotPolicies failed.")
	}

	snapshotPolicyId = *snapshotPolicies[0].SnapshotPolicyId

	d.SetId(snapshotPolicyId)

	return resourceTencentCloudVpcSnapshotPolicyRead(d, meta)
}

func resourceTencentCloudVpcSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_snapshot_policy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	snapshotPolicyId := d.Id()

	snapshotPolicies, err := service.DescribeVpcSnapshotPoliciesById(ctx, snapshotPolicyId)
	if err != nil {
		return err
	}

	if snapshotPolicies == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcSnapshotPolicies` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	for _, snapshotPolicy := range snapshotPolicies {
		if snapshotPolicy.SnapshotPolicyName != nil {
			_ = d.Set("snapshot_policy_name", snapshotPolicy.SnapshotPolicyName)
		}

		if snapshotPolicy.BackupType != nil {
			_ = d.Set("backup_type", snapshotPolicy.BackupType)
		}

		if snapshotPolicy.KeepTime != nil {
			_ = d.Set("keep_time", snapshotPolicy.KeepTime)
		}

		if snapshotPolicy.CreateNewCos != nil {
			_ = d.Set("create_new_cos", snapshotPolicy.CreateNewCos)
		}

		if snapshotPolicy.CosRegion != nil {
			_ = d.Set("cos_region", snapshotPolicy.CosRegion)
		}

		if snapshotPolicy.CosBucket != nil {
			_ = d.Set("cos_bucket", snapshotPolicy.CosBucket)
		}

		if snapshotPolicy.SnapshotPolicyId != nil {
			_ = d.Set("snapshot_policy_id", snapshotPolicy.SnapshotPolicyId)
		}

		if snapshotPolicy.BackupPolicies != nil {
			backupPoliciesList := []interface{}{}
			for _, backupPolicies := range snapshotPolicy.BackupPolicies {
				backupPoliciesMap := map[string]interface{}{}

				if backupPolicies.BackupDay != nil {
					backupPoliciesMap["backup_day"] = backupPolicies.BackupDay
				}

				if backupPolicies.BackupTime != nil {
					backupPoliciesMap["backup_time"] = backupPolicies.BackupTime
				}

				backupPoliciesList = append(backupPoliciesList, backupPoliciesMap)
			}

			_ = d.Set("backup_policies", backupPoliciesList)
		}

		if snapshotPolicy.Enable != nil {
			_ = d.Set("enable", snapshotPolicy.Enable)
		}

		if snapshotPolicy.CreateTime != nil {
			_ = d.Set("create_time", snapshotPolicy.CreateTime)
		}
	}

	return nil
}

func resourceTencentCloudVpcSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_snapshot_policies.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := vpc.NewModifySnapshotPoliciesRequest()

	needChange := false

	snapshotPolicyId := d.Id()

	immutableArgs := []string{
		"backup_type", "create_new_cos", "cos_region", "cos_bucket",
		"snapshot_policy_id", "enable", "create_time",
	}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	mutableArgs := []string{
		"snapshot_policy_name", "keep_time", "backup_policies",
	}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		snapshotPolicy := vpc.BatchModifySnapshotPolicy{}

		snapshotPolicy.SnapshotPolicyId = &snapshotPolicyId

		if v, ok := d.GetOk("snapshot_policy_name"); ok {
			snapshotPolicy.SnapshotPolicyName = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("keep_time"); ok {
			snapshotPolicy.KeepTime = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("backup_policies"); ok {
			for _, item := range v.([]interface{}) {
				backupPoliciesMap := item.(map[string]interface{})
				backupPolicy := vpc.BackupPolicy{}
				if v, ok := backupPoliciesMap["backup_day"]; ok {
					backupPolicy.BackupDay = helper.String(v.(string))
				}
				if v, ok := backupPoliciesMap["backup_time"]; ok {
					backupPolicy.BackupTime = helper.String(v.(string))
				}
				snapshotPolicy.BackupPolicies = append(snapshotPolicy.BackupPolicies, &backupPolicy)
			}
		}
		request.SnapshotPoliciesInfo = append(request.SnapshotPoliciesInfo, &snapshotPolicy)

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ModifySnapshotPolicies(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update vpc snapshotPolicies failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudVpcSnapshotPolicyRead(d, meta)
}

func resourceTencentCloudVpcSnapshotPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_snapshot_policy.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	snapshotPolicyId := d.Id()

	if err := service.DeleteVpcSnapshotPoliciesById(ctx, snapshotPolicyId); err != nil {
		return err
	}

	return nil
}
