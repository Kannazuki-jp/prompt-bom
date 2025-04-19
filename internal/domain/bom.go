package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/Masterminds/semver/v3"
)

// Metadata構造体: BOMおよびComponentのメタデータを表す
// BOMではowner, licenseが必須。Componentではownerのみ必須（MVP仕様）
type Metadata struct {
	Owner   string // 所有者
	License string // ライセンス（BOMのみ必須）
}

// BOM構造体: プロンプトBOM全体を表す（MVP仕様準拠）
type BOM struct {
	Name        string          // BOM名
	Version     *semver.Version // セマンティックバージョン
	Model       string          // モデル名
	Description string          // 説明
	Metadata    Metadata        // メタデータ
	Components  []Component     // 含まれるコンポーネント一覧
}

// Component構造体: 各プロンプト部品を表す（MVP仕様準拠）
type Component struct {
	ID          string          // コンポーネントID
	Version     *semver.Version // セマンティックバージョン
	Hash        string          // sha256形式のハッシュ値
	Description string          // 説明
	Metadata    Metadata        // メタデータ（ownerのみ必須）
	// 将来拡張用: DependsOn []string // 依存関係
}

// ComputeSHA256: データのsha256ハッシュ値を 'sha256:<hex>' 形式で返す
func ComputeSHA256(data []byte) string {
	h := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(h[:])
}

// DetectCycle: 依存グラフにサイクル（循環依存）があるか検出する
// deps: map[コンポーネントID][]依存先コンポーネントID
// サイクルがあればtrueを返す
func DetectCycle(deps map[string][]string) bool {
	visited := make(map[string]bool)
	stack := make(map[string]bool)

	var visit func(string) bool
	visit = func(node string) bool {
		if stack[node] {
			return true // サイクル検出
		}
		if visited[node] {
			return false
		}
		visited[node] = true
		stack[node] = true
		for _, neighbor := range deps[node] {
			if visit(neighbor) {
				return true
			}
		}
		stack[node] = false
		return false
	}

	for node := range deps {
		if visit(node) {
			return true
		}
	}
	return false
}

// NewBOM: BOM構造体の生成ヘルパー（テスト用）
func NewBOM(name, version, model, desc, owner, license string, components []Component) (*BOM, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return nil, fmt.Errorf("BOMのバージョンが不正です: %w", err)
	}
	return &BOM{
		Name:        name,
		Version:     v,
		Model:       model,
		Description: desc,
		Metadata:    Metadata{Owner: owner, License: license},
		Components:  components,
	}, nil
}

// NewComponent: Component構造体の生成ヘルパー（テスト用）
func NewComponent(id, version, desc, owner string, data []byte) (*Component, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return nil, fmt.Errorf("コンポーネントのバージョンが不正です: %w", err)
	}
	hash := ComputeSHA256(data)
	return &Component{
		ID:          id,
		Version:     v,
		Hash:        hash,
		Description: desc,
		Metadata:    Metadata{Owner: owner},
	}, nil
}
