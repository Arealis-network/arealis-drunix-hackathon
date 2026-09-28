package main

type Asset struct {
    ID               string `json:"id"`
    Owner            string `json:"owner"`
    TotalTokens      int64  `json:"total_tokens"`
    AvailableTokens  int64  `json:"available_tokens"`
    Status           string `json:"status"`
    SettlementStatus string `json:"settlement_status"`
}