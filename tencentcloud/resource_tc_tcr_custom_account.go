/*
Provides a resource to create a tcr custom_account

Example Usage

Create custom account with specified duration days

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr-instance"
  instance_type = "basic"
  delete_bucket = true
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id    = tencentcloud_tcr_instance.example.id
  name           = "tf_test_tcr_namespace"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "tf_example_cve_id"
  }
}

resource "tencentcloud_tcr_custom_account" "example" {
  registry_id = tencentcloud_tcr_instance.example.id
  name        = "tf_example_account"
  permissions {
    resource = tencentcloud_tcr_namespace.example.name
    actions  = ["tcr:PushRepository", "tcr:PullRepository"]
  }
  description = "tf example for tcr custom account"
  duration    = 10
  disable     = false
  tags = {
    "createdBy" = "terraform"
  }
}
```

With specified expiration time

```hcl
resource "tencentcloud_tcr_custom_account" "example" {
  registry_id = tencentcloud_tcr_instance.example.id
  name        = "tf_example_account"
  permissions {
    resource = tencentcloud_tcr_namespace.example.name
    actions  = ["tcr:PushRepository", "tcr:PullRepository"]
  }
  description = "tf example for tcr custom account"
  expires_at  = 1676897989000 //time stamp
  disable     = false
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tcr custom_account can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_custom_account.custom_account custom_account_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTcrCustomAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrCustomAccountCreate,
		Read:   resourceTencentCloudTcrCustomAccountRead,
		Update: resourceTencentCloudTcrCustomAccountUpdate,
		Delete: resourceTencentCloudTcrCustomAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "custom account name.",
			},

			"permissions": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "strategy list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "resource path, currently only supports Namespace. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"actions": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Actions, currently only support: `tcr:PushRepository`, `tcr:PullRepository`. Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "custom account description.",
			},

			"duration": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "expiration date (unit: day), calculated from the current time, priority is higher than `ExpiresAt`.",
			},

			"expires_at": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "custom account expiration time (time stamp, unit: milliseconds).",
			},

			"disable": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "whether to disable custom accounts.",
			},

			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Password of the account.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTcrCustomAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_custom_account.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = tcr.NewCreateCustomAccountRequest()
		response   = tcr.NewCreateCustomAccountResponse()
		registryId string
	)
	if v, ok := d.GetOk("registry_id"); ok {
		request.RegistryId = helper.String(v.(string))
		registryId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("permissions"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			permission := tcr.Permission{}
			if v, ok := dMap["resource"]; ok {
				permission.Resource = helper.String(v.(string))
			}
			if v, ok := dMap["actions"]; ok {
				actionsSet := v.(*schema.Set).List()
				for i := range actionsSet {
					if actionsSet[i] != nil {
						actions := actionsSet[i].(string)
						permission.Actions = append(permission.Actions, &actions)
					}
				}
			}
			request.Permissions = append(request.Permissions, &permission)
		}
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("duration"); ok {
		request.Duration = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("expires_at"); ok {
		request.ExpiresAt = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("disable"); ok {
		request.Disable = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTCRClient().CreateCustomAccount(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tcr CustomAccount failed, reason:%+v", logId, err)
		return err
	}

	name := *response.Response.Name
	d.SetId(strings.Join([]string{registryId, name}, FILED_SP))

	pw := response.Response.Password
	if pw != nil {
		_ = d.Set("password", *pw)
	}

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::tcr:%s:uin/:instance/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTcrCustomAccountRead(d, meta)
}

func resourceTencentCloudTcrCustomAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_custom_account.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	name := idSplit[1]

	CustomAccount, err := service.DescribeTcrCustomAccountById(ctx, registryId, name)
	if err != nil {
		return err
	}

	if CustomAccount == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TcrCustomAccount` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("registry_id", registryId)

	if CustomAccount.Name != nil {
		_ = d.Set("name", CustomAccount.Name)
	}

	if CustomAccount.Permissions != nil {
		permissionsList := []interface{}{}
		for _, permission := range CustomAccount.Permissions {
			permissionsMap := map[string]interface{}{}

			if permission.Resource != nil {
				permissionsMap["resource"] = permission.Resource
			}

			if len(permission.Actions) > 0 {
				permissionsMap["actions"] = helper.StringsInterfaces(permission.Actions)
			}

			permissionsList = append(permissionsList, permissionsMap)
		}

		_ = d.Set("permissions", permissionsList)

	}

	if CustomAccount.Description != nil {
		_ = d.Set("description", CustomAccount.Description)
	}

	if CustomAccount.ExpiresAt != nil {
		_ = d.Set("expires_at", CustomAccount.ExpiresAt)
	}

	if CustomAccount.Disable != nil {
		_ = d.Set("disable", CustomAccount.Disable)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "tcr", "instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTcrCustomAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_custom_account.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tcr.NewModifyCustomAccountRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	name := idSplit[1]

	request.RegistryId = &registryId
	request.Name = &name

	immutableArgs := []string{"registry_id", "name", "permissions", "description", "duration", "expires_at", "disable"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("registry_id") {
		if v, ok := d.GetOk("registry_id"); ok {
			request.RegistryId = helper.String(v.(string))
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("permissions") {
		if v, ok := d.GetOk("permissions"); ok {
			for _, item := range v.([]interface{}) {
				permission := tcr.Permission{}
				dMap := item.(map[string]interface{})
				if v, ok := dMap["resource"]; ok {
					permission.Resource = helper.String(v.(string))
				}
				if v, ok := dMap["actions"]; ok {
					actionsSet := v.(*schema.Set).List()
					for i := range actionsSet {
						if actionsSet[i] != nil {
							actions := actionsSet[i].(string)
							permission.Actions = append(permission.Actions, &actions)
						}
					}
				}
				request.Permissions = append(request.Permissions, &permission)
			}
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("duration") {
		if v, ok := d.GetOkExists("duration"); ok {
			request.Duration = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("expires_at") {
		if v, ok := d.GetOkExists("expires_at"); ok {
			request.ExpiresAt = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("disable") {
		if v, ok := d.GetOkExists("disable"); ok {
			request.Disable = helper.Bool(v.(bool))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTCRClient().ModifyCustomAccount(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tcr CustomAccount failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("tcr", "instance", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTcrCustomAccountRead(d, meta)
}

func resourceTencentCloudTcrCustomAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_custom_account.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	name := idSplit[1]

	if err := service.DeleteTcrCustomAccountById(ctx, registryId, name); err != nil {
		return err
	}

	return nil
}
