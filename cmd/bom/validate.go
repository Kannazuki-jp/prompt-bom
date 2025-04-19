package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

// validateCmd: bom validateコマンドのCobraコマンド定義
var validateCmd = &cobra.Command{
	Use:   "validate <path>",
	Short: "BOM YAMLのスキーマ・必須フィールド検証を行う",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		schemaPath := "spec/prompt.bom.schema.json"

		// YAMLファイルをJSONに変換してロード
		jsonData, err := yamlFileToJSON(path)
		if err != nil {
			return fmt.Errorf("%s: YAMLパースエラー: %w", path, err)
		}

		// JSON Schema検証
		schemaLoader := gojsonschema.NewReferenceLoader("file://" + getAbsPath(schemaPath))
		jsonLoader := gojsonschema.NewBytesLoader(jsonData)
		result, err := gojsonschema.Validate(schemaLoader, jsonLoader)
		if err != nil {
			return fmt.Errorf("スキーマ検証エラー: %w", err)
		}
		if !result.Valid() {
			var errMsg string
			for _, e := range result.Errors() {
				// わかりやすいエラーメッセージ形式
				errMsg += fmt.Sprintf("Error: %s: %s\n", path, e.String())
			}
			// os.Exit(1) の代わりにエラーを返す
			return fmt.Errorf("%s", errMsg)
		}
		fmt.Println("OK: スキーマと必須フィールド検証に合格")
		return nil
	},
}

func init() {
	// validateコマンドをルートに登録する場合はmain.goで実施済み
}

// yamlFileToJSON: YAMLファイルをJSONバイト列に変換
func yamlFileToJSON(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var obj interface{}
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return nil, err
	}
	return json.Marshal(obj)
}

// getAbsPath: 相対パスを絶対パスに変換
// NOTE: テスト時はテスト用一時ディレクトリがカレントになるため注意
func getAbsPath(path string) string {
	// テスト以外で実行される場合、実行ファイルのパスを基準にするなど、より堅牢な方法を検討
	abs, err := filepath.Abs(path)
	if err != nil {
		// エラー時は元のパスを返す（ベストエフォート）
		return path
	}
	return abs
}
