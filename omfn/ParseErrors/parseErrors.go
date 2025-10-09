package perr

type TokenError struct {
	Begin int
	End   int
	Msg   string
}

func (err TokenError) Error() string {
	return err.Msg
}

func UnexpectedToken(tok string, exTok string) string {
	return "예상치 못 한 토큰: '" + tok + "'. '" + exTok + "'토큰을 예상함"
}
