package state

import "fmt"

// State определяет интерфейс состояния
type State interface {
	Play(p *Player)
	Stop(p *Player)
	Pause(p *Player)
}

// Player представляет собой контекст, который изменяет свое состояние
type Player struct {
	state State
}

func NewPlayer() *Player {
	return &Player{state: &StoppedState{}}
}

func (p *Player) SetState(state State) {
	p.state = state
}

func (p *Player) Play() {
	p.state.Play(p)
}

func (p *Player) Stop() {
	p.state.Stop(p)
}

func (p *Player) Pause() {
	p.state.Pause(p)
}

// PlayingState реализует состояние воспроизведения
type PlayingState struct{}

func (s *PlayingState) Play(p *Player) {
	fmt.Println("Already playing")
}

func (s *PlayingState) Stop(p *Player) {
	fmt.Println("Stopping the music")
	p.SetState(&StoppedState{})
}

func (s *PlayingState) Pause(p *Player) {
	fmt.Println("Pausing the music")
	p.SetState(&PausedState{})
}

// StoppedState реализует состояние остановки
type StoppedState struct{}

func (s *StoppedState) Play(p *Player) {
	fmt.Println("Playing the music")
	p.SetState(&PlayingState{})
}

func (s *StoppedState) Stop(p *Player) {
	fmt.Println("Already stopped")
}

func (s *StoppedState) Pause(p *Player) {
	fmt.Println("Can't pause. The music is stopped")
}

// PausedState реализует состояние паузы
type PausedState struct{}

func (s *PausedState) Play(p *Player) {
	fmt.Println("Resuming music")
	p.SetState(&PlayingState{})
}

func (s *PausedState) Stop(p *Player) {
	fmt.Println("Stopping the music from pause")
	p.SetState(&StoppedState{})
}

func (s *PausedState) Pause(p *Player) {
	fmt.Println("Already on pause")
}

// Основная функция
func main() {
	player := NewPlayer()

	// Попытки управления плеером
	player.Play()
	player.Pause()
	player.Play()
	player.Stop()
	player.Play()
	player.Stop()
}
