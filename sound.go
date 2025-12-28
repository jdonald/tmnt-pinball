package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	SampleRate  = 44100
	AudioBuffer = 4096
)

type SoundSystem struct {
	deviceID sdl.AudioDeviceID
	enabled  bool
}

type AudioCallback struct {
	toneFreq  float64
	tonePhase float64
	duration  int
	volume    float32
}

func NewSoundSystem() *SoundSystem {
	spec := &sdl.AudioSpec{
		Freq:     SampleRate,
		Format:   sdl.AUDIO_F32,
		Channels: 1,
		Samples:  AudioBuffer,
	}

	deviceID, err := sdl.OpenAudioDevice("", false, spec, nil, 0)
	if err != nil {
		return &SoundSystem{enabled: false}
	}

	sdl.PauseAudioDevice(deviceID, false)

	return &SoundSystem{
		deviceID: deviceID,
		enabled:  true,
	}
}

func (ss *SoundSystem) Close() {
	if ss.enabled {
		sdl.CloseAudioDevice(ss.deviceID)
	}
}

func (ss *SoundSystem) PlayTone(frequency float64, durationMs int, volume float32) {
	if !ss.enabled {
		return
	}

	samples := (SampleRate * durationMs) / 1000
	audioData := make([]float32, samples)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(SampleRate)
		// Generate a sine wave with volume envelope
		amplitude := volume
		if i < samples/10 {
			// Fade in
			amplitude *= float32(i) / float32(samples/10)
		} else if i > samples*9/10 {
			// Fade out
			amplitude *= float32(samples-i) / float32(samples/10)
		}
		audioData[i] = amplitude * float32(math.Sin(2*math.Pi*frequency*t))
	}

	ss.QueueAudio(audioData)
}

func (ss *SoundSystem) PlayCowabunga() {
	if !ss.enabled {
		return
	}

	// Synthesize "Cowabunga!" as a series of tones
	// This creates a distinctive sound pattern
	go func() {
		// "Co" - low to mid
		ss.PlayTone(220, 150, 0.3)
		sdl.Delay(150)

		// "wa" - mid to high
		ss.PlayTone(330, 150, 0.3)
		sdl.Delay(150)

		// "bun" - high
		ss.PlayTone(440, 200, 0.35)
		sdl.Delay(200)

		// "ga!" - drop and rise
		ss.PlayTone(330, 100, 0.3)
		sdl.Delay(100)
		ss.PlayTone(550, 200, 0.4)
	}()
}

func (ss *SoundSystem) PlayBallLaunch() {
	if !ss.enabled {
		return
	}
	// Rising "whoosh" sound
	go func() {
		for freq := 100.0; freq < 400.0; freq += 50 {
			ss.PlayTone(freq, 30, 0.2)
			sdl.Delay(30)
		}
	}()
}

func (ss *SoundSystem) PlayFlipperHit() {
	if !ss.enabled {
		return
	}
	// Quick "thwack" sound
	ss.PlayTone(150, 50, 0.25)
}

func (ss *SoundSystem) PlayBumperHit() {
	if !ss.enabled {
		return
	}
	// "Boing" sound - high pitched with quick decay
	ss.PlayTone(800, 80, 0.3)
}

func (ss *SoundSystem) PlayTargetHit() {
	if !ss.enabled {
		return
	}
	// "Ding" sound
	ss.PlayTone(1200, 100, 0.25)
}

func (ss *SoundSystem) PlayRailBounce() {
	if !ss.enabled {
		return
	}
	// Metallic bounce
	ss.PlayTone(600, 60, 0.15)
}

func (ss *SoundSystem) QueueAudio(data []float32) {
	if !ss.enabled {
		return
	}

	// Convert float32 slice to byte slice
	byteData := make([]byte, len(data)*4)
	for i, sample := range data {
		bits := math.Float32bits(sample)
		byteData[i*4] = byte(bits)
		byteData[i*4+1] = byte(bits >> 8)
		byteData[i*4+2] = byte(bits >> 16)
		byteData[i*4+3] = byte(bits >> 24)
	}

	sdl.QueueAudio(ss.deviceID, byteData)
}
