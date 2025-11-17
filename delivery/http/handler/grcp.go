package handler

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"router-template/delivery"
	"time"

	nt "router-template/protos/github.com/didirdt/router-template/protos/notes"
	pb "router-template/protos/github.com/didirdt/router-template/protos/rpc"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetProducts(ctx *gin.Context) {

	conn, err := grpc.NewClient("0.0.0.0:10541", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductsClient(conn)

	// Contact the server and print out its response.
	nctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	nctx, cancel = context.WithCancel(nctx)
	defer cancel()

	// r, err := c.GetProduct(nctx, &pb.GetProductReq{Id: 3})
	// if err != nil {
	// 	delivery.PrintError(err.Error())
	// 	ctx.String(http.StatusInternalServerError, "get Product error")
	// }

	// receive stream from rpc
	list := make([]*pb.Product, 0)
	stream, err := c.GetProducts(nctx, &pb.GetProductsReq{Limit: 40})
	if err != nil {
		delivery.PrintError(err.Error())
		ctx.String(http.StatusInternalServerError, "get list Products error")
		return
	}

	for {
		prod, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			delivery.PrintError(err.Error())
			ctx.String(http.StatusInternalServerError, "get list Products error")
			return
		}
		list = append(list, prod)
	}

	if err != nil {
		delivery.PrintError(err.Error())
		ctx.String(http.StatusInternalServerError, "get list Products error")
	} else {
		ctx.JSON(http.StatusOK, gin.H{"list": list})
	}
}

func GetNotes(ctx *gin.Context) {

	conn, err := grpc.NewClient("0.0.0.0:10542", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := nt.NewNotesClient(conn)

	nctx, _ := context.WithTimeout(context.Background(), 40*time.Second)

	r, err := c.GetNote(nctx, &nt.GetNoteReq{Id: 3})
	fmt.Print("note", r.GetContent())

	if err != nil {
		delivery.PrintError(err.Error())
		ctx.String(http.StatusInternalServerError, "get Note error")
	} else {
		ctx.JSON(http.StatusOK, r)
	}
}
