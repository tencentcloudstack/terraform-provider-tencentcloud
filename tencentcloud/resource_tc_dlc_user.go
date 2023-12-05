package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDlcUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcUserCreate,
		Read:   resourceTencentCloudDlcUserRead,
		Update: resourceTencentCloudDlcUserUpdate,
		Delete: resourceTencentCloudDlcUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "The sub-user uin that needs to be authorized.",
			},

			"user_description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User description information, easy to distinguish between different users.",
			},

			"user_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User Type. `ADMIN` or `COMMONN`.",
			},

			"user_alias": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User alias, the character length is less than 50.",
			},

			"work_group_ids": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "A collection of workgroup IDs bound to the user.",
			},
		},
	}
}

func resourceTencentCloudDlcUserCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_user.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = dlc.NewCreateUserRequest()
		userId  string
	)
	if v, ok := d.GetOk("user_id"); ok {
		userId = v.(string)
		request.UserId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_description"); ok {
		request.UserDescription = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_type"); ok {
		request.UserType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_alias"); ok {
		request.UserAlias = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDlcClient().CreateUser(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dlc user failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(userId)

	return resourceTencentCloudDlcUserRead(d, meta)
}

func resourceTencentCloudDlcUserRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_user.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DlcService{client: meta.(*TencentCloudClient).apiV3Conn}

	userId := d.Id()

	user, err := service.DescribeDlcUserById(ctx, userId)
	if err != nil {
		return err
	}

	if user == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DlcUser` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if user.UserId != nil {
		_ = d.Set("user_id", user.UserId)
	}

	if user.UserDescription != nil {
		_ = d.Set("user_description", user.UserDescription)
	}

	if user.UserType != nil {
		_ = d.Set("user_type", user.UserType)
	}

	if user.UserAlias != nil {
		_ = d.Set("user_alias", user.UserAlias)
	}

	if user.WorkGroupSet != nil {
		workGroups := make([]*int64, len(user.WorkGroupSet))

		for _, workGroup := range user.WorkGroupSet {
			workGroups = append(workGroups, workGroup.WorkGroupId)
		}
		_ = d.Set("work_group_ids", workGroups)
	}

	return nil
}

func resourceTencentCloudDlcUserUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_user.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dlc.NewModifyUserRequest()

	userId := d.Id()

	request.UserId = &userId

	immutableArgs := []string{"user_type", "user_alias"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("user_description") {
		if v, ok := d.GetOk("user_description"); ok {
			request.UserDescription = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDlcClient().ModifyUser(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dlc user failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDlcUserRead(d, meta)
}

func resourceTencentCloudDlcUserDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_user.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DlcService{client: meta.(*TencentCloudClient).apiV3Conn}
	userId := d.Id()

	if err := service.DeleteDlcUserById(ctx, userId); err != nil {
		return err
	}

	return nil
}
