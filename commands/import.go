package commands

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/api"
	cli "gopkg.in/urfave/cli.v2"
)

var schemes = map[string]map[string]string{
	"1password-de": map[string]string{
		"Titel":                     "name",
		"Typ":                       "type",
		"Username":                  "username",
		"Password":                  "password",
		"Änderungsdatum":            "modifiedAt",
		"Erstellungsdatum":          "createdAt",
		"Lizenziert für(reg_name)":  "regName",
		"Lizenzschlüssel(reg_code)": "regCode",
		"E-Mail(email)":             "email",
		"Organisation(org_name)":    "orgName",
		"SID(sid)":                  "sid",
		"Nachname(lastname)":        "lastname",
		"Vorname(firstname)":        "firstname",
		"URL":                       "urls",
	},
}

var typeSchemes = map[string]map[string]api.SecretType{
	"1password-de": map[string]api.SecretType{
		"Login":          api.SecretTypeLogin,
		"Softwarelizenz": api.SecretTypeLicence,
		"Sichere Notiz":  api.SecretTypeNote,
		"Password":       api.SecretTypePassword,
		"WLAN-Router":    api.SecretTypeWLAN,
	},
}

// ImportFlags holds all the values of the commandline flags relevant to the
// ImportCommand
var ImportFlags = &struct {
	Scheme string
}{}

// ImportCommand imports password from another password safe exported as csv
var ImportCommand = &cli.Command{
	Name:      "import",
	Usage:     "Import password from a csv file",
	ArgsUsage: "<csv-file>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "scheme",
			Value:       "",
			Usage:       "The output scheme (support: 1password-de)",
			Destination: &ImportFlags.Scheme,
		},
	},
	Action: withDetailedErrors(importFile),
}

func importFile(ctx *cli.Context) error {
	if ctx.Args().Len() != 1 {
		return errors.New("No file to import specificied")
	}
	file, err := os.Open(ctx.Args().First())
	if err != nil {
		return errors.Wrapf(err, "Unable to open: %s", ctx.Args().First())
	}
	defer file.Close()

	scheme := schemes[ImportFlags.Scheme]
	typeScheme := typeSchemes[ImportFlags.Scheme]

	lines := csv.NewReader(file)
	lines.LazyQuotes = true
	headers, err := lines.Read()
	if err != nil {
		return errors.Wrap(err, "Failed to read header")
	}
	fieldNames := make([]string, len(headers))
	for i, header := range headers {
		fieldName := header
		if translated, ok := scheme[fieldName]; ok {
			fieldName = translated
		}
		fieldNames[i] = fieldName
	}
	fmt.Println("The followind fields have been identified:")
	for _, fieldName := range fieldNames {
		fmt.Println(fieldName)
	}
	if !confirm("Continue import?") {
		return nil
	}

	logger := createLogger()
	client := createRemote(logger)

	if _, err := unlockStore(client); err != nil {
		return err
	}

	for {
		record, err := lines.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.Wrap(err, "CSV parse failed")
		}
		properties := map[string]string{}
		for i, col := range record {
			if col == "" {
				continue
			}
			properties[fieldNames[i]] = col
		}
		name, nameOk := properties["name"]
		secretType, typeOk := properties["type"]
		if !nameOk || !typeOk {
			fmt.Printf("Missing name and type information in: %v\n", properties)
			continue
		}
		delete(properties, "name")
		delete(properties, "type")
		delete(properties, "Tags")
		if translated, ok := typeScheme[secretType]; ok {
			secretType = string(translated)
		}
		id, err := importGenerateID(properties)
		if err != nil {
			fmt.Println(err)
			continue
		}

		timestamp := importExtractTimestamp(properties)
		delete(properties, "createdAt")
		delete(properties, "modifiedAt")

		urls := strings.Split(properties["urls"], ",")
		delete(properties, "urls")

		if err := client.Add(createClientContext(), id, api.SecretType(secretType), api.SecretVersion{
			Timestamp:  timestamp,
			Name:       name,
			URLs:       urls,
			Properties: properties,
		}); err != nil {
			fmt.Printf("Unable to import: %s\n", name)
			fmt.Println(err)
		} else {
			fmt.Printf("Imported: %s\n", name)
		}
	}

	return nil
}

func importGenerateID(properties map[string]string) (string, error) {
	hash := sha256.New()

	for name, value := range properties {
		if _, err := hash.Write([]byte(name)); err != nil {
			return "", errors.Wrap(err, "Failed to generate id")
		}
		if _, err := hash.Write([]byte(value)); err != nil {
			return "", errors.Wrap(err, "Failed to generate id")
		}
	}

	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil)), nil
}

func importExtractTimestamp(properties map[string]string) time.Time {
	if modifiedAt, ok := properties["modifiedAt"]; ok {
		unix, err := strconv.ParseInt(modifiedAt, 10, 64)
		if err == nil {
			return time.Unix(unix, 0)
		}
	}
	if createdAt, ok := properties["createdAt"]; ok {
		unix, err := strconv.ParseInt(createdAt, 10, 64)
		if err == nil {
			return time.Unix(unix, 0)
		}
	}
	return time.Now()
}
