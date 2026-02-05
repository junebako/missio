package core

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config は設定ファイルの構造を定義します
type Config struct {
	// 抽出するファイルパターン
	Include struct {
		// ファイル名のパターン（例: .env, master.key）
		Names []string `yaml:"names"`
		// 拡張子のパターン（例: .key, .pem）
		Extensions []string `yaml:"extensions"`
		// パスのパターン（例: config/credentials/*）
		Paths []string `yaml:"paths"`
	} `yaml:"include"`

	// 除外するファイルパターン
	Exclude struct {
		// ファイル名のパターン（例: .env.example）
		Names []string `yaml:"names"`
		// 拡張子のパターン（例: .md, .txt）
		Extensions []string `yaml:"extensions"`
		// パスのパターン（例: test/*, spec/*）
		Paths []string `yaml:"paths"`
	} `yaml:"exclude"`
}

// デフォルトの設定
var defaultConfig = Config{
	Include: struct {
		Names      []string `yaml:"names"`
		Extensions []string `yaml:"extensions"`
		Paths      []string `yaml:"paths"`
	}{
		Names: []string{
			".env",
			".env.",
			".envrc",
			"master.key",
			"credentials.yml.enc",
			"id_rsa",
			"id_ecdsa",
			"id_ed25519",
		},
		Extensions: []string{
			".key",
			".pem",
			".crt",
			".p12",
			".pfx",
			".jks",
			".keystore",
		},
		Paths: []string{
			"config/credentials/*",
			"config/master.key",
			".kamal/*",
		},
	},
	Exclude: struct {
		Names      []string `yaml:"names"`
		Extensions []string `yaml:"extensions"`
		Paths      []string `yaml:"paths"`
	}{
		Names: []string{
			".env.example",
			".env.sample",
			".env.template",
			"example",
			"sample",
			"template",
			"test",
			"spec",
		},
		Extensions: []string{
			".md",
			".txt",
			".example",
			".sample",
			".template",
		},
		Paths: []string{
			"test/*",
			"spec/*",
			"examples/*",
		},
	},
}

// LoadConfig は設定ファイルを読み込みます
func LoadConfig(rootDir string) (*Config, error) {
	// デフォルト設定をコピー
	config := defaultConfig

	// 設定ファイルのパスを構築
	configPath := filepath.Join(rootDir, "missio.yml")

	// 設定ファイルが存在しない場合はデフォルト設定を返す
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &config, nil
	}

	// 設定ファイルを読み込む
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// YAMLをパース
	var userConfig Config
	if err := yaml.Unmarshal(data, &userConfig); err != nil {
		return nil, err
	}

	// ユーザー設定をマージ
	if len(userConfig.Include.Names) > 0 {
		config.Include.Names = userConfig.Include.Names
	}
	if len(userConfig.Include.Extensions) > 0 {
		config.Include.Extensions = userConfig.Include.Extensions
	}
	if len(userConfig.Include.Paths) > 0 {
		config.Include.Paths = userConfig.Include.Paths
	}
	if len(userConfig.Exclude.Names) > 0 {
		config.Exclude.Names = userConfig.Exclude.Names
	}
	if len(userConfig.Exclude.Extensions) > 0 {
		config.Exclude.Extensions = userConfig.Exclude.Extensions
	}
	if len(userConfig.Exclude.Paths) > 0 {
		config.Exclude.Paths = userConfig.Exclude.Paths
	}

	return &config, nil
}
