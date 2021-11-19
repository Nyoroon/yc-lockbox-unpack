# yc-lockbox-struct

Library to populate struct fields from Yandex.Cloud Lockbox Secret Payload.

### How to use

1. Define struct:
    ```go
    type testSecret struct {
        Text   string `secretKey:"testText"`
        Binary []byte `secretKey:"testFile"`
    }
    ```
2. Request secret payload and extract secrets to variable:
    ```go
    payload, err := sdk.LockboxPayload().Payload().Get(ctx, &lockbox.GetPayloadRequest{SecretId: "mysecretid"})
    if err != nil {
        log.Fatalf("Unable to get secret: %v", err)
    }

    secret := testSecret{}
    if err := yc_lockbox_unpack.UnpackPayload(payload, &secret); err != nil {
        log.Fatalf("Unable to unpack secret: %v", err)
    }
    ```

#### See [example](example/)
