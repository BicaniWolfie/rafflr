package main

type RaffleEntry struct {
	Patron string
	Points int
}

type RolledEntry struct {
	points int
	rolls  []int
}
