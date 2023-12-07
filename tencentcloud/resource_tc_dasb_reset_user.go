package tencentcloud

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDasbResetUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDasbResetUserCreate,
		Read:   resourceTencentCloudDasbResetUserRead,
		Delete: resourceTencentCloudDasbResetUserDelete,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "User Id.",
			},
		},
	}
}

func resourceTencentCloudDasbResetUserCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_reset_user.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		request = dasb.NewResetUserRequest()
		userId  string
	)

	if v, ok := d.GetOk("user_id"); ok {
		request.IdSet = append(request.IdSet, helper.IntUint64(v.(int)))
		userId = strconv.Itoa(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDasbClient().ResetUser(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate dasb ResetUser failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(userId)

	return resourceTencentCloudDasbResetUserRead(d, meta)
}

func resourceTencentCloudDasbResetUserRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_reset_user.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDasbResetUserDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_reset_user.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
