package main

import (
	"crypto/x509"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
    mspID        = "Org1MSP"
    cryptoPath   = "../drunix/drunix-network/test-network/organizations/peerOrganizations/org1.example.com"
    certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts"
    keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore"
    tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
    peerEndpoint = "dns:///localhost:7051"
    gatewayPeer  = "peer0.org1.example.com"
)

func main() {
	fmt.Println("Connecting to Drunix...")

	connection := newGrpcConnection()
	defer connection.Close()

	id := newIdentity()
	sign := newSign()

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(connection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(fmt.Errorf("failed to connect to Drunix: %w", err))
	}
	defer gw.Close()

	network := gw.GetNetwork("mychannel")
	contract := network.GetContract("basic")

	result, err := contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		panic(fmt.Errorf("failed to query GetAllAssets: %w", err))
	}

	fmt.Println("Connected to Drunix successfully!")
	fmt.Println("GetAllAssets:")
	fmt.Println(string(result))
}

func newGrpcConnection() *grpc.ClientConn {
	certificatePEM, err := os.ReadFile(tlsCertPath)
	if err != nil {
		panic(fmt.Errorf("failed to read TLS certificate: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)

	transportCredentials := credentials.NewClientTLSFromCert(
		certPool,
		gatewayPeer,
	)

	connection, err := grpc.NewClient(
		peerEndpoint,
		grpc.WithTransportCredentials(transportCredentials),
	)
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

func newIdentity() *identity.X509Identity {
	certificatePEM, err := readFirstFile(certPath)
	if err != nil {
		panic(fmt.Errorf("failed to read identity certificate: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

func newSign() identity.Sign {
	privateKeyPEM, err := readFirstFile(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

func readFirstFile(dirPath string) ([]byte, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	fileNames, err := dir.Readdirnames(1)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(path.Join(dirPath, fileNames[0]))
}
