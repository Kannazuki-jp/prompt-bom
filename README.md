# prompt-bom

<!-- ロゴ（後で docs/img/logo.png に差し替え） -->
![prompt-bom logo](https://via.placeholder.com/300x100.png?text=prompt-bom)

<!-- バッジ（後で実際のCIやバージョンに合わせて更新） -->
[![Go Version](https://img.shields.io/badge/go-1.22+-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
<!-- [![Build Status](https://img.shields.io/github/actions/workflow/status/<your-org>/prompt-bom/go.yml?branch=main)](https://github.com/<your-org>/prompt-bom/actions) -->

プロンプト管理のためのBOM（部品表）CLIツール。

生成AIプロンプトの複雑化に対応し、BOM/SBOMの知見を応用したプロンプト管理基盤を提供します。
YAML形式のBOM定義とGo製CLIツールで最小限の管理機能をOSSとして実装し、
`init`/`validate`/`build`等のCLIコマンドを備え、拡張性・ガバナンス・監査要件にも対応します。

## 特徴

- **部品としてのプロンプト管理**: プロンプトを部品（コンポーネント）として定義・管理
- **YAMLによる宣言的定義**: `prompt.bom.yaml` でBOM構造を宣言的に記述
- **シンプルなCLI**: `init`, `validate`, `build` の基本コマンドで簡単操作
- **拡張性**: 将来的な外部ツール連携（Dify, LangChain等）を考慮した設計

## 目次

- [インストール](#インストール)
- [Quick Start (60秒体験)](#quick-start-60秒体験)
- [コマンド一覧](#コマンド一覧)
- [アーキテクチャ概要](#アーキテクチャ概要)
- [ライセンス](#ライセンス)
- [ロードマップ](#ロードマップ)

## インストール

### Go

```bash
go install github.com/kannazuki/prompt-bom/cmd/bom@latest
```
*(注意: モジュールパスは今後変更される可能性があります)*


## Quick Start 

1.  **BOMテンプレート生成:**

    ```bash
    bom init
    # -> prompt.bom.yaml を生成しました。
    ```

2.  **サンプルコンポーネント作成:**

    ```bash
    mkdir -p examples/components
    echo "これは部品Aです。" > examples/components/partA.md
    echo "これは部品Bです。" > examples/components/partB.md
    ```

3.  **`prompt.bom.yaml` を編集:**
    `components:` セクションに以下を追加します。

    ```yaml
    components:
      - id: "partA"
        version: "1.0.0"
        hash: "sha256:dummy_hash_a"
        description: "Part A prompt"
        metadata:
          owner: "your-team"
      - id: "partB"
        version: "1.0.0"
        hash: "sha256:dummy_hash_b"
        description: "Part B prompt"
        metadata:
          owner: "your-team"
    ```
    *(注意: `hash` はダミー値です。将来的に自動生成・検証機能が追加されます)*

4.  **BOM検証:**

    ```bash
    bom validate prompt.bom.yaml
    # -> OK: スキーマと必須フィールド検証に合格
    ```

5.  **プロンプト結合:**

    ```bash
    bom build prompt.bom.yaml
    # 標準出力に以下が表示される:
    # これは部品Aです。
    #
    # これは部品Bです。
    #
    ```

    ファイルに出力する場合:

    ```bash
    bom build prompt.bom.yaml -o final.prompt.txt
    # -> final.prompt.txt に結合結果を出力しました。
    ```

## コマンド一覧

| コマンド         | 説明                                    |
| ---------------- | --------------------------------------- |
| `bom init`       | BOMテンプレートYAML (`prompt.bom.yaml`) を生成 |
| `bom validate`   | BOM YAMLのスキーマ・必須フィールドを検証  |
| `bom build`      | BOMに基づきコンポーネントを結合して出力   |

詳細は [`docs/usage.md`](docs/usage.md) を参照してください。

## アーキテクチャ概要

```plaintext
prompt-bom/
├── cmd/bom/          # CLIエントリーポイントとコマンド実装
│   ├── main.go
│   ├── init.go
│   ├── validate.go
│   └── build.go
├── internal/
│   ├── domain/       # BOM/Component等のコア構造体、ビジネスロジック
│   │   └── bom.go
│   ├── app/          # (予定) アプリケーション層
│   └── infra/        # (予定) ファイルI/O, 外部連携等
├── spec/             # 仕様ファイル
│   └── prompt.bom.schema.json # BOM YAMLのJSON Schema
├── examples/         # 利用例
│   └── components/   # サンプルコンポーネントファイル (*.md)
├── docs/             # ドキュメント
│   ├── testing.md
│   └── img/          # (ロゴ画像)
├── go.mod
├── go.sum
└── README.md
```

## ライセンス

[MIT License](https://opensource.org/licenses/MIT)

## ロードマップ

MVP完了後、以下の機能拡張を計画しています。

- Dify連携 
- LangChain連携 
- 回帰テスト自動化 
- 階層化 


