/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/ishanmadhav/supernetes/cmd/superctl/cmdapi"
	"github.com/ishanmadhav/supernetes/internals/utils"
	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Takes a JSON file path name as input and applies it to cluster.",
	Long:  `Takes a JSON file and applies it to cluster right now. Will change implementation to support yaml files soon`,
	Run: func(cmd *cobra.Command, args []string) {
		fileName, _ := cmd.Flags().GetString("file")

		if fileName == "" {
			fmt.Println("Please provide a file name")
		} else {
			deployment, err := utils.ParseDeploymentFileJSON(fileName)
			if err != nil {
				fmt.Print(err)
			}

			err = cmdapi.CreateDeploymentAPI(deployment)
			if err != nil {
				fmt.Print(err)
			}

			fmt.Println(deployment)
			fmt.Println("Deployment created successfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.PersistentFlags().StringP("file", "f", "", "The file to apply to the cluster")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
