package flows

// ProjectLanguageGetter fetches the project language for a channel (e.g. from Flows API).
type ProjectLanguageGetter interface {
	GetProjectLanguage(channelUUID string) (string, error)
}
