package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				ValidateFunc: validateIntegerMin(2),
				Description:  "Instance type. Available values reference data source `tencentcloud_redis_zone_config` or [document](https://intl.cloud.tencent.com/document/product/239/32069), toggle immediately when modified.",
			},
			"redis_shard_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedIntValue([]int{1, 3, 5, 8, 12, 16, 24, 32, 40, 48, 64, 80, 96, 128}),
				Description:  "The number of instance shards; this parameter does not need to be configured for standard version instances; for cluster version instances, the number of shards ranges from: [`1`, `3`, `5`, `8`, `12`, `16`, `24 `, `32`, `40`, `48`, `64`, `80`, `96`, `128`].",
			},
			"redis_replicas_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validateAllowedIntValue([]int{1, 2, 3, 4, 5}),
				Description:  "The number of instance copies. This is not required for standalone and master slave versions and must equal to count of `replica_zone_ids`, Non-multi-AZ does not require `replica_zone_ids`; Redis memory version 4.0, 5.0, 6.2 standard architecture and cluster architecture support the number of copies in the range [1, 2, 3, 4, 5]; Redis 2.8 standard version and CKV standard version only support 1 copy.",
			},
			"replica_zone_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "ID of replica nodes available zone. This is not required for standalone and master slave versions. NOTE: Removing some of the same zone of replicas (e.g. removing 100001 of [100001, 100001, 100002]) will pick the first hit to remove.",
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
				Description: "Indicates whether the redis instance support no-auth access. NOTE: Only available in private cloud environment.",
			},
			"replicas_read_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether copy read-only is supported, Redis 2.8 Standard Edition and CKV Standard Edition do not support replica read-only, turn on replica read-only, the instance will automatically read and write separate, write requests are routed to the primary node, read requests are routed to the replica node, if you need to open replica read-only, the recommended number of replicas >=2.",
			},
			"mem_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateAllowedIntValue([]int{1024, 2048, 4096, 8192, 12288, 16384, 20480, 24576, 32768, 40960, 49152, 65536}),
				Description:  "The memory volume of an available instance(in MB), please refer to `tencentcloud_redis_zone_config.list[zone].shard_memories`. When redis is standard type, it represents total memory size of the instance; when Redis is cluster type, it represents memory size of per sharding.",
			},
			"vpc_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateStringLengthInRange(1, 100),
				Description:  "ID of the vpc with which the instance is to be associated. When the `operation_network` is `changeVpc` or `changeBaseToVpc`, this parameter needs to be configured.",
			},
			"subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateStringLengthInRange(1, 100),
				Description:  "Specifies which subnet the instance should belong to. When the `operation_network` is `changeVpc` or `changeBaseToVpc`, this parameter needs to be configured.",
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set: func(v interface{}) int {
					return helper.HashString(v.(string))
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
				Default:     6379,
				Description: "The port used to access a redis instance. The default value is 6379. When the `operation_network` is `changeVPort` or `changeVip`, this parameter needs to be configured.",
			},
			"params_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify params template id. If not set, will use default template.",
			},

			"operation_network": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(REDIS_MODIFY_NETWORK_CONFIG),
				Description:  "Refers to the category of the pre-modified network, including: `changeVip`: refers to switching the private network, including its intranet IPv4 address and port; `changeVpc`: refers to switching the subnet to which the private network belongs; `changeBaseToVpc`: refers to switching the basic network to a private network; `changeVPort`: refers to only modifying the instance network port.",
			},

			"recycle": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue(REDIS_RECYCLE_TIME),
				Description:  "Original intranet IPv4 address retention time: unit: day, value range: `0`, `1`, `2`, `3`, `7`, `15`.",
			},

			"ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "IP address of an instance. When the `operation_network` is `changeVip`, this parameter needs to be configured.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Instance tags.",
			},

			// Computed values
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
			"auto_renew_flag": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1, 2}),
				Default:      0,
				Description:  "Auto-renew flag. 0 - default state (manual renewal); 1 - automatic renewal; 2 - explicit no automatic renewal.",
			},
			"node_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Readonly Primary/Replica nodes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"master": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the node is master.",
						},
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the master or replica node.",
						},
						"zone_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the availability zone of the master or replica node.",
						},
					},
				},
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
	redisShardNum := 1
	if v, ok := d.GetOk("redis_shard_num"); ok {
		redisShardNum = v.(int)
	}
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
	autoRenewFlag := d.Get("auto_renew_flag").(int)
	paramsTemplateId := d.Get("params_template_id").(string)
	operation := d.Get("operation_network").(string)
	chargeTypeID := REDIS_CHARGE_TYPE_ID[chargeType]
	var replicasReadonly bool
	if v, ok := d.GetOk("replicas_read_only"); ok {
		replicasReadonly = v.(bool)
	}
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

	if operation != "" {
		return fmt.Errorf("This parameter `operation_network` is not required when redis is created")
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
		autoRenewFlag,
		replicasReadonly,
		paramsTemplateId,
	)

	if err != nil {
		return err
	}

	if len(instanceIds) == 0 {
		return fmt.Errorf("redis api CreateInstances return empty redis id")
	}
	var redisId = *instanceIds[0]
	_, _, _, err = redisService.CheckRedisOnlineOk(ctx, redisId, 20*readRetryTimeout)

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
		has, _, info, e = service.CheckRedisOnlineOk(ctx, d.Id(), readRetryTimeout*20)
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
	_ = d.Set("auto_renew_flag", info.AutoRenewFlag)
	slaveReadWeight := *info.SlaveReadWeight
	if slaveReadWeight == 0 {
		_ = d.Set("replicas_read_only", false)
	} else if slaveReadWeight == 100 {
		_ = d.Set("replicas_read_only", true)
	}

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
		var zoneIds []int
		var nodeInfos []interface{}
		for i := range info.NodeSet {
			nodeInfo := info.NodeSet[i]
			nodeInfos = append(nodeInfos, map[string]interface{}{
				"master":  *nodeInfo.NodeType == 0,
				"zone_id": *nodeInfo.ZoneId,
				"id":      *nodeInfo.NodeId,
			})
			if *nodeInfo.NodeType == 0 {
				continue
			}
			zoneIds = append(zoneIds, int(*nodeInfo.ZoneId))
		}

		_ = d.Set("node_info", nodeInfos)

		var zoneIdsEqual = false

		replicaZones, replicaZonesOk := d.GetOk("replica_zone_ids")
		if replicaZonesOk {
			oldIds := helper.InterfacesIntegers(replicaZones.([]interface{}))
			zoneIdsEqual = checkIdsEqual(oldIds, zoneIds)
		}

		if !zoneIdsEqual {
			_ = d.Set("replica_zone_ids", zoneIds)
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
	}

	// MemSize, ShardNum and ReplicaNum can only change one for each upgrade invoke
	if d.HasChange("mem_size") {

		_, newInter := d.GetChange("mem_size")
		newMemSize := newInter.(int)
		oShard, _ := d.GetChange("redis_shard_num")
		redisShardNum := oShard.(int)
		oReplica, _ := d.GetChange("redis_replicas_num")
		redisReplicasNum := oReplica.(int)

		if newMemSize < 1 {
			return fmt.Errorf("redis mem_size value cannot be set to less than 1")
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			_, err := redisService.UpgradeInstance(ctx, id, newMemSize, redisShardNum, redisReplicasNum, nil)
			if err != nil {
				// Upgrade memory will cause instance lock and cannot acknowledge by polling status, wait until lock release
				return retryError(err, redis.FAILEDOPERATION_UNKNOWN, redis.FAILEDOPERATION_SYSTEMERROR)
			}
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s redis upgrade instance error, reason:%s\n", logId, err.Error())
			return err
		}

		err = redisService.CheckRedisUpdateOk(ctx, id)

		if err != nil {
			log.Printf("[CRITAL]%s redis update mem size fail , reason:%s\n", logId, err.Error())
			return err
		}
	}

	// MemSize, ShardNum and ReplicaNum can only change one for each upgrade invoke
	if d.HasChange("redis_shard_num") {
		redisShardNum := d.Get("redis_shard_num").(int)
		oReplica, _ := d.GetChange("redis_replicas_num")
		redisReplicasNum := oReplica.(int)
		memSize := d.Get("mem_size").(int)
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			_, err := redisService.UpgradeInstance(ctx, id, memSize, redisShardNum, redisReplicasNum, nil)
			if err != nil {
				// Upgrade memory will cause instance lock and cannot acknowledge by polling status, wait until lock release
				return retryError(err, redis.FAILEDOPERATION_UNKNOWN, redis.FAILEDOPERATION_SYSTEMERROR)
			}
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s redis upgrade instance error, reason:%s\n", logId, err.Error())
			return err
		}

		err = redisService.CheckRedisUpdateOk(ctx, id)

		if err != nil {
			log.Printf("[CRITAL]%s redis update shard num fail , reason:%s\n", logId, err.Error())
			return err
		}
	}

	if d.HasChange("redis_replicas_num") || d.HasChange("replica_zone_ids") {
		err := resourceRedisNodeSetModify(ctx, &redisService, d)
		if err != nil {
			return err
		}
	}

	if d.HasChange("password") || d.HasChange("no_auth") {
		var (
			taskId   int64
			password = d.Get("password").(string)
			noAuth   = d.Get("no_auth").(bool)
			err      error
		)

		// After redis spec modified, reset password may not successfully response immediately.
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			taskId, err = redisService.ResetPassword(ctx, id, password, noAuth)
			if err != nil {
				log.Printf("[CRITAL]%s redis change password error, reason:%s\n", logId, err.Error())
				return retryError(err, redis.FAILEDOPERATION_SYSTEMERROR)
			}
			return nil
		})

		if err != nil {
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
	}

	if d.HasChange("params_template_id") {
		request := redis.NewApplyParamsTemplateRequest()
		request.InstanceIds = []*string{&id}
		request.TemplateId = helper.String(d.Get("params_template_id").(string))
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			_, err := redisService.ApplyParamsTemplate(ctx, request)
			if err != nil {
				return retryError(err, redis.FAILEDOPERATION_SYSTEMERROR, redis.RESOURCEUNAVAILABLE_INSTANCELOCKEDERROR)
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	if d.HasChange("project_id") {
		projectId := d.Get("project_id").(int)
		err := redisService.ModifyInstanceProjectId(ctx, id, int64(projectId))
		if err != nil {
			return err
		}
	}

	if d.HasChanges("security_groups") {
		sgs := d.Get("security_groups").(*schema.Set).List()
		var sgIds []*string
		for _, sgId := range sgs {
			sgIds = append(sgIds, helper.String(sgId.(string)))
		}
		err := redisService.ModifyDBInstanceSecurityGroups(ctx, "redis", d.Id(), sgIds)
		if err != nil {
			return err
		}
	}

	if d.HasChanges("type_id") {
		request := redis.NewUpgradeInstanceVersionRequest()
		typeId := d.Get("type_id").(int)
		request.InstanceId = &id
		request.TargetInstanceType = helper.String(strconv.Itoa(typeId))
		request.SwitchOption = helper.IntInt64(2)
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().UpgradeInstanceVersion(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate redis upgradeVersionOperation failed, reason:%+v", logId, err)
			return err
		}

		service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
		_, _, _, err = service.CheckRedisOnlineOk(ctx, id, 20*readRetryTimeout)
		if err != nil {
			log.Printf("[CRITAL]%s redis upgradeVersionOperation fail, reason:%s\n", logId, err.Error())
			return err
		}
	}

	if d.HasChange("vpc_id") || d.HasChange("subnet_id") || d.HasChange("port") || d.HasChange("recycle") || d.HasChange("ip") {
		if _, ok := d.GetOk("operation_network"); !ok {
			return fmt.Errorf("When modifying `vpc_id`, `subnet_id`, `port`, `recycle`, `ip`, the `operation_network` parameter is required")
		}

		request := redis.NewModifyNetworkConfigRequest()
		request.InstanceId = &id

		operation := d.Get("operation_network").(string)
		request.Operation = &operation

		switch operation {
		case REDIS_MODIFY_NETWORK_CONFIG[0]:
			if v, ok := d.GetOk("ip"); ok {
				request.Vip = helper.String(v.(string))
			} else {
				return fmt.Errorf("When `operation_network` is %v, this parameter must be filled in", operation)
			}

			if v, ok := d.GetOk("port"); ok {
				request.VPort = helper.IntInt64(v.(int))
			} else {
				return fmt.Errorf("When `operation_network` is %v, this parameter must be filled in", operation)
			}
		case REDIS_MODIFY_NETWORK_CONFIG[1], REDIS_MODIFY_NETWORK_CONFIG[2]:
			if v, ok := d.GetOk("vpc_id"); ok {
				request.VpcId = helper.String(v.(string))
			} else {
				return fmt.Errorf("When `operation_network` is %v, this parameter must be filled in", operation)
			}

			if v, ok := d.GetOk("subnet_id"); ok {
				request.SubnetId = helper.String(v.(string))
			} else {
				return fmt.Errorf("When `operation_network` is %v, this parameter must be filled in", operation)
			}
		case REDIS_MODIFY_NETWORK_CONFIG[3]:
			if v, ok := d.GetOk("port"); ok {
				request.VPort = helper.IntInt64(v.(int))
			} else {
				return fmt.Errorf("When `operation_network` is %v, this parameter must be filled in", operation)
			}
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().ModifyNetworkConfig(request)
			if e != nil {
				if _, ok := e.(*sdkErrors.TencentCloudSDKError); !ok {
					return resource.RetryableError(e)
				} else {
					return resource.NonRetryableError(e)
				}
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate redis networkConfig failed, reason:%+v", logId, err)
			return err
		}

		service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
		_, _, _, err = service.CheckRedisOnlineOk(ctx, id, 20*readRetryTimeout)
		if err != nil {
			log.Printf("[CRITAL]%s redis networkConfig fail, reason:%s\n", logId, err.Error())
			return err
		}

		_ = d.Set("operation_network", operation)
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		resourceName := BuildTagResourceName("redis", "instance", region, id)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
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
	has, _, info, err := service.CheckRedisOnlineOk(ctx, d.Id(), 20*readRetryTimeout)

	if err != nil {
		log.Printf("[CRITAL]%s redis querying before deleting task fail, reason:%s\n", logId, err.Error())
		return err
	}

	if !has {
		return nil
	}

	chargeType = REDIS_CHARGE_TYPE_NAME[*info.BillingMode]

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
			}
			return resource.RetryableError(fmt.Errorf("%s timeout.", action))
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
	}
	return nil
}

func checkIdsEqual(o []int, n []int) bool {
	if len(o) != len(n) {
		return false
	}

	sort.Ints(o)
	sort.Ints(n)

	for i, v := range o {
		if v != n[i] {
			return false
		}
	}
	return true
}

func resourceRedisNodeSetModify(ctx context.Context, service *RedisService, d *schema.ResourceData) error {
	id := d.Id()
	memSize := d.Get("mem_size").(int)
	shardNum := d.Get("redis_shard_num").(int)
	o, n := d.GetChange("replica_zone_ids")
	oz := helper.InterfacesIntegers(o.([]interface{}))
	nz := helper.InterfacesIntegers(n.([]interface{}))
	log.Printf("o = %v, n = %v", oz, nz)
	adds, lacks := GetListDiffs(oz, nz)

	var redisNodeInfos []*redis.RedisNodeInfo

	if len(adds) > 0 {
		_, _, info, err := service.CheckRedisOnlineOk(ctx, id, readRetryTimeout)
		if err != nil {
			return err
		}
		redisNodeInfos = info.NodeSet
		redisReplicaCount := len(redisNodeInfos) - 1

		log.Printf("%v will be add", adds)
		var addNodes []*redis.RedisNodeInfo
		for _, zoneId := range adds {
			addNodes = append(addNodes, &redis.RedisNodeInfo{
				NodeType: helper.IntInt64(1),
				ZoneId:   helper.IntUint64(zoneId),
			})
		}
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			_, err := service.UpgradeInstance(ctx, d.Id(), memSize, shardNum, redisReplicaCount+len(adds), addNodes)
			if err != nil {
				return retryError(err, redis.FAILEDOPERATION_UNKNOWN)
			}
			return nil
		})
		if err != nil {
			return err
		}
		err = service.CheckRedisUpdateOk(ctx, id)
		if err != nil {
			return err
		}
	}

	if len(lacks) > 0 {
		_, _, info, err := service.CheckRedisOnlineOk(ctx, id, readRetryTimeout)
		if err != nil {
			return err
		}
		redisNodeInfos = info.NodeSet
		redisReplicaCount := len(redisNodeInfos) - 1
		removeNodes := tencentCloudRedisGetRemoveNodesByIds(lacks[:], redisNodeInfos)
		replicasParam := redisReplicaCount - len(lacks)
		if replicasParam <= 0 {
			return fmt.Errorf("cannot delete replica %d which is your only replica on instance %s", removeNodes[0].NodeId, id)
		}
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			_, err := service.UpgradeInstance(ctx, id, memSize, shardNum, replicasParam, removeNodes)
			if err != nil {
				return retryError(err, redis.FAILEDOPERATION_UNKNOWN)
			}
			return nil
		})
		if err != nil {
			return err
		}
		err = service.CheckRedisUpdateOk(ctx, id)
		if err != nil {
			return err
		}
	}

	// Non-Multi-AZ modification redis_replicas_num
	if d.HasChange("redis_replicas_num") && len(oz) == 0 && len(nz) == 0 {
		_, replica := d.GetChange("redis_replicas_num")
		redisReplicasNum := replica.(int)
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			_, err := service.UpgradeInstance(ctx, id, memSize, shardNum, redisReplicasNum, nil)
			if err != nil {
				// Upgrade memory will cause instance lock and cannot acknowledge by polling status, wait until lock release
				return retryError(err, redis.FAILEDOPERATION_UNKNOWN, redis.FAILEDOPERATION_SYSTEMERROR)
			}
			return nil
		})
		if err != nil {
			return err
		}

		err = service.CheckRedisUpdateOk(ctx, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func tencentCloudRedisGetRemoveNodesByIds(ids []int, nodes []*redis.RedisNodeInfo) (result []*redis.RedisNodeInfo) {
	for i := range nodes {
		node := nodes[i]
		if *node.NodeType == 0 {
			continue
		}
		index := FindIntListIndex(ids, int(*node.ZoneId))
		if index == -1 {
			continue
		}
		result = append(result, node)
		ids = append(ids[:index], ids[index+1:]...)
	}
	return
}
