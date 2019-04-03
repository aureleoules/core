package constants

type Module string

var Modules []Module = []Module{
	Galleries,
	Projects,
	Articles,
	Videos,
	Music,
}

const (
	Galleries Module = "galleries"
	Projects  Module = "projects"
	Articles  Module = "articles"
	Videos    Module = "videos"
	Music     Module = "music"
)
