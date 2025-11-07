package daysteps

import (
	"errors"
	"fmt"
	"four_sprint_exam/internal/spentcalories"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	info := strings.Split(data, ",")
	if len(info) != 2 {
		return 0, 0, errors.New("неверный формат данных. Ожидается: шаги,время")
	}
	stepsInfo := info[0]
	steps, err := strconv.Atoi(stepsInfo)
	if err != nil {
		return 0, 0, fmt.Errorf("%v", err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("%v", err)
	}
	durationInfo := info[1]
	duration, err := time.ParseDuration(durationInfo)
	if err != nil {
		return 0, 0, fmt.Errorf("%v", err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("%v", err)
	}
	return steps, duration, err
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	if steps <= 0 {
		log.Println(err)
		return ""
	}
	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return fmt.Sprintf("%v", err)
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm, calories)

}
