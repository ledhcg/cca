package i18n

var esCatalog = map[Key]string{
	// Common / Generic
	KeyCommonDefault:   "predeterminado",
	KeyCommonNone:      "ninguno",
	KeyCommonCancelled: "Cancelado",
	KeyCommonOk:        "correcto",
	KeyCommonPresent:   "presente",
	KeyCommonActive:    "Activo",

	// Commands - ls
	KeyCmdLsHeaderProfile: "PERFIL",
	KeyCmdLsHeaderStatus:  "ESTADO",
	KeyCmdLsHeaderAccount: "CUENTA",
	KeyCmdLsHeaderPlan:    "PLAN",
	KeyCmdLsStatusIn:      "sesión iniciada",
	KeyCmdLsStatusOut:     "sesión cerrada",
	KeyCmdLsNoProfiles:    "Aún no hay perfiles adicionales. Cree uno: %scca new work --login%s",

	// Commands - info
	KeyCmdInfoConfigDir:       "Dir de configuración",
	KeyCmdInfoConfigDirNotSet: "(CLAUDE_CONFIG_DIR no establecido)",
	KeyCmdInfoLoggedIn:        "Sesión iniciada",
	KeyCmdInfoStatus:          "Estado",
	KeyCmdInfoProjectsOpened:  "%s · %d proyecto(s) abierto(s)",
	KeyCmdInfoPrivateSize:     "Tamaño privado",
	KeyCmdInfoShared:          "Compartidos",

	// Commands - doctor
	KeyCmdDoctorProfilesHeading: "Perfiles",
	KeyCmdDoctorBrokenLinks:     "enlace(s) roto(s): %s",
	KeyCmdDoctorMissingLinks:    "enlace(s) faltante(s): %s (ejecute cca sync --all)",
	KeyCmdDoctorLoggedInNoCred:  "sesión iniciada pero no se encontró %s: %s",
	KeyCmdDoctorHasCredNoLogin:  "tiene un %s pero no se pudo iniciar sesión",
	KeyCmdDoctorOrphansHeading:  "%s no pertenece a ningún perfil",
	KeyCmdDoctorOrphanHint:      "(¿perfil eliminado manualmente?)",

	// Commands - new
	KeyCmdNewMissingName:   "falta el nombre del perfil — ejemplo: cca new work",
	KeyCmdNewAlreadyExists: "el perfil '%s' ya existe en %s",
	KeyCmdNewCreated:       "Perfil creado %s%s%s → %s",
	KeyCmdNewLoginHint:     "Iniciar sesión:  %scca login %s%s",
	KeyCmdNewYoloNote:      "  %s⚡ defaultMode = bypassPermissions%s %s(solo este perfil)%s",

	// Commands - login / logout
	KeyCmdLoginOpeningBrowser: "Abriendo el navegador para iniciar sesión en el perfil '%s'…",
	KeyCmdLoginStillOut:       "'%s' aún no ha iniciado sesión",
	KeyCmdLogoutSuccess:       "Sesión cerrada en el perfil '%s'",

	// Commands - use / sh / exec
	KeyCmdUseWarnNotLoggedIn: "El perfil '%s' no ha iniciado sesión — ejecute: cca login %s",
	KeyCmdUseYoloNote:        "%s⚡ omitir permisos%s %s(--yolo → %s)%s",
	KeyCmdShSubshellBanner:   "Subshell con CLAUDE_CONFIG_DIR=%s — escriba exit para salir",
	KeyCmdExecMissingCmd:     "falta el comando a ejecutar — ejemplo: cca exec work -- claude auth status",

	// Commands - rm
	KeyCmdRmCannotRmDefault: "no se puede eliminar el perfil predeterminado — es ~/.claude",
	KeyCmdRmConfirmHeading:  "A punto de eliminar el perfil %s%s%s:",
	KeyCmdRmDirectoryLabel:  "directorio",
	KeyCmdRmAccountLabel:    "cuenta",
	KeyCmdRmHasSession:      "(aún mantiene una sesión activa)",
	KeyCmdRmSharedHint:      "Los enlaces compartidos solo se desenlazan — ~/.claude permanece intacto.",
	KeyCmdRmConfirmPrompt:   "¿Confirmar? [y/N] ",
	KeyCmdRmCancelled:       "Cancelado",
	KeyCmdRmRemoved:         "Perfil '%s' eliminado",

	// Commands - sync
	KeyCmdSyncCannotSyncDefault: "no se puede sincronizar el perfil predeterminado — es ~/.claude",
	KeyCmdSyncMissingName:       "falta el nombre del perfil — ejemplo: cca sync work   (o cca sync --all)",
	KeyCmdSyncNoProfiles:        "aún no hay perfiles para sincronizar",
	KeyCmdSyncUpToDate:          "(ya está actualizado)",

	// Commands - settings
	KeyCmdSettingsTitle:          "⚙  Configuración de cca",
	KeyCmdSettingsLangLabel:      "Idioma de visualización",
	KeyCmdSettingsSyncLabel:      "Estrategia de sincronización",
	KeyCmdSettingsSourceConfig:   "origen: config.json",
	KeyCmdSettingsSourceEnv:      "origen: variable de entorno",
	KeyCmdSettingsSourceDefault:  "origen: predeterminado del sistema",
	KeyCmdSettingsOptionsHeading: "Opciones disponibles:",
	KeyCmdSettingsOptLangMenu:    "cca settings lang          Seleccionar idioma mediante menú interactivo",
	KeyCmdSettingsOptLangCode:    "cca settings lang <code>   Cambiar idioma rápidamente (en, vi, zh, ja, es, auto)",
	KeyCmdSettingsPromptTitle:    "🌐  Seleccione el idioma de visualización",
	KeyCmdSettingsPromptChoice:   "Seleccione [1-%d] o 'q' para cancelar: ",
	KeyCmdSettingsInvalidChoice:  "Opción no válida. Elija un número entre 1 y %d.",
	KeyCmdSettingsLangUpdated:    "Idioma de visualización actualizado: %s",
	KeyCmdSettingsLangAutoSet:    "Idioma restablecido a automático (sistema actual: %s)",
	KeyCmdSettingsLangErrInvalid: "Idioma no compatible '%s'. Compatibles: %s, o 'auto'",

	// Guide & Usage
	KeyGuideUsage: `cca — varias cuentas de Claude Code en una sola máquina

Uso: cca <comando> [argumentos...]

Comandos:
  ls                    listar perfiles y en qué cuenta ha iniciado sesión cada uno
  new <nombre>           crear un nuevo perfil (--login, --yolo)
  use <nombre> [args…]   ejecutar Claude Code bajo un perfil
  login <nombre>         iniciar sesión en un perfil
  logout <nombre>        cerrar sesión en un perfil
  info <nombre>          mostrar detalles de un perfil
  sh <nombre>            abrir una subshell con el entorno del perfil
  exec <nombre> -- <cmd> ejecutar un comando arbitrario en el entorno del perfil
  rm <nombre>            eliminar un perfil (-y, --keep-keychain)
  sync [<nombre>|--all]  actualizar archivos compartidos (--strategy)
  doctor                 verificar enlaces, credenciales y perfiles huérfanos
  settings [lang]        administrar configuración e idioma de visualización
  config                 abrir ~/.claude-accounts en VS Code (--edit, --print)
  guide                  guía completa de uso
  install                agregar cca a PATH y configurar autocompletado
  version                imprimir la versión de cca

Ejemplos:
  cca new work --login       crear el perfil 'work' e iniciar sesión de inmediato
  cca work                   ejecutar Claude Code bajo el perfil 'work'
  cca work --resume          los argumentos posteriores se pasan directamente a claude
  cca ls                     ver qué perfil está conectado a qué cuenta
  cca sync --all             actualizar archivos compartidos tras instalar un plugin
  cca settings lang          cambiar el idioma de la interfaz
  cca work --yolo            alias de --dangerously-skip-permissions
  cca sh work                abrir una subshell con CLAUDE_CONFIG_DIR establecido
  cca exec work -- git log   ejecutar cualquier comando en el entorno del perfil

Vea también: cca guide`,

	KeyGuideFull: `
%[1]scca — varias cuentas de Claude Code en una máquina%[2]s

%[1]sINICIO RÁPIDO%[2]s
  %[3]scca new work --login%[2]s     crear el perfil 'work' e iniciar sesión en el navegador
  %[3]scca ls%[2]s                   ver qué perfil está conectado a qué cuenta
  %[3]scca work%[2]s                 ejecutar Claude Code bajo el perfil 'work'

  El perfil raíz (~/.claude) siempre está disponible como %[3]sdefault%[2]s — no es necesario crearlo.

%[1]sCÓMO FUNCIONA%[2]s
  Cada perfil es su propio directorio de configuración bajo ~/.claude-accounts/. La sesión
  reside en el Llavero (macOS) o en un archivo .credentials.json dentro del propio
  directorio del perfil (Linux/Windows) — cada perfil mantiene una sesión independiente,
  permitiendo trabajar en dos terminales en paralelo.

%[1]sEJECUTAR CLAUDE%[2]s
  %[3]scca work%[2]s                     sesión interactiva normal
  %[3]scca work --resume%[2]s            los argumentos adicionales se pasan directamente a claude
  %[3]scca work -p "pregunta"%[2]s       modo de salida directa (print mode)
  %[3]scca work --yolo%[2]s              alias de --dangerously-skip-permissions
  %[3]scca default --yolo%[2]s           cuenta predeterminada sin comprobación de permisos
  %[4]s--yolo solo surte efecto en cca; ejecutar ` + "`claude`" + ` directamente no se altera.%[2]s

%[1]sADMINISTRAR PERFILES%[2]s
  %[3]scca new <nombre> [--login] [--yolo]%[2]s   crear; --yolo preestablece bypassPermissions
  %[3]scca login <nombre>%[2]s / %[3]scca logout <nombre>%[2]s     iniciar / cerrar sesión
  %[3]scca info <nombre>%[2]s                     directorio, credencial, cuenta, tamaño
  %[3]scca rm <nombre>%[2]s                       eliminar un perfil (solicita confirmación)
  %[3]scca settings [lang]%[2]s                   configurar idioma y opciones

%[1]sARCHIVOS COMPARTIDOS%[2]s
  plugins, skills y agents se enlazan con ~/.claude — instale un plugin una vez
  y estará disponible para todos los perfiles. En Windows sin Modo Desarrollador,
  cca recurre a uniones/hard links/copias (consulte %[3]scca doctor%[2]s).

  %[3]scca sync --all%[2]s      actualizar archivos compartidos (tras instalar un plugin)
  %[3]scca config%[2]s          abrir ~/.claude-accounts en VS Code para editar config.json

%[1]sRESOLUCIÓN DE PROBLEMAS%[2]s
  %[3]scca doctor%[2]s   enlaces rotos, credenciales huérfanas o perfiles desincronizados.

%[1]sVEA TAMBIÉN%[2]s
  %[3]scca <comando> --help%[2]s   detalles de cada comando
  %[3]scca version%[2]s            imprimir la versión de cca
`,

	// Profile validation errors
	KeyProfileErrReserved: "'%s' es un nombre reservado — el perfil predeterminado ya es ~/.claude",
	KeyProfileErrInvalid:  "los nombres de perfil solo pueden contener letras, dígitos y . _ -, y deben comenzar con letra o dígito",
	KeyProfileErrNotFound: "no existe el perfil '%s'. Disponibles: %s",
	KeyProfileErrCreateIt: "  Créelo con: cca new %s",

	// General CLI errors
	KeyCliErrMissingProfileName: "falta el nombre del perfil — ejemplo: cca use work",
	KeyCliErrNotCommandOrProf:   "'%s' no es un comando ni un perfil existente.",
	KeyCliErrExistingProfiles:   "  Perfiles existentes: %s",
	KeyCliErrNoExtraProfiles:    "  No hay perfiles adicionales (solo 'default' = ~/.claude).",
	KeyCliErrCreateThisProfile:  "  Cree este perfil:  cca new %s --login",
	KeyCliErrSeeCommandsGuide:   "  Ver comandos:     cca --help   ·   Guía completa: cca guide",
}
