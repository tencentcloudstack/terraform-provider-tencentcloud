package cbs

import (
	"context"

	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCbsSnapshotSharePermission() *schema.Resource {
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
			"snapshot_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of the snapshot to be queried. You can obtain this by using [DescribeSnapshots](https://cloud.tencent.com/document/api/362/15647).",
			},
		},
	}
}

func resourceTencentCloudCbsSnapshotSharePermissionCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_snapshot_share_permission.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var snapshotId string
	var accountIdsSet []interface{}
	if v, ok := d.GetOk("account_ids"); ok {
		accountIdsSet = v.(*schema.Set).List()
	}

	if v, ok := d.GetOk("snapshot_id"); ok {
		snapshotId = v.(string)
	}

	service := CbsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err := service.ModifySnapshotsSharePermission(ctx, snapshotId, SNAPSHOT_SHARE_PERMISSION_SHARE, helper.InterfacesStrings(accountIdsSet))
	if err != nil {
		return err
	}
	d.SetId(snapshotId)

	return resourceTencentCloudCbsSnapshotSharePermissionRead(d, meta)
}

func resourceTencentCloudCbsSnapshotSharePermissionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_snapshot_share_permission.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := CbsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	snapshotId := d.Id()
	snapshotSharePermissions := []*cbs.SharePermission{}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCbsSnapshotSharePermissionById(ctx, snapshotId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		snapshotSharePermissions = result
		return nil
	})

	if err != nil {
		return err
	}

	accountIds := make([]string, 0)
	for _, snapshotSharePermission := range snapshotSharePermissions {
		accountIds = append(accountIds, *snapshotSharePermission.AccountId)
	}

	_ = d.Set("account_ids", accountIds)
	_ = d.Set("snapshot_id", snapshotId)
	return nil
}

func resourceTencentCloudCbsSnapshotSharePermissionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_snapshot_share_permission.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	snapshotId := d.Id()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := CbsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	if d.HasChange("account_ids") {
		old, new := d.GetChange("account_ids")
		oldSet := old.(*schema.Set)
		newSet := new.(*schema.Set)
		add := newSet.Difference(oldSet).List()
		remove := oldSet.Difference(newSet).List()
		if len(add) > 0 {
			addError := service.ModifySnapshotsSharePermission(ctx, snapshotId, SNAPSHOT_SHARE_PERMISSION_SHARE, helper.InterfacesStrings(add))
			if addError != nil {
				return addError
			}
		}

		if len(remove) > 0 {
			removeError := service.ModifySnapshotsSharePermission(ctx, snapshotId, SNAPSHOT_SHARE_PERMISSION_CANCEL, helper.InterfacesStrings(remove))
			if removeError != nil {
				return removeError
			}
		}
	}

	return resourceTencentCloudCbsSnapshotSharePermissionRead(d, meta)
}

func resourceTencentCloudCbsSnapshotSharePermissionDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_snapshot_share_permission.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := CbsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	snapshotId := d.Id()
	snapshotSharePermissions := []*cbs.SharePermission{}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCbsSnapshotSharePermissionById(ctx, snapshotId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		snapshotSharePermissions = result
		return nil
	})

	if err != nil {
		return err
	}

	accountIds := make([]string, 0)
	for _, snapshotSharePermission := range snapshotSharePermissions {
		accountIds = append(accountIds, *snapshotSharePermission.AccountId)
	}

	if err := service.ModifySnapshotsSharePermission(ctx, snapshotId, SNAPSHOT_SHARE_PERMISSION_CANCEL, accountIds); err != nil {
		return err
	}

	return nil
}
