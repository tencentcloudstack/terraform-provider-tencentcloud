/*
Provides a resource to create a Redis instance and set its attributes.

Example Usage

```hcl
resource "tencentcloud_redis_instance" "redis_instance_test"{
  availability_zone = "ap-hongkong-3"
  type              = "master_slave_redis"
  password          = "test12345789"
  mem_size          = 8192
  name              = "terrform_test"
  port              = 6379
}
```

Import

Redis instance can be imported, e.g.

```
$ terraform import tencentcloud_redis_instance.redislab redis-id
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func resourceTencentCloudRedisInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisInstanceCreate,
		Read:   resourceTencentCloudRedisInstanceRead,
		Update: resourceTencentCloudRedisInstanceUpdate,
		Delete: resourceTencentCloudRedisInstanceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The available zone ID of an instance to be created., refer to tencentcloud_redis_zone_config.list",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Instance name.",
			},
			"type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  REDIS_NAMES[REDIS_VERSION_MASTER_SLAVE_REDIS],
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					for _, name := range REDIS_NAMES {
						if name == value {
							return
						}
					}
					errors = append(errors, fmt.Errorf("this redis type %s not support now.", value))
					return
				},
				Description: "Instance type. Available values: master_slave_redis.",
			},
			"password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validateMysqlPassword,
				Description:  "Password for a Redis user，which should be 8 to 16 characters.",
			},
			"mem_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The memory volume of an available instance(in MB), refer to tencentcloud_redis_zone_config.list[zone].mem_sizes",
			},
			"vpc_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateStringLengthInRange(1, 100),
				Description:  "ID of the vpc with which the instance is to be associated.",
			},
			"subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateStringLengthInRange(1, 100),
				Description:  "Specifies which subnet the instance should belong to.",
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
				Description: "ID of security group. If both vpc_id and subnet_id are not set, this argument should not be set either. ",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Specifies which project the instance should belong to.",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     6379,
				Description: "The port used to access a redis instance. The default value is 6379. And this value can't be changed after creation, or the Redis instance will be recreated.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Instance tags.",
			},

			// Computed values
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address of an instance.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current status of an instance，maybe: init, processing, online, isolate and todelete.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: " The time when the instance was created.",
			},
		},
	}
}

func resourceTencentCloudRedisInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_instance.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	redisService := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
	tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
	region := meta.(*TencentCloudClient).apiV3Conn.Region

	availabilityZone := d.Get("availability_zone").(string)
	redisName := d.Get("name").(string)
	redisType := d.Get("type").(string)
	password := d.Get("password").(string)
	memSize := d.Get("mem_size").(int)
	vpcId := d.Get("vpc_id").(string)
	subnetId := d.Get("subnet_id").(string)
	securityGroups := d.Get("security_groups").(*schema.Set).List()
	projectId := d.Get("project_id").(int)
	port := d.Get("port").(int)
	tags := getTags(d, "tags")

	if availabilityZone != "" {
		if !strings.Contains(availabilityZone, region) {
			return fmt.Errorf("zone[%s] not in region[%s]", availabilityZone, region)
		}
	}

	requestSecurityGroup := make([]string, 0, len(securityGroups))

	for _, v := range securityGroups {
		requestSecurityGroup = append(requestSecurityGroup, v.(string))
	}

	dealId, err := redisService.CreateInstances(ctx,
		availabilityZone,
		redisType,
		password,
		vpcId,
		subnetId,
		redisName,
		int64(memSize),
		int64(projectId),
		int64(port),
		requestSecurityGroup)

	if err != nil {
		return err
	}

	if dealId == "" {
		return fmt.Errorf("redis api CreateInstances return empty redis id")
	}
	var redisId = dealId
	err = resource.Retry(60*time.Minute, func() *resource.RetryError {
		has, online, _, err := redisService.CheckRedisCreateOk(ctx, dealId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if !has {
			return resource.NonRetryableError(fmt.Errorf("redis instance not exists."))
		}
		if online {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("create redis task is processing"))
	})

	if err != nil {
		log.Printf("[CRITAL]%s create redis task fail, reason:%s\n ", logId, err.Error())
		return err
	}
	d.SetId(redisId)

	if len(tags) > 0 {
		resourceName := buildResourceNameWithUin("redis", "instance", redisId, region)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudRedisInstanceRead(d, meta)
}

func resourceTencentCloudRedisInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_instance.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
	has, _, info, err := service.CheckRedisCreateOk(ctx, d.Id())
	if info != nil {
		if *info.Status == REDIS_STATUS_ISOLATE || *info.Status == REDIS_STATUS_TODELETE {
			d.SetId("")
			return nil
		}
	}
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	statusName := REDIS_STATUS[*info.Status]
	if statusName == "" {
		err = fmt.Errorf("redis read unkwnow status %d", *info.ZoneId)
		log.Printf("[CRITAL]%s redis read status name error, reason:%s\n ", logId, err.Error())
		return err
	}
	d.Set("status", statusName)

	d.Set("name", *info.InstanceName)

	zoneName := REDIS_ZONE_ID2NAME[*info.ZoneId]
	if zoneName == "" {
		err = fmt.Errorf("redis read unkwnow zoneid %d", *info.ZoneId)
		log.Printf("[CRITAL]%s redis read zone name error, reason:%s\n ", logId, err.Error())
		return err
	}
	d.Set("availability_zone", zoneName)

	typeName := REDIS_NAMES[*info.Type]
	if typeName == "" {
		err = fmt.Errorf("redis read unkwnow type %d", *info.Type)
		log.Printf("[CRITAL]%s redis read type name error, reason:%s\n ", logId, err.Error())
		return err
	}
	d.Set("type", typeName)

	d.Set("mem_size", int64(*info.Size))

	d.Set("vpc_id", *info.UniqVpcId)
	d.Set("subnet_id", *info.UniqSubnetId)

	d.Set("project_id", *info.ProjectId)
	d.Set("port", *info.Port)
	d.Set("ip", *info.WanIp)
	d.Set("create_time", *info.Createtime)

	if d.Get("vpc_id").(string) != "" {
		securityGroups, err := service.DescribeInstanceSecurityGroup(ctx, d.Id())
		if err != nil {
			return err
		}
		if len(securityGroups) > 0 {
			d.Set("security_groups", securityGroups)
		}
	}

	tags := make(map[string]string, len(info.InstanceTags))
	for _, tag := range info.InstanceTags {
		if tag.TagKey == nil {
			return errors.New("redis tag key is nil")
		}
		if tag.TagValue == nil {
			return errors.New("redis tag value is nil")
		}

		tags[*tag.TagKey] = *tag.TagValue
	}
	d.Set("tags", tags)

	return nil
}

func resourceTencentCloudRedisInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_instance.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()

	d.Partial(true)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	// name\mem_size\password\project_id

	if d.HasChange("name") {
		name := d.Get("name").(string)
		if name == "" {
			name = id
		}
		err := service.ModifyInstanceName(ctx, id, name)
		if err != nil {
			return err
		}
		d.SetPartial("name")
	}

	if d.HasChange("mem_size") {

		oldInter, newInter := d.GetChange("mem_size")
		newMemSize := newInter.(int)
		oldMemSize := oldInter.(int)

		if oldMemSize >= newMemSize {
			return fmt.Errorf("redis mem_size can only increase")
		}

		if newMemSize < 1 {
			return fmt.Errorf("redis mem_size value cannot be set to less than 1")
		}
		redisId, err := service.UpgradeInstance(ctx, id, int64(newMemSize))

		if err != nil {
			log.Printf("[CRITAL]%s redis update mem size error, reason:%s\n ", logId, err.Error())
		}

		err = resource.Retry(600*time.Second, func() *resource.RetryError {
			_, _, info, err := service.CheckRedisCreateOk(ctx, redisId)

			if info != nil {
				status := REDIS_STATUS[*info.Status]
				if status == "" {
					return resource.NonRetryableError(fmt.Errorf("after update redis mem size, redis status is unknown ,status=%d", *info.Status))
				}
				if *info.Status == REDIS_STATUS_PROCESSING || *info.Status == REDIS_STATUS_INIT {
					return resource.RetryableError(fmt.Errorf("redis update processing."))
				}
				if *info.Status == REDIS_STATUS_ONLINE {
					return nil
				}
				return resource.NonRetryableError(fmt.Errorf("after update redis mem size, redis status is %s", status))
			}

			if err != nil {
				if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
					return resource.RetryableError(err)
				} else {
					return resource.NonRetryableError(err)
				}
			}
			return resource.NonRetryableError(fmt.Errorf("after update redis mem size, redis disappear"))
		})

		if err != nil {
			log.Printf("[CRITAL]%s redis update mem size fail , reason:%s\n ", logId, err.Error())
			return err
		}

		d.SetPartial("mem_size")
	}

	if d.HasChange("password") {
		password := d.Get("password").(string)
		taskid, err := service.ResetPassword(ctx, id, password)
		if err != nil {
			log.Printf("[CRITAL]%s redis change password error, reason:%s\n ", logId, err.Error())
			return err
		}
		err = resource.Retry(300*time.Second, func() *resource.RetryError {
			ok, err := service.DescribeTaskInfo(ctx, id, taskid)
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
			log.Printf("[CRITAL]%s redis change password fail, reason:%s\n ", logId, err.Error())
			return err
		}
		d.SetPartial("password")
	}

	if d.HasChange("project_id") {
		projectId := d.Get("project_id").(int)
		err := service.ModifyInstanceProjectId(ctx, id, int64(projectId))
		if err != nil {
			return err
		}
		d.SetPartial("project_id")
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}

		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := buildResourceNameWithUin("redis", "instance", id, region)

		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceTencentCloudRedisInstanceRead(d, meta)
}

func resourceTencentCloudRedisInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	defer logElapsed("resource.tencentcloud_redis_instance.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	var wait = func(action string, taskId int64) (errRet error) {

		errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			ok, err := service.DescribeTaskInfo(ctx, d.Id(), taskId)
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
				return resource.RetryableError(fmt.Errorf("%s timeout.", action))
			}
		})

		if errRet != nil {
			log.Printf("[CRITAL]%s redis %s fail, reason:%s\n ", logId, action, errRet.Error())
		}
		return errRet
	}

	action := "DestroyPostpaidInstance"
	taskId, err := service.DestroyPostpaidInstance(ctx, d.Id())
	if err != nil {
		log.Printf("[CRITAL]%s redis %s fail, reason:%s\n ", logId, action, err.Error())
		return err
	}
	if err = wait(action, taskId); err != nil {
		return err
	}

	action = "CleanUpInstance"
	taskId, err = service.CleanUpInstance(ctx, d.Id())
	if err != nil {
		log.Printf("[CRITAL]%s redis %s fail, reason:%s\n ", logId, action, err.Error())
		return err
	}

	return wait(action, taskId)
}
