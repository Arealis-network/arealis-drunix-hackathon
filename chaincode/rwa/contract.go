package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) CreateAsset(
	ctx contractapi.TransactionContextInterface,
	id string,
	owner string,
	totalTokens int64,
) error {
	existing, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read asset: %w", err)
	}

	if existing != nil {
		return fmt.Errorf("asset %s already exists", id)
	}

	if totalTokens <= 0 {
		return fmt.Errorf("total tokens must be greater than zero")
	}

	asset := Asset{
		ID:               id,
		Owner:            owner,
		TotalTokens:      totalTokens,
		AvailableTokens:  totalTokens,
		Status:           "AVAILABLE",
		SettlementStatus: "PENDING",
	}

	data, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to marshal asset: %w", err)
	}

	if err := ctx.GetStub().PutState(id, data); err != nil {
		return fmt.Errorf("failed to store asset: %w", err)
	}

	return nil
}

func (s *SmartContract) GetAsset(
	ctx contractapi.TransactionContextInterface,
	id string,
) (*Asset, error) {
	data, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read asset: %w", err)
	}

	if data == nil {
		return nil, fmt.Errorf("asset %s does not exist", id)
	}

	var asset Asset
	if err := json.Unmarshal(data, &asset); err != nil {
		return nil, fmt.Errorf("failed to unmarshal asset: %w", err)
	}

	return &asset, nil
}

func (s *SmartContract) LockAsset(
	ctx contractapi.TransactionContextInterface,
	id string,
) error {
	asset, err := s.GetAsset(ctx, id)
	if err != nil {
		return err
	}

	if asset.Status != "AVAILABLE" {
		return fmt.Errorf(
			"asset %s cannot be locked from status %s",
			id,
			asset.Status,
		)
	}

	asset.Status = "LOCKED"

	data, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to marshal asset: %w", err)
	}

	if err := ctx.GetStub().PutState(id, data); err != nil {
		return fmt.Errorf("failed to update asset: %w", err)
	}

	return nil
}
