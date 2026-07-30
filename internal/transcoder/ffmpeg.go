package transcoder

import (
	"os/exec"
	"path/filepath"
	"strconv"
)

func BuildFFmpegCommand(profile QualityProfile, inputPath, outputDir string) *exec.Cmd {
	args := []string{
		"-i", inputPath,
	}

	if profile.Resolution == "1080p" {
		args = append(args, "-vf", "scale=-1:1080")
	} else if profile.Resolution == "720p" {
		args = append(args, "-vf", "scale=-1:720")
	}

	args = append(args,
		"-c:v", "libx264",
		"-preset", "veryfast",
	)

	if profile.CRF != "" {
		args = append(args, "-crf", profile.CRF)
	}

	if profile.MaxRate != "" {
		args = append(args, "-maxrate", profile.MaxRate)
	}

	if profile.BufSize != "" {
		args = append(args, "-bufsize", profile.BufSize)
	}

	args = append(args,
		"-c:a", "aac",
		"-b:a", "128k",
		"-f", "hls",
		"-hls_time", "4",
		"-hls_list_size", "15",
		"-hls_flags", "delete_segments",
		"-start_number", "0",
		filepath.Join(outputDir, "index.m3u8"),
	)

	cmd := exec.Command("ffmpeg", args...)
	return cmd
}

func formatBitrate(bps int64) string {
	if bps >= 1000000 {
		return strconv.FormatInt(bps/1000000, 10) + " Mbps"
	}
	if bps >= 1000 {
		return strconv.FormatInt(bps/1000, 10) + " Kbps"
	}
	return strconv.FormatInt(bps, 10) + " bps"
}
