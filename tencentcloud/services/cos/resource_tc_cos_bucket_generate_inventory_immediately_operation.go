package cos

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func ResourceTencentCloudCosBucketGenerateInventoryImmediatelyOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCosBucketGenerateInventoryImmediatelyOperationCreate,
		Read:   resourceTencentCloudCosBucketGenerateInventoryImmediatelyOperationRead,
		Delete: resourceTencentCloudCosBucketGenerateInventoryImmediatelyOperationDelete,

		Schema: map[string]*schema.Schema{
			"inventory_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The id of inventory.",
			},
			"bucket": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Bucket.",
			},
		},
	}
}

func resourceTencentCloudCosBucketGenerateInventoryImmediatelyOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_generate_inventory_immediately_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	inventoryId := d.Get("inventory_id").(string)
	bucket := d.Get("bucket").(string)
	result, _, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTencentCosClient(bucket).Bucket.GetInventory(ctx, inventoryId)
	if err != nil {
		log.Printf("[CRITAL]%s get cos bucketInventory failed, reason:%+v", logId, err)
		return err
	}
	id := fmt.Sprintf("%s_instant_%s", result.ID, time.Now().Format("20060102150405"))
	inventoryOpt := &cos.BucketPostInventoryOptions{
		ID:                     id,
		IncludedObjectVersions: result.IncludedObjectVersions,
		Filter:                 result.Filter,
		OptionalFields:         result.OptionalFields,
		Destination:            result.Destination,
	}
	_, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTencentCosClient(bucket).Bucket.PostInventory(ctx, id, inventoryOpt)
	if err != nil {
		return err
	}
	items := strings.Split(result.Destination.Bucket, "::")
	targetBucket := items[len(items)-1]
	targetBucketItems := strings.Split(targetBucket, "-")
	objectId := fmt.Sprintf("%s/%s/%s/%s/%s/manifest.json", result.Destination.Prefix, targetBucketItems[len(targetBucketItems)-1], strings.Join(targetBucketItems[:len(targetBucketItems)-1], "-"), id, time.Now().Format("20060102"))
	err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		resp, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTencentCosClient(targetBucket).Object.Head(ctx, objectId, nil)
		if resp.StatusCode == 404 {
			return resource.RetryableError(fmt.Errorf("Inventory still creating!"))
		}
		if e != nil {
			return tccommon.RetryError(e)
		}

		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s PostInventory: object not exist, reason:%+v", logId, err)
		return err
	}
	d.SetId(id)

	return resourceTencentCloudCosBucketGenerateInventoryImmediatelyOperationRead(d, meta)
}

func resourceTencentCloudCosBucketGenerateInventoryImmediatelyOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_generate_inventory_immediately_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCosBucketGenerateInventoryImmediatelyOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_generate_inventory_immediately_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
