package main

type mino_t int

type pos_t struct {
	x int
	y int
}

type rot_t int

const (
	None rot_t = iota
	Single
	Double
	Triple
)

const (
	t_I mino_t = iota
	t_O
	t_T
	t_L
	t_J
	t_S
	t_Z
)

type Mino struct {
	shape    mino_t
	position *pos_t
	rotation rot_t
	matrix   [][]int
	stuck    bool
}

func (m *Mino) setMatrix() {
	switch m.shape {
	case t_I:
		m.matrix = [][]int{
			{1, 1, 1, 1},
		}
	case t_O:
		m.matrix = [][]int{
			{1, 1},
			{1, 1},
		}
	case t_T:
		m.matrix = [][]int{
			{1, 1, 1},
			{0, 1, 0},
		}
	case t_L:
		m.matrix = [][]int{
			{1, 0},
			{1, 0},
			{1, 1},
		}
	case t_J:
		m.matrix = [][]int{
			{0, 1},
			{0, 1},
			{1, 1},
		}
	case t_S:
		m.matrix = [][]int{
			{0, 1, 1},
			{1, 1, 0},
		}
	case t_Z:
		m.matrix = [][]int{
			{1, 1, 0},
			{0, 1, 1},
		}
	}
}

func (m *Mino) Fall() {
	m.position.y += 1
}
