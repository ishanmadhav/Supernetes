/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/ishanmadhav/supernetes/api"
	"github.com/ishanmadhav/supernetes/cmd/superctl/cmdapi"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
		deploymentName, _ := cmd.Flags().GetString("name")

		if deploymentName == "" {
			fmt.Println("Please provide a name for the deployment")
		} else {
			image, _ := cmd.Flags().GetString("image")
			port, _ := cmd.Flags().GetString("port")
			replicas, _ := cmd.Flags().GetUint("replicas")
			selector, _ := cmd.Flags().GetString("selector")

			//Make an http call to the superapiserver using this deployment object
			tempDeployment := api.Deployment{
				Name:     deploymentName,
				Image:    image,
				Port:     port,
				Replicas: replicas,
				Selector: selector,
			}

			fmt.Print(tempDeployment)
			err := cmdapi.CreateDeploymentAPI(tempDeployment)

			if err != nil {
				fmt.Println("There was an error")
				fmt.Print(err)
				return
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().String("name", "", "Name of the resource to create")
	createCmd.PersistentFlags().String("image", "", "Image of the deployment to create")
	createCmd.PersistentFlags().String("port", "", "Port of the deployment to create")
	createCmd.PersistentFlags().Uint("replicas", 1, "Replicas of the deployment to create")
	createCmd.PersistentFlags().String("selector", "", "Selector of the deployment to create")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
