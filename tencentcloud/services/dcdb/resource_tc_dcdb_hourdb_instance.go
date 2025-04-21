package dcdb

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDcdbHourdbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbHourdbInstanceCreate,
		Read:   resourceTencentCloudDcdbHourdbInstanceRead,
		Update: resourceTencentCloudDcdbHourdbInstanceUpdate,
		Delete: resourceTencentCloudDcdbHourdbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zones": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "available zone.",
			},

			"shard_memory": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "memory(GB) for each shard. It can be obtained by querying api DescribeShardSpec.",
			},

			"shard_storage": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "storage(GB) for each shard. It can be obtained by querying api DescribeShardSpec.",
			},

			"shard_node_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "node count for each shard. It can be obtained by querying api DescribeShardSpec.",
			},

			"shard_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "instance shard count.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "vpc id.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "subnet id, its required when vpcId is set.",
			},

			"db_version_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "db engine version, default to Percona 5.7.17.",
			},

			"security_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "security group id.",
			},

			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "project id.",
			},

			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "name of this instance.",
			},

			"dcn_region": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "DCN source region.",
			},

			"dcn_instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "DCN source instance ID.",
			},

			"ipv6_flag": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to support IPv6.",
			},

			"extranet_access": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to open the extranet access.",
			},

			"vip": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The field is required to specify VIP.",
			},

			"vipv6": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The field is required to specify VIPv6.",
			},

			"vport": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Intranet port.",
			},

			"resource_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "resource tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "tag value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDcdbHourdbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_hourdb_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service       = DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request       = dcdb.NewCreateHourDCDBInstanceRequest()
		response      *dcdb.CreateHourDCDBInstanceResponse
		instanceId    string
		dcnInstanceId string
		vpcId         string
		subnetId      string
		ipv6Flag      int
	)

	if v, ok := d.GetOk("zones"); ok {
		zonesSet := v.(*schema.Set).List()
		for i := range zonesSet {
			zones := zonesSet[i].(string)
			request.Zones = append(request.Zones, &zones)
		}
	}

	if v, ok := d.GetOkExists("shard_memory"); ok {
		request.ShardMemory = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("shard_storage"); ok {
		request.ShardStorage = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("shard_node_count"); ok {
		request.ShardNodeCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("shard_count"); ok {
		request.ShardCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
		vpcId = v.(string)
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
		subnetId = v.(string)
	}

	if v, ok := d.GetOk("db_version_id"); ok {
		request.DbVersionId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_id"); ok {
		request.SecurityGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		request.ProjectId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, _ := d.GetOkExists("ipv6_flag"); v != nil {
		request.Ipv6Flag = helper.IntInt64(v.(int))
		ipv6Flag = v.(int)
	}

	if v, ok := d.GetOk("resource_tags"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			resourceTag := dcdb.ResourceTag{}
			if v, ok := dMap["tag_key"]; ok {
				resourceTag.TagKey = helper.String(v.(string))
			}

			if v, ok := dMap["tag_value"]; ok {
				resourceTag.TagValue = helper.String(v.(string))
			}

			request.ResourceTags = append(request.ResourceTags, &resourceTag)
		}
	}

	if v, ok := d.GetOk("dcn_region"); ok {
		request.DcnRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dcn_instance_id"); ok {
		request.DcnInstanceId = helper.String(v.(string))
		dcnInstanceId = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDcdbClient().CreateHourDCDBInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create dcdb hourdbInstance failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dcdb hourdbInstance failed, reason:%+v", logId, err)
		return err
	}

	if len(response.Response.InstanceIds) != 1 {
		return fmt.Errorf("InstanceIds is error.")
	}

	instanceId = *response.Response.InstanceIds[0]

	d.SetId(instanceId)

	defaultInitParams := []*dcdb.DBParamValue{
		{
			Param: helper.String("character_set_server"),
			Value: helper.String("utf8mb4"),
		},
		{
			Param: helper.String("lower_case_table_names"),
			Value: helper.String("1"),
		},
		{
			Param: helper.String("sync_mode"),
			Value: helper.String("2"),
		},
		{
			Param: helper.String("innodb_page_size"),
			Value: helper.String("16384"),
		},
	}

	initRet, flowId, err := service.InitDcdbDbInstance(ctx, instanceId, defaultInitParams)
	if err != nil {
		return err
	}

	if !initRet {
		return fmt.Errorf("db instance init failed")
	}

	if flowId != nil {
		// need to wait init operation success
		// 0:success; 1:failed, 2:running
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"0"}, 3*tccommon.ReadRetryTimeout, time.Second, service.DcdbDbInstanceStateRefreshFunc(helper.UInt64Int64(*flowId), []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	if dcnInstanceId != "" {
		// need to wait dcn init processing complete
		// 0:none; 1:creating, 2:running
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"2"}, 3*tccommon.ReadRetryTimeout, time.Second, service.DcdbDcnStateRefreshFunc(instanceId, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	if v, ok := d.GetOkExists("extranet_access"); ok && v != nil {
		flag := v.(bool)
		err := service.SetDcdbExtranetAccess(ctx, instanceId, ipv6Flag, flag)
		if err != nil {
			return err
		}
	}

	var (
		vip   string
		vipv6 string
	)

	if v, ok := d.GetOk("vip"); ok {
		vip = v.(string)
	}

	if v, ok := d.GetOk("vipv6"); ok {
		vipv6 = v.(string)
	}

	if vip != "" || vipv6 != "" {
		if vpcId == "" || subnetId == "" {
			return fmt.Errorf("`vpc_id` and `subnet_id` cannot be empty when setting `vip` or `vipv6` fields!")
		}

		err := service.SetNetworkVip(ctx, instanceId, vpcId, subnetId, vip, vipv6)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudDcdbHourdbInstanceRead(d, meta)
}

func resourceTencentCloudDcdbHourdbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_hourdb_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	hourdbInstances, err := service.DescribeDcdbHourdbInstance(ctx, instanceId)
	if err != nil {
		return err
	}

	if hourdbInstances == nil {
		d.SetId("")
		return fmt.Errorf("resource `hourdbInstance` %s does not exist", instanceId)
	}

	if *hourdbInstances.TotalCount > 1 || len(hourdbInstances.Instances) > 1 {
		d.SetId("")
		return fmt.Errorf("the count of resource `hourdbInstance` shoud not beyond one. count: %v\n", hourdbInstances.TotalCount)
	}

	hourdbInstance := hourdbInstances.Instances[0]
	if hourdbInstance.ShardDetail[0] != nil { // Memory and Storage is params for one shard
		shard := hourdbInstance.ShardDetail[0]
		if shard.Memory != nil {
			_ = d.Set("shard_memory", shard.Memory)
		}

		if shard.Storage != nil {
			_ = d.Set("shard_storage", shard.Storage)
		}
	}

	if hourdbInstance.NodeCount != nil {
		_ = d.Set("shard_node_count", hourdbInstance.NodeCount)
	}

	if hourdbInstance.ShardCount != nil {
		_ = d.Set("shard_count", hourdbInstance.ShardCount)
	}

	if hourdbInstance.UniqueVpcId != nil {
		_ = d.Set("vpc_id", hourdbInstance.UniqueVpcId)
	}

	if hourdbInstance.UniqueSubnetId != nil {
		_ = d.Set("subnet_id", hourdbInstance.UniqueSubnetId)
	}

	if hourdbInstance.DbVersionId != nil {
		_ = d.Set("db_version_id", hourdbInstance.DbVersionId)
	}

	if hourdbInstance.ProjectId != nil {
		_ = d.Set("project_id", hourdbInstance.ProjectId)
	}

	if hourdbInstance.InstanceName != nil {
		_ = d.Set("instance_name", hourdbInstance.InstanceName)
	}

	if hourdbInstance.Ipv6Flag != nil {
		_ = d.Set("ipv6_flag", hourdbInstance.Ipv6Flag)
	}

	if hourdbInstance.WanStatus != nil {
		//0-未开通；1-已开通；2-关闭；3-开通中
		if *hourdbInstance.WanStatus == DCDB_WAN_STATUS_UNOPEN || *hourdbInstance.WanStatus == DCDB_WAN_STATUS_CLOSED {
			_ = d.Set("extranet_access", false)
		}

		if *hourdbInstance.WanStatus == DCDB_WAN_STATUS_OPENED {
			_ = d.Set("extranet_access", true)
		}
	}

	if hourdbInstance.ResourceTags != nil {
		resourceTagsList := []interface{}{}
		for _, resourceTags := range hourdbInstance.ResourceTags {
			resourceTagsMap := map[string]interface{}{}
			if resourceTags.TagKey != nil {
				resourceTagsMap["tag_key"] = resourceTags.TagKey
			}

			if resourceTags.TagValue != nil {
				resourceTagsMap["tag_value"] = resourceTags.TagValue
			}

			resourceTagsList = append(resourceTagsList, resourceTagsMap)
		}
		_ = d.Set("resource_tags", resourceTagsList)
	}

	if sg, err := service.DescribeDcdbSecurityGroup(ctx, instanceId); err == nil {
		sgId := ""
		if sg != nil && len(sg.Groups) > 0 {
			sgId = *sg.Groups[0].SecurityGroupId
		}

		_ = d.Set("security_group_id", sgId)
	} else {
		return err
	}

	// set dcn id and region
	if dcns, err := service.DescribeDcnDetailById(ctx, instanceId); err == nil {
		for _, dcn := range dcns {
			var master *dcdb.DcnDetailItem
			if *dcn.DcnFlag == DCDB_DCN_FLAG_MASTER {
				master = dcn
				_ = d.Set("dcn_region", master.Region)
				_ = d.Set("dcn_instance_id", master.InstanceId)
			}
		}
	} else {
		return err
	}

	// set vip, vipv6 and vport
	if detail, err := service.DescribeDcdbDbInstanceDetailById(ctx, instanceId); err == nil {
		if detail != nil {
			_ = d.Set("vip", detail.Vip)
			_ = d.Set("vipv6", detail.Vip6)
			_ = d.Set("vport", detail.Vport)

			if detail.MasterZone != nil {
				zones := []*string{detail.MasterZone}
				if detail.SlaveZones != nil {
					zones = append(zones, detail.SlaveZones...)
				}

				_ = d.Set("zones", zones)
			}
		}
	} else {
		return err
	}

	return nil
}

func resourceTencentCloudDcdbHourdbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_hourdb_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	immutableArgs := []string{"zones", "shard_node_count", "shard_count", "db_version_id", "dcn_region", "dcn_instance_id", "resource_tags"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if v, ok := d.GetOkExists("extranet_access"); ok && v != nil {
		flag := v.(bool)
		var ipv6Flag int
		if v, _ := d.GetOk("ipv6_flag"); v != nil {
			ipv6Flag = v.(int)
		}

		err := service.SetDcdbExtranetAccess(ctx, instanceId, ipv6Flag, flag)
		if err != nil {
			return err
		}

		time.Sleep(2 * time.Second)
	}

	if d.HasChange("project_id") {
		if projectId, ok := d.GetOk("project_id"); ok {
			request := dcdb.NewModifyDBInstancesProjectRequest()
			request.InstanceIds = []*string{&instanceId}
			request.ProjectId = helper.IntInt64(projectId.(int))
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDcdbClient().ModifyDBInstancesProject(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s operate dcdb modifyInstanceProjectOperation failed, reason:%+v", logId, err)
				return err
			}
		}

		time.Sleep(2 * time.Second)
	}

	if d.HasChange("vpc_id") || d.HasChange("subnet_id") || d.HasChange("vip") || d.HasChange("vipv6") {
		var (
			vip      string
			vipv6    string
			vpcId    string
			subnetId string
		)

		if v, ok := d.GetOk("vip"); ok {
			vip = v.(string)
		}

		if v, ok := d.GetOk("vipv6"); ok {
			vipv6 = v.(string)
		}

		if v, ok := d.GetOk("vpc_id"); ok {
			vpcId = v.(string)
		}

		if v, ok := d.GetOk("subnet_id"); ok {
			subnetId = v.(string)
		}

		if vpcId == "" || subnetId == "" {
			return fmt.Errorf("`vpc_id` and `subnet_id` cannot be empty when updating network configs!")
		}

		err := service.SetNetworkVip(ctx, instanceId, vpcId, subnetId, vip, vipv6)
		if err != nil {
			return err
		}
	}

	if d.HasChange("instance_name") {
		request := dcdb.NewModifyDBInstanceNameRequest()
		if v, ok := d.GetOk("instance_name"); ok {
			request.InstanceName = helper.String(v.(string))
		}

		request.InstanceId = &instanceId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDcdbClient().ModifyDBInstanceName(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify dcdb hourdbInstance name failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("shard_memory") || d.HasChange("shard_storage") {
		// get ShardInstanceIds
		hourdbInstances, err := service.DescribeDcdbHourdbInstance(ctx, instanceId)
		if err != nil {
			return err
		}

		if hourdbInstances == nil {
			d.SetId("")
			return fmt.Errorf("resource `hourdbInstance` %s does not exist", instanceId)
		}

		if *hourdbInstances.TotalCount > 1 || len(hourdbInstances.Instances) > 1 {
			d.SetId("")
			return fmt.Errorf("the count of resource `hourdbInstance` shoud not beyond one. count: %v\n", hourdbInstances.TotalCount)
		}

		hourdbInstance := hourdbInstances.Instances[0]
		shardInstanceIds := make([]string, 0, len(hourdbInstance.ShardDetail))
		if len(hourdbInstance.ShardDetail) > 0 { // Memory and Storage is params for one shard
			for _, item := range hourdbInstance.ShardDetail {
				shardInstanceIds = append(shardInstanceIds, *item.ShardInstanceId)
			}
		}

		var (
			shardMemory  int
			shardStorage int
		)

		if v, ok := d.GetOkExists("shard_memory"); ok {
			shardMemory = v.(int)
		}

		if v, ok := d.GetOkExists("shard_storage"); ok {
			shardStorage = v.(int)
		}

		request := dcdb.NewUpgradeHourDCDBInstanceRequest()
		request.InstanceId = &instanceId
		request.UpgradeType = helper.String("EXPAND")
		request.ExpandShardConfig = &dcdb.ExpandShardConfig{
			ShardInstanceIds: helper.Strings(shardInstanceIds),
			ShardMemory:      helper.IntInt64(shardMemory),
			ShardStorage:     helper.IntInt64(shardStorage),
		}

		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDcdbClient().UpgradeHourDCDBInstance(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify dcdb hourdbInstance config failed, reason:%+v", logId, err)
			return err
		}

		// wait
		err = resource.Retry(tccommon.ReadRetryTimeout*10, func() *resource.RetryError {
			dbInstances, errResp := service.DescribeDcdbDbInstance(ctx, instanceId)
			if errResp != nil {
				return tccommon.RetryError(errResp, tccommon.InternalError)
			}

			if dbInstances.Instances[0] == nil {
				return resource.NonRetryableError(fmt.Errorf("DescribeDcdbDbInstance return result(dcdb instance) is nil!"))
			}

			dbInstance := dbInstances.Instances[0]
			if *dbInstance.Status == 2 {
				return nil
			}

			return resource.RetryableError(fmt.Errorf("dcdb instance status is %v, retry...", *dbInstance.Status))
		})

		if err != nil {
			return err
		}
	}

	return resourceTencentCloudDcdbHourdbInstanceRead(d, meta)
}

func resourceTencentCloudDcdbHourdbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_hourdb_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := service.DeleteDcdbHourdbInstanceById(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
