package config

var (
	ContextKeyAuth = contextKey("auth")
)

type contextKey string

//func (c contextKey) toString() string {
//	return string(c)
//}
//
//func (c contextKey) toInt() (int, error) {
//	v, err := strconv.Atoi(string(c))
//	if err != nil {
//		return 0, fmt.Errorf("failed to get an integer value: %w", err)
//	}
//
//	return v, nil
//}
