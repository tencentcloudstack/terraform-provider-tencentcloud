/*
Provide a resource to configure kubernetes cluster authentication info.

~> **NOTE:** Only available for cluster version >= 1.20

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "cluster_cidr" {
  default = "172.16.0.0/16"
}

variable "default_instance_type" {
  default = "S1.SMALL1"
}

data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  cluster_cidr            = "10.31.0.0/16"
  cluster_max_pod_num     = 32
  cluster_name            = "keep"
  cluster_desc            = "test cluster desc"
  cluster_version         = "1.20.6"
  cluster_max_service_num = 32

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
  }

  cluster_deploy_type = "MANAGED_CLUSTER"
}

resource "tencentcloud_kubernetes_auth_attachment" "test_auth_attach" {
  cluster_id                           = tencentcloud_kubernetes_cluster.managed_cluster.id
  jwks_uri                             = "https://${tencentcloud_kubernetes_cluster.managed_cluster.id}.ccs.tencent-cloud.com/openid/v1/jwks"
  issuer                               = "https://${tencentcloud_kubernetes_cluster.managed_cluster.id}.ccs.tencent-cloud.com"
  auto_create_discovery_anonymous_auth = true
}
```

Use the TKE default issuer and jwks_uri

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "cluster_cidr" {
  default = "172.16.0.0/16"
}

variable "default_instance_type" {
  default = "S1.SMALL1"
}

data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  cluster_cidr            = "10.31.0.0/16"
  cluster_max_pod_num     = 32
  cluster_name            = "keep"
  cluster_desc            = "test cluster desc"
  cluster_version         = "1.20.6"
  cluster_max_service_num = 32

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
  }

  cluster_deploy_type = "MANAGED_CLUSTER"
}

# if you want to use tke default issuer and jwks_uri, please set use_tke_default to true and set issuer to empty string.
resource "tencentcloud_kubernetes_auth_attachment" "test_use_tke_default_auth_attach" {
  cluster_id                           = tencentcloud_kubernetes_cluster.managed_cluster.id
  auto_create_discovery_anonymous_auth = true
  use_tke_default                      = true
}
```

Use OIDC Config
```
resource "tencentcloud_kubernetes_auth_attachment" "test_auth_attach" {
  cluster_id                              = tencentcloud_kubernetes_cluster.managed_cluster.id
  use_tke_default                         = true
  auto_create_discovery_anonymous_auth    = true
  auto_create_oidc_config                 = true
  auto_install_pod_identity_webhook_addon = true
}

data "tencentcloud_cam_oidc_config" "oidc_config" {
  name       = tencentcloud_kubernetes_cluster.managed_cluster.id
  depends_on = [
    tencentcloud_kubernetes_auth_attachment.test_auth_attach
  ]
}

output "identity_key" {
  value = data.tencentcloud_cam_oidc_config.oidc_config.identity_key
}

output "identity_url" {
  value = data.tencentcloud_cam_oidc_config.oidc_config.identity_url
}

```

*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTKEAuthAttachment() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of clusters.",
			},
			"use_tke_default": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"issuer", "jwks_uri"},
				Description:   "If set to `true`, the issuer and jwks_uri will be generated automatically by tke, please do not set issuer and jwks_uri.",
			},
			"issuer": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"use_tke_default"},
				Description:   "Specify service-account-issuer. If use_tke_default is set to `true`, please do not set this field.",
			},
			"jwks_uri": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"use_tke_default"},
				Description:   "Specify service-account-jwks-uri. If use_tke_default is set to `true`, please do not set this field.",
			},
			"auto_create_discovery_anonymous_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to `true`, the rbac rule will be created automatically which allow anonymous user to access '/.well-known/openid-configuration' and '/openid/v1/jwks'.",
			},
			"auto_create_oidc_config": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Creating an identity provider.",
			},
			"auto_create_client_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Creating ClientId of the identity provider.",
			},
			"auto_install_pod_identity_webhook_addon": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Creating the PodIdentityWebhook component. if `auto_create_oidc_config` is true, this field must set true.",
			},
			"tke_default_issuer": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The default issuer of tke. If use_tke_default is set to `true`, this parameter will be set to the default value.",
			},
			"tke_default_jwks_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The default jwks_uri of tke. If use_tke_default is set to `true`, this parameter will be set to the default value.",
			},
		},
		Create: resourceTencentCloudTKEAuthAttachmentCreate,
		Update: resourceTencentCloudTKEAuthAttachmentUpdate,
		Read:   resourceTencentCloudTKEAuthAttachmentRead,
		Delete: resourceTencentCloudTKEAuthAttachmentDelete,
	}
}

func resourceTencentCloudTKEAuthAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.resource_tc_kubernetes_auth_attachment.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	id := d.Get("cluster_id").(string)

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	request := tke.NewModifyClusterAuthenticationOptionsRequest()
	request.ClusterId = &id
	request.ServiceAccounts = &tke.ServiceAccountAuthenticationOptions{}

	if v, ok := d.GetOk("auto_create_discovery_anonymous_auth"); ok {
		request.ServiceAccounts.AutoCreateDiscoveryAnonymousAuth = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("use_tke_default"); ok && v.(bool) {
		request.ServiceAccounts.UseTKEDefault = helper.Bool(true)
	} else {
		if v, ok := d.GetOk("issuer"); ok {
			request.ServiceAccounts.Issuer = helper.String(v.(string))
		}
		if v, ok := d.GetOk("jwks_uri"); ok {
			request.ServiceAccounts.JWKSURI = helper.String(v.(string))
		}
	}
	request.OIDCConfig = &tke.OIDCConfigAuthenticationOptions{}
	if v, ok := d.GetOkExists("auto_create_oidc_config"); ok {
		request.OIDCConfig.AutoCreateOIDCConfig = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("auto_create_client_id"); ok {
		clientsSet := v.(*schema.Set).List()
		for i := range clientsSet {
			client := clientsSet[i].(string)
			request.OIDCConfig.AutoCreateClientId = append(request.OIDCConfig.AutoCreateClientId, &client)
		}
	}
	if v, ok := d.GetOkExists("auto_install_pod_identity_webhook_addon"); ok {
		request.OIDCConfig.AutoInstallPodIdentityWebhookAddon = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		err := service.ModifyClusterAuthenticationOptions(ctx, request)
		if err != nil {
			return retryError(err, tke.RESOURCEUNAVAILABLE_CLUSTERSTATE)
		}
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceTencentCloudTKEAuthAttachmentRead(d, meta)
}
func resourceTencentCloudTKEAuthAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.resource_tc_kubernetes_auth_attachment.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	info, oidc, err := service.WaitForAuthenticationOptionsUpdateSuccess(ctx, id)

	if err != nil {
		d.SetId("")
		return err
	}

	d.SetId(id)

	if v, ok := d.GetOk("use_tke_default"); ok && v.(bool) {
		_ = d.Set("tke_default_issuer", info.Issuer)
		_ = d.Set("tke_default_jwks_uri", info.JWKSURI)
	} else {
		_ = d.Set("jwks_uri", info.JWKSURI)
		_ = d.Set("issuer", info.Issuer)
	}

	if oidc.AutoCreateOIDCConfig != nil {
		_ = d.Set("auto_create_oidc_config", oidc.AutoCreateOIDCConfig)
	}

	if oidc.AutoCreateClientId != nil {
		_ = d.Set("auto_create_client_id", oidc.AutoCreateClientId)
	}

	if oidc.AutoInstallPodIdentityWebhookAddon != nil {
		_ = d.Set("auto_install_pod_identity_webhook_addon", oidc.AutoInstallPodIdentityWebhookAddon)
	}
	return nil
}

func resourceTencentCloudTKEAuthAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.resource_tc_kubernetes_auth_attachment.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	request := tke.NewModifyClusterAuthenticationOptionsRequest()
	request.ClusterId = &id
	request.ServiceAccounts = &tke.ServiceAccountAuthenticationOptions{}

	useTkeDefault := false
	if v, ok := d.GetOk("use_tke_default"); ok {
		request.ServiceAccounts.UseTKEDefault = helper.Bool(v.(bool))
		useTkeDefault = v.(bool)
	}

	if !useTkeDefault {
		if d.HasChange("jwks_uri") {
			request.ServiceAccounts.JWKSURI = helper.String(d.Get("jwks_uri").(string))
		}
		if d.HasChange("issuer") {
			issuer := d.Get("issuer").(string)
			request.ServiceAccounts.Issuer = helper.String(issuer)
		}
	}

	if d.HasChange("auto_create_discovery_anonymous_auth") {
		if v, ok := d.GetOk("auto_create_discovery_anonymous_auth"); ok {
			request.ServiceAccounts.AutoCreateDiscoveryAnonymousAuth = helper.Bool(v.(bool))
		}
	}

	request.OIDCConfig = &tke.OIDCConfigAuthenticationOptions{}
	if d.HasChange("auto_create_oidc_config") {
		if v, ok := d.GetOkExists("auto_create_oidc_config"); ok {
			request.OIDCConfig.AutoCreateOIDCConfig = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("auto_create_client_id") {
		if v, ok := d.GetOkExists("auto_create_client_id"); ok {
			clientsSet := v.(*schema.Set).List()
			for i := range clientsSet {
				client := clientsSet[i].(string)
				request.OIDCConfig.AutoCreateClientId = append(request.OIDCConfig.AutoCreateClientId, &client)
			}
		}
	}

	if d.HasChange("auto_install_pod_identity_webhook_addon") {
		if v, ok := d.GetOkExists("auto_install_pod_identity_webhook_addon"); ok {
			request.OIDCConfig.AutoInstallPodIdentityWebhookAddon = helper.Bool(v.(bool))
		}
	}

	err := resource.Retry(3*writeRetryTimeout, func() *resource.RetryError {
		err := service.ModifyClusterAuthenticationOptions(ctx, request)
		if err != nil {
			return retryError(err, tke.RESOURCEUNAVAILABLE_CLUSTERSTATE)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudTKEAuthAttachmentRead(d, meta)
}

func resourceTencentCloudTKEAuthAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.resource_tc_kubernetes_auth_attachment.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	request := tke.NewModifyClusterAuthenticationOptionsRequest()
	request.ClusterId = &id
	request.ServiceAccounts = &tke.ServiceAccountAuthenticationOptions{
		JWKSURI: helper.String(""),
		Issuer:  helper.String(DefaultAuthenticationOptionsIssuer),
	}

	if err := service.ModifyClusterAuthenticationOptions(ctx, request); err != nil {
		return err
	}

	_, _, err := service.WaitForAuthenticationOptionsUpdateSuccess(ctx, id)

	if err != nil {
		return err
	}

	return nil
}
