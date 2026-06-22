package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoEdgeKVNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoEdgeKVNamespaceCreate,
		Read:   resourceTencentCloudTeoEdgeKVNamespaceRead,
		Update: resourceTencentCloudTeoEdgeKVNamespaceUpdate,
		Delete: resourceTencentCloudTeoEdgeKVNamespaceDelete,
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
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Namespace name. Supports 1-50 characters, allowed characters are a-z, A-Z, 0-9, -, and - cannot be used alone or consecutively, cannot be placed at the beginning or end. The name must be unique within the same site.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Namespace description. Used to describe the purpose or business meaning of the namespace. Maximum 256 characters.",
			},
			"capacity": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "KV storage available capacity in bytes.",
			},
			"capacity_used": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "KV storage used capacity in bytes.",
			},
			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time in ISO 8601 format (UTC).",
			},
			"modified_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modification time in ISO 8601 format (UTC).",
			},
		},
	}
}

func resourceTencentCloudTeoEdgeKVNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_edge_k_v_namespace.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewCreateEdgeKVNamespaceRequest()
		zoneId  string
		ns      string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("namespace"); ok {
		request.Namespace = helper.String(v.(string))
		ns = v.(string)
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateEdgeKVNamespaceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo_edge_k_v_namespace failed, Response is nil"))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo_edge_k_v_namespace failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create teo_edge_k_v_namespace, logId: %s, id: %s", logId, logId, d.Id())
	d.SetId(strings.Join([]string{zoneId, ns}, tccommon.FILED_SP))

	return resourceTencentCloudTeoEdgeKVNamespaceRead(d, meta)
}

func resourceTencentCloudTeoEdgeKVNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_edge_k_v_namespace.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewDescribeEdgeKVNamespacesRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	ns := idSplit[1]

	request.ZoneId = helper.String(zoneId)
	request.Limit = helper.Int64(1000)
	fuzzy := false
	request.Filters = []*teov20220901.AdvancedFilter{
		{
			Name:   helper.String("namespace"),
			Values: []*string{helper.String(ns)},
			Fuzzy:  &fuzzy,
		},
	}

	var response *teov20220901.DescribeEdgeKVNamespacesResponse
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeEdgeKVNamespacesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read teo_edge_k_v_namespace failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response == nil || response.Response == nil || len(response.Response.KVNamespaces) == 0 {
		log.Printf("[CRUD] teo_edge_k_v_namespace id=%s", d.Id())
		d.SetId("")
		return nil
	}

	// Find the exact matching namespace from the response list
	var target *teov20220901.KVNamespace
	for _, item := range response.Response.KVNamespaces {
		if item.Namespace != nil && *item.Namespace == ns {
			target = item
			break
		}
	}

	if target == nil {
		log.Printf("[CRUD] teo_edge_k_v_namespace id=%s", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	if target.Namespace != nil {
		_ = d.Set("namespace", target.Namespace)
	}

	if target.Remark != nil {
		_ = d.Set("remark", target.Remark)
	}

	if target.Capacity != nil {
		_ = d.Set("capacity", target.Capacity)
	}

	if target.CapacityUsed != nil {
		_ = d.Set("capacity_used", target.CapacityUsed)
	}

	if target.CreatedOn != nil {
		_ = d.Set("created_on", target.CreatedOn)
	}

	if target.ModifiedOn != nil {
		_ = d.Set("modified_on", target.ModifiedOn)
	}

	return nil
}

func resourceTencentCloudTeoEdgeKVNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_edge_k_v_namespace.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewModifyEdgeKVNamespaceRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	ns := idSplit[1]

	request.ZoneId = helper.String(zoneId)
	request.Namespace = helper.String(ns)

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		} else {
			request.Remark = helper.String("")
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyEdgeKVNamespaceWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Update teo_edge_k_v_namespace failed, Response is nil"))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update teo_edge_k_v_namespace failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudTeoEdgeKVNamespaceRead(d, meta)
}

func resourceTencentCloudTeoEdgeKVNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_edge_k_v_namespace.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewDeleteEdgeKVNamespaceRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	ns := idSplit[1]

	request.ZoneId = helper.String(zoneId)
	request.Namespace = helper.String(ns)

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteEdgeKVNamespaceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete teo_edge_k_v_namespace failed, Response is nil"))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo_edge_k_v_namespace failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
