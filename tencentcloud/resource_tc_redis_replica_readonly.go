/*
Provides a resource to create a redis replica_readonly

Example Usage

```hcl
resource "tencentcloud_redis_replica_readonly" "replica_readonly" {
  instance_id = "crs-c1nl9rpv"
  readonly_policy =
}
```

Import

redis replica_readonly can be imported using the id, e.g.

```
terraform import tencentcloud_redis_replica_readonly.replica_readonly replica_readonly_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudRedisReplicaReadonly() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisReplicaReadonlyCreate,
		Read:   resourceTencentCloudRedisReplicaReadonlyRead,
		Update: resourceTencentCloudRedisReplicaReadonlyUpdate,
		Delete: resourceTencentCloudRedisReplicaReadonlyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"readonly_policy": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Routing policy: Enter `master` or `replication`, which indicates the master node or slave node.",
			},
		},
	}
}

func resourceTencentCloudRedisReplicaReadonlyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_replica_readonly.create")()
	defer inconsistentCheck(d, meta)()

	var (
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudRedisReplicaReadonlyUpdate(d, meta)
}

func resourceTencentCloudRedisReplicaReadonlyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_replica_readonly.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	instance, err := service.DescribeRedisInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisReplicaReadonly` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instance.InstanceId != nil {
		_ = d.Set("instance_id", instance.InstanceId)
	}

	if instance.SlaveReadWeight != nil {
		if *instance.SlaveReadWeight == 100 {
			_ = d.Set("readonly_policy", []string{"master"})
		}
	}

	return nil
}

func resourceTencentCloudRedisReplicaReadonlyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_replica_readonly.update")()
	defer inconsistentCheck(d, meta)()

	// logId := getLogId(contextNil)

	// var (
	// 	disableRequest  = redis.NewDisableReplicaReadonlyRequest()
	// 	disableResponse = redis.NewDisableReplicaReadonlyResponse()
	// 	enableRequest   = redis.NewEnableReplicaReadonlyRequest()
	// 	enableResponse  = redis.NewEnableReplicaReadonlyResponse()
	// )

	// replicaReadonlyId := d.Id()

	// request.InstanceId = &instanceId

	// immutableArgs := []string{"instance_id", "readonly_policy"}

	// for _, v := range immutableArgs {
	// 	if d.HasChange(v) {
	// 		return fmt.Errorf("argument `%s` cannot be changed", v)
	// 	}
	// }

	// err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
	// 	result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().DisableReplicaReadonly(request)
	// 	if e != nil {
	// 		return retryError(e)
	// 	} else {
	// 		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
	// 	}
	// 	return nil
	// })
	// if err != nil {
	// 	log.Printf("[CRITAL]%s update redis replicaReadonly failed, reason:%+v", logId, err)
	// 	return err
	// }

	// service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	// conf := BuildStateChangeConf([]string{}, []string{"succeed"}, 60*readRetryTimeout, time.Second, service.RedisReplicaReadonlyStateRefreshFunc(d.Id(), []string{}))

	// if _, e := conf.WaitForState(); e != nil {
	// 	return e
	// }

	// service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	// conf := BuildStateChangeConf([]string{}, []string{"succeed"}, 60*readRetryTimeout, time.Second, service.RedisReplicaReadonlyStateRefreshFunc(d.Id(), []string{}))

	// if _, e := conf.WaitForState(); e != nil {
	// 	return e
	// }

	return resourceTencentCloudRedisReplicaReadonlyRead(d, meta)
}

func resourceTencentCloudRedisReplicaReadonlyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_replica_readonly.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
