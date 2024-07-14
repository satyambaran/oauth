package structs

type Env struct {
    DbUrl                  string
    AccessTokenExpiryHour  int
    RefreshTokenExpiryHour int
    AccessTokenSecret      string
    RefreshTokenSecret     string
}
