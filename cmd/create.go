/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"

	pluralize "github.com/gertd/go-pluralize"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new migration file",
	Long:  `Create a new migration file.`,
	Run:   makeMigrationFile,
	Args:  cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var MigrationName string

var MigrationRoot string = "db/migrations"

var timestampFormat = "20060102150405"

func makeMigrationFile(cmd *cobra.Command, args []string) {
	MigrationName = args[0]
	version := time.Now().Format(timestampFormat)
	filename := fmt.Sprintf("%v_create_%v.%v", version, MigrationName, "go")
	fs := afero.NewBasePathFs(afero.NewOsFs(), "db/migrations"+"/")

	createFile(fs, "stubs/migration.stub", filename)
	fmt.Println("create called", filename)
}

func createFile(fs afero.Fs, stubPath, filePath string) {

	fs.Create(filePath)

	_, filename, _, _ := runtime.Caller(1)
	stubPath = path.Join(path.Dir(filename), stubPath)

	contents, _ := fileContents(stubPath)
	contents = replaceStub(contents, strings.Title(MigrationName))

	err := overwrite("db/migrations"+"/"+filePath, contents)
	if err != nil {
		fmt.Println(err)
	}
}

func fileContents(filePath string) (string, error) {
	a := afero.NewOsFs()
	contents, err := afero.ReadFile(a, filePath)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

func overwrite(file string, message string) error {
	a := afero.NewOsFs()
	return afero.WriteFile(a, file, []byte(message), 0666)
}

func replaceStub(content string, name string) string {
	content = strings.Replace(content, "{{TableName}}", strings.ToLower(Plural(name)), -1)
	content = strings.Replace(content, "{{Table}}", Plural(name), -1)
	return content
}

func Plural(name string) string {
	pluralize := pluralize.NewClient()

	return pluralize.Plural(name)
}

func Singular(name string) string {
	pluralize := pluralize.NewClient()
	return pluralize.Singular(name)
}

func Lower(name string) string {
	return strings.ToLower(name)
}
