package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	gridWidth, gridHeight    = 12, 12
	initialPosX, initialPosY = gridWidth / 2, gridHeight / 2
)

var (
	meatX, meatY int
)

const (
	snakeHead = "[S]"
	snakeBody = "[+]"
	meatCell  = " M "
	emptyCell = " . "
)

const (
	Right = "right"
	Left  = "left"
	Up    = "up"
	Down  = "down"
)

var directions = map[string][2]int{
	Right: {0, 1},
	Left:  {0, -1},
	Up:    {-1, 0},
	Down:  {1, 0},
}

type Coordinates struct {
	x, y int
}

type SnakeSegment struct {
	Coordinates
	isHead bool
	markup string
}

var snake = []SnakeSegment{
	{Coordinates{initialPosX, initialPosY}, true, snakeHead},
}

func main() {
	clearScreen()

	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := keyboard.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	direction := Right
	directionChan := make(chan string)

	go handleDirectionInput(directionChan)
	generateMeat()

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case newDirection := <-directionChan:
			if isValidDirectionChange(direction, newDirection) {
				direction = newDirection
			}
		case <-ticker.C:
			moveSnake(direction)
			if snake[0].x == meatX && snake[0].y == meatY {
				growSnake()
				generateMeat()
			}
			displayGrid()
		}
	}
}

func moveSnake(direction string) {
	head := snake[0]
	newHead := SnakeSegment{
		Coordinates: moveAndWrap(head.x, head.y, directions[direction]),
		isHead:      true,
		markup:      snakeHead,
	}

	if collision(newHead) {
		log.Fatalf("Game over, your score: %d", len(snake))
	}

	snake[0].isHead = false
	snake[0].markup = snakeBody
	snake = append([]SnakeSegment{newHead}, snake...)
	snake = snake[:len(snake)-1]
}

func moveAndWrap(x, y int, dir [2]int) Coordinates {
	return Coordinates{
		x: (x + dir[0] + gridHeight) % gridHeight,
		y: (y + dir[1] + gridWidth) % gridWidth,
	}
}

func collision(head SnakeSegment) bool {
	for _, segment := range snake {
		if segment.x == head.x && segment.y == head.y {
			return true
		}
	}
	return false
}

func growSnake() {
	tail := snake[len(snake)-1]
	newTail := SnakeSegment{
		Coordinates: tail.Coordinates,
		isHead:      false,
		markup:      snakeBody,
	}
	snake = append(snake, newTail)
}

func displayGrid() {
	moveCursorToTopLeft()
	board := createEmptyBoard()

	for _, segment := range snake {
		board[segment.x][segment.y] = segment.markup
	}

	board[meatX][meatY] = meatCell
	printBoard(board)
}

func createEmptyBoard() [gridWidth][gridHeight]string {
	var board [gridWidth][gridHeight]string
	for i := 0; i < gridHeight; i++ {
		for j := 0; j < gridWidth; j++ {
			board[i][j] = emptyCell
		}
	}
	return board
}

func handleDirectionInput(directionChan chan<- string) {
	for {
		if _, key, err := keyboard.GetKey(); err == nil {
			switch key {
			case keyboard.KeyArrowRight:
				directionChan <- Right
			case keyboard.KeyArrowLeft:
				directionChan <- Left
			case keyboard.KeyArrowUp:
				directionChan <- Up
			case keyboard.KeyArrowDown:
				directionChan <- Down
			case keyboard.KeyCtrlC:
				os.Exit(0)
			}
		} else {
			log.Fatal(err)
		}
	}
}

func isValidDirectionChange(current, new string) bool {
	switch current {
	case Right:
		return new != Left
	case Left:
		return new != Right
	case Up:
		return new != Down
	case Down:
		return new != Up
	}
	return true
}

func printBoard(board [gridWidth][gridHeight]string) {
	printHorizontalBorder()

	for i := 0; i < gridHeight; i++ {
		fmt.Print("|")
		for j := 0; j < gridWidth; j++ {
			fmt.Print(board[i][j])
		}
		fmt.Print("|\n")
	}

	printHorizontalBorder()
}

func printHorizontalBorder() {
	fmt.Print("+")
	for i := 0; i < gridWidth; i++ {
		fmt.Print("---")
	}
	fmt.Println("+")
}

func clearScreen() {
	fmt.Print("\033[2J")
}

func moveCursorToTopLeft() {
	fmt.Print("\033[H")
}

func generateMeat() {
	for {
		meatX, meatY = rand.Intn(gridHeight), rand.Intn(gridWidth)
		if !isOccupied(meatX, meatY) {
			break
		}
	}
}

func isOccupied(x, y int) bool {
	for _, segment := range snake {
		if segment.x == x && segment.y == y {
			return true
		}
	}
	return false
}
