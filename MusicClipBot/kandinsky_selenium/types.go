package kandinsky_selenium

type KandinskyData struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

type KandinskyConfig struct {
	Data KandinskyData `yaml:"kandinsky"`
}

type VideoParametrs struct {
	UserDirPath         string
	CurrerntGPTRespID   int
	Prompts             [][]string
	Lyrics              string
	VideoQnty           int
	IsAnimated          bool
	SongNumber          int
	FusionBrainEmail    string
	FusionBrainPassword string
}
