package main

// ItemType represents different types of items
type ItemType int

const (
	ItemGold ItemType = iota
	ItemPotion
	ItemWeapon
	ItemArmor
	ItemTreasure
	ItemKey
)

// Item represents an item in the game
type Item struct {
	X, Y        int      // Position in the dungeon
	Type        ItemType // Type of item
	Name        string   // Name of the item
	Description string   // Description of the item
	Value       int      // Value (gold, healing amount, damage, etc.)
	Symbol      rune     // Symbol to display on the map
	Collected   bool     // Whether the item has been collected
}

// NewHealthPotion creates a new health potion
func NewHealthPotion(x, y int) Item {
	return Item{
		X:          x,
		Y:          y,
		Type:       ItemPotion,
		Name:       "Health Potion",
		Description: "Restores 10 health points",
		Value:      10,
		Symbol:     '!',
		Collected:  false,
	}
}

// NewWeapon creates a new weapon
func NewWeapon(x, y int, name string, damage int) Item {
	return Item{
		X:          x,
		Y:          y,
		Type:       ItemWeapon,
		Name:       name,
		Description: "Increases attack by " + string(damage),
		Value:      damage,
		Symbol:     '/',
		Collected:  false,
	}
}

// NewArmor creates a new armor
func NewArmor(x, y int, name string, defense int) Item {
	return Item{
		X:          x,
		Y:          y,
		Type:       ItemArmor,
		Name:       name,
		Description: "Increases defense by " + string(defense),
		Value:      defense,
		Symbol:     '[',
		Collected:  false,
	}
}

// NewGold creates a new gold pile
func NewGold(x, y int, amount int) Item {
	return Item{
		X:          x,
		Y:          y,
		Type:       ItemGold,
		Name:       "Gold",
		Description: "Worth " + string(amount) + " gold",
		Value:      amount,
		Symbol:     '$',
		Collected:  false,
	}
}

// NewKey creates a new key
func NewKey(x, y int) Item {
	return Item{
		X:          x,
		Y:          y,
		Type:       ItemKey,
		Name:       "Key",
		Description: "Can unlock doors",
		Value:      1,
		Symbol:     'k',
		Collected:  false,
	}
}
