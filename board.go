package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

const (
	board_width_blocks  int = 10
	board_height_blocks int = 20

	board_width_chars  int = 10 * 2
	board_height_chars int = 20

	small_width_blocks  int = 4
	small_height_blocks int = 4

	small_width_chars  int = 4 * 2
	small_height_chars int = 4
)

type Board struct {
	/* Temporary matrix for the current real mino, while being in this
	   matrix, the mino does not trigger rows cleaning */
	upperMatrix [][]int
	/* Matrix for the ghost (max-fallen) mino related to the current real mino */
	ghostMatrix [][]int
	/* Matrix for the real minoes stuck and solidified, while in this matrix,
	   the minoes can trigger rows cleaning */
	solidMatrix [][]int
	/* Matrix for the current reserved mino */
	reservedMatrix [][]int
	/* Matrix for the next coming mino */
	futureMatrix [][]int
}

func NewBoard() *Board {
	var b = Board{
		upperMatrix:    [][]int{},
		ghostMatrix:    [][]int{},
		solidMatrix:    [][]int{},
		reservedMatrix: [][]int{},
		futureMatrix:   [][]int{},
	}

	tbprint(0, 0, termbox.ColorWhite, termbox.ColorBlack, "New board crated")

	b.upperMatrix = make([][]int, board_height_blocks)
	for x := 0; x < len(b.upperMatrix); x++ {
		b.upperMatrix[x] = make([]int, board_width_blocks)
	}

	b.solidMatrix = make([][]int, board_height_blocks)
	for x := 0; x < len(b.solidMatrix); x++ {
		b.solidMatrix[x] = make([]int, board_width_blocks)
	}

	b.ghostMatrix = make([][]int, board_height_blocks)
	for x := 0; x < len(b.ghostMatrix); x++ {
		b.ghostMatrix[x] = make([]int, board_width_blocks)
	}

	b.reservedMatrix = make([][]int, small_height_blocks)
	for x := 0; x < len(b.reservedMatrix); x++ {
		b.reservedMatrix[x] = make([]int, small_width_blocks)
	}

	b.futureMatrix = make([][]int, small_height_blocks)
	for x := 0; x < len(b.futureMatrix); x++ {
		b.futureMatrix[x] = make([]int, small_width_blocks)
	}

	return &b
}

func (b *Board) Clear() {
	for by := 0; by < len(b.solidMatrix); by++ {
		for bx := 0; bx < len(b.solidMatrix[0]); bx++ {
			b.solidMatrix[by][bx] = 0
		}
	}
}

func (b *Board) ClearRow(r int) {
	// Shift rows above the cleared row down by one slot
	for y := r - 1; y >= 0; y-- {
		for x := 0; x < len(b.solidMatrix[0]); x++ {
			b.solidMatrix[y+1][x] = b.solidMatrix[y][x]
		}
	}

	// Clear the specified row
	for x := 0; x < len(b.solidMatrix[0]); x++ {
		b.solidMatrix[0][x] = 0
	}
}

func (b *Board) Draw() {
	b.drawRect(board_width_chars+2, board_height_chars+2, 0, 0)
	b.drawRect(small_width_chars+2, small_height_chars+2, -17, -8+16)
	b.drawRect(small_width_chars+2, small_height_chars+2, +17, -8)

	maxw, maxh := termbox.Size()
	offx_board := (maxw - board_width_chars) / 2
	offy_board := (maxh - board_height_chars) / 2
	offx_small_stored := (maxw-small_width_chars)/2 - 17
	offy_small_stored := (maxh-small_height_chars)/2 - 8 + 16
	offx_small_future := (maxw-small_width_chars)/2 + 17
	offy_small_future := (maxh-small_height_chars)/2 - 8
	tbprint(offx_small_stored, offy_small_stored-2, termbox.ColorGreen, termbox.ColorDefault, "Reserved")
	tbprint(offx_small_future+1, offy_small_future+5, termbox.ColorGreen, termbox.ColorDefault, "Future")

	b.drawMatrix(b.ghostMatrix, board_height_chars, board_width_chars, offx_board, offy_board, '(', ')', termbox.ColorWhite, termbox.ColorDefault, false)
	b.drawMatrix(b.upperMatrix, board_height_chars, board_width_chars, offx_board, offy_board, '[', ']', termbox.ColorDefault, termbox.ColorDefault, true)
	b.drawMatrix(b.solidMatrix, board_height_chars, board_width_chars, offx_board, offy_board, '[', ']', termbox.ColorDefault, termbox.ColorDefault, true)
	b.drawMatrix(b.reservedMatrix, small_height_chars, small_width_chars, offx_small_stored, offy_small_stored, '[', ']', termbox.ColorDefault, termbox.ColorDefault, true)
	b.drawMatrix(b.futureMatrix, small_height_chars, small_width_chars, offx_small_future, offy_small_future, '[', ']', termbox.ColorDefault, termbox.ColorDefault, true)
}

func (b *Board) drawRect(width, height int, offx, offy int) {

	termWidth, termHeight := termbox.Size()

	rectX := (termWidth-width)/2 + offx
	rectY := (termHeight-height)/2 + offy

	for x := rectX; x < rectX+width; x++ {
		termbox.SetCell(x, rectY, '-', termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(x, rectY+height-1, '-', termbox.ColorDefault, termbox.ColorDefault)
	}

	for y := rectY + 1; y < rectY+height-1; y++ {
		termbox.SetCell(rectX, y, '|', termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(rectX+width-1, y, '|', termbox.ColorDefault, termbox.ColorDefault)
	}

	termbox.SetCell(rectX, rectY, '+', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(rectX+width-1, rectY, '+', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(rectX, rectY+height-1, '+', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(rectX+width-1, rectY+height-1, '+', termbox.ColorDefault, termbox.ColorDefault)
}

func (b *Board) drawMatrix(matrix [][]int, height, width, offx, offy int, op, cl rune, fg, bg termbox.Attribute, chooseFg bool) {

	for by := 0; by < height; by += 1 {
		for bx := 0; bx < width; bx += 1 {
			if matrix[by][bx/2] > 0 {
				if chooseFg {
					fg = b.shapeToColor(matrix[by][bx/2])
				}
				if bx%2 == 0 {
					termbox.SetCell(offx+bx, offy+by, op, fg, bg)
				} else {
					termbox.SetCell(offx+bx, offy+by, cl, fg, bg)
				}
			}
		}
	}
}

func (b *Board) shapeToColor(v int) termbox.Attribute {
	var c termbox.Attribute
	switch v - 1 {
	case int(t_I):
		c = termbox.ColorYellow
	case int(t_J):
		c = termbox.ColorGreen
	case int(t_L):
		c = termbox.ColorRed
	case int(t_O):
		c = termbox.ColorCyan
	case int(t_S):
		c = termbox.ColorBlue
	case int(t_T):
		c = termbox.ColorLightRed
	case int(t_Z):
		c = termbox.ColorBlue | termbox.ColorYellow
	default:
		c = termbox.ColorRed
	}
	// tbprint(1, 2, termbox.ColorRed, termbox.ColorDefault, fmt.Sprintf("color: %d", c))
	return c
}

func (b *Board) FixBlocks(m *Mino) {
	for my := range m.matrix {
		for mx := range m.matrix[my] {
			b.solidMatrix[m.position.y+my][m.position.x+mx] += m.matrix[my][mx] * (int(m.shape) + 1)
		}
	}
}

func (b *Board) ClearMatrix(matrix [][]int) {
	for my := range matrix {
		for mx := range matrix[my] {
			matrix[my][mx] = 0
		}
	}
}

// ProjectMino sets a given mino in the given matrix between board's game matrixes with an xy pos offset
func (b *Board) ProjectMino(m *Mino, matrix [][]int, posy, posx int, clear bool) error {
	if m == nil {
		return nil
	}
	if clear {
		b.ClearMatrix(matrix)
	}
	for my := 0; my < len(m.matrix); my++ {
		for mx := 0; mx < len(m.matrix[my]); mx++ {
			if m.matrix[my][mx] == 0 {
				continue
			}
			if b.solidMatrix[posy+my][posx+mx] > 0 {
				return fmt.Errorf("block x:%d y:%d, busy", posx+mx, posy+my)
			}
			matrix[posy+my][posx+mx] = m.matrix[my][mx] * (int(m.shape) + 1)
		}
	}
	return nil
}
