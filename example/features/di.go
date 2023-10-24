package features

import (
	"go-app/appbase/di"
	"go-app/example/features/orders"

	"github.com/rs/zerolog/log"
)

func GetOrderApi() orders.OrderApi {
	var item orders.OrderApi
	err := di.Container().Invoke(func(e orders.OrderApi) {
		item = e
	})
	if err != nil {
		log.Fatal().Msg("failed to resolve OrderApi")
	}
	return item
}
