package cmd

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/thetredev/steamcmd-cli/server"
	"github.com/thetredev/steamcmd-cli/shared"
)

var certsCmd = &cobra.Command{
	Use:     "certs",
	Version: shared.Version,
	Short:   "Cert stuff",
	Long:    `A longer description`,
	Run:     certsCallback,
}

func init() {
	serverCmd.AddCommand(certsCmd)
}

func certsCallback(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func stringToInt(input string, inputType string, maxStringLength int, maxValue int) int {
	if len(input) > maxStringLength {
		log.Fatalf("Invalid length for '%s' value '%s' (max length: %d).", inputType, input, maxStringLength)
	}

	value, err := strconv.Atoi(input)

	if err != nil {
		log.Fatalf("Could convert string '%s' to an int.", input)
	}

	if value > maxValue {
		log.Fatalf("Value '%d' is large than the maximum of '%d'.", value, maxValue)
	}

	return value
}

func stringToDate(input string) time.Time {
	parts := strings.Split(input, "-")

	if len(parts) < 3 {
		log.Fatal("Flag 'valid-until' passed in invalid format.")
	}

	month := time.Month(stringToInt(parts[1], "month", 2, 12))
	day := stringToInt(parts[2], "day", 2, 31)

	if month == 2 && day > 28 {
		log.Fatalf("February only has 28 days. Aborting.")
	}

	return time.Date(
		stringToInt(parts[0], "year", 4, 9999), month, day,
		1, 1, 1,
		0,
		time.UTC,
	)
}

func certCmdParseFlags(cmd *cobra.Command) *server.CertificateParameters {
	company, _ := cmd.Flags().GetString("company")
	country, _ := cmd.Flags().GetString("country")
	province, _ := cmd.Flags().GetString("province")
	locality, _ := cmd.Flags().GetString("locality")
	streetAddress, _ := cmd.Flags().GetString("street-address")
	postalCode, _ := cmd.Flags().GetString("postal-code")
	validUntilString, _ := cmd.Flags().GetString("valid-until")
	keyLength, _ := cmd.Flags().GetInt("key-length")
	hostnames, _ := cmd.Flags().GetStringArray("hostname")

	validUntil := stringToDate(validUntilString)

	return &server.CertificateParameters{
		Company:       company,
		Country:       country,
		Province:      province,
		Locality:      locality,
		StreetAddress: streetAddress,
		PostalCode:    postalCode,
		ValidUntil:    validUntil.UTC(),
		KeyLength:     keyLength,
		Hostnames:     hostnames,
	}
}

func certCmdAddFlags(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false

	cmd.PersistentFlags().String(
		"company",
		"steamcmd-cli",
		"certificate company",
	)

	cmd.PersistentFlags().String(
		"country",
		"US",
		"certificate country",
	)

	cmd.PersistentFlags().String(
		"province",
		"datacenter",
		"certificate province",
	)

	cmd.PersistentFlags().String(
		"locality",
		"server1",
		"certificate locality",
	)

	cmd.PersistentFlags().String(
		"street-address",
		"rack1",
		"certificate street address",
	)

	cmd.PersistentFlags().String(
		"postal-code",
		"row1",
		"certificate postal code",
	)

	cmd.PersistentFlags().String(
		"valid-until",
		"2999-12-31",
		"Date until the certificate should be valid (UTC).\nHours, minutes, and seconds will always be '01'.",
	)

	cmd.PersistentFlags().Int(
		"key-length",
		4096,
		"certificate key length",
	)

	cmd.PersistentFlags().StringArray(
		"hostname",
		[]string{},
		"certificate hostname list.\nDNS names and IPv4/IPv6 addresses are valid, IP networks are NOT.",
	)
}
