package exporter

type Exporter interface {
	Run(<-chan struct{}) error
}
