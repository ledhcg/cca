package cli

import (
	"fmt"

	"github.com/ledhcg/cca/internal/i18n"
	"github.com/ledhcg/cca/internal/ui"
)

func printUsage() {
	fmt.Println(i18n.T(i18n.KeyGuideUsage))
}

func printGuide() {
	b, o, c, d := ui.C.Bold, ui.C.Off, ui.C.Cyan, ui.C.Dim
	fmt.Printf(i18n.T(i18n.KeyGuideFull), b, o, c, d)
}
