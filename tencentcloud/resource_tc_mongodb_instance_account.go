/*
Provides a resource to create a mongodb instance_account

Example Usage

```hcl
resource "tencentcloud_mongodb_instance_account" "instance_account" {
  instance_id = "cmgo-lxaz2c9b"
  user_name = "test_account"
  password = "xxxxxxxx"
  mongo_user_password = "xxxxxxxxx"
  user_desc = "test account"
  auth_role {
    mask = 0
    namespace = "*"
  }
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMongodbInstanceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMongodbInstanceAccountCreate,
		Read:   resourceTencentCloudMongodbInstanceAccountRead,
		Update: resourceTencentCloudMongodbInstanceAccountUpdate,
		Delete: resourceTencentCloudMongodbInstanceAccountDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Instance ID, the format is: cmgo-9d0p6umb.Same as the instance ID displayed in the cloud database console page.",
			},

			"user_name": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "The new account name. Its format requirements are as follows: character range [1,32]. Characters in the range of [A,Z], [a,z], [1,9] as well as underscore _ and dash - can be input.",
			},

			"password": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "New account password. Password complexity requirements are as follows: character length range [8,32]. Contains at least letters, numbers and special characters (exclamation point!, at@, pound sign #, percent sign %, caret ^, asterisk *, parentheses (), underscore _).",
			},

			"mongo_user_password": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "The password corresponding to the mongouser account. mongouser is the system default account, which is the password set when creating an instance.",
			},

			"user_desc": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Account remarks.",
			},

			"auth_role": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "The read and write permission information of the account.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mask": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Permission information of the current account. 0: No permission. 1: read-only. 2: Write only. 3: Read and write.",
						},
						"namespace": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Refers to the name of the database with the current account permissions.*: Indicates all databases. db.name: Indicates the database of a specific name.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMongodbInstanceAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mongodb_instance_account.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mongodb.NewCreateAccountUserRequest()
		response   = mongodb.NewCreateAccountUserResponse()
		instanceId string
		userName   string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_name"); ok {
		userName = v.(string)
		request.UserName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	if v, ok := d.GetOk("mongo_user_password"); ok {
		request.MongoUserPassword = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_desc"); ok {
		request.UserDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("auth_role"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			auth := mongodb.Auth{}
			if v, ok := dMap["mask"]; ok {
				auth.Mask = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["namespace"]; ok {
				auth.NameSpace = helper.String(v.(string))
			}
			request.AuthRole = append(request.AuthRole, &auth)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMongodbClient().CreateAccountUser(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mongodb instanceAccount failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + FILED_SP + userName)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := MongodbService{client: meta.(*TencentCloudClient).apiV3Conn}

	if response != nil && response.Response != nil {
		if err = service.DescribeAsyncRequestInfo(ctx, helper.UInt64ToStr(*response.Response.FlowId), 3*readRetryTimeout); err != nil {
			return err
		}
	}

	return resourceTencentCloudMongodbInstanceAccountRead(d, meta)
}

func resourceTencentCloudMongodbInstanceAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mongodb_instance_account.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MongodbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userName := idSplit[1]

	instanceAccount, err := service.DescribeMongodbInstanceAccountById(ctx, instanceId, userName)
	if err != nil {
		return err
	}

	if instanceAccount == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MongodbInstanceAccount` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if instanceAccount.UserName != nil {
		_ = d.Set("user_name", instanceAccount.UserName)
	}

	if instanceAccount.UserDesc != nil {
		_ = d.Set("user_desc", instanceAccount.UserDesc)
	}

	if instanceAccount.AuthRole != nil {
		authRoleList := []interface{}{}
		for _, authRole := range instanceAccount.AuthRole {
			authRoleMap := map[string]interface{}{}

			if authRole.Mask != nil {
				authRoleMap["mask"] = authRole.Mask
			}

			if authRole.NameSpace != nil {
				authRoleMap["namespace"] = authRole.NameSpace
			}

			authRoleList = append(authRoleList, authRoleMap)
		}

		_ = d.Set("auth_role", authRoleList)

	}

	return nil
}

func resourceTencentCloudMongodbInstanceAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mongodb_instance_account.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mongodb.NewSetAccountUserPrivilegeRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userName := idSplit[1]

	request.InstanceId = &instanceId
	request.UserName = &userName

	immutableArgs := []string{"user_desc"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("auth_role") {
		if v, ok := d.GetOk("auth_role"); ok {
			for _, item := range v.([]interface{}) {
				auth := mongodb.Auth{}
				dMap := item.(map[string]interface{})
				if v, ok := dMap["mask"]; ok {
					auth.Mask = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["namespace"]; ok {
					auth.NameSpace = helper.String(v.(string))
				}
				request.AuthRole = append(request.AuthRole, &auth)
			}
		}

		var response *mongodb.SetAccountUserPrivilegeResponse
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseMongodbClient().SetAccountUserPrivilege(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mongodb instanceAccount failed, reason:%+v", logId, err)
			return err
		}

		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := MongodbService{client: meta.(*TencentCloudClient).apiV3Conn}

		if response != nil && response.Response != nil {
			if err = service.DescribeAsyncRequestInfo(ctx, helper.UInt64ToStr(*response.Response.FlowId), 3*readRetryTimeout); err != nil {
				return err
			}
		}
	}

	if d.HasChange("password") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := MongodbService{client: meta.(*TencentCloudClient).apiV3Conn}
		password := d.Get("password").(string)
		err := service.ResetInstancePassword(ctx, instanceId, userName, password)
		if err != nil {
			return err
		}

		d.SetPartial("password")
	}

	return resourceTencentCloudMongodbInstanceAccountRead(d, meta)
}

func resourceTencentCloudMongodbInstanceAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mongodb_instance_account.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MongodbService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userName := idSplit[1]

	var mongoUserPassword string
	if v, ok := d.GetOk("mongo_user_password"); ok {
		mongoUserPassword = v.(string)
	}

	if err := service.DeleteMongodbInstanceAccountById(ctx, instanceId, userName, mongoUserPassword); err != nil {
		return err
	}

	return nil
}
