package main

import (
	"math"
	"math/rand"

	"github.com/MadAppGang/kdbush"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Sprite struct {
	x   float64
	y   float64
	vx  float64
	vy  float64
	img *ebiten.Image
}

func newSprite(img *ebiten.Image, maxH, maxW int) *Sprite {
	x := float64(rand.Intn(maxW))
	y := float64(rand.Intn(maxH))
	vx := rand.Float64() * maxVel
	vy := rand.Float64() * maxVel
	return &Sprite{
		x:   x,
		y:   y,
		vx:  vx,
		vy:  vy,
		img: img,
	}
}

func (s *Sprite) update(ax, ay float64) {
	if math.Abs(s.vx+ax) < maxVel {
		s.vx += ax
	}
	if math.Abs(s.vy+ay) < maxVel {
		s.vy += ay
	}
	if s.x+s.vx < 0 || s.x+s.vx > screenWidth {
		s.vx *= -1
	}
	if s.y+s.vy < 0 || s.y+s.vy > screenHeight {
		s.vy *= -1
	}
	s.x += s.vx
	s.y += s.vy
}

func getColor(vx, vy float64) (float32, float32, float32, float32) {
	vel := float32(math.Abs(vx)+math.Abs(vy)) / 2
	r := vel / maxVel
	g := (maxVel - vel) / maxVel
	return r, g, 0.0, 1.0
}

func (s *Sprite) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	sx, sy := s.img.Size()
	op.GeoM.Translate(s.x-float64(sx)/2, s.y-float64(sy)/2)
	r, g, b, a := getColor(s.vx, s.vy)
	op.ColorScale.Scale(r, g, b, a)
	screen.DrawImage(s.img, op)
}

// naive aproach
func (s *Sprite) drawDetectionLines(sprites []*Sprite, screen *ebiten.Image) {
	var path vector.Path
	count := 1
	for _, st := range sprites {
		if count > maxNeighbours {
			break
		}
		if s.x == st.x && s.y == st.y {
			continue
		}
		if (s.x-st.x)*(s.x-st.x)+(s.y-st.y)*(s.y-st.y) > searchRad*searchRad {
			continue
		}
		path.MoveTo(float32(s.x), float32(s.y))
		path.LineTo(float32(st.x), float32(st.y))
		count++
	}
	op := &vector.StrokeOptions{}
	op.LineCap = vector.LineCapButt
	op.LineJoin = vector.LineJoinRound
	op.Width = 1
	vs, is := path.AppendVerticesAndIndicesForStroke([]ebiten.Vertex{}, []uint16{}, op)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
	}
	screen.DrawTriangles(vs, is, blankSubImage, &ebiten.DrawTrianglesOptions{
		AntiAlias: true,
	})
}

// uses kd-tree for optimization: huge improvement
func drawAllDetectionLines(sprites []*Sprite, screen *ebiten.Image) {
	points := [maxSprites]kdbush.Point{}
	for i, s := range sprites {
		points[i] = &kdbush.SimplePoint{X: s.x, Y: s.y}
	}
	tree := kdbush.NewBush(points[:], 16)

	var path vector.Path
	count := 0
	for idx, p := range points {
		matches := tree.Within(p, searchRad)
		X, Y := p.Coordinates()
		for _, i := range matches {
			x, y := points[i].Coordinates()
			path.MoveTo(float32(X), float32(Y))
			path.LineTo(float32(x), float32(y))
		}

		count++
		// drawing in batches due to reaching ebiten.MaxIndicesCount
		if count >= maxBatch || idx == maxSprites-1 {
			count = 0
			op := &vector.StrokeOptions{}
			op.LineCap = vector.LineCapButt
			op.LineJoin = vector.LineJoinRound
			op.Width = 1
			vs, is := path.AppendVerticesAndIndicesForStroke([]ebiten.Vertex{}, []uint16{}, op)
			for i := range vs {
				vs[i].SrcX = 1
				vs[i].SrcY = 1
				vs[i].ColorR = 230.0 / 255.0
				vs[i].ColorG = float32(150*i/len(vs)) / 255.0
				vs[i].ColorB = 0
			}
			screen.DrawTriangles(vs, is, blankSubImage, &ebiten.DrawTrianglesOptions{
				AntiAlias: true,
			})
			path = vector.Path{}
		}
	}
}
