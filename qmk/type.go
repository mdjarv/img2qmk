package qmk

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type Animation struct {
	Name      string
	Frames    [][]byte
	FrameRate int
}

const headerTemplate = `#ifndef {{.Guard}}
#define {{.Guard}}

typedef struct {
    const char* frames;          // Pointer to frame data stored in PROGMEM
    uint16_t frame_size;         // Size of each frame in bytes
    const uint16_t* delays;      // Pointer to array of delays for each frame
    uint8_t frame_count;         // Number of frames
    uint8_t idx;                 // Current frame index
} Animation;

#endif // {{.Guard}}`

type HeaderData struct {
	Name  string
	Guard string
}

func PrintType() error {
	data := HeaderData{
		Name:  "Animation",
		Guard: strings.ToUpper("ANIMATION_H"),
	}

	tmpl, err := template.New("header").Parse(headerTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(os.Stdout, data)
}

// Define the Go template for the C code
const cTemplate = `#include "animation.h"

// Frame data stored in PROGMEM
static const char {{.Name}}_frames[] PROGMEM = {
{{- range $index, $frame := .Frames}}
    {{range $byteIndex, $byte := $frame}}{{if $byteIndex}}, {{end}}0x{{printf "%02X" $byte}}{{end}}, // Frame {{$index}}
{{- end}}
};

// Delays for each frame
static const uint16_t {{.Name}}_delays[] = { 
{{- range $index, $delay := .Delays}}
    {{- if $index}},{{end}}{{.}}
{{- end}}};

// Animation instance
static Animation {{.Name}}_animation = {
    {{.Name}}_frames,
    sizeof({{.Name}}_frames) / {{.FrameCount}},
    {{.Name}}_delays,
    {{.FrameCount}},
    0
};
`

func (anim Animation) Print() error {
	if len(anim.Frames) == 0 {
		return errors.New("animation must have at least one frame")
	}
	if anim.FrameRate <= 0 {
		return errors.New("frame rate must be positive")
	}

	// Prepare data for the template
	delays := make([]int, len(anim.Frames))
	for i := range delays {
		delays[i] = anim.FrameRate
	}

	data := struct {
		Name       string
		Frames     [][]byte
		Delays     []int
		FrameCount int
	}{
		Name:       anim.Name,
		Frames:     anim.Frames,
		Delays:     delays,
		FrameCount: len(anim.Frames),
	}

	// Parse and execute the template
	tmpl, err := template.New("cCode").Parse(cTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	fmt.Println(buf.String())
	return nil
}
