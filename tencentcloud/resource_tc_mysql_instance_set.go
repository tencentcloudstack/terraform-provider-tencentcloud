/*
Provides a mysql instance resource to create master database instances.

~> **NOTE:** If this mysql has readonly instance, the terminate operation of the mysql does NOT take effect immediately, maybe takes for several hours. so during that time, VPCs associated with that mysql instance can't be terminated also.

Example Usage

```hcl
resource "tencentcloud_mysql_instance_set" "default" {
  internet_service = 1
  engine_version   = "5.7"
  charge_type = "POSTPAID"
  root_password     = "********"
  slave_deploy_mode = 0
  first_slave_zone  = "ap-guangzhou-4"
  second_slave_zone = "ap-guangzhou-4"
  slave_sync_mode   = 1
  availability_zone = "ap-guangzhou-4"
  project_id        = 201901010001
  instance_name     = "myTestMysql"
  mem_size          = 128000
  volume_size       = 250
  vpc_id            = "vpc-12mt3l31"
  subnet_id         = "subnet-9uivyb1g"
  intranet_port     = 3306
  security_groups   = ["sg-ot8eclwz"]

  tags = {
    name = "test"
  }

  parameters = {
    character_set_server = "UTF8"
    max_connections = "1000"
  }
}
```

Import

MySQL instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_mysql_instance.foo cdb-12345678"
```

*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMysqlInstanceSet() *schema.Resource {
	specialInfo := map[string]*schema.Schema{
		"parameters": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "List of parameters to use.",
		},
		"internet_service": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			Default:      0,
			Description:  "Indicates whether to enable the access to an instance from public network: 0 - No, 1 - Yes.",
		},
		"engine_version": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			ValidateFunc: validateAllowedStringValue(MYSQL_SUPPORTS_ENGINE),
			Default:      MYSQL_SUPPORTS_ENGINE[len(MYSQL_SUPPORTS_ENGINE)-2],
			Description:  "The version number of the database engine to use. Supported versions include 5.5/5.6/5.7/8.0, and default is 5.7.",
		},

		"availability_zone": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Computed:    true,
			Description: "Indicates which availability zone will be used.",
		},
		"root_password": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validateMysqlPassword,
			Description:  "Password of root account. This parameter can be specified when you purchase master instances, but it should be ignored when you purchase read-only instances or disaster recovery instances.",
		},
		"slave_deploy_mode": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			Default:      0,
			Description:  "Availability zone deployment method. Available values: 0 - Single availability zone; 1 - Multiple availability zones.",
		},
		"first_slave_zone": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "Zone information about first slave instance.",
		},
		"second_slave_zone": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "Zone information about second slave instance.",
		},
		"slave_sync_mode": {
			Type:         schema.TypeInt,
			ForceNew:     true,
			Optional:     true,
			ValidateFunc: validateAllowedIntValue([]int{0, 1, 2}),
			Default:      0,
			Description:  "Data replication mode. 0 - Async replication; 1 - Semisync replication; 2 - Strongsync replication.",
		},
		"project_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "Project ID, default value is 0.",
		},

		// Computed values
		"internet_host": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "host for public access.",
		},
		"internet_port": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Access port for public access.",
		},
	}

	basic := TencentMsyqlBasicInfo()
	for k, v := range basic {
		specialInfo[k] = v
	}
	return &schema.Resource{
		Create: resourceTencentCloudMysqlInstanceSetCreate,
		Read:   resourceTencentCloudMysqlInstanceSetRead,
		Update: resourceTencentCloudMysqlInstanceSetUpdate,
		Delete: resourceTencentCloudMysqlInstanceSetDelete,
		Schema: specialInfo,
		Importer: &schema.ResourceImporter{
			State: helper.ImportWithDefaultValue(map[string]interface{}{
				"charge_type":       MYSQL_CHARGE_TYPE_POSTPAID,
				"prepaid_period":    1,
				"auto_renew_flag":   0,
				"intranet_port":     3306,
				"force_delete":      false,
				"internet_service":  0,
				"engine_version":    MYSQL_SUPPORTS_ENGINE[len(MYSQL_SUPPORTS_ENGINE)-2],
				"slave_deploy_mode": 0,
				"slave_sync_mode":   0,
				"project_id":        0,
			}),
		},
	}
}

func resourceTencentCloudMysqlInstanceSetCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_instance.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	payType := getPayType(d).(int)

	if payType == MysqlPayByMonth {
		err := mysqlCreateInstancePayByMonth(ctx, d, meta)
		if err != nil {
			return err
		}
	} else if payType == MysqlPayByUse {
		err := mysqlCreateInstancePayByUse(ctx, d, meta)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("mysql not support this pay type yet.")
	}

	mysqlID := d.Id()

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		resourceName := BuildTagResourceName("cdb", "instanceId", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	err := resource.Retry(7*readRetryTimeout, func() *resource.RetryError {
		mysqlInfo, err := mysqlService.DescribeDBInstanceById(ctx, mysqlID)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if mysqlInfo == nil {
			err = fmt.Errorf("mysqlid %s instance not exists", mysqlID)
			return resource.NonRetryableError(err)
		}
		if *mysqlInfo.Status == MYSQL_STATUS_DELIVING {
			return resource.RetryableError(fmt.Errorf("create mysql task  status is MYSQL_STATUS_DELIVING(%d)", MYSQL_STATUS_DELIVING))
		}
		if *mysqlInfo.Status == MYSQL_STATUS_RUNNING {
			return nil
		}
		err = fmt.Errorf("create mysql    task status is %v,we won't wait for it finish", *mysqlInfo.Status)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mysql  task fail, reason:%s\n ", logId, err.Error())
		return err
	}

	//internet service
	internetService := d.Get("internet_service").(int)
	if internetService == 1 {
		asyncRequestId, err := mysqlService.OpenWanService(ctx, d.Id())
		if err != nil {
			return err
		}
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			taskStatus, message, err := mysqlService.DescribeAsyncRequestInfo(ctx, asyncRequestId)
			if err != nil {
				if _, ok := err.(*errors.TencentCloudSDKError); !ok {
					return resource.RetryableError(err)
				} else {
					return resource.NonRetryableError(err)
				}
			}
			if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
				return nil
			}
			if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
				return resource.RetryableError(fmt.Errorf("create account task  status is %s", taskStatus))
			}
			err = fmt.Errorf("open internet service task status is %s,we won't wait for it finish ,it show message:%s", ",",
				message)
			return resource.NonRetryableError(err)
		})

		if err != nil {
			log.Printf("[CRITAL]%s open internet service   fail, reason:%s\n ", logId, err.Error())
			return err
		}
	}

	return resourceTencentCloudMysqlInstanceSetRead(d, meta)
}

func resourceTencentCloudMysqlInstanceSetRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	var mysqlInfo *cdb.InstanceInfo
	var e error
	var onlineHas = true
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		mysqlInfo, e = tencentMsyqlBasicInfoRead(ctx, d, meta, true)
		if e != nil {
			if mysqlService.NotFoundMysqlInstance(e) {
				d.SetId("")
				onlineHas = false
				return nil
			}
			return retryError(e)
		}
		if mysqlInfo == nil {
			d.SetId("")
			onlineHas = false
			return nil
		}
		_ = d.Set("project_id", int(*mysqlInfo.ProjectId))
		_ = d.Set("engine_version", mysqlInfo.EngineVersion)
		if *mysqlInfo.WanStatus == 1 {
			_ = d.Set("internet_service", 1)
			_ = d.Set("internet_host", mysqlInfo.WanDomain)
			_ = d.Set("internet_port", int(*mysqlInfo.WanPort))
		} else {
			_ = d.Set("internet_service", 0)
			_ = d.Set("internet_host", "")
			_ = d.Set("internet_port", 0)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Fail to get basic info from mysql, reaseon %s", err.Error())
	}
	if !onlineHas {
		return nil
	}
	parametersMap, ok := d.Get("parameters").(map[string]interface{})
	if !ok {
		log.Printf("[INFO] %v  config error,parameters is not map[string]interface{}\n", logId)
	} else {
		var cares []string
		for k := range parametersMap {
			cares = append(cares, k)
		}

		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			caresParameters, e := mysqlService.DescribeCaresParameters(ctx, d.Id(), cares)
			if e != nil {
				if mysqlService.NotFoundMysqlInstance(e) {
					d.SetId("")
					onlineHas = false
					return nil
				}
				return retryError(e)
			}
			if e := d.Set("parameters", caresParameters); e != nil {
				log.Printf("[CRITAL]%s provider set caresParameters fail, reason:%s\n ", logId, e.Error())
				return resource.NonRetryableError(e)
			}
			_ = d.Set("availability_zone", mysqlInfo.Zone)
			return nil
		})
		if err != nil {
			return fmt.Errorf("Describe CaresParameters Fail, reason:%s", err.Error())
		}
		if !onlineHas {
			return nil
		}
	}
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		backConfig, e := mysqlService.DescribeDBInstanceConfig(ctx, d.Id())
		if e != nil {
			if mysqlService.NotFoundMysqlInstance(e) {
				d.SetId("")
				onlineHas = false
				return nil
			}
			return retryError(e)
		}
		_ = d.Set("slave_sync_mode", int(*backConfig.Response.ProtectMode))
		_ = d.Set("slave_deploy_mode", int(*backConfig.Response.DeployMode))
		if backConfig.Response.SlaveConfig != nil && *backConfig.Response.SlaveConfig.Zone != "" {
			if _, ok := d.GetOk("first_slave_zone"); ok {
				_ = d.Set("first_slave_zone", *backConfig.Response.SlaveConfig.Zone)
			}
		}
		if backConfig.Response.BackupConfig != nil && *backConfig.Response.BackupConfig.Zone != "" {
			if _, ok := d.GetOk("second_slave_zone"); ok {
				_ = d.Set("second_slave_zone", *backConfig.Response.BackupConfig.Zone)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Describe DBInstanceConfig Fail, reason:%s", err.Error())
	}
	return nil
}

/*
   [master] and [dr] and [ro] all need update
*/

/*
 [master] need set
*/

func resourceTencentCloudMysqlInstanceSetUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_instance.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	payType := getPayType(d).(int)

	d.Partial(true)
	if payType == MysqlPayByMonth {
		err := mysqlUpdateInstancePayByMonth(ctx, d, meta)
		if err != nil {
			return err
		}
	} else if payType == MysqlPayByUse {
		err := mysqlUpdateInstancePayByUse(ctx, d, meta)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("mysql not support this pay type yet.")
	}
	d.Partial(false)
	time.Sleep(7 * time.Second)

	return resourceTencentCloudMysqlInstanceSetRead(d, meta)
}

func resourceTencentCloudMysqlInstanceSetDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_instance.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, err := mysqlService.IsolateDBInstance(ctx, d.Id())
		if err != nil {
			//for the pay order wait
			return retryError(err, InternalError)
		}
		return nil
	})

	if err != nil {
		return err
	}

	var hasDeleted = false

	payType := getPayType(d).(int)
	forceDelete := d.Get("force_delete").(bool)
	err = resource.Retry(7*readRetryTimeout, func() *resource.RetryError {
		mysqlInfo, err := mysqlService.DescribeDBInstanceById(ctx, d.Id())

		if err != nil {
			if _, ok := err.(*errors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if mysqlInfo == nil {
			hasDeleted = true
			return nil
		}
		if *mysqlInfo.Status == MYSQL_STATUS_ISOLATING || *mysqlInfo.Status == MYSQL_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("mysql isolating."))
		}
		if *mysqlInfo.Status == MYSQL_STATUS_ISOLATED {
			return nil
		}
		return resource.NonRetryableError(fmt.Errorf("after IsolateDBInstance mysql Status is %d", *mysqlInfo.Status))
	})

	if hasDeleted {
		return nil
	}
	if err != nil {
		return err
	}

	if payType == MysqlPayByMonth && !forceDelete {
		return nil
	}

	err = mysqlService.OfflineIsolatedInstances(ctx, d.Id())
	if err != nil {
		return err
	}

	err = resource.Retry(7*readRetryTimeout, func() *resource.RetryError {
		mysqlInfo, err := mysqlService.DescribeIsolatedDBInstanceById(ctx, d.Id())
		if err != nil {
			if _, ok := err.(*errors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if mysqlInfo == nil {
			return nil
		} else {
			if mysqlInfo.RoGroups != nil && len(mysqlInfo.RoGroups) > 0 {
				log.Printf("[WARN]this mysql has RoGroups , RoGroups is released asynchronously, and the bound resource is not now fully released now\n")
				return nil
			}
			return resource.RetryableError(fmt.Errorf("after OfflineIsolatedInstances mysql Status is %d", *mysqlInfo.Status))
		}
	})
	return err
}
