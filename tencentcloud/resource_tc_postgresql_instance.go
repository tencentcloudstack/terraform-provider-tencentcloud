/*
Use this resource to create postgresql instance.

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-1"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "guagua_vpc_instance_test"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "guagua_vpc_subnet_test"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create postgresql
resource "tencentcloud_postgresql_instance" "foo" {
  name              = "example"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  engine_version    = "10.4"
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  memory            = 2
  storage           = 10

  tags = {
    test = "tf"
  }
}
```

Create a multi available zone bucket

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

variable "standby_availability_zone" {
  default = "ap-guangzhou-7"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "guagua_vpc_instance_test"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "guagua_vpc_subnet_test"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create postgresql
resource "tencentcloud_postgresql_instance" "foo" {
  name              = "example"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  engine_version    = "10.4"
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  memory            = 2
  storage           = 10

  db_node_set {
    role = "Primary"
    zone = var.availability_zone
  }
  db_node_set {
    zone = var.standby_availability_zone
  }

  tags = {
    test = "tf"
  }
}
```

create pgsql with kms key
```
resource "tencentcloud_postgresql_instance" "pg" {
  name              = "tf_postsql_instance"
  availability_zone = "ap-guangzhou-6"
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = "vpc-86v957zb"
  subnet_id         = "subnet-enm92y0m"
  engine_version    = "11.12"
  #  db_major_vesion   = "11"
  db_kernel_version = "v11.12_r1.3"
  need_support_tde  = 1
  kms_key_id        = "788c606a-c7b7-11ec-82d1-5254001e5c4e"
  kms_region        = "ap-guangzhou"
  root_password     = "xxxxxxxxxx"
  charset           = "LATIN1"
  project_id        = 0
  memory            = 4
  storage           = 100

  backup_plan {
    min_backup_start_time        = "00:10:11"
    max_backup_start_time        = "01:10:11"
    base_backup_retention_period = 7
    backup_period                = ["tuesday", "wednesday"]
  }

  tags = {
    tf = "test"
  }
}
```

upgrade kernel version
```
resource "tencentcloud_postgresql_instance" "test" {
  name = "tf_postsql_instance_update"
  availability_zone = data.tencentcloud_availability_zones_by_product.zone.zones[5].name
  charge_type	    = "POSTPAID_BY_HOUR"
  vpc_id  	  		= local.vpc_id
  subnet_id 		= local.subnet_id
  engine_version	= "13.3"
  root_password	    = "*"
  charset 			= "LATIN1"
  project_id 		= 0
  public_access_switch = false
  security_groups   = [local.sg_id]
  memory 			= 4
  storage 			= 250
  backup_plan {
	min_backup_start_time 		 = "01:10:11"
	max_backup_start_time		 = "02:10:11"
	base_backup_retention_period = 5
	backup_period 			     = ["monday", "thursday", "sunday"]
  }

  db_kernel_version = "v13.3_r1.4" # eg:from v13.3_r1.1 to v13.3_r1.4

  tags = {
	tf = "teest"
  }
}
```

Import

postgresql instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_postgresql_instance.foo postgres-cda1iex1
```
*/
package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudPostgresqlInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlInstanceCreate,
		Read:   resourceTencentCloudPostgresqlInstanceRead,
		Update: resourceTencentCloudPostgresqlInstanceUpdate,
		Delete: resourceTencentCLoudPostgresqlInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "Name of the postgresql instance.",
			},
			"charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     COMMON_PAYTYPE_POSTPAID,
				Description: "Pay type of the postgresql instance. Values `POSTPAID_BY_HOUR` (Default), `PREPAID`. It support to update the type from `POSTPAID_BY_HOUR` to `PREPAID`.",
			},
			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specify Prepaid period in month. Default `1`. Values: `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.",
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
				ForceNew:    true,
				Optional:    true,
				Default:     "10.4",
				Description: "Version of the postgresql database engine. Valid values: `10.4`, `11.8`, `12.4`.",
			},
			"db_major_vesion": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "`db_major_vesion` will be deprecated, use `db_major_version` instead.",
				Description: "PostgreSQL major version number. Valid values: 10, 11, 12, 13. " +
					"If it is specified, an instance running the latest kernel of PostgreSQL DBMajorVersion will be created.",
			},
			"db_major_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "PostgreSQL major version number. Valid values: 10, 11, 12, 13. " +
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
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Project id, default value is `0`.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Availability zone. NOTE: If value modified but included in `db_node_set`, the diff will be suppressed.",
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
				ForceNew:    true,
				Optional:    true,
				Default:     "root",
				Description: "Instance root account name. This parameter is optional, Default value is `root`.",
			},
			"root_password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validateMysqlPassword,
				Description:  "Password of root account. This parameter can be specified when you purchase master instances, but it should be ignored when you purchase read-only instances or disaster recovery instances.",
			},
			"charset": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      POSTGRESQL_DB_CHARSET_UTF8,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(POSTGRESQL_DB_CHARSET),
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
				Description: "Specify DB backup plan.",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_backup_start_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specify earliest backup start time, format `hh:mm:ss`.",
						},
						"max_backup_start_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specify latest backup start time, format `hh:mm:ss`.",
						},
						"base_backup_retention_period": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specify days of the retention.",
						},
						"backup_period": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of backup period per week, available values: `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, `sunday`. NOTE: At least specify two days.",
							Elem:        &schema.Schema{Type: schema.TypeString},
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
					},
				},
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
	defer logElapsed("resource.tencentcloud_postgresql_instance.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	postgresqlService := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		name           = d.Get("name").(string)
		dbVersion      = d.Get("engine_version").(string)
		payType        = d.Get("charge_type").(string)
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

	var instanceId, specVersion, specCode string
	var outErr, inErr error
	var allowVersion, allowMemory []string

	var (
		dbMajorVersion  = ""
		dbKernelVersion = ""
		needSupportTde  = 0
		kmsKeyId        = ""
		kmsRegion       = ""
		period          = 1
		autoRenewFlag   = 0
		autoVoucher     = 0
		voucherIds      []*string
	)

	if v, ok := d.GetOk("period"); ok {
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
	if v, ok := d.GetOk("need_support_tde"); ok {
		needSupportTde = v.(int)
	}
	if v, ok := d.GetOk("kms_key_id"); ok {
		kmsKeyId = v.(string)
	}
	if v, ok := d.GetOk("kms_region"); ok {
		kmsRegion = v.(string)
	}
	if v, ok := d.Get("auto_renew_flag").(int); ok {
		autoRenewFlag = v
	}
	if v, ok := d.Get("auto_voucher").(int); ok {
		autoVoucher = v
	}
	if v, ok := d.GetOk("voucher_ids"); ok {
		voucherIds = helper.InterfacesStringsPoint(v.([]interface{}))
	}
	requestSecurityGroup := make([]string, 0, len(securityGroups))

	for _, v := range securityGroups {
		requestSecurityGroup = append(requestSecurityGroup, v.(string))
	}

	// get specCode with engine_version and memory
	outErr = resource.Retry(readRetryTimeout*5, func() *resource.RetryError {
		speccodes, inErr := postgresqlService.DescribeSpecinfos(ctx, zone)
		if inErr != nil {
			return retryError(inErr)
		}
		for _, info := range speccodes {
			if !IsContains(allowVersion, *info.Version) {
				allowVersion = append(allowVersion, *info.Version)
			}
			if *info.Version == dbVersion {
				specVersion = *info.Version
				memoryString := fmt.Sprintf("%d", int(*info.Memory)/1024)
				if !IsContains(allowMemory, memoryString) {
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
	if outErr != nil {
		return outErr
	}

	if specVersion == "" {
		return fmt.Errorf(`The "engine_version" value: "%s" is invalid, Valid values are one of: "%s"`, dbVersion, strings.Join(allowVersion, `", "`))
	}

	if specCode == "" {
		return fmt.Errorf(`The "memory" value: %d is invalid, Valid values are one of: %s`, memory, strings.Join(allowMemory, `, `))
	}

	var dbNodeSet []*postgresql.DBNode
	if len(nodeSet) > 0 {

		for i := range nodeSet {
			var (
				item = nodeSet[i].(map[string]interface{})
				role = item["role"].(string)
				zone = item["zone"].(string)
				node = &postgresql.DBNode{
					Role: &role,
					Zone: &zone,
				}
			)
			dbNodeSet = append(dbNodeSet, node)
		}

		// check if availability_zone and node_set consists
		if include, z, nzs := checkZoneSetInclude(d); !include {
			return fmt.Errorf("`availability_zone`: %s is not included in `db_node_set`: %s", z, nzs)
		}
	}

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
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
			autoVoucher,
			voucherIds,
		)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(instanceId)

	// check creation done
	err := resource.Retry(20*readRetryTimeout, func() *resource.RetryError {
		instance, has, err := postgresqlService.DescribePostgresqlInstanceById(ctx, instanceId)
		if err != nil {
			return retryError(err)
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
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = postgresqlService.ModifyPublicService(ctx, true, instanceId)
			if inErr != nil {
				return retryError(inErr)
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
	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		inErr := postgresqlService.ModifyPostgresqlInstanceName(ctx, instanceId, name)
		if inErr != nil {
			return retryError(inErr)
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

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		resourceName := BuildTagResourceName("postgres", "DBInstanceId", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	// set pg params
	paramEntrys := make(map[string]string)
	if v, ok := d.GetOkExists("max_standby_archive_delay"); ok {
		paramEntrys["max_standby_archive_delay"] = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOkExists("max_standby_streaming_delay"); ok {
		paramEntrys["max_standby_streaming_delay"] = strconv.Itoa(v.(int))
	}

	if len(paramEntrys) != 0 {
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr := postgresqlService.ModifyPgParams(ctx, instanceId, paramEntrys)
			if inErr != nil {
				return retryError(inErr)
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
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			err := postgresqlService.ModifyBackupPlan(ctx, request)
			if err != nil {
				return retryError(err, postgresql.OPERATIONDENIED_INSTANCESTATUSLIMITOPERROR)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudPostgresqlInstanceRead(d, meta)
}

func resourceTencentCloudPostgresqlInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_instance.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	postgresqlService := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId := d.Id()
	d.Partial(true)

	if err := helper.ImmutableArgsChek(d,
		// "charge_type",
		// "period",
		// "auto_renew_flag",
		// "auto_voucher",
		"voucher_ids",
	); err != nil {
		return err
	}

	if d.HasChange("charge_type") {
		var (
			chargeTypeOld string
			chargeTypeNew string
			period        = 1
			autoRenew     = 0
			autoVoucher   = 0
			request       = postgresql.NewModifyDBInstanceChargeTypeRequest()
		)

		old, new := d.GetChange("charge_type")
		if old != nil {
			chargeTypeOld = old.(string)
		}
		if new != nil {
			chargeTypeNew = new.(string)
		}

		// The real Update operation rather than the operation from create
		if chargeTypeOld != "" && chargeTypeOld != chargeTypeNew {
			if v, ok := d.GetOk("period"); ok {
				log.Printf("period set")
				period = v.(int)
			} else {
				log.Printf("period not set")
			}

			if v, ok := d.GetOk("auto_renew_flag"); ok {
				log.Printf("auto_renew_flag set")
				autoRenew = v.(int)
			} else {
				log.Printf("auto_renew_flag not set")
			}

			if v, ok := d.GetOk("auto_voucher"); ok {
				log.Printf("auto_voucher set")
				autoVoucher = v.(int)
			} else {
				log.Printf("auto_voucher not set")
			}

			request.DBInstanceId = &instanceId
			request.InstanceChargeType = &chargeTypeNew
			request.Period = helper.IntInt64(period)
			request.AutoRenewFlag = helper.IntInt64(autoRenew)
			request.AutoVoucher = helper.IntInt64(autoVoucher)

			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().ModifyDBInstanceChargeType(request)
				if e != nil {
					return retryError(e)
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
			service := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}
			conf := BuildStateChangeConf([]string{}, []string{"running"}, 2*readRetryTimeout, time.Second, service.PostgresqlDBInstanceStateRefreshFunc(instanceId, []string{}))
			if _, e := conf.WaitForState(); e != nil {
				return e
			}
		}
	}

	var outErr, inErr, checkErr error
	// update vpc and subnet
	if d.HasChange("vpc_id") || d.HasChange("subnet_id") {
		var (
			vpcOld    string
			vpcNew    string
			subnetOld string
			subnetNew string
			vipOld    string
			vipNew    string
		)

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

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().CreateDBInstanceNetworkAccess(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create postgresql Instance NetworkAccess failed, reason:%+v", logId, err)
			return err
		}

		service := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}
		// wait for new network enabled
		conf := BuildStateChangeConf([]string{}, []string{"opened"}, 3*readRetryTimeout, time.Second, service.PostgresqlDBInstanceNetworkAccessStateRefreshFunc(instanceId, vpcNew, subnetNew, vipOld, "", []string{}))
		if object, e := conf.WaitForState(); e != nil {
			return e
		} else {
			// find the vip assiged by system
			ret := object.(*postgresql.DBInstanceNetInfo)
			vipNew = *ret.Ip
		}

		// wait unit network changing operation of instance done
		conf = BuildStateChangeConf([]string{}, []string{"running"}, 3*readRetryTimeout, time.Second, service.PostgresqlDBInstanceStateRefreshFunc(instanceId, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}

		// delete the old one
		if v, ok := d.GetOk("private_access_ip"); ok {
			vipOld = v.(string)
		}
		if err := service.DeletePostgresqlDBInstanceNetworkAccessById(ctx, instanceId, vpcOld, subnetOld, vipOld); err != nil {
			return err
		}

		// wait for old network removed
		conf = BuildStateChangeConf([]string{}, []string{"closed"}, 3*readRetryTimeout, time.Second, service.PostgresqlDBInstanceNetworkAccessStateRefreshFunc(instanceId, vpcOld, subnetOld, vipNew, vipOld, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}

		// wait unit network changing operation of instance done
		conf = BuildStateChangeConf([]string{}, []string{"running"}, 3*readRetryTimeout, time.Second, service.PostgresqlDBInstanceStateRefreshFunc(instanceId, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}

		// refresh the private ip with new one
		_ = d.Set("private_access_ip", vipNew)
	}

	// update name
	if d.HasChange("name") {
		name := d.Get("name").(string)
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = postgresqlService.ModifyPostgresqlInstanceName(ctx, instanceId, name)
			if inErr != nil {
				return retryError(inErr)
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
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = postgresqlService.UpgradePostgresqlInstance(ctx, instanceId, memory, storage)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}
		// Wait for status to processing
		_ = resource.Retry(time.Second*10, func() *resource.RetryError {
			instance, _, err := postgresqlService.DescribePostgresqlInstanceById(ctx, instanceId)
			if err != nil {
				return retryError(err)
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

	// update project id
	if d.HasChange("project_id") {
		projectId := d.Get("project_id").(int)
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = postgresqlService.ModifyPostgresqlInstanceProjectId(ctx, instanceId, projectId)
			if inErr != nil {
				return retryError(inErr)
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
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = postgresqlService.ModifyPublicService(ctx, public_access_switch, instanceId)
			if inErr != nil {
				return retryError(inErr)
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
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = postgresqlService.SetPostgresqlInstanceRootPassword(ctx, instanceId, d.Get("root_password").(string))
			if inErr != nil {
				return retryError(inErr)
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

		// Only redis service support modify Generic DB instance security groups
		service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
		ids := d.Get("security_groups").(*schema.Set).List()
		var sgIds []*string
		for _, id := range ids {
			sgIds = append(sgIds, helper.String(id.(string)))
		}
		err := service.ModifyDBInstanceSecurityGroups(ctx, "postgres", d.Id(), sgIds)
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
		}
	}

	if d.HasChange("db_node_set") {

		if include, z, nzs := checkZoneSetInclude(d); !include {
			return fmt.Errorf("`availability_zone`: %s is not included in `db_node_set`: %s", z, nzs)
		}

		nodeSet := d.Get("db_node_set").(*schema.Set).List()
		request := postgresql.NewModifyDBInstanceDeploymentRequest()
		request.DBInstanceId = helper.String(d.Id())
		request.SwitchTag = helper.Int64(0)
		for i := range nodeSet {
			var (
				node = nodeSet[i].(map[string]interface{})
				role = node["role"].(string)
				zone = node["zone"].(string)
			)
			request.DBNodeSet = append(request.DBNodeSet, &postgresql.DBNode{
				Role: &role,
				Zone: &zone,
			})
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			if err := postgresqlService.ModifyDBInstanceDeployment(ctx, request); err != nil {
				return retryError(err, postgresql.OPERATIONDENIED_INSTANCESTATUSLIMITOPERROR)
			}
			return nil
		})

		if err != nil {
			return err
		}

		err = resource.Retry(readRetryTimeout*10, func() *resource.RetryError {
			instance, _, err := postgresqlService.DescribePostgresqlInstanceById(ctx, d.Id())
			if err != nil {
				return retryError(err)
			}
			if IsContains(POSTGRESQL_RETRYABLE_STATUS, *instance.DBInstanceStatus) {
				return resource.RetryableError(fmt.Errorf("instance status is %s, retrying", *instance.DBInstanceStatus))
			}
			return nil
		})

		if err != nil {
			return err
		}

	}

	if d.HasChange("zone") {
		log.Printf("[WARN] argument `zone` modified, skip process")

	}

	if d.HasChange("db_kernel_version") {
		upgradeVersion := d.Get("db_kernel_version").(string)

		upgradeRequest := postgresql.NewUpgradeDBInstanceKernelVersionRequest()
		// upgradeResponse:= postgresql.NewUpgradeDBInstanceKernelVersionResponse()
		upgradeRequest.DBInstanceId = &instanceId
		upgradeRequest.TargetDBKernelVersion = &upgradeVersion

		// only support for the immediate upgrade policy
		switchTag := POSTGRESQL_KERNEL_UPGRADE_IMMEDIATELY
		upgradeRequest.SwitchTag = helper.IntUint64(switchTag)

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().UpgradeDBInstanceKernelVersion(upgradeRequest)
			if e != nil {
				tcErr := e.(*sdkErrors.TencentCloudSDKError)

				if tcErr.Code == "FailedOperation.FailedOperationError" {
					// upgrade version invalid.
					return resource.NonRetryableError(fmt.Errorf("Upgrade kernel version failed: %v", e.Error()))
				}
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, upgradeRequest.GetAction(), upgradeRequest.ToJsonString(), result.ToJsonString())
			}

			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create dcdb dbInstance failed, reason:%+v", logId, err)
			return err
		}

		// only wait for immediately upgrade mode

		conf := BuildStateChangeConf([]string{}, []string{"running", "isolated", "offline"}, 10*readRetryTimeout, time.Second, postgresqlService.PostgresqlUpgradeKernelVersionRefreshFunc(d.Id(), []string{}))

		if _, e := conf.WaitForState(); e != nil {
			return e
		}

	}

	if d.HasChange("tags") {

		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))

		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		resourceName := BuildTagResourceName("postgres", "DBInstanceId", tcClient.Region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}

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
	if d.HasChange("db_major_vesion") || d.HasChange("db_major_version") {
		return fmt.Errorf("Not support change db major version.")
	}

	if d.HasChange("need_support_tde") || d.HasChange("kms_key_id") || d.HasChange("kms_region") {
		return fmt.Errorf("Not support change params contact with data transparent encryption.")
	}
	if len(paramEntrys) != 0 {
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr := postgresqlService.ModifyPgParams(ctx, instanceId, paramEntrys)
			if inErr != nil {
				return retryError(inErr)
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

func resourceTencentCloudPostgresqlInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_instance.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		instance *postgresql.DBInstance
		has      bool
		outErr,
		inErr error
	)
	// Check if import
	postgresqlService := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instance, has, inErr = postgresqlService.DescribePostgresqlInstanceById(ctx, d.Id())
		if inErr != nil {
			ee, ok := inErr.(*sdkErrors.TencentCloudSDKError)
			if ok && (ee.GetCode() == "ResourceNotFound.InstanceNotFoundError" || ee.GetCode() == "InvalidParameter") {
				return nil
			}
			return retryError(inErr)
		}
		if instance != nil && IsContains(POSTGRESQL_RETRYABLE_STATUS, *instance.DBInstanceStatus) {
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
	_ = d.Set("vpc_id", instance.VpcId)
	_ = d.Set("subnet_id", instance.SubnetId)
	_ = d.Set("engine_version", instance.DBVersion)
	_ = d.Set("db_kernel_version", instance.DBKernelVersion)
	_ = d.Set("db_major_vesion", instance.DBMajorVersion)
	_ = d.Set("db_major_version", instance.DBMajorVersion)
	_ = d.Set("name", instance.DBInstanceName)
	_ = d.Set("charset", instance.DBCharset)
	if rootUser != "" {
		_ = d.Set("root_user", &rootUser)
	}

	if *instance.PayType == POSTGRESQL_PAYTYPE_PREPAID || *instance.PayType == COMMON_PAYTYPE_PREPAID {
		_ = d.Set("charge_type", COMMON_PAYTYPE_PREPAID)
	} else {
		_ = d.Set("charge_type", COMMON_PAYTYPE_POSTPAID)
	}

	// net status
	public_access_switch := false
	if len(instance.DBInstanceNetInfo) > 0 {
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
	}
	_ = d.Set("public_access_switch", public_access_switch)

	// security groups
	// Only redis service support modify Generic DB instance security groups
	redisService := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
	sg, err := redisService.DescribeDBSecurityGroups(ctx, "postgres", d.Id())
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
				"role": item.Role,
				"zone": item.Zone,
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
	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "postgres", "DBInstanceId", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	// backup plans (only specified will rewrite)
	if _, ok := d.GetOk("backup_plan"); ok {
		request := postgresql.NewDescribeBackupPlansRequest()
		request.DBInstanceId = helper.String(d.Id())
		response, err := postgresqlService.DescribeBackupPlans(ctx, request)

		if err != nil {
			return err
		}

		var backupPlan *postgresql.BackupPlan

		if len(response) > 0 {
			backupPlan = response[0]
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
			_ = d.Set("backup_plan", []interface{}{planMap})
		}

	}

	// pg params
	var parmas map[string]string
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		parmas, inErr = postgresqlService.DescribePgParams(ctx, d.Id())
		if inErr != nil {
			ee, ok := inErr.(*sdkErrors.TencentCloudSDKError)
			if ok && ee.GetCode() == "ResourceNotFound.InstanceNotFoundError" {
				return nil
			}
			return retryError(inErr)
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

func resourceTencentCLoudPostgresqlInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_instance.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Id()
	postgresqlService := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	var outErr, inErr error
	var has bool

	outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, inErr = postgresqlService.DescribePostgresqlInstanceById(ctx, d.Id())
		if inErr != nil {
			// ResourceNotFound.InstanceNotFoundError
			ee, ok := inErr.(*sdkErrors.TencentCloudSDKError)
			if ok && ee.GetCode() == "ResourceNotFound.InstanceNotFoundError" {
				return nil
			}
			return retryError(inErr, postgresql.OPERATIONDENIED_INSTANCESTATUSLIMITOPERROR)
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

	outErr = postgresqlService.DeletePostgresqlInstance(ctx, instanceId)
	if outErr != nil {
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = postgresqlService.DeletePostgresqlInstance(ctx, instanceId)
			if inErr != nil {
				// ResourceNotFound.InstanceNotFoundError
				ee, ok := inErr.(*sdkErrors.TencentCloudSDKError)
				if ok && ee.GetCode() == "ResourceNotFound.InstanceNotFoundError" {
					return nil
				}
				return retryError(inErr)
			}
			return nil
		})
	}

	if outErr != nil {
		return outErr
	}

	outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, inErr = postgresqlService.DescribePostgresqlInstanceById(ctx, d.Id())
		if inErr != nil {
			// ResourceNotFound.InstanceNotFoundError
			ee, ok := inErr.(*sdkErrors.TencentCloudSDKError)
			if ok && ee.GetCode() == "ResourceNotFound.InstanceNotFoundError" {
				return nil
			}
			return retryError(inErr)
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
