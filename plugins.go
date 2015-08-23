package main

// Interfaces
type PluginNotifier interface {
	Register(value chan []byte)
	NotifyPlugins(value []byte)
}

type PluginHandler interface {
	Update()
	GetChannel() chan []byte
}

// Implementations
type PluginNotifierService struct {
	plugins []chan []byte
}

func (pl *PluginNotifierService) Register(c chan []byte) {
	pl.plugins = append(pl.plugins, c)
}

func (pl *PluginNotifierService) NotifyPlugins(m []byte) {
	for _, c := range pl.plugins {
		c <- m
	}
}