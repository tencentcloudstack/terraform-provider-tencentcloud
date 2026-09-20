package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoInferenceAPIToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoInferenceAPITokenCreate,
		Read:   resourceTencentCloudTeoInferenceAPITokenRead,
		Update: resourceTencentCloudTeoInferenceAPITokenUpdate,
		Delete: resourceTencentCloudTeoInferenceAPITokenDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the inference API token. Max length: 30 characters.",
			},
			"token_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the inference API token.",
			},
			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Content of the inference API token.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the inference API token in ISO 8601 format.",
			},
		},
	}
}

func resourceTencentCloudTeoInferenceAPITokenCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = teo.NewCreateInferenceAPITokenRequest()
		response = teo.NewCreateInferenceAPITokenResponse()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		name := v.(string)
		if len(name) > 30 {
			return fmt.Errorf("name exceeds maximum length of 30 characters")
		}
		request.Name = helper.String(name)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateInferenceAPITokenWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TokenId == nil || *result.Response.TokenId == "" {
			return resource.NonRetryableError(fmt.Errorf("Create teo inference_api_token failed, Response is nil or TokenId is empty."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s create teo inference_api_token failed, reason:%+v", logId, err)
		return err
	}

	tokenId := *response.Response.TokenId

	log.Printf("[DEBUG]%s create teo inference_api_token success, logId=%s, tokenId=%s", logId, logId, tokenId)

	zoneId := d.Get("zone_id").(string)
	d.SetId(strings.Join([]string{zoneId, tokenId}, tccommon.FILED_SP))

	return resourceTencentCloudTeoInferenceAPITokenRead(d, meta)
}

func resourceTencentCloudTeoInferenceAPITokenRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	tokenId := idSplit[1]

	request := teo.NewDescribeInferenceAPITokensRequest()
	request.ZoneId = helper.String(zoneId)
	request.Limit = helper.Int64(100)

	var token *teo.InferenceAPIToken
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeInferenceAPITokensWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil || len(result.Response.Tokens) == 0 {
			return resource.NonRetryableError(fmt.Errorf("teo inference_api_token not found"))
		}

		for _, item := range result.Response.Tokens {
			if item.TokenId != nil && *item.TokenId == tokenId {
				token = item
				break
			}
		}

		if token == nil {
			return resource.NonRetryableError(fmt.Errorf("teo inference_api_token not found"))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRUD] teo inference_api_token id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if token == nil {
		log.Printf("[CRUD] teo inference_api_token id=%s", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)
	if token.TokenId != nil {
		_ = d.Set("token_id", token.TokenId)
	}
	if token.Name != nil {
		_ = d.Set("name", token.Name)
	}
	if token.Content != nil {
		_ = d.Set("content", token.Content)
	}
	if token.CreateTime != nil {
		_ = d.Set("create_time", token.CreateTime)
	}

	return nil
}

func resourceTencentCloudTeoInferenceAPITokenUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{"name"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	return resourceTencentCloudTeoInferenceAPITokenRead(d, meta)
}

func resourceTencentCloudTeoInferenceAPITokenDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teo.NewDeleteInferenceAPITokenRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	tokenId := idSplit[1]

	request.ZoneId = helper.String(zoneId)
	request.TokenId = helper.String(tokenId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteInferenceAPITokenWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete teo inference_api_token failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s delete teo inference_api_token failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
