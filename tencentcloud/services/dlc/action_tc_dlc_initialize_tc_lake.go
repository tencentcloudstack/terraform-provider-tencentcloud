package dlc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// _ action.Action is intentionally omitted: the factory function
// NewDlcInitializeTCLake returns action.Action, which already performs the
// implicit compile-time check at the return statement. Only the
// extended-interface assertion is kept here to verify the Configure method
// signature (promoted from fw.ActionWithConfigure) at compile time.
var _ action.ActionWithConfigure = &DlcInitializeTCLake{}

// NewDlcInitializeTCLake is the factory referenced by
// tencentcloud/framework/registry.go to register this action.
func NewDlcInitializeTCLake() action.Action {
	return &DlcInitializeTCLake{}
}

// DlcInitializeTCLake implements action.Action for
// tencentcloud_dlc_initialize_tc_lake.
type DlcInitializeTCLake struct {
	fw.ActionWithConfigure
}

// DlcInitializeTCLakeModel maps the schema attributes. The InitializeTCLake
// API takes no request parameters, so this model is empty (kept for symmetry
// with the reference action implementations).
type DlcInitializeTCLakeModel struct {
}

func (a *DlcInitializeTCLake) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_dlc_initialize_tc_lake"
}

func (a *DlcInitializeTCLake) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an action to activate (open) the DLC (Data Lake Compute) TCLake service " +
			"via the InitializeTCLake API. This is a one-time, idempotent initialization operation; " +
			"no cloud-side state is persisted after the action completes.",
		Attributes: map[string]schema.Attribute{},
	}
}

// Invoke is called to run the logic of the action.
func (a *DlcInitializeTCLake) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data DlcInitializeTCLakeModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if a.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	service := NewDlcService(a.Client())
	response, err := service.InitializeTCLake(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("initializing DLC TCLake"),
			err.Error(),
		)
		return
	}

	// Framework actions cannot expose output attributes, so the API response
	// fields are logged for observability only.
	if response != nil && response.Response != nil {
		instanceId := ""
		if response.Response.InstanceId != nil {
			instanceId = *response.Response.InstanceId
		}
		isSuccess := false
		if response.Response.IsSuccess != nil {
			isSuccess = *response.Response.IsSuccess
		}
		log.Printf("[DEBUG] dlc initialize_tc_lake invoked, instance_id=%s, is_success=%v", instanceId, isSuccess)
	}
}
