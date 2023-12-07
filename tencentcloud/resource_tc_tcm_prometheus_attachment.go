package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcm/v20210413"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTcmPrometheusAttachment() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTcmPrometheusAttachmentRead,
		Create: resourceTencentCloudTcmPrometheusAttachmentCreate,
		Delete: resourceTencentCloudTcmPrometheusAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"mesh_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Mesh ID.",
			},

			"prometheus": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				ForceNew:    true,
				Description: "Prometheus configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Vpc id for TMP.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Subnet id for TMP.",
						},
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Region for TMP.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Existed TMP id, auto create TMP if empty.",
						},
						"custom_prom": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Third party prometheus.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_public_addr": {
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "Whether it is public address, default false.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
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
										Computed:    true,
										Description: "Username of the prometheus, used in basic authentication type.",
									},
									"password": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Sensitive:   true,
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
		request = tcm.NewLinkPrometheusRequest()
		meshID  string
	)

	if v, ok := d.GetOk("mesh_id"); ok {
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
		if CustomPromMap, ok := helper.InterfaceToMap(dMap, "custom_prom"); ok {
			customPromConfig := tcm.CustomPromConfig{}
			if v, ok := CustomPromMap["is_public_addr"]; ok {
				customPromConfig.IsPublicAddr = helper.Bool(v.(bool))
			}
			if v, ok := CustomPromMap["vpc_id"]; ok {
				customPromConfig.VpcId = helper.String(v.(string))
			}
			if v, ok := CustomPromMap["url"]; ok {
				customPromConfig.Url = helper.String(v.(string))
			}
			if v, ok := CustomPromMap["auth_type"]; ok {
				customPromConfig.AuthType = helper.String(v.(string))
			}
			if v, ok := CustomPromMap["username"]; ok {
				customPromConfig.Username = helper.String(v.(string))
			}
			if v, ok := CustomPromMap["password"]; ok {
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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tcm prometheusAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(meshID)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TcmService{client: meta.(*TencentCloudClient).apiV3Conn}
	err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
		mesh, errRet := service.DescribeTcmMesh(ctx, meshID)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if mesh.Mesh.Status == nil || mesh.Mesh.Status.TPS == nil {
			return nil
		}
		if *mesh.Mesh.Status.TPS.State == "PENDING" {
			return resource.RetryableError(fmt.Errorf("mesh status is %v, retry...", *mesh.Mesh.State))
		}
		return nil
	})
	if err != nil {
		return err
	}

	return resourceTencentCloudTcmPrometheusAttachmentRead(d, meta)
}

func resourceTencentCloudTcmPrometheusAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcm_prometheus_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcmService{client: meta.(*TencentCloudClient).apiV3Conn}

	meshId := d.Id()

	response, err := service.DescribeTcmMesh(ctx, meshId)

	if err != nil {
		return err
	}

	mesh := response.Mesh
	if mesh == nil {
		d.SetId("")
		return fmt.Errorf("resource `prometheusAttachment` %s does not exist", meshId)
	}

	if mesh.MeshId != nil {
		_ = d.Set("mesh_id", mesh.MeshId)
	}

	prometheus := mesh.Config.Prometheus
	if prometheus != nil {
		prometheusMap := map[string]interface{}{}
		if prometheus.VpcId != nil {
			prometheusMap["vpc_id"] = prometheus.VpcId
		}
		if prometheus.SubnetId != nil {
			prometheusMap["subnet_id"] = prometheus.SubnetId
		}
		if prometheus.Region != nil {
			prometheusMap["region"] = prometheus.Region
		}
		if prometheus.InstanceId != nil {
			prometheusMap["instance_id"] = prometheus.InstanceId
		}
		if prometheus.CustomProm != nil {
			customPromMap := map[string]interface{}{}
			if prometheus.CustomProm.IsPublicAddr != nil {
				customPromMap["is_public_addr"] = prometheus.CustomProm.IsPublicAddr
			}
			if prometheus.CustomProm.VpcId != nil {
				customPromMap["vpc_id"] = prometheus.CustomProm.VpcId
			}
			if prometheus.CustomProm.Url != nil {
				customPromMap["url"] = prometheus.CustomProm.Url
			}
			if prometheus.CustomProm.AuthType != nil {
				customPromMap["auth_type"] = prometheus.CustomProm.AuthType
			}
			if prometheus.CustomProm.Username != nil {
				customPromMap["username"] = prometheus.CustomProm.Username
			}
			if prometheus.CustomProm.Password != nil {
				customPromMap["password"] = prometheus.CustomProm.Password
			}

			prometheusMap["custom_prom"] = []interface{}{customPromMap}
		}

		_ = d.Set("prometheus", []interface{}{prometheusMap})
	}
	prom := mesh.Status.TPS
	if prom != nil && *prom.Type == "tmp" {
		prometheusMap := map[string]interface{}{}
		if prom.VpcId != nil {
			prometheusMap["vpc_id"] = prom.VpcId
		}
		// if prometheus.SubnetId != nil {
		// 	prometheusMap["subnet_id"] = prometheus.SubnetId
		// }
		if prom.Region != nil {
			prometheusMap["region"] = prom.Region
		}
		if prom.InstanceId != nil {
			prometheusMap["instance_id"] = prom.InstanceId
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

	meshID := d.Id()

	if err := service.DeleteTcmPrometheusAttachmentById(ctx, meshID); err != nil {
		return err
	}

	return nil
}
