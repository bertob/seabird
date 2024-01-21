package ui

import (
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type PrefsWindow struct {
	*adw.PreferencesWindow
	navigationView *adw.NavigationView
}

func NewPreferencesWindow() *PrefsWindow {
	p := PrefsWindow{PreferencesWindow: adw.NewPreferencesWindow()}

	content := gtk.NewBox(gtk.OrientationVertical, 0)
	p.navigationView = adw.NewNavigationView()
	p.navigationView.Add(adw.NewNavigationPage(content, "main"))
	p.SetContent(p.navigationView)

	header := adw.NewHeaderBar()
	view := adw.NewViewSwitcher()
	view.SetPolicy(adw.ViewSwitcherPolicyWide)
	header.SetTitleWidget(view)
	content.Append(header)

	stack := adw.NewViewStack()
	stack.AddTitledWithIcon(p.createGeneralPage(), "general", "General", "document-properties-symbolic")
	content.Append(stack)
	view.SetStack(stack)

	p.ConnectUnrealize(func() {
		if err := application.prefs.Save(); err != nil {
			ShowErrorDialog(&p.Window.Window, "Could not save preferences", err)
			return
		}
	})

	return &p
}

func (p *PrefsWindow) createGeneralPage() gtk.Widgetter {
	page := adw.NewPreferencesPage()

	general := adw.NewPreferencesGroup()
	theme := adw.NewComboRow()
	theme.SetTitle("Theme")
	themes := gtk.NewStringList([]string{"Dark", "Light"})
	theme.SetModel(themes)
	general.Add(theme)

	clusters := adw.NewPreferencesGroup()
	clusters.SetTitle("Clusters")
	addCluster := gtk.NewButton()
	addCluster.AddCSSClass("flat")
	addCluster.SetIconName("list-add")
	addCluster.ConnectClicked(func() {
		p.navigationView.Push(NewClusterPrefPage(&p.Window.Window, nil).NavigationPage)
	})

	clusters.SetHeaderSuffix(addCluster)
	for _, c := range application.prefs.Clusters {
		cluster := c
		row := adw.NewActionRow()
		row.SetActivatable(true)
		row.ConnectActivated(func() {
			p.navigationView.Push(NewClusterPrefPage(&p.Window.Window, cluster).NavigationPage)
		})
		row.SetTitle(cluster.Name)
		row.AddSuffix(gtk.NewImageFromIconName("go-next-symbolic"))
		clusters.Add(row)
	}

	page.Add(general)
	page.Add(clusters)

	return page
}

func (p *PrefsWindow) other() gtk.Widgetter {
	return gtk.NewLabel("other")
}
