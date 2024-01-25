package behavior

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/imkira/go-observer/v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

type basePreferences struct {
	ColorScheme adw.ColorScheme
	Clusters    []ClusterPreferences
}

type Preferences struct {
	*basePreferences
	Clusters []observer.Property[ClusterPreferences]
}

func PrefPath() string {
	cd, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return path.Join(cd, "seabird", "prefs.json")
}

type ClusterPreferences struct {
	Name       string
	Host       string
	TLS        rest.TLSClientConfig
	Navigation struct {
		Favourites []schema.GroupVersionResource
	}
}

func LoadPreferences() (*Preferences, error) {
	if _, err := os.Stat(PrefPath()); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			prefs := Preferences{basePreferences: &basePreferences{}}
			prefs.Defaults()
			return &prefs, nil
		}
		return nil, err
	}

	f, err := os.Open(PrefPath())
	if err != nil {
		return nil, err
	}

	var base basePreferences
	if err := json.NewDecoder(f).Decode(&base); err != nil {
		return nil, err
	}
	base.Defaults()

	prefs := Preferences{
		basePreferences: &base,
	}
	for _, cluster := range base.Clusters {
		prefs.Clusters = append(prefs.Clusters, observer.NewProperty(cluster))
	}

	return &prefs, nil
}

func (c *basePreferences) Defaults() {
	for i, _ := range c.Clusters {
		c.Clusters[i].Defaults()
	}
}

func (c *ClusterPreferences) Defaults() {
	if len(c.Navigation.Favourites) == 0 {
		c.Navigation.Favourites = []schema.GroupVersionResource{
			{
				Version:  corev1.SchemeGroupVersion.Version,
				Resource: "pods",
			},
			{
				Version:  corev1.SchemeGroupVersion.Version,
				Resource: "configmaps",
			},
			{
				Version:  corev1.SchemeGroupVersion.Version,
				Resource: "secrets",
			},
			{
				Version:  corev1.SchemeGroupVersion.Version,
				Resource: "namespaces",
			},
			{
				Version:  corev1.SchemeGroupVersion.Version,
				Resource: "services",
			},
			{
				Group:    appsv1.SchemeGroupVersion.Group,
				Version:  appsv1.SchemeGroupVersion.Version,
				Resource: "deployments",
			},
			{
				Group:    appsv1.SchemeGroupVersion.Group,
				Version:  appsv1.SchemeGroupVersion.Version,
				Resource: "statefulsets",
			},
		}
	}
}

func (c *Preferences) Save() error {
	if err := os.MkdirAll(path.Dir(PrefPath()), os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(PrefPath())
	if err != nil {
		return err
	}

	c.basePreferences.Clusters = []ClusterPreferences{}
	for _, v := range c.Clusters {
		c.basePreferences.Clusters = append(c.basePreferences.Clusters, v.Value())
	}

	if err := json.NewEncoder(f).Encode(c.basePreferences); err != nil {
		return err
	}

	return nil
}
