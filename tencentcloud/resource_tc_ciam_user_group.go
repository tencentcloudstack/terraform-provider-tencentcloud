package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ciam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ciam/v20220331"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCiamUserGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiamUserGroupCreate,
		Read:   resourceTencentCloudCiamUserGroupRead,
		Update: resourceTencentCloudCiamUserGroupUpdate,
		Delete: resourceTencentCloudCiamUserGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_store_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "User Store ID.",
			},
			"display_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "User Group Name.",
			},
			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User Group Description.",
			},
		},
	}
}

func resourceTencentCloudCiamUserGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ciam_user_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = ciam.NewCreateUserGroupRequest()
		response    = ciam.NewCreateUserGroupResponse()
		userStoreId string
		userGroupId string
	)
	if v, ok := d.GetOk("user_store_id"); ok {
		userStoreId = v.(string)
		request.UserStoreId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("display_name"); ok {
		request.DisplayName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiamClient().CreateUserGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ciam userGroup failed, reason:%+v", logId, err)
		return err
	}

	userGroupId = *response.Response.UserGroupId

	d.SetId(userStoreId + FILED_SP + userGroupId)

	return resourceTencentCloudCiamUserGroupRead(d, meta)
}

func resourceTencentCloudCiamUserGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ciam_user_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiamService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	userStoreId := idSplit[0]
	userGroupId := idSplit[1]

	userGroup, err := service.DescribeCiamUserGroupById(ctx, userStoreId, userGroupId)
	if err != nil {
		return err
	}

	if userGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CiamUserGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if userGroup.DisplayName != nil {
		_ = d.Set("display_name", userGroup.DisplayName)
	}

	if userGroup.UserStoreId != nil {
		_ = d.Set("user_store_id", userGroup.UserStoreId)
	}

	if userGroup.Description != nil {
		_ = d.Set("description", userGroup.Description)
	}

	return nil
}

func resourceTencentCloudCiamUserGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ciam_user_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := ciam.NewUpdateUserGroupRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	userStoreId := idSplit[0]
	userGroupId := idSplit[1]

	request.UserStoreId = &userStoreId
	request.UserGroupId = &userGroupId

	needChange := false
	mutableArgs := []string{
		"display_name", "user_store_id", "description",
	}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {

		if v, ok := d.GetOk("display_name"); ok {
			request.DisplayName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("user_store_id"); ok {
			request.UserStoreId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiamClient().UpdateUserGroup(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update ciam userGroup failed, reason:%+v", logId, err)
			return err
		}

	}
	return resourceTencentCloudCiamUserGroupRead(d, meta)
}

func resourceTencentCloudCiamUserGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ciam_user_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiamService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	userStoreId := idSplit[0]
	userGroupId := idSplit[1]

	if err := service.DeleteCiamUserGroupById(ctx, userStoreId, userGroupId); err != nil {
		return err
	}

	return nil
}
