package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"arealis-drunix-hackathon/drunix"
	"arealis-drunix-hackathon/handler"
	"arealis-drunix-hackathon/repository"
	"arealis-drunix-hackathon/services/correlation"
	"arealis-drunix-hackathon/services/orchestration"
	"arealis-drunix-hackathon/services/sequencer"
)

func main() {
	fmt.Println("Starting Ergos demo...")

	// ------------------------------------------------------------
	// 1. Connect to Drunix
	// ------------------------------------------------------------

	fmt.Println("\nConnecting to Drunix...")

	client, err := drunix.NewClient(drunix.Config{
		MSPID:         "Org1MSP",
		CryptoPath:    "../drunix/drunix-network/test-network/organizations/peerOrganizations/org1.example.com",
		ChannelName:   "mychannel",
		ChaincodeName: "rwa",
		PeerEndpoint:  "dns:///localhost:7051",
		GatewayPeer:   "peer0.org1.example.com",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Println("Connected to Drunix successfully!")

	fmt.Println("\nCreating RWA asset...")

	if err := client.CreateRWAAsset(
		"SOLAR-001",
		"issuer-001",
		100000,
	); err != nil {
		log.Fatal(err)
	}

	fmt.Println("RWA asset created!")

	// ------------------------------------------------------------
	// 2. Create Ergos components
	// ------------------------------------------------------------

	store := repository.NewCorrelationStore()
	barrier := correlation.NewBarrier(store)
	orchestrator := orchestration.NewService(barrier)
	seq := sequencer.NewSequencer()

	// ------------------------------------------------------------
	// 3. Read agent events
	// ------------------------------------------------------------

	raw, err := os.ReadFile("examples/solar-asset/events.json")
	if err != nil {
		log.Fatal(err)
	}

	var rawEvents []json.RawMessage

	if err := json.Unmarshal(raw, &rawEvents); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nReceived %d agent events\n", len(rawEvents))

	// ------------------------------------------------------------
	// 4. Process agent events
	// ------------------------------------------------------------

	for _, rawEvent := range rawEvents {
		event, err := handler.ParseEvent(rawEvent)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf(
			"Agent event: %-30s Asset: %s\n",
			event.Type,
			event.AssetID,
		)

		intent, ready, err := orchestrator.ProcessEvent(event)
		if err != nil {
			log.Fatal(err)
		}

		if !ready {
			fmt.Println("  → Correlation barrier not ready")
			continue
		}

		fmt.Printf(
			"  → Barrier satisfied for transaction %s\n",
			intent.CorrelationID,
		)

		seq.Enqueue(*intent)

		fmt.Printf(
			"  → Queued intent: %s %s\n",
			intent.Action,
			intent.AssetID,
		)
	}

	// ------------------------------------------------------------
	// 5. Get next transaction from sequencer
	// ------------------------------------------------------------

	intent, ok := seq.Next("SOLAR-001")
	if !ok {
		log.Fatal("no transaction intent available")
	}

	fmt.Printf(
		"\nSequencer produced intent:\n"+
			"  Correlation ID: %s\n"+
			"  Asset ID:       %s\n"+
			"  Action:         %s\n",
		intent.CorrelationID,
		intent.AssetID,
		intent.Action,
	)

	// ------------------------------------------------------------
	// 6. Execute against Drunix
	// ------------------------------------------------------------

	fmt.Println("\nSubmitting transaction to Drunix...")

	if err := client.LockRWAAsset(intent.AssetID); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transaction committed successfully!")

	// ------------------------------------------------------------
	// 7. Read authoritative state from Drunix
	// ------------------------------------------------------------

	asset, err := client.GetRWAAsset(intent.AssetID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nFinal Drunix ledger state:\n%s\n", asset)

	fmt.Println("\nErgos demo completed successfully!")
}
