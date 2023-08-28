package noteui

type BackMsg struct{}
type TableMsg struct{}
type CreateMsg struct{ RootDir string }
type OpenMsg struct{ TitleNote string }
type EditMsg struct{ Filepath string }
type DeleteMsg struct{ TitleNote string }
type UpdateTableMsg struct{}
