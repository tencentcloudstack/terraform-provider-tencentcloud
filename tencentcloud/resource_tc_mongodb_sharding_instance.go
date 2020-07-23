/*
Provide a resource to create a Mongodb sharding instance.

Example Usage

```hcl
resource "tencentcloud_mongodb_sharding_instance" "mongodb" {
  instance_name   = "mongodb"
  shard_quantity  = 2
  nodes_per_shard = 3
  memory          = 4
  volume          = 100
  engine_version  = "MONGO_3_WT"
  machine_type    = "GIO"
  available_zone  = "ap-guangzhou-3"
  vpc_id          = "vpc-mz3efvbw"
  subnet_id       = "subnet-lk0svi3p"
  project_id      = 0
  password        = "password1234"
}
```

Import

Mongodb sharding instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_mongodb_sharding_instance.mongodb cmgo-41s6jwy4
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMongodbShardingInstance() *schema.Resource {
	mongodbShardingInstanceInfo := map[string]*schema.Schema{
		"shard_quantity": {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateIntegerInRange(2, 20),
			Description:  "Number of sharding.",
		},
		"nodes_per_shard": {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateIntegerInRange(3, 5),
			Description:  "Number of nodes per shard, at least 3(one master and two slaves).",
		},
	}
	basic := TencentMongodbBasicInfo()
	for k, v := range basic {
		mongodbShardingInstanceInfo[k] = v
	}

	return &schema.Resource{
		Create: resourceMongodbShardingInstanceCreate,
		Read:   resourceMongodbShardingInstanceRead,
		Update: resourceMongodbShardingInstanceUpdate,
		Delete: resourceMongodbShardingInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: mongodbShardingInstanceInfo,
	}
}

func mongodbAllShardingInstanceReqSet(requestInter interface{}, d *schema.ResourceData) error {
	requestByMonth, okByMonth := requestInter.(*mongodb.CreateDBInstanceRequest)
	requestByUse, _ := requestInter.(*mongodb.CreateDBInstanceHourRequest)

	var (
		replicateSetNum       = d.Get("shard_quantity").(int)
		nodeNum               = d.Get("nodes_per_shard").(int)
		goodsNum              = 1
		clusterType           = MONGODB_CLUSTER_TYPE_SHARD
		memoryInterface       = d.Get("memory").(int)
		volumeInterface       = d.Get("volume").(int)
		mongoVersionInterface = d.Get("engine_version").(string)
		zoneInterface         = d.Get("available_zone").(string)
		machine               = d.Get("machine_type").(string)
		password              = d.Get("password").(string)
		instanceType          = MONGO_INSTANCE_TYPE_FORMAL
	)

	if machine == MONGODB_MACHINE_TYPE_GIO {
		machine = MONGODB_MACHINE_TYPE_HIO
	} else if machine == MONGODB_MACHINE_TYPE_TGIO {
		machine = MONGODB_MACHINE_TYPE_HIO10G
	}

	if okByMonth {
		if v, ok := d.GetOk("prepaid_period"); ok {
			requestByMonth.Period = helper.IntUint64(v.(int))
		} else {
			return fmt.Errorf("prepaid_period must be specified for a PREPAID instance")
		}
		requestByMonth.AutoRenewFlag = helper.IntUint64(d.Get("auto_renew_flag").(int))

		requestByMonth.ReplicateSetNum = helper.IntUint64(replicateSetNum)
		requestByMonth.NodeNum = helper.IntUint64(nodeNum)
		requestByMonth.GoodsNum = helper.IntUint64(goodsNum)
		requestByMonth.ClusterType = &clusterType
		requestByMonth.Memory = helper.IntUint64(memoryInterface)
		requestByMonth.Volume = helper.IntUint64(volumeInterface)
		requestByMonth.MongoVersion = &mongoVersionInterface
		requestByMonth.Zone = &zoneInterface
		requestByMonth.MachineCode = &machine
		requestByMonth.Password = &password
		requestByMonth.Clone = helper.IntInt64(instanceType)

		if v, ok := d.GetOk("vpc_id"); ok {
			requestByMonth.VpcId = helper.String(v.(string))
		}
		if v, ok := d.GetOk("subnet_id"); ok {
			requestByMonth.SubnetId = helper.String(v.(string))
		}
		err := fmt.Errorf("you have to set vpc_id and subnet_id both")
		if (requestByMonth.VpcId != nil && requestByMonth.SubnetId == nil) || (requestByMonth.VpcId == nil && requestByMonth.SubnetId != nil) {
			return err
		}
		if v, ok := d.GetOk("project_id"); ok {
			requestByMonth.ProjectId = helper.IntInt64(v.(int))
		}
		if v, ok := d.GetOk("security_groups"); ok {
			securityGroups := v.(*schema.Set).List()
			requestByMonth.SecurityGroup = make([]*string, 0, len(securityGroups))
			for _, v := range securityGroups {
				requestByMonth.SecurityGroup = append(requestByMonth.SecurityGroup, helper.String(v.(string)))
			}
		}
	} else {
		requestByUse.ReplicateSetNum = helper.IntUint64(replicateSetNum)
		requestByUse.NodeNum = helper.IntUint64(nodeNum)
		requestByUse.GoodsNum = helper.IntUint64(goodsNum)
		requestByUse.ClusterType = &clusterType
		requestByUse.Memory = helper.IntUint64(memoryInterface)
		requestByUse.Volume = helper.IntUint64(volumeInterface)
		requestByUse.MongoVersion = &mongoVersionInterface
		requestByUse.Zone = &zoneInterface
		requestByUse.MachineCode = &machine
		requestByUse.Password = &password
		requestByUse.Clone = helper.IntInt64(instanceType)

		if v, ok := d.GetOk("vpc_id"); ok {
			requestByUse.VpcId = helper.String(v.(string))
		}
		if v, ok := d.GetOk("subnet_id"); ok {
			requestByUse.SubnetId = helper.String(v.(string))
		}
		err := fmt.Errorf("you have to set vpc_id and subnet_id both")
		if (requestByUse.VpcId != nil && requestByUse.SubnetId == nil) || (requestByUse.VpcId == nil && requestByUse.SubnetId != nil) {
			return err
		}
		if v, ok := d.GetOk("project_id"); ok {
			requestByUse.ProjectId = helper.IntInt64(v.(int))
		}
		if v, ok := d.GetOk("security_groups"); ok {
			securityGroups := v.(*schema.Set).List()
			requestByUse.SecurityGroup = make([]*string, 0, len(securityGroups))
			for _, v := range securityGroups {
				requestByUse.SecurityGroup = append(requestByUse.SecurityGroup, helper.String(v.(string)))
			}
		}
	}

	return nil
}

func mongodbCreateShardingInstanceByUse(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	logId := getLogId(ctx)
	request := mongodb.NewCreateDBInstanceHourRequest()

	if err := mongodbAllShardingInstanceReqSet(request, d); err != nil {
		return err
	}

	var response *mongodb.CreateDBInstanceHourResponse
	err := resource.Retry(6*writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMongodbClient().CreateDBInstanceHour(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s\n", logId, request.GetAction(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.InstanceIds) < 1 {
		return fmt.Errorf("mongodb instance id is nil")
	}
	d.SetId(*response.Response.InstanceIds[0])

	return nil
}

func mongodbCreateShardingInstanceByMonth(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	logId := getLogId(ctx)
	request := mongodb.NewCreateDBInstanceRequest()

	if err := mongodbAllShardingInstanceReqSet(request, d); err != nil {
		return err
	}

	var response *mongodb.CreateDBInstanceResponse
	err := resource.Retry(6*writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMongodbClient().CreateDBInstance(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s\n", logId, request.GetAction(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.InstanceIds) < 1 {
		return fmt.Errorf("mongodb instance id is nil")
	}
	d.SetId(*response.Response.InstanceIds[0])

	return nil
}

func resourceMongodbShardingInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mongodb_sharding_instance.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	client := meta.(*TencentCloudClient).apiV3Conn
	mongodbService := MongodbService{client: client}
	tagService := TagService{client: client}
	region := client.Region

	// check security group info
	if d.Get("engine_version").(string) == MONGODB_ENGINE_VERSION_4_WT {
		if _, ok := d.GetOk("security_groups"); ok {
			return fmt.Errorf("[CRITAL] for instance which `engine_version` is `MONGO_40_WT`, `security_groups` is not supported")
		}
	}

	chargeType := d.Get("charge_type")
	if chargeType == MONGODB_CHARGE_TYPE_POSTPAID {
		_, ok := d.GetOk("prepaid_period")
		_, ok1 := d.GetOk("auto_renew_flag")
		if ok || ok1 {
			return fmt.Errorf("prepaid_period and auto_renew_flag don't make sense for POSTPAID instance, please remove them from your template")
		}
		if err := mongodbCreateShardingInstanceByUse(ctx, d, meta); err != nil {
			return err
		}
	} else {
		if err := mongodbCreateShardingInstanceByMonth(ctx, d, meta); err != nil {
			return err
		}
	}

	instanceId := d.Id()

	_, has, err := mongodbService.DescribeInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("[CRITAL]%s creating mongodb sharding instance failed, instance doesn't exist\n", logId)
	}

	// setting instance name
	instanceName := d.Get("instance_name").(string)
	err = mongodbService.ModifyInstanceName(ctx, instanceId, instanceName)
	if err != nil {
		return err
	}

	_, has, err = mongodbService.DescribeInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("[CRITAL]%s creating mongodb instance failed, instance doesn't exist\n", logId)
	}
	d.SetId(instanceId)

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		resourceName := BuildTagResourceName("mongodb", "instance", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceMongodbShardingInstanceRead(d, meta)
}

func resourceMongodbShardingInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mongodb_sharding_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Id()

	mongodbService := MongodbService{client: meta.(*TencentCloudClient).apiV3Conn}
	instance, has, err := mongodbService.DescribeInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	if nilFields := CheckNil(instance, map[string]string{
		"InstanceName":      "instance name",
		"ProjectId":         "project id",
		"ClusterType":       "cluster type",
		"Zone":              "available zone",
		"VpcId":             "vpc id",
		"SubnetId":          "subnet id",
		"Status":            "status",
		"Vip":               "vip",
		"Vport":             "vport",
		"CreateTime":        "create time",
		"MongoVersion":      "engine version",
		"Memory":            "memory",
		"Volume":            "volume",
		"MachineType":       "machine type",
		"ReplicationSetNum": "shard quantity",
		"SecondaryNum":      "secondary number",
	}); len(nilFields) > 0 {
		return fmt.Errorf("mongodb %v are nil", nilFields)
	}

	_ = d.Set("shard_quantity", instance.ReplicationSetNum)
	_ = d.Set("nodes_per_shard", *instance.SecondaryNum+1)
	_ = d.Set("instance_name", instance.InstanceName)
	_ = d.Set("memory", *instance.Memory/1024/(*instance.ReplicationSetNum))
	_ = d.Set("volume", *instance.Volume/1024/(*instance.ReplicationSetNum))
	_ = d.Set("engine_version", instance.MongoVersion)
	_ = d.Set("charge_type", MONGODB_CHARGE_TYPE[*instance.PayMode])
	if MONGODB_CHARGE_TYPE[*instance.PayMode] == MONGODB_CHARGE_TYPE_PREPAID {
		_ = d.Set("auto_renew_flag", *instance.AutoRenewFlag)
	}

	switch *instance.MachineType {
	case MONGODB_MACHINE_TYPE_TGIO:
		_ = d.Set("machine_type", MONGODB_MACHINE_TYPE_HIO10G)

	case MONGODB_MACHINE_TYPE_GIO:
		_ = d.Set("machine_type", MONGODB_MACHINE_TYPE_HIO)

	default:
		_ = d.Set("machine_type", *instance.MachineType)
	}

	_ = d.Set("available_zone", instance.Zone)
	_ = d.Set("vpc_id", instance.VpcId)
	_ = d.Set("subnet_id", instance.SubnetId)
	_ = d.Set("project_id", instance.ProjectId)
	_ = d.Set("status", instance.Status)
	_ = d.Set("vip", instance.Vip)
	_ = d.Set("vport", instance.Vport)
	_ = d.Set("create_time", instance.CreateTime)

	tags := make(map[string]string, len(instance.Tags))
	for _, tag := range instance.Tags {
		if tag.TagKey == nil {
			return errors.New("mongodb tag key is nil")
		}
		if tag.TagValue == nil {
			return errors.New("mongodb tag value is nil")
		}

		tags[*tag.TagKey] = *tag.TagValue
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceMongodbShardingInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mongodb_sharding_instance.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Id()

	client := meta.(*TencentCloudClient).apiV3Conn
	mongodbService := MongodbService{client: client}
	tagService := TagService{client: client}
	region := client.Region

	d.Partial(true)

	if d.HasChange("memory") || d.HasChange("volume") {
		// precheck
		if !(d.HasChange("memory") && d.HasChange("volume")) {
			return fmt.Errorf("[CRITAL] updating memory and volume of mongodb instance failed, memory and volume must upgrade/downgrade at same time")
		}
		oldMemory, newMemory := d.GetChange("memory")
		oldVolume, newVolume := d.GetChange("volume")
		if (newMemory.(int)-oldMemory.(int))^(newVolume.(int)-oldVolume.(int)) < 0 {
			return fmt.Errorf("[CRITAL] updating memory and volume of mongodb instance failed, memory and volume must upgrade/downgrade at same time")
		}
		memory := d.Get("memory").(int)
		volume := d.Get("volume").(int)
		err := mongodbService.UpgradeInstance(ctx, instanceId, memory, volume)
		if err != nil {
			return err
		}

		// it will take time to wait for memory and volume change even describe request succeeded even the status returned in describe response is running
		errUpdate := resource.Retry(20*readRetryTimeout, func() *resource.RetryError {
			infos, has, e := mongodbService.DescribeInstanceById(ctx, instanceId)
			if e != nil {
				return resource.NonRetryableError(e)
			}
			if !has {
				return resource.NonRetryableError(fmt.Errorf("[CRITAL]%s updating mongodb sharding instance failed, instance doesn't exist\n", logId))
			}

			memoryDes := *infos.Memory / 1024 / (*infos.ReplicationSetNum)
			volumeDes := *infos.Volume / 1024 / (*infos.ReplicationSetNum)
			if d.Get("memory").(int) != int(memoryDes) || d.Get("volume").(int) != int(volumeDes) {
				return resource.RetryableError(fmt.Errorf("[CRITAL] updating mongodb sharding instance, current memory and volume values: %d, %d, waiting for them becoming new value: %d, %d", memoryDes, volumeDes, d.Get("memory").(int), d.Get("volume").(int)))
			}
			return nil
		})
		if errUpdate != nil {
			return fmt.Errorf("[CRITAL] updating mongodb sharding instance failed, memory and volume values don't change")
		}

		d.SetPartial("memory")
		d.SetPartial("volume")
	}

	if d.HasChange("instance_name") {
		instanceName := d.Get("instance_name").(string)
		err := mongodbService.ModifyInstanceName(ctx, instanceId, instanceName)
		if err != nil {
			return err
		}
		d.SetPartial("instance_name")
	}

	if d.HasChange("project_id") {
		projectId := d.Get("project_id").(int)
		err := mongodbService.ModifyProjectId(ctx, instanceId, projectId)
		if err != nil {
			return err
		}
		d.SetPartial("project_id")
	}

	if d.HasChange("password") {
		password := d.Get("password").(string)
		err := mongodbService.ResetInstancePassword(ctx, instanceId, "mongouser", password)
		if err != nil {
			return err
		}

		d.SetPartial("password")
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		resourceName := BuildTagResourceName("mongodb", "instance", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

		d.SetPartial("tags")
	}

	if d.HasChange("prepaid_period") {
		return fmt.Errorf("setting of the field[prepaid_period] does not make sense after the initialization")
	}

	if d.HasChange("auto_renew_flag") {
		autoRenewFlag := d.Get("auto_renew_flag").(int)
		period := d.Get("prepaid_period").(int)
		err := mongodbService.ModifyAutoRenewFlag(ctx, instanceId, period, autoRenewFlag)
		if err != nil {
			return err
		}
		d.SetPartial("auto_renew_flag")
	}

	d.Partial(false)

	return resourceMongodbShardingInstanceRead(d, meta)
}

func resourceMongodbShardingInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mongodb_sharding_instance.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Id()
	mongodbService := MongodbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	instanceDetail, has, err := mongodbService.DescribeInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}
	if MONGODB_CHARGE_TYPE[*instanceDetail.PayMode] == MONGODB_CHARGE_TYPE_PREPAID {
		return fmt.Errorf("PREPAID instances are not allowed to be deleted now, please isolate them on console")
	}

	err = mongodbService.IsolateInstance(ctx, instanceId)
	if err != nil {
		return err
	}

	return nil
}
