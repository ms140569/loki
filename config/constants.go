package config

// Systemwide constants used all over the place.
const (
	InnerMagic        = "LOKI" // Fixed magic value provided in any record to verify successful decryption (with hash)
	BinaryName        = "loki"
	DefaultDirectory  = ".loki"
	FileSuffix        = ".loki"
	ConfigFilename    = ".config"
	MasterFilename    = ".master"
	LokiBaseEnv       = "LOKI_BASE"
	LokiEditorEnv     = "EDITOR"
	LokiLoglevelEnv   = "LOKI_LOGLEVEL"
	CommunicationFile = "/tmp/loki-%d.sock"
	ConfigTemplate    = "configfile.tmpl"
	ConfigTemplateGit = "configfile-git.tmpl"
	RequestMagic      = "req"
	ShutdownMagic     = "shutdown"
	AgentLogfile      = "/tmp/loki-apentd.log"
	KeyLength         = 32

	MagicLabel    = "Magic       : "
	MD5Label      = "MD5         : "
	TitleLabel    = "Title       : "
	AccountLabel  = "Account     : "
	PasswordLabel = "Password    : "
	TagsLabel     = "Tags        : "
	URLLabel      = "Url         : "
	NotesLabel    = "Notes       : "

	ExitCodeOK      = 0
	ExitCodeFailure = 1
)
