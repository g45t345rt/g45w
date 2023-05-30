package animation

import (
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

func ButtonClick() *Animation {
	sequence := gween.NewSequence(
		gween.New(1, .98, .1, ease.Linear),
		gween.New(.98, 1, .4, ease.OutBounce),
	)
	return NewAnimation(false, sequence)
}

func ButtonEnter() *Animation {
	sequence := gween.NewSequence(
		gween.New(1, .99, .1, ease.Linear),
	)
	return NewAnimation(false, sequence)
}

func ButtonLeave() *Animation {
	sequence := gween.NewSequence(
		gween.New(.99, 1, .1, ease.Linear),
	)
	return NewAnimation(false, sequence)
}
