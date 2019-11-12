/*
Copyright © 2019 AltoStack <info@altostack.io>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	dynamoTableName     string
	dynamoBatchSize     int64
	dynamoAppendRestore bool
	s3BucketName        string
	s3BucketFolderName  string
	s3DateSuffix        bool
	waitTime            int64
	origin              string
	destination         string
)

var rootCmd = &cobra.Command{
	Use:   "dynamodump",
	Short: "AWS DynamoDB Backup and Restores",
	Long: `
Dynamodump allows for easier and cheaper actions and restores of DynamoDB Tables.
Backups are compatible with the AWS DataPipeline functionality.
	   
It is also capable of restoring a actions from s3 to a given table both from this
tool or from a actions generated using the AWS DataPipeline functionality.
to quickly create a Cobra application.
  `,
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(func() {
		viper.SetEnvPrefix("dyn") // prefix that ENVIRONMENT variables will use.
		viper.AutomaticEnv()      // read in environment variables that match
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

		postInitCommands(rootCmd.Commands())
	})

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
}

func postInitCommands(commands []*cobra.Command) {
	for _, cmd := range commands {
		presetRequiredFlags(cmd)
		if cmd.HasSubCommands() {
			postInitCommands(cmd.Commands())
		}
	}
}

func presetRequiredFlags(cmd *cobra.Command) {
	viper.BindPFlags(cmd.Flags())
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			cmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}
