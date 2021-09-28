package helpers

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
)

func IndonesiaCurrency(value float64) string {
	return fmt.Sprintf("Rp %s", strings.ReplaceAll(humanize.Commaf(value), ",", "."))
}
