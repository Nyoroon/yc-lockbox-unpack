# Example

How to run:
```shell
$ YC_TOKEN=$(yc iam create-token) go run main.go -id secretid
2021/11/20 01:25:09 Creating SDK
2021/11/20 01:25:09 Requesting Secret Payload
2021/11/20 01:25:09 {
  "test_text": "my text",
  "test_file": "iVBORw0KGgoAAAANSUhEUgAAAV0AAAA9CAIAAACBaS0OAAAgAElEQVR4nIy9B1ccV7r1P9..."
}
```
