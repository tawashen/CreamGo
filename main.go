package go
import (
	"fmt"
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	playerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#000000"))
	treeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	waterStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))
	mtStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#8B4513"))
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error %v", err)
	}
}

type model struct {
	playerX int
	playerY int
	mapData [][]rune
	width int
	height int
}

func initialModel() model {
	m := model {
		playerX: 10,
		playerY: 10,
		width: 19,
		height: 19,
	}
	m.generateMap()
	return m
}

func (m *model) generateMap() {
	tiles := []rune{'T', '~', '^', ' ', ' ', ' '}
	m.mapData = make([][]rune, m.height)
	for y := 0; y < m.height; y++ {
		row := make([][]rune, m.width)
		for x := 0; x < m.width; x++ {
			row[x] := tiles[rand.Intn(len(tiles))]
		}
		m.mapData[y] = row
	}
}

func (m model) Init() tea.Cmd {
	retrun nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KsyMsg:
			switch msg.String() {
			case "ctrl+c", "q", "esc":
				return m, tea.Quit
			case "up":
				if m.playerY > 0 {
					m.playerY--
				}
			case "down":
				if m.playerY < m.height-1 {
					m.playerY++
				}
			case "left":
				if m.playerX > 0 {
					m.playerX--
				}
			case "right":
				if m.playerX < m.width -1 {
					m.playerX++
				}
		}
	}
 	return m, nil
}

func (m model) View() string {
	var s strings.Builder

	s.WriteString(playerStyle.Render("くりぃむ大戦 \n\n"))

    for y := 0; y < m.height; y++ {
    	for x := 0; x < m.width; x++ {
		if x == m.playerX && y == m.playerY {
			s.WriteString(playerStyle.Render("🙋"))
			continue
		}

			// マップチップの描画
			char := m.mapData[y][x]
			switch char {
			case 'T':
				s.WriteString(treeStyle.Render("🌲"))
			case '~':
				s.WriteString(waterStyle.Render("🌊"))
			case '^':
				s.WriteString(mtStyle.Render("🌋"))
			default:
				s.WriteString("  ") // 半角スペース2つ（全角1マス分）
			}
		}
		s.WriteString("\n")
	}

	s.WriteString(fmt.Sprintf("\n座標: (%d, %d)", m.playerX, m.playerY))
	return s.String()
	}


 
type Monster struct {
	ID   int
	Name string
	HP int
	MP int
	Special []string
	DOt string // ANSIエスケープシーケンス済みの文字列
}

func Encount(num int) 


