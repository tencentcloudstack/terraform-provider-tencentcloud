/*
Provides a resource to create a bi embed_interval

Example Usage

```hcl
resource "tencentcloud_bi_embed_interval" "embed_interval" {
  project_id = 123
  page_id = 123
  b_i_token = "abc"
  extra_param = ""
  scope = "page"
}
```

Import

bi embed_interval can be imported using the id, e.g.

```
terraform import tencentcloud_bi_embed_interval.embed_interval embed_interval_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bi "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bi/v20220105"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudBiEmbedInterval() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBiEmbedIntervalCreate,
		Read:   resourceTencentCloudBiEmbedIntervalRead,
		Delete: resourceTencentCloudBiEmbedIntervalDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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

			"b_i_token": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Token that needs to be applied for extension.",
			},

			"extra_param": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Alternate fields.",
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

func resourceTencentCloudBiEmbedIntervalCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_embed_interval.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = bi.NewApplyEmbedIntervalRequest()
		response = bi.NewApplyEmbedIntervalResponse()
		bIToken  string
	)
	if v, _ := d.GetOk("project_id"); v != nil {
		request.ProjectId = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("page_id"); v != nil {
		request.PageId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("b_i_token"); ok {
		bIToken = v.(string)
		request.BIToken = helper.String(v.(string))
	}

	if v, ok := d.GetOk("extra_param"); ok {
		request.ExtraParam = helper.String(v.(string))
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

	bIToken = *response.Response.BIToken
	d.SetId(bIToken)

	return resourceTencentCloudBiEmbedIntervalRead(d, meta)
}

func resourceTencentCloudBiEmbedIntervalRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_embed_interval.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudBiEmbedIntervalDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bi_embed_interval.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
