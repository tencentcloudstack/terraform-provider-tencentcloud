package waf

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wafv20180125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWafObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafObjectCreate,
		Read:   resourceTencentCloudWafObjectRead,
		Update: resourceTencentCloudWafObjectUpdate,
		Delete: resourceTencentCloudWafObjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"object_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Modifies the object identifier.",
			},

			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "New instance ID: considered a successful modification if identical to an already bound instance.",
			},

			"status": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "New WAF switch status, considered successful if identical to existing status.",
			},

			"proxy": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether to enable proxy. 0: do not enable; 1: use the first IP address in XFF as the client IP address; 2: use remote_addr as the client IP address; 3: obtain the client IP address from the specified header field that is given in `ip_headers`.",
			},

			"ip_headers": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "This parameter indicates a custom header and is required when `proxy` is set to 3.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"member_app_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The ID of the member to whom the listener belongs.",
			},

			"member_uin": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Uin of the listener member.",
			},
		},
	}
}

func resourceTencentCloudWafObjectCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_object.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = wafv20180125.NewModifyObjectRequest()
		objectId string
	)

	if v, ok := d.GetOk("object_id"); ok {
		request.ObjectId = helper.String(v.(string))
		objectId = v.(string)
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("member_app_id"); ok {
		request.MemberAppId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("member_uin"); ok {
		request.MemberUin = helper.String(v.(string))
	}

	request.OpType = helper.String("InstanceId")
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyObjectWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create waf object failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(objectId)
	return resourceTencentCloudWafObjectUpdate(d, meta)
}

func resourceTencentCloudWafObjectRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_object.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service  = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		objectId = d.Id()
	)

	role, err := service.DescribeOrganizationRole(ctx)
	if err != nil {
		return err
	}

	if role == nil {
		role = helper.String("Member")
	}

	respData, err := service.DescribeWafObjectById(ctx, objectId, role)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_waf_object` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.ObjectId != nil {
		_ = d.Set("object_id", respData.ObjectId)
	}

	if respData.InstanceId != nil {
		_ = d.Set("instance_id", respData.InstanceId)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.Proxy != nil {
		_ = d.Set("proxy", respData.Proxy)
		if *respData.Proxy == 3 && len(respData.IpHeaders) > 0 {
			_ = d.Set("ip_headers", respData.IpHeaders)
		}
	}

	if respData.MemberAppId != nil {
		_ = d.Set("member_app_id", respData.MemberAppId)
	}

	if respData.MemberUin != nil {
		_ = d.Set("member_uin", respData.MemberUin)
	}

	return nil
}

func resourceTencentCloudWafObjectUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_object.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		objectId = d.Id()
	)

	if d.HasChange("status") {
		request := wafv20180125.NewModifyObjectRequest()
		if v, ok := d.GetOkExists("status"); ok {
			request.Status = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("member_app_id"); ok {
			request.MemberAppId = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("member_uin"); ok {
			request.MemberUin = helper.String(v.(string))
		}

		request.ObjectId = &objectId
		request.OpType = helper.String("Status")
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyObjectWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update waf object status failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if d.HasChange("instance_id") {
		request := wafv20180125.NewModifyObjectRequest()
		if v, ok := d.GetOk("instance_id"); ok {
			request.InstanceId = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("member_app_id"); ok {
			request.MemberAppId = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("member_uin"); ok {
			request.MemberUin = helper.String(v.(string))
		}

		request.ObjectId = &objectId
		request.OpType = helper.String("InstanceId")
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyObjectWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update waf object instance id failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if d.HasChange("proxy") || d.HasChange("ip_headers") {
		request := wafv20180125.NewModifyObjectRequest()
		if v, ok := d.GetOkExists("proxy"); ok {
			request.Proxy = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("ip_headers"); ok {
			ipHeadersSet := v.(*schema.Set).List()
			for i := range ipHeadersSet {
				ipHeaders := ipHeadersSet[i].(string)
				request.IpHeaders = append(request.IpHeaders, helper.String(ipHeaders))
			}
		}

		if v, ok := d.GetOkExists("member_app_id"); ok {
			request.MemberAppId = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("member_uin"); ok {
			request.MemberUin = helper.String(v.(string))
		}

		request.ObjectId = &objectId
		request.OpType = helper.String("Proxy")
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyObjectWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update waf object proxy failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudWafObjectRead(d, meta)
}

func resourceTencentCloudWafObjectDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_object.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		objectId = d.Id()
	)

	// stop first
	stopReq := wafv20180125.NewModifyObjectRequest()
	if v, ok := d.GetOkExists("member_app_id"); ok {
		stopReq.MemberAppId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("member_uin"); ok {
		stopReq.MemberUin = helper.String(v.(string))
	}

	stopReq.ObjectId = &objectId
	stopReq.Status = helper.IntInt64(0)
	stopReq.OpType = helper.String("Status")
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyObjectWithContext(ctx, stopReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, stopReq.GetAction(), stopReq.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update waf object status failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// delete second
	request := wafv20180125.NewModifyObjectRequest()
	if v, ok := d.GetOkExists("member_app_id"); ok {
		request.MemberAppId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("member_uin"); ok {
		request.MemberUin = helper.String(v.(string))
	}

	request.ObjectId = &objectId
	request.OpType = helper.String("UnbindInstance")
	reqErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyObjectWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete waf object failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
