package main

import (
	"context"
	"enttest/controllers"
	"enttest/ent"
	"enttest/seeder"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	seedData,
		generateSchema,
		startListening,
		listenPort,
		postgresHost,
		postgresPort,
		postgresUser,
		postgresDBName,
		postgresPassword,
		postgresSSLMode := LoadEnvironmentVariables(".env")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		postgresHost,
		postgresPort,
		postgresUser,
		postgresDBName,
		postgresPassword,
		postgresSSLMode)
	client, err := ent.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()
	ctx := context.Background()

	if generateSchema {
		// Run the auto migration tool.
		if err := client.Schema.Create(ctx); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}
	}

	if seedData {
		if err = SeedData(ctx, client); err != nil {
			log.Fatal(err)
		}
	}

	if startListening {
		controllers.RegisterControllers(ctx, client)
		http.ListenAndServe(fmt.Sprintf(":%s", listenPort), nil)
	}
}

func LoadEnvironmentVariables(envFile string) (bool, bool, bool, string, string, string, string, string, string, string) {
	godotenv.Load(envFile)

	seedData, err := strconv.ParseBool(goDotEnvVariable("SEED_DATA"))
	if err != nil {
		log.Fatal("Please configure .env file with SEED_DATA variable (bool)")
		log.Fatal(err)
		panic(nil)
	}

	generateSchema, err := strconv.ParseBool(goDotEnvVariable("GENERATE_SCHEMA"))
	if err != nil {
		log.Fatal("Please configure .env file with GENERATE_SCHEMA variable (bool)")
		log.Fatal(err)
		panic(nil)
	}

	startListening, err := strconv.ParseBool(goDotEnvVariable("LISTEN_SERVER"))
	if err != nil {
		log.Fatal("Please configure .env file with LISTEN_SERVER variable (bool)")
		log.Fatal(err)
		panic(nil)
	}

	listenPort := goDotEnvVariable("LISTEN_PORT")
	if listenPort == "" {
		log.Fatal("Please configure .env file with LISTEN_PORT variable (string)")
		panic(nil)
	}

	postgresHost := goDotEnvVariable("PGHOST")
	if postgresHost == "" {
		log.Fatal("Please configure .env file with PGHOST variable (string)")
		panic(nil)
	}

	postgresPort := goDotEnvVariable("PGPORT")
	if postgresPort == "" {
		log.Fatal("Please configure .env file with PGPORT variable (string)")
		panic(nil)
	}

	postgresUser := goDotEnvVariable("PGUSER")
	if postgresUser == "" {
		log.Fatal("Please configure .env file with PGUSER variable (string)")
		panic(nil)
	}

	postgresDBName := goDotEnvVariable("PGDBNAME")
	if postgresDBName == "" {
		log.Fatal("Please configure .env file with PGDBNAME variable (string)")
		panic(nil)
	}

	postgresPassword := goDotEnvVariable("PGPASSWORD")
	if postgresPassword == "" {
		log.Fatal("Please configure .env file with PGPASSWORD variable (string)")
		panic(nil)
	}

	postgresSSLMode := goDotEnvVariable("PGSSLMODE")
	if postgresSSLMode == "" {
		log.Fatal("Please configure .env file with PGSSLMODE variable (string)")
		panic(nil)
	}

	return seedData, generateSchema, startListening, listenPort, postgresHost, postgresPort, postgresUser, postgresDBName, postgresPassword, postgresSSLMode
}

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func SeedData(ctx context.Context,
	client *ent.Client) error {
	data := seeder.GetSeedData(ctx, client)

	bulk := make([]*ent.AircraftCreate, len(data))
	for i, aircraftInfo := range data {
		bulk[i] = aircraftInfo
	}

	_, err := client.Aircraft.CreateBulk(bulk...).Save(ctx)

	return err
}
