package output

type OutputBuilder interface {
	GenerateOutput(interface{}) (string, error)
}

type Config struct {
	Format string
}

type Output struct {
	Config Config
}

func (o *Output) Generate(elements interface{}) (string, error) {
	if elements == nil {
		return "", nil
	}

	var v OutputBuilder
	switch o.Config.Format {
	case "json":
		v = &JSONOutput{
			Config: o.Config,
		}
	case "shell":
		v = &ShellOutput{
			Config: o.Config,
		}
	default:
		v = &InvalidOutput{
			Config: o.Config,
		}
	}

	return v.GenerateOutput(elements)
}
