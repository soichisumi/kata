package _go

import (
	"contrib.go.opencensus.io/integrations/ocsql"
	"go.opencensus.io/trace"
	"os"
)

func WrapDBDriverWithOCSQL(driverName string, sampleProbability float64) (string, error){
	var o ocsql.TraceOption
	switch os.Getenv("env") {
	case "bench":
		_o, err := benchConfig()
		if err != nil {
			return "", err
		}
		o = _o
	default:
		_o, err := defaultConfig(sampleProbability)
		if err != nil {
			return "", err
		}
		o = _o
	}
	return ocsql.Register(driverName, o)
}

func defaultConfig(sampleProbability float64) (ocsql.TraceOption, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	return ocsql.WithOptions(ocsql.TraceOptions{
		Sampler:      trace.ProbabilitySampler(sampleProbability),
		DefaultAttributes: []trace.Attribute{
			trace.StringAttribute("hostname", hostname),
			trace.StringAttribute("env", os.Getenv("env")),
		},
		Query:        true,
		QueryParams:  true,
	}), nil
}

func benchConfig() (ocsql.TraceOption, error) {
	return ocsql.WithOptions(ocsql.TraceOptions{
		Sampler:      trace.NeverSample(),
	}), nil
}