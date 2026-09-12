package i18n

var viCatalog = map[Key]string{
	// Common / Generic
	KeyCommonDefault:   "mặc định",
	KeyCommonNone:      "không có",
	KeyCommonCancelled: "Đã hủy",
	KeyCommonOk:        "tốt",
	KeyCommonPresent:   "đã có",
	KeyCommonActive:    "Đang kích hoạt",

	// Commands - ls
	KeyCmdLsHeaderProfile: "HỒ SƠ",
	KeyCmdLsHeaderStatus:  "TRẠNG THÁI",
	KeyCmdLsHeaderAccount: "TÀI KHOẢN",
	KeyCmdLsHeaderPlan:    "GÓI CƯỚC",
	KeyCmdLsStatusIn:      "đã đăng nhập",
	KeyCmdLsStatusOut:     "chưa đăng nhập",
	KeyCmdLsNoProfiles:    "Chưa có thêm hồ sơ nào. Tạo hồ sơ mới: %scca new work --login%s",

	// Commands - info
	KeyCmdInfoConfigDir:       "Thư mục cấu hình",
	KeyCmdInfoConfigDirNotSet: "(chưa thiết lập CLAUDE_CONFIG_DIR)",
	KeyCmdInfoLoggedIn:        "Đã đăng nhập",
	KeyCmdInfoStatus:          "Trạng thái",
	KeyCmdInfoProjectsOpened:  "%s · %d dự án đã mở",
	KeyCmdInfoPrivateSize:     "Dung lượng riêng",
	KeyCmdInfoShared:          "Dùng chung",

	// Commands - doctor
	KeyCmdDoctorProfilesHeading: "Danh sách hồ sơ",
	KeyCmdDoctorBrokenLinks:     "liên kết bị hỏng: %s",
	KeyCmdDoctorMissingLinks:    "thiếu liên kết: %s (chạy: cca sync --all)",
	KeyCmdDoctorLoggedInNoCred:  "đã đăng nhập nhưng không tìm thấy %s: %s",
	KeyCmdDoctorHasCredNoLogin:  "có %s nhưng không thể đăng nhập",
	KeyCmdDoctorOrphansHeading:  "%s không thuộc về hồ sơ nào",
	KeyCmdDoctorOrphanHint:      "(hồ sơ đã bị xóa thủ công?)",

	// Commands - new
	KeyCmdNewMissingName:   "thiếu tên hồ sơ — ví dụ: cca new work",
	KeyCmdNewAlreadyExists: "hồ sơ '%s' đã tồn tại tại %s",
	KeyCmdNewCreated:       "Đã tạo hồ sơ %s%s%s → %s",
	KeyCmdNewLoginHint:     "Đăng nhập:  %scca login %s%s",
	KeyCmdNewYoloNote:      "  %s⚡ defaultMode = bypassPermissions%s %s(chỉ áp dụng cho hồ sơ này)%s",

	// Commands - login / logout
	KeyCmdLoginOpeningBrowser: "Đang mở trình duyệt để đăng nhập hồ sơ '%s'…",
	KeyCmdLoginStillOut:       "'%s' vẫn chưa được đăng nhập",
	KeyCmdLogoutSuccess:       "Đã đăng xuất hồ sơ '%s'",

	// Commands - use / sh / exec
	KeyCmdUseWarnNotLoggedIn: "Hồ sơ '%s' chưa đăng nhập — hãy chạy: cca login %s",
	KeyCmdUseYoloNote:        "%s⚡ bỏ qua quyền phê duyệt%s %s(--yolo → %s)%s",
	KeyCmdShSubshellBanner:   "Subshell với CLAUDE_CONFIG_DIR=%s — gõ exit để thoát",
	KeyCmdExecMissingCmd:     "thiếu lệnh cần chạy — ví dụ: cca exec work -- claude auth status",

	// Commands - rm
	KeyCmdRmCannotRmDefault: "không thể xóa hồ sơ mặc định — đó chính là ~/.claude",
	KeyCmdRmConfirmHeading:  "Chuẩn bị xóa hồ sơ %s%s%s:",
	KeyCmdRmDirectoryLabel:  "thư mục",
	KeyCmdRmAccountLabel:    "tài khoản",
	KeyCmdRmHasSession:      "(vẫn đang lưu phiên đăng nhập)",
	KeyCmdRmSharedHint:      "Các liên kết dùng chung chỉ bị gỡ liên kết — bản gốc ~/.claude được giữ nguyên.",
	KeyCmdRmConfirmPrompt:   "Xác nhận? [y/N] ",
	KeyCmdRmCancelled:       "Đã hủy",
	KeyCmdRmRemoved:         "Đã xóa hồ sơ '%s'",

	// Commands - sync
	KeyCmdSyncCannotSyncDefault: "không thể đồng bộ hồ sơ mặc định — đó chính là ~/.claude",
	KeyCmdSyncMissingName:       "thiếu tên hồ sơ — ví dụ: cca sync work   (hoặc cca sync --all)",
	KeyCmdSyncNoProfiles:        "chưa có hồ sơ nào để đồng bộ",
	KeyCmdSyncUpToDate:          "(đã cập nhật mới nhất)",

	// Commands - settings
	KeyCmdSettingsTitle:          "⚙  Cài đặt cca",
	KeyCmdSettingsLangLabel:      "Ngôn ngữ hiển thị",
	KeyCmdSettingsSyncLabel:      "Chiến lược đồng bộ",
	KeyCmdSettingsSourceConfig:   "nguồn: config.json",
	KeyCmdSettingsSourceEnv:      "nguồn: biến môi trường",
	KeyCmdSettingsSourceDefault:  "nguồn: mặc định hệ thống",
	KeyCmdSettingsOptionsHeading: "Tùy chọn khả dụng:",
	KeyCmdSettingsOptLangMenu:    "cca settings lang          Chọn ngôn ngữ qua menu tương tác",
	KeyCmdSettingsOptLangCode:    "cca settings lang <code>   Đổi nhanh ngôn ngữ (en, vi, zh, ja, es, auto)",
	KeyCmdSettingsPromptTitle:    "🌐  Chọn ngôn ngữ hiển thị",
	KeyCmdSettingsPromptChoice:   "Nhập lựa chọn [1-%d] hoặc 'q' để hủy: ",
	KeyCmdSettingsInvalidChoice:  "Lựa chọn không hợp lệ. Vui lòng chọn từ 1 đến %d.",
	KeyCmdSettingsLangUpdated:    "Đã cập nhật ngôn ngữ hiển thị: %s",
	KeyCmdSettingsLangAutoSet:    "Đã đặt lại ngôn ngữ tự động (hệ thống hiện tại: %s)",
	KeyCmdSettingsLangErrInvalid: "Ngôn ngữ '%s' không được hỗ trợ. Các ngôn ngữ hỗ trợ: %s, hoặc 'auto'",

	// Guide & Usage
	KeyGuideUsage: `cca — quản lý nhiều tài khoản Claude Code trên một máy

Sử dụng: cca <lệnh> [tham số...]

Các lệnh:
  ls                    liệt kê các hồ sơ và trạng thái đăng nhập
  new <tên>              tạo một hồ sơ mới (--login, --yolo)
  use <tên> [tham số…]   chạy Claude Code dưới hồ sơ chỉ định
  login <tên>            đăng nhập cho một hồ sơ
  logout <tên>           đăng xuất một hồ sơ
  info <tên>             xem chi tiết thông tin hồ sơ
  sh <tên>               mở subshell trong môi trường của hồ sơ
  exec <tên> -- <lệnh>   chạy một lệnh bất kỳ trong môi trường của hồ sơ
  rm <tên>               xóa một hồ sơ (-y, --keep-keychain)
  sync [<tên>|--all]     làm mới các tệp dùng chung (--strategy)
  doctor                 kiểm tra liên kết, credentials, hồ sơ mồ côi
  settings [lang]        quản lý cài đặt & ngôn ngữ hiển thị
  config                 mở thư mục ~/.claude-accounts trong VS Code (--edit, --print)
  guide                  xem cẩm nang hướng dẫn đầy đủ
  install                thêm cca vào PATH + cài đặt completion
  version                in ra phiên bản cca hiện tại

Ví dụ:
  cca new work --login       tạo hồ sơ 'work' và đăng nhập ngay
  cca work                   chạy Claude Code với hồ sơ 'work'
  cca work --resume          truyền các tham số còn lại thẳng vào claude
  cca ls                     xem hồ sơ nào đang đăng nhập tài khoản nào
  cca sync --all             đồng bộ lại plugin/skill sau khi cài plugin mới
  cca settings lang          thay đổi ngôn ngữ hiển thị
  cca work --yolo            bí danh cho --dangerously-skip-permissions
  cca sh work                mở terminal con với CLAUDE_CONFIG_DIR được đặt sẵn
  cca exec work -- git log   chạy lệnh trong môi trường của hồ sơ

Xem thêm: cca guide`,

	KeyGuideFull: `
%[1]scca — quản lý nhiều tài khoản Claude Code trên một máy tính%[2]s

%[1]sKHỞI ĐỘNG NHANH%[2]s
  %[3]scca new work --login%[2]s     tạo hồ sơ 'work', mở trình duyệt để đăng nhập
  %[3]scca ls%[2]s                   xem hồ sơ nào đang đăng nhập tài khoản nào
  %[3]scca work%[2]s                 chạy Claude Code với hồ sơ 'work'

  Hồ sơ gốc (~/.claude) luôn có sẵn dưới tên %[3]sdefault%[2]s — không cần tạo thêm.

%[1]sCÁCH THỨC HOẠT ĐỘNG%[2]s
  Mỗi hồ sơ là một thư mục cấu hình riêng tại ~/.claude-accounts/. Phiên đăng nhập
  được lưu trong Keychain (macOS) hoặc file .credentials.json bên trong thư mục
  hồ sơ (Linux/Windows) — mỗi hồ sơ giữ một phiên độc lập, giúp bạn chạy song song
  hai tài khoản trong hai cửa sổ terminal hoàn toàn tách biệt.

%[1]sCHẠY CLAUDE CODE%[2]s
  %[3]scca work%[2]s                     phiên làm việc tương tác thông thường
  %[3]scca work --resume%[2]s            các tham số phía sau được truyền thẳng tới claude
  %[3]scca work -p "câu hỏi"%[2]s        chế độ in trực tiếp (print mode)
  %[3]scca work --yolo%[2]s              bí danh của --dangerously-skip-permissions
  %[3]scca default --yolo%[2]s           chạy tài khoản mặc định, bỏ qua phê duyệt quyền
  %[4]s--yolo chỉ có tác dụng khi chạy qua cca; lệnh ` + "`claude`" + ` trực tiếp không đổi.%[2]s

%[1]sQUẢN LÝ HỒ SƠ%[2]s
  %[3]scca new <tên> [--login] [--yolo]%[2]s   tạo hồ sơ; --yolo bật sẵn bypassPermissions
  %[3]scca login <tên>%[2]s / %[3]scca logout <tên>%[2]s     đăng nhập / đăng xuất hồ sơ
  %[3]scca info <tên>%[2]s                     xem thư mục, credential, tài khoản, dung lượng
  %[3]scca rm <tên>%[2]s                       xóa hồ sơ (sẽ hỏi xác nhận trước khi xóa)
  %[3]scca settings [lang]%[2]s                cài đặt ngôn ngữ hiển thị và tùy chọn

%[1]sTỆP DÙNG CHUNG (SHARED FILES)%[2]s
  plugins, skills, agents được tạo liên kết ngược về ~/.claude — bạn chỉ cần cài đặt
  plugin một lần và mọi hồ sơ đều thấy. Trên Windows nếu chưa bật Developer Mode,
  cca sẽ tự động chuyển sang junction/hard link/copy (xem %[3]scca doctor%[2]s).

  %[3]scca sync --all%[2]s      làm mới tệp dùng chung (chạy sau khi cài plugin mới)
  %[3]scca config%[2]s          mở ~/.claude-accounts trong VS Code để xem config.json

%[1]sCHẨN ĐOÁN LỖI%[2]s
  %[3]scca doctor%[2]s   phát hiện liên kết hỏng, hồ sơ thiếu credential, hoặc trôi lệch.

%[1]sXEM THÊM%[2]s
  %[3]scca <lệnh> --help%[2]s   chi tiết từng lệnh
  %[3]scca version%[2]s         in phiên bản hiện tại của cca
`,

	// Profile validation errors
	KeyProfileErrReserved: "'%s' là tên dành riêng — hồ sơ mặc định vốn là ~/.claude",
	KeyProfileErrInvalid:  "tên hồ sơ chỉ được chứa chữ cái, chữ số, và . _ - đồng thời phải bắt đầu bằng chữ cái hoặc chữ số",
	KeyProfileErrNotFound: "không tìm thấy hồ sơ '%s'. Hiện có: %s",
	KeyProfileErrCreateIt: "  Tạo hồ sơ mới bằng: cca new %s",

	// General CLI errors
	KeyCliErrMissingProfileName: "thiếu tên hồ sơ — ví dụ: cca use work",
	KeyCliErrNotCommandOrProf:   "'%s' không phải là lệnh và cũng không phải là hồ sơ hiện có.",
	KeyCliErrExistingProfiles:   "  Các hồ sơ hiện có: %s",
	KeyCliErrNoExtraProfiles:    "  Chưa có thêm hồ sơ nào (chỉ có 'default' = ~/.claude).",
	KeyCliErrCreateThisProfile:  "  Tạo hồ sơ này với:  cca new %s --login",
	KeyCliErrSeeCommandsGuide:   "  Xem danh sách lệnh: cca --help   ·   Hướng dẫn đầy đủ: cca guide",
}
