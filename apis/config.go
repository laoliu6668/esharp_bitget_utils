package apis

import (
	"fmt"

	root "github.com/laoliu6668/esharp_bitget_utils"
	"github.com/laoliu6668/esharp_bitget_utils/util"
)

const Gateway = "api.bitget.com"

func GetFlag() string {
	return fmt.Sprintf("%s_%s", root.ExchangeName, util.GetFuncName(1))
}
