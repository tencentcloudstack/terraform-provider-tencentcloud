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

func ResourceTencentCloudTeoEdgeKVPut() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoEdgeKVPutCreate,
		Read:   resourceTencentCloudTeoEdgeKVPutRead,
		Update: resourceTencentCloudTeoEdgeKVPutUpdate,
		Delete: resourceTencentCloudTeoEdgeKVPutDelete,
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
			"expiration": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Expiration time of the key-value pair, absolute time, Unix timestamp in seconds. Must be greater than or equal to current time + 60.",
			},
			"expiration_ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Time-to-live of the key-value pair, relative time in seconds. Must be greater than or equal to 60.",
			},
		},
	}
}

func resourceTencentCloudTeoEdgeKVPutCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_edge_k_v_put.create")()
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

	if v, ok := d.GetOkExists("expiration"); ok {
		request.Expiration = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOkExists("expiration_ttl"); ok {
		request.ExpirationTTL = helper.Int64(int64(v.(int)))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().EdgeKVPutWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo edge kv put failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo edge kv put failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(strings.Join([]string{zoneId, ns, key}, tccommon.FILED_SP))
	return resourceTencentCloudTeoEdgeKVPutRead(d, meta)
}

func resourceTencentCloudTeoEdgeKVPutRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_edge_k_v_put.read")()
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
		log.Printf("[CRITAL]%s read teo edge kv put failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response == nil || response.Response == nil || len(response.Response.Data) == 0 {
		log.Printf("[WARN]%s resource `tencentcloud_teo_edge_k_v_put` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	kvPair := response.Response.Data[0]
	if kvPair.Value == nil || *kvPair.Value == "" {
		log.Printf("[WARN]%s resource `tencentcloud_teo_edge_k_v_put` [%s] value is empty, marking as deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("namespace", ns)
	_ = d.Set("key", key)
	_ = d.Set("value", kvPair.Value)

	return nil
}

func resourceTencentCloudTeoEdgeKVPutUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_edge_k_v_put.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{"value", "expiration", "expiration_ttl"}
	if err := helper.ImmutableArgsChek(d, immutableArgs...); err != nil {
		return err
	}

	return resourceTencentCloudTeoEdgeKVPutRead(d, meta)
}

func resourceTencentCloudTeoEdgeKVPutDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_edge_k_v_put.delete")()
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
			return resource.NonRetryableError(fmt.Errorf("Delete teo edge kv put failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo edge kv put failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
