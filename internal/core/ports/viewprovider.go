package ports

type ViewProvider interface {
	ShowPrefsPage()
	ShowMainPage()
	Provider
}
