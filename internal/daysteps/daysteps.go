package daysteps

import (
	"errors"
	"fmt"
	"spentcalories"
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
	stepsInfo := strings.TrimSpace(info[0])
	steps, err := strconv.Atoi(stepsInfo)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка в данных шагов %v", err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов <= 0")
	}
	durationInfo := strings.TrimSpace(info[1])
	duration, err := time.ParseDuration(durationInfo)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка в данных времени %v", err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("продолжительность времени <= 0")
	}
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, _, err := parsePackage(data)
	if err != nil {
		return fmt.Sprintf("ошибка %v", err)
	}
	if steps <= 0 {
		return fmt.Sprint("")
	}
	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm

	switch {
	case spentcalories.TrainingInfo() == "ходьба":
		calories := spentcalories.WalkingSpentCalories()
		return fmt.Sprintf("Количество шагов: %d.\nДистанция составила: %.2f км.\nВы сожгли %.2f ккал.", steps, distanceKm, calories)

	case spentcalories.TrainingInfo() == "бег":
		calories := spentcalories.RunningSpentCalories()
		return fmt.Sprintf("Количество шагов: %d.\nДистанция составила: %.2f км.\nВы сожгли %.2f ккал.", steps, distanceKm, calories)

	}
}
