package handler

type App struct {
}

func Init() *App {
	return &App{}
}

type GroupIds struct {
	Ids []int64 `json:"ids" param:"ids" query:"ids" form:"ids"`
}

type SingleId struct {
	Id int64 `json:"id" param:"id" query:"id" form:"id"`
}
