package command

import "fmt"

// Command определяет интерфейс команды
type Command interface {
	Execute()
}

// Light представляет собой получателя
type Light struct {
	status bool
}

func (l *Light) On() {
	l.status = true
	fmt.Println("Light is On")
}

func (l *Light) Off() {
	l.status = false
	fmt.Println("Light is Off")
}

// LightOnCommand представляет конкретную команду для включения света
type LightOnCommand struct {
	light *Light
}

func NewLightOnCommand(light *Light) *LightOnCommand {
	return &LightOnCommand{light: light}
}

func (c *LightOnCommand) Execute() {
	c.light.On()
}

// LightOffCommand представляет конкретную команду для выключения света
type LightOffCommand struct {
	light *Light
}

func NewLightOffCommand(light *Light) *LightOffCommand {
	return &LightOffCommand{light: light}
}

func (c *LightOffCommand) Execute() {
	c.light.Off()
}

// RemoteControl представляет вызывателя
type RemoteControl struct {
	command Command
}

func (r *RemoteControl) SetCommand(command Command) {
	r.command = command
}

func (r *RemoteControl) PressButton() {
	r.command.Execute()
}

// Основная функция
func main() {
	// Создаем получателя
	light := &Light{status: false}

	// Создаем команды
	lightOn := NewLightOnCommand(light)
	lightOff := NewLightOffCommand(light)

	// Создаем вызывателя
	remote := &RemoteControl{}

	// Включаем свет
	remote.SetCommand(lightOn)
	remote.PressButton()

	// Выключаем свет
	remote.SetCommand(lightOff)
	remote.PressButton()
}
