package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aburizalpurnama/grpc-server/proto"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	e := echo.New()
	e.GET("/accounts", func(c echo.Context) error {
		grpcServerAddr := fmt.Sprintf("%s:%s", os.Getenv("GRPC_SERVER_HOST"), os.Getenv("GRPC_SERVER_PORT"))
		conn, err := grpc.Dial(grpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			e.Logger.Fatalf("can't connect to grpc server: %v", err)
		}
		defer conn.Close()

		client := proto.NewAccountsClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		resp, err := client.SelectAccount(ctx, &proto.SelectAccountRequest{})
		if err != nil {
			e.Logger.Fatalf("could't select accounts: %v", err)
		}

		return c.JSON(http.StatusOK, &resp)
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("APP_PORT"))))
}
