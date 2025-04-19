package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// MVP仕様に基づく最小テンプレートYAML
const defaultBOMYAML = `schema_version: "1.0.0"
bom:
  name: "my-bom"
  version: "0.1.0"
  model: "gpt-4o-2025-05"
  description: "BOM description"
  metadata:
    owner: "your-team"
    license: "MIT"
components: []
`

// initCmd: bom initコマンドのCobraコマンド定義
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "BOMテンプレートYAMLを生成する",
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")
		if output == "" {
			output = "prompt.bom.yaml"
		}
		// 既存ファイルがある場合は上書き確認（MVPでは単純に上書き）
		err := os.WriteFile(output, []byte(defaultBOMYAML), 0644)
		if err != nil {
			return fmt.Errorf("YAMLファイルの書き込みに失敗しました: %w", err)
		}
		cmd.Printf("%s を生成しました。\n", output)
		return nil
	},
}

func init() {
	// --output フラグを追加
	initCmd.Flags().StringP("output", "o", "", "出力ファイル名 (デフォルト: prompt.bom.yaml)")
}
