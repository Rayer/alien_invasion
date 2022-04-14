/*
Copyright © 2022 Rayer

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	alien_invastion "alien-invastion"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "alien-invastion <mapfile path> <alien count>",
	Short: "A game of alien invasion",
	Long: `A game of alien invasion. Given a map file and a number of aliens, aliens will try to invade the cities in the map.
Roles are : 
1. Alien will enter a random city
2. Alien will try to enter an adjacent city
3. When 2 aliens enters same city, they will fight, and result the city being destroyed.
4. If city is destroyed, all path lead to, and leads from this city, will be removed, preventing other aliens from entering or exiting.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return cmd.Help()
		}
		parser := alien_invastion.StreamParser{}
		gameMap, errors := parser.ParseFile(args[0])
		if errors != nil && len(errors) > 0 {
			return fmt.Errorf("%v", errors)
		}
		alienCount, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		aliens := make([]*alien_invastion.Alien, 0)
		for i := 0; i < alienCount; i++ {
			aliens = append(aliens, alien_invastion.NewAlien())
		}
		err = gameMap.AssignAliens(aliens)
		if err != nil {
			return err
		}
		for {
			if gameMap.Update() == false {
				break
			}
		}
		fmt.Println(gameMap.DumpMap())
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.alien-invastion.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
