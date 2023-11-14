/*
Provides a resource to create a cbs snapshot_share_permission

Example Usage

```hcl
resource "tencentcloud_cbs_snapshot_share_permission" "snapshot_share_permission" {
  account_ids =
  permission = "SHARE"
  snapshot_ids =
}
```

Import

cbs snapshot_share_permission can be imported using the id, e.g.

```
terraform import tencentcloud_cbs_snapshot_share_permission.snapshot_share_permission snapshot_share_permission_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCbsSnapshotSharePermission() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCbsSnapshotSharePermissionCreate,
		Read:   resourceTencentCloudCbsSnapshotSharePermissionRead,
		Update: resourceTencentCloudCbsSnapshotSharePermissionUpdate,
		Delete: resourceTencentCloudCbsSnapshotSharePermissionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"account_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of account IDs with which a snapshot is shared. For the format of array-type parameters, see[API Introduction](https://cloud.tencent.com/document/api/213/568). You can find the account ID in[Account Information](https://console.cloud.tencent.com/developer).",
			},

			"permission": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Operations. Valid values SHARE, sharing an image; CANCEL, cancelling the sharing of an image.",
			},

			"snapshot_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The ID of the snapshot to be queried. You can obtain this by using [DescribeSnapshots](https://cloud.tencent.com/document/api/362/15647).",
			},
		},
	}
}

func resourceTencentCloudCbsSnapshotSharePermissionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_snapshot_share_permission.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cbs.NewModifySnapshotsSharePermissionRequest()
		response   = cbs.NewModifySnapshotsSharePermissionResponse()
		snapshotId string
	)
	if v, ok := d.GetOk("account_ids"); ok {
		accountIdsSet := v.(*schema.Set).List()
		for i := range accountIdsSet {
			accountIds := accountIdsSet[i].(string)
			request.AccountIds = append(request.AccountIds, &accountIds)
		}
	}

	if v, ok := d.GetOk("permission"); ok {
		request.Permission = helper.String(v.(string))
	}

	if v, ok := d.GetOk("snapshot_ids"); ok {
		snapshotIdsSet := v.(*schema.Set).List()
		for i := range snapshotIdsSet {
			snapshotIds := snapshotIdsSet[i].(string)
			request.SnapshotIds = append(request.SnapshotIds, &snapshotIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCbsClient().ModifySnapshotsSharePermission(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cbs SnapshotSharePermission failed, reason:%+v", logId, err)
		return err
	}

	snapshotId = *response.Response.SnapshotId
	d.SetId(snapshotId)

	return resourceTencentCloudCbsSnapshotSharePermissionRead(d, meta)
}

func resourceTencentCloudCbsSnapshotSharePermissionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_snapshot_share_permission.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CbsService{client: meta.(*TencentCloudClient).apiV3Conn}

	snapshotSharePermissionId := d.Id()

	SnapshotSharePermission, err := service.DescribeCbsSnapshotSharePermissionById(ctx, snapshotId)
	if err != nil {
		return err
	}

	if SnapshotSharePermission == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CbsSnapshotSharePermission` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if SnapshotSharePermission.AccountIds != nil {
		_ = d.Set("account_ids", SnapshotSharePermission.AccountIds)
	}

	if SnapshotSharePermission.Permission != nil {
		_ = d.Set("permission", SnapshotSharePermission.Permission)
	}

	if SnapshotSharePermission.SnapshotIds != nil {
		_ = d.Set("snapshot_ids", SnapshotSharePermission.SnapshotIds)
	}

	return nil
}

func resourceTencentCloudCbsSnapshotSharePermissionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_snapshot_share_permission.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cbs.NewModifySnapshotsSharePermissionRequest()

	snapshotSharePermissionId := d.Id()

	request.SnapshotId = &snapshotId

	immutableArgs := []string{"account_ids", "permission", "snapshot_ids"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("account_ids") {
		if v, ok := d.GetOk("account_ids"); ok {
			accountIdsSet := v.(*schema.Set).List()
			for i := range accountIdsSet {
				accountIds := accountIdsSet[i].(string)
				request.AccountIds = append(request.AccountIds, &accountIds)
			}
		}
	}

	if d.HasChange("permission") {
		if v, ok := d.GetOk("permission"); ok {
			request.Permission = helper.String(v.(string))
		}
	}

	if d.HasChange("snapshot_ids") {
		if v, ok := d.GetOk("snapshot_ids"); ok {
			snapshotIdsSet := v.(*schema.Set).List()
			for i := range snapshotIdsSet {
				snapshotIds := snapshotIdsSet[i].(string)
				request.SnapshotIds = append(request.SnapshotIds, &snapshotIds)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCbsClient().ModifySnapshotsSharePermission(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cbs SnapshotSharePermission failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCbsSnapshotSharePermissionRead(d, meta)
}

func resourceTencentCloudCbsSnapshotSharePermissionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cbs_snapshot_share_permission.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CbsService{client: meta.(*TencentCloudClient).apiV3Conn}
	snapshotSharePermissionId := d.Id()

	if err := service.DeleteCbsSnapshotSharePermissionById(ctx, snapshotId); err != nil {
		return err
	}

	return nil
}
