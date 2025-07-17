package whisperuiapi

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func createNewEventDir(audioForTranscriptPath string) (*os.File, error) {

	confData := getWhisperData()

	baseAudioName := filepath.Base(audioForTranscriptPath)

	audioHash := getFileID(baseAudioName)

	audioFile, err := os.Open(audioForTranscriptPath)

	if err != nil {
		logrus.Errorf("err_info: can't open event audio file, audio_id: %s, err_text: %v", audioHash, err)
		return nil, err
	}

	defer audioFile.Close()

	newGradioDir := fmt.Sprintf("%s/%s", confData.AppData, audioHash)

	err = os.MkdirAll(newGradioDir, os.ModePerm)

	if err != nil {
		logrus.Errorf("err_info: can't make new event dir, audio_id: %s, err_text: %v", audioHash, err)
		return nil, err
	}

	newAudioInGradioDir, err := os.Create(filepath.Join(newGradioDir, baseAudioName))

	if err != nil {
		logrus.Errorf("err_info: can't make new file in event dir, audio_id: %s, err_text: %v", audioHash, err)
		return nil, err
	}

	defer newAudioInGradioDir.Close()

	_, err = io.Copy(newAudioInGradioDir, audioFile)

	if err != nil {
		logrus.Errorf("err_info: can't copy audio file to gradio event dir, audio_id: %s, err_text: %v", audioHash, err)
		return nil, err
	}

	logrus.Debugf("event dir for gradio was created, path_to_event_audio: %s", newAudioInGradioDir.Name())

	return newAudioInGradioDir, nil

}

func getFileID(baseAudioName string) string {
	newHash := sha512.New()

	newHash.Write([]byte(baseAudioName))

	hash := newHash.Sum(nil)

	return hex.EncodeToString(hash)
}
