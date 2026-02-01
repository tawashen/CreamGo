package main

import (
	"fmt"
	"math/rand"
	"strings"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// カスタムメッセージ型（時間待ち用）
type DelayMsg struct{}

// より詳細なメッセージ型の例
/*
// 1. 連鎖演出用メッセージ
type AttackAnimationMsg struct{}
type DamageCalculationMsg struct{}
type BattleEndMsg struct{}

// 2. 条件分岐用メッセージ（データ付き）
type VictoryMsg struct {
	ExpGained  int
	GoldGained int
}
type LowHPWarningMsg struct{}
type CriticalHitMsg struct{}
type NormalAttackEndMsg struct{}

// 3. アニメーション用メッセージ
type MonsterAnimationMsg struct{}
type StartAttackAnimationMsg struct{}

// 4. 音声・効果用メッセージ
type SoundEffectMsg struct {
	Sound string  // "attack.wav", "victory.wav" など
}
type ScreenShakeMsg struct{}
type ShowDamageMsg struct {
	Damage int
}

// 5. 複雑な演出用メッセージ
type PlayAttackSoundMsg struct{}
type ShowDamageNumbersMsg struct {
	Damage int
}
type AttackSequenceEndMsg struct{}

// 6. モンスターアニメーション用の拡張
type Monster struct {
	ID      int
	Name    string
	HP      int
	MaxHP   int     // 最大HP追加
	MP      int
	Attack  int
	Defend  int
	Special []string
	Dot     string
	AnimationDots []string  // アニメーション用の複数ドット
}

// 7. モデルにアニメーション関連フィールド追加
type model struct {
	// ... 既存フィールド
	animationFrame int      // 現在のアニメーションフレーム
	monsterDots    []string // モンスターのアニメーションドット配列
}

// 実際の音声ライブラリ使用例（疑似コード）
/*
import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/wav"
	"github.com/faiface/beep/speaker"
)

func playSound(filename string) {
	// WAVファイルを読み込んで再生
	// 実装は音声ライブラリのドキュメント参照
}
*/
*/

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
		Items: []Item{
			{Name: "薬草", Kind: "Heal", Power: 10, Value: 8},
		},
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



func (m *model) UseItem(item Item) tea.Cmd {
	switch item.Kind {
	case "Heal":
		m.Msg = fmt.Sprintf("%sを使った！あなたの体力は%d回復した\n\n\n", item.Name, item.Power)
		m.HP += item.Power
		m.Action = "waiting"

		return tea.Tick(time.Second, func(time.Time) tea.Msg {
			return DelayMsg{}
		})
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
	case DelayMsg:
		// 時間待ち後の処理
		if m.Scene == "battle" && m.Action == "waiting" {
			m.Action = "menu"  // メニューに戻る
			m.Msg = ""         // メッセージをクリア
			
			// 敵が倒れた場合の処理
			if m.CurrentMonster != nil && m.CurrentMonster.HP <= 0 {
				m.Msg = fmt.Sprintf("%sを倒した！\n\n\n", m.CurrentMonster.Name)
				// さらに1秒後にフィールドに戻る場合
				 m.Action = "victory"
				 return m, tea.Tick(time.Second, func(time.Time) tea.Msg {
				     return DelayMsg{}
				 })
			}
		}

		if m.Scene == "battle" && m.Action == "victory" {
			m.Scene = "field"
		}

		if m.Scene == 

		return m, nil
		
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
			case "menu":
				switch msg.String() {
				case "1":
					// 攻撃を即座に実行
					return m, m.Battle()
				case "2":
					m.Action = "selectitem"
				case "3":
					m.Action = "selectspecial"
				case "4":
					m.Action = "escape"
				}
			
			case "selectitem":
				index, err := strconv.Atoi(msg.String())
				if err == nil && index >= 1 && index <= len(m.Items) {
					selectedItem := m.Items[index-1]
					m.UseItem(selectedItem)
				}
			
			case "escape":
				m.Scene = "field"
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
			case "selectItem":
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

func (m *model) Battle() tea.Cmd {
	if m.CurrentMonster == nil {
		m.Msg = "敵がいません"
		return nil
	}

	monster := m.CurrentMonster
	weaponPower := 0
	if m.Weapon != nil {
		weaponPower = m.Weapon.Power
	}
	
	damage := (m.Attack + weaponPower) - monster.Defend
	if damage <= 0 {
		damage = 1
	}
	
	monster.HP -= damage
	msg := fmt.Sprintf("攻撃！ %sに%dのダメージ！\n\n\n", monster.Name, damage)
	m.Msg = msg
	m.Action = "waiting"  // 待機状態に変更
	
	// 1秒後にDelayMsgを送信
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return DelayMsg{}
	})
	
	// 実用的な使用例集:
	/*
	// 1. 連鎖的な演出（Update内でmsg送信を連鎖）
	// Battle関数で最初のメッセージ送信
	return tea.Tick(500*time.Millisecond, func(time.Time) tea.Msg {
		return AttackAnimationMsg{}
	})
	
	// Update関数で連鎖処理
	case AttackAnimationMsg:
		m.Msg = "攻撃！"
		return m, tea.Tick(300*time.Millisecond, func(time.Time) tea.Msg {
			return DamageCalculationMsg{}
		})
	
	case DamageCalculationMsg:
		m.Msg = fmt.Sprintf("%dのダメージ！", damage)
		return m, tea.Tick(800*time.Millisecond, func(time.Time) tea.Msg {
			return BattleEndMsg{}
		})
	
	case BattleEndMsg:
		m.Action = "menu"
		m.Msg = ""
		return m, nil
	
	// 2. 条件分岐でメッセージを変化（実践的！）
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		if monster.HP <= 0 {
			return VictoryMsg{ExpGained: 100, GoldGained: 50}
		} else if monster.HP < monster.MaxHP/4 {
			return LowHPWarningMsg{}  // 瀕死警告
		} else if damage > 50 {
			return CriticalHitMsg{}   // クリティカル
		}
		return NormalAttackEndMsg{}   // 通常攻撃終了
	})
	
	// 3. tea.Everyでアニメーション（条件付き制御）
	// モンスターのアニメーション開始
	if m.Action == "monster_animation" {
		return tea.Every(200*time.Millisecond, func(time.Time) tea.Msg {
			return MonsterAnimationMsg{}
		})
	}
	
	// Update関数でアニメーション制御
	case MonsterAnimationMsg:
		if m.Action != "monster_animation" {
			// アニメーション停止条件
			return m, nil
		}
		// アニメーションフレーム更新
		m.animationFrame = (m.animationFrame + 1) % len(m.monsterDots)
		m.CurrentMonster.Dot = m.monsterDots[m.animationFrame]
		return m, nil  // tea.Everyが自動で次のフレームを送信
	
	// アニメーション停止
	case SomeOtherMsg:
		m.Action = "menu"  // これでMonsterAnimationMsgが停止
	
	// 4. tea.Batchの実用例
	// 複数の効果を同時実行
	return tea.Batch(
		// 効果音再生（500ms後）
		tea.Tick(500*time.Millisecond, func(time.Time) tea.Msg {
			return SoundEffectMsg{Sound: "attack.wav"}
		}),
		// 画面シェイク（100ms後）
		tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
			return ScreenShakeMsg{}
		}),
		// ダメージ表示（1秒後）
		tea.Tick(time.Second, func(time.Time) tea.Msg {
			return ShowDamageMsg{Damage: damage}
		}),
	)
	
	// 5. 音声について（重要！）
	// Bubble Tea自体は音を出せませんが、メッセージで音声再生を指示できます
	case SoundEffectMsg:
		// 外部ライブラリで音声再生
		// 例: github.com/faiface/beep
		// playSound(msg.Sound)
		return m, nil
	
	// 実際の音声実装例（疑似コード）
	/*
	import "github.com/faiface/beep/speaker"
	
	case SoundEffectMsg:
		switch msg.Sound {
		case "attack.wav":
			// 攻撃音再生
			playAttackSound()
		case "victory.wav":
			// 勝利音再生
			playVictorySound()
		}
		return m, nil
	*/
	
	// 6. 複雑な演出の組み合わせ例
	// 攻撃 → アニメーション → 音 → ダメージ表示 → 結果
	return tea.Batch(
		tea.Tick(0, func(time.Time) tea.Msg {
			return StartAttackAnimationMsg{}
		}),
		tea.Tick(300*time.Millisecond, func(time.Time) tea.Msg {
			return PlayAttackSoundMsg{}
		}),
		tea.Tick(600*time.Millisecond, func(time.Time) tea.Msg {
			return ShowDamageNumbersMsg{Damage: damage}
		}),
		tea.Tick(1200*time.Millisecond, func(time.Time) tea.Msg {
			return AttackSequenceEndMsg{}
		}),
	)
	*/
}
