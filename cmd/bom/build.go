package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// buildCmd: bom buildコマンドのCobraコマンド定義
var buildCmd = &cobra.Command{
	Use:   "build <path>",
	Short: "BOMのcomponents順にプロンプトを結合し出力する",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bomPath := args[0]
		output, _ := cmd.Flags().GetString("output")

		// BOM YAMLを読み込む
		bom, err := loadBOM(bomPath)
		if err != nil {
			return fmt.Errorf("%s: BOM読み込みエラー: %w", bomPath, err)
		}

		// components順に/examples/components/<id>.prompt.txtを結合
		var result string
		for _, comp := range bom.Components {
			filePath := filepath.Join("examples", "components", comp.ID+".md")
			data, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("%s: コンポーネントファイル読み込み失敗: %w", filePath, err)
			}
			result += string(data) + "\n"
		}

		// 出力先に書き込み or STDOUT
		if output != "" {
			err := os.WriteFile(output, []byte(result), 0644)
			if err != nil {
				return fmt.Errorf("出力ファイル書き込み失敗: %w", err)
			}
			fmt.Printf("%s に結合結果を出力しました。\n", output)
		} else {
			fmt.Print(result)
		}
		return nil
	},
}

func init() {
	// --output フラグを追加
	buildCmd.Flags().StringP("output", "o", "", "出力ファイル名 (省略時はSTDOUT)")
}

// loadBOM: BOM YAMLを読み込んで構造体にデコード
func loadBOM(path string) (*struct {
	Components []struct {
		ID string `yaml:"id"`
	} `yaml:"components"`
}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var bom struct {
		Components []struct {
			ID string `yaml:"id"`
		} `yaml:"components"`
	}
	if err := yaml.Unmarshal(data, &bom); err != nil {
		return nil, err
	}
	return &bom, nil
}
