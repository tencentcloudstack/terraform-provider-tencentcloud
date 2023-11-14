/*
Provides a resource to create a mongodb instance_property

Example Usage

```hcl
resource "tencentcloud_mongodb_instance_property" "instance_property" {
  instance_id = "cmgo-9d0p6umb"
  user_name = "test_account"
  password = "Abc@123..."
  mongo_user_password = "Abc@123."
  user_desc = "test account"
  auth_role {
		mask =
		name_space = ""

  }
}
```

Import

mongodb instance_property can be imported using the id, e.g.

```
terraform import tencentcloud_mongodb_instance_property.instance_property instance_property_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudMongodbInstanceProperty() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMongodbInstancePropertyCreate,
		Read:   resourceTencentCloudMongodbInstancePropertyRead,
		Update: resourceTencentCloudMongodbInstancePropertyUpdate,
		Delete: resourceTencentCloudMongodbInstancePropertyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, the format is: cmgo-9d0p6umb.Same as the instance ID displayed in the cloud database console page.",
			},

			"user_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The new account name. Its format requirements are as follows: character range [1,32]. Characters in the range of [A,Z], [a,z], [1,9] as well as underscore _ and dash - can be input.",
			},

			"password": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "New account password. Password complexity requirements are as follows: character length range [8,32]. Contains at least letters, numbers and special characters (exclamation point !, at@, pound sign #, percent sign %, caret ^, asterisk *, parentheses () , underscore _).",
			},

			"mongo_user_password": {
				Required:    true,
				Type:        schema.TypeString,
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
						"name_space": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Refers to the name of the database with the current account permissions.* : Indicates all databases. db.name: Indicates the database of a specific name.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMongodbInstancePropertyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mongodb_instance_property.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mongodb.NewCreateAccountUserRequest()
		response   = mongodb.NewCreateAccountUserResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_name"); ok {
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
			if v, ok := dMap["name_space"]; ok {
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
		log.Printf("[CRITAL]%s create mongodb instanceProperty failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudMongodbInstancePropertyRead(d, meta)
}

func resourceTencentCloudMongodbInstancePropertyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mongodb_instance_property.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MongodbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instancePropertyId := d.Id()

	instanceProperty, err := service.DescribeMongodbInstancePropertyById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instanceProperty == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MongodbInstanceProperty` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instanceProperty.InstanceId != nil {
		_ = d.Set("instance_id", instanceProperty.InstanceId)
	}

	if instanceProperty.UserName != nil {
		_ = d.Set("user_name", instanceProperty.UserName)
	}

	if instanceProperty.Password != nil {
		_ = d.Set("password", instanceProperty.Password)
	}

	if instanceProperty.MongoUserPassword != nil {
		_ = d.Set("mongo_user_password", instanceProperty.MongoUserPassword)
	}

	if instanceProperty.UserDesc != nil {
		_ = d.Set("user_desc", instanceProperty.UserDesc)
	}

	if instanceProperty.AuthRole != nil {
		authRoleList := []interface{}{}
		for _, authRole := range instanceProperty.AuthRole {
			authRoleMap := map[string]interface{}{}

			if instanceProperty.AuthRole.Mask != nil {
				authRoleMap["mask"] = instanceProperty.AuthRole.Mask
			}

			if instanceProperty.AuthRole.NameSpace != nil {
				authRoleMap["name_space"] = instanceProperty.AuthRole.NameSpace
			}

			authRoleList = append(authRoleList, authRoleMap)
		}

		_ = d.Set("auth_role", authRoleList)

	}

	return nil
}

func resourceTencentCloudMongodbInstancePropertyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mongodb_instance_property.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mongodb.NewSetAccountUserPrivilegeRequest()

	instancePropertyId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "user_name", "password", "mongo_user_password", "user_desc", "auth_role"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("instance_id") {
		if v, ok := d.GetOk("instance_id"); ok {
			request.InstanceId = helper.String(v.(string))
		}
	}

	if d.HasChange("user_name") {
		if v, ok := d.GetOk("user_name"); ok {
			request.UserName = helper.String(v.(string))
		}
	}

	if d.HasChange("auth_role") {
		if v, ok := d.GetOk("auth_role"); ok {
			for _, item := range v.([]interface{}) {
				auth := mongodb.Auth{}
				if v, ok := dMap["mask"]; ok {
					auth.Mask = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["name_space"]; ok {
					auth.NameSpace = helper.String(v.(string))
				}
				request.AuthRole = append(request.AuthRole, &auth)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMongodbClient().SetAccountUserPrivilege(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mongodb instanceProperty failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMongodbInstancePropertyRead(d, meta)
}

func resourceTencentCloudMongodbInstancePropertyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mongodb_instance_property.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MongodbService{client: meta.(*TencentCloudClient).apiV3Conn}
	instancePropertyId := d.Id()

	if err := service.DeleteMongodbInstancePropertyById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
