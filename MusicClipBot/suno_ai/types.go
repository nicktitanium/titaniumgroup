package suno_ai

type SongData struct {
	SongData []RequestData `json:"data"`
}

type RequestData struct {
	Title    string  `json:"title"`
	Lyrics   string  `json:"lyric"`
	AudioURL string  `json:"audio_url"`
	Duration float64 `json:"duration"`
}

type AceDataConf struct {
	Data AceDataToken `yaml:"acedata"`
}

type AceDataToken struct {
	Token string `yaml:"token"`
}
