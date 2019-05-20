package samplesort

import (
	"math"
	"sort"
)

type Param interface {
	Values() []float64
}

type Simple float64

func (s Simple) Values() []float64 { return []float64{normalize(float64(s))} }

type Sample struct {
	Path      string
	Group     int
	Low_level struct {
		Average_loudness                Simple `json:"average_loudness"`
		Dynamic_complexity              Simple `json:"dynamic_complexity"`
		Barkbands_crest                 Scalar `json:"barkbands_crest"`
		Barkbands_flatness_db           Scalar `json:"barkbands_flatness_db"`
		Barkbands_kurtosis              Scalar `json:"barkbands_kurtosis"`
		Barkbands_skewness              Scalar `json:"barkbands_skewness"`
		Barkbands_spread                Scalar `json:"barkbands_spread"`
		Dissonance                      Scalar `json:"dissonance"`
		Erbbands_crest                  Scalar `json:"erbbands_crest"`
		Erbbands_flatness_db            Scalar `json:"erbbands_flatness_db"`
		Erbbands_kurtosis               Scalar `json:"erbbands_kurtosis"`
		Erbbands_skewness               Scalar `json:"erbbands_skewness"`
		Erbbands_spread                 Scalar `json:"erbbands_spread"`
		Hfc                             Scalar `json:"hfc"`
		Melbands_crest                  Scalar `json:"melbands_crest"`
		Melbands_flatness_db            Scalar `json:"melbands_flatness_db"`
		Melbands_kurtosis               Scalar `json:"melbands_kurtosis"`
		Melbands_skewness               Scalar `json:"melbands_skewness"`
		Melbands_spread                 Scalar `json:"melbands_spread"`
		Pitch_salience                  Scalar `json:"pitch_salience"`
		Silence_rate_20dB               Scalar `json:"silence_rate_20db"`
		Silence_rate_30dB               Scalar `json:"silence_rate_30db"`
		Silence_rate_60dB               Scalar `json:"silence_rate_60db"`
		Spectral_centroid               Scalar `json:"spectral_centroid"`
		Spectral_complexity             Scalar `json:"spectral_complexity"`
		Spectral_decrease               Scalar `json:"spectral_decrease"`
		Spectral_energy                 Scalar `json:"spectral_energy"`
		Spectral_energyband_high        Scalar `json:"spectral_energyband_high"`
		Spectral_energyband_low         Scalar `json:"spectral_energyband_low"`
		Spectral_energyband_middle_high Scalar `json:"spectral_energyband_middle_high"`
		Spectral_energyband_middle_low  Scalar `json:"spectral_energyband_middle_low"`
		Spectral_entropy                Scalar `json:"spectral_entropy"`
		Spectral_flux                   Scalar `json:"spectral_flux"`
		Spectral_kurtosis               Scalar `json:"spectral_kurtosis"`
		Spectral_rms                    Scalar `json:"spectral_rms"`
		Spectral_rolloff                Scalar `json:"spectral_rolloff"`
		Spectral_skewness               Scalar `json:"spectral_skewness"`
		Spectral_spread                 Scalar `json:"spectral_spread"`
		Spectral_strongpeak             Scalar `json:"spectral_strongpeak"`
		Zerocrossingrate                Scalar `json:"zerocrossingrate"`
	} `json:"lowlevel"`
}

func (s Sample) Values() []float64 {
	result := make([]float64, 0, 2+9*38)
	result = append(result, s.Low_level.Average_loudness.Values()...)
	result = append(result, s.Low_level.Dynamic_complexity.Values()...)
	result = append(result, s.Low_level.Barkbands_crest.Values()...)
	result = append(result, s.Low_level.Barkbands_flatness_db.Values()...)
	result = append(result, s.Low_level.Barkbands_kurtosis.Values()...)
	result = append(result, s.Low_level.Barkbands_skewness.Values()...)
	result = append(result, s.Low_level.Barkbands_spread.Values()...)
	result = append(result, s.Low_level.Dissonance.Values()...)
	result = append(result, s.Low_level.Erbbands_crest.Values()...)
	result = append(result, s.Low_level.Erbbands_flatness_db.Values()...)
	result = append(result, s.Low_level.Erbbands_kurtosis.Values()...)
	result = append(result, s.Low_level.Erbbands_skewness.Values()...)
	result = append(result, s.Low_level.Erbbands_spread.Values()...)
	result = append(result, s.Low_level.Hfc.Values()...)
	result = append(result, s.Low_level.Melbands_crest.Values()...)
	result = append(result, s.Low_level.Melbands_flatness_db.Values()...)
	result = append(result, s.Low_level.Melbands_kurtosis.Values()...)
	result = append(result, s.Low_level.Melbands_skewness.Values()...)
	result = append(result, s.Low_level.Melbands_spread.Values()...)
	result = append(result, s.Low_level.Pitch_salience.Values()...)
	result = append(result, s.Low_level.Silence_rate_20dB.Values()...)
	result = append(result, s.Low_level.Silence_rate_30dB.Values()...)
	result = append(result, s.Low_level.Silence_rate_60dB.Values()...)
	result = append(result, s.Low_level.Spectral_centroid.Values()...)
	result = append(result, s.Low_level.Spectral_complexity.Values()...)
	result = append(result, s.Low_level.Spectral_decrease.Values()...)
	result = append(result, s.Low_level.Spectral_energy.Values()...)
	result = append(result, s.Low_level.Spectral_energyband_high.Values()...)
	result = append(result, s.Low_level.Spectral_energyband_low.Values()...)
	result = append(result, s.Low_level.Spectral_energyband_middle_high.Values()...)
	result = append(result, s.Low_level.Spectral_energyband_middle_low.Values()...)
	result = append(result, s.Low_level.Spectral_entropy.Values()...)
	result = append(result, s.Low_level.Spectral_flux.Values()...)
	result = append(result, s.Low_level.Spectral_kurtosis.Values()...)
	result = append(result, s.Low_level.Spectral_rms.Values()...)
	result = append(result, s.Low_level.Spectral_rolloff.Values()...)
	result = append(result, s.Low_level.Spectral_skewness.Values()...)
	result = append(result, s.Low_level.Spectral_spread.Values()...)
	result = append(result, s.Low_level.Spectral_strongpeak.Values()...)
	result = append(result, s.Low_level.Zerocrossingrate.Values()...)
	return result
}

type Scalar struct {
	Dmean  float64 `json:"dmean"`
	Dmean2 float64 `json:"dmean2"`
	Dvar   float64 `json:"dvar"`
	Dvar2  float64 `json:"dvar2"`
	Max    float64 `json:"max"`
	Mean   float64 `json:"mean"`
	Median float64 `json:"median"`
	Min    float64 `json:"min"`
	Var    float64 `json:"var"`
}

func (s Scalar) Values() []float64 {
	return []float64{
		normalize(s.Dmean),
		normalize(s.Dmean2),
		normalize(s.Dvar),
		normalize(s.Dvar2),
		normalize(s.Max),
		normalize(s.Mean),
		normalize(s.Median),
		normalize(s.Min),
		normalize(s.Var),
	}
}

type Collection struct {
	Samples []*Sample
}

func (c *Collection) build(samples <-chan *Sample) {
	c.Samples = make([]*Sample, 0)
	for sample := range samples {
		c.Samples = append(c.Samples, sample)
	}
}

func (c *Collection) features() [][]float64 {
	result := make([][]float64, len(c.Samples), len(c.Samples))
	for i, sample := range c.Samples {
		result[i] = sample.Values()
	}
	return result
}

func (c *Collection) categorize(groups []int) {
	if len(c.Samples) != len(groups) {
		panic("collection and analysis are different")
	}
	for i, group := range groups {
		c.Samples[i].Group = group
	}
	sort.Slice(c.Samples, func(i, j int) bool {
		return c.Samples[i].Group < c.Samples[j].Group
	})
}

func normalize(input float64) (output float64) {
	return math.Round(input/PARAM_NORMALIZE) * PARAM_NORMALIZE
}
