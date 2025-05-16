package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Game states
const (
	StateMainMenu = iota
	StatePlaying
	StateInventory
	StateGameOver
)

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
	
	// Initialize game state
	gameState := StatePlaying
	
	// Create a new dungeon
	dungeon := NewDungeon(80, 24)
	
	// Create a new player in the first room
	var player *Player
	if len(dungeon.Rooms) > 0 {
		// Place player in the center of the first room
		room := dungeon.Rooms[0]
		player = NewPlayer(room.X+room.Width/2, room.Y+room.Height/2)
	} else {
		// Fallback if no rooms were generated
		player = NewPlayer(1, 1)
	}

	// Create a reader for user input
	reader := bufio.NewReader(os.Stdin)

	// Display welcome message and instructions
	fmt.Println("=== Welcome to Dungeon Crawler ===")
	printHelp()

	// Main game loop
	for {
		// Handle different game states
		switch gameState {
		case StatePlaying:
			// Display the dungeon and player status
			dungeon.Print(player)
			player.DisplayStatus()
			
			// Process player input
			fmt.Print("\nEnter command: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			
			// Process the command
			switch input {
			case "q", "quit":
				fmt.Println("Thanks for playing! Goodbye!")
				return
				
			case "w", "up":
				player.Move(0, -1, dungeon)
				dungeon.MoveEnemies(player) // Enemies move after player
				
			case "s", "down":
				player.Move(0, 1, dungeon)
				dungeon.MoveEnemies(player)
				
			case "a", "left":
				player.Move(-1, 0, dungeon)
				dungeon.MoveEnemies(player)
				
			case "d", "right":
				player.Move(1, 0, dungeon)
				dungeon.MoveEnemies(player)
				
			case "i", "inventory":
				gameState = StateInventory
				
			case ">":
				// Check if player is on stairs
				if dungeon.GetTileAt(player.X, player.Y) == StairsDown {
					// Generate a new dungeon level
					dungeon = NewDungeon(80, 24)
					dungeon.Level = dungeon.Level + 1
					
					// Place player in the first room of the new level
					if len(dungeon.Rooms) > 0 {
						room := dungeon.Rooms[0]
						player.X = room.X + room.Width/2
						player.Y = room.Y + room.Height/2
					} else {
						player.X, player.Y = 1, 1
					}
					
					fmt.Printf("You descend to dungeon level %d...\n", dungeon.Level)
				} else {
					fmt.Println("There are no stairs here.")
				}
				
			case "h", "help":
				printHelp()
				
			case "r", "rest":
				// Rest to recover health (with risk)
				if rand.Intn(3) == 0 {
					// 1/3 chance of enemy encounter during rest
					fmt.Println("Your rest is interrupted by a wandering monster!")
					// Spawn a random enemy near the player
					spawnEnemyNearPlayer(player, dungeon)
				} else {
					// Recover some health
					healAmount := 2 + rand.Intn(3)
					player.Health += healAmount
					if player.Health > player.MaxHealth {
						player.Health = player.MaxHealth
					}
					fmt.Printf("You rest and recover %d health points.\n", healAmount)
					dungeon.MoveEnemies(player) // Enemies still move while resting
				}
				
			default:
				fmt.Println("Unknown command. Type 'h' or 'help' for instructions.")
			}
			
			// Check if player is dead
			if player.Health <= 0 {
				gameState = StateGameOver
			}
			
		case StateInventory:
			// Display inventory
			fmt.Println("\n=== Inventory ===")
			player.DisplayInventory()
			fmt.Println("\nEnter item number to use it, or 'b' to go back:")
			
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			
			if input == "b" || input == "back" {
				gameState = StatePlaying
			} else {
				// Try to parse item index
				var itemIndex int
				_, err := fmt.Sscanf(input, "%d", &itemIndex)
				if err == nil && itemIndex > 0 && itemIndex <= len(player.Inventory) {
					player.UseItem(itemIndex - 1) // Convert to 0-based index
				} else {
					fmt.Println("Invalid item selection.")
				}
			}
			
		case StateGameOver:
			// Game over screen
			fmt.Println("\n=== GAME OVER ===")
			fmt.Printf("You died on dungeon level %d.\n", dungeon.Level)
			fmt.Printf("Final score: %d gold collected.\n", player.Gold)
			fmt.Println("\nPress 'r' to restart or 'q' to quit:")
			
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			
			if input == "r" || input == "restart" {
				// Restart the game
				dungeon = NewDungeon(80, 24)
				if len(dungeon.Rooms) > 0 {
					room := dungeon.Rooms[0]
					player = NewPlayer(room.X+room.Width/2, room.Y+room.Height/2)
				} else {
					player = NewPlayer(1, 1)
				}
				gameState = StatePlaying
			} else if input == "q" || input == "quit" {
				fmt.Println("Thanks for playing! Goodbye!")
				return
			}
		}
	}
}

// printHelp displays the game instructions
func printHelp() {
	fmt.Println("\n=== Instructions ===")
	fmt.Println("Movement: w/up, a/left, s/down, d/right")
	fmt.Println("Actions:")
	fmt.Println("  i - Open inventory")
	fmt.Println("  > - Descend stairs (when standing on them)")
	fmt.Println("  r - Rest to recover health")
	fmt.Println("  h - Show this help")
	fmt.Println("  q - Quit game")
	fmt.Println("\nSymbols:")
	fmt.Println("  @ - Player")
	fmt.Println("  . - Floor")
	fmt.Println("  # - Wall")
	fmt.Println("  + - Door")
	fmt.Println("  $ - Treasure")
	fmt.Println("  ^ - Trap")
	fmt.Println("  > - Stairs down")
	fmt.Println("  g/o/T/s - Enemies (goblin, orc, troll, skeleton)")
	fmt.Println("\nCombat: Move into enemies to attack them")
	fmt.Println()
}

// spawnEnemyNearPlayer creates a random enemy near the player
func spawnEnemyNearPlayer(player *Player, dungeon *Dungeon) {
	// Define possible spawn positions (adjacent to player)
	positions := []struct{ dx, dy int }{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0},           {1, 0},
		{-1, 1},  {0, 1},  {1, 1},
	}
	
	// Try each position
	for _, pos := range positions {
		x, y := player.X+pos.dx, player.Y+pos.dy
		
		// Check if position is valid
		if dungeon.IsWalkable(x, y) && dungeon.GetEnemyAt(x, y) == nil {
			// Create a random enemy
			enemyTypes := []struct {
				name   string
				symbol rune
				health int
				damage int
			}{
				{"Goblin", 'g', 3, 1},
				{"Rat", 'r', 1, 1},
			}
			
			enemyType := enemyTypes[rand.Intn(len(enemyTypes))]
			
			// Create and add the enemy
			enemy := &Enemy{
				X:       x,
				Y:       y,
				Health:  enemyType.health,
				Symbol:  enemyType.symbol,
				Name:    enemyType.name,
				Damage:  enemyType.damage,
				Hostile: true,
			}
			
			dungeon.Enemies = append(dungeon.Enemies, enemy)
			fmt.Printf("A %s appears!\n", enemy.Name)
			return
		}
	}
}
