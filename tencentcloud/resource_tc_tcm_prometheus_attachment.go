/*
Provides a resource to create a tcm prometheus_attachment

Example Usage

```hcl
resource "tencentcloud_tcm_prometheus_attachment" "prometheus_attachment" {
  mesh_i_d = "mesh-xxxxxxxx"
  prometheus {
		vpc_id = "vpc-xxx"
		subnet_id = "subnet-xxx"
		region = "sh"
		instance_id = "prom-xxx"
		custom_prom {
			is_public_addr = false
			vpc_id = "vpc-xxx"
			url = "http://x.x.x.x:9090"
			auth_type = "none, basic"
			username = "test"
			password = "test"
		}

  }
}
```

Import

tcm prometheus_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_tcm_prometheus_attachment.prometheus_attachment prometheus_attachment_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcm/v20210413"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTcmPrometheusAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcmPrometheusAttachmentCreate,
		Read:   resourceTencentCloudTcmPrometheusAttachmentRead,
		Delete: resourceTencentCloudTcmPrometheusAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"mesh_i_d": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Mesh ID.",
			},

			"prometheus": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Prometheus configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Vpc id for TMP.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Subnet id for TMP.",
						},
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Region for TMP.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Existed TMP id, auto create TMP if empty.",
						},
						"custom_prom": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Third party prometheus.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_public_addr": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether it is public address, default false.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Vpc id.",
									},
									"url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Url of the prometheus.",
									},
									"auth_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Authentication type of the prometheus.",
									},
									"username": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Username of the prometheus, used in basic authentication type.",
									},
									"password": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Password of the prometheus, used in basic authentication type.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTcmPrometheusAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_prometheus_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tcm.NewLinkPrometheusRequest()
		response = tcm.NewLinkPrometheusResponse()
		meshID   string
	)
	if v, ok := d.GetOk("mesh_i_d"); ok {
		meshID = v.(string)
		request.MeshID = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "prometheus"); ok {
		prometheusConfig := tcm.PrometheusConfig{}
		if v, ok := dMap["vpc_id"]; ok {
			prometheusConfig.VpcId = helper.String(v.(string))
		}
		if v, ok := dMap["subnet_id"]; ok {
			prometheusConfig.SubnetId = helper.String(v.(string))
		}
		if v, ok := dMap["region"]; ok {
			prometheusConfig.Region = helper.String(v.(string))
		}
		if v, ok := dMap["instance_id"]; ok {
			prometheusConfig.InstanceId = helper.String(v.(string))
		}
		if customPromMap, ok := helper.InterfaceToMap(dMap, "custom_prom"); ok {
			customPromConfig := tcm.CustomPromConfig{}
			if v, ok := customPromMap["is_public_addr"]; ok {
				customPromConfig.IsPublicAddr = helper.Bool(v.(bool))
			}
			if v, ok := customPromMap["vpc_id"]; ok {
				customPromConfig.VpcId = helper.String(v.(string))
			}
			if v, ok := customPromMap["url"]; ok {
				customPromConfig.Url = helper.String(v.(string))
			}
			if v, ok := customPromMap["auth_type"]; ok {
				customPromConfig.AuthType = helper.String(v.(string))
			}
			if v, ok := customPromMap["username"]; ok {
				customPromConfig.Username = helper.String(v.(string))
			}
			if v, ok := customPromMap["password"]; ok {
				customPromConfig.Password = helper.String(v.(string))
			}
			prometheusConfig.CustomProm = &customPromConfig
		}
		request.Prometheus = &prometheusConfig
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTcmClient().LinkPrometheus(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tcm PrometheusAttachment failed, reason:%+v", logId, err)
		return err
	}

	meshID = *response.Response.MeshID
	d.SetId(meshID)

	return resourceTencentCloudTcmPrometheusAttachmentRead(d, meta)
}

func resourceTencentCloudTcmPrometheusAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_prometheus_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcmService{client: meta.(*TencentCloudClient).apiV3Conn}

	prometheusAttachmentId := d.Id()

	PrometheusAttachment, err := service.DescribeTcmPrometheusAttachmentById(ctx, meshID)
	if err != nil {
		return err
	}

	if PrometheusAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TcmPrometheusAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if PrometheusAttachment.MeshID != nil {
		_ = d.Set("mesh_i_d", PrometheusAttachment.MeshID)
	}

	if PrometheusAttachment.Prometheus != nil {
		prometheusMap := map[string]interface{}{}

		if PrometheusAttachment.Prometheus.VpcId != nil {
			prometheusMap["vpc_id"] = PrometheusAttachment.Prometheus.VpcId
		}

		if PrometheusAttachment.Prometheus.SubnetId != nil {
			prometheusMap["subnet_id"] = PrometheusAttachment.Prometheus.SubnetId
		}

		if PrometheusAttachment.Prometheus.Region != nil {
			prometheusMap["region"] = PrometheusAttachment.Prometheus.Region
		}

		if PrometheusAttachment.Prometheus.InstanceId != nil {
			prometheusMap["instance_id"] = PrometheusAttachment.Prometheus.InstanceId
		}

		if PrometheusAttachment.Prometheus.CustomProm != nil {
			customPromMap := map[string]interface{}{}

			if PrometheusAttachment.Prometheus.CustomProm.IsPublicAddr != nil {
				customPromMap["is_public_addr"] = PrometheusAttachment.Prometheus.CustomProm.IsPublicAddr
			}

			if PrometheusAttachment.Prometheus.CustomProm.VpcId != nil {
				customPromMap["vpc_id"] = PrometheusAttachment.Prometheus.CustomProm.VpcId
			}

			if PrometheusAttachment.Prometheus.CustomProm.Url != nil {
				customPromMap["url"] = PrometheusAttachment.Prometheus.CustomProm.Url
			}

			if PrometheusAttachment.Prometheus.CustomProm.AuthType != nil {
				customPromMap["auth_type"] = PrometheusAttachment.Prometheus.CustomProm.AuthType
			}

			if PrometheusAttachment.Prometheus.CustomProm.Username != nil {
				customPromMap["username"] = PrometheusAttachment.Prometheus.CustomProm.Username
			}

			if PrometheusAttachment.Prometheus.CustomProm.Password != nil {
				customPromMap["password"] = PrometheusAttachment.Prometheus.CustomProm.Password
			}

			prometheusMap["custom_prom"] = []interface{}{customPromMap}
		}

		_ = d.Set("prometheus", []interface{}{prometheusMap})
	}

	return nil
}

func resourceTencentCloudTcmPrometheusAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_prometheus_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcmService{client: meta.(*TencentCloudClient).apiV3Conn}
	prometheusAttachmentId := d.Id()

	if err := service.DeleteTcmPrometheusAttachmentById(ctx, meshID); err != nil {
		return err
	}

	return nil
}
