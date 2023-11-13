/*
Provides a resource to create a tcm tracing_config

Example Usage

```hcl
resource "tencentcloud_tcm_tracing_config" "tracing_config" {
  mesh_id = "mesh-xxxxxxxx"
  enable = true
  a_p_m {
		enable = true
		region = "ap-shanghai"
		instance_id = "apm-xxx"

  }
  sampling =
  zipkin {
		address = "10.10.10.10:9411"

  }
}
```

Import

tcm tracing_config can be imported using the id, e.g.

```
terraform import tencentcloud_tcm_tracing_config.tracing_config tracing_config_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcm/v20210413"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTcmTracingConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcmTracingConfigCreate,
		Read:   resourceTencentCloudTcmTracingConfigRead,
		Update: resourceTencentCloudTcmTracingConfigUpdate,
		Delete: resourceTencentCloudTcmTracingConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"mesh_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Mesh ID.",
			},

			"enable": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether enable tracing.",
			},

			"a_p_m": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "APM config.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether enable APM.",
						},
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Region.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Instance id of the APM.",
						},
					},
				},
			},

			"sampling": {
				Optional:    true,
				Type:        schema.TypeFloat,
				Description: "Tracing sampling, 0.0-1.0.",
			},

			"zipkin": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Third party zipkin config.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Zipkin address.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTcmTracingConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_tracing_config.create")()
	defer inconsistentCheck(d, meta)()

	var meshId string
	if v, ok := d.GetOk("mesh_id"); ok {
		meshId = v.(string)
	}

	d.SetId(meshId)

	return resourceTencentCloudTcmTracingConfigUpdate(d, meta)
}

func resourceTencentCloudTcmTracingConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_tracing_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcmService{client: meta.(*TencentCloudClient).apiV3Conn}

	tracingConfigId := d.Id()

	TracingConfig, err := service.DescribeTcmTracingConfigById(ctx, meshId)
	if err != nil {
		return err
	}

	if TracingConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TcmTracingConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if TracingConfig.MeshId != nil {
		_ = d.Set("mesh_id", TracingConfig.MeshId)
	}

	if TracingConfig.Enable != nil {
		_ = d.Set("enable", TracingConfig.Enable)
	}

	if TracingConfig.APM != nil {
		aPMMap := map[string]interface{}{}

		if TracingConfig.APM.Enable != nil {
			aPMMap["enable"] = TracingConfig.APM.Enable
		}

		if TracingConfig.APM.Region != nil {
			aPMMap["region"] = TracingConfig.APM.Region
		}

		if TracingConfig.APM.InstanceId != nil {
			aPMMap["instance_id"] = TracingConfig.APM.InstanceId
		}

		_ = d.Set("a_p_m", []interface{}{aPMMap})
	}

	if TracingConfig.Sampling != nil {
		_ = d.Set("sampling", TracingConfig.Sampling)
	}

	if TracingConfig.Zipkin != nil {
		zipkinMap := map[string]interface{}{}

		if TracingConfig.Zipkin.Address != nil {
			zipkinMap["address"] = TracingConfig.Zipkin.Address
		}

		_ = d.Set("zipkin", []interface{}{zipkinMap})
	}

	return nil
}

func resourceTencentCloudTcmTracingConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_tracing_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tcm.NewModifyTracingConfigRequest()

	tracingConfigId := d.Id()

	request.MeshId = &meshId

	immutableArgs := []string{"mesh_id", "enable", "a_p_m", "sampling", "zipkin"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("enable") {
		if v, ok := d.GetOkExists("enable"); ok {
			request.Enable = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("a_p_m") {
		if dMap, ok := helper.InterfacesHeadMap(d, "a_p_m"); ok {
			aPM := tcm.APM{}
			if v, ok := dMap["enable"]; ok {
				aPM.Enable = helper.Bool(v.(bool))
			}
			if v, ok := dMap["region"]; ok {
				aPM.Region = helper.String(v.(string))
			}
			if v, ok := dMap["instance_id"]; ok {
				aPM.InstanceId = helper.String(v.(string))
			}
			request.APM = &aPM
		}
	}

	if d.HasChange("sampling") {
		if v, ok := d.GetOkExists("sampling"); ok {
			request.Sampling = helper.Float64(v.(float64))
		}
	}

	if d.HasChange("zipkin") {
		if dMap, ok := helper.InterfacesHeadMap(d, "zipkin"); ok {
			tracingZipkin := tcm.TracingZipkin{}
			if v, ok := dMap["address"]; ok {
				tracingZipkin.Address = helper.String(v.(string))
			}
			request.Zipkin = &tracingZipkin
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTcmClient().ModifyTracingConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tcm TracingConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTcmTracingConfigRead(d, meta)
}

func resourceTencentCloudTcmTracingConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_tracing_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
