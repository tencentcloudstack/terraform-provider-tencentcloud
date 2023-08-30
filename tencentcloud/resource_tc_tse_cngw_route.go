/*
Provides a resource to create a tse cngw_route

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_tse_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "tf_tse_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_tse_cngw_gateway" "cngw_gateway" {
  description                = "terraform test1"
  enable_cls                 = true
  engine_region              = "ap-guangzhou"
  feature_version            = "STANDARD"
  gateway_version            = "2.5.1"
  ingress_class_name         = "tse-nginx-ingress"
  internet_max_bandwidth_out = 0
  name                       = "terraform-gateway1"
  trade_type                 = 0
  type                       = "kong"

  node_config {
    number        = 2
    specification = "1c2g"
  }

  vpc_config {
    subnet_id = tencentcloud_subnet.subnet.id
    vpc_id    = tencentcloud_vpc.vpc.id
  }

  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tse_cngw_service" "cngw_service" {
  gateway_id    = tencentcloud_tse_cngw_gateway.cngw_gateway.id
  name          = "terraform-test"
  path          = "/test"
  protocol      = "http"
  retries       = 5
  timeout       = 60000
  upstream_type = "HostIP"

  upstream_info {
    algorithm             = "round-robin"
    auto_scaling_cvm_port = 0
    host                  = "arunma.cn"
    port                  = 8012
    slow_start            = 0
  }
}

resource "tencentcloud_tse_cngw_route" "cngw_route" {
  destination_ports = []
  force_https       = false
  gateway_id        = tencentcloud_tse_cngw_gateway.cngw_gateway.id
  hosts = [
    "192.168.0.1:9090",
  ]
  https_redirect_status_code = 426
  paths = [
    "/user",
  ]
  headers {
	key = "req"
	value = "terraform"
  }
  preserve_host = false
  protocols = [
    "http",
    "https",
  ]
  route_name = "terraform-route"
  service_id = tencentcloud_tse_cngw_service.cngw_service.service_id
  strip_path = true
}

```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTseCngwRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTseCngwRouteCreate,
		Read:   resourceTencentCloudTseCngwRouteRead,
		Update: resourceTencentCloudTseCngwRouteUpdate,
		Delete: resourceTencentCloudTseCngwRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "gateway ID.",
			},

			"service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of the service which the route belongs to.",
			},

			"route_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "the name of the route, unique in the instance.",
			},

			"methods": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "route methods. Reference value:`GET`,`POST`,`DELETE`,`PUT`,`OPTIONS`,`PATCH`,`HEAD`,`ANY`,`TRACE`,`COPY`,`MOVE`,`PROPFIND`,`PROPPATCH`,`MKCOL`,`LOCK`,`UNLOCK`.",
			},

			"hosts": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "host list.",
			},

			"paths": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "path list.",
			},

			"protocols": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "the protocol list of route.Reference value:`https`,`http`.",
			},

			"preserve_host": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "whether to keep the host when forwarding to the backend.",
			},

			"https_redirect_status_code": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "https redirection status code.",
			},

			"strip_path": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "whether to strip path when forwarding to the backend.",
			},

			"force_https": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "whether to enable forced HTTPS, no longer use.",
			},

			"destination_ports": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "destination port for Layer 4 matching.",
			},

			"headers": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "the headers of route.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "key of header.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "value of header.",
						},
					},
				},
			},

			"route_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "the id of the route, unique in the instance.",
			},
		},
	}
}

func resourceTencentCloudTseCngwRouteCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_route.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tse.NewCreateCloudNativeAPIGatewayRouteRequest()
		gatewayId string
		serviceID string
		routeName string
	)
	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
		request.GatewayId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_id"); ok {
		serviceID = v.(string)
		request.ServiceID = helper.String(v.(string))
	}

	if v, ok := d.GetOk("route_name"); ok {
		routeName = v.(string)
		request.RouteName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("methods"); ok {
		methodsSet := v.(*schema.Set).List()
		for i := range methodsSet {
			methods := methodsSet[i].(string)
			request.Methods = append(request.Methods, &methods)
		}
	}

	if v, ok := d.GetOk("hosts"); ok {
		hostsSet := v.(*schema.Set).List()
		for i := range hostsSet {
			hosts := hostsSet[i].(string)
			request.Hosts = append(request.Hosts, &hosts)
		}
	}

	if v, ok := d.GetOk("paths"); ok {
		pathsSet := v.(*schema.Set).List()
		for i := range pathsSet {
			paths := pathsSet[i].(string)
			request.Paths = append(request.Paths, &paths)
		}
	}

	if v, ok := d.GetOk("protocols"); ok {
		protocolsSet := v.(*schema.Set).List()
		for i := range protocolsSet {
			protocols := protocolsSet[i].(string)
			request.Protocols = append(request.Protocols, &protocols)
		}
	}

	if v, ok := d.GetOkExists("preserve_host"); ok {
		request.PreserveHost = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("https_redirect_status_code"); ok {
		request.HttpsRedirectStatusCode = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("strip_path"); ok {
		request.StripPath = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("force_https"); ok {
		request.ForceHttps = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("destination_ports"); ok {
		destinationPortsSet := v.(*schema.Set).List()
		for i := range destinationPortsSet {
			destinationPorts := destinationPortsSet[i].(int)
			request.DestinationPorts = append(request.DestinationPorts, helper.IntUint64(destinationPorts))
		}
	}

	if v, ok := d.GetOk("headers"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			kVMapping := tse.KVMapping{}
			if v, ok := dMap["key"]; ok {
				kVMapping.Key = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				kVMapping.Value = helper.String(v.(string))
			}
			request.Headers = append(request.Headers, &kVMapping)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().CreateCloudNativeAPIGatewayRoute(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tse cngwRoute failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(gatewayId + FILED_SP + serviceID + FILED_SP + routeName)

	return resourceTencentCloudTseCngwRouteRead(d, meta)
}

func resourceTencentCloudTseCngwRouteRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_route.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	serviceID := idSplit[1]
	routeName := idSplit[2]

	cngwRoute, err := service.DescribeTseCngwRouteById(ctx, gatewayId, serviceID, routeName)
	if err != nil {
		return err
	}

	if cngwRoute == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TseCngwRoute` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("gateway_id", gatewayId)
	_ = d.Set("service_id", serviceID)
	_ = d.Set("route_name", routeName)

	if cngwRoute.Methods != nil {
		_ = d.Set("methods", cngwRoute.Methods)
	}

	if cngwRoute.Hosts != nil {
		_ = d.Set("hosts", cngwRoute.Hosts)
	}

	if cngwRoute.Paths != nil {
		_ = d.Set("paths", cngwRoute.Paths)
	}

	if cngwRoute.Protocols != nil {
		_ = d.Set("protocols", cngwRoute.Protocols)
	}

	if cngwRoute.PreserveHost != nil {
		_ = d.Set("preserve_host", cngwRoute.PreserveHost)
	}

	if cngwRoute.HttpsRedirectStatusCode != nil {
		_ = d.Set("https_redirect_status_code", cngwRoute.HttpsRedirectStatusCode)
	}

	if cngwRoute.StripPath != nil {
		_ = d.Set("strip_path", cngwRoute.StripPath)
	}

	if cngwRoute.ForceHttps != nil {
		_ = d.Set("force_https", cngwRoute.ForceHttps)
	}

	if cngwRoute.DestinationPorts != nil {
		_ = d.Set("destination_ports", cngwRoute.DestinationPorts)
	}

	// if cngwRoute.Headers != nil {
	// 	headersList := []interface{}{}
	// 	for _, headers := range cngwRoute.Headers {
	// 		headersMap := map[string]interface{}{}

	// 		if cngwRoute.Headers.Key != nil {
	// 			headersMap["key"] = cngwRoute.Headers.Key
	// 		}

	// 		if cngwRoute.Headers.Value != nil {
	// 			headersMap["value"] = cngwRoute.Headers.Value
	// 		}

	// 		headersList = append(headersList, headersMap)
	// 	}

	// 	_ = d.Set("headers", headersList)

	// }

	if cngwRoute.ID != nil {
		_ = d.Set("route_id", cngwRoute.ID)
	}

	return nil
}

func resourceTencentCloudTseCngwRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_route.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := tse.NewModifyCloudNativeAPIGatewayRouteRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	serviceID := idSplit[1]
	routeName := idSplit[2]

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}
	cngwRoute, err := service.DescribeTseCngwRouteById(ctx, gatewayId, serviceID, routeName)
	if err != nil {
		return err
	}

	if cngwRoute == nil {
		return fmt.Errorf("The result of querying %s is empty", routeName)
	}

	request.GatewayId = &gatewayId
	request.ServiceID = &serviceID
	request.RouteName = &routeName
	request.RouteID = cngwRoute.ID

	immutableArgs := []string{"gateway_id", "service_id", "route_name"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("methods") {
		if v, ok := d.GetOk("methods"); ok {
			methodsSet := v.(*schema.Set).List()
			for i := range methodsSet {
				methods := methodsSet[i].(string)
				request.Methods = append(request.Methods, &methods)
			}
		}
	}

	if d.HasChange("hosts") {
		if v, ok := d.GetOk("hosts"); ok {
			hostsSet := v.(*schema.Set).List()
			for i := range hostsSet {
				hosts := hostsSet[i].(string)
				request.Hosts = append(request.Hosts, &hosts)
			}
		}
	}

	if d.HasChange("paths") {
		if v, ok := d.GetOk("paths"); ok {
			pathsSet := v.(*schema.Set).List()
			for i := range pathsSet {
				paths := pathsSet[i].(string)
				request.Paths = append(request.Paths, &paths)
			}
		}
	}

	if d.HasChange("protocols") {
		if v, ok := d.GetOk("protocols"); ok {
			protocolsSet := v.(*schema.Set).List()
			for i := range protocolsSet {
				protocols := protocolsSet[i].(string)
				request.Protocols = append(request.Protocols, &protocols)
			}
		}
	}

	if d.HasChange("preserve_host") {
		if v, ok := d.GetOkExists("preserve_host"); ok {
			request.PreserveHost = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("https_redirect_status_code") {
		if v, ok := d.GetOkExists("https_redirect_status_code"); ok {
			request.HttpsRedirectStatusCode = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("strip_path") {
		if v, ok := d.GetOkExists("strip_path"); ok {
			request.StripPath = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("force_https") {
		if v, ok := d.GetOkExists("force_https"); ok {
			request.ForceHttps = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("destination_ports") {
		if v, ok := d.GetOk("destination_ports"); ok {
			destinationPortsSet := v.(*schema.Set).List()
			for i := range destinationPortsSet {
				destinationPorts := destinationPortsSet[i].(int)
				request.DestinationPorts = append(request.DestinationPorts, helper.IntUint64(destinationPorts))
			}
		}
	}

	if d.HasChange("headers") {
		if v, ok := d.GetOk("headers"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				kVMapping := tse.KVMapping{}
				if v, ok := dMap["key"]; ok {
					kVMapping.Key = helper.String(v.(string))
				}
				if v, ok := dMap["value"]; ok {
					kVMapping.Value = helper.String(v.(string))
				}
				request.Headers = append(request.Headers, &kVMapping)
			}
		}
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().ModifyCloudNativeAPIGatewayRoute(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tse cngwRoute failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTseCngwRouteRead(d, meta)
}

func resourceTencentCloudTseCngwRouteDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_route.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	// serviceID := idSplit[1]
	routeName := idSplit[2]

	if err := service.DeleteTseCngwRouteById(ctx, gatewayId, routeName); err != nil {
		return err
	}

	return nil
}
