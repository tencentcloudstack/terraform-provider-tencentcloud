/*
Use this resource to create tcr long term token.

Example Usage

```hcl
resource "tencentcloud_tcr_token" "foo" {
  instance_id		= "cls-cda1iex1"
  description		= "test"
}
```

Import

tcr token can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcr_token.foo cls-cda1iex1#namespace#buv3h3j96j2d1rk1cllg
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

func resourceTencentCloudTcrToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrTokenCreate,
		Read:   resourceTencentCloudTcrTokenRead,
		Update: resourceTencentCloudTcrTokenUpdate,
		Delete: resourceTencentCLoudTcrTokenDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the TCR instance.",
			},
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicate to enable this token or not.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Description of the token. Valid length is [0~255].",
			},
			//computed
			"token_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Sub ID of the TCR token. The full ID of token format like `instance_id#token_id`.",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content of the token.",
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User name of the token.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time.",
			},
		},
	}
}

func resourceTencentCloudTcrTokenCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_token.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		instanceId               = d.Get("instance_id").(string)
		enable                   = d.Get("enable").(bool)
		description              = d.Get("description").(string)
		tokenId, token, userName string
		outErr, inErr            error
	)

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		tokenId, token, userName, inErr = tcrService.CreateTCRLongTermToken(ctx, instanceId, description)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(instanceId + FILED_SP + tokenId)
	_ = d.Set("token", token)
	_ = d.Set("user_name", userName)

	if !enable {
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = tcrService.ModifyTCRLongTermToken(ctx, instanceId, tokenId, enable)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}
	}

	return resourceTencentCloudTcrTokenRead(d, meta)
}

func resourceTencentCloudTcrTokenUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_token.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	resourceId := d.Id()
	items := strings.Split(resourceId, FILED_SP)
	if len(items) != 2 {
		return fmt.Errorf("invalid ID %s", resourceId)
	}

	instanceId := items[0]
	tokenId := items[1]

	if d.HasChange("enable") {
		enable := d.Get("enable").(bool)
		var outErr, inErr error
		tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}
		outErr = tcrService.ModifyTCRLongTermToken(ctx, instanceId, tokenId, enable)
		if outErr != nil {
			outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				inErr = tcrService.ModifyTCRLongTermToken(ctx, instanceId, tokenId, enable)
				if inErr != nil {
					return retryError(inErr)
				}
				return nil
			})
		}
		if outErr != nil {
			return outErr
		}
	}

	return resourceTencentCloudTcrTokenRead(d, meta)
}

func resourceTencentCloudTcrTokenRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_token.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	resourceId := d.Id()
	items := strings.Split(resourceId, FILED_SP)
	if len(items) != 2 {
		return fmt.Errorf("invalid ID %s", resourceId)
	}

	instanceId := items[0]
	tokenId := items[1]

	var outErr, inErr error
	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}
	token, has, outErr := tcrService.DescribeTCRLongTermTokenById(ctx, instanceId, tokenId)
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			token, has, inErr = tcrService.DescribeTCRLongTermTokenById(ctx, instanceId, tokenId)
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

	_ = d.Set("token_id", tokenId)
	_ = d.Set("create_time", token.CreatedAt)
	_ = d.Set("description", token.Desc)
	_ = d.Set("enable", token.Enabled)
	_ = d.Set("instance_id", instanceId)

	return nil
}

func resourceTencentCLoudTcrTokenDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_token.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	resourceId := d.Id()
	items := strings.Split(resourceId, FILED_SP)
	if len(items) != 2 {
		return fmt.Errorf("invalid ID %s", resourceId)
	}

	instanceId := items[0]
	tokenId := items[1]

	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}

	var inErr, outErr error
	var has bool

	outErr = tcrService.DeleteTCRLongTermToken(ctx, instanceId, tokenId)
	if outErr != nil {
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = tcrService.DeleteTCRLongTermToken(ctx, instanceId, tokenId)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
	}

	if outErr != nil {
		return outErr
	}

	outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, inErr = tcrService.DescribeTCRLongTermTokenById(ctx, instanceId, tokenId)
		if inErr != nil {
			return retryError(inErr)
		}
		if has {
			inErr = fmt.Errorf("delete tcr namespace %s fail, namespace still exists from SDK DescribeTcrNamespaceById", resourceId)
			return resource.RetryableError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}

	return nil
}
