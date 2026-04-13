package handler

import "strconv"

func parseChatID(chatIDStr string) (uint, error) {
	chatID, err := strconv.ParseUint(chatIDStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(chatID), nil
}

func parseLimit(limitStr string) int {
	if limitStr == "" {
		return limitDefault
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return limitDefault
	}

	if limit <= 0 {
		return limitDefault
	}

	if limit > limitMax {
		return limitMax
	}

	return limit
}
