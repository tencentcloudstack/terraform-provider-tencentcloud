/*
Provides a resource for an AS (Auto scaling) notification.

Example Usage

```hcl
resource "tencentcloud_autoscaling_notification" "aslab" {
  scaling_group_id              = "sg-12af45"
  notification_type             = ["SCALE_OUT_FAILED", "SCALE_IN_SUCCESSFUL", "SCALE_IN_FAILED", "REPLACE_UNHEALTHY_INSTANCE_FAILED"]
  notification_user_group_ids   = ["ASGID"]
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
)

func resourceTencentCloudAsNotification() *schema.Resource {
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
			"notification_type": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "A list of Notification Types that trigger notifications. Acceptable values are SCALE_OUT_FAILED, SCALE_IN_SUCCESSFUL, SCALE_IN_FAILED, REPLACE_UNHEALTHY_INSTANCE_SUCCESSFUL and REPLACE_UNHEALTHY_INSTANCE_FAILED.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateAllowedStringValue(SCALING_GROUP_NOTIFICATION_TYPE),
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
	logId := GetLogId(nil)

	request := as.NewCreateNotificationConfigurationRequest()
	request.AutoScalingGroupId = stringToPointer(d.Get("scaling_group_id").(string))
	notificationTypes := d.Get("notification_types").([]interface{})
	request.NotificationTypes = make([]*string, 0, len(notificationTypes))
	for _, value := range notificationTypes {
		request.NotificationTypes = append(request.NotificationTypes, stringToPointer(value.(string)))
	}
	userGroupIds := d.Get("notification_user_group_ids").([]interface{})
	request.NotificationUserGroupIds = make([]*string, 0, len(userGroupIds))
	for _, value := range userGroupIds {
		request.NotificationUserGroupIds = append(request.NotificationUserGroupIds, stringToPointer(value.(string)))
	}

	response, err := meta.(*TencentCloudClient).apiV3Conn.UseAsClient().CreateNotificationConfiguration(request)
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
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	notificationId := d.Id()
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	notification, err := asService.DescribeNotificationById(ctx, notificationId)
	if err != nil {
		return err
	}

	d.Set("scaling_group_id", *notification.AutoScalingGroupId)
	d.Set("notification_type", flattenStringList(notification.NotificationTypes))
	d.Set("notification_user_group_ids", flattenStringList(notification.NotificationUserGroupIds))

	return nil
}

func resourceTencentCloudAsNotificationUpdate(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)

	request := as.NewModifyNotificationConfigurationRequest()
	notificationId := d.Id()
	request.AutoScalingNotificationId = &notificationId
	if d.HasChange("notification_type") {
		notificationTypes := d.Get("notification_types").([]interface{})
		request.NotificationTypes = make([]*string, 0, len(notificationTypes))
		for _, value := range notificationTypes {
			request.NotificationTypes = append(request.NotificationTypes, stringToPointer(value.(string)))
		}
	}
	if d.HasChange("notification_user_group_ids") {
		userGroupIds := d.Get("notification_user_group_ids").([]interface{})
		request.NotificationUserGroupIds = make([]*string, 0, len(userGroupIds))
		for _, value := range userGroupIds {
			request.NotificationUserGroupIds = append(request.NotificationUserGroupIds, stringToPointer(value.(string)))
		}
	}

	response, err := meta.(*TencentCloudClient).apiV3Conn.UseAsClient().ModifyNotificationConfiguration(request)
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
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	notificationId := d.Id()
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := asService.DeleteNotification(ctx, notificationId)
	if err != nil {
		return err
	}

	return nil
}
