package mongodb

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

//internal version: replace import begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
//internal version: replace import end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

func ResourceTencentCloudMongodbInstance() *schema.Resource {
	mongodbInstanceInfo := map[string]*schema.Schema{
		"standby_instance_list": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of standby instances' info.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"standby_instance_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Indicates the ID of standby instance.",
					},
					"standby_instance_region": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Indicates the region of standby instance.",
					},
				},
			},
		},
		"node_num": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "The number of nodes in each replica set. Default value: 3.",
		},
		"add_node_list": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Add node attribute list.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"role": {
						Type:     schema.TypeString,
						Required: true,
						Description: "The node role that needs to be added.\n" +
							"- SECONDARY: Mongod node;\n" +
							"- READONLY: read-only node;\n" +
							"- MONGOS: Mongos node.",
					},
					"zone": {
						Type:     schema.TypeString,
						Required: true,
						Description: "The availability zone corresponding to the node.\n" +
							"- single availability zone, where all nodes are in the same availability zone;\n" +
							"- multiple availability zones: the current standard specification is the distribution of three availability zones, and the master and slave nodes are not in the same availability zone. You should pay attention to configuring the availability zone corresponding to the new node, and the rule that the number of nodes in any two availability zones is greater than the third availability zone must be met after the addition.",
					},
				},
			},
		},
		"remove_node_list": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Add node attribute list.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"role": {
						Type:     schema.TypeString,
						Required: true,
						Description: "The node role that needs to be deleted.\n" +
							"- SECONDARY: Mongod node;\n" +
							"- READONLY: read-only node;\n" +
							"- MONGOS: Mongos node.",
					},
					"node_name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The node ID to delete. The shard cluster must specify the name of the node to be deleted by a group of shards, and the rest of the shards should be grouped and aligned.",
					},
					"zone": {
						Type:     schema.TypeString,
						Required: true,
						Description: "The availability zone corresponding to the node.\n" +
							"- single availability zone, where all nodes are in the same availability zone;\n" +
							"- multiple availability zones: the current standard specification is the distribution of three availability zones, and the master and slave nodes are not in the same availability zone. You should pay attention to configuring the availability zone corresponding to the new node, and the rule that the number of nodes in any two availability zones is greater than the third availability zone must be met after the addition.",
					},
				},
			},
		},
		"availability_zone_list": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "If cloud database instances are deployed in multiple availability zones, specify a list of multiple availability zones.\n" +
				"	- To deploy an instance with multiple availability zones, the parameter Zone specifies the primary availability zone information of the instance; Availability ZoneList specifies all availability zone information, including the primary availability zone. The input format is as follows: [ap-Guangzhou-2,ap-Guangzhou-3,ap-Guangzhou-4].\n" +
				"	- You can obtain availability zone information planned in different regions of the cloud database through the interface DescribeSpecInfo, so as to specify effective availability zones.\n" +
				"	- Multiple availability zone deployment nodes can only be deployed in 3 different availability zones. Deploying most nodes of a cluster in the same availability zone is not supported. For example, a 3-node cluster does not support 2 nodes deployed in the same zone.",
		},
		"hidden_zone": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The availability zone to which the Hidden node belongs. This parameter is required in cross-AZ instance deployment.",
		},
		"maintenance_start": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Maintenance window start time. The value range is any full point or half point from `00:00-23:00`, such as 00:00 or 00:30.",
		},
		"maintenance_end": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			Description: "Maintenance window end time.\n" +
				"	- The value range is any full point or half point from `00:00-23:00`, and the maintenance time duration is at least 30 minutes and at most 3 hours.\n" +
				"	- The end time must be based on the start time backwards.",
		},
	}
	basic := TencentMongodbBasicInfo()
	conflictList := []string{"mongos_cpu", "mongos_memory", "mongos_node_num"}
	for _, item := range conflictList {
		delete(basic, item)
	}
	for k, v := range basic {
		mongodbInstanceInfo[k] = v
	}

	return &schema.Resource{
		Create: resourceTencentCloudMongodbInstanceCreate,
		Read:   resourceTencentCloudMongodbInstanceRead,
		Update: resourceTencentCloudMongodbInstanceUpdate,
		Delete: resourceTencentCloudMongodbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: mongodbInstanceInfo,
	}
}

func mongodbAllInstanceReqSet(requestInter interface{}, d *schema.ResourceData) error {
	var (
		replicateSetNum       = 1
		nodeNum               = 3
		goodsNum              = 1
		clusterType           = MONGODB_CLUSTER_TYPE_REPLSET
		memoryInterface       = d.Get("memory").(int)
		volumeInterface       = d.Get("volume").(int)
		mongoVersionInterface = d.Get("engine_version").(string)
		zoneInterface         = d.Get("available_zone").(string)
		machine               = d.Get("machine_type").(string)
		instanceType          = MONGO_INSTANCE_TYPE_FORMAL
		projectId             = d.Get("project_id").(int)
		password              string
	)

	if machine == MONGODB_MACHINE_TYPE_GIO {
		machine = MONGODB_MACHINE_TYPE_HIO
	} else if machine == MONGODB_MACHINE_TYPE_TGIO {
		machine = MONGODB_MACHINE_TYPE_HIO10G
	}

	if v, ok := d.GetOk("node_num"); ok {
		nodeNum = v.(int)
	}

	if v, ok := d.GetOk("password"); ok {
		password = v.(string)
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
	if v, ok := d.GetOk("availability_zone_list"); ok {
		availabilityZoneList := helper.InterfacesStringsPoint(v.([]interface{}))
		value.FieldByName("AvailabilityZoneList").Set(reflect.ValueOf(availabilityZoneList))
	}
	if v, ok := d.GetOk("hidden_zone"); ok {
		value.FieldByName("HiddenZone").Set(reflect.ValueOf(helper.String(v.(string))))
	}
	return nil
}

func mongodbCreateInstanceByUse(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	logId := tccommon.GetLogId(ctx)
	request := mongodb.NewCreateDBInstanceHourRequest()

	if err := mongodbAllInstanceReqSet(request, d); err != nil {
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

func mongodbCreateInstanceByMonth(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	logId := tccommon.GetLogId(ctx)
	request := mongodb.NewCreateDBInstanceRequest()

	if err := mongodbAllInstanceReqSet(request, d); err != nil {
		return err
	}

	var response *mongodb.CreateDBInstanceResponse
	var err error
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().CreateDBInstance(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			//internal version: replace bpass begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
			//internal version: replace bpass end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
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

func resourceTencentCloudMongodbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance.create")()

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
		if err := mongodbCreateInstanceByUse(ctx, d, meta); err != nil {
			return err
		}
	} else {
		if err := mongodbCreateInstanceByMonth(ctx, d, meta); err != nil {
			return err
		}
	}

	instanceId := d.Id()

	//internal version: replace setTag begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	//internal version: replace setTag end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	_, has, err := mongodbService.DescribeInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("[CRITAL]%s creating mongodb instance failed, instance doesn't exist", logId)
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

	maintenanceStart, okMaintenanceStart := d.GetOk("maintenance_start")
	maintenanceEnd, okMaintenanceEnd := d.GetOk("maintenance_end")
	if okMaintenanceStart && okMaintenanceEnd {
		err := mongodbService.SetInstanceMaintenance(ctx, instanceId, maintenanceStart.(string), maintenanceEnd.(string))
		if err != nil {
			return err
		}
	}
	// mongodbService.SetInstanceMaintenance()
	//internal version: replace begin begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	//internal version: replace begin end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		resourceName := tccommon.BuildTagResourceName("mongodb", "instance", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	//internal version: replace end begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	//internal version: replace end end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	return resourceTencentCloudMongodbInstanceRead(d, meta)
}

func resourceTencentCloudMongodbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()

	instanceId := d.Id()

	mongodbService := MongodbService{client}
	tagService := svctag.NewTagService(client)
	instance, has, err := mongodbService.DescribeInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	if nilFields := tccommon.CheckNil(instance, map[string]string{
		"InstanceName": "instance name",
		"ProjectId":    "project id",
		"Zone":         "available zone",
		"VpcId":        "vpc id",
		"SubnetId":     "subnet id",
		"Status":       "status",
		"Vip":          "vip",
		"Vport":        "vport",
		"CreateTime":   "create time",
		"MongoVersion": "engine version",
		"Memory":       "memory",
		"Volume":       "volume",
		"MachineType":  "machine type",
	}); len(nilFields) > 0 {
		return fmt.Errorf("mongodb %v are nil", nilFields)
	}

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
	_ = d.Set("node_num", *instance.SecondaryNum+1)
	if instance.MaintenanceStart != nil && len(*instance.MaintenanceStart) == 8 {
		_ = d.Set("maintenance_start", (*instance.MaintenanceStart)[:5])
	}
	if instance.MaintenanceEnd != nil && len(*instance.MaintenanceEnd) == 8 {
		_ = d.Set("maintenance_end", (*instance.MaintenanceEnd)[:5])
	}

	replicateSets, err := mongodbService.DescribeDBInstanceNodeProperty(ctx, instanceId)
	if err != nil {
		return err
	}
	if len(replicateSets) > 0 {
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

	// standby instance list
	var standbyInsList []map[string]string
	for _, v := range instance.StandbyInstances {
		standbyInsList = append(
			standbyInsList,
			map[string]string{
				"standby_instance_id":     *v.InstanceId,
				"standby_instance_region": *v.Region,
			},
		)
	}

	// if not standby instance, need set `standby_instance_list`
	if _, ok := d.GetOk("father_instance_id"); !ok {
		_ = d.Set("standby_instance_list", standbyInsList)
	}

	tags, _ := tagService.DescribeResourceTags(ctx, "mongodb", "instance", client.Region, instanceId)

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudMongodbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	mongodbService := MongodbService{client: client}
	tagService := svctag.NewTagService(client)
	region := client.Region

	d.Partial(true)

	immutableArgs := []string{"availability_zone_list", "hidden_zone"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("memory") || d.HasChange("volume") || d.HasChange("node_num") {
		memory := d.Get("memory").(int)
		volume := d.Get("volume").(int)
		params := make(map[string]interface{})
		var nodeNum int
		if v, ok := d.GetOk("node_num"); ok {
			nodeNum = v.(int)
			params["node_num"] = nodeNum
		}
		if v, ok := d.GetOk("add_node_list"); ok {
			addNodeList := v.([]interface{})
			params["add_node_list"] = addNodeList
		}
		if v, ok := d.GetOk("remove_node_list"); ok {
			removeNodeList := v.([]interface{})
			params["remove_node_list"] = removeNodeList
		}
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
					if sdkError, ok := err.(*errors.TencentCloudSDKError); ok {
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

	if d.HasChange("vpc_id") || d.HasChange("subnet_id") {
		vpcId := d.Get("vpc_id").(string)
		subnetId := d.Get("subnet_id").(string)

		err := mongodbService.ModifyNetworkAddress(ctx, instanceId, vpcId, subnetId)
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

	if d.HasChange("maintenance_start") || d.HasChange("maintenance_end") {
		maintenanceStart, okMaintenanceStart := d.GetOk("maintenance_start")
		maintenanceEnd, okMaintenanceEnd := d.GetOk("maintenance_end")
		if okMaintenanceStart && okMaintenanceEnd {
			err := mongodbService.SetInstanceMaintenance(ctx, instanceId, maintenanceStart.(string), maintenanceEnd.(string))
			if err != nil {
				return err
			}
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

	return resourceTencentCloudMongodbInstanceRead(d, meta)
}

func resourceTencentCloudMongodbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance.delete")()

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
