package i18n

// Key is a typed string identifier for translation catalog entries.
type Key = string

const (
	// Common / Generic
	KeyCommonDefault   Key = "common.default"
	KeyCommonNone      Key = "common.none"
	KeyCommonCancelled Key = "common.cancelled"
	KeyCommonOk        Key = "common.ok"
	KeyCommonPresent   Key = "common.present"
	KeyCommonActive    Key = "common.active"

	// Commands - ls
	KeyCmdLsHeaderProfile Key = "cmd.ls.header.profile"
	KeyCmdLsHeaderStatus  Key = "cmd.ls.header.status"
	KeyCmdLsHeaderAccount Key = "cmd.ls.header.account"
	KeyCmdLsHeaderPlan    Key = "cmd.ls.header.plan"
	KeyCmdLsStatusIn      Key = "cmd.ls.status.logged_in"
	KeyCmdLsStatusOut     Key = "cmd.ls.status.logged_out"
	KeyCmdLsNoProfiles    Key = "cmd.ls.no_profiles_hint"

	// Commands - info
	KeyCmdInfoConfigDir       Key = "cmd.info.config_dir"
	KeyCmdInfoConfigDirNotSet Key = "cmd.info.config_dir_not_set"
	KeyCmdInfoLoggedIn        Key = "cmd.info.logged_in"
	KeyCmdInfoStatus          Key = "cmd.info.status"
	KeyCmdInfoProjectsOpened  Key = "cmd.info.projects_opened"
	KeyCmdInfoPrivateSize     Key = "cmd.info.private_size"
	KeyCmdInfoShared          Key = "cmd.info.shared"

	// Commands - doctor
	KeyCmdDoctorProfilesHeading Key = "cmd.doctor.profiles_heading"
	KeyCmdDoctorBrokenLinks     Key = "cmd.doctor.broken_links"
	KeyCmdDoctorMissingLinks    Key = "cmd.doctor.missing_links"
	KeyCmdDoctorLoggedInNoCred  Key = "cmd.doctor.logged_in_no_cred"
	KeyCmdDoctorHasCredNoLogin  Key = "cmd.doctor.has_cred_no_login"
	KeyCmdDoctorOrphansHeading  Key = "cmd.doctor.orphans_heading"
	KeyCmdDoctorOrphanHint      Key = "cmd.doctor.orphan_hint"

	// Commands - new
	KeyCmdNewMissingName   Key = "cmd.new.missing_name"
	KeyCmdNewAlreadyExists Key = "cmd.new.already_exists"
	KeyCmdNewCreated       Key = "cmd.new.created"
	KeyCmdNewLoginHint     Key = "cmd.new.login_hint"
	KeyCmdNewYoloNote      Key = "cmd.new.yolo_note"

	// Commands - login / logout
	KeyCmdLoginOpeningBrowser Key = "cmd.login.opening_browser"
	KeyCmdLoginStillOut       Key = "cmd.login.still_logged_out"
	KeyCmdLogoutSuccess       Key = "cmd.logout.success"

	// Commands - use / sh / exec
	KeyCmdUseWarnNotLoggedIn Key = "cmd.use.warn_not_logged_in"
	KeyCmdUseYoloNote        Key = "cmd.use.yolo_note"
	KeyCmdShSubshellBanner   Key = "cmd.sh.subshell_banner"
	KeyCmdExecMissingCmd     Key = "cmd.exec.missing_cmd"

	// Commands - rm
	KeyCmdRmCannotRmDefault Key = "cmd.rm.cannot_rm_default"
	KeyCmdRmConfirmHeading  Key = "cmd.rm.confirm_heading"
	KeyCmdRmDirectoryLabel  Key = "cmd.rm.directory_label"
	KeyCmdRmAccountLabel    Key = "cmd.rm.account_label"
	KeyCmdRmHasSession      Key = "cmd.rm.has_session"
	KeyCmdRmSharedHint      Key = "cmd.rm.shared_hint"
	KeyCmdRmConfirmPrompt   Key = "cmd.rm.confirm_prompt"
	KeyCmdRmCancelled       Key = "cmd.rm.cancelled"
	KeyCmdRmRemoved         Key = "cmd.rm.removed"

	// Commands - sync
	KeyCmdSyncCannotSyncDefault Key = "cmd.sync.cannot_sync_default"
	KeyCmdSyncMissingName       Key = "cmd.sync.missing_name"
	KeyCmdSyncNoProfiles        Key = "cmd.sync.no_profiles"
	KeyCmdSyncUpToDate          Key = "cmd.sync.up_to_date"

	// Commands - settings
	KeyCmdSettingsTitle          Key = "cmd.settings.title"
	KeyCmdSettingsLangLabel      Key = "cmd.settings.lang_label"
	KeyCmdSettingsSyncLabel      Key = "cmd.settings.sync_label"
	KeyCmdSettingsSourceConfig   Key = "cmd.settings.source_config"
	KeyCmdSettingsSourceEnv      Key = "cmd.settings.source_env"
	KeyCmdSettingsSourceDefault  Key = "cmd.settings.source_default"
	KeyCmdSettingsOptionsHeading Key = "cmd.settings.options_heading"
	KeyCmdSettingsOptLangMenu    Key = "cmd.settings.opt_lang_menu"
	KeyCmdSettingsOptLangCode    Key = "cmd.settings.opt_lang_code"
	KeyCmdSettingsPromptTitle    Key = "cmd.settings.prompt_title"
	KeyCmdSettingsPromptChoice   Key = "cmd.settings.prompt_choice"
	KeyCmdSettingsInvalidChoice  Key = "cmd.settings.invalid_choice"
	KeyCmdSettingsLangUpdated    Key = "cmd.settings.lang_updated"
	KeyCmdSettingsLangAutoSet    Key = "cmd.settings.lang_auto_set"
	KeyCmdSettingsLangErrInvalid Key = "cmd.settings.lang_err_invalid"

	// Guide & Usage
	KeyGuideUsage Key = "guide.usage"
	KeyGuideFull  Key = "guide.full"

	// Profile validation errors
	KeyProfileErrReserved Key = "profile.err.reserved"
	KeyProfileErrInvalid  Key = "profile.err.invalid"
	KeyProfileErrNotFound Key = "profile.err.not_found"
	KeyProfileErrCreateIt Key = "profile.err.create_it"

	// General CLI errors
	KeyCliErrMissingProfileName Key = "cli.err.missing_profile_name"
	KeyCliErrNotCommandOrProf   Key = "cli.err.not_command_or_prof"
	KeyCliErrExistingProfiles   Key = "cli.err.existing_profiles"
	KeyCliErrNoExtraProfiles    Key = "cli.err.no_extra_profiles"
	KeyCliErrCreateThisProfile  Key = "cli.err.create_this_profile"
	KeyCliErrSeeCommandsGuide   Key = "cli.err.see_commands_guide"
)
