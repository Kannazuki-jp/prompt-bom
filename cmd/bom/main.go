package main

import (
	"os"

	"github.com/spf13/cobra"
)

// ルートコマンド
var rootCmd = &cobra.Command{
	Use:   "bom",
	Short: "プロンプトBOM管理CLIツール",
}

func main() {
	// サブコマンドを登録
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(buildCmd)

	// コマンド実行
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
