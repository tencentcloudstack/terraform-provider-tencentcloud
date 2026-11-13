// registry.go is the central, SDKv2-style manifest of every
// terraform-plugin-framework reference shipped by this provider.
//
// Every reference is declared as a `map[name]factory` literal, mirroring
// tencentcloud/provider.go's ResourcesMap / DataSourcesMap:
//
//	"tencentcloud_config_aggregate_compliance_pack": config.NewConfigAggregateCompliancePackResource,
//
// Why a map instead of a slice:
//   - the name Terraform reports for the reference (the type name users write
//     in HCL, or the bare function name) sits next to the constructor, so the
//     file can be searched by name instead of by Go symbol;
//   - duplicate registrations become a compile error, because Go rejects
//     duplicate map keys, instead of the runtime "Duplicate Resource Type
//     Defined" diagnostic raised when the provider schema is first requested;
//   - the key doubles as the expected Metadata TypeName (the Metadata Name for
//     provider-defined functions), which registry_test.go verifies for every
//     entry.
//
// To add a reference: add the services subpackage import (if not already
// present) and exactly one entry in the matching map. No edit to provider.go is
// required.
package framework

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ags"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/antiddos"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/bdrc"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cam"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/config"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dlc"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ioa"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/kms"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mongodb"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ssl"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ssm"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tco"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// resourcesMap maps every framework managed resource type name to its factory.
var resourcesMap = map[string]func() resource.Resource{
	"tencentcloud_ags_sandbox_tool":                 ags.NewAgsSandboxToolResource,
	"tencentcloud_cam_update_password_rule_config":  cam.NewCamUpdatePasswordRuleConfigResource,
	"tencentcloud_config_aggregate_compliance_pack": config.NewConfigAggregateCompliancePackResource,
	"tencentcloud_monitor_dispense_external_rule":   monitor.NewMonitorDispenseExternalRuleResource,
	"tencentcloud_ssm_secret_version_v2":            ssm.NewSsmSecretVersionV2Resource,
	"tencentcloud_ioa_delivery_server":              ioa.NewIoaDeliveryServerResource,
}

// dataSourcesMap maps every framework data source type name to its factory.
var dataSourcesMap = map[string]func() datasource.DataSource{
	"tencentcloud_config_list_aggregate_compliance_packs": config.NewConfigListAggregateCompliancePacksDataSource,
}

// functionsMap maps every framework function name to its factory. Function
// names carry no provider prefix (calls already read
// `provider::tencentcloud::<name>`), so the key is the bare name reported by
// Metadata, for example "parse_resource_id".
var functionsMap = map[string]func() function.Function{}

// ephemeralResourcesMap maps every framework ephemeral resource type name to its
// factory.
var ephemeralResourcesMap = map[string]func() ephemeral.EphemeralResource{
	"tencentcloud_ssm_secret_version": ssm.NewSsmSecretVersionEphemeralResource,
}

// listResourcesMap maps every framework list resource type name to its factory.
var listResourcesMap = map[string]func() list.ListResource{}

// actionsMap maps every framework action type name to its factory.
var actionsMap = map[string]func() action.Action{
	"tencentcloud_antiddos_unblock_resources":                antiddos.NewAntiddosUnblockResources,
	"tencentcloud_bdrc_run_copy_pair_tasks":                  bdrc.NewBdrcRunCopyPairTasks,
	"tencentcloud_dlc_initialize_tc_lake":                    dlc.NewDlcInitializeTCLake,
	"tencentcloud_kms_instance":                              kms.NewKmsInstance,
	"tencentcloud_mongodb_restore_db_instance":               mongodb.NewMongodbRestoreDbInstance,
	"tencentcloud_organization_provision_role_configuration": tco.NewTcoProvisionRoleConfiguration,
	"tencentcloud_ssl_upload_certificate_instance":           ssl.NewSslUploadCertificateInstance,
	"tencentcloud_ssm_service":                               ssm.NewSsmServiceAction,
	"tencentcloud_teo_confirm_origin_acl_update":             teo.NewTeoConfirmOriginAclUpdate,
}

// factoriesOf converts a name → factory registry map into the factory slice the
// framework provider callbacks expect. The order is whatever Go's map iteration
// yields: the framework stores these factories in maps keyed by type name, so
// registration order has no effect on the provider.
func factoriesOf[T any](factoriesMap map[string]func() T) []func() T {
	factories := make([]func() T, 0, len(factoriesMap))
	for _, factory := range factoriesMap {
		factories = append(factories, factory)
	}
	return factories
}

// frameworkResources returns every framework Resource factory.
func frameworkResources() []func() resource.Resource {
	return factoriesOf(resourcesMap)
}

// frameworkDataSources returns every framework DataSource factory.
func frameworkDataSources() []func() datasource.DataSource {
	return factoriesOf(dataSourcesMap)
}

// frameworkFunctions returns every framework Function factory.
func frameworkFunctions() []func() function.Function {
	return factoriesOf(functionsMap)
}

// frameworkEphemeralResources returns every framework EphemeralResource factory.
func frameworkEphemeralResources() []func() ephemeral.EphemeralResource {
	return factoriesOf(ephemeralResourcesMap)
}

// frameworkListResources returns every framework ListResource factory.
func frameworkListResources() []func() list.ListResource {
	return factoriesOf(listResourcesMap)
}

// frameworkActions returns every framework Action factory.
func frameworkActions() []func() action.Action {
	return factoriesOf(actionsMap)
}
