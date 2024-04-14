package app

// InitService this initializes all the busines logic services
func InitService(a *App) {
	a.Url = InitUrl(&UrlImplOpts{
		App:    a,
		Logger: a.Logger,
	})
}
