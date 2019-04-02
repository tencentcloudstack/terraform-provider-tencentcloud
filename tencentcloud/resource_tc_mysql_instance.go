package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
)

func TencentMsyqlBasicInfo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"pay_type": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			Default:      0,
		},
		"period": {
			Type:         schema.TypeInt,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validateAllowedIntValue(MYSQL_AVAILABLE_PERIOD),
		},
		"auto_renew_flag": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			Default:      0,
		},
		"engine_version": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validateAllowedStringValue(MYSQL_SUPPORTS_ENGINE),
		},
		"mem_size": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"volume_size": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"vpc_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"subnet_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"internet_service": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			Default:      0,
		},
		"gtid": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			Default:      0,
		},
		"project_id": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"security_groups": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Set: func(v interface{}) int {
				return hashcode.String(v.(string))
			},
		},

		"parameters": {
			Type:     schema.TypeMap,
			Optional: true,
		},
		"tags": {
			Type:     schema.TypeMap,
			Optional: true,
		},

		// Computed values
		"intranet_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"intranet_port": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"locked": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"task_status": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"internet_host": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"internet_port": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func resourceTencentCloudMysqlInstance() *schema.Resource {
	specialInfo := map[string]*schema.Schema{
		"availability_zone": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"root_password": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validateMysqlPassword,
		},
		"slave_deploy_mode": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			Default:      0,
		},
		"first_slave_zone": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"second_slave_zone": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"slave_sync_mode": {
			Type:         schema.TypeInt,
			ValidateFunc: validateAllowedIntValue([]int{0, 1, 2}),
			Optional:     true,
			Default:      0,
		},
	}

	basic := TencentMsyqlBasicInfo()
	for k, v := range basic {
		specialInfo[k] = v
	}
	return &schema.Resource{
		Create: resourceTencentCloudMysqlInstanceCreate,
		Read:   resourceTencentCloudMysqlInstanceRead,
		Update: resourceTencentCloudMysqlInstanceUpdate,
		Delete: resourceTencentCloudMysqlInstanceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: specialInfo,
	}
}

func resourceTencentCloudMysqlInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId("cdb-8xtme2cj")
	return resourceTencentCloudMysqlInstanceRead(d, meta)
}

func tencentMsyqlBasicInfoRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (mysqlInfo *cdb.InstanceInfo,
	errRet error) {

	logId := GetLogId(ctx)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	mysqlInfo, errRet = mysqlService.DescribeDBInstanceById(ctx, d.Id())

	if errRet != nil {
		errRet = fmt.Errorf("Describe mysql instance fails, reaseon %s", errRet.Error())
		return
	}
	if mysqlInfo == nil {
		d.SetId("")
		return
	}

	d.Set("instance_name", *mysqlInfo.InstanceName)
	d.Set("pay_type", int(*mysqlInfo.PayType))

	tempInt, _ := d.Get("period").(int)
	if tempInt == 0 {
		d.Set("period", 1)
	}

	if *mysqlInfo.AutoRenew == MYSQL_RENEW_CLOSE {
		*mysqlInfo.AutoRenew = MYSQL_RENEW_NOUSE
	}
	d.Set("auto_renew_flag", int(*mysqlInfo.AutoRenew))

	d.Set("engine_version", *mysqlInfo.EngineVersion)
	d.Set("mem_size", *mysqlInfo.Memory)
	d.Set("volume_size", *mysqlInfo.Volume)
	d.Set("vpc_id", *mysqlInfo.UniqVpcId)
	d.Set("subnet_id", *mysqlInfo.UniqSubnetId)

	if *mysqlInfo.WanStatus == 1 {
		d.Set("internet_service", 1)
		d.Set("internet_host", *mysqlInfo.WanDomain)
		d.Set("internet_port", int(*mysqlInfo.WanPort))
	} else {
		d.Set("internet_service", 0)
		d.Set("internet_host", "")
		d.Set("internet_port", 0)
	}

	isGTIDOpen, err := mysqlService.CheckDBGTIDOpen(ctx, d.Id())
	if err != nil {
		errRet = err
		return
	}
	d.Set("gtid", int(isGTIDOpen))
	d.Set("project_id", int(*mysqlInfo.ProjectId))

	securityGroups, err := mysqlService.DescribeDBSecurityGroups(ctx, d.Id())
	if err != nil {
		errRet = err
		return
	}
	d.Set("security_groups", securityGroups)

	parametersMap, ok := d.Get("parameters").(map[string]interface{})
	if !ok {
		log.Printf("[INFO] %d  config error,parameters is not map[string]interface{}\n", logId)
	} else {
		var cares []string
		for k, _ := range parametersMap {
			cares = append(cares, k)
		}
		caresParameters, err := mysqlService.DescribeCaresParameters(ctx, d.Id(), cares)
		if err != nil {
			errRet = err
			return
		}
		if err := d.Set("parameters", caresParameters); err != nil {
			log.Printf("[CRITAL]%s provider set caresParameters fail, reason:%s\n ", logId, err.Error())
		}
	}
	tags, err := mysqlService.DescribeTagsOfInstanceId(ctx, d.Id())
	if err != nil {
		errRet = err
		return
	}
	if err := d.Set("tags", tags); err != nil {
		log.Printf("[CRITAL]%s provider set tags fail, reason:%s\n ", logId, err.Error())
	}

	d.Set("intranet_ip", *mysqlInfo.Vip)
	d.Set("intranet_port", int(*mysqlInfo.Vport))

	if *mysqlInfo.CdbError != 0 {
		d.Set("locked", 1)
	} else {
		d.Set("locked", 0)
	}
	d.Set("status", *mysqlInfo.Status)
	d.Set("task_status", *mysqlInfo.TaskStatus)
	return
}

func resourceTencentCloudMysqlInstanceRead(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	mysqlInfo, err := tencentMsyqlBasicInfoRead(ctx, d, meta)
	if err != nil {
		return err
	}
	d.Set("availability_zone", *mysqlInfo.Zone)

	backConfig, err := mysqlService.DescribeDBInstanceConfig(ctx, d.Id())
	d.Set("slave_sync_mode", int(*backConfig.Response.ProtectMode))
	d.Set("slave_deploy_mode", int(*backConfig.Response.DeployMode))

	if backConfig.Response.SlaveConfig != nil && *backConfig.Response.SlaveConfig.Zone != "" {
		d.Set("first_slave_zone", *backConfig.Response.SlaveConfig.Zone)
	}
	if backConfig.Response.BackupConfig != nil && *backConfig.Response.BackupConfig.Zone != "" {
		d.Set("second_slave_zone", *backConfig.Response.BackupConfig.Zone)
	}

	return nil
}
func resourceTencentCloudMysqlInstanceUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}
func resourceTencentCloudMysqlInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
