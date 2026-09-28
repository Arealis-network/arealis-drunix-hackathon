package main

import (
	"fmt"
	"log"

	"arealis-drunix-hackathon/drunix"
)

func main() {
	fmt.Println("Connecting to Drunix...")

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

	err = client.CreateRWAAsset(
		"SOLAR-001",
		"issuer-001",
		100000,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("RWA asset created!")

	asset, err := client.GetRWAAsset("SOLAR-001")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nBefore lock:\n%s\n", asset)

	fmt.Println("\nLocking RWA asset...")

	err = client.LockRWAAsset("SOLAR-001")
	if err != nil {
		log.Fatal(err)
	}

	asset, err = client.GetRWAAsset("SOLAR-001")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nAfter lock:\n%s\n", asset)
}