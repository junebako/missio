package core

import (
	"os"
	"path/filepath"
	"strings"
)

// Scanner は指定されたディレクトリ内の秘匿ファイルを検出します
type Scanner struct {
	rootDir string
	logger  *Logger
	config  *Config
}

// 除外するディレクトリ名のリスト
var excludeDirs = []string{
	".git",          // Gitリポジトリ
	"node_modules",  // Node.js依存関係
	"vendor",        // 依存関係（Go, PHP, Ruby等）
	"tmp",           // 一時ファイル
	"cache",         // キャッシュファイル
	"log",           // ログファイル
	"logs",          // ログファイル
	"coverage",      // テストカバレッジ
	"dist",          // ビルド成果物
	"build",         // ビルド成果物
	".bundle",       // Bundler
	".gradle",       // Gradle
	".idea",         // IntelliJ IDEA
	".vscode",       // Visual Studio Code
	"__pycache__",   // Python
	".pytest_cache", // Python
}

// NewScanner は新しいScannerインスタンスを作成します
func NewScanner(rootDir string, verbose bool, maxDepth int) (*Scanner, error) {
	// 設定ファイルを読み込む
	config, err := LoadConfig(rootDir)
	if err != nil {
		return nil, err
	}

	return &Scanner{
		rootDir: rootDir,
		logger:  NewLogger(rootDir, verbose, maxDepth),
		config:  config,
	}, nil
}

// Scan はディレクトリを走査し、秘匿ファイルのリストを返します
func (s *Scanner) Scan() ([]string, error) {
	var files []string

	err := filepath.Walk(s.rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// 除外ディレクトリをスキップ
			for _, dir := range excludeDirs {
				if info.Name() == dir {
					return filepath.SkipDir
				}
			}
			s.logger.LogProgress(path)
			return nil
		}

		s.logger.IncrementScanned()
		s.logger.LogProgress(path)

		// パスを相対パスに変換
		relPath, err := filepath.Rel(s.rootDir, path)
		if err != nil {
			return err
		}

		// 秘匿ファイルかどうかをチェック（パターンマッチ優先）
		if s.isSecretFile(relPath) {
			files = append(files, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.logger.LogSummary(files)
	return files, nil
}

// isSecretFile はファイルが秘匿情報を含むかどうかを判定します
func (s *Scanner) isSecretFile(relPath string) bool {
	filename := filepath.Base(relPath)
	ext := filepath.Ext(relPath)
	lowerFilename := strings.ToLower(filename)
	lowerPath := strings.ToLower(relPath)

	// 除外パターンをチェック
	for _, pattern := range s.config.Exclude.Names {
		if strings.Contains(lowerFilename, strings.ToLower(pattern)) {
			return false
		}
	}

	for _, pattern := range s.config.Exclude.Extensions {
		if strings.EqualFold(ext, pattern) {
			return false
		}
	}

	for _, pattern := range s.config.Exclude.Paths {
		if matchPathPattern(pattern, lowerPath) {
			return false
		}
	}

	// 秘匿ファイルパターンをチェック
	for _, pattern := range s.config.Include.Names {
		if strings.Contains(lowerFilename, strings.ToLower(pattern)) {
			return true
		}
	}

	for _, pattern := range s.config.Include.Extensions {
		if strings.EqualFold(ext, pattern) {
			return true
		}
	}

	for _, pattern := range s.config.Include.Paths {
		if matchPathPattern(pattern, lowerPath) {
			return true
		}
	}

	return false
}

// matchPathPattern はパスパターンを相対パスの各サブパスに対してマッチングします。
// 例: パターン ".kamal/*" は "project/.kamal/secrets" にもマッチします。
func matchPathPattern(pattern, path string) bool {
	for {
		if matched, _ := filepath.Match(pattern, path); matched {
			return true
		}
		// 先頭のディレクトリを1つ削って次のサブパスを試す
		i := strings.IndexByte(path, filepath.Separator)
		if i < 0 {
			return false
		}
		path = path[i+1:]
	}
}
