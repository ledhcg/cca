package i18n

var zhCatalog = map[Key]string{
	// Common / Generic
	KeyCommonDefault:   "默认",
	KeyCommonNone:      "无",
	KeyCommonCancelled: "已取消",
	KeyCommonOk:        "正常",
	KeyCommonPresent:   "存在",
	KeyCommonActive:    "已激活",

	// Commands - ls
	KeyCmdLsHeaderProfile: "配置文件",
	KeyCmdLsHeaderStatus:  "状态",
	KeyCmdLsHeaderAccount: "账号",
	KeyCmdLsHeaderPlan:    "方案",
	KeyCmdLsStatusIn:      "已登录",
	KeyCmdLsStatusOut:     "未登录",
	KeyCmdLsNoProfiles:    "暂无额外配置文件。创建一个: %scca new work --login%s",

	// Commands - info
	KeyCmdInfoConfigDir:       "配置目录",
	KeyCmdInfoConfigDirNotSet: "(未设置 CLAUDE_CONFIG_DIR)",
	KeyCmdInfoLoggedIn:        "已登录",
	KeyCmdInfoStatus:          "状态",
	KeyCmdInfoProjectsOpened:  "%s · 已打开 %d 个项目",
	KeyCmdInfoPrivateSize:     "私有大小",
	KeyCmdInfoShared:          "共享",

	// Commands - doctor
	KeyCmdDoctorProfilesHeading: "配置文件列表",
	KeyCmdDoctorBrokenLinks:     "失效链接: %s",
	KeyCmdDoctorMissingLinks:    "缺失链接: %s (运行: cca sync --all)",
	KeyCmdDoctorLoggedInNoCred:  "已登录但未找到 %s: %s",
	KeyCmdDoctorHasCredNoLogin:  "拥有 %s 但无法登录",
	KeyCmdDoctorOrphansHeading:  "%s 未关联任何配置文件",
	KeyCmdDoctorOrphanHint:      "(配置文件已被手动删除?)",

	// Commands - new
	KeyCmdNewMissingName:   "缺少配置文件名称 — 例如: cca new work",
	KeyCmdNewAlreadyExists: "配置文件 '%s' 已存在于 %s",
	KeyCmdNewCreated:       "已创建配置文件 %s%s%s → %s",
	KeyCmdNewLoginHint:     "登录:  %scca login %s%s",
	KeyCmdNewYoloNote:      "  %s⚡ defaultMode = bypassPermissions%s %s(仅限此配置文件)%s",

	// Commands - login / logout
	KeyCmdLoginOpeningBrowser: "正在打开浏览器登录配置文件 '%s'…",
	KeyCmdLoginStillOut:       "'%s' 仍处于未登录状态",
	KeyCmdLogoutSuccess:       "已注销配置文件 '%s'",

	// Commands - use / sh / exec
	KeyCmdUseWarnNotLoggedIn: "配置文件 '%s' 尚未登录 — 请运行: cca login %s",
	KeyCmdUseYoloNote:        "%s⚡ 跳过权限检查%s %s(--yolo → %s)%s",
	KeyCmdShSubshellBanner:   "子 Shell 环境 CLAUDE_CONFIG_DIR=%s — 输入 exit 退出",
	KeyCmdExecMissingCmd:     "缺少要运行的命令 — 例如: cca exec work -- claude auth status",

	// Commands - rm
	KeyCmdRmCannotRmDefault: "无法删除默认配置文件 — 这是 ~/.claude 根配置",
	KeyCmdRmConfirmHeading:  "准备删除配置文件 %s%s%s:",
	KeyCmdRmDirectoryLabel:  "目录",
	KeyCmdRmAccountLabel:    "账号",
	KeyCmdRmHasSession:      "(仍保存有登录会话)",
	KeyCmdRmSharedHint:      "共享链接仅会被解除关联 — ~/.claude 原始内容不受影响。",
	KeyCmdRmConfirmPrompt:   "确认删除? [y/N] ",
	KeyCmdRmCancelled:       "已取消",
	KeyCmdRmRemoved:         "已删除配置文件 '%s'",

	// Commands - sync
	KeyCmdSyncCannotSyncDefault: "无法同步默认配置文件 — 这是 ~/.claude 根配置",
	KeyCmdSyncMissingName:       "缺少配置文件名称 — 例如: cca sync work   (或 cca sync --all)",
	KeyCmdSyncNoProfiles:        "尚无要同步的配置文件",
	KeyCmdSyncUpToDate:          "(已是最新状态)",

	// Commands - settings
	KeyCmdSettingsTitle:          "⚙  cca 设置",
	KeyCmdSettingsLangLabel:      "显示语言",
	KeyCmdSettingsSyncLabel:      "同步策略",
	KeyCmdSettingsSourceConfig:   "来源: config.json",
	KeyCmdSettingsSourceEnv:      "来源: 环境变量",
	KeyCmdSettingsSourceDefault:  "来源: 系统默认",
	KeyCmdSettingsOptionsHeading: "可用选项:",
	KeyCmdSettingsOptLangMenu:    "cca settings lang          通过交互式菜单选择语言",
	KeyCmdSettingsOptLangCode:    "cca settings lang <code>   快速设置语言 (en, vi, zh, ja, es, auto)",
	KeyCmdSettingsPromptTitle:    "🌐  选择显示语言",
	KeyCmdSettingsPromptChoice:   "请输入选项 [1-%d] 或输入 'q' 取消: ",
	KeyCmdSettingsInvalidChoice:  "选项无效。请输入 1 到 %d 之间的数字。",
	KeyCmdSettingsLangUpdated:    "显示语言已更新: %s",
	KeyCmdSettingsLangAutoSet:    "语言已重置为自动 (当前系统: %s)",
	KeyCmdSettingsLangErrInvalid: "不支持的语言 '%s'。支持的语言: %s，或 'auto'",

	// Commands - update
	KeyCmdUpdateChecking:        "正在检查更新…",
	KeyCmdUpdateAlreadyLatest:   "cca 已是最新版本 (%s)",
	KeyCmdUpdateAvailableNotice: "发现 cca 新版本: %s → %s",
	KeyCmdUpdateRunHint:         "运行 'cca update' 进行更新。",
	KeyCmdUpdateDownloading:     "正在下载 cca %s…",
	KeyCmdUpdateSuccess:         "已成功将 cca 更新至 %s",
	KeyCmdUpdateCheckFailed:     "检查更新失败: %v",
	KeyCmdUpdateApplyFailed:     "应用更新失败: %v",
	KeyCmdUpdateNoAsset:         "发布版本 %[3]s 中没有适用于 %[1]s/%[2]s 的预编译二进制文件",

	// Commands - handoff
	KeyCmdHandoffMissingTarget:  "缺少目标配置文件名称 — 例如: cca handoff work",
	KeyCmdHandoffSameProfile:    "源配置文件与目标配置文件不能相同",
	KeyCmdHandoffNoSessionFound: "在配置文件 '%s' 中未找到此项目的会话",
	KeyCmdHandoffBanner:         "正在将会话 %s (%s) 从 '%s' 交接至 '%s'…",
	KeyCmdHandoffCopyFailed:     "交接会话失败",

	// Commands - session
	KeyCmdSessionUnknownSubcmd:     "未知的会话命令 '%s' — 尝试: ls, cp, mv, rm",
	KeyCmdSessionMissingFromTo:     "--from 和 --to 均为必填项 — 例如: cca session cp --from work --to personal --all",
	KeyCmdSessionSameProfile:       "--from 和 --to 不能是同一个配置文件",
	KeyCmdSessionRequireIdOrAll:    "必须指定 --id <sessionId> 或 --all",
	KeyCmdSessionNoSessionsFound:   "在配置文件 '%s' 中未找到此项目的会话",
	KeyCmdSessionNoSessionsAllHint: "(使用 --all 查看所有项目的会话)",
	KeyCmdSessionCopySuccess:       "已将会话 %s 从 '%s' 复制到 '%s'",
	KeyCmdSessionCopyAllSuccess:    "已将 %d 个会话从 '%s' 复制到 '%s'",
	KeyCmdSessionMoveSuccess:       "已将会话 %s 从 '%s' 移动到 '%s'",
	KeyCmdSessionMoveAllSuccess:    "已将 %d 个会话从 '%s' 移动到 '%s'",
	KeyCmdSessionRmSuccess:         "已从 '%s' 中删除会话 %s",
	KeyCmdSessionRmAllSuccess:      "已从 '%s' 中删除 %d 个会话",
	KeyCmdSessionRmConfirmSingle:   "即将从配置文件 '%s' 中删除会话 %s。确认吗? [y/N] ",
	KeyCmdSessionRmConfirmAll:      "即将从配置文件 '%s' 中删除此项目的 %d 个会话。确认吗? [y/N] ",
	KeyCmdSessionHeaderID:          "会话ID",
	KeyCmdSessionHeaderProject:     "项目",
	KeyCmdSessionHeaderTitle:       "标题",
	KeyCmdSessionHeaderMessages:    "消息数",
	KeyCmdSessionHeaderSize:        "大小",
	KeyCmdSessionHeaderModified:    "修改时间",

	// Guide & Usage
	KeyGuideUsage: `cca — 在同一台机器上管理多个 Claude Code 账号

用法: cca <命令> [参数...]

命令列表:
  ls                    列出所有配置文件及其登录的账号
  new <名称>             创建新的配置文件 (--login, --yolo)
  use <名称> [参数…]     在指定配置文件下运行 Claude Code
  handoff <目标> [参数…] 将当前会话交接给另一配置文件并继续 (--fork)
  session <操作>         管理会话: ls, cp, mv, rm (--from, --to, --all, --id)
  login <名称>           登录指定配置文件
  logout <名称>          注销指定配置文件
  info <名称>            查看配置文件详细信息
  sh <名称>              在配置文件的环境变量中打开子 Shell
  exec <名称> -- <命令>   在配置文件的环境中运行任意命令
  rm <名称>              删除配置文件 (-y, --keep-keychain)
  sync [<名称>|--all]    刷新共享文件 (--strategy)
  doctor                 检查软链接、凭证和孤立配置文件
  settings [lang]        管理设置与显示语言
  update [--check]       检查并自动更新 cca
  config                 在 VS Code 中打开 ~/.claude-accounts (--edit, --print)
  guide                  查看完整使用指南
  install                将 cca 添加到 PATH 并配置自动补全
  version                显示 cca 当前版本

示例:
  cca new work --login       创建 'work' 配置文件并立即登录
  cca work                   在 'work' 配置文件下运行 Claude Code
  cca work --resume          后续参数将直接透传给 claude
  cca handoff work           将会话快速交接给 'work' 并继续工作
  cca session ls             列出当前目录下的会话
  cca session cp --from default --to work --all   复制全部会话
  cca ls                     查看各个配置文件登录的账号
  cca sync --all             安装新插件后刷新共享文件
  cca settings lang          修改显示语言
  cca work --yolo            --dangerously-skip-permissions 的别名
  cca sh work                在预设好 CLAUDE_CONFIG_DIR 的子 Shell 中运行

示例:
  cca new work --login       创建 'work' 配置文件并立即登录
  cca work                   在 'work' 配置文件下运行 Claude Code
  cca work --resume          后续参数将直接透传给 claude
  cca ls                     查看各个配置文件登录的账号
  cca sync --all             安装新插件后刷新共享文件
  cca settings lang          修改显示语言
  cca work --yolo            --dangerously-skip-permissions 的别名
  cca sh work                在预设好 CLAUDE_CONFIG_DIR 的子 Shell 中运行
  cca exec work -- git log   在配置文件环境中运行命令

查看更多: cca guide`,

	KeyGuideFull: `
%[1]scca — 在同一台电脑上管理多个 Claude Code 账号%[2]s

%[1]s快速上手%[2]s
  %[3]scca new work --login%[2]s     创建配置文件 'work' 并打开浏览器登录
  %[3]scca ls%[2]s                   查看每个配置文件登录的账号
  %[3]scca work%[2]s                 在配置文件 'work' 下运行 Claude Code

  根配置文件 (~/.claude) 始终作为 %[3]sdefault%[2]s 可用 — 无需手动创建。

%[1]s工作原理%[2]s
  每个配置文件都是 ~/.claude-accounts/ 下的独立配置目录。登录会话存储在
  系统 Keychain (macOS) 或配置文件目录下的 .credentials.json 文件中 (Linux/Windows)。
  每个配置文件保持独立会话，让您可以在两个终端窗口中并排独立运行两个账号。

%[1]s运行 CLAUDE%[2]s
  %[3]scca work%[2]s                     常规交互式会话
  %[3]scca work --resume%[2]s            附加参数将直接传递给 claude
  %[3]scca work -p "问题"%[2]s           单次直接输出模式
  %[3]scca work --yolo%[2]s              --dangerously-skip-permissions 的别名
  %[3]scca default --yolo%[2]s           根账号跳过权限审批
  %[4]s--yolo 仅在 cca 内部生效；直接运行 ` + "`claude`" + ` 保持不变。%[2]s

%[1]s管理配置文件%[2]s
  %[3]scca new <名称> [--login] [--yolo]%[2]s   创建配置文件；--yolo 预设跳过权限
  %[3]scca login <名称>%[2]s / %[3]scca logout <名称>%[2]s     登录 / 注销配置文件
  %[3]scca info <名称>%[2]s                     查看目录、凭证、账号、占用大小
  %[3]scca rm <名称>%[2]s                       删除配置文件 (操作前需确认)
  %[3]scca settings [lang]%[2]s                 配置显示语言和其他选项

%[1]s共享文件%[2]s
  插件 (plugins)、技能 (skills)、智能体 (agents) 软链接回 ~/.claude — 插件仅需
  安装一次，所有配置文件即刻可用。在 Windows 未开启开发者模式时，
  cca 会自动回退为目录联接/硬链接/文件复制 (参见 %[3]scca doctor%[2]s)。

  %[3]scca sync --all%[2]s      刷新共享文件 (安装新插件后运行)
  %[3]scca config%[2]s          在 VS Code 中打开 ~/.claude-accounts 编辑 config.json

%[1]s故障排查%[2]s
  %[3]scca doctor%[2]s   检测失效链接、凭证不一致或异常配置文件。

%[1]s参见%[2]s
  %[3]scca <命令> --help%[2]s   查看单项命令帮助
  %[3]scca version%[2]s         显示当前 cca 版本
`,

	// Profile validation errors
	KeyProfileErrReserved: "'%s' 是保留名称 — 默认配置文件已固定为 ~/.claude",
	KeyProfileErrInvalid:  "配置文件名称只能包含字母、数字和 . _ -，且必须以字母或数字开头",
	KeyProfileErrNotFound: "未找到配置文件 '%s'。当前可用: %s",
	KeyProfileErrCreateIt: "  使用此命令创建: cca new %s",

	// General CLI errors
	KeyCliErrMissingProfileName: "缺少配置文件名称 — 例如: cca use work",
	KeyCliErrNotCommandOrProf:   "'%s' 既不是已知命令也不是现有配置文件。",
	KeyCliErrExistingProfiles:   "  现有配置文件: %s",
	KeyCliErrNoExtraProfiles:    "  暂无额外配置文件 (仅有默认 'default' = ~/.claude)。",
	KeyCliErrCreateThisProfile:  "  创建此配置文件:  cca new %s --login",
	KeyCliErrSeeCommandsGuide:   "  查看命令帮助:   cca --help   ·   完整使用指南: cca guide",
}
