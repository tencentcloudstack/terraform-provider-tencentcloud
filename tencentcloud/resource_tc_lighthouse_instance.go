/*
Provides a resource to create a lighthouse instance.

Example Usage

```hcl
resource "tencentcloud_lighthouse_instance" "lighthouse" {
  bundle_id    = "bundle2022_gen_01"
  blueprint_id = "lhbp-f1lkcd41"

  period     = 1
  renew_flag = "NOTIFY_AND_AUTO_RENEW"

  instance_name = "hello world"
  zone          = "ap-guangzhou-3"

  containers {
    container_image = "ccr.ccs.tencentyun.com/qcloud/nginx"
    container_name = "nginx"
    envs {
      key = "key"
      value = "value"
    }
    envs {
      key = "key2"
      value = "value2"
    }
    publish_ports {
      host_port = 80
      container_port = 80
      ip = "127.0.0.1"
      protocol = "tcp"
    }
    publish_ports {
      host_port = 36000
      container_port = 36000
      ip = "127.0.0.1"
      protocol = "tcp"
    }
    volumes {
      container_path = "/data"
      host_path = "/tmp"
    }
    volumes {
      container_path = "/var"
      host_path = "/tmp"
    }
    command = "ls -l"
  }

  containers {
    container_image = "ccr.ccs.tencentyun.com/qcloud/resty"
    container_name = "resty"
    envs {
      key = "key2"
      value = "value2"
    }
    publish_ports {
      host_port = 80
      container_port = 80
      ip = "127.0.0.1"
      protocol = "udp"
    }

    volumes {
      container_path = "/var"
      host_path = "/tmp"
    }
    command = "echo \"hello\""
  }
}
```

*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudLighthouseInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseInstanceCreate,
		Read:   resourceTencentCloudLighthouseInstanceRead,
		Delete: resourceTencentCloudLighthouseInstanceDelete,
		Update: resourceTencentCloudLighthouseInstanceUpdate,
		Schema: map[string]*schema.Schema{
			"bundle_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the Lighthouse package.",
			},
			"is_update_bundle_id_auto_voucher": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the voucher is deducted automatically when update bundle id. Value range: `true`: indicates automatic deduction of vouchers, `false`: does not automatically deduct vouchers. Default value: `false`.",
			},
			"blueprint_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the Lighthouse image.",
			},
			"period": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Subscription period in months. Valid values: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36, 48, 60.",
			},
			"renew_flag": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Auto-Renewal flag. Valid values: NOTIFY_AND_AUTO_RENEW: notify upon expiration and renew automatically; NOTIFY_AND_MANUAL_RENEW: notify upon expiration but do not renew automatically. You need to manually renew DISABLE_NOTIFY_AND_AUTO_RENEW: neither notify upon expiration nor renew automatically. " +
					"Default value: NOTIFY_AND_MANUAL_RENEW.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the Lighthouse instance.",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "List of availability zones. A random AZ is selected by default.",
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: "Whether the request is a dry run only." +
					"true: dry run only. The request will not create instance(s). A dry run can check whether all the required parameters are specified, whether the request format is right, whether the request exceeds service limits, and whether the specified CVMs are available. If the dry run fails, the corresponding error code will be returned.If the dry run succeeds, the RequestId will be returned." +
					"false (default value): send a normal request and create instance(s) if all the requirements are met.",
			},
			"client_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A unique string supplied by the client to ensure that the request is idempotent. Its maximum length is 64 ASCII characters. If this parameter is not specified, the idem-potency of the request cannot be guaranteed.",
			},
			"login_configuration": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Login password of the instance. It is only available for Windows instances. If it is not specified, it means that the user choose to set the login password after the instance creation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_generate_password": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "whether auto generate password. if false, need set password.",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Login password.",
						},
					},
				},
			},
			"permit_default_key_pair_login": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue([]string{"YES", "NO"}),
				Description:  "Whether to allow login using the default key pair. `YES`: allow login; `NO`: disable login. Default: `YES`.",
			},
			"isolate_data_disk": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to return the mounted data disk. `true`: returns both the instance and the mounted data disk; `false`: returns the instance and no longer returns its mounted data disk. Default: `true`.",
			},
			"containers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Configuration of the containers to create.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container_image": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Container image address.",
						},
						"container_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Container name.",
						},
						"envs": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of environment variables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Environment variable key.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Environment variable value.",
									},
								},
							},
						},
						"publish_ports": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of mappings of container ports and host ports.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host_port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Host port.",
									},
									"container_port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Container port.",
									},
									"ip": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "External IP. It defaults to 0.0.0.0.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The protocol defaults to tcp. Valid values: tcp, udp and sctp.",
									},
								},
							},
						},
						"volumes": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of container mount volumes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"container_path": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Container path.",
									},
									"host_path": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Host path.",
									},
								},
							},
						},
						"command": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The command to run.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudLighthouseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = lighthouse.NewCreateInstancesRequest()
		instanceId string
	)

	if v, ok := d.GetOk("bundle_id"); ok {
		request.BundleId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("blueprint_id"); ok {
		request.BlueprintId = helper.String(v.(string))
	}

	instanceChargePrepaid := lighthouse.InstanceChargePrepaid{}
	if v, ok := d.GetOk("period"); ok {
		instanceChargePrepaid.Period = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("renew_flag"); ok {
		instanceChargePrepaid.RenewFlag = helper.String(v.(string))
	}
	request.InstanceChargePrepaid = &instanceChargePrepaid

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("zone"); ok {
		request.Zones = append(request.Zones, helper.String(v.(string)))
	}

	if v, ok := d.GetOk("dry_run"); ok {
		request.DryRun = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("client_token"); ok {
		request.ClientToken = helper.String(v.(string))
	}

	if loginConfigurationMap, ok := helper.InterfacesHeadMap(d, "login_configuration"); ok {
		loginConfiguration := lighthouse.LoginConfiguration{}
		if v, ok := loginConfigurationMap["auto_generate_password"]; ok {
			loginConfiguration.AutoGeneratePassword = helper.String(v.(string))
		}
		if v, ok := loginConfigurationMap["password"]; ok {
			loginConfiguration.Password = helper.String(v.(string))
		}
		request.LoginConfiguration = &loginConfiguration
	}

	if v, ok := d.GetOk("containers"); ok {
		for _, container := range v.([]interface{}) {
			dockerContainerConfiguration := lighthouse.DockerContainerConfiguration{}
			containerMap := container.(map[string]interface{})
			if v, ok := containerMap["container_image"]; ok {
				dockerContainerConfiguration.ContainerImage = helper.String(v.(string))
			}
			if v, ok := containerMap["container_name"]; ok {
				dockerContainerConfiguration.ContainerName = helper.String(v.(string))
			}
			if v, ok := containerMap["envs"]; ok {
				for _, env := range v.([]interface{}) {
					containerEnv := lighthouse.ContainerEnv{}
					envMap := env.(map[string]interface{})
					if v, ok := envMap["key"]; ok {
						containerEnv.Key = helper.String(v.(string))
					}
					if v, ok := envMap["value"]; ok {
						containerEnv.Value = helper.String(v.(string))
					}
					dockerContainerConfiguration.Envs = append(dockerContainerConfiguration.Envs, &containerEnv)
				}
			}
			if v, ok := containerMap["publish_ports"]; ok {
				for _, publishPort := range v.([]interface{}) {
					dockerContainerPublishPort := lighthouse.DockerContainerPublishPort{}
					publishPortMap := publishPort.(map[string]interface{})
					if v, ok := publishPortMap["host_port"]; ok {
						dockerContainerPublishPort.HostPort = helper.IntInt64(v.(int))
					}
					if v, ok := publishPortMap["container_port"]; ok {
						dockerContainerPublishPort.ContainerPort = helper.IntInt64(v.(int))
					}
					if v, ok := publishPortMap["ip"]; ok {
						dockerContainerPublishPort.Ip = helper.String(v.(string))
					}
					if v, ok := publishPortMap["protocol"]; ok {
						dockerContainerPublishPort.Protocol = helper.String(v.(string))
					}
					dockerContainerConfiguration.PublishPorts = append(dockerContainerConfiguration.PublishPorts, &dockerContainerPublishPort)
				}
			}
			if v, ok := containerMap["volumes"]; ok {
				for _, volume := range v.([]interface{}) {
					dockerContainerVolume := lighthouse.DockerContainerVolume{}
					volumeMap := volume.(map[string]interface{})
					if v, ok := volumeMap["container_path"]; ok {
						dockerContainerVolume.ContainerPath = helper.String(v.(string))
					}
					if v, ok := volumeMap["host_path"]; ok {
						dockerContainerVolume.HostPath = helper.String(v.(string))
					}
					dockerContainerConfiguration.Volumes = append(dockerContainerConfiguration.Volumes, &dockerContainerVolume)
				}
			}
			if v, ok := containerMap["command"]; ok {
				dockerContainerConfiguration.Command = helper.String(v.(string))
			}
			request.Containers = append(request.Containers, &dockerContainerConfiguration)
		}
	}

	result, err := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().CreateInstances(request)

	if err != nil {
		log.Printf("[CRITAL]%s create lighthouse instance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *result.Response.InstanceIdSet[0]

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	lighthouseService := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := lighthouseService.DescribeLighthouseInstanceById(ctx, instanceId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if instance != nil && (*instance.InstanceState == "RUNNING") {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("lighthouse instance status is %s, retry...", *instance.InstanceState))
	})
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("permit_default_key_pair_login"); ok {
		permitLogin := v.(string)
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := lighthouseService.ModifyInstancesLoginKeyPairAttribute(ctx, instanceId, permitLogin)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update lighthouse instanceLoginKeyPair failed, reason:%+v", logId, err)
			return err
		}
	}

	d.SetId(instanceId)

	return resourceTencentCloudLighthouseInstanceRead(d, meta)
}

func resourceTencentCloudLighthouseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	lighthouseService := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	id := d.Id()

	instance, err := lighthouseService.DescribeLighthouseInstanceById(ctx, id)

	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		return fmt.Errorf("resource `lighthouse instance` %s does not exist", id)
	}

	if instance.BundleId != nil {
		_ = d.Set("bundle_id", instance.BundleId)
	}

	if instance.BlueprintId != nil {
		_ = d.Set("blueprint_id", instance.BlueprintId)
	}

	if instance.InstanceChargeType != nil {
		_ = d.Set("renew_flag", instance.RenewFlag)
	}

	if instance.InstanceName != nil {
		_ = d.Set("instance_name", instance.InstanceName)
	}

	if instance.Zone != nil {
		_ = d.Set("zone", instance.Zone)
	}

	instanceLoginKeyPair, err := lighthouseService.DescribeLighthouseInstanceLoginKeyPairById(ctx, id)
	if err != nil {
		return err
	}

	if instanceLoginKeyPair != nil && instanceLoginKeyPair.PermitLogin != nil {
		_ = d.Set("permit_default_key_pair_login", instanceLoginKeyPair.PermitLogin)
	}
	return nil
}

func resourceTencentCloudLighthouseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = lighthouse.NewModifyInstancesAttributeRequest()
	)
	id := d.Id()

	request.InstanceIds = append(request.InstanceIds, helper.String(id))

	if d.HasChange("instance_name") {
		if v, ok := d.GetOk("instance_name"); ok {
			request.InstanceName = helper.String(v.(string))
		}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().ModifyInstancesAttribute(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if err != nil {
			return err
		}

		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}
		err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
			instance, errRet := service.DescribeLighthouseInstanceById(ctx, id)
			if errRet != nil {
				return retryError(errRet, InternalError)
			}
			if instance.LatestOperationState == nil {
				return resource.RetryableError(fmt.Errorf("waiting for instance operation update"))
			}
			if *instance.LatestOperationState == "OPERATING" {
				return resource.RetryableError(fmt.Errorf("waiting for instance %s operation", id))
			}
			if *instance.LatestOperationState == "FAILED" {
				return resource.NonRetryableError(fmt.Errorf("failed operation"))
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	if d.HasChange("bundle_id") {
		_, new := d.GetChange("bundle_id")
		request := lighthouse.NewModifyInstancesBundleRequest()
		request.InstanceIds = helper.StringsStringsPoint([]string{id})
		request.BundleId = helper.String(new.(string))
		autoVoucher := d.Get("is_update_bundle_id_auto_voucher").(bool)
		request.AutoVoucher = &autoVoucher
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().ModifyInstancesBundle(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update lighthouse instanceModifyBundle failed, reason:%+v", logId, err)
			return err
		}

		service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

		conf := BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*readRetryTimeout, time.Second, service.LighthouseInstanceStateRefreshFunc(id, []string{}))

		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	if d.HasChange("blueprint_id") {
		_, new := d.GetChange("blueprint_id")
		request := lighthouse.NewResetInstanceRequest()
		request.InstanceId = helper.String(id)
		request.BlueprintId = helper.String(new.(string))

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().ResetInstance(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate lighthouse resetInstance failed, reason:%+v", logId, err)
			return err
		}

		service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

		conf := BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*readRetryTimeout, time.Second, service.LighthouseInstanceStateRefreshFunc(id, []string{}))

		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	if d.HasChange("period") {
		old, _ := d.GetChange("period")
		_ = d.Set("period", old)
		return fmt.Errorf("`period` do not support change now.")
	}

	if d.HasChange("renew_flag") {
		_, new := d.GetChange("renew_flag")
		request := lighthouse.NewModifyInstancesRenewFlagRequest()
		request.InstanceIds = helper.StringsStringsPoint([]string{id})
		request.RenewFlag = helper.String(new.(string))
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().ModifyInstancesRenewFlag(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate lighthouse modifyInstanceRenewFlag failed, reason:%+v", logId, err)
			return err
		}
		service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

		conf := BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*readRetryTimeout, time.Second, service.LighthouseInstanceStateRefreshFunc(id, []string{}))

		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	if d.HasChange("zone") {
		old, _ := d.GetChange("zone")
		_ = d.Set("zone", old)
		return fmt.Errorf("`zone` do not support change now.")
	}

	if d.HasChange("dry_run") {
		old, _ := d.GetChange("dry_run")
		_ = d.Set("dry_run", old)
		return fmt.Errorf("`dry_run` do not support change now.")
	}

	if d.HasChange("client_token") {
		old, _ := d.GetChange("client_token")
		_ = d.Set("client_token", old)
		return fmt.Errorf("`client_token` do not support change now.")
	}

	if d.HasChange("login_configuration.0.auto_generate_password") {
		old, _ := d.GetChange("login_configuration")
		_ = d.Set("login_configuration", old)
		return fmt.Errorf("`auto_generate_password` do not support change now.")
	}
	if d.HasChange("login_configuration.0.password") {
		_, new := d.GetChange("login_configuration")
		var newLoginConfiguration map[string]interface{}
		if len(new.([]interface{})) > 0 {
			newLoginConfiguration = new.([]interface{})[0].(map[string]interface{})
		}
		newPassword := newLoginConfiguration["password"].(string)
		request := lighthouse.NewResetInstancesPasswordRequest()
		request.InstanceIds = helper.StringsStringsPoint([]string{id})
		request.Password = helper.String(newPassword)
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().ResetInstancesPassword(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate lighthouse resetInstancesPassword failed, reason:%+v", logId, err)
			return err
		}
		service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

		conf := BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*readRetryTimeout, time.Second, service.LighthouseInstanceStateRefreshFunc(id, []string{}))

		if _, e := conf.WaitForState(); e != nil {
			return e
		}

		_ = d.Set("login_configuration", new)
	}

	if d.HasChange("containers") {
		old, _ := d.GetChange("containers")
		_ = d.Set("containers", old)
		return fmt.Errorf("`containers` do not support change now.")
	}

	if d.HasChange("permit_default_key_pair_login") {
		_, new := d.GetChange("permit_default_key_pair_login")
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := service.ModifyInstancesLoginKeyPairAttribute(ctx, id, new.(string))
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update lighthouse instanceLoginKeyPair failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudLighthouseInstanceRead(d, meta)
}

func resourceTencentCloudLighthouseInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}
	id := d.Id()
	isolateDataDisk := d.Get("isolate_data_disk").(bool)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if err := service.IsolateLighthouseInstanceById(ctx, id, isolateDataDisk); err != nil {
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeLighthouseInstanceById(ctx, id)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if instance.LatestOperationState == nil {
			return resource.RetryableError(fmt.Errorf("waiting for instance operation update"))
		}
		if *instance.LatestOperationState == "FAILED" {
			return resource.NonRetryableError(fmt.Errorf("failed operation"))
		}
		if *instance.InstanceState == "SHUTDOWN" && *instance.LatestOperationState != "OPERATING" {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("instance status is %s, retry...", *instance.InstanceState))
	})

	if err != nil {
		return err
	}

	if err := service.DeleteLighthouseInstanceById(ctx, id); err != nil {
		return err
	}
	return nil
}
