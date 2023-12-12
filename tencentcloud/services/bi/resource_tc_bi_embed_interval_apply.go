package bi

import (
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bi "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bi/v20220105"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBiEmbedIntervalApply() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_bi_embed_interval_apply.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBiClient().ApplyEmbedInterval(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_bi_embed_interval_apply.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudBiEmbedIntervalApplyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bi_embed_interval_apply.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
