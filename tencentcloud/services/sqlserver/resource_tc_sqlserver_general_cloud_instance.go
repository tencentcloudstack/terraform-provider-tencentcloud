package sqlserver

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSqlserverGeneralCloudInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverGeneralCloudInstanceCreate,
		Read:   resourceTencentCloudSqlserverGeneralCloudInstanceRead,
		Update: resourceTencentCloudSqlserverGeneralCloudInstanceUpdate,
		Delete: resourceTencentCloudSqlserverGeneralCloudInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "Name of the SQL Server instance.",
			},
			"zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance AZ, such as ap-guangzhou-1 (Guangzhou Zone 1). Purchasable AZs for an instance can be obtained through the DescribeZones API.",
			},
			"memory": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Memory, unit: GB.",
			},
			"storage": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "instance disk storage, unit: GB.",
			},
			"cpu": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Cpu, unit: CORE.",
			},
			"machine_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The host disk type of the purchased instance, CLOUD_HSSD-enhanced SSD cloud disk for virtual machines, CLOUD_TSSD-extremely fast SSD cloud disk for virtual machines, CLOUD_BSSD-universal SSD cloud disk for virtual machines.",
			},
			"instance_charge_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Payment mode, the value supports PREPAID (prepaid), POSTPAID (postpaid).",
			},
			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "project ID.",
			},
			"subnet_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "VPC subnet ID, in the form of subnet-bdoe83fa; SubnetId and VpcId need to be set at the same time or not set at the same time.",
			},
			"vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "VPC network ID, in the form of vpc-dsp338hz; SubnetId and VpcId need to be set at the same time or not set at the same time.",
			},
			"period": {
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 48),
				Description:  "Purchase instance period, the default value is 1, which means one month. The value cannot exceed 48. Valid only when the 'instance_charge_type' parameter value is 'PREPAID'.",
			},
			"db_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "sqlserver version, currently all supported versions are: 2008R2 (SQL Server 2008 R2 Enterprise), 2012SP3 (SQL Server 2012 Enterprise), 201202 (SQL Server 2012 Standard), 2014SP2 (SQL Server 2014 Enterprise), 201402 (SQL Server 2014 Standard), 2016SP1 (SQL Server 2016 Enterprise), 201602 (SQL Server 2016 Standard), 2017 (SQL Server 2017 Enterprise), 201702 (SQL Server 2017 Standard), 2019 (SQL Server 2019 Enterprise), 201902 (SQL Server 2019 Standard). Each region supports different versions for sale, and the version information that can be sold in each region can be pulled through the DescribeProductConfig interface. If left blank, the default version is 2008R2.",
			},
			"auto_renew_flag": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Automatic renewal flag: 0-normal renewal 1-automatic renewal, the default is 1 automatic renewal. Valid only when purchasing a prepaid instance. Valid only when the 'instance_charge_type' parameter value is 'PREPAID'.",
			},
			"security_group_list": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Security group list, fill in the security group ID in the form of sg-xxx.",
			},
			"weekly": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Maintainable time window configuration, in weeks, indicates the days of the week that allow maintenance, 1-7 represent Monday to weekend respectively.",
			},
			"start_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Maintainable time window configuration, daily maintainable start time.",
			},
			"span": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Maintainable time window configuration, duration, unit: hour.",
			},
			"resource_tags": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "A collection of tags bound to the new instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tag value.",
						},
					},
				},
			},
			"collation": {
				Optional:    true,
				Type:        schema.TypeString,
				Default:     "Chinese_PRC_CI_AS",
				Description: "System character set collation, default: Chinese_PRC_CI_AS.",
			},
			"time_zone": {
				Optional:    true,
				Type:        schema.TypeString,
				Default:     "China Standard Time",
				Description: "System time zone, default: China Standard Time.",
			},
			"ha_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Deprecated:  "It has been deprecated from version 1.81.2.",
				Description: "Upgrade the high-availability architecture of sqlserver, upgrade from mirror disaster recovery to always on cluster disaster recovery, only support 2017 and above and support always on high-availability instances, do not support downgrading to mirror disaster recovery, CLUSTER-upgrade to always on capacity Disaster, if not filled, the high-availability architecture will not be modified.",
			},
		},
	}
}

func resourceTencentCloudSqlserverGeneralCloudInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_general_cloud_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service      = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request      = sqlserver.NewCreateCloudDBInstancesRequest()
		instanceId   string
		instanceName string
		dealId       string
	)

	if v, ok := d.GetOk("name"); ok {
		instanceName = v.(string)
	}

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("memory"); ok {
		request.Memory = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("storage"); ok {
		request.Storage = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("cpu"); ok {
		request.Cpu = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("machine_type"); ok {
		request.MachineType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
		if v.(string) == SQLSERVER_TYPE_PREPAID {
			if v, ok := d.GetOk("period"); ok {
				request.Period = helper.IntInt64(v.(int))
			}

			if v, ok := d.GetOk("auto_renew_flag"); ok {
				request.AutoRenewFlag = helper.IntInt64(v.(int))
			}
		}
	}

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_version"); ok {
		request.DBVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_list"); ok {
		securityGroupListSet := v.(*schema.Set).List()
		for i := range securityGroupListSet {
			securityGroupList := securityGroupListSet[i].(string)
			request.SecurityGroupList = append(request.SecurityGroupList, &securityGroupList)
		}
	}

	if v, ok := d.GetOk("weekly"); ok {
		weeklySet := v.(*schema.Set).List()
		for i := range weeklySet {
			weekly := weeklySet[i].(int)
			request.Weekly = append(request.Weekly, helper.IntInt64(weekly))
		}
	}

	if v, ok := d.GetOk("start_time"); ok {
		request.StartTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("span"); ok {
		request.Span = helper.IntInt64(v.(int))
	}

	request.MultiZones = helper.Bool(true)

	if v, ok := d.GetOk("resource_tags"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			resourceTag := sqlserver.ResourceTag{}
			if v, ok := dMap["tag_key"]; ok {
				resourceTag.TagKey = helper.String(v.(string))
			}
			if v, ok := dMap["tag_value"]; ok {
				resourceTag.TagValue = helper.String(v.(string))
			}
			request.ResourceTags = append(request.ResourceTags, &resourceTag)
		}
	}

	if v, ok := d.GetOk("collation"); ok {
		request.Collation = helper.String(v.(string))
	}

	if v, ok := d.GetOk("time_zone"); ok {
		request.TimeZone = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().CreateCloudDBInstances(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		dealId = *result.Response.DealName
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver generalCloudInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId, err = service.GetInfoFromDeal(ctx, dealId)
	if err != nil {
		return err
	}

	// set name
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		inErr := service.ModifySqlserverInstanceName(ctx, instanceId, instanceName)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}

		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverGeneralCloudInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralCloudInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_general_cloud_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	generalCloudInstance, err := service.DescribeSqlserverGeneralCloudInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if generalCloudInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverGeneralCloudInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if generalCloudInstance.Name != nil {
		_ = d.Set("name", generalCloudInstance.Name)
	}

	if generalCloudInstance.Zone != nil {
		_ = d.Set("zone", generalCloudInstance.Zone)
	}

	if generalCloudInstance.Memory != nil {
		_ = d.Set("memory", generalCloudInstance.Memory)
	}

	if generalCloudInstance.Storage != nil {
		_ = d.Set("storage", generalCloudInstance.Storage)
	}

	if generalCloudInstance.Cpu != nil {
		_ = d.Set("cpu", generalCloudInstance.Cpu)
	}

	if generalCloudInstance.Type != nil {
		_ = d.Set("machine_type", generalCloudInstance.Type)
	}

	if generalCloudInstance.PayMode != nil {
		if *generalCloudInstance.PayMode == 0 {
			_ = d.Set("instance_charge_type", SQLSERVER_TYPE_POSTPAID)
		} else {
			_ = d.Set("instance_charge_type", SQLSERVER_TYPE_PREPAID)
		}
	}

	if generalCloudInstance.ProjectId != nil {
		_ = d.Set("project_id", generalCloudInstance.ProjectId)
	}

	if generalCloudInstance.UniqSubnetId != nil {
		_ = d.Set("subnet_id", generalCloudInstance.UniqSubnetId)
	}

	if generalCloudInstance.UniqVpcId != nil {
		_ = d.Set("vpc_id", generalCloudInstance.UniqVpcId)
	}

	if generalCloudInstance.Version != nil {
		_ = d.Set("db_version", generalCloudInstance.Version)
	}

	if generalCloudInstance.RenewFlag != nil {
		_ = d.Set("auto_renew_flag", generalCloudInstance.RenewFlag)
	}

	if generalCloudInstance.ResourceTags != nil {
		resourceTagsList := []interface{}{}
		for _, resourceTags := range generalCloudInstance.ResourceTags {
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

	if generalCloudInstance.Collation != nil {
		_ = d.Set("collation", generalCloudInstance.Collation)
	}

	if generalCloudInstance.TimeZone != nil {
		_ = d.Set("time_zone", generalCloudInstance.TimeZone)
	}

	maintenanceSpan, err := service.DescribeMaintenanceSpanById(ctx, instanceId)
	if err != nil {
		return err
	}

	if maintenanceSpan == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlservereMaintenanceSpan` [%s] not found, please check if it has been deleted.", logId, d.Id())
		return nil
	}

	if maintenanceSpan.Span != nil {
		_ = d.Set("span", maintenanceSpan.Span)
	}

	if maintenanceSpan.StartTime != nil {
		_ = d.Set("start_time", maintenanceSpan.StartTime)
	}

	if maintenanceSpan.Weekly != nil {
		_ = d.Set("weekly", maintenanceSpan.Weekly)
	}

	securityGroupList, err := service.DescribeInstanceSecurityGroups(ctx, instanceId)
	if err != nil {
		return err
	}

	if securityGroupList == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlservereSecurityGroups` [%s] not found, please check if it has been deleted.", logId, d.Id())
		return nil
	}

	if securityGroupList != nil {
		_ = d.Set("security_group_list", securityGroupList)
	}

	return nil
}

func resourceTencentCloudSqlserverGeneralCloudInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_general_cloud_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		ctx              = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		client           = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		sqlserverService = SqlserverService{client: client}
		request          = sqlserver.NewUpgradeDBInstanceRequest()
		instanceId       = d.Id()
		waitSwitch       int64
		dealId           string
		instanceName     string
	)

	request.InstanceId = &instanceId
	immutableArgs := []string{"zone", "machine_type", "instance_charge_type", "project_id", "subnet_id", "vpc_id", "period", "security_group_list", "weekly", "start_time", "span", "resource_tags", "collation", "time_zone"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			instanceName = v.(string)

			// set name
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				inErr := sqlserverService.ModifySqlserverInstanceName(ctx, instanceId, instanceName)
				if inErr != nil {
					return tccommon.RetryError(inErr)
				}

				return nil
			})

			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("memory") {
		if v, ok := d.GetOk("memory"); ok {
			request.Memory = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("storage") {
		if v, ok := d.GetOk("storage"); ok {
			request.Storage = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("cpu") {
		if v, ok := d.GetOk("cpu"); ok {
			request.Cpu = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("db_version") {
		if v, ok := d.GetOk("db_version"); ok {
			request.DBVersion = helper.String(v.(string))
		}
	}

	waitSwitch = 0
	request.WaitSwitch = &waitSwitch

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().UpgradeDBInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		dealId = *result.Response.DealName
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver generalCloudInstance failed, reason:%+v", logId, err)
		return err
	}

	_, err = sqlserverService.GetInfoFromDeal(ctx, dealId)
	if err != nil {
		return err
	}

	return resourceTencentCloudSqlserverGeneralCloudInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralCloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_general_cloud_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	if err := service.TerminateSqlserverInstanceById(ctx, instanceId); err != nil {
		return err
	}

	if err := service.DeleteSqlserverInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
