/*
Use this data source to query detailed information of tcm mesh

Example Usage

```hcl
data "tencentcloud_tcm_mesh" "mesh" {
  mesh_id = ["mesh-xxxxxx"]
  mesh_name = ["KEEP_MASH"]
  tags = ["key"]
  mesh_cluster = ["cls-xxxx"]
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tcm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcm/v20210413"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTcmMesh() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcmMeshRead,
		Schema: map[string]*schema.Schema{
			"mesh_id": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Mesh instance Id.",
			},

			"mesh_name": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Display name.",
			},

			"tags": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "tag key.",
			},

			"mesh_cluster": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Mesh name.",
			},

			"mesh_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The mesh information is queriedNote: This field may return null, indicating that a valid value is not available.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mesh_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mesh instance Id.",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mesh name.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mesh version.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mesh type.  Value range:- `STANDALONE`: Standalone mesh- `HOSTED`: hosted the mesh.",
						},
						"config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Mesh configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"istio": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Istio configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"outbound_traffic_policy": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Outbound traffic policy.",
												},
												"disable_policy_checks": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Disable policy checks.",
												},
												"enable_pilot_http": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Enable HTTP/1.0 support.",
												},
												"disable_http_retry": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Disable http retry.",
												},
												"smart_dns": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "SmartDNS configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"istio_meta_dns_capture": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Enable dns proxy.",
															},
															"istio_meta_dns_auto_allocate": {
																Type:        schema.TypeBool,
																Computed:    true,
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
							Computed:    true,
							Description: "A list of associated tags.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag key.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag value.",
									},
									"passthrough": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Passthrough to other related product.",
									},
								},
							},
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTcmMeshRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tcm_mesh.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string][]*string)
	if v, ok := d.GetOk("mesh_id"); ok {
		meshIdSet := v.(*schema.Set).List()
		paramMap["MeshId"] = helper.InterfacesStringsPoint(meshIdSet)
	}

	if v, ok := d.GetOk("mesh_name"); ok {
		meshName := v.(*schema.Set).List()
		paramMap["MeshName"] = helper.InterfacesStringsPoint(meshName)
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsSet := v.(*schema.Set).List()
		paramMap["Tags"] = helper.InterfacesStringsPoint(tagsSet)
	}

	if v, ok := d.GetOk("mesh_cluster"); ok {
		meshClusterSet := v.(*schema.Set).List()
		paramMap["MeshCluster"] = helper.InterfacesStringsPoint(meshClusterSet)
	}

	service := TcmService{client: meta.(*TencentCloudClient).apiV3Conn}

	var meshList []*tcm.Mesh

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTcmMeshByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		meshList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(meshList))
	tmpList := make([]map[string]interface{}, 0, len(meshList))

	if meshList != nil {
		for _, mesh := range meshList {
			meshMap := map[string]interface{}{}

			if mesh.MeshId != nil {
				meshMap["mesh_id"] = mesh.MeshId
			}

			if mesh.DisplayName != nil {
				meshMap["display_name"] = mesh.DisplayName
			}

			if mesh.Version != nil {
				meshMap["version"] = mesh.Version
			}

			if mesh.Type != nil {
				meshMap["type"] = mesh.Type
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

				meshMap["config"] = []interface{}{configMap}
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

				meshMap["tag_list"] = tagListList
			}

			ids = append(ids, *mesh.MeshId)
			tmpList = append(tmpList, meshMap)
		}

		err := d.Set("mesh_list", tmpList)
		if err != nil {
			return err
		}
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
