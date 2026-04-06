package main

import (
    "fmt"
    "time"
    "github.com/jhionan/multichain-staking/internal/auth"
)

func main() {
    svc := auth.NewJWTService("this-is-a-32-character-secret-key-for-testing!!")
    token, err := svc.Sign(auth.Claims{Wallet: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", Role: auth.RoleAdmin}, 24*time.Hour)
    if err != nil {
        panic(err)
    }
    fmt.Print(token)
}
