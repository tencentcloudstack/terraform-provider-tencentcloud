/*
Provides a resource to create a bi embed_interval

Example Usage

```hcl
resource "tencentcloud_bi_embed_interval_apply" "embed_interval" {
  project_id = 11015030
  page_id    = 10520483
  bi_token   = "4192d65b-d674-4117-9a59-xxxxxxxxx"
  scope      = "page"
}
```
*/
package tencentcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bi "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bi/v20220105"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudBiEmbedIntervalApply() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBiEmbedIntervalApplyCreate,
		Read:   resourceTencentCloudBiEmbedIntervalApplyRead,
		Delete: resourceTencentCloudBiEmbedIntervalApplyDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Sharing project id, required.",
			},

			"page_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Sharing page id, this is empty value 0 when embedding the board.",
			},

			"bi_token": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Token that needs to be applied for extension.",
			},

			"scope": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Choose panel or page.",
			},
		},
	}
}

func resourceTencentCloudBiEmbedIntervalApplyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_embed_interval_apply.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = bi.NewApplyEmbedIntervalRequest()
		response = bi.NewApplyEmbedIntervalResponse()
		biToken  string
	)
	if v, ok := d.GetOkExists("project_id"); ok {
		request.ProjectId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("page_id"); ok {
		request.PageId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("bi_token"); ok {
		biToken = v.(string)
		request.BIToken = helper.String(v.(string))
	}

	if v, ok := d.GetOk("scope"); ok {
		request.Scope = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseBiClient().ApplyEmbedInterval(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate bi embedInterval failed, reason:%+v", logId, err)
		return err
	}

	if !*response.Response.Data.Result {
		return fmt.Errorf("There was an error in token application, err: %v", response.Response.Msg)
	}

	d.SetId(biToken)

	return resourceTencentCloudBiEmbedIntervalApplyRead(d, meta)
}

func resourceTencentCloudBiEmbedIntervalApplyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_embed_interval_apply.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudBiEmbedIntervalApplyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_embed_interval_apply.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
