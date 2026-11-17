package fw

import (
	// Blank-import every schema/*planmodifier and schema/*default subpackage
	// so `go mod vendor` keeps them in the vendor directory even when no
	// resource references them directly yet. These leaf subpackages provide
	// the built-in plan modifiers (e.g. stringplanmodifier.RequiresReplace(),
	// stringplanmodifier.UseStateForUnknown()) and attribute default
	// constructors (e.g. stringdefault.StaticString()) that new Framework
	// resources commonly need. The framework itself never imports them in
	// production code, so without this file `go mod vendor` prunes any
	// subpackage that is not imported by the main module, and generated code
	// that references them fails to compile with:
	//
	//   cannot find module providing package
	//   .../resource/schema/stringplanmodifier: import lookup disabled by
	//   -mod=vendor
	//
	// This is a monorepo-wide dependency pin: keep this import list in sync
	// with the subpackages actually shipped by terraform-plugin-framework.
	//
	//nolint:revive,nolintlint // blank imports are the whole point of this file
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicdefault"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicplanmodifier"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/float32default"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/float32planmodifier"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberdefault"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	_ "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)
