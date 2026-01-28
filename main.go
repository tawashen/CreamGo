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
	Special []string
	Dot     string // ANSIエスケープシーケンス済みの文字列
}

type model struct {
	playerX int
	playerY int
	mapData [][]rune
	width   int
	height  int
	scene   string
	turn string
	action string
}
// 修正版: アイテム機能を使うためにitemsフィールドが必要です
/*
type model struct {
	playerX int
	playerY int
	mapData [][]rune
	width   int
	height  int
	scene   string
	turn    string
	action  string
	items   []Item  // アイテムリストを追加
}
*/

func initialModel() model {
	m := model{
		playerX: 10,
		playerY: 10,
		width:   19,
		height:  19,
		scene:   "field",   // カンマ追加
		turn:    "player",  // カンマ追加
		action:  "menu",    // カンマ追加（最後のフィールドでも推奨）
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

func (m *model) UseItem(Item string) model {
	switch Item.Kind {
	case "Heal":
		
	}

}
// 修正版: 複数の問題があります
/*
1. 引数の型: string → Item
2. 引数名: Item → item (大文字で始まる変数名は推奨されません)
3. stringにはKindフィールドがありません
4. 戻り値が必要です

func (m *model) UseItem(item Item) model {
	switch item.Kind {
	case "Heal":
		// HP回復処理など
	}
	return *m  // modelを返す
}
*/


func (m *model) generateMap() {
	tiles := []rune{'T', '~', '^', ' ', ' ', ' '}
	m.mapData = make([][]rune, m.height)
	for y := 0; y < m.height; y++ {
		row := make([]rune, m.width)
		for x := 0; x < m.width; x++ {
			row[x] = tiles[rand.Intn(len(tiles))]
		}
		m.mapData[y] = row
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.scene == "field" {
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
			if m.playerX < m.width-1 {
				m.playerX++
			}
		case "b":
			m.scene = "battle"
		}
		}

		if m.scene == "battle" && m.turn == "player" {
			switch m.action {
			case "Menu":
			switch msg.String() {
			case "1":
				m.action = "Attack"
			case "2":
				m.action = "SelectItem"
			case "3":
				m.action = "SelecSpecial"
			case "4":
				m.action = "Escape"
			}

			case "SelectItem":
				idex, err := strconv.Atoi(msg.String())
				if err == nil && idex >= 1 && idx <= len(m.items) {
					SelectedItem := m.items[idx-1]
					m.action = "UseItem"
					m.UseItem(SelectedItem)
				}
			// 修正版: 複数の変数名の間違いがあります
			/*
			case "SelectItem":
				index, err := strconv.Atoi(msg.String())  // idex → index
				if err == nil && index >= 1 && index <= len(m.items) {  // idx → index
					selectedItem := m.items[index-1]  // SelectedItem → selectedItem
					m.action = "UseItem"
					m.UseItem(selectedItem)
				}
			// 注意: m.itemsフィールドをmodelに追加する必要があります
			*/
			}
		}
	}
		return m, nil
}


func (m model) View() string {
	var s strings.Builder

	s.WriteString(playerStyle.Render("くりぃむ大戦 \n\n"))

	if m.scene == "field" {
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
	}

	if m.scene == "battle" {
		s.WriteString(playerStyle.Render(monsterList[0].Dot))
		s.WriteString(playerStyle.Render(monsterList[0].Name))
		s.WriteString("\n")
	}

	s.WriteString(fmt.Sprintf("\n座標: (%d, %d)", m.playerX, m.playerY))
	return s.String()
}


func PickMonster(num int) Monster {
	return monsterList[num]
}

func (m *model) Battle() model {
	m.scene = "battle"
	monster := PickMonster(0)

}
// 修正版: 未使用変数と戻り値の問題があります
/*
func (m *model) Battle() model {
	m.scene = "battle"
	monster := PickMonster(0)
	// monsterを使用するか、_ := PickMonster(0) にする
	// 例: m.currentMonster = monster (currentMonsterフィールドを追加する場合)
	
	return *m  // 戻り値を返す必要があります
}
*/
