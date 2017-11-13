package relabel

import (
	"fmt"
	"strings"

	"github.com/prometheus/common/model"
	"github.com/rajatjindal/prom-to-kairosdb/config"
)

//Process the samples and apply RelabelConfig
func Process(metric model.Metric, cfgs ...*config.RelabelConfig) model.Metric {
	for _, cfg := range cfgs {
		metric = relabel(metric, cfg)
		if metric == nil {
			return nil
		}
	}
	return metric
}

func relabel(metric model.Metric, cfg *config.RelabelConfig) model.Metric {
	values := make([]string, 0, len(cfg.SourceLabels))
	for _, labelName := range cfg.SourceLabels {
		values = append(values, string(metric[labelName]))
	}
	valueOfSourceLabels := strings.Join(values, cfg.Separator)

	switch cfg.Action {
	case config.RelabelDrop:
		if cfg.Regex.MatchString(valueOfSourceLabels) {
			return nil
		}
	case config.RelabelKeep:
		if !cfg.Regex.MatchString(valueOfSourceLabels) {
			return nil
		}
	case config.RelabelAddPrefix:
		if cfg.Regex.MatchString(valueOfSourceLabels) {
			metric[model.MetricNameLabel] = model.LabelValue(fmt.Sprintf("%s%s", cfg.Prefix, metric[model.MetricNameLabel]))
		}
	case config.RelabelLabelDrop:
		for labelName := range metric {
			if cfg.Regex.MatchString(string(labelName)) {
				delete(metric, labelName)
			}
		}
	case config.RelabelLabelKeep:
		for labelName := range metric {
			if !cfg.Regex.MatchString(string(labelName)) {
				delete(metric, labelName)
			}
		}
	default:
		fmt.Printf("warn: retrieval.relabel: unknown relabel action type %q", cfg.Action)
	}
	return metric
}