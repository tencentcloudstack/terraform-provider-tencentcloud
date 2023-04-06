/*
Provides a resource to create a redis account

Example Usage

```hcl
resource "tencentcloud_redis_account" "account" {
  instance_id 	   = "crs-xxxxxx"
  account_name 	   = "account_test"
  account_password = "test1234"
  remark 		   = "master"
  readonly_policy  = ["master"]
  privilege 	   = "rw"
}
```

Import

redis account can be imported using the id, e.g.

```
terraform import tencentcloud_redis_account.account crs-xxxxxx#account_test
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudRedisAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisAccountCreate,
		Read:   resourceTencentCloudRedisAccountRead,
		Update: resourceTencentCloudRedisAccountUpdate,
		Delete: resourceTencentCloudRedisAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"account_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The account name.",
			},

			"account_password": {
				Required:    true,
				Type:        schema.TypeString,
				Sensitive:   true,
				Description: "1: Length 8-30 digits, it is recommended to use a password of more than 12 digits; 2: Cannot start with `/`; 3: Include at least two items: a.Lowercase letters `a-z`; b.Uppercase letters `A-Z` c.Numbers `0-9`;  d.`()`~!@#$%^&*-+=_|{}[]:;<>,.?/`.",
			},

			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Remark.",
			},

			"readonly_policy": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Routing policy: Enter master or replication, which indicates the master node or slave node, cannot be empty when modifying operations.",
			},

			"privilege": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Read and write policy: Enter R and RW to indicate read-only, read-write, cannot be empty when modifying operations.",
			},
		},
	}
}

func resourceTencentCloudRedisAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_account.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request     = redis.NewCreateInstanceAccountRequest()
		response    = redis.NewCreateInstanceAccountResponse()
		instanceId  string
		accountName string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("account_name"); ok {
		accountName = v.(string)
		request.AccountName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("account_password"); ok {
		request.AccountPassword = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("readonly_policy"); ok {
		readonlyPolicySet := v.(*schema.Set).List()
		for i := range readonlyPolicySet {
			readonlyPolicy := readonlyPolicySet[i].(string)
			request.ReadonlyPolicy = append(request.ReadonlyPolicy, &readonlyPolicy)
		}
	}

	if v, ok := d.GetOk("privilege"); ok {
		request.Privilege = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().CreateInstanceAccount(request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.Code == "FailedOperation.SystemError" {
					return resource.NonRetryableError(e)
				}
			}
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create redis account failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + FILED_SP + accountName)

	taskId := *response.Response.TaskId
	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
	err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
		ok, err := service.DescribeTaskInfo(ctx, instanceId, taskId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("create account is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s redis create account fail, reason:%s\n", logId, err.Error())
		return err
	}

	conf := BuildStateChangeConf(
		[]string{},
		[]string{"2"},
		6*readRetryTimeout,
		time.Second,
		service.RedisAccountStateRefreshFunc(instanceId, accountName, []string{}),
	)

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudRedisAccountRead(d, meta)
}

func resourceTencentCloudRedisAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_account.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	accountName := idSplit[1]

	account, err := service.DescribeRedisAccountById(ctx, instanceId, accountName)
	if err != nil {
		return err
	}

	if account == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisAccount` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if account.InstanceId != nil {
		_ = d.Set("instance_id", account.InstanceId)
	}

	if account.AccountName != nil {
		_ = d.Set("account_name", account.AccountName)
	}

	// if account.AccountPassword != nil {
	// 	_ = d.Set("account_password", account.AccountPassword)
	// }

	if account.Remark != nil {
		_ = d.Set("remark", account.Remark)
	}

	if account.ReadonlyPolicy != nil {
		_ = d.Set("readonly_policy", account.ReadonlyPolicy)
	}

	if account.Privilege != nil {
		_ = d.Set("privilege", account.Privilege)
	}

	return nil
}

func resourceTencentCloudRedisAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_account.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := redis.NewModifyInstanceAccountRequest()
	response := redis.NewModifyInstanceAccountResponse()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	accountName := idSplit[1]

	request.InstanceId = &instanceId
	request.AccountName = &accountName

	immutableArgs := []string{"instance_id", "account_name"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("account_password") {
		if v, ok := d.GetOk("account_password"); ok {
			request.AccountPassword = helper.String(v.(string))
		}
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
	}

	if v, ok := d.GetOk("readonly_policy"); ok {
		readonlyPolicySet := v.(*schema.Set).List()
		for i := range readonlyPolicySet {
			readonlyPolicy := readonlyPolicySet[i].(string)
			request.ReadonlyPolicy = append(request.ReadonlyPolicy, &readonlyPolicy)
		}
	}

	if v, ok := d.GetOk("privilege"); ok {
		request.Privilege = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().ModifyInstanceAccount(request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.Code == "FailedOperation.SystemError" {
					return resource.NonRetryableError(e)
				}
			}
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update redis account failed, reason:%+v", logId, err)
		return err
	}

	taskId := *response.Response.TaskId
	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
	err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
		ok, err := service.DescribeTaskInfo(ctx, instanceId, taskId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("change account is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s redis change account fail, reason:%s\n", logId, err.Error())
		return err
	}

	conf := BuildStateChangeConf(
		[]string{},
		[]string{"2"},
		6*readRetryTimeout,
		time.Second,
		service.RedisAccountStateRefreshFunc(instanceId, accountName, []string{}),
	)

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudRedisAccountRead(d, meta)
}

func resourceTencentCloudRedisAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_account.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	accountName := idSplit[1]

	taskId, err := service.DeleteRedisAccountById(ctx, instanceId, accountName)
	if err != nil {
		return err
	}

	err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
		ok, err := service.DescribeTaskInfo(ctx, instanceId, taskId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("delete account is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s redis delete account fail, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
