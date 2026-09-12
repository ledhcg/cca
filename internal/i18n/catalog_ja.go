package i18n

var jaCatalog = map[Key]string{
	// Common / Generic
	KeyCommonDefault:   "デフォルト",
	KeyCommonNone:      "なし",
	KeyCommonCancelled: "キャンセルしました",
	KeyCommonOk:        "正常",
	KeyCommonPresent:   "存在",
	KeyCommonActive:    "有効",

	// Commands - ls
	KeyCmdLsHeaderProfile: "プロファイル",
	KeyCmdLsHeaderStatus:  "ステータス",
	KeyCmdLsHeaderAccount: "アカウント",
	KeyCmdLsHeaderPlan:    "プラン",
	KeyCmdLsStatusIn:      "ログイン中",
	KeyCmdLsStatusOut:     "未ログイン",
	KeyCmdLsNoProfiles:    "追加プロファイルがまだありません。作成: %scca new work --login%s",

	// Commands - info
	KeyCmdInfoConfigDir:       "設定ディレクトリ",
	KeyCmdInfoConfigDirNotSet: "(CLAUDE_CONFIG_DIR 未設定)",
	KeyCmdInfoLoggedIn:        "ログイン済み",
	KeyCmdInfoStatus:          "状態",
	KeyCmdInfoProjectsOpened:  "%s · %d 個のプロジェクトを開きました",
	KeyCmdInfoPrivateSize:     "専有サイズ",
	KeyCmdInfoShared:          "共有項目",

	// Commands - doctor
	KeyCmdDoctorProfilesHeading: "プロファイル一覧",
	KeyCmdDoctorBrokenLinks:     "リンク切れ: %s",
	KeyCmdDoctorMissingLinks:    "未作成リンク: %s (実行: cca sync --all)",
	KeyCmdDoctorLoggedInNoCred:  "ログイン中ですが %s が見つかりません: %s",
	KeyCmdDoctorHasCredNoLogin:  "%s はありますがログインできません",
	KeyCmdDoctorOrphansHeading:  "%s どのプロファイルにも属していません",
	KeyCmdDoctorOrphanHint:      "(手動で削除されたプロファイル?)",

	// Commands - new
	KeyCmdNewMissingName:   "プロファイル名が指定されていません — 例: cca new work",
	KeyCmdNewAlreadyExists: "プロファイル '%s' は既に %s に存在します",
	KeyCmdNewCreated:       "プロファイルを作成しました %s%s%s → %s",
	KeyCmdNewLoginHint:     "ログイン:  %scca login %s%s",
	KeyCmdNewYoloNote:      "  %s⚡ defaultMode = bypassPermissions%s %s(このプロファイルのみ)%s",

	// Commands - login / logout
	KeyCmdLoginOpeningBrowser: "ブラウザを開いてプロファイル '%s' にログインしています…",
	KeyCmdLoginStillOut:       "'%s' はまだログインしていません",
	KeyCmdLogoutSuccess:       "プロファイル '%s' からログアウトしました",

	// Commands - use / sh / exec
	KeyCmdUseWarnNotLoggedIn: "プロファイル '%s' は未ログインです — 実行: cca login %s",
	KeyCmdUseYoloNote:        "%s⚡ 権限スキップ%s %s(--yolo → %s)%s",
	KeyCmdShSubshellBanner:   "CLAUDE_CONFIG_DIR=%s のサブシェル — exit で戻る",
	KeyCmdExecMissingCmd:     "実行するコマンドがありません — 例: cca exec work -- claude auth status",

	// Commands - rm
	KeyCmdRmCannotRmDefault: "デフォルトプロファイルは削除できません — これは ~/.claude です",
	KeyCmdRmConfirmHeading:  "プロファイル %s%s%s を削除しようとしています:",
	KeyCmdRmDirectoryLabel:  "ディレクトリ",
	KeyCmdRmAccountLabel:    "アカウント",
	KeyCmdRmHasSession:      "(ログインセッションが残っています)",
	KeyCmdRmSharedHint:      "共有リンクのみが解除されます — ~/.claude 本体は変更されません。",
	KeyCmdRmConfirmPrompt:   "削除しますか? [y/N] ",
	KeyCmdRmCancelled:       "キャンセルしました",
	KeyCmdRmRemoved:         "プロファイル '%s' を削除しました",

	// Commands - sync
	KeyCmdSyncCannotSyncDefault: "デフォルトプロファイルは同期できません — これは ~/.claude です",
	KeyCmdSyncMissingName:       "プロファイル名がありません — 例: cca sync work   (または cca sync --all)",
	KeyCmdSyncNoProfiles:        "同期対象のプロファイルがありません",
	KeyCmdSyncUpToDate:          "(既に最新です)",

	// Commands - settings
	KeyCmdSettingsTitle:          "⚙  cca 設定",
	KeyCmdSettingsLangLabel:      "表示言語",
	KeyCmdSettingsSyncLabel:      "同期戦略",
	KeyCmdSettingsSourceConfig:   "ソース: config.json",
	KeyCmdSettingsSourceEnv:      "ソース: 環境変数",
	KeyCmdSettingsSourceDefault:  "ソース: システム標準",
	KeyCmdSettingsOptionsHeading: "利用可能な操作:",
	KeyCmdSettingsOptLangMenu:    "cca settings lang          対話メニューで言語を選択",
	KeyCmdSettingsOptLangCode:    "cca settings lang <code>   言語を直接設定 (en, vi, zh, ja, es, auto)",
	KeyCmdSettingsPromptTitle:    "🌐  表示言語を選択してください",
	KeyCmdSettingsPromptChoice:   "選択肢 [1-%d] を入力 (キャンセルは 'q'): ",
	KeyCmdSettingsInvalidChoice:  "無効な選択です。1 から %d の数字を入力してください。",
	KeyCmdSettingsLangUpdated:    "表示言語を更新しました: %s",
	KeyCmdSettingsLangAutoSet:    "言語設定を自動にリセットしました (現在のシステム: %s)",
	KeyCmdSettingsLangErrInvalid: "サポートされていない言語 '%s' です。対応言語: %s、または 'auto'",

	// Commands - update
	KeyCmdUpdateChecking:        "アップデートを確認中…",
	KeyCmdUpdateAlreadyLatest:   "cca は既に最新バージョンです (%s)",
	KeyCmdUpdateAvailableNotice: "cca の新しいバージョンが利用可能です: %s → %s",
	KeyCmdUpdateRunHint:         "'cca update' を実行して更新してください。",
	KeyCmdUpdateDownloading:     "cca %s をダウンロード中…",
	KeyCmdUpdateSuccess:         "cca を %s に正常に更新しました",
	KeyCmdUpdateCheckFailed:     "アップデートの確認に失敗しました: %v",
	KeyCmdUpdateApplyFailed:     "アップデートの適用に失敗しました: %v",
	KeyCmdUpdateNoAsset:         "リリース %[3]s に %[1]s/%[2]s 用のビルド済みバイナリがありません",

	// Guide & Usage
	KeyGuideUsage: `cca — 1台のマシンで複数の Claude Code アカウントを管理

使い方: cca <コマンド> [引数...]

コマンド一覧:
  ls                    プロファイルとログイン中のアカウント一覧
  new <名前>             新規プロファイルを作成 (--login, --yolo)
  use <名前> [引数…]     指定プロファイルで Claude Code を起動
  login <名前>           プロファイルにログイン
  logout <名前>          プロファイルからログアウト
  info <名前>            プロファイル詳細を表示
  sh <名前>              プロファイル環境でサブシェルを開く
  exec <名前> -- <cmd>   プロファイル環境で任意のコマンドを実行
  rm <名前>              プロファイルを削除 (-y, --keep-keychain)
  sync [<名前>|--all]    共有ファイルを更新 (--strategy)
  doctor                 リンク、認証情報、孤立プロファイルを検査
  settings [lang]        設定と表示言語を管理
  update [--check]       アップデートの確認と自動更新
  config                 ~/.claude-accounts を VS Code で開く (--edit, --print)
  guide                  完全ガイドを表示
  install                cca を PATH に追加し補完を設定
  version                cca のバージョンを表示

使用例:
  cca new work --login       'work' プロファイルを作成してすぐログイン
  cca work                   'work' プロファイルで Claude Code を実行
  cca work --resume          追加の引数はそのまま claude に渡されます
  cca ls                     アカウントのログイン状況を確認
  cca sync --all             プラグイン追加後に共有ファイルを再同期
  cca settings lang          表示言語を変更
  cca work --yolo            --dangerously-skip-permissions の別名
  cca sh work                CLAUDE_CONFIG_DIR 設定済みのサブシェルを開く
  cca exec work -- git log   プロファイル環境でコマンドを実行

詳細: cca guide`,

	KeyGuideFull: `
%[1]scca — 1台のPCで複数の Claude Code アカウントを切り替え%[2]s

%[1]sクイックスタート%[2]s
  %[3]scca new work --login%[2]s     'work' を作成し、ブラウザでログイン
  %[3]scca ls%[2]s                   どのアカウントにログイン中か確認
  %[3]scca work%[2]s                 'work' プロファイルで Claude Code を実行

  ルート設定 (~/.claude) は常に %[3]sdefault%[2]s として利用可能です。

%[1]s仕組み%[2]s
  各プロファイルは ~/.claude-accounts/ 下の個別設定ディレクトリです。
  ログインセッションはキーチェーン (macOS) またはディレクトリ内の
  .credentials.json (Linux/Windows) に保存され、2つの端末で2つの
  アカウントを完全に並行して実行できます。

%[1]sCLAUDE の実行%[2]s
  %[3]scca work%[2]s                     通常の対話セッション
  %[3]scca work --resume%[2]s            以降の引数は直接 claude に渡されます
  %[3]scca work -p "質問"%[2]s           直接出力モード
  %[3]scca work --yolo%[2]s              --dangerously-skip-permissions の別名
  %[3]scca default --yolo%[2]s           デフォルトアカウントで権限チェックをスキップ
  %[4]s--yolo は cca 経由時のみ有効です。直接 ` + "`claude`" + ` を実行した場合は変わりません。%[2]s

%[1]sプロファイル管理%[2]s
  %[3]scca new <名前> [--login] [--yolo]%[2]s   作成; --yolo でパーミッション自動スキップ
  %[3]scca login <名前>%[2]s / %[3]scca logout <名前>%[2]s     ログイン / ログアウト
  %[3]scca info <名前>%[2]s                     ディレクトリ、認証情報、容量を確認
  %[3]scca rm <名前>%[2]s                       削除 (確認プロンプトあり)
  %[3]scca settings [lang]%[2]s                 言語設定や各種オプションを変更

%[1]s共有ファイル%[2]s
  plugins、skills、agents は ~/.claude へシンボリックリンクされます。
  1度インストールすれば全プロファイルから参照できます。Windows で
  開発者モードが無効な場合はジャンクションやコピーへ自動フォールバックします
  (参照: %[3]scca doctor%[2]s)。

  %[3]scca sync --all%[2]s      共有ファイルを更新 (プラグイン追加後に実行)
  %[3]scca config%[2]s          VS Code で ~/.claude-accounts の config.json を開く

%[1]sトラブルシューティング%[2]s
  %[3]scca doctor%[2]s   リンク切れや孤立した認証情報を検査

%[1]s関連情報%[2]s
  %[3]scca <コマンド> --help%[2]s   各コマンドの詳細
  %[3]scca version%[2]s            現在のバージョンを表示
`,

	// Profile validation errors
	KeyProfileErrReserved: "'%s' は予約名です — デフォルトプロファイルは ~/.claude です",
	KeyProfileErrInvalid:  "プロファイル名には英数字および . _ - のみ使用でき、英数字で始まる必要があります",
	KeyProfileErrNotFound: "プロファイル '%s' が見つかりません。利用可能: %s",
	KeyProfileErrCreateIt: "  作成コマンド: cca new %s",

	// General CLI errors
	KeyCliErrMissingProfileName: "プロファイル名が指定されていません — 例: cca use work",
	KeyCliErrNotCommandOrProf:   "'%s' はコマンドでも既存プロファイルでもありません。",
	KeyCliErrExistingProfiles:   "  既存のプロファイル: %s",
	KeyCliErrNoExtraProfiles:    "  追加プロファイルはありません ('default' = ~/.claude のみ)。",
	KeyCliErrCreateThisProfile:  "  プロファイル作成:  cca new %s --login",
	KeyCliErrSeeCommandsGuide:   "  コマンド一覧:     cca --help   ·   完全ガイド: cca guide",
}
