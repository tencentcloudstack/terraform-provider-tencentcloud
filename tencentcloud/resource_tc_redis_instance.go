/*
Provides a resource to create a Redis instance and set its attributes.

Example Usage

```hcl
data "tencentcloud_redis_zone_config" "zone" {
}

resource "tencentcloud_redis_instance" "redis_instance_test_2" {
  availability_zone  = data.tencentcloud_redis_zone_config.zone.list[0].zone
  type_id            = data.tencentcloud_redis_zone_config.zone.list[0].type_id
  password           = "test12345789"
  mem_size           = 8192
  redis_shard_num    = data.tencentcloud_redis_zone_config.zone.list[0].redis_shard_nums[0]
  redis_replicas_num = data.tencentcloud_redis_zone_config.zone.list[0].redis_replicas_nums[0]
  name               = "terrform_test"
  port               = 6379
}
```

Using multi replica zone set
```
data "tencentcloud_availability_zones" "az" {

}

variable "redis_replicas_num" {
  default = 3
}

resource "tencentcloud_redis_instance" "red1" {
  availability_zone  = data.tencentcloud_availability_zones.az.zones[0].name
  charge_type        = "POSTPAID"
  mem_size           = 1024
  name               = "test-redis"
  port               = 6379
  project_id         = 0
  redis_replicas_num = var.redis_replicas_num
  redis_shard_num    = 1
  security_groups    = [
    "sg-d765yoec",
  ]
  subnet_id          = "subnet-ie01x91v"
  type_id            = 6
  vpc_id             = "vpc-k4lrsafc"
  password = "a12121312334"

  replica_zone_ids = [
    for i in range(var.redis_replicas_num)
    : data.tencentcloud_availability_zones.az.zones[i % length(data.tencentcloud_availability_zones.az.zones)].id ]
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
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudRedisInstance() *schema.Resource {
	types := []string{}
	for _, v := range REDIS_NAMES {
		types = append(types, "`"+v+"`")
	}
	sort.Strings(types)
	typeStr := strings.Trim(strings.Join(types, ","), ",")

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
				Description: "The available zone ID of an instance to be created, please refer to `tencentcloud_redis_zone_config.list`.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Instance name.",
			},
			"type_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateIntegerMin(2),
				Description:  "Instance type. Available values reference data source `tencentcloud_redis_zone_config` or [document](https://intl.cloud.tencent.com/document/product/239/32069).",
			},
			"redis_shard_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     1,
				Description: "The number of instance shard. This is not required for standalone and master slave versions.",
			},
			"redis_replicas_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     1,
				Description: "The number of instance copies. This is not required for standalone and master slave versions.",
			},
			"replica_zone_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of replica nodes available zone. This is not required for standalone and master slave versions.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					for _, name := range REDIS_NAMES {
						if name == value {
							return
						}
					}
					errors = append(errors, fmt.Errorf("this redis type %s not support now", value))
					return
				},
				Deprecated:  "It has been deprecated from version 1.33.1. Please use 'type_id' instead.",
				Description: "Instance type. Available values: " + typeStr + ", specific region support specific types, need to refer data `tencentcloud_redis_zone_config`.",
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validateMysqlPassword,
				Description:  "Password for a Redis user, which should be 8 to 16 characters. NOTE: Only `no_auth=true` specified can make password empty.",
			},
			"no_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Indicates whether the redis instance support no-auth access. NOTE: Only available in private cloud environment.",
			},
			"mem_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The memory volume of an available instance(in MB), please refer to `tencentcloud_redis_zone_config.list[zone].mem_sizes`. When redis is standard type, it represents total memory size of the instance; when Redis is cluster type, it represents memory size of per sharding.",
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
				Description: "ID of security group. If both vpc_id and subnet_id are not set, this argument should not be set either.",
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
				Description: "Current status of an instance, maybe: init, processing, online, isolate and todelete.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the instance was created.",
			},
			// payment
			"charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      REDIS_CHARGE_TYPE_POSTPAID,
				ValidateFunc: validateAllowedStringValue([]string{REDIS_CHARGE_TYPE_POSTPAID, REDIS_CHARGE_TYPE_PREPAID}),
				Description:  "The charge type of instance. Valid values: `PREPAID` and `POSTPAID`. Default value is `POSTPAID`. Note: TencentCloud International only supports `POSTPAID`. Caution that update operation on this field will delete old instances and create new with new charge type.",
			},
			"prepaid_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue(REDIS_PREPAID_PERIOD),
				Description:  "The tenancy (time unit is month) of the prepaid instance, NOTE: it only works when charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.",
			},
			"force_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicate whether to delete Redis instance directly or not. Default is false. If set true, the instance will be deleted instead of staying recycle bin. Note: only works for `PREPAID` instance.",
			},
		},
	}
}

func resourceTencentCloudRedisInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_instance.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	client := meta.(*TencentCloudClient).apiV3Conn
	redisService := RedisService{client: client}
	tagService := TagService{client: client}
	region := client.Region

	availabilityZone := d.Get("availability_zone").(string)
	redisName := d.Get("name").(string)
	redisType := d.Get("type").(string)
	typeId := int64(d.Get("type_id").(int))
	redisShardNum := d.Get("redis_shard_num").(int)
	redisReplicasNum := d.Get("redis_replicas_num").(int)
	password := d.Get("password").(string)
	noAuth := d.Get("no_auth").(bool)
	memSize := d.Get("mem_size").(int)
	vpcId := d.Get("vpc_id").(string)
	subnetId := d.Get("subnet_id").(string)
	securityGroups := d.Get("security_groups").(*schema.Set).List()
	projectId := d.Get("project_id").(int)
	port := d.Get("port").(int)
	chargeType := d.Get("charge_type").(string)
	chargeTypeID := REDIS_CHARGE_TYPE_ID[chargeType]
	var chargePeriod uint64 = 1
	if chargeType == REDIS_CHARGE_TYPE_PREPAID {
		if period, ok := d.GetOk("prepaid_period"); ok {
			chargePeriod = uint64(period.(int))
		} else {
			return fmt.Errorf("instance charge type prepaid period can not be empty when charge type is %s", chargeType)
		}
	}

	if (typeId == 0 && redisType == "") || (typeId != 0 && redisType != "") {
		return fmt.Errorf("`type_id` and `type` set one item and only one item")
	}

	if password == "" && !noAuth {
		return fmt.Errorf("`password` must not be empty unless `no_auth` is `true`")
	}

	if noAuth && (vpcId == "" || subnetId == "") {
		return fmt.Errorf("cannot set `no_auth=true` if `vpc_id` and `subnet_id` is empty")
	}

	for id, name := range REDIS_NAMES {
		if redisType == name {
			typeId = id
			break
		}
	}

	sellConfigures, err := redisService.DescribeRedisZoneConfig(ctx)
	if err != nil {
		return fmt.Errorf("api[DescribeRedisZoneConfig]fail, return %s", err.Error())
	}
	var regionItem *redis.RegionConf
	var zoneItem *redis.ZoneCapacityConf
	var redisItem *redis.ProductConf
	for _, regionItem = range sellConfigures {
		if *regionItem.RegionId == region {
			break
		}
	}
	if regionItem == nil {
		return fmt.Errorf("all redis in this region `%s` be sold out", region)
	}
	for _, zones := range regionItem.ZoneSet {
		if *zones.IsSaleout {
			continue
		}
		if *zones.ZoneName == availabilityZone {
			zoneItem = zones
			break
		}
	}
	if zoneItem == nil {
		return fmt.Errorf("all redis in this zone `%s` be sold out", availabilityZone)
	}

	for _, reds := range zoneItem.ProductSet {
		if *reds.Type == typeId {
			redisItem = reds
			break
		}
	}
	if redisItem == nil {
		return fmt.Errorf("redis type_id `%d` be sold out or this type_id is not supports", typeId)
	}
	var redisShardNums []string
	var redisReplicasNums []string
	var numErrors []string
	for _, v := range redisItem.ShardNum {
		redisShardNums = append(redisShardNums, *v)
	}
	for _, v := range redisItem.ReplicaNum {
		redisReplicasNums = append(redisReplicasNums, *v)
	}
	if !IsContains(redisShardNums, fmt.Sprintf("%d", redisShardNum)) {
		numErrors = append(numErrors, fmt.Sprintf("redis_shard_num : %s", strings.Join(redisShardNums, ",")))
	}

	if !IsContains(redisReplicasNums, fmt.Sprintf("%d", redisReplicasNum)) {
		numErrors = append(numErrors, fmt.Sprintf(" redis_replicas_num : %s", strings.Join(redisReplicasNums, ",")))
	}

	if len(numErrors) > 0 {
		return fmt.Errorf("redis type_id `%d` only supports %s", typeId, strings.Join(numErrors, ","))
	}

	requestSecurityGroup := make([]string, 0, len(securityGroups))

	for _, v := range securityGroups {
		requestSecurityGroup = append(requestSecurityGroup, v.(string))
	}

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	nodeInfo := make([]*redis.RedisNodeInfo, 0)
	if raw, ok := d.GetOk("replica_zone_ids"); ok {
		zoneIds := raw.([]interface{})

		masterZoneId, err := service.getZoneId(availabilityZone)
		if err != nil {
			return err
		}

		// insert master node
		nodeInfo = append(nodeInfo, &redis.RedisNodeInfo{
			NodeType: helper.Int64(0),
			ZoneId:   helper.Int64Uint64(masterZoneId),
		})

		for _, v := range zoneIds {
			id := v.(int)
			nodeInfo = append(nodeInfo, &redis.RedisNodeInfo{
				NodeType: helper.Int64(1),
				ZoneId:   helper.IntUint64(id),
			})
		}
	}

	instanceIds, err := redisService.CreateInstances(ctx,
		availabilityZone,
		typeId,
		password,
		vpcId,
		subnetId,
		redisName,
		int64(memSize),
		int64(projectId),
		int64(port),
		requestSecurityGroup,
		redisShardNum,
		redisReplicasNum,
		chargeTypeID,
		chargePeriod,
		nodeInfo,
		noAuth,
	)

	if err != nil {
		return err
	}

	if len(instanceIds) == 0 {
		return fmt.Errorf("redis api CreateInstances return empty redis id")
	}
	var redisId = *instanceIds[0]
	err = resource.Retry(20*readRetryTimeout, func() *resource.RetryError {
		has, online, _, err := redisService.CheckRedisOnlineOk(ctx, redisId)
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
		log.Printf("[CRITAL]%s create redis task fail, reason:%s\n", logId, err.Error())
		return err
	}
	d.SetId(redisId)

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		resourceName := BuildTagResourceName("redis", "instance", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudRedisInstanceRead(d, meta)
}

func resourceTencentCloudRedisInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
	var onlineHas = true
	var (
		has  bool
		info *redis.InstanceSet
		e    error
	)
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		has, _, info, e = service.CheckRedisOnlineOk(ctx, d.Id())
		if info != nil {
			if *info.Status == REDIS_STATUS_ISOLATE || *info.Status == REDIS_STATUS_TODELETE {
				d.SetId("")
				onlineHas = false
				return nil
			}
		}
		if e != nil {
			return resource.NonRetryableError(e)
		}
		if !has {
			d.SetId("")
			onlineHas = false
			return nil
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Fail to get info from redis, reaseon %s", err.Error())
	}
	if !onlineHas {
		return nil
	}

	statusName := REDIS_STATUS[*info.Status]
	if statusName == "" {
		err = fmt.Errorf("redis read unkwnow status %d", *info.Status)
		log.Printf("[CRITAL]%s redis read status name error, reason:%s\n", logId, err.Error())
		return err
	}
	_ = d.Set("status", statusName)

	_ = d.Set("name", *info.InstanceName)

	zoneName, err := service.getZoneName(*info.ZoneId)
	if err != nil {
		return err
	}
	// not set field type_id
	// process import case
	if d.Get("type_id").(int) == 0 && d.Get("type").(string) != "" {
		typeName := REDIS_NAMES[*info.Type]
		if typeName == "" {
			err = fmt.Errorf("redis read unkwnow type %d", *info.Type)
			log.Printf("[CRITAL]%s redis read type name error, reason:%s\n", logId, err.Error())
			return err
		}
		_ = d.Set("type", typeName)
	} else {
		_ = d.Set("type_id", info.Type)
	}

	_ = d.Set("redis_shard_num", info.RedisShardNum)
	_ = d.Set("redis_replicas_num", info.RedisReplicasNum)
	_ = d.Set("availability_zone", zoneName)
	_ = d.Set("mem_size", info.RedisShardSize)
	_ = d.Set("vpc_id", info.UniqVpcId)
	_ = d.Set("subnet_id", info.UniqSubnetId)
	_ = d.Set("project_id", info.ProjectId)
	_ = d.Set("port", info.Port)
	_ = d.Set("ip", info.WanIp)
	_ = d.Set("create_time", info.Createtime)

	// only true or user explicit declared will set for import case.
	if _, ok := d.GetOk("no_auth"); ok || *info.NoAuth {
		_ = d.Set("no_auth", info.NoAuth)
	}

	if d.Get("vpc_id").(string) != "" {
		securityGroups, err := service.DescribeInstanceSecurityGroup(ctx, d.Id())
		if err != nil {
			return err
		}
		if len(securityGroups) > 0 {
			_ = d.Set("security_groups", securityGroups)
		}
	}

	if info.NodeSet != nil {
		var zoneIds []uint64
		for i := range info.NodeSet {
			nodeInfo := info.NodeSet[i]
			if *nodeInfo.NodeType == 0 {
				continue
			}
			zoneIds = append(zoneIds, *nodeInfo.ZoneId)
		}

		if err := d.Set("replica_zone_ids", zoneIds); err != nil {
			log.Printf("[WARN] replica_zone_ids set error: %s", err.Error())
		}
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "redis", "instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	_ = d.Set("charge_type", REDIS_CHARGE_TYPE_NAME[*info.BillingMode])
	return nil
}

func resourceTencentCloudRedisInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_instance.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	client := meta.(*TencentCloudClient).apiV3Conn
	redisService := RedisService{client: client}
	tagService := TagService{client: client}
	region := client.Region

	d.Partial(true)

	unsupportedUpdateFields := []string{
		"prepaid_period",
	}
	for _, field := range unsupportedUpdateFields {
		if d.HasChange(field) {
			return fmt.Errorf("tencentcloud_redis_instance update on %s is not support yet", field)
		}
	}

	// name\mem_size\password\project_id

	if d.HasChange("name") {
		name := d.Get("name").(string)
		if name == "" {
			name = id
		}
		err := redisService.ModifyInstanceName(ctx, id, name)
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
		redisShardNum := d.Get("redis_shard_num").(int)
		redisReplicasNum := d.Get("redis_replicas_num").(int)
		_, err := redisService.UpgradeInstance(ctx, id, int64(newMemSize), int64(redisShardNum), int64(redisReplicasNum))

		if err != nil {
			log.Printf("[CRITAL]%s redis update mem size error, reason:%s\n", logId, err.Error())
			return err
		}

		err = resource.Retry(4*readRetryTimeout, func() *resource.RetryError {
			_, _, info, err := redisService.CheckRedisOnlineOk(ctx, id)

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
			log.Printf("[CRITAL]%s redis update mem size fail , reason:%s\n", logId, err.Error())
			return err
		}

		d.SetPartial("mem_size")
	}

	if d.HasChange("password") {
		password := d.Get("password").(string)
		taskId, err := redisService.ResetPassword(ctx, id, password)
		if err != nil {
			log.Printf("[CRITAL]%s redis change password error, reason:%s\n", logId, err.Error())
			return err
		}
		err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
			ok, err := redisService.DescribeTaskInfo(ctx, id, taskId)
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
		d.SetPartial("password")
	}

	if d.HasChange("project_id") {
		projectId := d.Get("project_id").(int)
		err := redisService.ModifyInstanceProjectId(ctx, id, int64(projectId))
		if err != nil {
			return err
		}
		d.SetPartial("project_id")
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		resourceName := BuildTagResourceName("redis", "instance", region, id)
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
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	// Collect infos before deleting action
	var chargeType string
	errQuery := resource.Retry(20*readRetryTimeout, func() *resource.RetryError {
		has, online, info, err := service.CheckRedisOnlineOk(ctx, d.Id())
		if err != nil {
			log.Printf("[CRITAL]%s redis querying before deleting fail, reason:%s\n", logId, err.Error())
			return resource.NonRetryableError(err)
		}
		if !has {
			return nil
		}
		if online {
			chargeType = REDIS_CHARGE_TYPE_NAME[*info.BillingMode]
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("Deleting ERROR: Creating redis task is processing."))
		}
	})
	if errQuery != nil {
		log.Printf("[CRITAL]%s redis querying before deleting task fail, reason:%s\n", logId, errQuery.Error())
		return errQuery
	}

	var wait = func(action string, taskInfo interface{}) (errRet error) {

		errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			var ok bool
			var err error
			switch v := taskInfo.(type) {
			case int64:
				ok, err = service.DescribeTaskInfo(ctx, d.Id(), v)
			case string:
				ok, _, err = service.DescribeInstanceDealDetail(ctx, v)
			}
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
			log.Printf("[CRITAL]%s redis %s fail, reason:%s\n", logId, action, errRet.Error())
		}
		return errRet
	}

	forceDelete := d.Get("force_delete").(bool)
	if chargeType == REDIS_CHARGE_TYPE_POSTPAID {
		forceDelete = true
		taskId, err := service.DestroyPostpaidInstance(ctx, d.Id())
		if err != nil {
			log.Printf("[CRITAL]%s redis %s fail, reason:%s\n", logId, "DestroyPostpaidInstance", err.Error())
			return err
		}
		if err = wait("DestroyPostpaidInstance", taskId); err != nil {
			return err
		}

	} else {
		if _, err := service.DestroyPrepaidInstance(ctx, d.Id()); err != nil {
			log.Printf("[CRITAL]%s redis %s fail, reason:%s\n", logId, "DestroyPrepaidInstance", err.Error())
			return err
		}

		// Deal info only support create and renew and resize, need to check destroy status by describing api.
		if errDestroyChecking := resource.Retry(20*readRetryTimeout, func() *resource.RetryError {
			has, isolated, err := service.CheckRedisDestroyOk(ctx, d.Id())
			if err != nil {
				log.Printf("[CRITAL]%s CheckRedisDestroyOk fail, reason:%s\n", logId, err.Error())
				return resource.NonRetryableError(err)
			}
			if !has || isolated {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("instance is not ready to be destroyed"))
		}); errDestroyChecking != nil {
			log.Printf("[CRITAL]%s redis querying before deleting task fail, reason:%s\n", logId, errDestroyChecking.Error())
			return errDestroyChecking
		}
	}

	if forceDelete {
		taskId, err := service.CleanUpInstance(ctx, d.Id())
		if err != nil {
			log.Printf("[CRITAL]%s redis %s fail, reason:%s\n", logId, "CleanUpInstance", err.Error())
			return err
		}

		return wait("CleanUpInstance", taskId)
	} else {
		return nil
	}
}
