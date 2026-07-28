package transcoder

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID string

	Video string

	Quality string

	Directory string

	Cmd *exec.Cmd

	LastSeen time.Time
}

type TranscodeManager struct {
	mu sync.Mutex

	sessions         map[string]*Session
	activeTranscodes map[string]*Session

	workers chan struct{}
}

func NewTranscodeManager(workers int) *TranscodeManager {
	return &TranscodeManager{
		sessions:         make(map[string]*Session),
		activeTranscodes: make(map[string]*Session),
		workers:          make(chan struct{}, workers),
	}
}

func (m *TranscodeManager) Start(videoId, video, quality string) (*Session, error) {
	key := videoId + ":" + quality

	m.workers <- struct{}{}

	m.mu.Lock()
	// Check if active transcode exists
	if session, exists := m.activeTranscodes[key]; exists {
		m.mu.Unlock()
		return session, nil // Return existing session
	}

	id := uuid.NewString()

	dir := filepath.Join("/tmp/hls/", id)

	os.MkdirAll(dir, 0755)

	session := &Session{
		ID:        id,
		Directory: dir,
		Video:     video,
		Quality:   quality,
		LastSeen:  time.Now(),
	}

	cmd := buildFFmpeg(video, quality, dir)

	session.Cmd = cmd

	m.sessions[id] = session
	m.activeTranscodes[key] = session

	m.mu.Unlock()

	go func() {
		// Combine stdout and stderr for simpler logging
		output, err := cmd.CombinedOutput()

		if err != nil {
			log.Printf("[%s] FFmpeg failed: %v\nOutput: %s", session.ID, err, string(output))
		} else {
			log.Printf("[%s] FFmpeg completed\nOutput: %s", session.ID, string(output))
		}

		m.mu.Lock()
		delete(m.activeTranscodes, key)
		m.mu.Unlock()

		<-m.workers
		os.RemoveAll(dir)
	}()

	err := waitForPlaylist(dir) // Blocks until index.m3u8 exists or timeout

	if err != nil {
		return nil, err
	}

	return session, nil
}

func waitForPlaylist(dir string) error {
	playlist := filepath.Join(dir, "index.m3u8")

	for i := 0; i < 50; i++ { // 5 second timeout
		if _, err := os.Stat(playlist); err == nil {
			return nil
		}

		time.Sleep(100 * time.Millisecond)
	}
	return errors.New("playlist creation timeout")
}

func (m *TranscodeManager) Cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, s := range m.sessions {

		if time.Since(s.LastSeen) > 2*time.Minute {

			s.Cmd.Process.Kill()

			os.RemoveAll(s.Directory)

		}

	}

}

func buildFFmpeg(
	input string,
	quality string,
	outputDir string,
) *exec.Cmd {
	return exec.Command(

		"ffmpeg",

		"-i", input,

		"-vf", "scale=1280:720",

		"-c:v", "libx264",

		"-preset", "veryfast",

		"-crf", "23",

		"-c:a", "aac",

		"-b:a", "128k",

		"-f", "hls",

		"-hls_time", "4",

		"-hls_list_size", "15",

		"-hls_flags", "delete_segments",

		filepath.Join(outputDir, "index.m3u8"),
	)
}
