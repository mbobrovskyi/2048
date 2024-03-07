package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
	"time"
)

type Game struct {
	screenWidth  int
	screenHeight int

	input      *Input
	board      *Board
	boardImage *ebiten.Image
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (g *Game) Update() error {
	g.input.Update()
	if err := g.board.Update(g.input); err != nil {
		return err
	}
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screenWidth, g.screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.boardImage == nil {
		w, h := g.board.Size()
		g.boardImage = ebiten.NewImage(w, h)
	}
	screen.Fill(backgroundColor)
	g.board.Draw(g.boardImage)
	op := &ebiten.DrawImageOptions{}
	sw, sh := screen.Size()
	bw, bh := g.boardImage.Size()
	x := (sw - bw) / 2
	y := (sh - bh) / 2
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(g.boardImage, op)
}

func NewGame(screenWidth, screenHeight, boardSize int) (*Game, error) {
	g := &Game{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,

		input: NewInput(),
	}

	var err error

	g.board, err = NewBoard(boardSize)
	if err != nil {
		return nil, err
	}

	return g, nil
}
