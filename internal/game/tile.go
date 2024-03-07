package game

import (
	"errors"
	"image/color"
	"log"
	"math/rand"
	"sort"
	"strconv"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var (
	mplusSmallFont  font.Face
	mplusNormalFont font.Face
	mplusBigFont    font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusSmallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	tileImage.Fill(color.White)
}

type TileData struct {
	value int
	x     int
	y     int
}

type Tile struct {
	current TileData

	next TileData

	movingCount       int
	startPoppingCount int
	poppingCount      int
}

func (t *Tile) Pos() (int, int) {
	return t.current.x, t.current.y
}

func (t *Tile) NextPos() (int, int) {
	return t.next.x, t.next.y
}

func (t *Tile) Value() int {
	return t.current.value
}

func (t *Tile) NextValue() int {
	return t.next.value
}

func NewTile(value int, x, y int) *Tile {
	return &Tile{
		current: TileData{
			value: value,
			x:     x,
			y:     y,
		},
		startPoppingCount: maxPoppingCount,
	}
}

func (t *Tile) IsMoving() bool {
	return 0 < t.movingCount
}

func (t *Tile) stopAnimation() {
	if 0 < t.movingCount {
		t.current = t.next
		t.next = TileData{}
	}
	t.movingCount = 0
	t.startPoppingCount = 0
	t.poppingCount = 0
}

func tileAt(tiles map[*Tile]struct{}, x, y int) *Tile {
	var result *Tile
	for t := range tiles {
		if t.current.x != x || t.current.y != y {
			continue
		}
		if result != nil {
			panic("not reach")
		}
		result = t
	}
	return result
}

func currentOrNextTileAt(tiles map[*Tile]struct{}, x, y int) *Tile {
	var result *Tile
	for t := range tiles {
		if 0 < t.movingCount {
			if t.next.x != x || t.next.y != y || t.next.value == 0 {
				continue
			}
		} else {
			if t.current.x != x || t.current.y != y {
				continue
			}
		}
		if result != nil {
			panic("not reach")
		}
		result = t
	}
	return result
}

const (
	maxMovingCount  = 5
	maxPoppingCount = 6
)

func MoveTiles(tiles map[*Tile]struct{}, size int, dir Dir) bool {
	vx, vy := dir.Vector()
	var tx []int
	var ty []int
	for i := 0; i < size; i++ {
		tx = append(tx, i)
		ty = append(ty, i)
	}
	if vx > 0 {
		sort.Sort(sort.Reverse(sort.IntSlice(tx)))
	}
	if vy > 0 {
		sort.Sort(sort.Reverse(sort.IntSlice(ty)))
	}

	moved := false
	for _, j := range ty {
		for _, i := range tx {
			t := tileAt(tiles, i, j)

			if t == nil {
				continue
			}

			if t.next != (TileData{}) {
				panic("not reach")
			}

			if t.IsMoving() {
				panic("not reach")
			}

			ii := i
			jj := j

			for {
				ni := ii + vx
				nj := jj + vy
				if ni < 0 || ni >= size || nj < 0 || nj >= size {
					break
				}
				tt := currentOrNextTileAt(tiles, ni, nj)
				if tt == nil {
					ii = ni
					jj = nj
					moved = true
					continue
				}
				if t.current.value != tt.current.value {
					break
				}
				if 0 < tt.movingCount && tt.current.value != tt.next.value {
					break
				}
				ii = ni
				jj = nj
				moved = true
				break
			}

			next := TileData{}
			next.value = t.current.value

			if tt := currentOrNextTileAt(tiles, ii, jj); tt != t && tt != nil {
				next.value = t.current.value + tt.current.value
				tt.next.value = 0
				tt.next.x = ii
				tt.next.y = jj
				tt.movingCount = maxMovingCount
			}

			next.x = ii
			next.y = jj

			if t.current != next {
				t.next = next
				t.movingCount = maxMovingCount
			}
		}
	}

	if !moved {
		for t := range tiles {
			t.next = TileData{}
			t.movingCount = 0
		}
	}

	return moved
}

func addRandomTile(tiles map[*Tile]struct{}, size int) error {
	cells := make([]bool, size*size)
	for t := range tiles {
		if t.IsMoving() {
			panic("not reach")
		}
		i := t.current.x + t.current.y*size
		cells[i] = true
	}
	availableCells := []int{}
	for i, b := range cells {
		if b {
			continue
		}
		availableCells = append(availableCells, i)
	}
	if len(availableCells) == 0 {
		return errors.New("twenty48: there is no space to add a new tile")
	}

	c := availableCells[rand.Intn(len(availableCells))]
	v := 2

	if rand.Intn(10) == 0 {
		v = 4
	}

	x := c % size
	y := c / size

	t := NewTile(v, x, y)

	tiles[t] = struct{}{}

	return nil
}

func (t *Tile) Update() error {
	switch {
	case 0 < t.movingCount:
		t.movingCount--
		if t.movingCount == 0 {
			if t.current.value != t.next.value && 0 < t.next.value {
				t.poppingCount = maxPoppingCount
			}
			t.current = t.next
			t.next = TileData{}
		}
	case 0 < t.startPoppingCount:
		t.startPoppingCount--
	case 0 < t.poppingCount:
		t.poppingCount--
	}
	return nil
}

func mean(a, b int, rate float64) int {
	return int(float64(a)*(1-rate) + float64(b)*rate)
}

func meanF(a, b float64, rate float64) float64 {
	return a*(1-rate) + b*rate
}

const (
	tileSize   = 80
	tileMargin = 4
)

var (
	tileImage = ebiten.NewImage(tileSize, tileSize)
)

func (t *Tile) Draw(boardImage *ebiten.Image) {
	i, j := t.current.x, t.current.y
	ni, nj := t.next.x, t.next.y
	v := t.current.value

	if v == 0 {
		return
	}

	op := &ebiten.DrawImageOptions{}

	x := i*tileSize + (i+1)*tileMargin
	y := j*tileSize + (j+1)*tileMargin
	nx := ni*tileSize + (ni+1)*tileMargin
	ny := nj*tileSize + (nj+1)*tileMargin

	switch {
	case 0 < t.movingCount:
		rate := 1 - float64(t.movingCount)/maxMovingCount
		x = mean(x, nx, rate)
		y = mean(y, ny, rate)
	case 0 < t.startPoppingCount:
		rate := 1 - float64(t.startPoppingCount)/float64(maxPoppingCount)
		scale := meanF(0.0, 1.0, rate)
		op.GeoM.Translate(float64(-tileSize/2), float64(-tileSize/2))
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(tileSize/2), float64(tileSize/2))
	case 0 < t.poppingCount:
		const maxScale = 1.2
		rate := 0.0
		if maxPoppingCount*2/3 <= t.poppingCount {
			rate = 1 - float64(t.poppingCount-2*maxPoppingCount/3)/float64(maxPoppingCount/3)
		} else {
			rate = float64(t.poppingCount) / float64(maxPoppingCount*2/3)
		}
		scale := meanF(1.0, maxScale, rate)
		op.GeoM.Translate(float64(-tileSize/2), float64(-tileSize/2))
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(tileSize/2), float64(tileSize/2))
	}

	op.GeoM.Translate(float64(x), float64(y))
	op.ColorM.ScaleWithColor(tileBackgroundColor(v))
	boardImage.DrawImage(tileImage, op)
	str := strconv.Itoa(v)

	f := mplusBigFont

	switch {
	case 3 < len(str):
		f = mplusSmallFont
	case 2 < len(str):
		f = mplusNormalFont
	}

	bound, _ := font.BoundString(f, str)

	w := (bound.Max.X - bound.Min.X).Ceil()
	h := (bound.Max.Y - bound.Min.Y).Ceil()
	x = x + (tileSize-w)/2
	y = y + (tileSize-h)/2 + h

	text.Draw(boardImage, str, f, x, y, tileColor(v))
}
