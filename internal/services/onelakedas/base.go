// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package onelakedas

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "OneLake Data Access Security",
	Type:           "onelake_data_access_security",
	Names:          "OneLake Data Access Securities",
	Types:          "onelake_data_access_securities",
	DocsURL:        "https://learn.microsoft.com/rest/api/fabric/core/onelake-data-access-security",
	IsPreview:      true,
	IsSPNSupported: true,
}
