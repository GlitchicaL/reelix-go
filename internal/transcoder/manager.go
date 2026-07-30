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
	ID          string
	Video       string
	Quality     string
	Directory   string
	Cmd         *exec.Cmd
	LastSeen    time.Time
	CreatedAt   time.Time
	State       string
	AccessCount int
}

type TranscodeManager struct {
	mu       sync.Mutex
	sessions map[string]*Session
	workers  chan struct{}
}

func NewTranscodeManager(workers int) *TranscodeManager {
	mgr := &TranscodeManager{
		sessions: make(map[string]*Session),
		workers:  make(chan struct{}, workers),
	}

	log.Printf("[transcoder] Manager initialized with %d workers", workers)

	StartCleanupManager(mgr, 5*time.Minute, 30*time.Minute)

	return mgr
}

func (m *TranscodeManager) Start(videoId, video, quality string) (*Session, error) {
	key := videoId + ":" + quality

	m.workers <- struct{}{}

	m.mu.Lock()

	if session, exists := m.sessions[key]; exists {
		m.mu.Unlock()
		session.LastSeen = time.Now()
		session.State = "active"
		session.AccessCount++
		log.Printf("[%s] Accessed (count: %d)", session.ID, session.AccessCount)
		return session, nil
	}

	id := uuid.NewString()
	dir := filepath.Join("/tmp/hls/", id)

	if err := os.MkdirAll(dir, 0755); err != nil {
		m.mu.Unlock()
		<-m.workers
		return nil, err
	}

	profile, exists := GetQualityProfile(quality)
	if !exists {
		m.mu.Unlock()
		<-m.workers
		return nil, errors.New("unknown quality profile: " + quality)
	}

	session := &Session{
		ID:        id,
		Directory: dir,
		Video:     video,
		Quality:   quality,
		LastSeen:  time.Now(),
		CreatedAt: time.Now(),
		State:     "pending",
	}

	cmd := BuildFFmpegCommand(profile, video, dir)
	session.Cmd = cmd

	m.sessions[key] = session

	log.Printf("[%s] Created: video=%s, quality=%s, dir=%s", session.ID, video, quality, dir)
	log.Printf("[%s] FFmpeg started with profile: %s (%s, CRF %s)", session.ID, quality, profile.MaxRate, profile.CRF)

	m.mu.Unlock()

	go func() {
		output, err := cmd.CombinedOutput()

		if err != nil {
			log.Printf("[%s] FFmpeg failed: %v\nOutput: %s", session.ID, err, string(output))
		} else {
			log.Printf("[%s] FFmpeg completed\nOutput: %s", session.ID, string(output))
		}

		m.mu.Lock()
		session.State = "expired"
		m.mu.Unlock()

		<-m.workers
	}()

	err := waitForPlaylist(dir)

	if err != nil {
		m.mu.Lock()
		delete(m.sessions, key)
		m.mu.Unlock()
		<-m.workers
		return nil, err
	}

	m.mu.Lock()
	session.State = "ready"
	log.Printf("[%s] Ready (playlist available)", session.ID)
	m.mu.Unlock()

	return session, nil
}

func (m *TranscodeManager) Touch(sessionID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, session := range m.sessions {
		if session.ID == sessionID {
			session.LastSeen = time.Now()
			session.State = "active"
			session.AccessCount++
			log.Printf("[%s] Accessed (count: %d)", session.ID, session.AccessCount)
			return true
		}
	}

	return false
}

func (m *TranscodeManager) GetSession(sessionID string) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, session := range m.sessions {
		if session.ID == sessionID {
			return session
		}
	}

	return nil
}

func waitForPlaylist(dir string) error {
	playlist := filepath.Join(dir, "index.m3u8")

	for i := 0; i < 50; i++ {
		if _, err := os.Stat(playlist); err == nil {
			return nil
		}

		time.Sleep(100 * time.Millisecond)
	}
	return errors.New("playlist creation timeout")
}
