package shuffle

import (
	"math/rand"
	"time"

	"github.com/rast1025/sekretsanta_bot/internal/models"
)

func Randomize(users []models.User) map[models.User]models.User {
	result := make(map[models.User]models.User, len(users))
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := append([]models.User{}, users...)
	for {
		valid := true
		r.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})
		for i := 0; i < len(shuffled); i++ {
			if shuffled[i].Username == users[i].Username {
				valid = false
				break
			}
		}
		if valid {
			for i := 0; i < len(shuffled); i++ {
				result[shuffled[i]] = users[i]
			}
			return result
		}
	}
}
