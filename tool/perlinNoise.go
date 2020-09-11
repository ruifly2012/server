package tool

import (
	"testing"

	perlin "github.com/aquilax/go-perlin"
)

const (
	seed = 123
)

func TestPerlinNoise1D(t *testing.T) {
	expected := 0.0
	p := perlin.NewPerlin(2, 2, 3, seed)
	noise := p.Noise1D(10)
	if noise != expected {
		t.Fail()
		t.Logf("Wrong node result: given: %f, expected: %f", noise, expected)
	}
}

func TestPerlinNoise2D(t *testing.T) {
	expected := 0.0
	p := perlin.NewPerlin(2, 2, 3, seed)
	noise := p.Noise2D(10, 10)
	if noise != expected {
		t.Fail()
		t.Logf("Wrong node result: given: %f, expected: %f", noise, expected)
	}
}

func TestPerlinNoise3D(t *testing.T) {
	expected := 0.0
	p := perlin.NewPerlin(2, 2, 3, seed)
	noise := p.Noise3D(10, 10, 10)
	if noise != expected {
		t.Fail()
		t.Logf("Wrong node result: given: %f, expected: %f", noise, expected)
	}
}

func BenchmarkPerlinNoise1D(b *testing.B) {
	p := perlin.NewPerlin(2, 2, 3, seed)
	for n := 0; n < b.N; n++ {
		p.Noise1D(10)
	}
}

func BenchmarkPerlinNoise2D(b *testing.B) {
	p := perlin.NewPerlin(2, 2, 3, seed)
	for n := 0; n < b.N; n++ {
		p.Noise2D(10, 10)
	}
}
func BenchmarkPerlinNoise3D(b *testing.B) {
	p := perlin.NewPerlin(2, 2, 3, seed)
	for n := 0; n < b.N; n++ {
		p.Noise3D(10, 10, 10)
	}
}
