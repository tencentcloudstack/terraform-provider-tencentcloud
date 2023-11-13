/*
Provides a resource to create a cvm modify_image_share_permission

Example Usage

```hcl
resource "tencentcloud_cvm_modify_image_share_permission" "modify_image_share_permission" {
  image_id = "img-gvbnzy6f"
  account_ids =
  permission = "SHARE"
}
```

Import

cvm modify_image_share_permission can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_modify_image_share_permission.modify_image_share_permission modify_image_share_permission_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCvmModifyImageSharePermission() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmModifyImageSharePermissionCreate,
		Read:   resourceTencentCloudCvmModifyImageSharePermissionRead,
		Delete: resourceTencentCloudCvmModifyImageSharePermissionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"image_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Image ID such as `img-gvbnzy6f`. You can only specify an image in the NORMAL state.",
			},

			"account_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of account IDs with which an image is shared.",
			},

			"permission": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Operations. Valid values: SHARE, sharing an image; CANCEL, cancelling an image sharing.",
			},
		},
	}
}

func resourceTencentCloudCvmModifyImageSharePermissionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_modify_image_share_permission.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cvm.NewModifyImageSharePermissionRequest()
		response = cvm.NewModifyImageSharePermissionResponse()
		imageId  string
	)
	if v, ok := d.GetOk("image_id"); ok {
		imageId = v.(string)
		request.ImageId = helper.String(v.(string))
	}

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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().ModifyImageSharePermission(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm modifyImageSharePermission failed, reason:%+v", logId, err)
		return err
	}

	imageId = *response.Response.ImageId
	d.SetId(imageId)

	return resourceTencentCloudCvmModifyImageSharePermissionRead(d, meta)
}

func resourceTencentCloudCvmModifyImageSharePermissionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_modify_image_share_permission.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmModifyImageSharePermissionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_modify_image_share_permission.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
