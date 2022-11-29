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
  cls {
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
$ terraform import tencentcloud_tcm_access_log_config.access_log_config accessLogConfig_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tcm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcm/v20210413"
)

func resourceTencentCloudTcmAccessLogConfig() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTcmAccessLogConfigRead,
		Create: resourceTencentCloudTcmAccessLogConfigCreate,
		Update: resourceTencentCloudTcmAccessLogConfigUpdate,
		Delete: resourceTencentCloudTcmAccessLogConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"mesh_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Mesh ID.",
			},

			"selected_range": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Log template, istio/trace/custome.",
			},

			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether enable log.",
			},

			"cls": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Log encoding, TEXT or JSON.",
			},

			"format": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Log format.",
			},

			"enable_stdout": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether enable stdout.",
			},

			"enable_server": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether enable third party grpc server.",
			},

			"address": {
				Type:        schema.TypeString,
				Optional:    true,
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

	meshName := d.Id()

	accessLogConfig, err := service.DescribeTcmAccessLogConfig(ctx, meshName)

	if err != nil {
		return err
	}

	if accessLogConfig == nil {
		d.SetId("")
		return fmt.Errorf("resource `accessLogConfig` %s does not exist", meshName)
	}

	_ = d.Set("mesh_name", meshName)

	if accessLogConfig.SelectedRange != nil {
		selectedRangeMap := map[string]interface{}{}
		if accessLogConfig.SelectedRange.Items != nil {
			itemsList := []interface{}{}
			for _, items := range accessLogConfig.SelectedRange.Items {
				itemsMap := map[string]interface{}{}
				if items.Namespace != nil {
					itemsMap["namespace"] = items.Namespace
				}
				if items.Gateways != nil {
					itemsMap["gateways"] = items.Gateways
				}

				itemsList = append(itemsList, itemsMap)
			}
			selectedRangeMap["items"] = itemsList
		}
		if accessLogConfig.SelectedRange.All != nil {
			selectedRangeMap["all"] = accessLogConfig.SelectedRange.All
		}

		_ = d.Set("selected_range", []interface{}{selectedRangeMap})
	}

	if accessLogConfig.Template != nil {
		_ = d.Set("template", accessLogConfig.Template)
	}

	if accessLogConfig.Enable != nil {
		_ = d.Set("enable", accessLogConfig.Enable)
	}

	if accessLogConfig.CLS != nil {
		cLSMap := map[string]interface{}{}
		if accessLogConfig.CLS.Enable != nil {
			cLSMap["enable"] = accessLogConfig.CLS.Enable
		}
		if accessLogConfig.CLS.LogSet != nil {
			cLSMap["log_set"] = accessLogConfig.CLS.LogSet
		}
		if accessLogConfig.CLS.Topic != nil {
			cLSMap["topic"] = accessLogConfig.CLS.Topic
		}

		_ = d.Set("cls", []interface{}{cLSMap})
	}

	if accessLogConfig.Encoding != nil {
		_ = d.Set("encoding", accessLogConfig.Encoding)
	}

	if accessLogConfig.Format != nil {
		_ = d.Set("format", accessLogConfig.Format)
	}

	if accessLogConfig.EnableStdout != nil {
		_ = d.Set("enable_stdout", accessLogConfig.EnableStdout)
	}

	if accessLogConfig.EnableServer != nil {
		_ = d.Set("enable_server", accessLogConfig.EnableServer)
	}

	if accessLogConfig.Address != nil {
		_ = d.Set("address", accessLogConfig.Address)
	}

	return nil
}

func resourceTencentCloudTcmAccessLogConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_access_log_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tcm.NewModifyAccessLogConfigRequest()

	meshName := d.Id()

	request.MeshId = &meshName

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTcmClient().ModifyAccessLogConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tcm accessLogConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTcmAccessLogConfigRead(d, meta)
}

func resourceTencentCloudTcmAccessLogConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_access_log_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
