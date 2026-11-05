# AGS local verification and SDK handoff

The Provider manages Tools. Runtime instances belong to a separate SDK caller.
These helpers keep candidate installation and example validation independent
of credentials, cloud provisioning and a Registry release.

From the repository root:

```sh
source examples/ags/e2e/env.sh
make build
make check-mux
go test -mod=vendor -count=1 ./tencentcloud/services/ags ./tencentcloud/framework ./tencentcloud/internal/sharedmeta

python3 examples/ags/e2e/validate_examples.py \
  --binary "$GOBIN/terraform-provider-tencentcloud" \
  --terraform /absolute/path/to/terraform \
  --work-root /absolute/path/outside/the/repository/ags-example-validation
```

Use Terraform1.17.0 for the local candidate examples. Validation installs an
unpublished `0.0.1-local` mirror, checks five examples and the actual Protocol5
mux schema, and verifies a create-only plan. It never applies that plan. Keep
state, saved plans, credentials and observations outside Git.

The Go tests exercise the public SDK wire format, Framework state retention,
import, cancellation/deletion boundaries, pagination, semantic equivalence and
provider alias isolation without cloud requests. Complete successful exact-ID
empty reads remove established managed state; explicit errors and incomplete
responses retain it. Successful permission-filtered empties are indistinguishable
from external deletion, as with other Provider resources.

## Independent runtime smoke

The `runtime` directory is a separate Go module using public `ags-go-sdk v0.1.5`.
Build it independently:

```sh
cd examples/ags/e2e/runtime
go build -mod=readonly -o /absolute/path/outside/the/repository/ags-runtime-smoke .
```

The helper requires `--cloud`, a private ownership ledger and an observation
path. The ledger binds `run_id`, `region`, `tool_id` and the exact
`tf-ags-<run_id>-...` Tool name. Credentials use the SDK's standard environment
variables, including a session token when applicable. The `read` mode performs
an independent Tool query; `code` starts a dedicated instance and requires actual
remote stdout42. Cleanup modes operate only on ledger-owned identities and are
separate from Provider behavior. Review the ledger and ensure an actual compatible
code Tool before opting in; the Terraform resource itself never starts or stops
instances.

Historical development cloud results were obtained on a different unpublished
candidate. This PR does not claim complete V1 cloud acceptance on its rebased
binary. Remaining environment/compatibility limits are documented in the PR,
including Description clearing, VPC quota and the service's current custom-only
CLS support. Real cloud scenarios require compatible images and existing roles,
networking and storage references; ordinary tests do not provision them.
