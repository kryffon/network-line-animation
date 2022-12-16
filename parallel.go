package main

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// container for screen
type MScreen struct {
	mu     sync.Mutex
	screen *ebiten.Image
}

func (s *Sprite) parallel_draw(ms *MScreen) {
	op := &ebiten.DrawImageOptions{}
	sx, sy := s.img.Size()
	op.GeoM.Translate(s.x-float64(sx)/2, s.y-float64(sy)/2)
	r, g, b, a := getColor(s.vx, s.vy)
	op.ColorScale.Scale(r, g, b, a)

	ms.mu.Lock()
	ms.screen.DrawImage(s.img, op)
	ms.mu.Unlock()
}

func (s *Sprite) parallel_drawDetectionLines(sprites []*Sprite, ms *MScreen) {
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
	ms.mu.Lock()
	ms.screen.DrawTriangles(vs, is, blankSubImage, &ebiten.DrawTrianglesOptions{
		AntiAlias: true,
	})
	ms.mu.Unlock()
}

// not a huge performance improvement
func drawParrallel(sprites []*Sprite, screen *ebiten.Image) {
	s_len := len(sprites)

	// create lockable screen object
	ms := &MScreen{
		screen: screen,
	}
	// create buffered channels and woker function
	sprite_chan := make(chan *Sprite, 4)
	results := make(chan int, s_len)
	worker := func(wid int, sprite_chan chan *Sprite) {
		for s := range sprite_chan {
			s.parallel_draw(ms)
			s.parallel_drawDetectionLines(sprites, ms)
			results <- wid
		}
	}
	// run all workers
	for i := 0; i < 4; i++ {
		go worker(i, sprite_chan)
	}
	// push all sprites to buffered channel
	for _, s := range sprites {
		sprite_chan <- s
	}
	// close buffer to signal end
	close(sprite_chan)
	// wait for all jobs to complete
	for a := 1; a <= s_len; a++ {
		<-results
	}
}
