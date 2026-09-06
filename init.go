package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	pbfriends "github.com/PretendoNetwork/grpc/go/friends"
	"github.com/PretendoNetwork/nex-go/v2"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/Protarium-Network/splatoon-testfire-nex/globals"
)

func requiredEnv(name string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		globals.Logger.Errorf("%s environment variable not set", name)
		os.Exit(1)
	}

	return value
}

func requiredPort(name string) int {
	value := requiredEnv(name)
	port, err := strconv.Atoi(value)
	if err != nil || port < 1 || port > 65535 {
		globals.Logger.Errorf("%s is not a valid UDP port: %s", name, value)
		os.Exit(1)
	}

	return port
}

func configureRemoteAccountServices() {
	accountHost := requiredEnv("PN_GLOBAL_TESTFIRE_ACCOUNT_GRPC_HOST")
	accountPort := requiredPort("PN_GLOBAL_TESTFIRE_ACCOUNT_GRPC_PORT")
	accountAPIKey := os.Getenv("PN_GLOBAL_TESTFIRE_ACCOUNT_GRPC_API_KEY")
	if strings.TrimSpace(accountAPIKey) == "" {
		globals.Logger.Warning("Account gRPC API key is empty")
	}

	common_globals.ConnectToAccountGRPC(accountHost, uint16(accountPort), accountAPIKey)

	friendsHost := requiredEnv("PN_GLOBAL_TESTFIRE_FRIENDS_GRPC_HOST")
	friendsPort := requiredPort("PN_GLOBAL_TESTFIRE_FRIENDS_GRPC_PORT")
	friendsAPIKey := os.Getenv("PN_GLOBAL_TESTFIRE_FRIENDS_GRPC_API_KEY")
	if strings.TrimSpace(friendsAPIKey) == "" {
		globals.Logger.Warning("Friends gRPC API key is empty")
	}

	var err error
	globals.GRPCFriendsClientConnection, err = grpc.NewClient(
		fmt.Sprintf("dns:%s:%d", friendsHost, friendsPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		globals.Logger.Criticalf("Failed to connect to friends gRPC server: %v", err)
		os.Exit(1)
	}

	globals.GRPCFriendsClient = pbfriends.NewFriendsClient(globals.GRPCFriendsClientConnection)
	globals.GRPCFriendsCommonMetadata = metadata.Pairs("X-API-Key", friendsAPIKey)
}

func init() {
	globals.Logger = plogger.NewLogger()

	if err := godotenv.Load(); err != nil {
		globals.Logger.Warning("Error loading .env file")
	}

	// Validate listener configuration early, before opening background services.
	requiredPort("PN_GLOBAL_TESTFIRE_AUTHENTICATION_SERVER_PORT")
	requiredEnv("PN_GLOBAL_TESTFIRE_SECURE_SERVER_HOST")
	requiredPort("PN_GLOBAL_TESTFIRE_SECURE_SERVER_PORT")

	kerberosPassword := make([]byte, 0x10)
	if _, err := rand.Read(kerberosPassword); err != nil {
		globals.Logger.Criticalf("Error generating Kerberos password: %v", err)
		os.Exit(1)
	}

	globals.KerberosPassword = string(kerberosPassword)
	globals.InitAccounts()

	globals.LocalAuthMode = os.Getenv("PN_GLOBAL_TESTFIRE_LOCAL_MODE") == "1"
	if globals.LocalAuthMode {
		globals.Logger.Warning("Local Testfire mode is enabled; use it only in an isolated preservation environment")
	} else {
		configureRemoteAccountServices()
	}

	postgresURI := requiredEnv("PN_GLOBAL_TESTFIRE_POSTGRES_URI")
	var err error
	globals.Postgres, err = sql.Open("postgres", postgresURI)
	if err != nil {
		globals.Logger.Criticalf("Failed to configure PostgreSQL: %v", err)
		os.Exit(1)
	}
	if err := globals.Postgres.Ping(); err != nil {
		globals.Logger.Criticalf("Failed to reach PostgreSQL: %v", err)
		os.Exit(1)
	}
	globals.Logger.Success("Configured PostgreSQL matchmaking storage")

	healthCheckPort := strings.TrimSpace(os.Getenv("PN_GLOBAL_TESTFIRE_HEALTH_CHECK_PORT"))
	if healthCheckPort == "" {
		return
	}

	port, err := strconv.Atoi(healthCheckPort)
	if err != nil || port < 1 || port > 65535 {
		globals.Logger.Errorf("PN_GLOBAL_TESTFIRE_HEALTH_CHECK_PORT is not a valid UDP port: %s", healthCheckPort)
		os.Exit(1)
	}

	nex.EnableBasicUDPHealthCheck(port)
}
