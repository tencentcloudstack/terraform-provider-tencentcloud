package postgresql

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccrs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/crs"

	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresqlReadonlyInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlReadOnlyInstanceCreate,
		Read:   resourceTencentCloudPostgresqlReadOnlyInstanceRead,
		Update: resourceTencentCloudPostgresqlReadOnlyInstanceUpdate,
		Delete: resourceTencentCLoudPostgresqlReadOnlyInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"db_version": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "PostgreSQL kernel version, which must be the same as that of the primary instance.",
			},
			"storage": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance storage capacity in GB.",
			},
			"memory": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Memory size(in GB). Allowed value must be larger than `memory` that data source `tencentcloud_postgresql_specinfos` provides.",
			},
			"master_db_instance_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "ID of the primary instance to which the read-only replica belongs.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The instance ID of this readonly resource.",
			},
			"zone": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Availability zone ID, which can be obtained through the Zone field in the returned value of the DescribeZones API.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Project ID.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "VPC ID.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC subnet ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance name.",
			},
			"security_groups_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "ID of security group.",
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      COMMON_PAYTYPE_POSTPAID,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(POSTGRESQL_PAYTYPE),
				Description:  "instance billing mode. Valid values: PREPAID (monthly subscription), POSTPAID_BY_HOUR (pay-as-you-go).",
			},
			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specify Prepaid period in month. Default `1`. Values: `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.",
			},
			"auto_renew_flag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Auto renew flag, `1` for enabled. NOTES: Only support prepaid instance.",
			},
			"auto_voucher": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether to use voucher, `1` for enabled.",
			},
			"voucher_ids": {
				Type:         schema.TypeList,
				Optional:     true,
				RequiredWith: []string{"auto_voucher"},
				Description:  "Specify Voucher Ids if `auto_voucher` was `1`, only support using 1 vouchers for now.",
				Elem:         &schema.Schema{Type: schema.TypeString},
			},
			"need_support_ipv6": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to support IPv6 address access. Valid values: 1 (yes), 0 (no).",
			},
			//"tag_list": {
			//	Type:        schema.TypeMap,
			//	Optional:    true,
			//	Description: "The information of tags to be associated with instances. This parameter is left empty by default..",
			//},
			"read_only_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "RO group ID.",
			},
			// Computed values
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the postgresql instance.",
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

func resourceTencentCloudPostgresqlReadOnlyInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_instance.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request           = postgresql.NewCreateReadOnlyDBInstanceRequest()
		response          *postgresql.CreateReadOnlyDBInstanceResponse
		postgresqlService = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		zone              string
		dbVersion         string
		memory            int
	)
	if v, ok := d.GetOk("db_version"); ok {
		dbVersion = v.(string)
		request.DBVersion = helper.String(dbVersion)
	}
	if v, ok := d.GetOk("storage"); ok {
		request.Storage = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("memory"); ok {
		memory = v.(int)
	}
	if v, ok := d.GetOk("master_db_instance_id"); ok {
		request.MasterDBInstanceId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("zone"); ok {
		zone = v.(string)
		request.Zone = helper.String(zone)
	}
	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("auto_renew_flag"); ok {
		request.AutoRenewFlag = helper.IntInt64(v.(int))
	}
	if v, ok := d.Get("period").(int); ok && v > 0 {
		request.Period = helper.IntUint64(v)
	}
	if v, ok := d.Get("auto_voucher").(int); ok && v > 0 {
		request.AutoVoucher = helper.IntUint64(v)
	}
	if v, ok := d.GetOk("voucher_ids"); ok {
		request.VoucherIds = helper.InterfacesStringsPoint(v.([]interface{}))
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
	if v, ok := d.GetOk("need_support_ipv6"); ok {
		request.NeedSupportIpv6 = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("read_only_group_id"); ok {
		request.ReadOnlyGroupId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("security_groups_ids"); ok {
		securityGroupsIds := v.(*schema.Set).List()
		request.SecurityGroupIds = make([]*string, 0, len(securityGroupsIds))
		for _, item := range securityGroupsIds {
			request.SecurityGroupIds = append(request.SecurityGroupIds, helper.String(item.(string)))
		}
	}
	//if tags := helper.GetTags(d, "tag_list"); len(tags) > 0 {
	//	for k, v := range tags {
	//		request.TagList = &postgresql.Tag{
	//			TagKey:   &k,
	//			TagValue: &v,
	//		}
	//	}
	//}

	// get specCode with db_version and memory
	var allowVersion, allowMemory []string
	var specVersion, specCode string
	err := resource.Retry(tccommon.ReadRetryTimeout*5, func() *resource.RetryError {
		speccodes, inErr := postgresqlService.DescribeSpecinfos(ctx, zone)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		for _, info := range speccodes {
			if !tccommon.IsContains(allowVersion, *info.Version) {
				allowVersion = append(allowVersion, *info.Version)
			}
			if *info.Version == dbVersion {
				specVersion = *info.Version
				memoryString := fmt.Sprintf("%d", int(*info.Memory)/1024)
				if !tccommon.IsContains(allowMemory, memoryString) {
					allowMemory = append(allowMemory, memoryString)
				}
				if int(*info.Memory)/1024 == memory {
					specCode = *info.SpecCode
					break
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if specVersion == "" {
		return fmt.Errorf(`The "db_version" value: "%s" is invalid, Valid values are one of: "%s"`, dbVersion, strings.Join(allowVersion, `", "`))
	}
	if specCode == "" {
		return fmt.Errorf(`The "storage" value: %d is invalid, Valid values are one of: %s`, memory, strings.Join(allowMemory, `, `))
	}
	request.SpecCode = helper.String(specCode)

	request.InstanceCount = helper.IntUint64(1)
	request.Period = helper.IntUint64(1)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().CreateReadOnlyDBInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil

	})
	if err != nil {
		return err
	}
	instanceId := *response.Response.DBInstanceIdSet[0]
	d.SetId(instanceId)

	// check creation done
	err = resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, has, err := postgresqlService.DescribePostgresqlInstanceById(ctx, instanceId)
		if err != nil {
			return tccommon.RetryError(err)
		} else if has && *instance.DBInstanceStatus == "running" {
			return nil
		} else if !has {
			return resource.NonRetryableError(fmt.Errorf("create postgresql instance fail"))
		} else {
			return resource.RetryableError(fmt.Errorf("creating readonly postgresql instance %s , status %s ", instanceId, *instance.DBInstanceStatus))
		}
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudPostgresqlReadOnlyInstanceRead(d, meta)
}

func resourceTencentCloudPostgresqlReadOnlyInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_instance.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()
	postgresqlService := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	instance, has, err := postgresqlService.DescribePostgresqlInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("db_version", instance.DBVersion)
	_ = d.Set("storage", instance.DBInstanceStorage)
	_ = d.Set("memory", instance.DBInstanceMemory)
	_ = d.Set("master_db_instance_id", instance.MasterDBInstanceId)
	_ = d.Set("zone", instance.Zone)
	_ = d.Set("project_id", instance.ProjectId)

	if *instance.PayType == POSTGRESQL_PAYTYPE_PREPAID || *instance.PayType == COMMON_PAYTYPE_PREPAID {
		_ = d.Set("instance_charge_type", COMMON_PAYTYPE_PREPAID)
	} else {
		_ = d.Set("instance_charge_type", COMMON_PAYTYPE_POSTPAID)
	}

	_ = d.Set("auto_renew_flag", instance.AutoRenew)
	_ = d.Set("vpc_id", instance.VpcId)
	_ = d.Set("subnet_id", instance.SubnetId)
	_ = d.Set("name", instance.DBInstanceName)
	_ = d.Set("need_support_ipv6", instance.SupportIpv6)
	// set readonly group when DescribeReadOnlyGroups ready for filter by the readonly group id
	// _ = d.Set("read_only_group_id", readonlyGroup.Id)

	// security groups
	// Only redis service support modify Generic DB instance security groups
	redisService := svccrs.NewRedisService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	sg, err := redisService.DescribeDBSecurityGroups(ctx, "postgres", d.Id())
	if err != nil {
		return err
	}
	if len(sg) > 0 {
		_ = d.Set("security_groups_ids", sg)
	}

	//tags := make(map[string]string, len(instance.TagList))
	//for _, tag := range instance.TagList {
	//	tags[*tag.TagKey] = *tag.TagValue
	//}
	//_ = d.Set("tag_list", tags)

	// computed
	_ = d.Set("create_time", instance.CreateTime)

	if len(instance.DBInstanceNetInfo) > 0 {
		for _, v := range instance.DBInstanceNetInfo {
			// private or inner will not appear at same time, private for instance with vpc
			if (*v.NetType == "private" || *v.NetType == "inner") && *v.Ip != "" {
				_ = d.Set("private_access_ip", v.Ip)
				_ = d.Set("private_access_port", v.Port)
			}
		}
	}

	return nil
}

func resourceTencentCloudPostgresqlReadOnlyInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_instance.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	postgresqlService := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	instanceId := d.Id()
	d.Partial(true)

	if err := helper.ImmutableArgsChek(d,
		"charge_type",
		"period",
		"auto_renew_flag",
		"auto_voucher",
		"voucher_ids",
	); err != nil {
		return err
	}

	if d.HasChange("read_only_group_id") {
		var (
			masterInstanceId string
			roGroupIdOld     string
			roGroupIdNew     string
			request          = postgresql.NewModifyDBInstanceReadOnlyGroupRequest()
			service          = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		)

		masterInstanceId = d.Get("master_db_instance_id").(string)
		old, new := d.GetChange("read_only_group_id")
		if old != nil {
			roGroupIdOld = old.(string)
		}
		if new != nil {
			roGroupIdNew = new.(string)
		}

		request.DBInstanceId = &instanceId
		request.ReadOnlyGroupId = &roGroupIdOld
		request.NewReadOnlyGroupId = &roGroupIdNew

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().ModifyDBInstanceReadOnlyGroup(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate postgresql ChangeDbInstanceReadOnlyGroupOperation failed, reason:%+v", logId, err)
			return err
		}

		conf := tccommon.BuildStateChangeConf([]string{}, []string{"ok"}, 2*tccommon.ReadRetryTimeout, time.Second, service.PostgresqlReadonlyGroupStateRefreshFunc(masterInstanceId, roGroupIdNew, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	var outErr, inErr, checkErr error
	// update name
	if d.HasChange("name") {
		name := d.Get("name").(string)
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = postgresqlService.ModifyPostgresqlInstanceName(ctx, instanceId, name)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}
		// check update name done
		checkErr = postgresqlService.CheckDBInstanceStatus(ctx, instanceId)
		if checkErr != nil {
			return checkErr
		}

	}

	// upgrade storage and memory size
	if d.HasChange("memory") || d.HasChange("storage") {
		memory := d.Get("memory").(int)
		storage := d.Get("storage").(int)
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = postgresqlService.UpgradePostgresqlInstance(ctx, instanceId, memory, storage)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}
		// check update storage and memory done
		checkErr = postgresqlService.CheckDBInstanceStatus(ctx, instanceId)
		if checkErr != nil {
			return checkErr
		}

	}

	// update project id
	if d.HasChange("project_id") {
		projectId := d.Get("project_id").(int)
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = postgresqlService.ModifyPostgresqlInstanceProjectId(ctx, instanceId, projectId)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

		// check update project id done
		checkErr = postgresqlService.CheckDBInstanceStatus(ctx, instanceId)
		if checkErr != nil {
			return checkErr
		}

	}

	if d.HasChange("security_groups_ids") {

		// Only redis service support modify Generic DB instance security groups
		service := svccrs.NewRedisService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		ids := d.Get("security_groups_ids").(*schema.Set).List()
		var sgIds []*string
		for _, id := range ids {
			sgIds = append(sgIds, helper.String(id.(string)))
		}
		err := service.ModifyDBInstanceSecurityGroups(ctx, "postgres", d.Id(), sgIds)
		if err != nil {
			return err
		}

	}

	//if d.HasChange("tags") {
	//
	//	oldValue, newValue := d.GetChange("tags")
	//	replaceTags, deleteTags := diffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))
	//
	//	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	//	tagService := svctag.NewTagService(tcClient)
	//	resourceName := tccommon.BuildTagResourceName("postgres", "DBInstanceId", tcClient.Region, d.Id())
	//	err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
	//	if err != nil {
	//		return err
	//	}
	//
	//}

	d.Partial(false)

	return resourceTencentCloudPostgresqlReadOnlyInstanceRead(d, meta)
}

func resourceTencentCLoudPostgresqlReadOnlyInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_instance.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()
	postgresqlService := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	// isolate
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := postgresqlService.IsolatePostgresqlInstance(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Wait for status to isolated
	_ = resource.Retry(tccommon.ReadRetryTimeout*5, func() *resource.RetryError {
		instance, _, err := postgresqlService.DescribePostgresqlInstanceById(ctx, instanceId)
		if err != nil {
			return tccommon.RetryError(err)
		}

		if *instance.DBInstanceStatus == POSTGRESQL_STAUTS_ISOLATED {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("waiting for readonly instance isolating"))
	})

	// delete
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := postgresqlService.DeletePostgresqlInstance(ctx, instanceId)
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
