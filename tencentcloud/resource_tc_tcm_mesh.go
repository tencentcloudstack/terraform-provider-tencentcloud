/*
Provides a resource to create a tcm mesh

Example Usage

```hcl
resource "tencentcloud_tcm_mesh" "basic" {
  display_name = "test_mesh"
  mesh_version = "1.12.5"
  type = "HOSTED"
  config {
    istio {
      outbound_traffic_policy = "ALLOW_ANY"
      disable_policy_checks = true
      enable_pilot_http = true
      disable_http_retry = true
      smart_dns {
        istio_meta_dns_capture = true
        istio_meta_dns_auto_allocate = true
      }
    }
  }
  tag_list {
    key = "key"
    value = "value"
    passthrough = true
  }
}

```
Import

tcm mesh can be imported using the id, e.g.
```
$ terraform import tencentcloud_tcm_mesh.mesh mesh_id
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
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTcmMesh() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTcmMeshRead,
		Create: resourceTencentCloudTcmMeshCreate,
		Update: resourceTencentCloudTcmMeshUpdate,
		Delete: resourceTencentCloudTcmMeshDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"mesh_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mesh ID.",
			},

			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Mesh name.",
			},

			"mesh_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Mesh version.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Mesh type.",
			},

			"config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Mesh configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"istio": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Istio configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"outbound_traffic_policy": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Outbound traffic policy.",
									},
									"disable_policy_checks": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Disable policy checks.",
									},
									"enable_pilot_http": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Enable HTTP/1.0 support.",
									},
									"disable_http_retry": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Disable http retry.",
									},
									"smart_dns": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "SmartDNS configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"istio_meta_dns_capture": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Enable dns proxy.",
												},
												"istio_meta_dns_auto_allocate": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Enable auto allocate address.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"tag_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of associated tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
						"passthrough": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Passthrough to other related product.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTcmMeshCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_mesh.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tcm.NewCreateMeshRequest()
		response *tcm.CreateMeshResponse
		meshId   string
	)

	if v, ok := d.GetOk("display_name"); ok {
		request.DisplayName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("mesh_version"); ok {
		request.MeshVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "config"); ok {
		meshConfig := tcm.MeshConfig{}
		if IstioMap, ok := helper.InterfaceToMap(dMap, "istio"); ok {
			istioConfig := tcm.IstioConfig{}
			if v, ok := IstioMap["outbound_traffic_policy"]; ok {
				istioConfig.OutboundTrafficPolicy = helper.String(v.(string))
			}
			if v, ok := IstioMap["disable_policy_checks"]; ok {
				istioConfig.DisablePolicyChecks = helper.Bool(v.(bool))
			}
			if v, ok := IstioMap["enable_pilot_http"]; ok {
				istioConfig.EnablePilotHTTP = helper.Bool(v.(bool))
			}
			if v, ok := IstioMap["disable_http_retry"]; ok {
				istioConfig.DisableHTTPRetry = helper.Bool(v.(bool))
			}
			if SmartDNSMap, ok := helper.InterfaceToMap(IstioMap, "smart_dns"); ok {
				smartDNSConfig := tcm.SmartDNSConfig{}
				if v, ok := SmartDNSMap["istio_meta_dns_capture"]; ok {
					smartDNSConfig.IstioMetaDNSCapture = helper.Bool(v.(bool))
				}
				if v, ok := SmartDNSMap["istio_meta_dns_auto_allocate"]; ok {
					smartDNSConfig.IstioMetaDNSAutoAllocate = helper.Bool(v.(bool))
				}
				istioConfig.SmartDNS = &smartDNSConfig
			}
			meshConfig.Istio = &istioConfig
		}

		request.Config = &meshConfig
	}

	if v, ok := d.GetOk("tag_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			tag := tcm.Tag{}
			if v, ok := dMap["key"]; ok {
				tag.Key = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				tag.Value = helper.String(v.(string))
			}
			if v, ok := dMap["passthrough"]; ok {
				tag.Passthrough = helper.Bool(v.(bool))
			}

			request.TagList = append(request.TagList, &tag)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTcmClient().CreateMesh(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tcm mesh failed, reason:%+v", logId, err)
		return err
	}

	meshId = *response.Response.MeshId

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TcmService{client: meta.(*TencentCloudClient).apiV3Conn}
	err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
		mesh, errRet := service.DescribeTcmMesh(ctx, meshId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if *mesh.Mesh.State == "PENDING" || *mesh.Mesh.State == "CREATING" {
			return resource.RetryableError(fmt.Errorf("mesh status is %v, retry...", *mesh.Mesh.State))
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(meshId)
	return resourceTencentCloudTcmMeshRead(d, meta)
}

func resourceTencentCloudTcmMeshRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_mesh.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcmService{client: meta.(*TencentCloudClient).apiV3Conn}

	meshId := d.Id()

	meshResponse, err := service.DescribeTcmMesh(ctx, meshId)

	if err != nil {
		return err
	}

	mesh := meshResponse.Mesh
	if mesh == nil {
		d.SetId("")
		return fmt.Errorf("resource `mesh` %s does not exist", meshId)
	}

	if mesh.MeshId != nil {
		_ = d.Set("mesh_id", mesh.MeshId)
	}

	if mesh.DisplayName != nil {
		_ = d.Set("display_name", mesh.DisplayName)
	}

	if mesh.Version != nil {
		_ = d.Set("mesh_version", mesh.Version)
	}

	if mesh.Type != nil {
		_ = d.Set("type", mesh.Type)
	}

	if mesh.Config != nil {
		configMap := map[string]interface{}{}
		if mesh.Config.Istio != nil {
			istioMap := map[string]interface{}{}
			if mesh.Config.Istio.OutboundTrafficPolicy != nil {
				istioMap["outbound_traffic_policy"] = mesh.Config.Istio.OutboundTrafficPolicy
			}
			if mesh.Config.Istio.DisablePolicyChecks != nil {
				istioMap["disable_policy_checks"] = mesh.Config.Istio.DisablePolicyChecks
			}
			if mesh.Config.Istio.EnablePilotHTTP != nil {
				istioMap["enable_pilot_http"] = mesh.Config.Istio.EnablePilotHTTP
			}
			if mesh.Config.Istio.DisableHTTPRetry != nil {
				istioMap["disable_http_retry"] = mesh.Config.Istio.DisableHTTPRetry
			}
			if mesh.Config.Istio.SmartDNS != nil {
				smartDNSMap := map[string]interface{}{}
				if mesh.Config.Istio.SmartDNS.IstioMetaDNSCapture != nil {
					smartDNSMap["istio_meta_dns_capture"] = mesh.Config.Istio.SmartDNS.IstioMetaDNSCapture
				}
				if mesh.Config.Istio.SmartDNS.IstioMetaDNSAutoAllocate != nil {
					smartDNSMap["istio_meta_dns_auto_allocate"] = mesh.Config.Istio.SmartDNS.IstioMetaDNSAutoAllocate
				}

				istioMap["smart_dns"] = []interface{}{smartDNSMap}
			}

			configMap["istio"] = []interface{}{istioMap}
		}

		_ = d.Set("config", []interface{}{configMap})
	}

	if mesh.TagList != nil {
		tagListList := []interface{}{}
		for _, tagList := range mesh.TagList {
			tagListMap := map[string]interface{}{}
			if tagList.Key != nil {
				tagListMap["key"] = tagList.Key
			}
			if tagList.Value != nil {
				tagListMap["value"] = tagList.Value
			}
			if tagList.Passthrough != nil {
				tagListMap["passthrough"] = tagList.Passthrough
			}

			tagListList = append(tagListList, tagListMap)
		}
		_ = d.Set("tag_list", tagListList)
	}

	return nil
}

func resourceTencentCloudTcmMeshUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_mesh.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	request := tcm.NewModifyMeshRequest()

	meshId := d.Id()
	request.MeshId = &meshId

	if d.HasChange("mesh_id") {
		return fmt.Errorf("`mesh_id` do not support change now.")
	}

	if d.HasChange("display_name") {
		if v, ok := d.GetOk("display_name"); ok {
			request.DisplayName = helper.String(v.(string))
		}
	}

	if d.HasChange("mesh_version") {
		return fmt.Errorf("`mesh_version` do not support change now.")
	}

	if d.HasChange("type") {
		return fmt.Errorf("`type` do not support change now.")
	}

	if d.HasChange("config") {
		if dMap, ok := helper.InterfacesHeadMap(d, "config"); ok {
			meshConfig := tcm.MeshConfig{}
			if IstioMap, ok := helper.InterfaceToMap(dMap, "istio"); ok {
				istioConfig := tcm.IstioConfig{}
				if v, ok := IstioMap["outbound_traffic_policy"]; ok {
					istioConfig.OutboundTrafficPolicy = helper.String(v.(string))
				}
				if v, ok := IstioMap["disable_policy_checks"]; ok {
					istioConfig.DisablePolicyChecks = helper.Bool(v.(bool))
				}
				if v, ok := IstioMap["enable_pilot_http"]; ok {
					istioConfig.EnablePilotHTTP = helper.Bool(v.(bool))
				}
				if v, ok := IstioMap["disable_http_retry"]; ok {
					istioConfig.DisableHTTPRetry = helper.Bool(v.(bool))
				}
				if SmartDNSMap, ok := helper.InterfaceToMap(IstioMap, "smart_dns"); ok {
					smartDNSConfig := tcm.SmartDNSConfig{}
					if v, ok := SmartDNSMap["istio_meta_dns_capture"]; ok {
						smartDNSConfig.IstioMetaDNSCapture = helper.Bool(v.(bool))
					}
					if v, ok := SmartDNSMap["istio_meta_dns_auto_allocate"]; ok {
						smartDNSConfig.IstioMetaDNSAutoAllocate = helper.Bool(v.(bool))
					}
					istioConfig.SmartDNS = &smartDNSConfig
				}
				meshConfig.Istio = &istioConfig
			}
			request.Config = &meshConfig
		}
	}

	if d.HasChange("tag_list") {
		return fmt.Errorf("`tag_list` do not support change now.")
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTcmClient().ModifyMesh(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tcm mesh failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTcmMeshRead(d, meta)
}

func resourceTencentCloudTcmMeshDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_mesh.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcmService{client: meta.(*TencentCloudClient).apiV3Conn}

	meshId := d.Id()

	if err := service.DeleteTcmMeshById(ctx, meshId); err != nil {
		return err
	}

	err := resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
		mesh, errRet := service.DescribeTcmMesh(ctx, meshId)
		if errRet != nil {
			if isExpectError(errRet, []string{"ResourceNotFound"}) {
				return nil
			}
			return retryError(errRet, InternalError)
		}
		if mesh != nil {
			if *mesh.Mesh.State == "DELETING" {
				return resource.RetryableError(fmt.Errorf("mesh status is %v, retry...", *mesh.Mesh.State))
			}
			if *mesh.Mesh.State == "DELETE_FAILED" {
				return resource.NonRetryableError(fmt.Errorf("mesh status is %v, retry...", *mesh.Mesh.State))
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
