/*
Provides a resource to create a elasticsearch index

Example Usage

```hcl
resource "tencentcloud_elasticsearch_index" "index" {
  instance_id = "es-xxxxxx"
  index_type = "normal"
  index_name = "test-es-index"
  index_meta_json = "{\"mappings\":{},\"settings\":{\"index.number_of_replicas\":1,\"index.number_of_shards\":1,\"index.refresh_interval\":\"30s\"}}"
}
```

Import

elasticsearch index can be imported using the id, e.g.

```
terraform import tencentcloud_elasticsearch_index.index index_id
```
*/
package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	elasticsearch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/wI2L/jsondiff"
)

func resourceTencentCloudElasticsearchIndex() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudElasticsearchIndexCreate,
		Read:   resourceTencentCloudElasticsearchIndexRead,
		Update: resourceTencentCloudElasticsearchIndexUpdate,
		Delete: resourceTencentCloudElasticsearchIndexDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "es instance id.",
			},

			"index_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "type of the index to be created. auto: autonomous index. normal: indicates a common index.",
			},

			"index_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "index name to create.",
			},

			"index_meta_json": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Create index metadata JSON, such as mappings, settings.",
			},
		},
	}
}

func resourceTencentCloudElasticsearchIndexCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_index.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = elasticsearch.NewCreateIndexRequest()
		instanceId string
		indexName  string
		indexType  string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOk("index_type"); ok {
		indexType = v.(string)
		request.IndexType = helper.String(indexType)
	}

	if v, ok := d.GetOk("index_name"); ok {
		indexName = v.(string)
		request.IndexName = helper.String(indexName)
	}

	if v, ok := d.GetOk("index_meta_json"); ok {
		request.IndexMetaJson = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseEsClient().CreateIndex(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create elasticsearch index failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + FILED_SP + indexType + FILED_SP + indexName)

	return resourceTencentCloudElasticsearchIndexRead(d, meta)
}

func resourceTencentCloudElasticsearchIndexRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_index.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ElasticsearchService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	index, err := service.DescribeElasticsearchIndexByName(ctx, idSplit[0], idSplit[1], idSplit[2])
	if err != nil {
		return err
	}

	if index == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ElasticsearchIndex` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if index.ClusterId != nil {
		_ = d.Set("instance_id", index.ClusterId)
	}

	if index.IndexType != nil {
		_ = d.Set("index_type", index.IndexType)
	}

	if index.IndexName != nil {
		_ = d.Set("index_name", index.IndexName)
	}

	if index.IndexMetaJson != nil {
		// settings string value to int
		indexMetaJsonMap := make(map[string]interface{})
		err := json.Unmarshal([]byte(*index.IndexMetaJson), &indexMetaJsonMap)
		if err != nil {
			return err
		}
		if v, ok := indexMetaJsonMap["settings"]; ok {
			settingsMap := v.(map[string]interface{})

			if v, ok := settingsMap["index.number_of_replicas"]; ok {
				if _, ok := v.(string); ok {
					intValue, err := strconv.Atoi(v.(string))
					if err != nil {
						return err
					}
					settingsMap["index.number_of_replicas"] = intValue
				}

			}
			if v, ok := settingsMap["index.number_of_shards"]; ok {
				if _, ok := v.(string); ok {
					intValue, err := strconv.Atoi(v.(string))
					if err != nil {
						return err
					}
					settingsMap["index.number_of_shards"] = intValue
				}
			}
			indexMetaJsonMap["settings"] = settingsMap
		}
		newIndexMetaJson, err := json.Marshal(indexMetaJsonMap)
		if err != nil {
			return err
		}
		_ = d.Set("index_meta_json", string(newIndexMetaJson))
	}

	return nil
}

func resourceTencentCloudElasticsearchIndexUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_index.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := elasticsearch.NewUpdateIndexRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	request.InstanceId = &idSplit[0]

	immutableArgs := []string{"instance_id", "index_name"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.IndexType = helper.String(idSplit[1])
	request.IndexName = helper.String(idSplit[2])

	if d.HasChange("index_meta_json") {
		o, n := d.GetChange("index_meta_json")
		oldJson := o.(string)
		newJson := n.(string)
		patch, err := jsondiff.CompareJSON(
			[]byte(oldJson),
			[]byte(newJson),
		)
		log.Printf("patch: %v, %v\n%v", oldJson, newJson, patch)

		if err != nil {
			return err
		}
		result := make(map[string]interface{})
		for _, op := range patch {
			operationType := op.Type
			value := op.Value
			path := string(op.Path)
			if operationType == "remove" {
				err := setMapLinkKey(result, path, value)
				if err != nil {
					return err
				}
			}
			if operationType == "add" || operationType == "replace" {
				err := setMapLinkKey(result, path, value)
				if err != nil {
					return err
				}
			}
		}
		resultJson, err := json.Marshal(result)
		if err != nil {
			return err
		}
		updateBody := string(resultJson)
		request.UpdateMetaJson = helper.String(updateBody)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseEsClient().UpdateIndex(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update elasticsearch index failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudElasticsearchIndexRead(d, meta)
}

func resourceTencentCloudElasticsearchIndexDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_index.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ElasticsearchService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	if err := service.DeleteElasticsearchIndexByName(ctx, idSplit[0], idSplit[1], idSplit[2]); err != nil {
		return err
	}

	return nil
}

func setMapLinkKey(m map[string]interface{}, key string, value interface{}) error {
	targetMap := m
	keyLinks := strings.Split(key, "/")[1:]
	for _, subKey := range keyLinks[:len(keyLinks)-1] {
		if v, ok := targetMap[subKey]; ok {
			targetMap = v.(map[string]interface{})
		} else {
			subValue := make(map[string]interface{})
			targetMap[subKey] = subValue
			targetMap = subValue
		}
	}
	targetMap[keyLinks[len(keyLinks)-1]] = value
	return nil
}
