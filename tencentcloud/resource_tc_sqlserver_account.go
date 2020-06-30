/*
Use this resource to create SQL Server account

Example Usage

```hcl
resource "tencentcloud_sqlserver_account" "foo" {
  name = "example"
  availability_zone = var.availability_zone
  charge_type = "POSTPAID_BY_HOUR"
  vpc_id      = "vpc-409mvdvv"
  subnet_id = "subnet-nf9n81ps"
  engine_version		= "9.3.5"
  root_password                 = "1qaA2k1wgvfa3ZZZ"
  charset = "UTF8"
  project_id = 0
  memory = 2
  storage = 100
}
```

Import

sqlserver account can be imported using the id, e.g.

```
$ terraform import tencentcloud_sqlserver_account.foo postgres-cda1iex1
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudSqlserverAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverAccountCreate,
		Read:   resourceTencentCloudSqlserverAccountRead,
		Update: resourceTencentCloudSqlserverAccountUpdate,
		Delete: resourceTencentCLoudSqlserverAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the SQL Server account.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password of the SQL Server account.",
			},
			"is_admin": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate that the account is root account or not.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Remark of the SQL Server account.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Instance ID that the account belongs to.",
			},
			//computed
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Status of the SQL Server account. 1 for creating, 2 for running, 3 for modifying, 4 for resetting password, -1 for deleting.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the SQL Server account.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last updated time of the SQL Server account.",
			},
		},
	}
}

func resourceTencentCloudSqlserverAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_account.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		name       = d.Get("name").(string)
		password   = d.Get("password").(string)
		remark     = d.Get("remark").(string)
		isAdmin    = d.Get("is_admin").(bool)
		instanceId = d.Get("instance_id").(string)
	)

	var outErr, inErr error

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		inErr = sqlserverService.CreateSqlserverAccount(ctx, instanceId, name, password, remark, isAdmin)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(instanceId + FILED_SP + name)

	return resourceTencentCloudSqlserverAccountRead(d, meta)
}

func resourceTencentCloudSqlserverAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_account.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	id := d.Id()
	idStrs := strings.Split(id, FILED_SP)
	if len(idStrs) != 2 {
		return fmt.Errorf("invalid SQL Server account id %s", id)
	}
	instanceId := idStrs[0]
	name := idStrs[1]
	d.Partial(true)

	var outErr, inErr error

	//update is_admin
	if d.HasChange("is_admin") {
		return fmt.Errorf("is_admin is not allowed to change")
	}

	//update password
	if d.HasChange("password") {
		password := d.Get("password").(string)
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = sqlserverService.ResetSqlserverAccountPassword(ctx, instanceId, name, password)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

		d.SetPartial("password")
	}

	//update remark
	if d.HasChange("remark") {
		remark := d.Get("remark").(string)
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = sqlserverService.ModifySqlserverAccountRemark(ctx, instanceId, name, remark)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

		d.SetPartial("remark")
	}

	d.Partial(false)

	return resourceTencentCloudSqlserverAccountRead(d, meta)
}

func resourceTencentCloudSqlserverAccountRead(d *schema.ResourceData, meta interface{}) error {

	defer logElapsed("resource.tencentcloud_sqlserver_account.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()
	idStrs := strings.Split(id, FILED_SP)
	if len(idStrs) != 2 {
		return fmt.Errorf("invalid SQL Server account id %s", id)
	}
	instanceId := idStrs[0]
	name := idStrs[1]

	var outErr, inErr error
	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	account, has, outErr := sqlserverService.DescribeSqlserverAccountById(ctx, instanceId, name)
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			account, has, inErr = sqlserverService.DescribeSqlserverAccountById(ctx, instanceId, name)
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

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("name", account.Name)
	if *account.Remark == "--" {
		_ = d.Set("remark", "")
	} else {
		_ = d.Set("remark", account.Remark)
	}
	_ = d.Set("create_time", account.CreateTime)
	_ = d.Set("update_time", account.UpdateTime)
	_ = d.Set("status", account.Status)

	return nil
}

func resourceTencentCLoudSqlserverAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_account.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()
	idStrs := strings.Split(id, FILED_SP)
	if len(idStrs) != 2 {
		return fmt.Errorf("invalid SQL Server account id %s", id)
	}
	instanceId := idStrs[0]
	name := idStrs[1]

	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	var outErr, inErr error
	var has bool

	outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, inErr = sqlserverService.DescribeSqlserverAccountById(ctx, instanceId, name)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}

	if !has {
		return nil
	}

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		inErr = sqlserverService.DeleteSqlserverAccount(ctx, instanceId, name)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}

	outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, inErr = sqlserverService.DescribeSqlserverAccountById(ctx, instanceId, name)
		if inErr != nil {
			return retryError(inErr)
		}
		if has {
			inErr = fmt.Errorf("delete SQL Server account %s fail, account still exists from SDK DescribeSqlserverAccountById", id)
			return resource.RetryableError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}
	return nil
}
