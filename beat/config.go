package beat

type FileSizeConfig struct {
    Period *int64
    Paths *[]string
}

type ConfigSettings struct {
    Input FileSizeConfig
}
