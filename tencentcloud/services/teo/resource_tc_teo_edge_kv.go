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

func ResourceTencentCloudTeoEdgeKV() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoEdgeKVCreate,
		Read:   resourceTencentCloudTeoEdgeKVRead,
		Update: resourceTencentCloudTeoEdgeKVUpdate,
		Delete: resourceTencentCloudTeoEdgeKVDelete,
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
				Description: "Namespace name.",
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Key name, 1-512 characters, allowed characters are letters, numbers, hyphens and underscores.",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key value. Cannot be empty, maximum 1 MB.",
			},
		},
	}
}

func resourceTencentCloudTeoEdgeKVCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_edge_kv.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewEdgeKVPutRequest()
		zoneId  string
		ns      string
		key     string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("namespace"); ok {
		request.Namespace = helper.String(v.(string))
		ns = v.(string)
	}

	if v, ok := d.GetOk("key"); ok {
		request.Key = helper.String(v.(string))
		key = v.(string)
	}

	if v, ok := d.GetOk("value"); ok {
		request.Value = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().EdgeKVPutWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create tencentcloud_teo_edge_kv failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create tencentcloud_teo_edge_kv failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(strings.Join([]string{zoneId, ns, key}, tccommon.FILED_SP))
	return resourceTencentCloudTeoEdgeKVRead(d, meta)
}

func resourceTencentCloudTeoEdgeKVRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_edge_kv.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewEdgeKVGetRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	ns := idSplit[1]
	key := idSplit[2]

	request.ZoneId = helper.String(zoneId)
	request.Namespace = helper.String(ns)
	request.Keys = []*string{helper.String(key)}

	var response *teov20220901.EdgeKVGetResponse
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().EdgeKVGetWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read tencentcloud_teo_edge_kv failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response == nil || response.Response == nil || len(response.Response.Data) == 0 {
		log.Printf("[WARN]%s resource `tencentcloud_teo_edge_kv` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	kvPair := response.Response.Data[0]
	if kvPair.Value == nil || *kvPair.Value == "" {
		log.Printf("[WARN]%s resource `tencentcloud_teo_edge_kv` [%s] value is empty, marking as deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("namespace", ns)
	_ = d.Set("key", key)
	_ = d.Set("value", kvPair.Value)

	return nil
}

func resourceTencentCloudTeoEdgeKVUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_edge_kv.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewEdgeKVPutRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	request.ZoneId = helper.String(idSplit[0])
	request.Namespace = helper.String(idSplit[1])
	request.Key = helper.String(idSplit[2])

	if v, ok := d.GetOk("value"); ok {
		request.Value = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().EdgeKVPutWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Update tencentcloud_teo_edge_kv failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update tencentcloud_teo_edge_kv failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudTeoEdgeKVRead(d, meta)
}

func resourceTencentCloudTeoEdgeKVDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_edge_kv.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewEdgeKVDeleteRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	ns := idSplit[1]
	key := idSplit[2]

	request.ZoneId = helper.String(zoneId)
	request.Namespace = helper.String(ns)
	request.Keys = []*string{helper.String(key)}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().EdgeKVDeleteWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete tencentcloud_teo_edge_kv failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete tencentcloud_teo_edge_kv failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
