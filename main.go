package creamgo

import (
	"fmt"
	"math/rand"
	"strings"
	"strconv"

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


type Monster struct {
	ID      int
	Name    string
	HP      int
	MP      int
	Attack int
	Defend int
	Special []string
	Dot     string // ANSIエスケープシーケンス済みの文字列
}

type model struct {
	PlayerX int
	PlayerY int
	Attack int
	Defend int
	Weapon *Weapon
	Armor *Armor
	Gold int
	Items []Item
	Status []string
	MapData [][]rune
	Width   int
	Height  int
	Scene   string
	Turn string
	Action string
	CurrentMonster *Monster
	Msg string
}


func initialModel() model {
	m := model{
		PlayerX: 10,
		PlayerY: 10,
		Attack: 5,
		Defend: 5,
		Weapon: nil,
		Armor: nil,
		Gold: 0,
		Items: []{},
		Status: []{},
		Mapdata: [][]{},
		Width:   19,
		Height:  19,
		Scene:   "field",   // カンマ追加
		Turn:    "player",  // カンマ追加
		Action:  "menu",    // カンマ追加（最後のフィールドでも推奨）
		CurrentMonster: nil,
		Msg "",
	}
	m.generateMap()
	return m
}

type Item struct {
	Name string
	Kind string
	Power int
	Value int
}

func (m *model) UseItem(item Item) model {
	switch item.Kind {
	case "Heal":
		
	}

}



func (m *model) generateMap() {
	tiles := []rune{'T', '~', '^', ' ', ' ', ' '}
	m.MapData = make([][]rune, m.Height)
	for y := 0; y < m.Height; y++ {
		row := make([]rune, m.Width)
		for x := 0; x < m.Width; x++ {
			row[x] = tiles[rand.Intn(len(tiles))]
		}
		m.MapData[y] = row
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.Scene == "field" {
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up":
			if m.PlayerY > 0 {
				m.PlayerY--
			}
		case "down":
			if m.PlayerY < m.Height-1 {
				m.PlayerY++
			}
		case "left":
			if m.PlayerX > 0 {
				m.PlayerX--
			}
		case "right":
			if m.PlayerX < m.Width-1 {
				m.PlayerX++
			}
		case "b":
			m.Scene = "battle"
			m.CurrentMonster = monsterList[0]
		}
		}

		if m.Scene == "battle" && m.Turn == "player" {
			switch m.Action {
			case "Menu":
			switch msg.String() {
			case "1":
				m.Action = "Attack"
			case "2":
				m.Action = "SelectItem"
			case "3":
				m.Action = "SelecSpecial"
			case "4":
				m.Action = "Escape"
			}

			case "SelectItem":
				index, err := strconv.Atoi(msg.String())
				if err == nil && index >= 1 && idx <= len(m.Items) {
					SelectedItem := m.Items[index-1]
					m.Action = "UseItem"
					m.UseItem(SelectedItem)
				}
			}
		}
	}
		return m, nil
}


func (m model) View() string {
	var s strings.Builder

	s.WriteString(playerStyle.Render("くりぃむ大戦 \n\n"))

	if m.Scene == "field" {
		for y := 0; y < m.Height; y++ {
			for x := 0; x < m.Width; x++ {
				if x == m.PlayerX && y == m.PlayerY {
					s.WriteString(playerStyle.Render("🙋"))
					continue
				}

				// マップチップの描画
				char := m.MapData[y][x]
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
	}

	if m.scene == "battle" {
		s.WriteString(playerStyle.Render(monsterList[0].Dot))
		s.WriteString(playerStyle.Render(monsterList[0].Name))
		s.WriteString("\n")
		s.WriteString(playerStyle.Render(m.Msg))
		s.WriteString("\n")
	}

	s.WriteString(fmt.Sprintf("\n座標: (%d, %d)", m.PlayerX, m.PlayerY))
	return s.String()
}


//func PickMonster(num int) Monster {
//	return monsterList[num]
//}

func (m *model) Battle() tea.Model {
	swtich m.Action {
	case "Attack" :
		monster := m.CurrentMonster
		damage := (m.Attack + m.Weapon.Power) - monster.Defend
		if damage <= 0 {
			damage = 0
		}

		msg := fmt.Sprintf("攻撃！ %sに%dのダメージ！\n", monster.Name, damage)
		m.Msg = msg
	}
	return m
}

