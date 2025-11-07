package main

import (
	"bytes"
	"fmt"
	"strings"
)

type node struct {
	val         string
	left, right *node
}

type pos struct {
	x, y  int
	width int
	text  []string
}

var (
	// nodeWidth should match len(drawCircle()[0])
	nodeWidth   = 7
	gap         = 3               // extra gap between nodes (horizontal)
	unit        = nodeWidth + gap // the horizontal unit for inorder placement
	depthGap    = 6
	nextInorder = 0
	positions   = map[*node]pos{}
	maxX, maxY  = 0, 0
)

// returns an ASCII circle for a value - fixed width = nodeWidth
func drawCircle(val string) []string {
	// ensure the middle line fits into nodeWidth (7)
	// pattern width: "  ___  " = 7
	//               " /   \ " = 7
	//               "|  x  |" = 7  (x placed centered)
	//               " \___/ " = 7
	// For multi-char val, you may need to pad/truncate; here we center.
	inner := val
	// center the inner within 3 chars space (we allow up to 3 chars). pad if shorter
	// available inner width = nodeWidth - 4 -> 3
	maxInner := nodeWidth - 4
	if len(inner) > maxInner {
		inner = inner[:maxInner]
	}
	// center pad
	leftPad := (maxInner - len(inner)) / 2
	rightPad := maxInner - len(inner) - leftPad
	mid := fmt.Sprintf("| %s%s%s |", strings.Repeat(" ", leftPad), inner, strings.Repeat(" ", rightPad))
	return []string{
		"  ___  ",
		" /   \\ ",
		mid,
		" \\___/ ",
	}
}

// layOut does inorder placement using a fixed unit for x positions
func layOut(n *node, depth int) {
	if n == nil {
		return
	}
	layOut(n.left, depth+1)

	circle := drawCircle(n.val)
	x := nextInorder * unit
	y := depth * depthGap

	positions[n] = pos{x: x, y: y, width: len(circle[0]), text: circle}
	if x > maxX {
		maxX = x
	}
	if y > maxY {
		maxY = y
	}
	nextInorder++
	layOut(n.right, depth+1)
}

func draw(root *node) string {
	positions = make(map[*node]pos)
	nextInorder = 0
	maxX, maxY = 0, 0

	layOut(root, 0)

	width := maxX + unit + 20
	height := maxY + depthGap + 10

	lines := make([][]byte, height)
	for i := range lines {
		lines[i] = make([]byte, width)
		for j := range lines[i] {
			lines[i][j] = ' '
		}
	}

	writeStr := func(s string, x, y int) {
		if y < 0 || y >= len(lines) {
			return
		}
		line := lines[y]
		for i := 0; i < len(s) && x+i < len(line); i++ {
			if x+i >= 0 {
				line[x+i] = s[i]
			}
		}
	}

	var drawRec func(n *node)
	drawRec = func(n *node) {
		if n == nil {
			return
		}
		p := positions[n]
		// draw the circle lines
		for i, l := range p.text {
			writeStr(l, p.x, p.y+i)
		}

		centerX := p.x + p.width/2
		bottomY := p.y + len(p.text) // first row below the node

		// left connector: diagonal from parent to child
		if n.left != nil {
			c := positions[n.left]
			childCenter := c.x + c.width/2
			dy := c.y - bottomY
			if dy <= 0 {
				// if child is directly below or same row, draw vertical
				for y := bottomY; y <= c.y; y++ {
					if y >= 0 && y < height && centerX >= 0 && centerX < width {
						lines[y][centerX] = '|'
					}
				}
			} else {
				sign := -1
				if childCenter > centerX {
					sign = 1
				}
				// step from parent downwards, moving horizontally one per row
				for s := 1; s <= dy; s++ {
					x := centerX + s*sign
					y := bottomY + s - 1
					if x >= 0 && x < width && y >= 0 && y < height {
						if sign < 0 {
							lines[y][x] = '/'
						} else {
							lines[y][x] = '\\'
						}
					}
				}
			}
		}

		// right connector: diagonal from parent to child
		if n.right != nil {
			c := positions[n.right]
			childCenter := c.x + c.width/2
			dy := c.y - bottomY
			if dy <= 0 {
				for y := bottomY; y <= c.y; y++ {
					if y >= 0 && y < height && centerX >= 0 && centerX < width {
						lines[y][centerX] = '|'
					}
				}
			} else {
				sign := 1
				if childCenter < centerX {
					sign = -1
				}
				for s := 1; s <= dy; s++ {
					x := centerX + s*sign
					y := bottomY + s - 1
					if x >= 0 && x < width && y >= 0 && y < height {
						if sign < 0 {
							lines[y][x] = '/'
						} else {
							lines[y][x] = '\\'
						}
					}
				}
			}
		}

		drawRec(n.left)
		drawRec(n.right)
	}
	drawRec(root)

	var buf bytes.Buffer
	for _, row := range lines {
		trimmed := strings.TrimRight(string(row), " ")
		if trimmed != "" {
			buf.WriteString(trimmed)
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func sampleTree() *node {
	return &node{
		val: "1",
		left: &node{
			val: "2",
			left: &node{
				val: "10",
			},
			right: &node{val: "2"},
		},
		right: &node{
			val:   "3",
			left:  &node{val: "2"},
			right: &node{val: "1"},
		},
	}
}

func main() {
	root := sampleTree()
	fmt.Print(draw(root))
}
