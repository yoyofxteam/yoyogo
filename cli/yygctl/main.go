package main

import (
	"github.com/spf13/cobra"
	"yygctl/cmds"
)

var rootCmd = &cobra.Command{
	Use:   "yygctl",
	Short: "A generator for Cobra based Applications",
	Long: `yygctl is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a yoyogo application.`,
}

func main() {
	//templates.GetProjectByName("console").List()
	rootCmd.AddCommand(cmds.VersionCmd)
	rootCmd.AddCommand(cmds.RunCmd)
	rootCmd.AddCommand(cmds.BuildCmd)
	rootCmd.AddCommand(cmds.NewCmd)
	_ = rootCmd.Execute()

}
