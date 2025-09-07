#!/bin/bash

# マイグレーションファイル生成スクリプト
# 使用方法: ./migration.sh [app|vector] [title]
# 例: ./migration.sh app add_users_table

set -e

# 引数チェック
if [ $# -ne 2 ]; then
    echo "使用方法: $0 [app|vector] [title]"
    echo "例: $0 app add_users_table"
    exit 1
fi

TYPE=$1
TITLE=$2

# タイプの検証
if [ "$TYPE" != "app" ] && [ "$TYPE" != "vector" ]; then
    echo "エラー: タイプは 'app' または 'vector' である必要があります"
    exit 1
fi

# タイトルの検証（英数字とアンダースコアのみ許可）
if ! [[ "$TITLE" =~ ^[a-zA-Z0-9_]+$ ]]; then
    echo "エラー: タイトルは英数字とアンダースコアのみ使用できます"
    exit 1
fi

TIMESTAMP=$(date +%Y%m%d%H%M%S)

# ファイル名生成
UP_FILE="${TIMESTAMP}_${TITLE}.up.sql"
DOWN_FILE="${TIMESTAMP}_${TITLE}.down.sql"

# 対象ディレクトリ
TARGET_DIR="./${TYPE}"

# ディレクトリが存在しない場合は作成
mkdir -p "$TARGET_DIR"

# upファイル生成
cat > "${TARGET_DIR}/${UP_FILE}" << EOF
EOF

# downファイル生成
cat > "${TARGET_DIR}/${DOWN_FILE}" << EOF
EOF

echo "マイグレーションファイルが生成されました:"
echo "  Up:   ${TARGET_DIR}/${UP_FILE}"
echo "  Down: ${TARGET_DIR}/${DOWN_FILE}"
