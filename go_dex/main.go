package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"
)

type AMMPool struct {
    TokenAReserve float64
    TokenBReserve float64
    FeePercent    float64
    mutex         sync.Mutex
}

func (p *AMMPool) Swap(fromToken string, amount float64) (float64, error) {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    if amount <= 0 {
        return 0, fmt.Errorf("amount must be > 0")
    }

    dxWithFee := amount * (1 - p.FeePercent)
    var dy float64

    switch fromToken {
    case "A":
        newX := p.TokenAReserve + dxWithFee
        dy = p.TokenBReserve - ((p.TokenAReserve * p.TokenBReserve) / newX)
        if dy <= 0 {
            return 0, fmt.Errorf("invalid swap")
        }
        p.TokenAReserve += amount
        p.TokenBReserve -= dy
    case "B":
        newY := p.TokenBReserve + dxWithFee
        dy = p.TokenAReserve - ((p.TokenAReserve * p.TokenBReserve) / newY)
        if dy <= 0 {
            return 0, fmt.Errorf("invalid swap")
        }
        p.TokenBReserve += amount
        p.TokenAReserve -= dy
    default:
        return 0, fmt.Errorf("invalid token")
    }

    return dy, nil
}


type SwapRequest struct {
    FromToken string  `json:"from_token"`
    Amount    float64 `json:"amount"`
}

type SwapResponse struct {
    Received float64 `json:"received"`
}

func swapHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Use POST", http.StatusMethodNotAllowed)
        return
    }

    var req SwapRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    received, err := pool.Swap(req.FromToken, req.Amount)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    json.NewEncoder(w).Encode(SwapResponse{Received: received})
}

func poolStatusHandler(w http.ResponseWriter, r *http.Request) {
    pool.mutex.Lock()
    defer pool.mutex.Unlock()

    json.NewEncoder(w).Encode(map[string]float64{
        "TokenAReserve": pool.TokenAReserve,
        "TokenBReserve": pool.TokenBReserve,
    })
}

var pool = AMMPool{
    TokenAReserve: 1000,
    TokenBReserve: 1000,
    FeePercent:    0.003,
}

func main() {
    http.HandleFunc("/swap", swapHandler)
    http.HandleFunc("/pool_status", poolStatusHandler)

    fmt.Println("AMM running at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
