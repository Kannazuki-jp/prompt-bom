package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// executeCommand はCobraコマンドを実行し、標準出力とエラーをキャプチャして返すヘルパー
func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.Execute()
	return buf.String(), err
}

func TestInitCmd(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		outputFile string // 期待される出力ファイル名
		wantOutput string // 期待される標準出力メッセージの一部
		wantErr    bool
	}{
		{
			name:       "デフォルト出力",
			args:       []string{"init"},
			outputFile: "prompt.bom.yaml",
			wantOutput: "prompt.bom.yaml を生成しました。",
			wantErr:    false,
		},
		{
			name:       "出力ファイル指定",
			args:       []string{"init", "--output", "custom.init.yaml"},
			outputFile: "custom.init.yaml",
			wantOutput: "custom.init.yaml を生成しました。",
			wantErr:    false,
		},
	}

	// ルートコマンドを準備（サブコマンド登録）
	rootCmdTest := &cobra.Command{Use: "bom"}
	rootCmdTest.AddCommand(initCmd) // init.goのinitCmdを使う

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト実行前に期待される出力ファイルを削除
			_ = os.Remove(tt.outputFile)

			output, err := executeCommand(rootCmdTest, tt.args...)

			if (err != nil) != tt.wantErr {
				t.Errorf("executeCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(output, tt.wantOutput) {
				t.Errorf("executeCommand() output = %v, want %v", output, tt.wantOutput)
			}

			// ファイルが生成されたか確認
			if _, err := os.Stat(tt.outputFile); os.IsNotExist(err) {
				t.Errorf("期待されたファイル '%s' が生成されませんでした", tt.outputFile)
			}

			// ファイル内容の検証（オプション）
			content, _ := os.ReadFile(tt.outputFile)
			if !strings.Contains(string(content), "schema_version: \"1.0.0\"") {
				t.Errorf("生成されたファイルの内容が期待と異なります")
			}

			// テスト後にファイルを削除
			_ = os.Remove(tt.outputFile)
		})
	}
}
