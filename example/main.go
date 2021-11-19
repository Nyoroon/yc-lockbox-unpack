package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	yc_lockbox_unpack "github.com/Nyoroon/yc-lockbox-unpack"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/lockbox/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
)

type testSecret struct {
	Text   string `secretKey:"testText" json:"test_text"`
	Binary []byte `secretKey:"testFile" json:"test_file"`
}

var flagSecretID = flag.String("id", "", "YC Lockbox Secret ID")

func main() {
	flag.Parse()

	if *flagSecretID == "" {
		flag.Usage()
		os.Exit(1)
	}

	ctx := context.Background()

	log.Println("Creating SDK")
	sdk, err := ycsdk.Build(
		ctx,
		ycsdk.Config{
			Credentials: ycsdk.NewIAMTokenCredentials(os.Getenv("YC_TOKEN")),
		},
	)
	if err != nil {
		log.Fatalf("Unable to create sdk: %v", err)
	}

	log.Println("Requesting Secret Payload")
	payload, err := sdk.LockboxPayload().Payload().Get(ctx, &lockbox.GetPayloadRequest{SecretId: *flagSecretID})
	if err != nil {
		log.Fatalf("Unable to get secret: %v", err)
	}

	secret := testSecret{}
	if err := yc_lockbox_unpack.UnpackPayload(payload, &secret); err != nil {
		log.Fatalf("Unable to unpack secret: %v", err)
	}

	data, err := json.MarshalIndent(secret, "", "  ")
	if err != nil {
		log.Fatalf("Unable to marshal secret: %v", err)
	}

	log.Println(string(data))
}
