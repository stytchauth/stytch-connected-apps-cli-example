package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/stytchauth/stytch-connected-apps-cli-example/utils"
)

// getResourcesCmd represents the getResources command
var getResourcesCmd = &cobra.Command{
	Use:   "get_resources",
	Short: "Gets /api/resources, passing in a connected apps access token",
	Run: func(cmd *cobra.Command, args []string) {
		// Get access token from keyring
		access_token, _ := utils.LoadToken(utils.AccessToken)

		// GET /api/resources with the access token
		url := "http://localhost:3001/api/resources"
		var bearer = "Bearer " + access_token 

		req, err := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", bearer)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Failed to get /api/resources with error: ", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error while reading response bytes: ", err)
		}
		fmt.Println(string(body))
	},
}

func init() {
	rootCmd.AddCommand(getResourcesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getResourcesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getResourcesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
