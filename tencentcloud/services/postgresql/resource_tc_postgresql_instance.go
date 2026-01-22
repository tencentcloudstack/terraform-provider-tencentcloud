package postgresql

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

//internal version: replace import begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
//internal version: replace import end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

func ResourceTencentCloudPostgresqlInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlInstanceCreate,
		Read:   resourceTencentCloudPostgresqlInstanceRead,
		Update: resourceTencentCloudPostgresqlInstanceUpdate,
		Delete: resourceTencentCLoudPostgresqlInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: helper.ImportWithDefaultValue(map[string]interface{}{
				"delete_protection": false,
			}),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "Name of the postgresql instance.",
			},
			"charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     COMMON_PAYTYPE_POSTPAID,
				Description: "Pay type of the postgresql instance. Values `POSTPAID_BY_HOUR` (Default), `PREPAID`. It only support to update the type from `POSTPAID_BY_HOUR` to `PREPAID`.",
			},
			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specify Prepaid period in month. Default `1`. Values: `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`. This field is valid only when creating a `PREPAID` type instance, or updating the charge type from `POSTPAID_BY_HOUR` to `PREPAID`.",
			},
			"auto_renew_flag": {
				Type:        schema.TypeInt,
				Optional:    true,
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
			"engine_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Version of the postgresql database engine. Valid values: `10.4`, `10.17`, `10.23`, `11.8`, `11.12`, `11.22`, `12.4`, `12.7`, `12.18`, `13.3`, `14.2`, `14.11`, `15.1`, `16.0`.",
			},
			"db_major_vesion": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "`db_major_vesion` will be deprecated, use `db_major_version` instead.",
				ConflictsWith: []string{"db_major_version"},
				Description: "PostgreSQL major version number. Valid values: 10, 11, 12, 13, 14, 15, 16. " +
					"If it is specified, an instance running the latest kernel of PostgreSQL DBMajorVersion will be created.",
			},
			"db_major_version": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"db_major_vesion"},
				Description: "PostgreSQL major version number. Valid values: 10, 11, 12, 13, 14, 15, 16. " +
					"If it is specified, an instance running the latest kernel of PostgreSQL DBMajorVersion will be created.",
			},
			"db_kernel_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "PostgreSQL kernel version number. " +
					"If it is specified, an instance running kernel DBKernelVersion will be created. It supports updating the minor kernel version immediately.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of VPC.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of subnet.",
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
			"storage": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Volume size(in GB). Allowed value must be a multiple of 10. The storage must be set with the limit of `storage_min` and `storage_max` which data source `tencentcloud_postgresql_specinfos` provides.",
			},
			"memory": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Memory size(in GB). Allowed value must be larger than `memory` that data source `tencentcloud_postgresql_specinfos` provides.",
			},
			"cpu": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of CPU cores. Allowed value must be equal `cpu` that data source `tencentcloud_postgresql_specinfos` provides.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Project id, default value is `0`.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Availability zone. NOTE: This field could not be modified, please use `db_node_set` instead of modification. The changes on this field will be suppressed when using the `db_node_set`.",
				DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
					if o == "" {
						return false
					}
					raw, ok := d.GetOk("db_node_set")
					if !ok {
						return n == o
					}
					nodeZones := raw.(*schema.Set).List()
					for i := range nodeZones {
						item := nodeZones[i].(map[string]interface{})
						if zone, ok := item["zone"].(string); ok && zone == n {
							return true
						}
					}
					return n == o
				},
			},
			"root_user": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "root",
				Description: "Instance root account name. This parameter is optional, Default value is `root`.",
			},
			"root_password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: tccommon.ValidateMysqlPassword,
				Description:  "Password of root account. This parameter can be specified when you purchase master instances, but it should be ignored when you purchase read-only instances or disaster recovery instances.",
			},
			"charset": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      POSTGRESQL_DB_CHARSET_UTF8,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(POSTGRESQL_DB_CHARSET),
				Description:  "Charset of the root account. Valid values are `UTF8`,`LATIN1`.",
			},
			"need_support_tde": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to support data transparent encryption, 1: yes, 0: no (default).",
			},
			"kms_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "KeyId of the custom key.",
			},
			"kms_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Region of the custom key.",
			},
			"kms_cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specify the cluster served by KMS. If KMSClusterId is blank, use the KMS of the default cluster. If you choose to specify a KMS cluster, you need to pass in KMSClusterId.",
			},
			"public_access_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether to enable the access to an instance from public network or not.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The available tags within this postgresql.",
			},
			"max_standby_archive_delay": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "max_standby_archive_delay applies when WAL data is being read from WAL archive (and is therefore not current). Units are milliseconds if not specified.",
			},
			"max_standby_streaming_delay": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "max_standby_streaming_delay applies when WAL data is being received via streaming replication. Units are milliseconds if not specified.",
			},
			"backup_plan": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Specify DB backup plan.",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_backup_start_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specify earliest backup start time, format `hh:mm:ss`.",
						},
						"max_backup_start_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specify latest backup start time, format `hh:mm:ss`.",
						},
						"base_backup_retention_period": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Specify days of the retention.",
						},
						"backup_period": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "List of backup period per week, available values: `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, `sunday`. NOTE: At least specify two days.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"monthly_backup_retention_period": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specify days of the retention.",
						},
						"monthly_backup_period": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "If it is in monthly dimension, the format is numeric characters, such as [\"1\",\"2\"].",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"monthly_plan_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Monthly plan id.",
						},
					},
				},
			},
			"db_node_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Specify instance node info for disaster migration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "Standby",
							Description: "Indicates node type, available values:`Primary`, `Standby`. Default: `Standby`.",
						},
						"zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Indicates the node available zone.",
						},
						"dedicated_cluster_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Dedicated cluster ID.",
						},
					},
				},
			},
			"delete_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to enable instance deletion protection. Default: false.",
			},
			"wait_switch": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{POSTGRESQL_KERNEL_UPGRADE_IMMEDIATELY, POSTGRESQL_KERNEL_UPGRADE_MAINTAIN_WINDOW}),
				Description:  "Switch time after instance configurations are modified. `0`: Switch immediately; `2`: Switch during maintenance time window. Default: `0`. Note: This only takes effect when updating the `memory`, `storage`, `cpu`, `db_node_set`, `db_kernel_version` fields.",
			},
			// Computed values
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
			"uid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Uid of the postgresql instance.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the postgresql instance.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_instance.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	//internal version: replace clientCreate begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	postgresqlService := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	//internal version: replace clientCreate end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

	var (
		name      = d.Get("name").(string)
		dbVersion = d.Get("engine_version").(string)
		payType   = d.Get("charge_type").(string)
		//internal version: replace var begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
		//internal version: replace var end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
		projectId      = d.Get("project_id").(int)
		subnetId       = d.Get("subnet_id").(string)
		vpcId          = d.Get("vpc_id").(string)
		securityGroups = d.Get("security_groups").(*schema.Set).List()
		zone           = d.Get("availability_zone").(string)
		storage        = d.Get("storage").(int)
		memory         = d.Get("memory").(int) // Memory only used for query specCode which contains memory info
		username       = d.Get("root_user").(string)
		password       = d.Get("root_password").(string)
		charset        = d.Get("charset").(string)
		nodeSet        = d.Get("db_node_set").(*schema.Set).List()
	)

	// the sdk asks to set value with 1 when paytype is postpaid

	var instanceId, majorVersion, specVersion, specCode string
	var outErr, inErr error
	var allowMajorVersion, allowSpecVersion, allowSpec []string

	var (
		dbMajorVersion  = ""
		dbKernelVersion = ""
		needSupportTde  = 0
		kmsKeyId        = ""
		kmsRegion       = ""
		kmsClusterId    = ""
		period          = 1
		autoRenewFlag   = 0
		autoVoucher     = 0
		voucherIds      []*string
		cpu             int // cpu only used for query specCode which contains cpu info
	)

	if v, ok := d.GetOkExists("cpu"); ok {
		cpu = v.(int)
	}

	if v, ok := d.GetOkExists("period"); ok {
		log.Printf("period set")
		period = v.(int)
	} else {
		log.Printf("period not set")
	}

	if v, ok := d.GetOk("db_major_vesion"); ok {
		dbMajorVersion = v.(string)
	}

	if v, ok := d.GetOk("db_major_version"); ok {
		dbMajorVersion = v.(string)
	}

	if v, ok := d.GetOk("db_kernel_version"); ok {
		dbKernelVersion = v.(string)
	}

	if v, ok := d.GetOkExists("need_support_tde"); ok {
		needSupportTde = v.(int)
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		kmsKeyId = v.(string)
	}

	if v, ok := d.GetOk("kms_region"); ok {
		kmsRegion = v.(string)
	}

	if v, ok := d.GetOk("kms_cluster_id"); ok {
		kmsClusterId = v.(string)
	}

	if v, ok := d.GetOkExists("auto_renew_flag"); ok {
		autoRenewFlag = v.(int)
	}

	if v, ok := d.GetOkExists("auto_voucher"); ok {
		autoVoucher = v.(int)
	}

	if v, ok := d.GetOk("voucher_ids"); ok {
		voucherIds = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	requestSecurityGroup := make([]string, 0, len(securityGroups))

	for _, v := range securityGroups {
		requestSecurityGroup = append(requestSecurityGroup, v.(string))
	}

	if dbVersion == "" && dbMajorVersion == "" && dbKernelVersion == "" {
		dbVersion = "10.4"
	}

	// get specCode with engine_version and memory
	outErr = resource.Retry(tccommon.ReadRetryTimeout*5, func() *resource.RetryError {
		speccodes, inErr := postgresqlService.DescribeSpecinfos(ctx, zone)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}

		for _, info := range speccodes {
			if !tccommon.IsContains(allowSpecVersion, *info.Version) {
				allowSpecVersion = append(allowSpecVersion, *info.Version)
			}

			if !tccommon.IsContains(allowMajorVersion, *info.MajorVersion) {
				allowMajorVersion = append(allowMajorVersion, *info.MajorVersion)
			}

			if *info.MajorVersion == dbMajorVersion || *info.Version == dbVersion {
				majorVersion = *info.MajorVersion
				specVersion = *info.Version
				specString := fmt.Sprintf("(%d, %d)", int(*info.Memory)/1024, int(*info.Cpu))
				if !tccommon.IsContains(allowSpec, specString) {
					allowSpec = append(allowSpec, specString)
				}

				if cpu != 0 && int(*info.Cpu) == cpu && int(*info.Memory)/1024 == memory {
					specCode = *info.SpecCode
					break
				}

				if cpu == 0 && int(*info.Memory)/1024 == memory {
					specCode = *info.SpecCode
					break
				}
			}
		}

		return nil
	})

	if outErr != nil {
		return outErr
	}

	if majorVersion == "" && specVersion == "" {
		return fmt.Errorf(`The "db_major_version" value: "%s" is invalid, Valid values are one of: "%s", The "engine_version" value: "%s" is invalid, Valid values are one of: "%s"`, dbMajorVersion, strings.Join(allowMajorVersion, `", "`), dbVersion, strings.Join(allowSpecVersion, `", "`))
	}

	if specCode == "" {
		return fmt.Errorf(`The "memory" value: %d or the "cpu" value: %d is invalid, Valid combine values are one of: %s .`,
			memory, cpu, strings.Join(allowSpec, `; `))
	}

	var dbNodeSet []*postgresql.DBNode
	if len(nodeSet) > 0 {
		for i := range nodeSet {
			var (
				item               = nodeSet[i].(map[string]interface{})
				role               = item["role"].(string)
				zone               = item["zone"].(string)
				dedicatedClusterId = item["dedicated_cluster_id"].(string)
				node               *postgresql.DBNode
			)

			if dedicatedClusterId != "" {
				node = &postgresql.DBNode{
					Role:               &role,
					Zone:               &zone,
					DedicatedClusterId: &dedicatedClusterId,
				}
			} else {
				node = &postgresql.DBNode{
					Role: &role,
					Zone: &zone,
				}
			}

			dbNodeSet = append(dbNodeSet, node)
		}

		// check if availability_zone and node_set consists
		if include, z, nzs := checkZoneSetInclude(d); !include {
			return fmt.Errorf("`availability_zone`: %s is not included in `db_node_set`: %s", z, nzs)
		}
	}

	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		instanceId, inErr = postgresqlService.CreatePostgresqlInstance(ctx,
			name,
			dbVersion,
			dbMajorVersion,
			dbKernelVersion,
			payType,
			specCode,
			autoRenewFlag,
			projectId,
			period,
			subnetId,
			vpcId,
			zone,
			requestSecurityGroup,
			storage,
			username,
			password,
			charset,
			dbNodeSet,
			needSupportTde,
			kmsKeyId,
			kmsRegion,
			kmsClusterId,
			autoVoucher,
			voucherIds,
		)

		if inErr != nil {
			//internal version: replace bpass begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
			//internal version: replace bpass end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
			return tccommon.RetryError(inErr)
		}

		return nil
	})

	if outErr != nil {
		return outErr
	}

	d.SetId(instanceId)

	//internal version: replace setTag begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	//internal version: replace setTag end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

	// check creation done
	err := resource.Retry(20*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, has, err := postgresqlService.DescribePostgresqlInstanceById(ctx, instanceId)
		if err != nil {
			return tccommon.RetryError(err)
		}

		if !has {
			return resource.NonRetryableError(fmt.Errorf("create postgresql instance fail"))
		}

		if *instance.DBInstanceStatus == POSTGRESQL_STAUTS_RUNNING {
			memory = int(*instance.DBInstanceMemory)
			return nil
		}

		return resource.RetryableError(fmt.Errorf("creating postgresql instance %s , status %s ", instanceId, *instance.DBInstanceStatus))
	})

	if err != nil {
		return err
	}

	// check init status
	checkErr := postgresqlService.CheckDBInstanceStatus(ctx, instanceId)
	if checkErr != nil {
		return checkErr
	}

	// set open public access
	public_access_switch := false
	if v, ok := d.GetOkExists("public_access_switch"); ok {
		public_access_switch = v.(bool)
	}

	if public_access_switch {
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = postgresqlService.ModifyPublicService(ctx, true, instanceId)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}

			return nil
		})

		if outErr != nil {
			return outErr
		}
	}

	// check creation done
	checkErr = postgresqlService.CheckDBInstanceStatus(ctx, instanceId)
	if checkErr != nil {
		return checkErr
	}

	// set name
	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		inErr := postgresqlService.ModifyPostgresqlInstanceName(ctx, instanceId, name)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}

		return nil
	})

	if outErr != nil {
		return outErr
	}

	// check creation done
	checkErr = postgresqlService.CheckDBInstanceStatus(ctx, instanceId)
	if checkErr != nil {
		return checkErr
	}

	//internal version: replace null begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("postgres", "DBInstanceId", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	//internal version: replace null end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

	// set pg params
	paramEntrys := make(map[string]string)
	if v, ok := d.GetOkExists("max_standby_archive_delay"); ok {
		paramEntrys["max_standby_archive_delay"] = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOkExists("max_standby_streaming_delay"); ok {
		paramEntrys["max_standby_streaming_delay"] = strconv.Itoa(v.(int))
	}

	if len(paramEntrys) != 0 {
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr := postgresqlService.ModifyPgParams(ctx, instanceId, paramEntrys)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}

			return nil
		})

		if outErr != nil {
			return outErr
		}

		// 10s is required to synchronize the data(`DescribeParamsEvent`).
		time.Sleep(10 * time.Second)
	}

	// set backup plan
	if plan, ok := helper.InterfacesHeadMap(d, "backup_plan"); ok {
		if v, ok := plan["monthly_backup_period"].([]interface{}); ok && len(v) > 0 {
			request0 := postgresql.NewCreateBackupPlanRequest()
			request0.DBInstanceId = &instanceId
			request0.BackupPeriodType = helper.String("month")
			request0.PlanName = helper.String("custom_month")
			request0.BackupPeriod = helper.InterfacesStringsPoint(v)

			if v, ok := plan["min_backup_start_time"].(string); ok && v != "" {
				request0.MinBackupStartTime = &v
			}

			if v, ok := plan["max_backup_start_time"].(string); ok && v != "" {
				request0.MaxBackupStartTime = &v
			}

			if v, ok := plan["monthly_backup_retention_period"].(int); ok && v != 0 {
				request0.BaseBackupRetentionPeriod = helper.IntUint64(v)
			}

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().CreateBackupPlan(request0)
				if e != nil {
					return tccommon.RetryError(err, postgresql.OPERATIONDENIED_INSTANCESTATUSLIMITOPERROR)
				}
				return nil
			})

			if err != nil {
				return err
			}
		}
		request := postgresql.NewModifyBackupPlanRequest()
		request.DBInstanceId = &instanceId
		if v, ok := plan["min_backup_start_time"].(string); ok && v != "" {
			request.MinBackupStartTime = &v
		}

		if v, ok := plan["max_backup_start_time"].(string); ok && v != "" {
			request.MaxBackupStartTime = &v
		}

		if v, ok := plan["base_backup_retention_period"].(int); ok && v != 0 {
			request.BaseBackupRetentionPeriod = helper.IntUint64(v)
		}

		if v, ok := plan["backup_period"].([]interface{}); ok && len(v) > 0 {
			request.BackupPeriod = helper.InterfacesStringsPoint(v)
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			err := postgresqlService.ModifyBackupPlan(ctx, request)
			if err != nil {
				return tccommon.RetryError(err, postgresql.OPERATIONDENIED_INSTANCESTATUSLIMITOPERROR)
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return resourceTencentCloudPostgresqlInstanceRead(d, meta)
}

func resourceTencentCloudPostgresqlInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_instance.read")()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		instance      *postgresql.DBInstance
		has           bool
		outErr, inErr error
	)

	// Check if import
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

	// rootUser
	accounts, outErr := postgresqlService.DescribeRootUser(ctx, d.Id())
	if outErr != nil {
		return outErr
	}

	var rootUser string
	if len(accounts) > 0 {
		rootUser = *accounts[0].UserName
	} else if username, ok := d.GetOk("root_user"); ok {
		rootUser = username.(string)
	}

	_ = d.Set("project_id", int(*instance.ProjectId))
	_ = d.Set("availability_zone", instance.Zone)
	_ = d.Set("engine_version", instance.DBVersion)
	_ = d.Set("db_kernel_version", instance.DBKernelVersion)
	_ = d.Set("db_major_vesion", instance.DBMajorVersion)
	_ = d.Set("db_major_version", instance.DBMajorVersion)
	_ = d.Set("name", instance.DBInstanceName)
	_ = d.Set("charset", instance.DBCharset)

	// check net num
	if len(instance.DBInstanceNetInfo) == 3 {
		_ = d.Set("vpc_id", instance.DBInstanceNetInfo[0].VpcId)
		_ = d.Set("subnet_id", instance.DBInstanceNetInfo[0].SubnetId)
		_ = d.Set("private_access_ip", instance.DBInstanceNetInfo[0].Ip)
		_ = d.Set("private_access_port", instance.DBInstanceNetInfo[0].Port)

		// net status
		public_access_switch := false
		for _, v := range instance.DBInstanceNetInfo {
			if *v.NetType == "public" {
				// both 1 and opened used in SDK
				if *v.Status == "opened" || *v.Status == "1" {
					public_access_switch = true
				}

				_ = d.Set("public_access_host", v.Address)
				_ = d.Set("public_access_port", v.Port)
			}
		}

		_ = d.Set("public_access_switch", public_access_switch)
	} else if len(instance.DBInstanceNetInfo) == 2 {
		_ = d.Set("vpc_id", instance.VpcId)
		_ = d.Set("subnet_id", instance.SubnetId)

		// net status
		public_access_switch := false
		for _, v := range instance.DBInstanceNetInfo {
			if *v.NetType == "public" {
				// both 1 and opened used in SDK
				if *v.Status == "opened" || *v.Status == "1" {
					public_access_switch = true
				}

				_ = d.Set("public_access_host", v.Address)
				_ = d.Set("public_access_port", v.Port)
			}

			// private or inner will not appear at same time, private for instance with vpc
			if (*v.NetType == "private" || *v.NetType == "inner") && *v.Ip != "" {
				_ = d.Set("private_access_ip", v.Ip)
				_ = d.Set("private_access_port", v.Port)
			}
		}

		_ = d.Set("public_access_switch", public_access_switch)
	} else {
		return fmt.Errorf("DBInstanceNetInfo returned incorrect information.")
	}

	if rootUser != "" {
		_ = d.Set("root_user", &rootUser)
	}

	if *instance.PayType == POSTGRESQL_PAYTYPE_PREPAID || *instance.PayType == COMMON_PAYTYPE_PREPAID {
		_ = d.Set("charge_type", COMMON_PAYTYPE_PREPAID)
	} else {
		_ = d.Set("charge_type", COMMON_PAYTYPE_POSTPAID)
	}

	// security groups
	sg, err := postgresqlService.DescribeDBInstanceSecurityGroupsById(ctx, d.Id())
	if err != nil {
		return err
	}

	if len(sg) > 0 {
		_ = d.Set("security_groups", sg)
	}

	attrRequest := postgresql.NewDescribeDBInstanceAttributeRequest()
	attrRequest.DBInstanceId = helper.String(d.Id())
	ins, err := postgresqlService.DescribeDBInstanceAttribute(ctx, attrRequest)
	if err != nil {
		return err
	}

	nodeSet := ins.DBNodeSet
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

	// computed
	_ = d.Set("create_time", instance.CreateTime)
	_ = d.Set("memory", instance.DBInstanceMemory)
	_ = d.Set("storage", instance.DBInstanceStorage)
	_ = d.Set("cpu", instance.DBInstanceCpu)

	// kms
	kmsRequest := postgresql.NewDescribeEncryptionKeysRequest()
	kmsRequest.DBInstanceId = helper.String(d.Id())
	_ = d.Set("need_support_tde", instance.IsSupportTDE)
	has, kms, err := postgresqlService.DescribeDBEncryptionKeys(ctx, kmsRequest)
	if err != nil {
		return err
	}

	if has {
		_ = d.Set("kms_key_id", kms.KeyId)
		_ = d.Set("kms_region", kms.KeyRegion)
		_ = d.Set("kms_cluster_id", kms.KMSClusterId)
	}

	// Uid, must use
	var filters = make([]*postgresql.Filter, 0, 10)
	idFilter := &postgresql.Filter{
		Name:   helper.String("db-instance-id"),
		Values: []*string{helper.String(d.Id())},
	}

	filters = append(filters, idFilter)
	instanceList, err := postgresqlService.DescribePostgresqlInstances(ctx, filters)
	if err != nil {
		log.Printf("[CRITAL]%s query postgreSql  error, reason:%s\n", logId, err.Error())
		return err
	}

	if len(instanceList) == 0 {
		return fmt.Errorf("no postgresql instance find by id: %s\n.", d.Id())
	}

	if len(instanceList) > 1 {
		return fmt.Errorf("find more than one postgresql instance find by id: %s\n.", d.Id())
	}

	_ = d.Set("uid", instanceList[0].Uid)

	// ignore spec_code
	// qcs::postgres:ap-guangzhou:uin/123435236:DBInstanceId/postgres-xxx
	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "postgres", "DBInstanceId", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	// backup plans (only specified will rewrite)
	bkpRequest := postgresql.NewDescribeBackupPlansRequest()
	bkpRequest.DBInstanceId = helper.String(d.Id())
	bkpResponse, err := postgresqlService.DescribeBackupPlans(ctx, bkpRequest)
	if err != nil {
		return err
	}

	var backupPlan, monthlyBackupPlan *postgresql.BackupPlan
	if len(bkpResponse) > 0 {
		backupPlan = bkpResponse[0]
		for _, plan := range bkpResponse {
			if plan != nil && plan.BackupPeriodType != nil {
				if *plan.BackupPeriodType == "month" {
					monthlyBackupPlan = plan
				}
				if *plan.BackupPeriodType == "week" {
					backupPlan = plan
				}
			}
		}
	}

	if backupPlan != nil {
		planMap := map[string]interface{}{}
		if backupPlan.MinBackupStartTime != nil {
			planMap["min_backup_start_time"] = backupPlan.MinBackupStartTime
		}

		if backupPlan.MaxBackupStartTime != nil {
			planMap["max_backup_start_time"] = backupPlan.MaxBackupStartTime
		}

		if backupPlan.BaseBackupRetentionPeriod != nil {
			planMap["base_backup_retention_period"] = backupPlan.BaseBackupRetentionPeriod
		}

		if backupPlan.BackupPeriod != nil {
			strSlice := []string{}
			// set period list from BackupPeriods string, eg:"BackupPeriod": "[\"tuesday\",\"wednesday\"]",
			err := json.Unmarshal([]byte(*backupPlan.BackupPeriod), &strSlice)
			if err != nil {
				return fmt.Errorf("BackupPeriod:[%s] has invalid format,Unmarshal failed! error: %v", *backupPlan.BackupPeriod, err.Error())
			}

			planMap["backup_period"] = strSlice
		}

		if monthlyBackupPlan != nil && monthlyBackupPlan.PlanId != nil {
			planMap["monthly_plan_id"] = monthlyBackupPlan.PlanId
		}
		if monthlyBackupPlan != nil && monthlyBackupPlan.BackupPeriod != nil {
			strSlice := []string{}
			err := json.Unmarshal([]byte(*monthlyBackupPlan.BackupPeriod), &strSlice)
			if err != nil {
				return fmt.Errorf("BackupPeriod:[%s] has invalid format,Unmarshal failed! error: %v", *backupPlan.BackupPeriod, err.Error())
			}

			planMap["monthly_backup_period"] = strSlice
		}
		if monthlyBackupPlan != nil && monthlyBackupPlan.BaseBackupRetentionPeriod != nil {
			planMap["monthly_backup_retention_period"] = monthlyBackupPlan.BaseBackupRetentionPeriod
		}

		_ = d.Set("backup_plan", []interface{}{planMap})
	}

	// pg params
	var parmas map[string]string
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		parmas, inErr = postgresqlService.DescribePgParams(ctx, d.Id())
		if inErr != nil {
			ee, ok := inErr.(*sdkErrors.TencentCloudSDKError)
			if ok && ee.GetCode() == "ResourceNotFound.InstanceNotFoundError" {
				return nil
			}

			return tccommon.RetryError(inErr)
		}

		return nil
	})

	if err != nil {
		return err
	}

	maxStandbyStreamingDelayValue, err := strconv.Atoi(parmas["max_standby_streaming_delay"])
	if err != nil {
		return err
	}

	maxStandbyArchiveDelayValue, err := strconv.Atoi(parmas["max_standby_archive_delay"])
	if err != nil {
		return err
	}

	_ = d.Set("max_standby_streaming_delay", maxStandbyStreamingDelayValue)
	_ = d.Set("max_standby_archive_delay", maxStandbyArchiveDelayValue)

	return nil
}

func resourceTencentCloudPostgresqlInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_instance.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	//internal version: replace clientUpdate begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	postgresqlService := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	//internal version: replace clientUpdate end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	instanceId := d.Id()
	d.Partial(true)

	if err := helper.ImmutableArgsChek(d,
		"auto_voucher",
		"voucher_ids",
		"root_user",
	); err != nil {
		return err
	}

	if d.HasChange("need_support_tde") || d.HasChange("kms_key_id") || d.HasChange("kms_region") || d.HasChange("kms_cluster_id") {
		return fmt.Errorf("Not support change params contact with data transparent encryption.")
	}

	waitSwitch := POSTGRESQL_KERNEL_UPGRADE_IMMEDIATELY
	if v, ok := d.GetOk("wait_switch"); ok {
		waitSwitch = v.(int)
	}

	if d.HasChange("charge_type") {
		var (
			request       = postgresql.NewModifyDBInstanceChargeTypeRequest()
			chargeTypeOld string
			chargeTypeNew string
		)

		old, new := d.GetChange("charge_type")
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
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"running"}, 2*tccommon.ReadRetryTimeout, time.Second, service.PostgresqlDBInstanceStateRefreshFunc(instanceId, []string{}))
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
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"running"}, 2*tccommon.ReadRetryTimeout, time.Second, service.PostgresqlDBInstanceStateRefreshFunc(instanceId, []string{}))
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

		if v, ok := d.GetOkExists("auto_voucher"); ok {
			request.AutoVoucher = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("voucher_ids"); ok {
			request.VoucherIds = helper.InterfacesStringsPoint(v.([]interface{}))
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
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"running"}, 2*tccommon.ReadRetryTimeout, time.Second, service.PostgresqlDBInstanceStateRefreshFunc(instanceId, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	var outErr, inErr, checkErr error
	// update vpc and subnet
	if d.HasChange("vpc_id") || d.HasChange("subnet_id") {
		var (
			postgresqlService = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
			instance          *postgresql.DBInstance
			has               bool
			vpcOld            string
			vpcNew            string
			subnetOld         string
			subnetNew         string
			vipOld            string
			vipNew            string
		)

		// check net first
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
			conf := tccommon.BuildStateChangeConf([]string{}, []string{"opened"}, 3*tccommon.ReadRetryTimeout, time.Second, postgresqlService.PostgresqlDBInstanceNetworkAccessStateRefreshFunc(instanceId, vpcNew, subnetNew, vipOld, "", []string{}))
			if object, e := conf.WaitForState(); e != nil {
				return e
			} else {
				// find the vip assiged by system
				ret := object.(*postgresql.DBInstanceNetInfo)
				vipNew = *ret.Ip
			}

			// wait unit network changing operation of instance done
			conf = tccommon.BuildStateChangeConf([]string{}, []string{"running"}, 3*tccommon.ReadRetryTimeout, time.Second, postgresqlService.PostgresqlDBInstanceStateRefreshFunc(instanceId, []string{}))
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
			conf = tccommon.BuildStateChangeConf([]string{}, []string{"closed"}, 3*tccommon.ReadRetryTimeout, time.Second, postgresqlService.PostgresqlDBInstanceNetworkAccessStateRefreshFunc(instanceId, vpcOld, subnetOld, vipNew, vipOld, []string{}))
			if _, e := conf.WaitForState(); e != nil {
				return e
			}

			// wait unit network changing operation of instance done
			conf = tccommon.BuildStateChangeConf([]string{}, []string{"running"}, 3*tccommon.ReadRetryTimeout, time.Second, postgresqlService.PostgresqlDBInstanceStateRefreshFunc(instanceId, []string{}))
			if _, e := conf.WaitForState(); e != nil {
				return e
			}

			// refresh the private ip with new one
			_ = d.Set("private_access_ip", vipNew)
		}
	}

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
	if d.HasChange("memory") || d.HasChange("storage") || d.HasChange("cpu") {
		memory := d.Get("memory").(int)
		storage := d.Get("storage").(int)
		var cpu int
		if v, ok := d.GetOkExists("cpu"); ok {
			cpu = v.(int)
		}

		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = postgresqlService.UpgradePostgresqlInstance(ctx, instanceId, memory, storage, cpu, waitSwitch)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}

			return nil
		})

		if outErr != nil {
			return outErr
		}

		if waitSwitch == POSTGRESQL_KERNEL_UPGRADE_IMMEDIATELY {
			// Wait for status to processing
			_ = resource.Retry(time.Second*10, func() *resource.RetryError {
				instance, _, err := postgresqlService.DescribePostgresqlInstanceById(ctx, instanceId)
				if err != nil {
					return tccommon.RetryError(err)
				}

				if *instance.DBInstanceStatus == POSTGRESQL_STAUTS_RUNNING {
					return resource.RetryableError(fmt.Errorf("waiting for upgrade status change"))
				}

				return nil
			})

			time.Sleep(time.Second * 5)
			// check update storage and memory done
			checkErr = postgresqlService.CheckDBInstanceStatus(ctx, instanceId, 60)
			if checkErr != nil {
				return checkErr
			}
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

	// update public access
	if d.HasChange("public_access_switch") {
		public_access_switch := false
		if v, ok := d.GetOkExists("public_access_switch"); ok {
			public_access_switch = v.(bool)
		}

		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = postgresqlService.ModifyPublicService(ctx, public_access_switch, instanceId)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}

			return nil
		})

		if outErr != nil {
			return outErr
		}

		// check update public service done
		checkErr = postgresqlService.CheckDBInstanceStatus(ctx, instanceId)
		if checkErr != nil {
			return checkErr
		}
	}

	// update root password
	if d.HasChange("root_password") {
		// to avoid other updating process conflicts with updating password, set the password updating with the last step, there is no way to figure out whether changing password is done
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = postgresqlService.SetPostgresqlInstanceRootPassword(ctx, instanceId, d.Get("root_user").(string), d.Get("root_password").(string))
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}

			return nil
		})

		if outErr != nil {
			return outErr
		}

		// check update password done
		checkErr = postgresqlService.CheckDBInstanceStatus(ctx, instanceId)
		if checkErr != nil {
			return checkErr
		}
	}

	if d.HasChange("security_groups") {
		ids := d.Get("security_groups").(*schema.Set).List()
		var sgIds []*string
		for _, id := range ids {
			sgIds = append(sgIds, helper.String(id.(string)))
		}

		err := postgresqlService.ModifyDBInstanceSecurityGroupsById(ctx, d.Id(), sgIds)
		if err != nil {
			return err
		}
	}

	if d.HasChange("backup_plan") {
		if plan, ok := helper.InterfacesHeadMap(d, "backup_plan"); ok {
			request := postgresql.NewModifyBackupPlanRequest()
			request.DBInstanceId = &instanceId
			if v, ok := plan["min_backup_start_time"].(string); ok && v != "" {
				request.MinBackupStartTime = &v
			}

			if v, ok := plan["max_backup_start_time"].(string); ok && v != "" {
				request.MaxBackupStartTime = &v
			}

			if v, ok := plan["base_backup_retention_period"].(int); ok && v != 0 {
				request.BaseBackupRetentionPeriod = helper.IntUint64(v)
			}

			if v, ok := plan["backup_period"].([]interface{}); ok && len(v) > 0 {
				request.BackupPeriod = helper.InterfacesStringsPoint(v)
			}

			err := postgresqlService.ModifyBackupPlan(ctx, request)
			if err != nil {
				return err
			}

			request1 := postgresql.NewModifyBackupPlanRequest()
			request1.DBInstanceId = &instanceId
			var hasMonthlybackupPeriod bool
			if v, ok := plan["min_backup_start_time"].(string); ok && v != "" {
				request1.MinBackupStartTime = &v
			}

			if v, ok := plan["max_backup_start_time"].(string); ok && v != "" {
				request1.MaxBackupStartTime = &v
			}

			if v, ok := plan["monthly_backup_retention_period"].(int); ok && v != 0 {
				request1.BaseBackupRetentionPeriod = helper.IntUint64(v)
			}

			if v, ok := plan["monthly_backup_period"].([]interface{}); ok && len(v) > 0 {
				request1.BackupPeriod = helper.InterfacesStringsPoint(v)
				hasMonthlybackupPeriod = true
			}

			var monthlyPlanId string
			if v, ok := plan["monthly_plan_id"].(string); ok && v != "" {
				request1.PlanId = helper.String(v)
				monthlyPlanId = v
			}
			if !hasMonthlybackupPeriod && monthlyPlanId != "" {
				request0 := postgresql.NewDeleteBackupPlanRequest()
				request0.DBInstanceId = &instanceId
				request0.PlanId = &monthlyPlanId
				err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().DeleteBackupPlan(request0)
					if e != nil {
						return tccommon.RetryError(e)
					}
					return nil
				})

				if err != nil {
					return err
				}
			} else if hasMonthlybackupPeriod && monthlyPlanId == "" {
				request00 := postgresql.NewCreateBackupPlanRequest()
				request00.DBInstanceId = &instanceId
				request00.BackupPeriodType = helper.String("month")
				request00.PlanName = helper.String("custom_month")
				if v, ok := plan["monthly_backup_period"].([]interface{}); ok && len(v) > 0 {
					request00.BackupPeriod = helper.InterfacesStringsPoint(v)
				}

				if v, ok := plan["min_backup_start_time"].(string); ok && v != "" {
					request00.MinBackupStartTime = &v
				}

				if v, ok := plan["max_backup_start_time"].(string); ok && v != "" {
					request00.MaxBackupStartTime = &v
				}

				if v, ok := plan["monthly_backup_retention_period"].(int); ok && v != 0 {
					request00.BaseBackupRetentionPeriod = helper.IntUint64(v)
				}
				err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().CreateBackupPlan(request00)
					if e != nil {
						return tccommon.RetryError(e)
					}
					return nil
				})

				if err != nil {
					return err
				}
			} else {
				request1.PlanId = helper.String(monthlyPlanId)
				err = postgresqlService.ModifyBackupPlan(ctx, request1)
				if err != nil {
					return err
				}
			}

		}
	}

	if d.HasChange("db_node_set") {
		if include, z, nzs := checkZoneSetInclude(d); !include {
			return fmt.Errorf("`availability_zone`: %s is not included in `db_node_set`: %s", z, nzs)
		}

		nodeSet := d.Get("db_node_set").(*schema.Set).List()
		request := postgresql.NewModifyDBInstanceDeploymentRequest()
		request.DBInstanceId = helper.String(d.Id())
		request.SwitchTag = helper.IntInt64(waitSwitch)
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

		if waitSwitch == POSTGRESQL_KERNEL_UPGRADE_IMMEDIATELY {
			err = resource.Retry(tccommon.ReadRetryTimeout*10, func() *resource.RetryError {
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
	}

	if d.HasChange("availability_zone") {
		log.Printf("[WARN] argument `availability_zone` modified, skip process")
		return fmt.Errorf("The `availability_zone` cannot be modified, please use `db_node_set` instead of it.")
	}

	if d.HasChange("db_kernel_version") || d.HasChange("db_major_vesion") || d.HasChange("db_major_version") || d.HasChange("engine_version") {
		oldKVInterface, newKVInterface := d.GetChange("db_kernel_version")
		oldMVInterface, newMVInterface := d.GetChange("db_major_vesion")
		oldMVInterface, newMVInterface = d.GetChange("db_major_version")
		oldEVInterface, newEVInterface := d.GetChange("engine_version")

		oldKVValue := oldKVInterface.(string)
		newKVValue := newKVInterface.(string)
		oldMVValue := oldMVInterface.(string)
		newMVValue := newMVInterface.(string)
		oldEVValue := oldEVInterface.(string)
		newEVValue := newEVInterface.(string)

		// check old version value
		oldParamMap := make(map[string]interface{})
		oldParamMap["DBVersion"] = oldEVValue
		oldParamMap["DBMajorVersion"] = oldMVValue
		oldParamMap["DBKernelVersion"] = oldKVValue
		oldResult, err := postgresqlService.DescribePostgresqlDbVersionsByFilter(ctx, oldParamMap)
		if err != nil {
			return err
		}

		if oldResult == nil || len(oldResult) == 0 {
			return fmt.Errorf("Current postgresql instance engine_version: %s, db_major_version: %s and db_kernel_version: %s. has no available upgrade target verison.", oldEVValue, oldMVValue, oldKVValue)
		}

		// check new version value
		newParamMap := make(map[string]interface{})
		newParamMap["DBVersion"] = newEVValue
		newParamMap["DBMajorVersion"] = newMVValue
		newParamMap["DBKernelVersion"] = newKVValue
		newResult, err := postgresqlService.DescribePostgresqlDbVersionsByFilter(ctx, newParamMap)
		if err != nil {
			return err
		}

		if newResult == nil || len(newResult) == 0 {
			return fmt.Errorf("The expected modifications of engine_version: %s, db_major_version: %s and db_kernel_version: %s are illegal, available upgrade target verison cannot be found.", newEVValue, newMVValue, newKVValue)
		}

		if oldMVValue != newMVValue {
			// use UpgradeDBInstanceMajorVersion
			upgradeRequest := postgresql.NewUpgradeDBInstanceMajorVersionRequest()
			upgradeRequest.DBInstanceId = &instanceId
			upgradeRequest.TargetDBKernelVersion = &newKVValue
			upgradeRequest.UpgradeTimeOption = helper.IntInt64(waitSwitch)
			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().UpgradeDBInstanceMajorVersion(upgradeRequest)
				if e != nil {
					tcErr := e.(*sdkErrors.TencentCloudSDKError)
					if tcErr.Code == "FailedOperation.FailedOperationError" {
						// upgrade version invalid.
						return resource.NonRetryableError(fmt.Errorf("Upgrade major version failed: %v", e.Error()))
					}

					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, upgradeRequest.GetAction(), upgradeRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s upgrade major version failed, reason:%+v", logId, err)
				return err
			}
		} else {
			// use UpgradeDBInstanceKernelVersion
			upgradeRequest := postgresql.NewUpgradeDBInstanceKernelVersionRequest()
			upgradeRequest.DBInstanceId = &instanceId
			upgradeRequest.TargetDBKernelVersion = &newKVValue
			upgradeRequest.SwitchTag = helper.IntUint64(waitSwitch)
			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().UpgradeDBInstanceKernelVersion(upgradeRequest)
				if e != nil {
					tcErr := e.(*sdkErrors.TencentCloudSDKError)
					if tcErr.Code == "FailedOperation.FailedOperationError" {
						// upgrade version invalid.
						return resource.NonRetryableError(fmt.Errorf("Upgrade kernel version failed: %v", e.Error()))
					}

					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, upgradeRequest.GetAction(), upgradeRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s upgrade kernel version failed, reason:%+v", logId, err)
				return err
			}
		}

		if waitSwitch == POSTGRESQL_KERNEL_UPGRADE_IMMEDIATELY {
			// only wait for immediately upgrade mode
			conf := tccommon.BuildStateChangeConf([]string{}, []string{"running", "isolated", "offline"}, 10*tccommon.ReadRetryTimeout, time.Second, postgresqlService.PostgresqlUpgradeKernelVersionRefreshFunc(d.Id(), []string{}))
			if _, e := conf.WaitForState(); e != nil {
				return e
			}
		}
	}

	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))

		//internal version: replace null begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		//internal version: replace null end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
		resourceName := tccommon.BuildTagResourceName("postgres", "DBInstanceId", tcClient.Region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
		//internal version: replace waitTag begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
		//internal version: replace waitTag end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	}

	paramEntrys := make(map[string]string)
	if d.HasChange("max_standby_archive_delay") {
		if v, ok := d.GetOkExists("max_standby_archive_delay"); ok {
			paramEntrys["max_standby_archive_delay"] = strconv.Itoa(v.(int))
		}
	}

	if d.HasChange("max_standby_streaming_delay") {
		if v, ok := d.GetOkExists("max_standby_streaming_delay"); ok {
			paramEntrys["max_standby_streaming_delay"] = strconv.Itoa(v.(int))
		}
	}

	if len(paramEntrys) != 0 {
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr := postgresqlService.ModifyPgParams(ctx, instanceId, paramEntrys)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}

			return nil
		})

		if outErr != nil {
			return outErr
		}

		// 10s is required to synchronize the data(`DescribeParamsEvent`).
		time.Sleep(10 * time.Second)
	}

	d.Partial(false)

	return resourceTencentCloudPostgresqlInstanceRead(d, meta)
}

func resourceTencentCLoudPostgresqlInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_instance.delete")()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		postgresqlService = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		outErr, inErr     error
		has               bool
		deleteProtection  bool
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
	_ = resource.Retry(tccommon.ReadRetryTimeout*5, func() *resource.RetryError {
		instance, _, err := postgresqlService.DescribePostgresqlInstanceById(ctx, instanceId)
		if err != nil {
			return tccommon.RetryError(err)
		}

		if *instance.DBInstanceStatus == POSTGRESQL_STAUTS_ISOLATED {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("waiting for instance isolating"))
	})

	if v, ok := d.GetOkExists("delete_protection"); ok {
		deleteProtection = v.(bool)
	}

	// delete protection
	if deleteProtection {
		return nil
	}

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

	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
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

// check availability_zone included in db_node_set
func checkZoneSetInclude(d *schema.ResourceData) (included bool, zone string, nodeZoneList []string) {
	zone = d.Get("availability_zone").(string)
	dbNodeSet := d.Get("db_node_set").(*schema.Set).List()
	for i := range dbNodeSet {
		item := dbNodeSet[i].(map[string]interface{})
		nodeZone := item["zone"].(string)
		if nodeZone == zone {
			included = true
		}

		nodeZoneList = append(nodeZoneList, nodeZone)
	}

	return
}
