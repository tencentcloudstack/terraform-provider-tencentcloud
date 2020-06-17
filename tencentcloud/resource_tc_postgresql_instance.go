/*
Use this resource to create postgresql instance

Example Usage

```hcl
resource "tencentcloud_postgresql_instance" "foo" {
  name = "example"
  availability_zone = var.availability_zone
  charge_type = "POSTPAID_BY_HOUR"
  vpc_id      = "vpc-409mvdvv"
  subnet_id = "subnet-nf9n81ps"
  engine_version		= "9.3.5"
  root_password                 = "1qaA2k1wgvfa3ZZZ"
  charset = "UTF8"
  spec_code = "cdb.pg.z1.2g"
  project_id = 0
  memory = 2
  storage = 100
}

Import

postgresql instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_postgresql_instance.foo postgres-cda1iex1
```
*/
package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Type:         schema.TypeString,
				Optional:     true,
				Default:      POSTGRESQL_PAYTYPE_POSTPAID,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(POSTGRESQL_PAYTYPE),
				Description:  "Pay type of the postgresql instance. For now, only `POSTPAID_BY_HOUR` is valid.",
			},
			"engine_version": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(POSTSQL_DB_VERSION),
				Default:      POSTSQL_DB_VERSION[len(POSTSQL_DB_VERSION)-1],
				Description:  "Version of the postgresql database engine. Allowed values are `9.3.5`, `9.5.4`, `10.4`.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "ID of VPC.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "ID of subnet.",
			},
			"storage": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Disk size (in GB). Allowed value is range from 10 to 1000, and the value must be a multiple of 10.",
			},
			"memory": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Memory size (in MB).",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Project ID, default value is 0.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "Availability zone.",
			},
			"spec_code": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The id of specification of the postgresql instance, like `cdb.pg.z1.2g`, which can be queried with data source `tencentcloud_postgresql_specinfos`.",
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
				ValidateFunc: validateAllowedStringValue(POSTSQL_DB_CHARSET),
				Description:  "Charset of the root account. Valid values are `UTF8`,`LATIN1`.",
			},
			"public_access_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether to enable the access to an instance from public network or not.",
			},
			//Computed values
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
				Description: "Ip for private access.",
			},
			"private_access_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Port for private access.",
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
		name      = d.Get("name").(string)
		dbVersion = d.Get("engine_version").(string)
		payType   = d.Get("charge_type").(string)
		specCode  = d.Get("spec_code").(string)
		projectId = d.Get("project_id").(int)
		subnetId  = d.Get("subnet_id").(string)
		vpcId     = d.Get("vpc_id").(string)
		zone      = d.Get("availability_zone").(string)
		storage   = d.Get("storage").(int)
	)

	var period = 1
	//the sdk asks to set value with 1 when paytype is postpaid

	if payType == COMMON_PAYTYPE_PREPAID {
		payType = POSTGRESQL_PAYTYPE_PREPAID
	} else {
		payType = POSTGRESQL_PAYTYPE_POSTPAID
	}

	var instanceId string
	var outErr, inErr error

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		instanceId, inErr = postgresqlService.CreatePostgresqlInstance(ctx, name, dbVersion, payType, specCode, 0, projectId, period, subnetId, vpcId, zone, storage)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	//check creation done
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instance, has, err := postgresqlService.DescribePostgresqlInstanceById(ctx, instanceId)
		if err != nil {
			return retryError(err)
		} else if has && *instance.DBInstanceStatus == "init" {
			return nil
		} else if !has {
			return resource.NonRetryableError(fmt.Errorf("create postgresql instance fail"))
		} else {
			return resource.RetryableError(fmt.Errorf("creating postgresql instance %s , status %s ", instanceId, *instance.DBInstanceStatus))
		}
	})

	if err != nil {
		return err
	}
	d.SetId(instanceId)

	var (
		password = d.Get("root_password").(string)
		charset  = d.Get("charset").(string)
	)

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		inErr = postgresqlService.InitPostgresqlInstance(ctx, instanceId, password, charset)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}
	//check init status

	//set open public access
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

	//set name
	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		inErr := postgresqlService.ModifyPostgresqlInstanceName(ctx, instanceId, d.Get("name").(string))
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	//check creation done
	checkErr := postgresqlService.CheckDBInstanceStatus(ctx, instanceId)
	if checkErr != nil {
		return checkErr
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

	var outErr, inErr, checkErr error
	//update name
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
		//check update public service done
		checkErr = postgresqlService.CheckDBInstanceStatus(ctx, instanceId)
		if checkErr != nil {
			return checkErr
		}
		d.SetPartial("name")
	}

	//upgrade storage and memory size
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
		//check update public service done
		checkErr = postgresqlService.CheckDBInstanceStatus(ctx, instanceId)
		if checkErr != nil {
			return checkErr
		}
		d.SetPartial("memory")
		d.SetPartial("storage")
	}

	//update project id
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

		//check update project id done
		checkErr = postgresqlService.CheckDBInstanceStatus(ctx, instanceId)
		if checkErr != nil {
			return checkErr
		}
		d.SetPartial("project_id")
	}

	//update public access
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
		//check update public service done
		checkErr = postgresqlService.CheckDBInstanceStatus(ctx, instanceId)
		if checkErr != nil {
			return checkErr
		}
		d.SetPartial("public_access_switch")
	}

	//update root password
	if d.HasChange("root_password") {
		//to avoid other updating process conflicts with updating password, set the password updating with the last step, there is no way to figure out whether changing password is done
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
		//check update password done
		checkErr = postgresqlService.CheckDBInstanceStatus(ctx, instanceId)
		if checkErr != nil {
			return checkErr
		}
		d.SetPartial("root_password")
	}

	d.Partial(false)

	return resourceTencentCloudPostgresqlInstanceRead(d, meta)
}

func resourceTencentCloudPostgresqlInstanceRead(d *schema.ResourceData, meta interface{}) error {

	defer logElapsed("resource.tencentcloud_postgresql_instance.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var outErr, inErr error
	postgresqlService := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	instance, has, outErr := postgresqlService.DescribePostgresqlInstanceById(ctx, d.Id())
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			instance, has, inErr = postgresqlService.DescribePostgresqlInstanceById(ctx, d.Id())
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
	}
	if outErr != nil {
		return outErr
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("project_id", int(*instance.ProjectId))
	_ = d.Set("availability_zone", instance.Zone)
	_ = d.Set("vpc_id", instance.VpcId)
	_ = d.Set("subnet_id", instance.SubnetId)
	_ = d.Set("engine_version", instance.DBVersion)
	_ = d.Set("name", instance.DBInstanceName)
	_ = d.Set("charset", instance.DBCharset)

	if *instance.PayType == POSTGRESQL_PAYTYPE_PREPAID {
		_ = d.Set("charge_type", COMMON_PAYTYPE_PREPAID)
	} else {
		_ = d.Set("charge_type", COMMON_PAYTYPE_POSTPAID)
	}
	//net status
	public_access_switch := false
	if len(instance.DBInstanceNetInfo) > 0 {
		for _, v := range instance.DBInstanceNetInfo {

			if *v.NetType == "public" {
				//both 1 and opened used in SDK
				if *v.Status == "opened" || *v.Status == "1" {
					public_access_switch = true
				}
				_ = d.Set("public_access_host", v.Address)
				_ = d.Set("public_access_port", v.Port)
			}
			if *v.NetType == "private" {
				_ = d.Set("private_access_ip", v.Ip)
				_ = d.Set("private_access_port", v.Port)
			}
		}
	}
	_ = d.Set("public_access_switch", public_access_switch)

	//computed
	_ = d.Set("create_time", instance.CreateTime)
	_ = d.Set("status", instance.DBInstanceStatus)
	_ = d.Set("memory", instance.DBInstanceMemory)
	_ = d.Set("storage", instance.DBInstanceStorage)

	//ignore spec_code
	return nil
}

func resourceTencentCLoudPostgresqlInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_instance.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Id()
	postgresqlService := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	var outErr, inErr, checkErr error
	//check status
	checkErr = postgresqlService.CheckDBInstanceStatus(ctx, instanceId)

	if checkErr != nil {
		return checkErr
	}

	outErr = postgresqlService.DeletePostgresqlInstance(ctx, instanceId)

	if outErr != nil {
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = postgresqlService.DeletePostgresqlInstance(ctx, instanceId)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
	}

	if outErr != nil {
		return outErr
	}

	_, has, outErr := postgresqlService.DescribePostgresqlInstanceById(ctx, instanceId)
	if outErr != nil || has {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, has, inErr = postgresqlService.DescribePostgresqlInstanceById(ctx, d.Id())
			if inErr != nil {
				return retryError(inErr)
			}
			if has {
				inErr = fmt.Errorf("delete postgresql instance %s fail, instance still exists from SDK DescribePostgresqlInstanceById", instanceId)
				return resource.RetryableError(inErr)
			}
			return nil
		})
	}
	if outErr != nil {
		return outErr
	}
	return nil
}
