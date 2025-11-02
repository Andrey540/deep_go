package main

import (
	"encoding/json"
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const nameLength = 42

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	var nameArr [nameLength]byte
	copy(nameArr[:], name)
	return func(person *GamePerson) {
		person.PersonName = PersonName{
			Data:   nameArr,
			Length: byte(len(name)),
		}
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.Coordinates = PersonCoordinates{
			X: int32(x),
			Y: int32(y),
			Z: int32(z),
		}
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.Wealth = uint32(gold<<1) | person.Wealth&0x1
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.PersonStrength = setValue(uint32(mana), person.PersonStrength, 10, 22)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.PersonStrength = setValue(uint32(health), person.PersonStrength, 10, 12)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.Person = setValue(uint8(respect), person.Person, 4, 4)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.PersonStrength = setValue(uint32(strength), person.PersonStrength, 4, 8)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.PersonStrength = setValue(uint32(experience), person.PersonStrength, 4, 4)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.PersonStrength = setValue(uint32(level), person.PersonStrength, 4, 0)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.Wealth = person.Wealth | 0x1
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.Person = setValue(uint8(1), person.Person, 1, 3)
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.Person = setValue(uint8(1), person.Person, 1, 2)
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.Person = setValue(uint8(personType), person.Person, 2, 0)
	}
}

func getBits[T uint32 | uint8](count, offset int) T {
	var val T
	for i := 0; i < count; i++ {
		val |= 1 << i
	}
	return val << offset
}

func clearBits[T uint32 | uint8](count, offset int) T {
	return ^getBits[T](count, offset)
}

func setValue[T uint32 | uint8](value, storage T, count, offset int) T {
	return value<<offset | storage&clearBits[T](count, offset)
}

func getValue[T uint32 | uint8](value T, count, offset int) T {
	return value & getBits[T](count, offset) >> offset
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type PersonName struct {
	Data   [nameLength]byte `json:"Data"`
	Length byte             `json:"Length"`
}

func (n *PersonName) Name() string {
	return string(n.Data[:n.Length])
}

type PersonCoordinates struct {
	X int32 `json:"X"`
	Y int32 `json:"Y"`
	Z int32 `json:"Z"`
}

type GamePerson struct {
	Coordinates PersonCoordinates `json:"Coordinates"`
	PersonName  PersonName        `json:"PersonName"`

	/*
		respect    4 bits
		hasGun     1 bit
		hasFamily  1 bit
		personType 2 bits
	*/
	Person uint8 `json:"Person"`

	/*
		gold     31 bits
		hasHouse 1 bit
	*/
	Wealth uint32 `json:"Wealth"`

	/*
		mana       10 bits
		health     10 bits
		strength   4 bits
		experience 4 bits
		level      4 bits
	*/
	PersonStrength uint32 `json:"PersonStrength"`
}

func NewGamePerson(options ...Option) GamePerson {
	person := &GamePerson{}
	for _, option := range options {
		option(person)
	}
	return *person
}

func (p *GamePerson) Name() string {
	return p.PersonName.Name()
}

func (p *GamePerson) X() int {
	return int(p.Coordinates.X)
}

func (p *GamePerson) Y() int {
	return int(p.Coordinates.Y)
}

func (p *GamePerson) Z() int {
	return int(p.Coordinates.Z)
}

func (p *GamePerson) Gold() int {
	return int(p.Wealth >> 1)
}

func (p *GamePerson) Mana() int {
	return int(getValue(p.PersonStrength, 10, 22))
}

func (p *GamePerson) Health() int {
	return int(getValue(p.PersonStrength, 10, 12))
}

func (p *GamePerson) Respect() int {
	return int(getValue(p.Person, 4, 4))
}

func (p *GamePerson) Strength() int {
	return int(getValue(p.PersonStrength, 4, 8))
}

func (p *GamePerson) Experience() int {
	return int(getValue(p.PersonStrength, 4, 4))
}

func (p *GamePerson) Level() int {
	return int(getValue(p.PersonStrength, 4, 0))
}

func (p *GamePerson) HasHouse() bool {
	return p.Wealth&0x1 == 0x1
}

func (p *GamePerson) HasGun() bool {
	return getValue(p.Person, 1, 3) == 0x1
}

func (p *GamePerson) HasFamily() bool {
	return getValue(p.Person, 1, 2) == 0x1
}

func (p *GamePerson) Type() int {
	return int(getValue(p.Person, 2, 0))
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())

	serialized, err := json.Marshal(person)
	assert.Nil(t, err)

	restoredPerson := GamePerson{}
	err = json.Unmarshal(serialized, &restoredPerson)
	assert.Nil(t, err)
	assert.Equal(t, person, restoredPerson)
}
