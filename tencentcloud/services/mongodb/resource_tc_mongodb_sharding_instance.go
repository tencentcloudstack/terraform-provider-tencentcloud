package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudMongodbShardingInstance() *schema.Resource {
	mongodbShardingInstanceInfo := map[string]*schema.Schema{
		"shard_quantity": {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: tccommon.ValidateIntegerInRange(2, 20),
			Description:  "Number of sharding.",
		},
		"nodes_per_shard": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "Number of nodes per shard, at least 3(one master and two slaves). Allow value[3, 5, 7].",
		},
		"availability_zone_list": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: `A list of nodes deployed in multiple availability zones. For more information, please use the API DescribeSpecInfo.
			- Multi-availability zone deployment nodes can only be deployed in 3 different availability zones. It is not supported to deploy most nodes of the cluster in the same availability zone. For example, a 3-node cluster does not support the deployment of 2 nodes in the same zone.
			- Version 4.2 and above are not supported.
			- Read-only disaster recovery instances are not supported.
			- Basic network cannot be selected.`,
		},
		"hidden_zone": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The availability zone to which the Hidden node belongs. This parameter is required in cross-AZ instance deployment.",
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
		projectId             = d.Get("project_id").(int)
	)

	if v, ok := d.GetOk("password"); ok && v.(string) != "" {
		password = v.(string)
	} else {
		return fmt.Errorf("`password` cannot be empty when creating")
	}

	if machine == MONGODB_MACHINE_TYPE_GIO {
		machine = MONGODB_MACHINE_TYPE_HIO
	} else if machine == MONGODB_MACHINE_TYPE_TGIO {
		machine = MONGODB_MACHINE_TYPE_HIO10G
	}

	getType := reflect.TypeOf(requestInter)
	value := reflect.ValueOf(requestInter).Elem()

	for k, v := range map[string]interface{}{
		"ReplicateSetNum": helper.IntUint64(replicateSetNum),
		"NodeNum":         helper.IntUint64(nodeNum),
		"GoodsNum":        helper.IntUint64(goodsNum),
		"ClusterType":     &clusterType,
		"Memory":          helper.IntUint64(memoryInterface),
		"Volume":          helper.IntUint64(volumeInterface),
		"MongoVersion":    &mongoVersionInterface,
		"Zone":            &zoneInterface,
		"MachineCode":     &machine,
		"Password":        &password,
		"Clone":           helper.IntInt64(instanceType),
		"ProjectId":       helper.IntInt64(projectId),
	} {
		value.FieldByName(k).Set(reflect.ValueOf(v))
	}

	var okVpc, okSubnet bool
	if v, ok := d.GetOk("vpc_id"); ok {
		okVpc = ok
		value.FieldByName("VpcId").Set(reflect.ValueOf(helper.String(v.(string))))
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		okSubnet = ok
		value.FieldByName("SubnetId").Set(reflect.ValueOf(helper.String(v.(string))))
	}
	if (okVpc && !okSubnet) || (!okVpc && okSubnet) {
		return fmt.Errorf("you have to set vpc_id and subnet_id both")
	}
	if v, ok := d.GetOk("security_groups"); ok {
		sliceReflect := helper.InterfacesStringsPoint(v.(*schema.Set).List())
		value.FieldByName("SecurityGroup").Set(reflect.ValueOf(sliceReflect))
	}

	if strings.Contains(getType.String(), "CreateDBInstanceRequest") {
		if v, ok := d.GetOk("prepaid_period"); ok {
			value.FieldByName("Period").Set(reflect.ValueOf(helper.IntUint64(v.(int))))
		} else {
			return fmt.Errorf("prepaid_period must be specified for a PREPAID instance")
		}
		value.FieldByName("AutoRenewFlag").Set(reflect.ValueOf(helper.IntUint64(d.Get("auto_renew_flag").(int))))
	}

	if v, ok := d.GetOk("mongos_memory"); ok {
		value.FieldByName("MongosMemory").Set(reflect.ValueOf(helper.IntUint64(v.(int))))
	}
	if v, ok := d.GetOk("mongos_cpu"); ok {
		value.FieldByName("MongosCpu").Set(reflect.ValueOf(helper.IntUint64(v.(int))))
	}
	if v, ok := d.GetOk("mongos_node_num"); ok {
		value.FieldByName("MongosNodeNum").Set(reflect.ValueOf(helper.IntUint64(v.(int))))
	}
	if v, ok := d.GetOk("availability_zone_list"); ok {
		availabilityZoneList := helper.InterfacesStringsPoint(v.([]interface{}))
		value.FieldByName("AvailabilityZoneList").Set(reflect.ValueOf(availabilityZoneList))
	}
	if v, ok := d.GetOk("hidden_zone"); ok {
		value.FieldByName("HiddenZone").Set(reflect.ValueOf(helper.String(v.(string))))
	}
	return nil
}

func mongodbCreateShardingInstanceByUse(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	logId := tccommon.GetLogId(ctx)
	request := mongodb.NewCreateDBInstanceHourRequest()

	if err := mongodbAllShardingInstanceReqSet(request, d); err != nil {
		return err
	}

	var response *mongodb.CreateDBInstanceHourResponse
	var err error
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().CreateDBInstanceHour(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(response.Response.InstanceIds) < 1 {
		return fmt.Errorf("mongodb instance id is nil")
	}
	d.SetId(*response.Response.InstanceIds[0])

	return nil
}

func mongodbCreateShardingInstanceByMonth(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	logId := tccommon.GetLogId(ctx)
	request := mongodb.NewCreateDBInstanceRequest()

	if err := mongodbAllShardingInstanceReqSet(request, d); err != nil {
		return err
	}

	var response *mongodb.CreateDBInstanceResponse
	var err error
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().CreateDBInstance(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(response.Response.InstanceIds) < 1 {
		return fmt.Errorf("mongodb instance id is nil")
	}
	d.SetId(*response.Response.InstanceIds[0])

	return nil
}

func resourceMongodbShardingInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_sharding_instance.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	mongodbService := MongodbService{client: client}
	tagService := svctag.NewTagService(client)
	region := client.Region

	chargeType := d.Get("charge_type")
	if chargeType == MONGODB_CHARGE_TYPE_POSTPAID {
		_, ok := d.GetOk("prepaid_period")
		_, ok1 := d.GetOk("auto_renew_flag")
		if ok || ok1 {
			return fmt.Errorf("prepaid_period and auto_renew_flag don't make sense for POSTPAID_BY_HOUR instance, please remove them from your template")
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
		return fmt.Errorf("[CRITAL]%s creating mongodb sharding instance failed, instance doesn't exist", logId)
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
		return fmt.Errorf("[CRITAL]%s creating mongodb instance failed, instance doesn't exist", logId)
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		resourceName := tccommon.BuildTagResourceName("mongodb", "instance", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceMongodbShardingInstanceRead(d, meta)
}

func resourceMongodbShardingInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_sharding_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()

	mongodbService := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	instance, has, err := mongodbService.DescribeInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	if nilFields := tccommon.CheckNil(instance, map[string]string{
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
	_ = d.Set("mongos_cpu", instance.MongosCpuNum)
	_ = d.Set("mongos_memory", *instance.MongosMemory/1024)
	_ = d.Set("mongos_node_num", instance.MongosNodeNum)
	_ = d.Set("auto_renew_flag", instance.AutoRenewFlag)

	groups, err := mongodbService.DescribeSecurityGroup(ctx, instanceId)
	if err != nil {
		return err
	}
	groupIds := make([]string, 0)
	for _, group := range groups {
		groupIds = append(groupIds, *group.SecurityGroupId)
	}
	if len(groupIds) > 1 {
		_ = d.Set("security_groups", groupIds)
	}
	replicateSets, err := mongodbService.DescribeDBInstanceNodeProperty(ctx, instanceId)
	if err != nil {
		return err
	}
	if len(replicateSets) > 1 {
		var hiddenZone string
		availabilityZoneList := make([]string, 0, 3)
		for _, replicate := range replicateSets[0].Nodes {
			itemZone := *replicate.Zone
			if *replicate.Hidden {
				hiddenZone = itemZone
			}
			availabilityZoneList = append(availabilityZoneList, itemZone)
		}
		_ = d.Set("hidden_zone", hiddenZone)
		_ = d.Set("availability_zone_list", availabilityZoneList)
	}
	tags := make(map[string]string, len(instance.Tags))
	for _, tag := range instance.Tags {
		if tag.TagKey == nil {
			return errors.New("mongodb tag key is nil")
		}
		if tag.TagValue == nil {
			return errors.New("mongodb tag value is nil")
		}
		if *tag.TagKey == "project" {
			continue
		}

		tags[*tag.TagKey] = *tag.TagValue
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceMongodbShardingInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_sharding_instance.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	mongodbService := MongodbService{client: client}
	tagService := svctag.NewTagService(client)
	region := client.Region

	d.Partial(true)
	if d.HasChange("availability_zone_list") || d.HasChange("hidden_zone") {
		return fmt.Errorf("setting of the field[availability_zone_list, hidden_zone] does not support update")
	}
	if d.HasChange("mongos_node_num") {
		return fmt.Errorf("setting of the field[mongos_node_num] does not support update")
	}
	if d.HasChange("mongos_cpu") && d.HasChange("mongos_memory") {
		if v, ok := d.GetOk("mongos_memory"); ok {
			dealId, err := mongodbService.ModifyMongosMemory(ctx, instanceId, v.(int))
			if err != nil {
				return err
			}
			if dealId == "" {
				return fmt.Errorf("deal id is empty")
			}

			errUpdate := resource.Retry(20*tccommon.ReadRetryTimeout, func() *resource.RetryError {
				dealResponseParams, err := mongodbService.DescribeDBInstanceDeal(ctx, dealId)
				if err != nil {
					if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
						if sdkError.Code == "InvalidParameter" && sdkError.Message == "deal resource not found." {
							return resource.RetryableError(err)
						}
					}
					return resource.NonRetryableError(err)
				}

				if *dealResponseParams.Status != MONGODB_STATUS_DELIVERY_SUCCESS {
					return resource.RetryableError(fmt.Errorf("mongodb status is not delivery success"))
				}
				return nil
			})
			if errUpdate != nil {
				return errUpdate
			}
		}
	}
	if d.HasChange("memory") || d.HasChange("volume") {
		memory := d.Get("memory").(int)
		volume := d.Get("volume").(int)
		params := make(map[string]interface{})

		var inMaintenance int
		if v, ok := d.GetOkExists("in_maintenance"); ok {
			inMaintenance = v.(int)
			params["in_maintenance"] = v.(int)
		}
		_, err := mongodbService.UpgradeInstance(ctx, instanceId, memory, volume, params)
		if err != nil {
			return err
		}

		// it will take time to wait for memory and volume change even describe request succeeded even the status returned in describe response is running
		if inMaintenance == 0 {
			errUpdate := resource.Retry(20*tccommon.ReadRetryTimeout, func() *resource.RetryError {
				infos, has, e := mongodbService.DescribeInstanceById(ctx, instanceId)
				if e != nil {
					return resource.NonRetryableError(e)
				}
				if !has {
					return resource.NonRetryableError(fmt.Errorf("[CRITAL]%s updating mongodb sharding instance failed, instance doesn't exist", logId))
				}

				memoryDes := *infos.Memory / 1024 / (*infos.ReplicationSetNum)
				volumeDes := *infos.Volume / 1024 / (*infos.ReplicationSetNum)
				if memory != int(memoryDes) || volume != int(volumeDes) {
					return resource.RetryableError(fmt.Errorf("[CRITAL] updating mongodb sharding instance, current memory and volume values: %d, %d, waiting for them becoming new value: %d, %d", memoryDes, volumeDes, d.Get("memory").(int), d.Get("volume").(int)))
				}
				return nil
			})
			if errUpdate != nil {
				return errUpdate
			}
		}
	}

	if d.HasChange("instance_name") {
		instanceName := d.Get("instance_name").(string)
		err := mongodbService.ModifyInstanceName(ctx, instanceId, instanceName)
		if err != nil {
			return err
		}

	}

	if d.HasChange("project_id") {
		projectId := d.Get("project_id").(int)
		err := mongodbService.ModifyProjectId(ctx, instanceId, projectId)
		if err != nil {
			return err
		}

	}

	if d.HasChange("password") {
		password := d.Get("password").(string)
		err := mongodbService.ResetInstancePassword(ctx, instanceId, "mongouser", password)
		if err != nil {
			return err
		}

	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		resourceName := tccommon.BuildTagResourceName("mongodb", "instance", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

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

	}

	if d.HasChange("security_groups") {
		securityGroups := d.Get("security_groups").(*schema.Set).List()
		securityGroupIds := make([]*string, 0, len(securityGroups))
		for _, securityGroup := range securityGroups {
			securityGroupIds = append(securityGroupIds, helper.String(securityGroup.(string)))
		}
		err := mongodbService.ModifySecurityGroups(ctx, instanceId, securityGroupIds)
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	if d.HasChange("engine_version") {
		request := mongodb.NewUpgradeDbInstanceVersionRequest()
		response := mongodb.NewUpgradeDbInstanceVersionResponse()
		request.InstanceId = &instanceId
		if v, ok := d.GetOk("engine_version"); ok {
			request.MongoVersion = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().UpgradeDbInstanceVersionWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.FlowId == nil {
				return resource.NonRetryableError(fmt.Errorf("Upgrade engine version failed, Response is nil."))
			}

			response = result
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s upgrade engine version failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		flowIdStr := helper.UInt64ToStr(*response.Response.FlowId)
		if err := mongodbService.DescribeAsyncRequestInfo(ctx, flowIdStr, 20*tccommon.ReadRetryTimeout); err != nil {
			return err
		}
	}

	return resourceMongodbShardingInstanceRead(d, meta)
}

func resourceMongodbShardingInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_sharding_instance.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()
	mongodbService := MongodbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
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
		err := mongodbService.TerminateDBInstances(ctx, instanceId)
		if err != nil {
			return err
		}
	} else {
		err := mongodbService.IsolateInstance(ctx, instanceId)
		if err != nil {
			return err
		}
	}

	err = mongodbService.OfflineIsolatedDBInstance(ctx, instanceId, true)
	if err != nil {
		log.Printf("[CRITAL]%s mongodb %s fail, reason:%s", logId, "OfflineIsolatedDBInstance", err.Error())
		return err
	}
	return nil
}
