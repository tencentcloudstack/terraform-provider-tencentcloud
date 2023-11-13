/*
Provides a resource to create a tcm access_log_config

Example Usage

```hcl
resource "tencentcloud_tcm_access_log_config" "access_log_config" {
  mesh_name = "mesh-xxxxxxxx"
  selected_range {
		items {
			namespace = "prod"
			gateways =
		}
		all = false

  }
  template = "istio"
  enable = true
  c_l_s {
		enable = true
		log_set = "mesh-xxx"
		topic = "accesslog"

  }
  encoding = "TEXT"
  format = "[%START_TIME%]"
  enable_stdout = false
  enable_server = false
  address = "xxx"
}
```

Import

tcm access_log_config can be imported using the id, e.g.

```
terraform import tencentcloud_tcm_access_log_config.access_log_config access_log_config_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcm/v20210413"
	"log"
)

func resourceTencentCloudTcmAccessLogConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcmAccessLogConfigCreate,
		Read:   resourceTencentCloudTcmAccessLogConfigRead,
		Update: resourceTencentCloudTcmAccessLogConfigUpdate,
		Delete: resourceTencentCloudTcmAccessLogConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"mesh_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Mesh ID.",
			},

			"selected_range": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Selected range.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"items": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Items.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"namespace": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Namespace.",
									},
									"gateways": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Ingress gateway list.",
									},
								},
							},
						},
						"all": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Select all if true, default false.",
						},
					},
				},
			},

			"template": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Log template, istio/trace/custom.",
			},

			"enable": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether enable log.",
			},

			"c_l_s": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "CLS config.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether enable CLS.",
						},
						"log_set": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Log set of CLS.",
						},
						"topic": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Log topic of CLS.",
						},
					},
				},
			},

			"encoding": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Log encoding, TEXT or JSON.",
			},

			"format": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Log format.",
			},

			"enable_stdout": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether enable stdout.",
			},

			"enable_server": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether enable third party grpc server.",
			},

			"address": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Third party grpc server address.",
			},
		},
	}
}

func resourceTencentCloudTcmAccessLogConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_access_log_config.create")()
	defer inconsistentCheck(d, meta)()

	var meshName string
	if v, ok := d.GetOk("mesh_name"); ok {
		meshName = v.(string)
	}

	d.SetId(meshName)

	return resourceTencentCloudTcmAccessLogConfigUpdate(d, meta)
}

func resourceTencentCloudTcmAccessLogConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_access_log_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcmService{client: meta.(*TencentCloudClient).apiV3Conn}

	accessLogConfigId := d.Id()

	AccessLogConfig, err := service.DescribeTcmAccessLogConfigById(ctx, meshName)
	if err != nil {
		return err
	}

	if AccessLogConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TcmAccessLogConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if AccessLogConfig.MeshName != nil {
		_ = d.Set("mesh_name", AccessLogConfig.MeshName)
	}

	if AccessLogConfig.SelectedRange != nil {
		selectedRangeMap := map[string]interface{}{}

		if AccessLogConfig.SelectedRange.Items != nil {
			itemsList := []interface{}{}
			for _, items := range AccessLogConfig.SelectedRange.Items {
				itemsMap := map[string]interface{}{}

				if items.Namespace != nil {
					itemsMap["namespace"] = items.Namespace
				}

				if items.Gateways != nil {
					itemsMap["gateways"] = items.Gateways
				}

				itemsList = append(itemsList, itemsMap)
			}

			selectedRangeMap["items"] = []interface{}{itemsList}
		}

		if AccessLogConfig.SelectedRange.All != nil {
			selectedRangeMap["all"] = AccessLogConfig.SelectedRange.All
		}

		_ = d.Set("selected_range", []interface{}{selectedRangeMap})
	}

	if AccessLogConfig.Template != nil {
		_ = d.Set("template", AccessLogConfig.Template)
	}

	if AccessLogConfig.Enable != nil {
		_ = d.Set("enable", AccessLogConfig.Enable)
	}

	if AccessLogConfig.CLS != nil {
		cLSMap := map[string]interface{}{}

		if AccessLogConfig.CLS.Enable != nil {
			cLSMap["enable"] = AccessLogConfig.CLS.Enable
		}

		if AccessLogConfig.CLS.LogSet != nil {
			cLSMap["log_set"] = AccessLogConfig.CLS.LogSet
		}

		if AccessLogConfig.CLS.Topic != nil {
			cLSMap["topic"] = AccessLogConfig.CLS.Topic
		}

		_ = d.Set("c_l_s", []interface{}{cLSMap})
	}

	if AccessLogConfig.Encoding != nil {
		_ = d.Set("encoding", AccessLogConfig.Encoding)
	}

	if AccessLogConfig.Format != nil {
		_ = d.Set("format", AccessLogConfig.Format)
	}

	if AccessLogConfig.EnableStdout != nil {
		_ = d.Set("enable_stdout", AccessLogConfig.EnableStdout)
	}

	if AccessLogConfig.EnableServer != nil {
		_ = d.Set("enable_server", AccessLogConfig.EnableServer)
	}

	if AccessLogConfig.Address != nil {
		_ = d.Set("address", AccessLogConfig.Address)
	}

	return nil
}

func resourceTencentCloudTcmAccessLogConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_access_log_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tcm.NewModifyAccessLogConfigRequest()

	accessLogConfigId := d.Id()

	request.MeshName = &meshName

	immutableArgs := []string{"mesh_name", "selected_range", "template", "enable", "c_l_s", "encoding", "format", "enable_stdout", "enable_server", "address"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTcmClient().ModifyAccessLogConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tcm AccessLogConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTcmAccessLogConfigRead(d, meta)
}

func resourceTencentCloudTcmAccessLogConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_access_log_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
