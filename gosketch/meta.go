package gosketch

// Meta meta.json
type Meta struct {
	Commit               string              `json:"commit"`
	PagesAndArtboards    map[string]MetaPage `json:"pagesAndArtboards"`
	Version              int                 `json:"version"`
	Fonts                []string            `json:"fonts"`
	CompatibilityVersion string              `json:"compatibilityVersion"`
	App                  string              `json:"app"`
	Autosaved            int64               `json:"autosaved"`
	Variant              string              `json:"variant"`
	Created              MetaCreated         `json:"created"`
	AppVersion           string              `json:"appVersion"`
	Build                int64               `json:"build"`
}

// MetaPage meta.json "pagesAndArtboards"
type MetaPage struct {
	Name      string                  `json:"name"`
	Artboards map[string]MetaArtboard `json:"artboards"`
}

// MetaArtboard meta.json "pagesAndArtboards > artboards"
type MetaArtboard struct {
	Name string
}

// MetaCreated meta.json "created"
type MetaCreated struct {
	Commit               string `json:"cpmmit"`
	AppVersion           string `json:"appVersion"`
	Build                int64  `json:"build"`
	App                  string `json:"app"`
	CompatibilityVersion int64  `json:"compatibilityVersion"`
	Version              int64  `json:"version"`
	Variant              string `json:"variant"`
}