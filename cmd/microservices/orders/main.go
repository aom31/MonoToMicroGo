package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Starting order microservice")

	ctx := cmd.Context()

	r, closeFn := createOrderMicroservices()

	defer closeFn()

	server := &http.Server{Addr: os.Getenv("SHOP_ORDER_SERVICE_BIND_ADDR"), Handler: r}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	<-ctx.Done()
	log.Println("closing order microservice")

	if err := server.Close(); err != nil {
		panic(err)
	}

}

func createOrderMicroservices() (router *chi.Mux, closeFn func()) {
	cmd.WaitForService(os.Getenv("SHOP_RABBITMQ_ADDR"))

	shopHTTPClient := orders.infra_product.NewHTTPClient(os.Getenv("SHOP_SERVICE_ADDR"))

	r := cmd.CreateRouter()

	orders_public_http.AddRoutes(r, ordersService, ordersRepo)
	orders_private_http.AddRoutes(r, ordersService, ordersRepo)

	return r, func() {

	}
}
