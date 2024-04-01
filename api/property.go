package api

import (
	"strings"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	corev1 "k8s.io/api/core/v1"
)

type Property interface {
	GetID() string
	GetPriority() int8
}

type TextProperty struct {
	ID        string
	Name      string
	Value     string
	Reference *corev1.ObjectReference
	Widget    func(gtk.Widgetter, *adw.NavigationView)
}

func (p *TextProperty) GetID() string {
	if p.ID == "" {
		return strings.ToLower(p.Name)
	}
	return p.ID
}

func (p *TextProperty) GetPriority() int8 {
	return 0
}

type GroupProperty struct {
	ID       string
	Priority int8
	Name     string
	Children []Property
	Widget   func(gtk.Widgetter, *adw.NavigationView)
}

func (p *GroupProperty) GetID() string {
	if p.ID == "" {
		return strings.ToLower(p.Name)
	}
	return p.ID
}

func (p *GroupProperty) GetPriority() int8 {
	return p.Priority
}
