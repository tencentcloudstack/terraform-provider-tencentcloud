/*
Provides a resource to create a bi embed_token

Example Usage

```hcl
resource "tencentcloud_bi_embed_token" "embed_token" {
  project_id = 123
  page_id = 123
  scope = "page"
  expire_time = "240"
  extra_param = ""
  user_corp_id = "abc"
  user_id = "abc"
}
```

Import

bi embed_token can be imported using the id, e.g.

```
terraform import tencentcloud_bi_embed_token.embed_token embed_token_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bi "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bi/v20220105"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudBiEmbedToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBiEmbedTokenCreate,
		Read:   resourceTencentCloudBiEmbedTokenRead,
		Update: resourceTencentCloudBiEmbedTokenUpdate,
		Delete: resourceTencentCloudBiEmbedTokenDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Share project id.",
			},

			"page_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Sharing page id, this is empty value 0 when embedding the board.",
			},

			"scope": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Page means embedding the page, and panel means embedding the entire board.",
			},

			"expire_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Expiration. Unit: Minutes Maximum value: 240. i.e. 4 hours Default: 240.",
			},

			"extra_param": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Alternate fields.",
			},

			"user_corp_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User enterprise ID (for multi-user only).",
			},

			"user_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "UserId (for multi-user only).",
			},
		},
	}
}

func resourceTencentCloudBiEmbedTokenCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_embed_token.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = bi.NewCreateEmbedTokenRequest()
		response = bi.NewCreateEmbedTokenResponse()
		pageId   int
	)
	if v, ok := d.GetOkExists("project_id"); ok {
		request.ProjectId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("page_id"); ok {
		pageId = v.(int)
		request.PageId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("scope"); ok {
		request.Scope = helper.String(v.(string))
	}

	if v, ok := d.GetOk("expire_time"); ok {
		request.ExpireTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("extra_param"); ok {
		request.ExtraParam = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_corp_id"); ok {
		request.UserCorpId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_id"); ok {
		request.UserId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseBiClient().CreateEmbedToken(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create bi embedToken failed, reason:%+v", logId, err)
		return err
	}

	pageId = *response.Response.PageId
	d.SetId(helper.Int64ToStr(int64(pageId)))

	return resourceTencentCloudBiEmbedTokenRead(d, meta)
}

func resourceTencentCloudBiEmbedTokenRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_embed_token.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := BiService{client: meta.(*TencentCloudClient).apiV3Conn}

	embedTokenId := d.Id()

	embedToken, err := service.DescribeBiEmbedTokenById(ctx, pageId)
	if err != nil {
		return err
	}

	if embedToken == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `BiEmbedToken` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if embedToken.ProjectId != nil {
		_ = d.Set("project_id", embedToken.ProjectId)
	}

	if embedToken.PageId != nil {
		_ = d.Set("page_id", embedToken.PageId)
	}

	if embedToken.Scope != nil {
		_ = d.Set("scope", embedToken.Scope)
	}

	if embedToken.ExpireTime != nil {
		_ = d.Set("expire_time", embedToken.ExpireTime)
	}

	if embedToken.ExtraParam != nil {
		_ = d.Set("extra_param", embedToken.ExtraParam)
	}

	if embedToken.UserCorpId != nil {
		_ = d.Set("user_corp_id", embedToken.UserCorpId)
	}

	if embedToken.UserId != nil {
		_ = d.Set("user_id", embedToken.UserId)
	}

	return nil
}

func resourceTencentCloudBiEmbedTokenUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_embed_token.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"project_id", "page_id", "scope", "expire_time", "extra_param", "user_corp_id", "user_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudBiEmbedTokenRead(d, meta)
}

func resourceTencentCloudBiEmbedTokenDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_embed_token.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := BiService{client: meta.(*TencentCloudClient).apiV3Conn}
	embedTokenId := d.Id()

	if err := service.DeleteBiEmbedTokenById(ctx, pageId); err != nil {
		return err
	}

	return nil
}
