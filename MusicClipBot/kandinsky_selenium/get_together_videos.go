package kandinsky_selenium

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

const (
	clipsDir    = "clips"
	songDir     = "songs"
	outputVideo = "videos.mp4"
	finalVideo  = "clip.mp4"
	videoList   = "videos.txt"
)

func (vp *VideoParametrs) GetTogether() error {
	// Получение списка видеофайлов
	time.Sleep(5 * time.Second)
	videoFiles, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", vp.UserDirPath, clipsDir))
	if err != nil {
		fmt.Print("Error reading user's clips dir: ", err)
		return err
	}

	// Проверка на наличие видеофайлов
	if len(videoFiles) == 0 {
		fmt.Print("Videofiles not found: ")
		return fmt.Errorf("Videofiles not found: ")
	}

	// Создание временного файла
	videoListFile, err := os.Create(fmt.Sprintf("%s%s", vp.UserDirPath, videoList))
	if err != nil {
		fmt.Print("Error creating temporary file: ")
		return err
	}
	defer videoListFile.Close()
	// Очистка файла после завершения

	// Запись видеофайлов в временный файл
	for _, file := range videoFiles {
		_, err := fmt.Fprintf(videoListFile, "file '%s/%s'\n", clipsDir, file.Name())
		if err != nil {
			fmt.Print("Error writing to temporary file: ")
			return err
		}
	}

	// Получение списка аудиофайлов
	audioFiles, err := ioutil.ReadDir(songDir)
	if err != nil {
		fmt.Print("Error with getting audiofiles list: ")
		return err
	}

	// Проверка на наличие аудиофайлов
	if len(audioFiles) == 0 {
		fmt.Println("Videofiles not found: ")
		return fmt.Errorf("Videofiles not found: ")
	}

	// Проверка, что указанный номер песни существует
	if vp.SongNumber < 0 || vp.SongNumber >= len(audioFiles) {
		fmt.Println("Invalid song number: ")
		return fmt.Errorf("Invalid song number: ")
	}
	audioFile := audioFiles[vp.SongNumber]

	err = SongAndVideosMerge(audioFile)

	if err != nil {
		return err
	}

	for _, file := range audioFiles {
		os.Remove(songDir + "/" + file.Name())
	}
	for _, file := range videoFiles {
		os.Remove(clipsDir + "/" + file.Name())
	}

	return nil
}

func SongAndVideosMerge(audioFile fs.FileInfo) error {
	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", "videos.txt", "-c", "copy", outputVideo)
	err := cmd.Run()
	if err != nil {
		fmt.Print("Error merging files: ")
		return err
	}
	defer os.Remove(outputVideo) // Очистка файла после завершения

	// Объединение видео и аудио
	cmd = exec.Command("ffmpeg", "-i", outputVideo, "-i", songDir+"/"+audioFile.Name(), "-shortest", "-c:v", "copy", "-c:a", "aac", finalVideo)
	err = cmd.Run()
	if err != nil {
		fmt.Print("Error cmd running command: ")
		return err
	}

	os.Remove(videoList)

	return nil

	// Удаление временных файлов
}
