package main

import (
	"strings"
)

type RoomID string

const (
	roomKitchen  RoomID = "кухня"
	roomCorridor RoomID = "коридор"
	roomRoom     RoomID = "комната"
	roomStreet   RoomID = "улица"
	roomHome     RoomID = "домой"
)

type Player struct {
	location     RoomID
	inventory    map[string]bool
	backpackOn   bool
	doorUnlocked bool
	roomState    roomState
}

type roomState struct {
	hasBackpack bool
	hasKeys     bool
	hasNotes    bool
}

var player Player

func initGame() {
	player = Player{
		location:     roomKitchen,
		inventory:    make(map[string]bool),
		backpackOn:   false,
		doorUnlocked: false,
		roomState: roomState{
			hasBackpack: true,
			hasKeys:     true,
			hasNotes:    true,
		},
	}
}

func handleCommand(command string) string {
	parts := splitCommand(command)
	if len(parts) == 0 {
		return "неизвестная команда"
	}

	switch parts[0] {
	case "осмотреться":
		return look()
	case "идти":
		if len(parts) < 2 {
			return "неизвестная команда"
		}
		return move(parts[1])
	case "взять":
		if len(parts) < 2 {
			return "неизвестная команда"
		}
		return take(parts[1])
	case "надеть":
		if len(parts) < 2 {
			return "неизвестная команда"
		}
		if parts[1] != "рюкзак" {
			return "неизвестная команда"
		}
		return wearBackpack()
	case "применить":
		if len(parts) < 3 {
			return "неизвестная команда"
		}
		return use(parts[1], parts[2])
	default:
		return "неизвестная команда"
	}
}

func splitCommand(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, " ")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func look() string {
	switch player.location {
	case roomKitchen:
		if player.backpackOn && player.inventory["ключи"] && player.inventory["конспекты"] {
			return "ты находишься на кухне, на столе: чай, надо идти в универ. можно пройти - коридор"
		}
		return "ты находишься на кухне, на столе: чай, надо собрать рюкзак и идти в универ. можно пройти - коридор"
	case roomCorridor:
		return "ничего интересного. можно пройти - кухня, комната, улица"
	case roomRoom:
		var table []string
		if player.roomState.hasKeys {
			table = append(table, "ключи")
		}
		if player.roomState.hasNotes {
			table = append(table, "конспекты")
		}

		hasAnyOnTable := len(table) > 0
		hasBackpack := player.roomState.hasBackpack

		if !hasAnyOnTable && !hasBackpack {
			return "пустая комната. можно пройти - коридор"
		}

		var b strings.Builder
		if hasAnyOnTable {
			b.WriteString("на столе: ")
			b.WriteString(strings.Join(table, ", "))
			if hasBackpack {
				b.WriteString(", ")
			} else {
				b.WriteString(". ")
			}
		}
		if hasBackpack {
			b.WriteString("на стуле: рюкзак. ")
		}
		b.WriteString("можно пройти - коридор")
		return b.String()
	case roomStreet:
		return "на улице весна. можно пройти - домой"
	default:
		return "неизвестная команда"
	}
}

func move(dst string) string {
	switch player.location {
	case roomKitchen:
		if dst == "коридор" {
			player.location = roomCorridor
			return look()
		}
		return "нет пути в " + dst
	case roomCorridor:
		switch dst {
		case "кухня":
			player.location = roomKitchen
			return "кухня, ничего интересного. можно пройти - коридор"
		case "комната":
			player.location = roomRoom
			return "ты в своей комнате. можно пройти - коридор"
		case "улица":
			if !player.doorUnlocked {
				return "дверь закрыта"
			}
			player.location = roomStreet
			return look()
		default:
			return "нет пути в " + dst
		}
	case roomRoom:
		if dst == "коридор" {
			player.location = roomCorridor
			return look()
		}
		return "нет пути в " + dst
	case roomStreet:
		if dst == "домой" {
			player.location = roomCorridor
			return look()
		}
		return "нет пути в " + dst
	default:
		return "нет пути в " + dst
	}
}

func take(item string) string {
	if player.location != roomRoom {
		// нет такого (в этой комнате этого предмета нет)
		return "нет такого"
	}
	if !player.backpackOn {
		return "некуда класть"
	}

	switch item {
	case "ключи":
		if !player.roomState.hasKeys {
			return "нет такого"
		}
		player.roomState.hasKeys = false
		player.inventory["ключи"] = true
		return "предмет добавлен в инвентарь: ключи"
	case "конспекты":
		if !player.roomState.hasNotes {
			return "нет такого"
		}
		player.roomState.hasNotes = false
		player.inventory["конспекты"] = true
		return "предмет добавлен в инвентарь: конспекты"
	default:
		return "нет такого"
	}
}

func wearBackpack() string {
	if player.location != roomRoom {
		return "неизвестная команда"
	}
	if !player.roomState.hasBackpack {
		// уже взяли раньше
		return "неизвестная команда"
	}
	player.roomState.hasBackpack = false
	player.backpackOn = true
	return "вы надели: рюкзак"
}

func use(item, target string) string {
	if !player.inventory[item] {
		return "нет предмета в инвентаре - " + item
	}
	if item == "ключи" && target == "дверь" &&
		(player.location == roomCorridor || player.location == roomStreet) {

		if player.doorUnlocked {
			player.doorUnlocked = false
			return "дверь закрыта"
		}
		player.doorUnlocked = true
		return "дверь открыта"
	}
	return "не к чему применить"
}
