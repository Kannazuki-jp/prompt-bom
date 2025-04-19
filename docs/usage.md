# Usage ドキュメント

本文書は `prompt-bom` CLIツールの詳細な使い方について説明します。

## Commands

### `bom init`

**目的:** 新しいプロンプトBOMプロジェクトの初期化、または既存のディレクトリに `prompt.bom.yaml` テンプレートを生成します。

**書式:**
```bash
bom init [flags]
```

**フラグ:**
- `-o, --output FILE`: 出力ファイル名を指定します。 (デフォルト: `prompt.bom.yaml`)

**実行例:**
```bash
# カレントディレクトリに prompt.bom.yaml を生成
bom init

# 指定したパスに my_prompt.yaml を生成
bom init --output my_prompt.yaml
```

### `bom validate`

**目的:** 指定された `prompt.bom.yaml` ファイルが、定義されたJSON Schemaに準拠しているか、また必須フィールドが満たされているかを検証します。

**書式:**
```bash
bom validate <path>
```

**引数:**
- `<path>`: 検証対象の `prompt.bom.yaml` ファイルへのパス。（必須）

**実行例:**
```bash
# prompt.bom.yaml を検証
bom validate prompt.bom.yaml
# -> OK: スキーマと必須フィールド検証に合格

# 不正なファイルを検証（エラー例）
bom validate invalid.yaml
# -> Error: invalid.yaml: bom: name is required (標準エラー出力)
```

### `bom build`

**目的:** 指定された `prompt.bom.yaml` の `components` セクションに記述された順序に従い、対応するコンポーネントファイル (`examples/components/<id>.md`) の内容を結合して出力します。

**書式:**
```bash
bom build <path> [flags]
```

**引数:**
- `<path>`: 処理対象の `prompt.bom.yaml` ファイルへのパス。（必須）

**フラグ:**
- `-o, --output FILE`: 結合結果を指定したファイルに出力します。省略した場合は標準出力に出力されます。

**実行例:**
```bash
# prompt.bom.yaml に基づき、結合結果を標準出力へ
bom build prompt.bom.yaml

# 結合結果を final.prompt.txt に出力
bom build prompt.bom.yaml --output final.prompt.txt
```

## Options

現在、各コマンドに固有のオプション（フラグ）があります。

- **`init`, `build` 共通:**
    - `-o, --output FILE`: 出力ファイル名を指定します。

## Error Reference

コマンド実行時に発生する可能性のある主なエラーです。

- **`bom validate` 時:**
    - `Error: <path>: YAMLパースエラー: open <path>: no such file or directory`
        - 原因: 指定されたYAMLファイルが見つかりません。
        - 対処: ファイルパスが正しいか確認してください。
    - `Error: <path>: <field> is required`
        - 原因: JSON Schemaで必須と定義されているフィールド（例: `bom.name`）がYAMLファイルに存在しません。
        - 対処: YAMLファイルに必要なフィールドを追加してください。
    - `Error: <path>: <field>: <reason>` (例: `version: Does not match pattern '^\d+\.\d+\.\d+$'`) 
        - 原因: フィールドの値がJSON Schemaで定義された型やパターンに一致しません。
        - 対処: スキーマに従ってフィールド値を修正してください。
- **`bom build` 時:**
    - `<path>: BOM読み込みエラー: ...`
        - 原因: 指定されたBOM YAMLファイルの読み込みまたはパースに失敗しました。
        - 対処: ファイルパスやYAMLの構文を確認してください。
    - `<component_path>: コンポーネントファイル読み込み失敗: open ... no such file or directory`
        - 原因: BOM YAMLで参照されているコンポーネントファイル (`examples/components/<id>.md`) が見つかりません。
        - 対処: コンポーネントファイルが正しい場所に存在するか、BOM YAMLの `id` が正しいか確認してください。

## FAQ

**Q1: `components` の `hash` フィールドは何のためにありますか？**

A1: 各コンポーネントプロンプトの内容に対するSHA256ハッシュ値を格納するためのフィールドです。MVP段階では形式チェックのみですが、将来的にはファイル内容の検証、改ざん検知、キャッシュなどに利用される予定です。

**Q2: コンポーネントファイル (`*.md`) はどこに置けば良いですか？**

A2: 現在の実装では、プロジェクトルート直下の `examples/components/` ディレクトリに `<id>.md` という名前で配置することを想定しています。`<id>` は `prompt.bom.yaml` の `components` セクションで指定した `id` と一致させる必要があります。

**Q3: BOMファイルやコンポーネントファイルのバージョン管理はどうすれば良いですか？**

A3: Gitなどのバージョン管理システムを利用して、`prompt.bom.yaml` と `examples/components/` 内のファイルを一緒に管理することを推奨します。 