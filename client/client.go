package main

import (
	"context"
	"log"
	pb "microstoria/proto"
	"time"

	"github.com/alexflint/go-arg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func logResponses(res *pb.EmailResponse, err error) {
	if err != nil {
		log.Fatalf("gRPC error: %v\n", err)
	}

	if res.EmailEntry == nil {
		log.Printf("email not found: %v\n", res)
	} else {
		log.Printf("gRPC response: %v\n", res.EmailEntry)
	}
}

func createEmail(client pb.MicrostoriaServiceClient, addr string) *pb.EmailEntry {
	log.Println("create email")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.CreateEmail(ctx, &pb.CreateEmailRequest{EmailAddr: addr})
	logResponses(res, err)

	return res.EmailEntry
}

func getEmail(client pb.MicrostoriaServiceClient, addr string) *pb.EmailEntry {
	log.Println("get email")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetEmail(ctx, &pb.GetEmailRequest{EmailAddr: addr})
	logResponses(res, err)

	return res.EmailEntry
}

func getEmailBatch(client pb.MicrostoriaServiceClient, count int, page int) {
	log.Println("get email")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetEmailBatch(ctx, &pb.GetEmailBatchRequest{Count: int32(count), Page: int32(page)})
	if err != nil {
		log.Fatalf("gRPC error: %v\n", err)
	}
	log.Println("response:")
	for i := 0; i < len(res.EmailEntries); i++ {
		log.Printf(" item [%v of %v]: %s", i+1, len(res.EmailEntries), res.EmailEntries[i])
	}
}

func updateEmail(client pb.MicrostoriaServiceClient, entry pb.EmailEntry) *pb.EmailEntry {
	log.Println("update email")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.UpdateEmail(ctx, &pb.UpdateEmailRequest{EmailEntry: &entry})
	logResponses(res, err)

	return res.EmailEntry
}

func deleteEmail(client pb.MicrostoriaServiceClient, addr string) *pb.EmailEntry {
	log.Println("delete email")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.DeleteEmail(ctx, &pb.DeleteEmailRequest{EmailAddr: addr})
	logResponses(res, err)

	return res.EmailEntry
}

var args struct {
	GrpcAddr string `arg:"env:MICROSTORIA_GRPC_ADDR"`
}

func main() {
	arg.MustParse(&args)
	if args.GrpcAddr == "" {
		args.GrpcAddr = ":8081"
	}

	conn, err := grpc.Dial(args.GrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to gRPC server: %v\n", err)
	}
	defer conn.Close()
	client := pb.NewMicrostoriaServiceClient(conn)

	newEmail := createEmail(client, "999@999.999")
	newEmail.ConfirmedAt = 10000
	updateEmail(client, *newEmail)
	deleteEmail(client, newEmail.Email)

	getEmailBatch(client, 3, 1)
	getEmailBatch(client, 3, 2)
	getEmailBatch(client, 3, 3)
}
