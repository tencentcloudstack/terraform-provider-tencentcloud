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
				Description: "Namespace name. Supports 1-50 characters, allowed characters are a-z, A-Z, 0-9, -.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Namespace description. Maximum 256 characters.",
			},
		},
	}
}

func resourceTencentCloudTeoEdgeKVNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_edge_k_v_namespace.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = teov20220901.NewCreateEdgeKVNamespaceRequest()
		zoneId    string
		namespace string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("namespace"); ok {
		request.Namespace = helper.String(v.(string))
		namespace = v.(string)
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateEdgeKVNamespaceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo edge kv namespace failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo edge kv namespace failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{zoneId, namespace}, tccommon.FILED_SP))
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
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	namespace := idSplit[1]

	request.ZoneId = &zoneId
	request.Limit = helper.Int64(1000)
	request.Filters = []*teov20220901.AdvancedFilter{
		{
			Name:   helper.String("namespace"),
			Values: []*string{&namespace},
			Fuzzy:  helper.Bool(false),
		},
	}

	var kvNamespace *teov20220901.KVNamespace
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeEdgeKVNamespacesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe teo edge kv namespaces failed, Response is nil."))
		}

		for _, ns := range result.Response.KVNamespaces {
			if ns.Namespace != nil && *ns.Namespace == namespace {
				kvNamespace = ns
				break
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read teo edge kv namespace failed, reason:%+v", logId, err)
		return err
	}

	if kvNamespace == nil {
		log.Printf("[WARN]%s resource `tencentcloud_teo_edge_k_v_namespace` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	if kvNamespace.Namespace != nil {
		_ = d.Set("namespace", kvNamespace.Namespace)
	}

	if kvNamespace.Remark != nil {
		_ = d.Set("remark", kvNamespace.Remark)
	}

	return nil
}

func resourceTencentCloudTeoEdgeKVNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_edge_k_v_namespace.update")()
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
	namespace := idSplit[1]

	if d.HasChange("remark") {
		request := teov20220901.NewModifyEdgeKVNamespaceRequest()
		request.ZoneId = &zoneId
		request.Namespace = &namespace

		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		} else {
			request.Remark = helper.String("")
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyEdgeKVNamespaceWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify teo edge kv namespace failed, Response is nil."))
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update teo edge kv namespace failed, reason:%+v", logId, err)
			return err
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
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	namespace := idSplit[1]

	request.ZoneId = &zoneId
	request.Namespace = &namespace

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteEdgeKVNamespaceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete teo edge kv namespace failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete teo edge kv namespace failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
