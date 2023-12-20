package as

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAsNotification() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsNotificationCreate,
		Read:   resourceTencentCloudAsNotificationRead,
		Update: resourceTencentCloudAsNotificationUpdate,
		Delete: resourceTencentCloudAsNotificationDelete,

		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of a scaling group.",
			},
			"notification_types": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "A list of Notification Types that trigger notifications. Acceptable values are `SCALE_OUT_FAILED`, `SCALE_IN_SUCCESSFUL`, `SCALE_IN_FAILED`, `REPLACE_UNHEALTHY_INSTANCE_SUCCESSFUL` and `REPLACE_UNHEALTHY_INSTANCE_FAILED`.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: tccommon.ValidateAllowedStringValue(SCALING_GROUP_NOTIFICATION_TYPE),
				},
			},
			"notification_user_group_ids": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A group of user IDs to be notified.",
			},
		},
	}
}

func resourceTencentCloudAsNotificationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_notification.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := as.NewCreateNotificationConfigurationRequest()
	request.AutoScalingGroupId = helper.String(d.Get("scaling_group_id").(string))
	notificationTypes := d.Get("notification_types").([]interface{})
	request.NotificationTypes = make([]*string, 0, len(notificationTypes))
	for _, value := range notificationTypes {
		request.NotificationTypes = append(request.NotificationTypes, helper.String(value.(string)))
	}
	userGroupIds := d.Get("notification_user_group_ids").([]interface{})
	request.NotificationUserGroupIds = make([]*string, 0, len(userGroupIds))
	for _, value := range userGroupIds {
		request.NotificationUserGroupIds = append(request.NotificationUserGroupIds, helper.String(value.(string)))
	}

	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().CreateNotificationConfiguration(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.AutoScalingNotificationId == nil {
		return fmt.Errorf("scaling policy id is nil")
	}
	d.SetId(*response.Response.AutoScalingNotificationId)

	return resourceTencentCloudAsNotificationRead(d, meta)
}

func resourceTencentCloudAsNotificationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_notification.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	notificationId := d.Id()
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		notification, has, e := asService.DescribeNotificationById(ctx, notificationId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if has == 0 {
			d.SetId("")
			return nil
		}
		_ = d.Set("scaling_group_id", *notification.AutoScalingGroupId)
		_ = d.Set("notification_types", helper.StringsInterfaces(notification.NotificationTypes))
		_ = d.Set("notification_user_group_ids", helper.StringsInterfaces(notification.NotificationUserGroupIds))
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudAsNotificationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_notification.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := as.NewModifyNotificationConfigurationRequest()
	notificationId := d.Id()
	request.AutoScalingNotificationId = &notificationId
	if d.HasChange("notification_types") {
		notificationTypes := d.Get("notification_types").([]interface{})
		request.NotificationTypes = make([]*string, 0, len(notificationTypes))
		for _, value := range notificationTypes {
			request.NotificationTypes = append(request.NotificationTypes, helper.String(value.(string)))
		}
	}
	if d.HasChange("notification_user_group_ids") {
		userGroupIds := d.Get("notification_user_group_ids").([]interface{})
		request.NotificationUserGroupIds = make([]*string, 0, len(userGroupIds))
		for _, value := range userGroupIds {
			request.NotificationUserGroupIds = append(request.NotificationUserGroupIds, helper.String(value.(string)))
		}
	}

	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().ModifyNotificationConfiguration(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func resourceTencentCloudAsNotificationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_notification.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	notificationId := d.Id()
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := asService.DeleteNotification(ctx, notificationId)
	if err != nil {
		return err
	}

	return nil
}
