/*
Provides a resource to create a redis modfiy_instance_password

Example Usage

```hcl
resource "tencentcloud_redis_modfiy_instance_password" "modfiy_instance_password" {
  instance_id  = "crs-c1nl9rpv"
  old_password = ""
  password 	   = ""
}
```

*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudRedisModfiyInstancePassword() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisModfiyInstancePasswordCreate,
		Read:   resourceTencentCloudRedisModfiyInstancePasswordRead,
		Update: resourceTencentCloudRedisModfiyInstancePasswordUpdate,
		Delete: resourceTencentCloudRedisModfiyInstancePasswordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"old_password": {
				Required:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
				Description: "The old password for the instance.",
			},

			"password": {
				Required:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
				Description: "The password for the instance.",
			},
		},
	}
}

func resourceTencentCloudRedisModfiyInstancePasswordCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_modfiy_instance_password.create")()
	defer inconsistentCheck(d, meta)()

	var (
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudRedisModfiyInstancePasswordUpdate(d, meta)
}

func resourceTencentCloudRedisModfiyInstancePasswordRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_modfiy_instance_password.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	modfiyInstancePassword, err := service.DescribeRedisInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if modfiyInstancePassword == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisModfiyInstancePassword` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if modfiyInstancePassword.InstanceId != nil {
		_ = d.Set("instance_id", modfiyInstancePassword.InstanceId)
	}

	return nil
}

func resourceTencentCloudRedisModfiyInstancePasswordUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_modfiy_instance_password.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := redis.NewModfiyInstancePasswordRequest()
	response := redis.NewModfiyInstancePasswordResponse()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if v, ok := d.GetOk("old_password"); ok {
		request.OldPassword = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().ModfiyInstancePassword(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update redis modfiyInstancePassword failed, reason:%+v", logId, err)
		return err
	}

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
	taskId := *response.Response.TaskId
	err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
		ok, err := service.DescribeTaskInfo(ctx, instanceId, taskId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("change password is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s redis change password fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudRedisModfiyInstancePasswordRead(d, meta)
}

func resourceTencentCloudRedisModfiyInstancePasswordDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_modfiy_instance_password.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
