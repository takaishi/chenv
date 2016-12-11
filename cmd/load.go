// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/getwe/figlet4go"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Printf("Failed load $HOME/.chenv.yaml, %s", err)
			return
		}

		key := args[0]
		nextEnv := getNextEnv(config, key)

		printBanner(key)

		for k, v := range nextEnv {
			os.Setenv(k.(string), v.(string))
		}

		shell := os.Getenv("SHELL")

		procAttr := new(os.ProcAttr)
		newEnv := []string{"CHENV=" + key}
		procAttr.Env = newEnv
		procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
		shellArgs := []string{
			//"--xtrace",
		}
		if p, err := os.StartProcess(shell, shellArgs, procAttr); err != nil {
			fmt.Println("[ERROR] Failed to StartProcess")
			fmt.Println("[ERROR] Failed to StartProcess, %s", err)
			return
		} else {
			_, err := p.Wait()
			if err != nil {
				fmt.Println("[ERROR]: failed wait: %s", err)
				return
			}
		}
	},
}

func loadConfig() (map[interface{}]interface{}, error) {
	buf, err := ioutil.ReadFile(os.Getenv("HOME") + "/.chenv.yaml")
	if err != nil {
		return nil, err
	}

	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		panic(err)
	}
	return m, nil
}

func getNextEnv(config map[interface{}]interface{}, key string) map[interface{}]interface{} {
	return config[key].(map[interface{}]interface{})
}

func printBanner(s string) {
	ascii := figlet4go.NewAsciiRender()
	renderStr, _ := ascii.Render(s)
	fmt.Println(renderStr)
}

func init() {
	RootCmd.AddCommand(loadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
