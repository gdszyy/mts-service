package main

import (
"fmt"
"log"
"github.com/gdsZyy/mts-service/internal/client"
)

func main() {
token := "YBb5b3gRMOGJrbTzpF"
bookmakerID, err := client.FetchBookmakerID(token, false)
if err != nil {
log.Fatalf("Error: %v", err)
}
fmt.Printf("Success! Bookmaker ID: %s\n", bookmakerID)
}
