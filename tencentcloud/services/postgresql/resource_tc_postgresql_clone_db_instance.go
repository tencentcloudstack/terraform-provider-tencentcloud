package postgresql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	postgresv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
)

func ResourceTencentCloudPostgresqlCloneDbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlCloneDbInstanceCreate,
		Read:   resourceTencentCloudPostgresqlCloneDbInstanceRead,
		Update: resourceTencentCloudPostgresqlCloneDbInstanceUpdate,
		Delete: resourceTencentCloudPostgresqlCloneDbInstanceDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the original instance to be cloned.",
			},

			"spec_code": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Purchasable code, which can be obtained from the `SpecCode` field in the return value of the [DescribeClasses](https://intl.cloud.tencent.com/document/api/409/89019?from_cn_redirect=1) API.",
			},

			"storage": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Instance storage capacity in GB.",
			},

			"period": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Purchase duration, in months.\n- Prepaid: Supports `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, and `36`.\n- Pay-as-you-go: Only supports `1`.",
			},

			"auto_renew_flag": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Renewal Flag:\n\n- `0`: manual renewal\n`1`: auto-renewal\n\nDefault value: 0.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC ID in the format of `vpc-xxxxxxx`, which can be obtained in the console or from the `unVpcId` field in the return value of the [DescribeVpcEx](https://intl.cloud.tencent.com/document/api/215/1372?from_cn_redirect=1) API.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC subnet ID in the format of `subnet-xxxxxxxx`, which can be obtained in the console or from the `unSubnetId` field in the return value of the [DescribeSubnets](https://intl.cloud.tencent.com/document/api/215/15784?from_cn_redirect=1) API.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the newly purchased instance, which can contain up to 60 letters, digits, or symbols (-_). If this parameter is not specified, \"Unnamed\" will be displayed by default.",
			},

			"instance_charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance billing type, which currently supports:\n\n- PREPAID: Prepaid, i.e., monthly subscription\n- POSTPAID_BY_HOUR: Pay-as-you-go, i.e., pay by consumption\n\nDefault value: PREPAID.",
			},

			"security_group_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Security group of the instance, which can be obtained from the `sgld` field in the return value of the [DescribeSecurityGroups](https://intl.cloud.tencent.com/document/api/215/15808?from_cn_redirect=1) API. If this parameter is not specified, the default security group will be bound.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Project ID.",
			},

			"tag_list": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "The information of tags to be bound with the instance, which is left empty by default. This parameter can be obtained from the `Tags` field in the return value of the [DescribeTags](https://intl.cloud.tencent.com/document/api/651/35316?from_cn_redirect=1) API.",
				Deprecated:  "It has been deprecated from version 1.83.10. Use `tags` instead.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The available tags within this postgresql.",
			},

			"db_node_set": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Deployment information of the instance node, which will display the information of each AZ when the instance node is deployed across multiple AZs.\nThe information of AZ can be obtained from the `Zone` field in the return value of the [DescribeZones](https://intl.cloud.tencent.com/document/api/409/16769?from_cn_redirect=1) API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Node type. Valid values:\n`Primary`;\n`Standby`.",
						},
						"zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "AZ where the node resides, such as ap-guangzhou-1.",
						},
						"dedicated_cluster_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Dedicated cluster ID.",
						},
					},
				},
			},

			"activity_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Campaign ID.",
			},

			"backup_set_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Basic backup set ID.",
			},

			"recovery_target_time": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Restoration point in time.",
			},

			"sync_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Primary-standby sync mode, which supports:\nSemi-sync: Semi-sync\nAsync: Asynchronous\nDefault value for the primary instance: Semi-sync\nDefault value for the read-only instance: Async.",
			},

			"deletion_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether deletion protection is enabled for the instance: `true` deletion protection enabled; `false` deletion protection disabled.",
			},

			// computed
			"new_db_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the cloned instance.",
			},

			"availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Availability zone.",
			},

			"public_access_host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Host for public access.",
			},
			"public_access_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Port for public access.",
			},
			"private_access_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP for private access.",
			},
			"private_access_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Port for private access.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlCloneDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_clone_db_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = postgresv20170312.NewCloneDBInstanceRequest()
		response = postgresv20170312.NewCloneDBInstanceResponse()
	)

	if v, ok := d.GetOk("db_instance_id"); ok {
		request.DBInstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("spec_code"); ok {
		request.SpecCode = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("storage"); ok {
		request.Storage = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("period"); ok {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_renew_flag"); ok {
		request.AutoRenewFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, helper.String(securityGroupIds))
		}
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		request.ProjectId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("tag_list"); ok {
		for _, item := range v.([]interface{}) {
			tagListMap := item.(map[string]interface{})
			tag := postgresv20170312.Tag{}
			if v, ok := tagListMap["tag_key"]; ok {
				tag.TagKey = helper.String(v.(string))
			}

			if v, ok := tagListMap["tag_value"]; ok {
				tag.TagValue = helper.String(v.(string))
			}

			request.TagList = append(request.TagList, &tag)
		}
	}

	if v, ok := d.GetOk("db_node_set"); ok {
		for _, item := range v.([]interface{}) {
			dBNodeSetMap := item.(map[string]interface{})
			dBNode := postgresv20170312.DBNode{}
			if v, ok := dBNodeSetMap["role"]; ok {
				dBNode.Role = helper.String(v.(string))
			}

			if v, ok := dBNodeSetMap["zone"]; ok {
				dBNode.Zone = helper.String(v.(string))
			}

			if v, ok := dBNodeSetMap["dedicated_cluster_id"]; ok {
				dBNode.DedicatedClusterId = helper.String(v.(string))
			}

			request.DBNodeSet = append(request.DBNodeSet, &dBNode)
		}
	}

	if v, ok := d.GetOkExists("activity_id"); ok {
		request.ActivityId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("backup_set_id"); ok {
		request.BackupSetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("recovery_target_time"); ok {
		request.RecoveryTargetTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sync_mode"); ok {
		request.SyncMode = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("deletion_protection"); ok {
		request.DeletionProtection = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresV20170312Client().CloneDBInstanceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create postgresql clone db instance failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create postgresql clone db instance failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.DBInstanceId == nil {
		return fmt.Errorf("DBInstanceId is nil.")
	}

	dBInstanceId := *response.Response.DBInstanceId
	d.SetId(dBInstanceId)

	// wait
	if _, err := (&resource.StateChangeConf{
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
		Pending:    []string{},
		Refresh:    resourcePostgresqlCloneDbInstanceCreateStateRefreshFunc_0_0(ctx, dBInstanceId),
		Target:     []string{"running"},
		Timeout:    d.Timeout(schema.TimeoutCreate),
	}).WaitForStateContext(ctx); err != nil {
		return err
	}

	// set tags
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("postgres", "DBInstanceId", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudPostgresqlCloneDbInstanceRead(d, meta)
}

func resourceTencentCloudPostgresqlCloneDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_clone_db_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		instance      *postgresql.DBInstance
		has           bool
		outErr, inErr error
	)

	postgresqlService := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, has, inErr = postgresqlService.DescribePostgresqlInstanceById(ctx, d.Id())
		if inErr != nil {
			ee, ok := inErr.(*sdkErrors.TencentCloudSDKError)
			if ok && (ee.GetCode() == "ResourceNotFound.InstanceNotFoundError" || ee.GetCode() == "InvalidParameter") {
				return nil
			}

			return tccommon.RetryError(inErr)
		}

		if instance != nil && tccommon.IsContains(POSTGRESQL_RETRYABLE_STATUS, *instance.DBInstanceStatus) {
			return resource.RetryableError(fmt.Errorf("instance %s is %s, retrying", *instance.DBInstanceId, *instance.DBInstanceStatus))
		}

		return nil
	})

	if outErr != nil {
		return outErr
	}

	if !has {
		d.SetId("")
		return nil
	}

	if instance.DBInstanceId != nil {
		_ = d.Set("new_db_instance_id", instance.DBInstanceId)
	}

	if instance.DBInstanceStorage != nil {
		_ = d.Set("storage", instance.DBInstanceStorage)
	}

	if instance.AutoRenew != nil {
		_ = d.Set("auto_renew_flag", instance.AutoRenew)
	}

	if len(instance.DBInstanceNetInfo) == 3 {
		_ = d.Set("vpc_id", instance.DBInstanceNetInfo[0].VpcId)
		_ = d.Set("subnet_id", instance.DBInstanceNetInfo[0].SubnetId)
		_ = d.Set("private_access_ip", instance.DBInstanceNetInfo[0].Ip)
		_ = d.Set("private_access_port", instance.DBInstanceNetInfo[0].Port)

		for _, v := range instance.DBInstanceNetInfo {
			if *v.NetType == "public" {
				_ = d.Set("public_access_host", v.Address)
				_ = d.Set("public_access_port", v.Port)
			}
		}
	} else if len(instance.DBInstanceNetInfo) == 2 {
		_ = d.Set("vpc_id", instance.VpcId)
		_ = d.Set("subnet_id", instance.SubnetId)

		for _, v := range instance.DBInstanceNetInfo {
			if *v.NetType == "public" {
				_ = d.Set("public_access_host", v.Address)
				_ = d.Set("public_access_port", v.Port)
			}

			// private or inner will not appear at same time, private for instance with vpc
			if (*v.NetType == "private" || *v.NetType == "inner") && *v.Ip != "" {
				_ = d.Set("private_access_ip", v.Ip)
				_ = d.Set("private_access_port", v.Port)
			}
		}
	} else {
		return fmt.Errorf("DBInstanceNetInfo returned incorrect information.")
	}

	if instance.DBInstanceName != nil {
		_ = d.Set("name", instance.DBInstanceName)
	}

	if instance.PayType != nil {
		if *instance.PayType == "postpaid" {
			_ = d.Set("instance_charge_type", "POSTPAID_BY_HOUR")
		} else if *instance.PayType == "prepaid" {
			_ = d.Set("instance_charge_type", "PREPAID")
		} else {
			_ = d.Set("instance_charge_type", instance.PayType)
		}
	}

	if instance.ProjectId != nil {
		_ = d.Set("project_id", instance.ProjectId)
	}

	if instance.DeletionProtection != nil {
		_ = d.Set("deletion_protection", instance.DeletionProtection)
	}

	if instance.Zone != nil {
		_ = d.Set("availability_zone", instance.Zone)
	}

	nodeSet := instance.DBNodeSet
	zoneSet := schema.NewSet(schema.HashString, nil)
	if nodeCount := len(nodeSet); nodeCount > 0 {
		var dbNodeSet = make([]interface{}, 0, nodeCount)
		for i := range nodeSet {
			item := nodeSet[i]
			node := map[string]interface{}{
				"role":                 item.Role,
				"zone":                 item.Zone,
				"dedicated_cluster_id": item.DedicatedClusterId,
			}

			zoneSet.Add(*item.Zone)
			dbNodeSet = append(dbNodeSet, node)
		}

		// skip default set (single AZ and zone includes)
		_, nodeSetOk := d.GetOk("db_node_set")
		importedMaz := zoneSet.Len() > 1 && zoneSet.Contains(*instance.Zone)

		if nodeSetOk || importedMaz {
			_ = d.Set("db_node_set", dbNodeSet)
		}
	}

	// security groups
	sg, err := postgresqlService.DescribeDBInstanceSecurityGroupsById(ctx, d.Id())
	if err != nil {
		return err
	}

	if len(sg) > 0 {
		_ = d.Set("security_group_ids", sg)
	}

	// ignore spec_code
	// qcs::postgres:ap-guangzhou:uin/123435236:DBInstanceId/postgres-xxx
	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "postgres", "DBInstanceId", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudPostgresqlCloneDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_clone_db_instance.update")()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		postgresqlService = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId        = d.Id()
	)

	if d.HasChange("instance_charge_type") {
		var (
			request       = postgresql.NewModifyDBInstanceChargeTypeRequest()
			chargeTypeOld string
			chargeTypeNew string
		)

		old, new := d.GetChange("instance_charge_type")
		if old != nil {
			chargeTypeOld = old.(string)
		}

		if new != nil {
			chargeTypeNew = new.(string)
		}

		if chargeTypeOld != "POSTPAID_BY_HOUR" || chargeTypeNew != "PREPAID" {
			return fmt.Errorf("It only support to update the charge type from `POSTPAID_BY_HOUR` to `PREPAID`.")
		}

		request.DBInstanceId = &instanceId
		request.InstanceChargeType = &chargeTypeNew
		request.Period = helper.IntInt64(1)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().ModifyDBInstanceChargeType(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s operate postgresql ModifyDbInstanceChargeType failed, reason:%+v", logId, err)
			return err
		}

		// wait unit charge type changing operation of instance done
		service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), time.Second, service.PostgresqlDBInstanceStateRefreshFunc(instanceId, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	if d.HasChange("auto_renew_flag") {
		request := postgresql.NewSetAutoRenewFlagRequest()
		request.DBInstanceIdSet = helper.Strings([]string{instanceId})
		if v, ok := d.GetOkExists("auto_renew_flag"); ok {
			request.AutoRenewFlag = helper.IntInt64(v.(int))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().SetAutoRenewFlag(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s operate postgresql SetAutoRenewFlag failed, reason:%+v", logId, err)
			return err
		}

		// wait unit charge type changing operation of instance done
		service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), time.Second, service.PostgresqlDBInstanceStateRefreshFunc(instanceId, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	if d.HasChange("period") {
		request := postgresql.NewRenewInstanceRequest()
		request.DBInstanceId = &instanceId
		if v, ok := d.GetOkExists("period"); ok {
			request.Period = helper.IntInt64(v.(int))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().RenewInstance(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s operate postgresql RenewInstance failed, reason:%+v", logId, err)
			return err
		}

		// wait unit charge type changing operation of instance done
		service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), time.Second, service.PostgresqlDBInstanceStateRefreshFunc(instanceId, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		outErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr := postgresqlService.ModifyPostgresqlInstanceName(ctx, instanceId, name)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}

			return nil
		})

		if outErr != nil {
			return outErr
		}

		// check update name done
		checkErr := postgresqlService.CheckDBInstanceStatus(ctx, instanceId)
		if checkErr != nil {
			return checkErr
		}
	}

	if d.HasChange("project_id") {
		projectId := d.Get("project_id").(int)
		outErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr := postgresqlService.ModifyPostgresqlInstanceProjectId(ctx, instanceId, projectId)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}

			return nil
		})

		if outErr != nil {
			return outErr
		}

		// check update project id done
		checkErr := postgresqlService.CheckDBInstanceStatus(ctx, instanceId)
		if checkErr != nil {
			return checkErr
		}
	}

	if d.HasChange("security_group_ids") {
		ids := d.Get("security_group_ids").(*schema.Set).List()
		var sgIds []*string
		for _, id := range ids {
			sgIds = append(sgIds, helper.String(id.(string)))
		}

		err := postgresqlService.ModifyDBInstanceSecurityGroupsById(ctx, d.Id(), sgIds)
		if err != nil {
			return err
		}
	}

	if d.HasChange("db_node_set") {
		if include, z, nzs := checkZoneSetInclude(d); !include {
			return fmt.Errorf("`availability_zone`: %s is not included in `db_node_set`: %s", z, nzs)
		}

		nodeSet := d.Get("db_node_set").(*schema.Set).List()
		request := postgresql.NewModifyDBInstanceDeploymentRequest()
		request.DBInstanceId = helper.String(d.Id())
		request.SwitchTag = helper.IntInt64(0)
		for i := range nodeSet {
			var (
				node               = nodeSet[i].(map[string]interface{})
				role               = node["role"].(string)
				zone               = node["zone"].(string)
				dedicatedClusterId = node["dedicated_cluster_id"].(string)
			)

			if dedicatedClusterId != "" {
				request.DBNodeSet = append(request.DBNodeSet, &postgresql.DBNode{
					Role:               &role,
					Zone:               &zone,
					DedicatedClusterId: &dedicatedClusterId,
				})
			} else {
				request.DBNodeSet = append(request.DBNodeSet, &postgresql.DBNode{
					Role: &role,
					Zone: &zone,
				})
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			if err := postgresqlService.ModifyDBInstanceDeployment(ctx, request); err != nil {
				return tccommon.RetryError(err, postgresql.OPERATIONDENIED_INSTANCESTATUSLIMITOPERROR)
			}
			return nil
		})

		if err != nil {
			return err
		}

		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			instance, _, err := postgresqlService.DescribePostgresqlInstanceById(ctx, d.Id())
			if err != nil {
				return tccommon.RetryError(err)
			}

			if tccommon.IsContains(POSTGRESQL_RETRYABLE_STATUS, *instance.DBInstanceStatus) {
				return resource.RetryableError(fmt.Errorf("instance status is %s, retrying", *instance.DBInstanceStatus))
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	if d.HasChange("vpc_id") || d.HasChange("subnet_id") {
		var (
			postgresqlService = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
			instance          *postgresql.DBInstance
			has               bool
			inErr             error
			vpcOld            string
			vpcNew            string
			subnetOld         string
			subnetNew         string
			vipOld            string
			vipNew            string
		)

		// check net first
		outErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			instance, has, inErr = postgresqlService.DescribePostgresqlInstanceById(ctx, d.Id())
			if inErr != nil {
				ee, ok := inErr.(*sdkErrors.TencentCloudSDKError)
				if ok && (ee.GetCode() == "ResourceNotFound.InstanceNotFoundError" || ee.GetCode() == "InvalidParameter") {
					return nil
				}

				return tccommon.RetryError(inErr)
			}

			if instance != nil && tccommon.IsContains(POSTGRESQL_RETRYABLE_STATUS, *instance.DBInstanceStatus) {
				return resource.RetryableError(fmt.Errorf("instance %s is %s, retrying", *instance.DBInstanceId, *instance.DBInstanceStatus))
			}

			return nil
		})

		if outErr != nil {
			return outErr
		}

		if !has {
			d.SetId("")
			return nil
		}

		// check net num
		if instance.DBInstanceNetInfo != nil && len(instance.DBInstanceNetInfo) > 2 {
			return fmt.Errorf("There are already %d network information for the current PostgreSQL instance %s. Please remove one before modifying the instance network information.", len(instance.DBInstanceNetInfo)-1, d.Id())
		} else {
			old, new := d.GetChange("vpc_id")
			if old != nil {
				vpcOld = old.(string)
			}

			if new != nil {
				vpcNew = new.(string)
			}

			old, new = d.GetChange("subnet_id")
			if old != nil {
				subnetOld = old.(string)
			}

			if new != nil {
				subnetNew = new.(string)
			}

			// Create new network first, then delete the old one
			request := postgresql.NewCreateDBInstanceNetworkAccessRequest()
			request.DBInstanceId = helper.String(instanceId)
			request.VpcId = helper.String(vpcNew)
			request.SubnetId = helper.String(subnetNew)
			// ip assigned by system
			request.IsAssignVip = helper.Bool(false)
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().CreateDBInstanceNetworkAccess(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s create postgresql Instance NetworkAccess failed, reason:%+v", logId, err)
				return err
			}

			// wait for new network enabled
			conf := tccommon.BuildStateChangeConf([]string{}, []string{"opened"}, d.Timeout(schema.TimeoutUpdate), time.Second, postgresqlService.PostgresqlDBInstanceNetworkAccessStateRefreshFunc(instanceId, vpcNew, subnetNew, vipOld, "", []string{}))
			if object, e := conf.WaitForState(); e != nil {
				return e
			} else {
				// find the vip assiged by system
				ret := object.(*postgresql.DBInstanceNetInfo)
				vipNew = *ret.Ip
			}

			// wait unit network changing operation of instance done
			conf = tccommon.BuildStateChangeConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), time.Second, postgresqlService.PostgresqlDBInstanceStateRefreshFunc(instanceId, []string{}))
			if _, e := conf.WaitForState(); e != nil {
				return e
			}

			// delete the old one
			if v, ok := d.GetOk("private_access_ip"); ok {
				vipOld = v.(string)
			}

			if err := postgresqlService.DeletePostgresqlDBInstanceNetworkAccessById(ctx, instanceId, vpcOld, subnetOld, vipOld); err != nil {
				return err
			}

			// wait for old network removed
			conf = tccommon.BuildStateChangeConf([]string{}, []string{"closed"}, d.Timeout(schema.TimeoutUpdate), time.Second, postgresqlService.PostgresqlDBInstanceNetworkAccessStateRefreshFunc(instanceId, vpcOld, subnetOld, vipNew, vipOld, []string{}))
			if _, e := conf.WaitForState(); e != nil {
				return e
			}

			// wait unit network changing operation of instance done
			conf = tccommon.BuildStateChangeConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), time.Second, postgresqlService.PostgresqlDBInstanceStateRefreshFunc(instanceId, []string{}))
			if _, e := conf.WaitForState(); e != nil {
				return e
			}

			// refresh the private ip with new one
			_ = d.Set("private_access_ip", vipNew)
		}
	}

	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("postgres", "DBInstanceId", tcClient.Region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}

	if d.HasChange("deletion_protection") {
		request := postgresql.NewModifyDBInstanceDeletionProtectionRequest()
		request.DBInstanceId = &instanceId
		if v, ok := d.GetOkExists("deletion_protection"); ok {
			request.DeletionProtection = helper.Bool(v.(bool))
		}

		outErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().ModifyDBInstanceDeletionProtection(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if outErr != nil {
			return outErr
		}
	}

	return resourceTencentCloudPostgresqlCloneDbInstanceRead(d, meta)
}

func resourceTencentCloudPostgresqlCloneDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_clone_db_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		postgresqlService = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		outErr, inErr     error
		has               bool
	)

	instanceId := d.Id()
	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, has, inErr = postgresqlService.DescribePostgresqlInstanceById(ctx, d.Id())
		if inErr != nil {
			// ResourceNotFound.InstanceNotFoundError
			ee, ok := inErr.(*sdkErrors.TencentCloudSDKError)
			if ok && ee.GetCode() == "ResourceNotFound.InstanceNotFoundError" {
				return nil
			}

			return tccommon.RetryError(inErr, postgresql.OPERATIONDENIED_INSTANCESTATUSLIMITOPERROR)
		}

		return nil
	})

	if outErr != nil {
		return outErr
	}

	if !has {
		return nil
	}

	outErr = postgresqlService.IsolatePostgresqlInstance(ctx, instanceId)
	if outErr != nil {
		return outErr
	}

	// Wait for status to isolated
	_ = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		instance, _, err := postgresqlService.DescribePostgresqlInstanceById(ctx, instanceId)
		if err != nil {
			return tccommon.RetryError(err)
		}

		if *instance.DBInstanceStatus == POSTGRESQL_STAUTS_ISOLATED {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("waiting for instance isolating"))
	})

	outErr = postgresqlService.DeletePostgresqlInstance(ctx, instanceId)
	if outErr != nil {
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = postgresqlService.DeletePostgresqlInstance(ctx, instanceId)
			if inErr != nil {
				// ResourceNotFound.InstanceNotFoundError
				ee, ok := inErr.(*sdkErrors.TencentCloudSDKError)
				if ok && ee.GetCode() == "ResourceNotFound.InstanceNotFoundError" {
					return nil
				}

				return tccommon.RetryError(inErr)
			}

			return nil
		})
	}

	if outErr != nil {
		return outErr
	}

	outErr = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, has, inErr = postgresqlService.DescribePostgresqlInstanceById(ctx, d.Id())
		if inErr != nil {
			// ResourceNotFound.InstanceNotFoundError
			ee, ok := inErr.(*sdkErrors.TencentCloudSDKError)
			if ok && ee.GetCode() == "ResourceNotFound.InstanceNotFoundError" {
				return nil
			}

			return tccommon.RetryError(inErr)
		}

		if has {
			inErr = fmt.Errorf("delete postgresql instance %s fail, instance still exists from SDK DescribePostgresqlInstanceById", instanceId)
			return resource.RetryableError(inErr)
		}

		return nil
	})

	if outErr != nil {
		return outErr
	}

	return nil
}

func resourcePostgresqlCloneDbInstanceCreateStateRefreshFunc_0_0(ctx context.Context, dBInstanceId string) resource.StateRefreshFunc {
	var req *postgresv20170312.DescribeDBInstanceAttributeRequest
	return func() (interface{}, string, error) {
		meta := tccommon.ProviderMetaFromContext(ctx)
		if meta == nil {
			return nil, "", fmt.Errorf("resource data can not be nil")
		}

		if req == nil {
			d := tccommon.ResourceDataFromContext(ctx)
			if d == nil {
				return nil, "", fmt.Errorf("resource data can not be nil")
			}

			_ = d
			req = postgresv20170312.NewDescribeDBInstanceAttributeRequest()
			req.DBInstanceId = helper.String(dBInstanceId)
		}

		resp, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresV20170312Client().DescribeDBInstanceAttributeWithContext(ctx, req)
		if err != nil {
			return nil, "", err
		}

		if resp == nil || resp.Response == nil {
			return nil, "", nil
		}

		state := fmt.Sprintf("%v", *resp.Response.DBInstance.DBInstanceStatus)
		return resp.Response, state, nil
	}
}
