/*
Provides a resource to create a tdmq rabbitmq_user

Example Usage

```hcl
resource "tencentcloud_tdmq_rabbitmq_user" "rabbitmq_user" {
  instance_id = ""
  user = ""
  password = ""
  description = ""
  max_connections =
  max_channels =
}
```

Import

tdmq rabbitmq_user can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_rabbitmq_user.rabbitmq_user rabbitmq_user_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTdmqRabbitmqUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqRabbitmqUserCreate,
		Read:   resourceTencentCloudTdmqRabbitmqUserRead,
		Update: resourceTencentCloudTdmqRabbitmqUserUpdate,
		Delete: resourceTencentCloudTdmqRabbitmqUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster instance ID.",
			},

			"user": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Username, used when logging in.",
			},

			"password": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Password, used when logging in.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Describe.",
			},

			"max_connections": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The maximum number of connections for this user, if not filled in, there is no limit.",
			},

			"max_channels": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The maximum number of channels for this user, if not filled in, there is no limit.",
			},
		},
	}
}

func resourceTencentCloudTdmqRabbitmqUserCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_user.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = tdmq.NewCreateRabbitMQUserRequest()
		response   = tdmq.NewCreateRabbitMQUserResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user"); ok {
		request.User = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("max_connections"); ok {
		request.MaxConnections = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("max_channels"); ok {
		request.MaxChannels = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().CreateRabbitMQUser(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tdmq rabbitmqUser failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudTdmqRabbitmqUserRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqUserRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_user.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	rabbitmqUserId := d.Id()

	rabbitmqUser, err := service.DescribeTdmqRabbitmqUserById(ctx, instanceId)
	if err != nil {
		return err
	}

	if rabbitmqUser == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqRabbitmqUser` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if rabbitmqUser.InstanceId != nil {
		_ = d.Set("instance_id", rabbitmqUser.InstanceId)
	}

	if rabbitmqUser.User != nil {
		_ = d.Set("user", rabbitmqUser.User)
	}

	if rabbitmqUser.Password != nil {
		_ = d.Set("password", rabbitmqUser.Password)
	}

	if rabbitmqUser.Description != nil {
		_ = d.Set("description", rabbitmqUser.Description)
	}

	if rabbitmqUser.MaxConnections != nil {
		_ = d.Set("max_connections", rabbitmqUser.MaxConnections)
	}

	if rabbitmqUser.MaxChannels != nil {
		_ = d.Set("max_channels", rabbitmqUser.MaxChannels)
	}

	return nil
}

func resourceTencentCloudTdmqRabbitmqUserUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_user.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tdmq.NewModifyRabbitMQUserRequest()

	rabbitmqUserId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "user", "password", "description", "max_connections", "max_channels"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("instance_id") {
		if v, ok := d.GetOk("instance_id"); ok {
			request.InstanceId = helper.String(v.(string))
		}
	}

	if d.HasChange("user") {
		if v, ok := d.GetOk("user"); ok {
			request.User = helper.String(v.(string))
		}
	}

	if d.HasChange("password") {
		if v, ok := d.GetOk("password"); ok {
			request.Password = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("max_connections") {
		if v, ok := d.GetOkExists("max_connections"); ok {
			request.MaxConnections = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("max_channels") {
		if v, ok := d.GetOkExists("max_channels"); ok {
			request.MaxChannels = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().ModifyRabbitMQUser(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tdmq rabbitmqUser failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqRabbitmqUserRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqUserDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_rabbitmq_user.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}
	rabbitmqUserId := d.Id()

	if err := service.DeleteTdmqRabbitmqUserById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
