/*
Provides a resource to create a cam access_key

Example Usage

```hcl
resource "tencentcloud_cam_access_key" "access_key" {
  target_uin = &lt;nil&gt;
}
```

Import

cam access_key can be imported using the id, e.g.

```
terraform import tencentcloud_cam_access_key.access_key access_key_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCamAccessKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamAccessKeyCreate,
		Read:   resourceTencentCloudCamAccessKeyRead,
		Update: resourceTencentCloudCamAccessKeyUpdate,
		Delete: resourceTencentCloudCamAccessKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"target_uin": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Specify user Uin, if not filled, the access key is created for the current user by default.",
			},
		},
	}
}

func resourceTencentCloudCamAccessKeyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_access_key.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = cam.NewCreateAccessKeyRequest()
		response    = cam.NewCreateAccessKeyResponse()
		accessKeyId string
	)
	if v, ok := d.GetOkExists("target_uin"); ok {
		request.TargetUin = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().CreateAccessKey(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cam AccessKey failed, reason:%+v", logId, err)
		return err
	}

	accessKeyId = *response.Response.AccessKeyId
	d.SetId(helper.String(accessKeyId))

	return resourceTencentCloudCamAccessKeyRead(d, meta)
}

func resourceTencentCloudCamAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_access_key.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	accessKeyId := d.Id()

	AccessKey, err := service.DescribeCamAccessKeyById(ctx, accessKeyId)
	if err != nil {
		return err
	}

	if AccessKey == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CamAccessKey` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if AccessKey.TargetUin != nil {
		_ = d.Set("target_uin", AccessKey.TargetUin)
	}

	return nil
}

func resourceTencentCloudCamAccessKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_access_key.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cam.NewUpdateAccessKeyRequest()

	accessKeyId := d.Id()

	request.AccessKeyId = &accessKeyId

	immutableArgs := []string{"target_uin"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("target_uin") {
		if v, ok := d.GetOkExists("target_uin"); ok {
			request.TargetUin = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().UpdateAccessKey(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cam AccessKey failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCamAccessKeyRead(d, meta)
}

func resourceTencentCloudCamAccessKeyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_access_key.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}
	accessKeyId := d.Id()

	if err := service.DeleteCamAccessKeyById(ctx, accessKeyId); err != nil {
		return err
	}

	return nil
}
