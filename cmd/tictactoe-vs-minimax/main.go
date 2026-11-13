// dev in
// https://github.com/PulsarG

package main

import (
	"tictactoe-vs-minimax/internal/di"

	"go.uber.org/fx"
)

func main() {
	fx.New(di.Module).Run()
}
