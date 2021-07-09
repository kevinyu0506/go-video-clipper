package main

import (
	"fmt"
    "log"
	"os"
	"os/exec"
)


// Set Enable or Disable (video|audio)
type Enable struct {
	video bool
	audio bool
}

type Revarse struct {
	video bool
	audio bool
}

type Crop struct {
	x      float64
	y      float32
	width  float64
	height float64
}

type Video struct {
	path      string
	video     bool
	audio     bool
	revarse   Revarse
	volume    float32
	overwrite bool
}

func (v *Video) Describe() {
    var args []string
    args = append(args, "-of", "json")
    //args = append(args, "-show_format_entry", "filename,duration")
    args = append(args, "-show_entries", "stream=codec_type")
    args = append(args, "-show_format")
    args = append(args, "-v", "quiet", v.path)

	description, err := exec.Command("ffprobe", args...).Output()
    if err != nil {
        log.Fatal(err)
        return
    }

    fmt.Printf(string(description))
}

func (v *Video) Clip() {
    var args []string
    args = append(args, "-ss", "01:23")
    args = append(args, "-i", v.path)
    args = append(args, "-vframes", "1")
    args = append(args, "-q:v", "2")
    args = append(args, "output.jpg")

	c := exec.Command("ffmpeg", args...)
    c.Run()

    fmt.Printf("Clip successfully")
}

func (v *Video) Output(output string) {
	var args []string

	// there is no data (no audio && no video) to set to file stream
	// so we will create empty file stream
	if (!v.audio || v.volume == 0) && !v.video && v.overwrite {
		f, _ := os.Create(output)
		f.Close()
		fmt.Println("-i " + v.path + " " + output + " -y") // args
		return
	}

	// set input file stream
	// ffmpeg -i "input.mp4" ...
	args = append(args, "-i")
	args = append(args, v.path)

	// disable video or audio
	// disable audio if volume == 0 or there no data for audio
	if !v.audio || v.volume <= 0 { args = append(args, "-an") }
	if !v.video { args = append(args, "-vn") }

	// set volume
	// args = append(args, "-af")
	// args = append(args, fmt."volume=")

	//if v.revarse.audio { args = append(args, "-vf") }
	//if v.revarse.video { args = append(args, "-vn") }


	args = append(args, output)
	if v.overwrite { args = append(args, "-y") }
	if !v.overwrite { args = append(args, "-n") }

	fmt.Println(args) // args
	c := exec.Command("ffmpeg", args...)
	c.Start()
}


func load(pathStream string) (*Video, error) {
	_, err := os.Stat(pathStream)
	if os.IsNotExist(err) {
        return nil, os.ErrNotExist
    }

	return &Video {
		path: pathStream,
		video: true,
		audio: true,
		volume: 10,
		overwrite: true,
	}, nil
}

func main() {
	v, _ := load("test.mp4")
	//v.video = false
	//v.audio = false
	//v.Output("./n.mp4")
    v.Describe()
    //v.Clip()
}
