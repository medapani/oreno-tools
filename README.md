# oreno-tools

Wails (Go + React) で作られた、開発者向けデスクトップユーティリティ集です。
ネットワーク、エンコード、JWT、証明書まわりの作業を 1 つのアプリで扱えます。

## 主な機能

- バイト変換
- 通信速度変換
- データ転送速度変換
- データ転送時間計算
- CIDR 計算
- Base64 変換
- JWT 検証・デコード・生成
- JSON/YAML 変換
- URL エンコード/デコード
- 基数変換/計算
- 鍵ペア作成/検証
- 自己署名証明書作成
- mTLS 証明書作成
- CRL 更新

## 技術スタック

- Backend: Go 1.23
- Frontend: React + TypeScript + Vite + Tailwind CSS
- Desktop Runtime: Wails v2

## 前提環境

- Go 1.23 以上
- Node.js / npm
- Wails CLI (`wails` コマンド)

Wails CLI 未導入の場合:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

## セットアップ

```bash
cd /path/to/oreno-tools
go mod tidy
cd frontend && npm install
cd ..
```

## 開発

```bash
wails dev
```

- フロントエンドは `wails.json` の設定により `npm run dev` で監視されます。
- Backend の Go コード変更も反映されます。

## クイックスタート

1. `wails dev` を実行してアプリを起動
2. 左メニューから使いたいツールを選択
3. 入力値を設定して変換/生成/検証を実行
4. 必要ならコピーまたはファイル保存

## 使い方例

### バイト変換

- 入力: `1024` + 単位 `MB`
- 期待する確認ポイント: `GB` や `GiB` など複数単位へ同時変換される

### データ転送時間計算

- 入力: データ容量 `10 GB`、転送速度 `100 Mbps`
- 期待する確認ポイント: 秒/分/時/日で理論転送時間を比較できる

### JWT 検証

- 入力: JWT 文字列と秘密鍵(または公開鍵)
- 期待する確認ポイント: ヘッダー/ペイロード整形表示、署名検証結果 (`valid`)

### 自己署名証明書作成

- 入力: `Common Name`、`Organization`、SAN、有効日数、アルゴリズム
- 期待する確認ポイント: PEM 形式の証明書と秘密鍵が生成される

### mTLS 証明書作成

- 入力: CA/Server/Client の Common Name、SAN、有効日数
- 期待する確認ポイント: CA・サーバー・クライアントの証明書セットが生成される

### CRL 更新

- 入力: CA 証明書、CA 秘密鍵、失効対象証明書
- 期待する確認ポイント: 更新済み CRL と失効シリアル一覧が生成される

## ビルド

Makefile を使う場合:

```bash
make build
```

- macOS 向けアプリをビルドします。

Windows 向けビルド:

```bash
make build-win
```

## Lint

```bash
make lint
```

`golangci-lint` が未導入の場合は Makefile で自動インストールされます。

## 出力と保存先

- Wails ビルド成果物: `build/bin/`
- 証明書/CRL ダウンロード保存先: `~/Downloads/oreno-tools-certs/`

## ディレクトリ構成

```text
oreno-tools/
	backend/      # Go 側の各ツール実装
	frontend/     # React UI
	build/        # Wails ビルド関連
	app.go        # Wails バインディング
	main.go       # アプリ起動エントリポイント
	wails.json    # Wails 設定
	Makefile
```

## セキュリティ上の注意

- 秘密鍵、JWT シークレット、証明書データは機密情報です。
- 共有端末での利用時は、保存ファイルとクリップボードの取り扱いに注意してください。

## トラブルシュート

- `wails: command not found`

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

- フロントエンド依存が不足している

```bash
cd frontend
npm install
```

- 開発モードで表示が崩れる/更新されない
	- `frontend/node_modules` を再インストールし、`wails dev` を再起動してください。
