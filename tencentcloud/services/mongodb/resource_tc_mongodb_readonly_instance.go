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
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudMongodbReadOnlyInstance() *schema.Resource {
	mongodbReadOnlyInstanceInfo := map[string]*schema.Schema{
		"father_instance_region": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
			Description: "Indicates the region of main instance.",
		},
		"father_instance_id": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
			Description: "Indicates the main instance ID of readonly instances.",
		},
		"cluster_type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			Description: "Instance schema type." +
				"	- REPLSET: Replset cluster;" +
				"	- SHARD: Shard cluster.",
		},
		"node_num": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "The number of nodes in each replica set. Default value: 3.",
		},
		"shard_quantity": {
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: tccommon.ValidateIntegerInRange(2, 20),
			Description:  "Number of sharding.",
		},
		"nodes_per_shard": {
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: tccommon.ValidateIntegerInRange(3, 5),
			Description:  "Number of nodes per shard, at least 3(one master and two slaves).",
		},
	}
	basic := TencentMongodbBasicInfo()
	conflictList := []string{"password"}
	for _, item := range conflictList {
		delete(basic, item)
	}
	for k, v := range basic {
		mongodbReadOnlyInstanceInfo[k] = v
	}

	return &schema.Resource{
		Create: resourceTencentCloudMongodbReadOnlyInstanceCreate,
		Read:   resourceTencentCloudMongodbReadOnlyInstanceRead,
		Update: resourceTencentCloudMongodbReadOnlyInstanceUpdate,
		Delete: resourceTencentCloudMongodbReadOnlyInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: mongodbReadOnlyInstanceInfo,
	}
}

func mongodbAllReadOnlyInstanceReqSet(requestInter interface{}, d *schema.ResourceData, masterInfo map[string]string) error {
	var (
		replicateSetNum       = 1
		nodeNum               = 3
		goodsNum              = 1
		clusterType           = MONGODB_CLUSTER_TYPE_REPLSET
		memoryInterface       = d.Get("memory").(int)
		volumeInterface       = d.Get("volume").(int)
		mongoVersionInterface = masterInfo["engine_version"]
		zoneInterface         = d.Get("available_zone").(string)
		machine               = masterInfo["machine_type"]
		fatherId              = masterInfo["father_instance_id"]
		instanceType          = MONGO_INSTANCE_TYPE_READONLY
		projectId             = d.Get("project_id").(int)
	)

	if v, ok := d.GetOk("shard_quantity"); ok {
		replicateSetNum = v.(int)
	}

	if v, ok := d.GetOk("nodes_per_shard"); ok {
		nodeNum = v.(int)
	}

	if v, ok := d.GetOk("cluster_type"); ok {
		clusterType = v.(string)
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
		"Clone":           helper.IntInt64(instanceType),
		"ProjectId":       helper.IntInt64(projectId),
		"Father":          &fatherId,
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
	if v, ok := d.GetOk("availability_zone_list"); ok {
		availabilityZoneList := helper.InterfacesStringsPoint(v.([]interface{}))
		value.FieldByName("AvailabilityZoneList").Set(reflect.ValueOf(availabilityZoneList))
	}
	if v, ok := d.GetOk("hidden_zone"); ok {
		value.FieldByName("HiddenZone").Set(reflect.ValueOf(helper.String(v.(string))))
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
	return nil
}

func mongodbCreateReadOnlyInstanceByUse(ctx context.Context, d *schema.ResourceData, meta interface{}, masterInfo map[string]string) error {
	logId := tccommon.GetLogId(ctx)
	request := mongodb.NewCreateDBInstanceHourRequest()

	if err := mongodbAllReadOnlyInstanceReqSet(request, d, masterInfo); err != nil {
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
		return fmt.Errorf("mongodb ReadOnly instance id is nil")
	}
	d.SetId(*response.Response.InstanceIds[0])

	return nil
}

func mongodbCreateReadOnlyInstanceByMonth(ctx context.Context, d *schema.ResourceData, meta interface{}, masterInfo map[string]string) error {
	logId := tccommon.GetLogId(ctx)
	request := mongodb.NewCreateDBInstanceRequest()

	if err := mongodbAllReadOnlyInstanceReqSet(request, d, masterInfo); err != nil {
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
		return fmt.Errorf("mongodb ReadOnly instance id is nil")
	}
	d.SetId(*response.Response.InstanceIds[0])

	return nil
}

func resourceTencentCloudMongodbReadOnlyInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_readonly_instance.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	client1 := *client
	mongodbService := MongodbService{client: client}
	mongodbService1 := MongodbService{client: &client1}
	tagService := svctag.NewTagService(client)
	region := client.Region

	// collect info from master instance
	var masterInfoMap = make(map[string]string, 3)
	if _, ok := d.GetOk("father_instance_id"); !ok {
		return fmt.Errorf("[CRITAL] father instance id must be specified for ReadOnly instance")
	}
	if _, ok := d.GetOk("father_instance_region"); !ok {
		return fmt.Errorf("[CRITAL] father instance region must be specified for ReadOnly instance")
	}
	fatherRegion := d.Get("father_instance_region").(string)
	mongodbService1.client.Region = fatherRegion
	masterInfoMap["father_instance_id"] = d.Get("father_instance_id").(string)
	masterInfo, has, err := mongodbService1.DescribeInstanceById(ctx, masterInfoMap["father_instance_id"])
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("[CRITAL] father instance can't be found for creating mongodb ReadOnly instance")
	}

	masterInfoMap["machine_type"] = *masterInfo.MachineType
	masterInfoMap["engine_version"] = *masterInfo.MongoVersion

	chargeType := d.Get("charge_type").(string)

	if chargeType == MONGODB_CHARGE_TYPE_POSTPAID {
		_, ok := d.GetOk("prepaid_period")
		_, ok1 := d.GetOk("auto_renew_flag")
		if ok || ok1 {
			return fmt.Errorf("prepaid_period and auto_renew_flag don't make sense for POSTPAID_BY_HOUR mongodb ReadOnly instance, please remove them from your template")
		}
		if err := mongodbCreateReadOnlyInstanceByUse(ctx, d, meta, masterInfoMap); err != nil {
			return err
		}
	} else {
		if err := mongodbCreateReadOnlyInstanceByMonth(ctx, d, meta, masterInfoMap); err != nil {
			return err
		}
	}

	instanceId := d.Id()

	_, has, err = mongodbService.DescribeInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("[CRITAL]%s creating mongodb ReadOnly instance failed, instance doesn't exist", logId)
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

	return resourceTencentCloudMongodbReadOnlyInstanceRead(d, meta)
}

func resourceTencentCloudMongodbReadOnlyInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_readonly_instance.read")()
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
	_ = d.Set("memory", *instance.Memory/1024)
	_ = d.Set("volume", *instance.Volume/1024)
	_ = d.Set("engine_version", instance.MongoVersion)
	_ = d.Set("charge_type", MONGODB_CHARGE_TYPE[*instance.PayMode])
	if MONGODB_CHARGE_TYPE[*instance.PayMode] == MONGODB_CHARGE_TYPE_PREPAID {
		_ = d.Set("auto_renew_flag", *instance.AutoRenewFlag)
	}

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
	_ = d.Set("machine_type", *instance.MachineType)
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

	// info of master info
	if instance.RelatedInstance != nil {
		_ = d.Set("father_instance_id", instance.RelatedInstance.InstanceId)
		_ = d.Set("father_instance_region", instance.RelatedInstance.Region)
	}

	if len(instance.Tags) > 0 {
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
	}

	return nil
}

func resourceTencentCloudMongodbReadOnlyInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_readonly_instance.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	mongodbService := MongodbService{client: client}
	tagService := svctag.NewTagService(client)
	region := client.Region

	if d.HasChange("engine_version") {
		return fmt.Errorf("setting of the field[engine_version] does not support update")
	}

	d.Partial(true)

	if d.HasChange("memory") || d.HasChange("volume") {
		memory := d.Get("memory").(int)
		volume := d.Get("volume").(int)
		params := make(map[string]interface{})
		var inMaintenance int
		if v, ok := d.GetOkExists("in_maintenance"); ok {
			inMaintenance = v.(int)
			params["in_maintenance"] = v.(int)
		}
		dealId, err := mongodbService.UpgradeInstance(ctx, instanceId, memory, volume, params)
		if err != nil {
			return err
		}
		if dealId == "" {
			return fmt.Errorf("deal id is empty")
		}

		if inMaintenance == 0 {
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

	return resourceTencentCloudMongodbReadOnlyInstanceRead(d, meta)
}

func resourceTencentCloudMongodbReadOnlyInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_readonly_instance.delete")()

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
		return fmt.Errorf("PREPAID instances are not allowed to be deleted now, please isolate them on console")
	}

	err = mongodbService.IsolateInstance(ctx, instanceId)
	if err != nil {
		return err
	}
	err = mongodbService.OfflineIsolatedDBInstance(ctx, instanceId, true)
	if err != nil {
		log.Printf("[CRITAL]%s mongodb %s fail, reason:%s", logId, "OfflineIsolatedDBInstance", err.Error())
		return err
	}
	//describe and check not exist
	_, has, errRet := mongodbService.DescribeInstanceById(ctx, instanceId)
	if errRet != nil {
		return errRet
	}
	if !has {
		return nil
	}
	return fmt.Errorf("[CRITAL]%s mongodb %s fail", logId, "OfflineIsolatedDBInstance")
}
