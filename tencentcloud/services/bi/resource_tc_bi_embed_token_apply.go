package bi

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bi "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bi/v20220105"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBiEmbedTokenApply() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBiEmbedTokenApplyCreate,
		Read:   resourceTencentCloudBiEmbedTokenApplyRead,
		Delete: resourceTencentCloudBiEmbedTokenApplyDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Share project id.",
			},

			"page_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Sharing page id, this is empty value 0 when embedding the board.",
			},

			"scope": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Page means embedding the page, and panel means embedding the entire board.",
			},

			"expire_time": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Expiration. Unit: Minutes Maximum value: 240. i.e. 4 hours Default: 240.",
			},

			"user_corp_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "User enterprise ID (for multi-user only).",
			},

			"user_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "UserId (for multi-user only).",
			},

			"ticket_num": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Access limit, the limit range is 1-99999, if it is empty, no access limit will be set.",
			},

			"bi_token": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Create the generated token.",
			},

			"create_at": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Create time.",
			},

			"udpate_at": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Upadte time.",
			},
		},
	}
}

func resourceTencentCloudBiEmbedTokenApplyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bi_embed_token_apply.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

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

	if v, ok := d.GetOk("user_corp_id"); ok {
		request.UserCorpId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_id"); ok {
		request.UserId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("ticket_num"); ok {
		request.TicketNum = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBiClient().CreateEmbedToken(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate bi embedToken failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(helper.Int64ToStr(int64(pageId)))

	if response.Response.Data != nil {
		token := response.Response.Data
		if token.BIToken != nil {
			_ = d.Set("bi_token", token.BIToken)
		}

		if token.CreatedAt != nil {
			_ = d.Set("create_at", token.CreatedAt)
		}

		if token.UpdatedAt != nil {
			_ = d.Set("udpate_at", token.UpdatedAt)
		}

		if token.TicketNum != nil {
			_ = d.Set("ticket_num", token.TicketNum)
		}
	}

	return resourceTencentCloudBiEmbedTokenApplyRead(d, meta)
}

func resourceTencentCloudBiEmbedTokenApplyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bi_embed_token_apply.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudBiEmbedTokenApplyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bi_embed_token_apply.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
