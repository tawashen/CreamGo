package main

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


type Weapon struct {
	Name  string
	Power int
	Value int
}

type Armor struct {
	Name    string
	Defense int
	Value   int
}

type model struct {
	PlayerX int
	PlayerY int
	HP int
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
	Messages []string  // メッセージ履歴用（推奨）
	Msg string        // 現在のメッセージ用
}



func initialModel() model {
	m := model{
		PlayerX: 10,
		PlayerY: 10,
		HP: 100,
		Attack: 5,
		Defend: 5,
		Weapon: nil,
		Armor: nil,
		Gold: 0,
		Items: []Item{},
		Status: []string{},
		MapData: [][]rune{},
		Width:   19,
		Height:  19,
		Scene:   "field",   // カンマ追加
		Turn:    "player",  // カンマ追加
		Action:  "menu",    // カンマ追加（最後のフィールドでも推奨）
		CurrentMonster: nil,
		Msg: "",
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

func (m *model) UseItem(item Item) {
	switch item.Kind {
	case "Heal":
		m.Msg = fmt.Sprintf("%sを使った！あなたの体力は%d回復した", item.Name, item.Power)
		m.HP += item.Power
		m.Action = "menu"
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
			m.CurrentMonster = &monsterList[0]
			m.Action = "menu"  // バトル開始時はメニュー状態に
			m.Turn = "player"  // プレイヤーのターンに設定
		}
		}

		if m.Scene == "battle" && m.Turn == "player" {
			switch m.Action {
			case "menu":  // "Menu" → "menu" (小文字に統一)
				switch msg.String() {
				case "1":
					m.Action = "Attack"  // 攻撃を選択
				case "2":
					m.Action = "selectitem"
				case "3":
					m.Action = "selectspecial"
				case "4":
					m.Action = "escape"
				}
			
			case "Attack":
				// 攻撃処理を実行
				m.Battle()
				m.Action = "menu"  // 処理後はメニューに戻る
			
			case "selectitem":
				index, err := strconv.Atoi(msg.String())
				if err == nil && index >= 1 && index <= len(m.Items) {
					selectedItem := m.Items[index-1]
					m.UseItem(selectedItem)
					// UseItem内でm.Action = "menu"が設定される
				}
			
			case "escape":
				m.Scene = "field"  // フィールドに戻る
				m.Action = "menu"
			}
		}
	}
		return m, nil
}


func (m model) View() string {
	var s strings.Builder

	s.WriteString(playerStyle.Render("くりぃむ大戦 \n\n"))

	if m.Scene == "field" {
		s.WriteString("\n")
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

	if m.Scene == "battle" {
		if m.CurrentMonster != nil {
			// モンスター情報表示
			s.WriteString("\n")
			s.WriteString(playerStyle.Render(m.CurrentMonster.Dot))
			s.WriteString("\n")
			
			text := fmt.Sprintf("%sが現れた！ (HP: %d)", 
				m.CurrentMonster.Name, m.CurrentMonster.HP)
			s.WriteString(playerStyle.Render(text))
			s.WriteString("\n\n")
			
			// メッセージ表示
			if m.Msg != "" {
				s.WriteString(playerStyle.Render(m.Msg))
				s.WriteString("\n\n")
			}
			
			// アクションに応じたメニュー表示
			switch m.Action {
			case "menu":
				s.WriteString("どうする？\n")
				s.WriteString("1. 攻撃\n")
				s.WriteString("2. アイテム\n")
				s.WriteString("3. 特技\n")
				s.WriteString("4. 逃げる\n")
			case "selectitem":
				s.WriteString("アイテムを選んでください:\n")
				for i, item := range m.Items {
					s.WriteString(fmt.Sprintf("%d. %s\n", i+1, item.Name))
				}
			}
		}
	}

	s.WriteString(fmt.Sprintf("\n座標: (%d, %d)", m.PlayerX, m.PlayerY))
	return s.String()
}


//func PickMonster(num int) Monster {
//	return monsterList[num]
//}

func (m *model) Battle() {
	switch m.Action {
	case "Attack":

		if m.CurrentMonster == nil {
			m.Msg = "敵がいません"
		}

		monster := m.CurrentMonster


		weaponPower := 0
		if m.Weapon != nil {  // nilチェック追加
			weaponPower = m.Weapon.Power
		}
		
		damage := (m.Attack + weaponPower) - monster.Defend
		if damage <= 0 {
			damage = 1  // 最低1ダメージ
		}
		
		monster.HP -= damage  // HPを減らす
		msg := fmt.Sprintf("攻撃！ %sに%dのダメージ！", monster.Name, damage)
		m.Msg = msg
		
		// メッセージ履歴に追加する場合:
		// m.Messages = append(m.Messages, msg)
		
		//m.Action = "menu"  // アクションをリセット
	}
	//return m
}

/*
		damage := (m.Attack + m.Weapon.Power) - monster.Defend
		if damage <= 0 {
			damage = 0
		}

		msg := fmt.Sprintf("攻撃！ %sに%dのダメージ！\n", monster.Name, damage)
		m.Msg = msg
	}
	//return m
}
	*/

// 修正版: 複数の構文エラーがあります
/*
func (m *model) Battle() tea.Model {
	switch m.Action {  // swtich → switch
	case "Attack":     // コロンを削除
		if m.CurrentMonster == nil {  // nilチェック追加
			m.Msg = "敵がいません"
			return m
		}
		
		monster := m.CurrentMonster
		weaponPower := 0
		if m.Weapon != nil {  // nilチェック追加
			weaponPower = m.Weapon.Power
		}
		
		damage := (m.Attack + weaponPower) - monster.Defend
		if damage <= 0 {
			damage = 1  // 最低1ダメージ
		}
		
		monster.HP -= damage  // HPを減らす
		msg := fmt.Sprintf("攻撃！ %sに%dのダメージ！", monster.Name, damage)
		m.Msg = msg
		
		// メッセージ履歴に追加する場合:
		// m.Messages = append(m.Messages, msg)
		
		m.Action = "menu"  // アクションをリセット
	}
	return m
}
*/

