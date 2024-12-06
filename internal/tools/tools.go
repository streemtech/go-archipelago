//go:build tools

package tools

import (
	//TODO2 remove codegen import after tool directive is added to mod.
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
	_ "github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)
