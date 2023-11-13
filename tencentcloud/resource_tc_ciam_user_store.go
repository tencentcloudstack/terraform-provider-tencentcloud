/*
Provides a resource to create a ciam user_store

Example Usage

```hcl
resource "tencentcloud_ciam_user_store" "user_store" {
  user_pool_name = ""
  user_pool_desc = ""
  user_pool_logo = ""
}
```

Import

ciam user_store can be imported using the id, e.g.

```
terraform import tencentcloud_ciam_user_store.user_store user_store_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ciam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ciam/v20220331"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCiamUserStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiamUserStoreCreate,
		Read:   resourceTencentCloudCiamUserStoreRead,
		Update: resourceTencentCloudCiamUserStoreUpdate,
		Delete: resourceTencentCloudCiamUserStoreDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_pool_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "User Store Name.",
			},

			"user_pool_desc": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User Store Description.",
			},

			"user_pool_logo": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User Store Logo.",
			},
		},
	}
}

func resourceTencentCloudCiamUserStoreCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ciam_user_store.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = ciam.NewCreateUserStoreRequest()
		response    = ciam.NewCreateUserStoreResponse()
		userStoreId string
	)
	if v, ok := d.GetOk("user_pool_name"); ok {
		request.UserPoolName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_pool_desc"); ok {
		request.UserPoolDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_pool_logo"); ok {
		request.UserPoolLogo = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiamClient().CreateUserStore(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ciam userStore failed, reason:%+v", logId, err)
		return err
	}

	userStoreId = *response.Response.UserStoreId
	d.SetId(userStoreId)

	return resourceTencentCloudCiamUserStoreRead(d, meta)
}

func resourceTencentCloudCiamUserStoreRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ciam_user_store.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiamService{client: meta.(*TencentCloudClient).apiV3Conn}

	userStoreId := d.Id()

	userStore, err := service.DescribeCiamUserStoreById(ctx, userStoreId)
	if err != nil {
		return err
	}

	if userStore == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CiamUserStore` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if userStore.UserPoolName != nil {
		_ = d.Set("user_pool_name", userStore.UserPoolName)
	}

	if userStore.UserPoolDesc != nil {
		_ = d.Set("user_pool_desc", userStore.UserPoolDesc)
	}

	if userStore.UserPoolLogo != nil {
		_ = d.Set("user_pool_logo", userStore.UserPoolLogo)
	}

	return nil
}

func resourceTencentCloudCiamUserStoreUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ciam_user_store.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := ciam.NewUpdateUserStoreRequest()

	userStoreId := d.Id()

	request.UserStoreId = &userStoreId

	immutableArgs := []string{"user_pool_name", "user_pool_desc", "user_pool_logo"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("user_pool_name") {
		if v, ok := d.GetOk("user_pool_name"); ok {
			request.UserPoolName = helper.String(v.(string))
		}
	}

	if d.HasChange("user_pool_desc") {
		if v, ok := d.GetOk("user_pool_desc"); ok {
			request.UserPoolDesc = helper.String(v.(string))
		}
	}

	if d.HasChange("user_pool_logo") {
		if v, ok := d.GetOk("user_pool_logo"); ok {
			request.UserPoolLogo = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiamClient().UpdateUserStore(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update ciam userStore failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiamUserStoreRead(d, meta)
}

func resourceTencentCloudCiamUserStoreDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ciam_user_store.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiamService{client: meta.(*TencentCloudClient).apiV3Conn}
	userStoreId := d.Id()

	if err := service.DeleteCiamUserStoreById(ctx, userStoreId); err != nil {
		return err
	}

	return nil
}
