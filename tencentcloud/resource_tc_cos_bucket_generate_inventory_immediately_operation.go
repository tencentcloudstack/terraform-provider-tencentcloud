/*
Provides a resource to generate a cos bucket inventory immediately

Example Usage

```hcl
resource "tencentcloud_cos_bucket_generate_inventory_immediately_operation" "generate_inventory_immediately" {
    inventory_id = "test"
    bucket = "keep-test-xxxxxx"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func resourceTencentCloudCosBucketGenerateInventoryImmediatelyOperation() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_cos_bucket_generate_inventory_immediately_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	inventoryId := d.Get("inventory_id").(string)
	bucket := d.Get("bucket").(string)
	result, _, err := meta.(*TencentCloudClient).apiV3Conn.UseTencentCosClient(bucket).Bucket.GetInventory(ctx, inventoryId)
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
	_, err = meta.(*TencentCloudClient).apiV3Conn.UseTencentCosClient(bucket).Bucket.PostInventory(ctx, id, inventoryOpt)
	if err != nil {
		return err
	}
	items := strings.Split(result.Destination.Bucket, "::")
	targetBucket := items[len(items)-1]
	targetBucketItems := strings.Split(targetBucket, "-")
	objectId := fmt.Sprintf("%s/%s/%s/%s/%s/manifest.json", result.Destination.Prefix, targetBucketItems[len(targetBucketItems)-1], strings.Join(targetBucketItems[:len(targetBucketItems)-1], "-"), id, time.Now().Format("20060102"))
	err = resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
		resp, e := meta.(*TencentCloudClient).apiV3Conn.UseTencentCosClient(targetBucket).Object.Head(ctx, objectId, nil)
		if resp.StatusCode == 404 {
			return resource.RetryableError(fmt.Errorf("Inventory still creating!"))
		}
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_cos_bucket_generate_inventory_immediately_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCosBucketGenerateInventoryImmediatelyOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cos_bucket_generate_inventory_immediately_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
