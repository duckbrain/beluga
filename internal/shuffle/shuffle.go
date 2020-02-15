package shuffle

type Interface Interface{
	Len() int
	Swap(i, j int)
}

func Shuffle(i Interface) {
	rand.Shuffle(i.Len(), i.Swap)
}