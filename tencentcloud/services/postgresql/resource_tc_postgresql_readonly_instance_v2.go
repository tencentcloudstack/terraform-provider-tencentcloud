package postgresql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
)

func ResourceTencentCloudPostgresqlReadonlyInstanceV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlReadonlyInstanceV2Create,
		Read:   resourceTencentCloudPostgresqlReadonlyInstanceV2Read,
		Update: resourceTencentCloudPostgresqlReadonlyInstanceV2Update,
		Delete: resourceTencentCloudPostgresqlReadonlyInstanceV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Availability zone ID, such as ap-guangzhou-3.",
			},
			"master_db_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the primary instance to which the read-only replica belongs.",
			},
			"spec_code": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specification code, which can be obtained via DescribeClasses.",
			},
			"storage": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Instance storage capacity in GB, the step is 10.",
			},
			"instance_count": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Number of instances to purchase, value range: [1-6]. Only the first instance ID is managed by this resource.",
			},
			"period": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Purchase duration in months. PREPAID supports 1,2,3,4,5,6,7,8,9,10,11,12,24,36; POSTPAID only supports 1.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VPC ID, such as vpc-xxxxxxxx.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VPC subnet ID, such as subnet-xxxxxxxx.",
			},
			"instance_charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance billing mode. Valid values: PREPAID, POSTPAID_BY_HOUR. Default: PREPAID.",
			},
			"auto_voucher": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether to use voucher automatically, 1 for yes, 0 for no. Default: 0.",
			},
			"voucher_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Voucher ID list, currently only one voucher is supported.",
			},
			"auto_renew_flag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Auto renew flag, 0 for manual renew, 1 for auto renew. Default: 0. Only supports PREPAID.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Project ID. Default: 0, means default project.",
			},
			"activity_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Activity ID.",
			},
			"read_only_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Read-only group ID.",
			},
			"tag_list": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag value.",
						},
					},
				},
				Description: "Instance tag info (legacy, single tag). It is recommended to use the new field `tags`.",
			},
			"security_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Security group IDs bound to the instance.",
			},
			"need_support_ipv6": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether to support IPv6 access, 1 for yes, 0 for no. Default: 0.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance name, only supports chinese/english/numbers/_/- with length less than 60.",
			},
			"db_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "PostgreSQL kernel version, no longer needed, it will keep the same as the primary instance.",
			},
			"dedicated_cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Dedicated cluster ID.",
			},
			"deletion_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable deletion protection, true for enable, false for disable.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Instance tags.",
			},
			"deal_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Order number list, each instance corresponds to one order.",
			},
			"bill_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Frozen flow number.",
			},
			"db_instance_id_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Created instance ID set, only returned in POSTPAID scenario.",
			},
			"billing_parameters": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Billing parameters for order placement, only returned when the input parameter BillingParameters has a value.",
			},
			"db_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance ID managed by this resource.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlReadonlyInstanceV2Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_instance_v2.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := postgresql.NewCreateReadOnlyDBInstanceRequest()
	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}
	if v, ok := d.GetOk("master_db_instance_id"); ok {
		request.MasterDBInstanceId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("spec_code"); ok {
		request.SpecCode = helper.String(v.(string))
	}
	if v, ok := d.GetOk("storage"); ok {
		request.Storage = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("instance_count"); ok {
		request.InstanceCount = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("period"); ok {
		request.Period = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("auto_voucher"); ok {
		request.AutoVoucher = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("voucher_ids"); ok {
		request.VoucherIds = helper.InterfacesStringsPoint(v.([]interface{}))
	}
	if v, ok := d.GetOk("auto_renew_flag"); ok {
		request.AutoRenewFlag = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("activity_id"); ok {
		request.ActivityId = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("read_only_group_id"); ok {
		request.ReadOnlyGroupId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("tag_list"); ok {
		tagList := v.([]interface{})
		if len(tagList) > 0 {
			tagMap := tagList[0].(map[string]interface{})
			tag := &postgresql.Tag{}
			if key, ok := tagMap["tag_key"]; ok {
				tag.TagKey = helper.String(key.(string))
			}
			if value, ok := tagMap["tag_value"]; ok {
				tag.TagValue = helper.String(value.(string))
			}
			request.TagList = tag
		}
	}
	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIds := v.([]interface{})
		request.SecurityGroupIds = make([]*string, 0, len(securityGroupIds))
		for _, item := range securityGroupIds {
			request.SecurityGroupIds = append(request.SecurityGroupIds, helper.String(item.(string)))
		}
	}
	if v, ok := d.GetOk("need_support_ipv6"); ok {
		request.NeedSupportIpv6 = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}
	if v, ok := d.GetOk("db_version"); ok {
		request.DBVersion = helper.String(v.(string))
	}
	if v, ok := d.GetOk("dedicated_cluster_id"); ok {
		request.DedicatedClusterId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("deletion_protection"); ok {
		request.DeletionProtection = helper.Bool(v.(bool))
	}

	postgresqlService := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	response, err := postgresqlService.CreatePostgresqlReadonlyInstanceV2(ctx, request)
	if err != nil {
		return err
	}
	if response == nil || response.Response == nil {
		return fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}
	if len(response.Response.DBInstanceIdSet) == 0 {
		log.Printf("[CRITAL]%s create postgresql_readonly_instance_v2 id=%s, DBInstanceIdSet is empty", logId, d.Id())
		return fmt.Errorf("create postgresql_readonly_instance_v2 failed, DBInstanceIdSet is empty")
	}
	instanceId := *response.Response.DBInstanceIdSet[0]
	d.SetId(instanceId)

	// wait for instance running
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		instance, has, e := postgresqlService.DescribePostgresqlReadonlyInstanceV2ById(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if !has {
			return resource.RetryableError(fmt.Errorf("creating postgresql_readonly_instance_v2 %s, instance not found", instanceId))
		}
		if instance == nil || instance.DBInstanceStatus == nil {
			return resource.RetryableError(fmt.Errorf("creating postgresql_readonly_instance_v2 %s, status is nil", instanceId))
		}
		if *instance.DBInstanceStatus == POSTGRESQL_STAUTS_RUNNING {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("creating postgresql_readonly_instance_v2 %s, status %s", instanceId, *instance.DBInstanceStatus))
	})
	if err != nil {
		return err
	}

	// set computed fields from create response
	if len(response.Response.DealNames) > 0 {
		dealNames := make([]interface{}, 0, len(response.Response.DealNames))
		for _, deal := range response.Response.DealNames {
			if deal != nil {
				dealNames = append(dealNames, *deal)
			}
		}
		_ = d.Set("deal_names", dealNames)
	}
	if response.Response.BillId != nil {
		_ = d.Set("bill_id", *response.Response.BillId)
	}
	if len(response.Response.DBInstanceIdSet) > 0 {
		dbInstanceIdSet := make([]interface{}, 0, len(response.Response.DBInstanceIdSet))
		for _, id := range response.Response.DBInstanceIdSet {
			if id != nil {
				dbInstanceIdSet = append(dbInstanceIdSet, *id)
			}
		}
		_ = d.Set("db_instance_id_set", dbInstanceIdSet)
	}
	if response.Response.BillingParameters != nil {
		_ = d.Set("billing_parameters", *response.Response.BillingParameters)
	}
	_ = d.Set("db_instance_id", instanceId)

	// handle tags via tag service if provided
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("postgres", "DBInstanceId", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudPostgresqlReadonlyInstanceV2Read(d, meta)
}

func resourceTencentCloudPostgresqlReadonlyInstanceV2Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_instance_v2.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()
	postgresqlService := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var instance *postgresql.DBInstance
	var has bool
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		inst, h, e := postgresqlService.DescribePostgresqlReadonlyInstanceV2ById(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instance = inst
		has = h
		return nil
	})
	if err != nil {
		return err
	}
	if !has || instance == nil {
		log.Printf("[CRUD] postgresql_readonly_instance_v2 id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if instance.Zone != nil {
		_ = d.Set("zone", *instance.Zone)
	}
	if instance.VpcId != nil {
		_ = d.Set("vpc_id", *instance.VpcId)
	}
	if instance.SubnetId != nil {
		_ = d.Set("subnet_id", *instance.SubnetId)
	}
	if instance.DBVersion != nil {
		_ = d.Set("db_version", *instance.DBVersion)
	}
	if instance.DBInstanceStorage != nil {
		_ = d.Set("storage", *instance.DBInstanceStorage)
	}
	if instance.DBInstanceName != nil {
		_ = d.Set("name", *instance.DBInstanceName)
	}
	if instance.ProjectId != nil {
		_ = d.Set("project_id", *instance.ProjectId)
	}
	if instance.SupportIpv6 != nil {
		_ = d.Set("need_support_ipv6", *instance.SupportIpv6)
	}
	if instance.DeletionProtection != nil {
		_ = d.Set("deletion_protection", *instance.DeletionProtection)
	}
	if instance.MasterDBInstanceId != nil {
		_ = d.Set("master_db_instance_id", *instance.MasterDBInstanceId)
	}
	if instance.AutoRenew != nil {
		_ = d.Set("auto_renew_flag", *instance.AutoRenew)
	}
	if instance.PayType != nil {
		if *instance.PayType == "prepaid" {
			_ = d.Set("instance_charge_type", "PREPAID")
		} else {
			_ = d.Set("instance_charge_type", "POSTPAID_BY_HOUR")
		}
	}
	if len(instance.DBNodeSet) > 0 {
		for _, node := range instance.DBNodeSet {
			if node != nil && node.DedicatedClusterId != nil {
				_ = d.Set("dedicated_cluster_id", *node.DedicatedClusterId)
				break
			}
		}
	}

	// read only group
	if instance.MasterDBInstanceId != nil {
		readOnlyGroupId, roErr := postgresqlService.DescribeReadOnlyGroupsById(ctx, *instance.MasterDBInstanceId, d.Id())
		if roErr != nil {
			return roErr
		}
		if readOnlyGroupId != nil {
			_ = d.Set("read_only_group_id", *readOnlyGroupId)
		}
	}

	// security groups
	sg, sgErr := postgresqlService.DescribeDBInstanceSecurityGroupsById(ctx, d.Id())
	if sgErr != nil {
		return sgErr
	}
	if len(sg) > 0 {
		_ = d.Set("security_group_ids", sg)
	}

	// tags via tag service
	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, tagErr := tagService.DescribeResourceTags(ctx, "postgres", "DBInstanceId", tcClient.Region, d.Id())
	if tagErr != nil {
		return tagErr
	}
	_ = d.Set("tags", tags)

	_ = d.Set("db_instance_id", d.Id())
	return nil
}

func resourceTencentCloudPostgresqlReadonlyInstanceV2Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_instance_v2.update")()

	if err := helper.ImmutableArgsChek(d,
		"vpc_id",
		"subnet_id",
		"instance_charge_type",
		"auto_voucher",
		"voucher_ids",
		"auto_renew_flag",
		"project_id",
		"activity_id",
		"read_only_group_id",
		"tag_list",
		"security_group_ids",
		"need_support_ipv6",
		"name",
		"db_version",
		"dedicated_cluster_id",
		"deletion_protection",
		"tags",
	); err != nil {
		return err
	}

	return resourceTencentCloudPostgresqlReadonlyInstanceV2Read(d, meta)
}

func resourceTencentCloudPostgresqlReadonlyInstanceV2Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_instance_v2.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()
	postgresqlService := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := postgresqlService.IsolatePostgresqlReadonlyInstanceV2(ctx, instanceId); err != nil {
		return err
	}

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		instance, has, e := postgresqlService.DescribePostgresqlReadonlyInstanceV2ById(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if !has || instance == nil {
			return nil
		}
		if instance.DBInstanceStatus != nil && *instance.DBInstanceStatus == POSTGRESQL_STAUTS_ISOLATED {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("waiting for postgresql_readonly_instance_v2 %s isolating", instanceId))
	})
	if err != nil {
		log.Printf("[CRITAL]%s wait postgresql_readonly_instance_v2 id=%s isolated fail, reason: %v", logId, instanceId, err)
		return err
	}
	return nil
}
