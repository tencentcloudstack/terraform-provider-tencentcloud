/*
Provides a resource to query TEO edge KV data

# Example Usage

```hcl

resource "tencentcloud_teo_edge_k_v_get" "edge_kv_get" {
    zone_id   = "zone-3j1xw7910arp"
    namespace = "ns-011"
    keys      = ["hello", "world"]
}

```
Import

teo edge_k_v_get can be imported using the zoneId#namespace#keysHash, e.g.
```
terraform import tencentcloud_teo_edge_k_v_get.edge_kv_get zone-3j1xw7910arp#ns-011#abc123
```
*/
package teo

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoEdgeKVGet() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoEdgeKVGetCreate,
		Read:   resourceTencentCloudTeoEdgeKVGetRead,
		Update: resourceTencentCloudTeoEdgeKVGetUpdate,
		Delete: resourceTencentCloudTeoEdgeKVGetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Namespace name. You can get the namespace list under the site through the DescribeEdgeKVNamespaces interface.",
			},

			"keys": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Key name list. The maximum length of the array is 20. Each key name cannot be empty, the length is 1-512 characters, and the allowed characters are letters, numbers, hyphens, and underscores.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems: 20,
			},

			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Key-value pair data list. The results are returned in the order of the input Keys. If a key does not exist, the Value field of the corresponding item returns an empty string.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key name. Each key name cannot be empty, the length is 1-512 characters, and the allowed characters are letters, numbers, hyphens, and underscores.",
						},

						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key value. Cannot be empty when entering parameters, supports up to 1 MB. When outputting, if the key does not exist, an empty string is returned.",
						},

						"expiration": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expiration time, following ISO 8601 standard, in the format YYYY-MM-DDThh:mm:ssZ (UTC time). When outputting, if it is an empty string, it means the key-value pair will never expire.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoEdgeKVGetCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_edge_k_v_get.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		zoneId   string
		namespace string
		keys      []interface{}
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("namespace"); ok {
		namespace = v.(string)
	}

	if v, ok := d.GetOk("keys"); ok {
		keys = v.([]interface{})
	}

	// Generate resource ID
	keysHash := hashKeys(keys)
	resourceId := fmt.Sprintf("%s#%s#%s", zoneId, namespace, keysHash)
	d.SetId(resourceId)

	// Query edge KV data
	err := resourceTencentCloudTeoEdgeKVGetRead(d, meta)
	if err != nil {
		return fmt.Errorf("[CRITAL]%s create teo EdgeKVGet failed, reason:%+v", logId, err)
	}

	return nil
}

func resourceTencentCloudTeoEdgeKVGetRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_edge_k_v_get.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	// Parse resource ID
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	namespace := idSplit[1]

	var (
		keys []interface{}
	)

	if v, ok := d.GetOk("keys"); ok {
		keys = v.([]interface{})
	}

	// Convert keys to []*string
	keysStrings := make([]*string, 0, len(keys))
	for _, k := range keys {
		keyString := k.(string)
		keysStrings = append(keysStrings, &keyString)
	}

	// Call EdgeKVGet API
	request := teo.NewEdgeKVGetRequest()
	request.ZoneId = &zoneId
	request.Namespace = &namespace
	request.Keys = keysStrings

	response, err := service.EdgeKVGet(ctx, request)
	if err != nil {
		return err
	}

	// Set data to state
	if response.Data != nil {
		dataList := make([]map[string]interface{}, 0, len(response.Data))
		for _, item := range response.Data {
			data := make(map[string]interface{})
			if item.Key != nil {
				data["key"] = *item.Key
			}
			if item.Value != nil {
				data["value"] = *item.Value
			}
			if item.Expiration != nil {
				data["expiration"] = *item.Expiration
			}
			dataList = append(dataList, data)
		}
		if err := d.Set("data", dataList); err != nil {
			return fmt.Errorf("set data error: %s", err)
		}
	}

	if err := d.Set("zone_id", zoneId); err != nil {
		return fmt.Errorf("set zone_id error: %s", err)
	}

	if err := d.Set("namespace", namespace); err != nil {
		return fmt.Errorf("set namespace error: %s", err)
	}

	return nil
}

func resourceTencentCloudTeoEdgeKVGetUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_edge_k_v_get.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	// For update operation, regenerate ID if zone_id, namespace, or keys changed
	if d.HasChange("zone_id") || d.HasChange("namespace") || d.HasChange("keys") {
		var (
			zoneId   string
			namespace string
			keys      []interface{}
		)

		if v, ok := d.GetOk("zone_id"); ok {
			zoneId = v.(string)
		}

		if v, ok := d.GetOk("namespace"); ok {
			namespace = v.(string)
		}

		if v, ok := d.GetOk("keys"); ok {
			keys = v.([]interface{})
		}

		// Generate new resource ID
		keysHash := hashKeys(keys)
		newResourceId := fmt.Sprintf("%s#%s#%s", zoneId, namespace, keysHash)
		d.SetId(newResourceId)
	}

	// Query edge KV data with new parameters
	err := resourceTencentCloudTeoEdgeKVGetRead(d, meta)
	if err != nil {
		return fmt.Errorf("[CRITAL]%s update teo EdgeKVGet failed, reason:%+v", logId, err)
	}

	return nil
}

func resourceTencentCloudTeoEdgeKVGetDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_edge_k_v_get.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	// For query resource, delete operation only removes from state
	d.SetId("")

	log.Printf("[WARN]%s resource `TeoEdgeKVGet` [%s] has been deleted from state.\n", logId, d.Id())

	return nil
}

// hashKeys generates a hash from keys list for unique identification
func hashKeys(keys []interface{}) string {
	if len(keys) == 0 {
		return ""
	}

	// Sort keys to ensure consistent hash
	keysStrings := make([]string, 0, len(keys))
	for _, k := range keys {
		keysStrings = append(keysStrings, k.(string))
	}

	// Join keys with delimiter and generate hash
	joined := strings.Join(keysStrings, ",")
	hash := md5.Sum([]byte(joined))

	return hex.EncodeToString(hash[:])
}

// EdgeKVGet is a method to call TEO EdgeKVGet API
func (s *TeoService) EdgeKVGet(ctx context.Context, request *teo.EdgeKVGetRequest) (*teo.EdgeKVGetResponse, error) {
	response, err := s.client.UseTeoClient().EdgeKVGet(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
