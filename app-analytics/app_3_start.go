package appanalytics

func (a *App) Start() error {
	return a.transportTCP.Start()
}
