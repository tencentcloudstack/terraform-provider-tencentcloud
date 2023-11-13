/*
Provides a resource to create a redis switch_master

Example Usage

```hcl
resource "tencentcloud_redis_switch_master" "switch_master" {
  instance_id = "crs-c1nl9rpv"
  group_id = 10001
}
```

Import

redis switch_master can be imported using the id, e.g.

```
terraform import tencentcloud_redis_switch_master.switch_master switch_master_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"log"
	"time"
)

func resourceTencentCloudRedisSwitchMaster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisSwitchMasterCreate,
		Read:   resourceTencentCloudRedisSwitchMasterRead,
		Update: resourceTencentCloudRedisSwitchMasterUpdate,
		Delete: resourceTencentCloudRedisSwitchMasterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"group_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Replication group ID, required for multi-AZ instances.",
			},
		},
	}
}

func resourceTencentCloudRedisSwitchMasterCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_switch_master.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudRedisSwitchMasterUpdate(d, meta)
}

func resourceTencentCloudRedisSwitchMasterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_switch_master.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	switchMasterId := d.Id()

	switchMaster, err := service.DescribeRedisSwitchMasterById(ctx, instanceId)
	if err != nil {
		return err
	}

	if switchMaster == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisSwitchMaster` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if switchMaster.InstanceId != nil {
		_ = d.Set("instance_id", switchMaster.InstanceId)
	}

	if switchMaster.GroupId != nil {
		_ = d.Set("group_id", switchMaster.GroupId)
	}

	return nil
}

func resourceTencentCloudRedisSwitchMasterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_switch_master.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := redis.NewChangeReplicaToMasterRequest()

	switchMasterId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "group_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().ChangeReplicaToMaster(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update redis switchMaster failed, reason:%+v", logId, err)
		return err
	}

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"succeed"}, 60*readRetryTimeout, time.Second, service.RedisSwitchMasterStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudRedisSwitchMasterRead(d, meta)
}

func resourceTencentCloudRedisSwitchMasterDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_switch_master.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
