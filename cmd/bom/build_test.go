package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// setupTestBuild はbuildテスト用の準備を行うヘルパー
func setupTestBuild(t *testing.T) (testDir string, cleanup func()) {
	t.Helper()
	testDir = t.TempDir() // テスト用一時ディレクトリ

	// 必要なディレクトリ構造作成
	componentsDir := filepath.Join(testDir, "examples", "components")
	if err := os.MkdirAll(componentsDir, 0755); err != nil {
		t.Fatalf("テスト用ディレクトリ構造の作成に失敗: %v", err)
	}

	// テスト用コンポーネントファイル (A, C)
	compAContent := "プロンプトAの内容です。\n"
	compCContent := "プロンプトCの内容です。\n"
	if err := os.WriteFile(filepath.Join(componentsDir, "A.md"), []byte(compAContent), 0644); err != nil {
		t.Fatalf("テスト用A.mdの作成に失敗: %v", err)
	}
	if err := os.WriteFile(filepath.Join(componentsDir, "C.md"), []byte(compCContent), 0644); err != nil {
		t.Fatalf("テスト用C.mdの作成に失敗: %v", err)
	}

	// カレントディレクトリを記録して変更
	origDir, _ := os.Getwd()
	if err := os.Chdir(testDir); err != nil {
		t.Fatalf("テストディレクトリへの移動に失敗: %v", err)
	}

	// 片付け関数を返す
	cleanup = func() {
		_ = os.Chdir(origDir)
		// t.TempDir()は自動でクリーンアップされる
	}

	return testDir, cleanup
}

// シンプルなBOM構造体
type simpleBOM struct {
	SchemaVersion string `yaml:"schema_version"`
	BOM           struct {
		Name        string `yaml:"name"`
		Version     string `yaml:"version"`
		Model       string `yaml:"model"`
		Description string `yaml:"description"`
		Metadata    struct {
			Owner   string `yaml:"owner"`
			License string `yaml:"license"`
		} `yaml:"metadata"`
	} `yaml:"bom"`
	Components []struct {
		ID          string `yaml:"id"`
		Version     string `yaml:"version"`
		Hash        string `yaml:"hash"`
		Description string `yaml:"description"`
		Metadata    struct {
			Owner string `yaml:"owner"`
		} `yaml:"metadata"`
	} `yaml:"components"`
}

// BOMファイル作成ヘルパー
func createBOMFile(t *testing.T, fileName string, componentIDs []string) string {
	t.Helper()

	bom := simpleBOM{
		SchemaVersion: "1.0.0",
	}
	bom.BOM.Name = "test-bom"
	bom.BOM.Version = "0.1.0"
	bom.BOM.Model = "test-model"
	bom.BOM.Description = "Test BOM"
	bom.BOM.Metadata.Owner = "test-owner"
	bom.BOM.Metadata.License = "MIT"

	for _, id := range componentIDs {
		component := struct {
			ID          string `yaml:"id"`
			Version     string `yaml:"version"`
			Hash        string `yaml:"hash"`
			Description string `yaml:"description"`
			Metadata    struct {
				Owner string `yaml:"owner"`
			} `yaml:"metadata"`
		}{
			ID:          id,
			Version:     "1.0.0",
			Hash:        "sha256:dummy",
			Description: fmt.Sprintf("Component %s", id),
		}
		component.Metadata.Owner = "test-owner"
		bom.Components = append(bom.Components, component)
	}

	// YAML書き出し
	data, err := yaml.Marshal(&bom)
	if err != nil {
		t.Fatalf("BOM YAMLの生成に失敗: %v", err)
	}

	if err := os.WriteFile(fileName, data, 0644); err != nil {
		t.Fatalf("BOMファイルの書き込みに失敗: %v", err)
	}

	return fileName
}

func TestBuildCmdDirectly(t *testing.T) {
	// テスト環境準備
	testDir, cleanup := setupTestBuild(t)
	defer cleanup()

	// 有効なBOMファイル作成
	validBOMFile := createBOMFile(t, "valid.bom.yaml", []string{"A", "C"})

	// 無効なBOMファイル作成（存在しないコンポーネント参照）
	invalidBOMFile := createBOMFile(t, "invalid.bom.yaml", []string{"A", "D"})

	// 出力ファイル
	outputFile := filepath.Join(testDir, "output.txt")

	tests := []struct {
		name        string
		bomFile     string
		outputFile  string
		expectError bool
		errorText   string // エラーテキストの一部
	}{
		{
			name:        "有効なBOMをファイルに出力",
			bomFile:     validBOMFile,
			outputFile:  outputFile,
			expectError: false,
		},
		{
			name:        "存在しないコンポーネント参照",
			bomFile:     invalidBOMFile,
			outputFile:  "",
			expectError: true,
			errorText:   "コンポーネントファイル読み込み失敗", // 部分文字列マッチでチェック
		},
		{
			name:        "存在しないBOMファイル",
			bomFile:     "nonexistent.yaml",
			outputFile:  "",
			expectError: true,
			errorText:   "BOM読み込みエラー", // 部分文字列マッチでチェック
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 事前に出力ファイルをクリア
			if tt.outputFile != "" {
				_ = os.Remove(tt.outputFile)
			}

			// コマンド実行準備
			rootCmd := &cobra.Command{Use: "bom"}
			tempBuildCmd := *buildCmd // Cobraコマンドのディープコピー
			rootCmd.AddCommand(&tempBuildCmd)

			args := []string{"build", tt.bomFile}
			if tt.outputFile != "" {
				args = append(args, "--output", tt.outputFile)
			}
			rootCmd.SetArgs(args)

			// 実行
			err := rootCmd.Execute()

			// 結果検証
			if tt.expectError {
				if err == nil {
					t.Errorf("エラーが期待されたが、成功した")
				} else if !strings.Contains(err.Error(), tt.errorText) {
					t.Errorf("エラーテキストが期待と異なる: got=%q, want=%q", err.Error(), tt.errorText)
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}

				// 出力ファイルの検証
				if tt.outputFile != "" {
					_, err := os.Stat(tt.outputFile)
					if os.IsNotExist(err) {
						t.Errorf("出力ファイルが生成されなかった: %s", tt.outputFile)
					}

					// ファイル内容の検証（オプション）
					content, err := os.ReadFile(tt.outputFile)
					if err != nil {
						t.Errorf("出力ファイルの読み込みに失敗: %v", err)
					}

					// Aプロンプトの内容が含まれているか
					if !strings.Contains(string(content), "プロンプトAの内容です") {
						t.Errorf("出力ファイルにプロンプトAの内容が含まれていない")
					}

					// ファイル削除
					_ = os.Remove(tt.outputFile)
				}
			}
		})
	}
}
