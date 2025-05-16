package main

import (
	"fmt"
	"math/rand"
)

// Player represents the player character in the game
type Player struct {
	X, Y      int     // Position coordinates
	Health    int     // Current health points
	MaxHealth int     // Maximum health points
	Attack    int     // Attack damage
	Defense   int     // Damage reduction
	Gold      int     // Gold collected
	Level     int     // Player level
	Exp       int     // Experience points
	Inventory []Item  // Items carried by the player
}

// NewPlayer creates a new player at the specified position
func NewPlayer(x, y int) *Player {
	return &Player{
		X:         x,
		Y:         y,
		Health:    20,
		MaxHealth: 20,
		Attack:    3,
		Defense:   1,
		Gold:      0,
		Level:     1,
		Exp:       0,
		Inventory: make([]Item, 0),
	}
}

// Move attempts to move the player in the specified direction
func (p *Player) Move(dx, dy int, d *Dungeon) {
	newX := p.X + dx
	newY := p.Y + dy

	// Check if there's an enemy at the target position
	if enemy := d.GetEnemyAt(newX, newY); enemy != nil {
		// Attack the enemy instead of moving
		p.AttackEnemy(enemy, d)
		return
	}

	// Check if the position is walkable
	if d.IsWalkable(newX, newY) {
		p.X = newX
		p.Y = newY
		
		// Check for items or special tiles at the new position
		p.CheckPosition(d)
	} else {
		fmt.Println("You can't move there!")
	}
}

// AttackEnemy handles combat with an enemy
func (p *Player) AttackEnemy(enemy *Enemy, d *Dungeon) {
	// Calculate damage dealt to enemy
	damage := p.Attack
	
	// Apply damage to enemy
	enemy.Health -= damage
	
	fmt.Printf("You attack the %s for %d damage!\n", enemy.Name, damage)
	
	// Check if enemy is defeated
	if enemy.Health <= 0 {
		fmt.Printf("You defeated the %s!\n", enemy.Name)
		
		// Award experience and possibly gold
		expGain := 5 + enemy.Damage * 2
		p.Exp += expGain
		fmt.Printf("You gained %d experience points.\n", expGain)
		
		// Check for level up
		p.CheckLevelUp()
		
		// Remove the enemy from the dungeon
		d.RemoveEnemy(enemy)
		
		// 50% chance to drop gold
		if rand.Intn(2) == 0 {
			goldAmount := 1 + rand.Intn(10)
			p.Gold += goldAmount
			fmt.Printf("You found %d gold!\n", goldAmount)
		}
	} else {
		// Enemy counterattack
		enemyDamage := enemy.Damage - p.Defense
		if enemyDamage < 1 {
			enemyDamage = 1 // Minimum damage is 1
		}
		
		p.Health -= enemyDamage
		fmt.Printf("The %s attacks you for %d damage!\n", enemy.Name, enemyDamage)
		
		// Check if player is defeated
		if p.Health <= 0 {
			fmt.Println("You have been defeated! Game over.")
		}
	}
}

// CheckPosition checks for items or special tiles at the player's position
func (p *Player) CheckPosition(d *Dungeon) {
	// Get the tile at the player's position
	tile := d.GetTileAt(p.X, p.Y)
	
	switch tile {
	case Treasure:
		// Collect treasure
		p.Gold += 10 + rand.Intn(20)
		fmt.Printf("You found some gold! You now have %d gold.\n", p.Gold)
		d.Grid[p.Y][p.X] = rune(Floor) // Replace with floor
		
	case Trap:
		// Trigger trap
		damage := 2 + rand.Intn(3)
		p.Health -= damage
		fmt.Printf("You triggered a trap! You take %d damage.\n", damage)
		d.Grid[p.Y][p.X] = rune(Floor) // Trap is now disarmed
		
		// Check if player died from trap
		if p.Health <= 0 {
			fmt.Println("You died from a trap! Game over.")
		}
		
	case Door:
		// Open door
		fmt.Println("You open the door.")
		d.Grid[p.Y][p.X] = rune(Floor) // Door is now open
		
	case StairsDown:
		// Go to next level
		fmt.Println("You found stairs leading down! Press '>' to descend to the next level.")
	}
	
	// Check for items
	if item := d.GetItemAt(p.X, p.Y); item != nil {
		p.CollectItem(item)
	}
}

// CollectItem adds an item to the player's inventory
func (p *Player) CollectItem(item *Item) {
	// Mark the item as collected
	item.Collected = true
	
	// Handle different item types
	switch item.Type {
	case ItemGold:
		p.Gold += item.Value
		fmt.Printf("You collected %d gold! You now have %d gold.\n", item.Value, p.Gold)
		
	case ItemPotion:
		// Add to inventory
		p.Inventory = append(p.Inventory, *item)
		fmt.Printf("You picked up a %s.\n", item.Name)
		
	case ItemWeapon:
		// Add to inventory
		p.Inventory = append(p.Inventory, *item)
		fmt.Printf("You picked up a %s.\n", item.Name)
		
	case ItemArmor:
		// Add to inventory
		p.Inventory = append(p.Inventory, *item)
		fmt.Printf("You picked up a %s.\n", item.Name)
	}
}

// UseItem uses an item from the inventory
func (p *Player) UseItem(itemIndex int) {
	// Check if the index is valid
	if itemIndex < 0 || itemIndex >= len(p.Inventory) {
		fmt.Println("Invalid item index.")
		return
	}
	
	// Get the item
	item := p.Inventory[itemIndex]
	
	// Handle different item types
	switch item.Type {
	case ItemPotion:
		// Heal the player
		healAmount := item.Value
		p.Health += healAmount
		if p.Health > p.MaxHealth {
			p.Health = p.MaxHealth
		}
		fmt.Printf("You drink the %s and heal for %d health points.\n", item.Name, healAmount)
		
		// Remove the item from inventory
		p.Inventory = append(p.Inventory[:itemIndex], p.Inventory[itemIndex+1:]...)
		
	case ItemWeapon:
		// Equip the weapon
		p.Attack = item.Value
		fmt.Printf("You equip the %s. Your attack is now %d.\n", item.Name, p.Attack)
		
	case ItemArmor:
		// Equip the armor
		p.Defense = item.Value
		fmt.Printf("You equip the %s. Your defense is now %d.\n", item.Name, p.Defense)
	}
}

// CheckLevelUp checks if the player has enough experience to level up
func (p *Player) CheckLevelUp() {
	// Simple level up formula: 100 * current level
	expNeeded := 100 * p.Level
	
	if p.Exp >= expNeeded {
		p.Level++
		p.Exp -= expNeeded
		p.MaxHealth += 5
		p.Health = p.MaxHealth
		p.Attack++
		
		fmt.Printf("Level up! You are now level %d.\n", p.Level)
		fmt.Printf("Your health increased to %d and your attack increased to %d.\n", p.MaxHealth, p.Attack)
		
		// Check if there's another level up available
		p.CheckLevelUp()
	}
}

// DisplayStatus shows the player's current stats
func (p *Player) DisplayStatus() {
	fmt.Printf("Health: %d/%d | Attack: %d | Defense: %d | Gold: %d | Level: %d | Exp: %d/%d\n",
		p.Health, p.MaxHealth, p.Attack, p.Defense, p.Gold, p.Level, p.Exp, 100*p.Level)
}

// DisplayInventory shows the player's inventory
func (p *Player) DisplayInventory() {
	if len(p.Inventory) == 0 {
		fmt.Println("Your inventory is empty.")
		return
	}
	
	fmt.Println("Inventory:")
	for i, item := range p.Inventory {
		fmt.Printf("%d. %s (%s)\n", i+1, item.Name, item.Description)
	}
}
