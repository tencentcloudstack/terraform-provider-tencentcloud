package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTsfPathRewrite() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfPathRewriteCreate,
		Read:   resourceTencentCloudTsfPathRewriteRead,
		Update: resourceTencentCloudTsfPathRewriteUpdate,
		Delete: resourceTencentCloudTsfPathRewriteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"path_rewrite_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "path rewrite rule ID.",
			},

			"gateway_group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "gateway deployment group ID.",
			},

			"regex": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "regular expression.",
			},

			"replacement": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "content to replace.",
			},

			"blocked": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Whether to shield the mapped path, Y: Yes N: No.",
			},

			"order": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "rule order, the smaller the higher the priority.",
			},
		},
	}
}

func resourceTencentCloudTsfPathRewriteCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_path_rewrite.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = tsf.NewCreatePathRewritesWithDetailRespRequest()
		pathRewrites  = []*tsf.PathRewriteCreateObject{}
		pathRewrite   = tsf.PathRewriteCreateObject{}
		response      = tsf.NewCreatePathRewritesWithDetailRespResponse()
		pathRewriteId string
	)
	if v, ok := d.GetOk("gateway_group_id"); ok {
		pathRewrite.GatewayGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("regex"); ok {
		pathRewrite.Regex = helper.String(v.(string))
	}

	if v, ok := d.GetOk("replacement"); ok {
		pathRewrite.Replacement = helper.String(v.(string))
	}

	if v, ok := d.GetOk("blocked"); ok {
		pathRewrite.Blocked = helper.String(v.(string))
	}

	if v, _ := d.GetOk("order"); v != nil {
		pathRewrite.Order = helper.IntInt64(v.(int))
	}

	pathRewrites = append(pathRewrites, &pathRewrite)
	request.PathRewrites = pathRewrites
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreatePathRewritesWithDetailResp(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf pathRewrite failed, reason:%+v", logId, err)
		return err
	}

	pathRewriteId = *response.Response.Result[0]
	d.SetId(pathRewriteId)

	return resourceTencentCloudTsfPathRewriteRead(d, meta)
}

func resourceTencentCloudTsfPathRewriteRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_path_rewrite.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	pathRewriteId := d.Id()

	pathRewrite, err := service.DescribeTsfPathRewriteById(ctx, pathRewriteId)
	if err != nil {
		return err
	}

	if pathRewrite == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfPathRewrite` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if pathRewrite.PathRewriteId != nil {
		_ = d.Set("path_rewrite_id", pathRewrite.PathRewriteId)
	}

	if pathRewrite.GatewayGroupId != nil {
		_ = d.Set("gateway_group_id", pathRewrite.GatewayGroupId)
	}

	if pathRewrite.Regex != nil {
		_ = d.Set("regex", pathRewrite.Regex)
	}

	if pathRewrite.Replacement != nil {
		_ = d.Set("replacement", pathRewrite.Replacement)
	}

	if pathRewrite.Blocked != nil {
		_ = d.Set("blocked", pathRewrite.Blocked)
	}

	if pathRewrite.Order != nil {
		_ = d.Set("order", pathRewrite.Order)
	}

	return nil
}

func resourceTencentCloudTsfPathRewriteUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_path_rewrite.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tsf.NewModifyPathRewriteRequest()

	pathRewriteId := d.Id()

	request.PathRewriteId = &pathRewriteId

	immutableArgs := []string{"path_rewrite_id", "gateway_group_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("regex") {
		if v, ok := d.GetOk("regex"); ok {
			request.Regex = helper.String(v.(string))
		}
	}

	if d.HasChange("replacement") {
		if v, ok := d.GetOk("replacement"); ok {
			request.Replacement = helper.String(v.(string))
		}
	}

	if d.HasChange("blocked") {
		if v, ok := d.GetOk("blocked"); ok {
			request.Blocked = helper.String(v.(string))
		}
	}

	if d.HasChange("order") {
		if v, ok := d.GetOk("order"); ok {
			request.Order = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().ModifyPathRewrite(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf pathRewrite failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTsfPathRewriteRead(d, meta)
}

func resourceTencentCloudTsfPathRewriteDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_path_rewrite.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	pathRewriteId := d.Id()

	if err := service.DeleteTsfPathRewriteById(ctx, pathRewriteId); err != nil {
		return err
	}

	return nil
}
