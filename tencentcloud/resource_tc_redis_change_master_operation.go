/*
Provides a resource to create a redis change_master_operation

Example Usage

```hcl
resource "tencentcloud_redis_change_master_operation" "change_master_operation" {
  instance_id = "crs-c1nl9rpv"
  group_id = "crs-rpl-c1nl9rpv"
}
```

Import

redis change_master_operation can be imported using the id, e.g.

```
terraform import tencentcloud_redis_change_master_operation.change_master_operation change_master_operation_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
	"time"
)

func resourceTencentCloudRedisChangeMasterOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisChangeMasterOperationCreate,
		Read:   resourceTencentCloudRedisChangeMasterOperationRead,
		Delete: resourceTencentCloudRedisChangeMasterOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of group.",
			},
		},
	}
}

func resourceTencentCloudRedisChangeMasterOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_change_master_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = redis.NewChangeMasterInstanceRequest()
		response   = redis.NewChangeMasterInstanceResponse()
		groupId    string
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
		request.GroupId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().ChangeMasterInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate redis changeMasterOperation failed, reason:%+v", logId, err)
		return err
	}

	groupId = *response.Response.GroupId
	d.SetId(strings.Join([]string{groupId, instanceId}, FILED_SP))

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"succeed"}, 30*readRetryTimeout, time.Second, service.RedisChangeMasterOperationStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudRedisChangeMasterOperationRead(d, meta)
}

func resourceTencentCloudRedisChangeMasterOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_change_master_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudRedisChangeMasterOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_change_master_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
