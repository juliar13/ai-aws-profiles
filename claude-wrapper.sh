#!/bin/bash

# Slack Webhook URL
SLACK_WEBHOOK_URL="https://hooks.slack.com/services/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"
SLACK_CHANNEL="#claude-code"

# プロジェクト情報
PROJECT_NAME="ai-aws-profiles"
PROJECT_DIR=$(pwd)

# Slack通知関数
send_slack_notification() {
    local message="$1"
    local color="$2"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    
    local payload=$(cat <<EOF
{
    "channel": "$SLACK_CHANNEL",
    "username": "Claude Code",
    "icon_emoji": ":robot_face:",
    "attachments": [
        {
            "color": "$color",
            "fields": [
                {
                    "title": "Project",
                    "value": "$PROJECT_NAME",
                    "short": true
                },
                {
                    "title": "Directory",
                    "value": "$PROJECT_DIR",
                    "short": true
                },
                {
                    "title": "Time",
                    "value": "$timestamp",
                    "short": true
                },
                {
                    "title": "Status",
                    "value": "$message",
                    "short": false
                }
            ]
        }
    ]
}
EOF
    )
    
    curl -X POST -H 'Content-type: application/json' \
         --data "$payload" \
         "$SLACK_WEBHOOK_URL" \
         --silent > /dev/null
}

# ログファイル
LOG_FILE="/tmp/claude-code-session-$$.log"

# 実行開始通知
echo "Starting Claude Code session..."

# Claude Codeを実行し、出力を監視
claude "$@" 2>&1 | tee "$LOG_FILE" | while IFS= read -r line; do
    echo "$line"
    
    # ユーザー許可が必要なパターンを検出
    if echo "$line" | grep -qE "(Do you want to|Would you like to|Proceed|Continue|Allow|Confirm|Permission required|Authorization needed|\[Y/n\]|\[y/N\])"; then
        send_slack_notification "🔐 Permission required: $line" "warning"
    elif echo "$line" | grep -qE "(Error|Failed|Exception|denied|refused)"; then
        send_slack_notification "❌ Error occurred: $line" "danger"
    fi
done

# 実行完了の判定
EXIT_CODE=${PIPESTATUS[0]}

if [ $EXIT_CODE -eq 0 ]; then
    send_slack_notification "✅ Claude Code session completed successfully" "good"
else
    send_slack_notification "❌ Claude Code session failed (exit code: $EXIT_CODE)" "danger"
fi

# ログファイルを削除
rm -f "$LOG_FILE"

exit $EXIT_CODE