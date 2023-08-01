/*
Provides a resource for an AS (Auto scaling) notification.

Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "as"
}

data "tencentcloud_images" "image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "TencentOS Server 3.2 (Final)"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-example"
  cidr_block        = "10.0.0.0/16"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
}

resource "tencentcloud_as_scaling_config" "example" {
  configuration_name = "tf-example"
  image_id           = data.tencentcloud_images.image.images.0.image_id
  instance_types     = ["SA1.SMALL1", "SA2.SMALL1", "SA2.SMALL2", "SA2.SMALL4"]
  instance_name_settings {
    instance_name = "test-ins-name"
  }
}

resource "tencentcloud_as_scaling_group" "example" {
  scaling_group_name = "tf-example"
  configuration_id   = tencentcloud_as_scaling_config.example.id
  max_size           = 1
  min_size           = 0
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_ids         = [tencentcloud_subnet.subnet.id]
}

resource "tencentcloud_cam_group" "example" {
  name   = "tf-example"
  remark = "desc."
}

resource "tencentcloud_as_notification" "as_notification" {
  scaling_group_id            = tencentcloud_as_scaling_group.example.id
  notification_types          = [
    "SCALE_OUT_SUCCESSFUL", "SCALE_OUT_FAILED", "SCALE_IN_FAILED", "REPLACE_UNHEALTHY_INSTANCE_SUCCESSFUL", "REPLACE_UNHEALTHY_INSTANCE_FAILED"
  ]
  notification_user_group_ids = [tencentcloud_cam_group.example.id]
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
			"notification_types": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "A list of Notification Types that trigger notifications. Acceptable values are `SCALE_OUT_FAILED`, `SCALE_IN_SUCCESSFUL`, `SCALE_IN_FAILED`, `REPLACE_UNHEALTHY_INSTANCE_SUCCESSFUL` and `REPLACE_UNHEALTHY_INSTANCE_FAILED`.",
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
	defer logElapsed("resource.tencentcloud_as_notification.create")()

	logId := getLogId(contextNil)

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
	defer logElapsed("resource.tencentcloud_as_notification.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	notificationId := d.Id()
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		notification, has, e := asService.DescribeNotificationById(ctx, notificationId)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_as_notification.update")()

	logId := getLogId(contextNil)

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
	defer logElapsed("resource.tencentcloud_as_notification.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

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
