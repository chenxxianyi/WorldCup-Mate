package services

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"testing"
)

// SEC-06C: real decode + JPEG re-encode — fake files must be rejected and
// the output must be a clean JPEG.
func TestReencodeAvatar(t *testing.T) {
	t.Run("valid png re-encoded as jpeg", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 4, 4))
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			t.Fatal(err)
		}
		out, err := reencodeAvatar(bytes.NewReader(buf.Bytes()))
		if err != nil {
			t.Fatalf("valid png rejected: %v", err)
		}
		if len(out) < 2 || out[0] != 0xFF || out[1] != 0xD8 {
			t.Errorf("output is not a JPEG (magic bytes %x)", out[:2])
		}
	})

	t.Run("transparent png gets white background", func(t *testing.T) {
		// Fully transparent red pixel: RGBA(255,0,0,0). Without the white
		// composite it would encode as black (RGB 0,0,0).
		img := image.NewRGBA(image.Rect(0, 0, 2, 2))
		img.Set(0, 0, color.RGBA{255, 0, 0, 0})
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			t.Fatal(err)
		}
		out, err := reencodeAvatar(bytes.NewReader(buf.Bytes()))
		if err != nil {
			t.Fatalf("transparent png rejected: %v", err)
		}
		// Decode the re-encoded JPEG and assert the pixel is not black.
		decoded, err := decodeImage(out)
		if err != nil {
			t.Fatalf("re-encoded output undecodable: %v", err)
		}
		r, g, b, _ := decoded.At(0, 0).RGBA()
		if r == 0 && g == 0 && b == 0 {
			t.Error("transparent pixel came out black (white composite missing)")
		}
	})

	t.Run("gif accepted (first frame)", func(t *testing.T) {
		pal := color.Palette{color.White, color.Black}
		img := image.NewPaletted(image.Rect(0, 0, 2, 2), pal)
		var buf bytes.Buffer
		if err := gif.Encode(&buf, img, nil); err != nil {
			t.Fatal(err)
		}
		if _, err := reencodeAvatar(bytes.NewReader(buf.Bytes())); err != nil {
			t.Fatalf("valid gif rejected: %v", err)
		}
	})

	t.Run("fake image rejected", func(t *testing.T) {
		if _, err := reencodeAvatar(bytes.NewReader([]byte("definitely not an image"))); err == nil {
			t.Error("fake image accepted")
		}
	})

	t.Run("empty file rejected", func(t *testing.T) {
		if _, err := reencodeAvatar(bytes.NewReader(nil)); err == nil {
			t.Error("empty file accepted")
		}
	})
}

// SEC-06C: decompression-bomb guard — dimension check boundaries.
func TestValidateAvatarDimensions(t *testing.T) {
	cases := []struct {
		dx, dy int
		want   bool
	}{
		{0, 0, false},
		{1, 1, true},
		{4096, 4096, true},  // 16.7M exactly
		{4097, 4097, false}, // over the cap
		{100000, 1000, false},
		{-1, 10, false},
	}
	for _, c := range cases {
		if got := validateAvatarDimensions(c.dx, c.dy); got != c.want {
			t.Errorf("validateAvatarDimensions(%d,%d) = %v, want %v", c.dx, c.dy, got, c.want)
		}
	}
}

// SEC-06C: webp VP8X canvas pre-check (x/image/webp has no DecodeConfig).
func TestWebpDimensions(t *testing.T) {
	mk := func(w, h int) []byte {
		b := make([]byte, 30)
		copy(b[0:4], "RIFF")
		copy(b[8:12], "WEBP")
		copy(b[12:16], "VP8X")
		b[24] = byte(w)
		b[25] = byte(w >> 8)
		b[26] = byte(w >> 16)
		b[27] = byte(h)
		b[28] = byte(h >> 8)
		b[29] = byte(h >> 16)
		return b
	}
	if w, h, ok := webpDimensions(mk(800, 600)); !ok || w != 800 || h != 600 {
		t.Errorf("webpDimensions(800x600) = (%d,%d,%v)", w, h, ok)
	}
	if w, h, ok := webpDimensions(mk(0x12_3456, 100)); !ok || w != 0x123456 || h != 100 {
		t.Errorf("webpDimensions(24bit) = (%d,%d,%v)", w, h, ok)
	}
	if _, _, ok := webpDimensions([]byte("not a webp")); ok {
		t.Error("non-webp data reported ok")
	}
	if _, _, ok := webpDimensions(mk(10, 10)[:20]); ok {
		t.Error("truncated header reported ok")
	}
}

// SEC-06B: filenames must be random and not contain user-identifying data.
func TestRandomHex(t *testing.T) {
	a := randomHex(16)
	b := randomHex(16)
	if a == b {
		t.Error("collision between two random hex values")
	}
	if len(a) != 32 {
		t.Errorf("len = %d, want 32", len(a))
	}
}
