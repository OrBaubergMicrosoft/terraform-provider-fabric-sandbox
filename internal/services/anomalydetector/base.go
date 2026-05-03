// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package anomalydetector

import (
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

const (
	FabricItemType            = fabcore.ItemTypeAnomalyDetector
	ItemDefinitionEmpty       = `{}`
	ItemDefinitionPathDocsURL = "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/anomalydetector-definition"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Anomaly Detector",
	Type:           "anomaly_detector",
	Names:          "Anomaly Detectors",
	Types:          "anomaly_detectors",
	DocsURL:        "https://learn.microsoft.com/rest/api/fabric/articles/item-management/definitions/anomalydetector-definition",
	IsPreview:      true,
	IsSPNSupported: true,
}

var itemDefinitionFormats = []fabricitem.DefinitionFormat{ //nolint:gochecknoglobals
	{
		Type:  fabricitem.DefinitionFormatDefault,
		API:   "",
		Paths: []string{"Configurations.json"},
	},
}


