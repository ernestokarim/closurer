package config

import (
	"github.com/ernestokarim/closurer/app"
)

type Config struct {
	Build string `xml:"build,attr"`

	Ignores []*IgnoreNode `xml:"ignore"`
	Map     *MapNode      `xml:"map"`
	Js      *JsNode       `xml:"js"`
	Gss     *GssNode      `xml:"gss"`
	Soy     *SoyNode      `xml:"soy"`
	Library *LibraryNode  `xml:"library"`
}

// ==================================================================

type JsNode struct {
	Root        string `xml:"root,attr"`
	Compiler    string `xml:"compiler,attr"`
	Language    string `xml:"language,attr"`
	Formatting  string `xml:"formatting,attr"`
	SideEffects string `xml:"side-effects,attr"`

	Checks   *ChecksNode     `xml:"checks"`
	Targets  []*JsTargetNode `xml:"target"`
	Inputs   []*InputNode    `xml:"input"`
	Externs  []*ExternNode   `xml:"extern"`
	Prepends []*PrependNode  `xml:"prepend"`
}

func (n *JsNode) CurTarget() *JsTargetNode {
	if n == nil {
		return nil
	}

	for _, t := range n.Targets {
		if t.Name == SelectedTarget {
			return t
		}
	}
	return nil
}

// ==================================================================

type IgnoreNode struct {
	Path string `xml:"path,attr"`
}

// ==================================================================

type ChecksNode struct {
	Errors   []*CheckNode `xml:"error"`
	Warnings []*CheckNode `xml:"warning"`
	Offs     []*CheckNode `xml:"off"`
}

// ==================================================================

type CheckNode struct {
	Name string `xml:"name,attr"`
}

// ==================================================================

type JsTargetNode struct {
	Name     string `xml:"name,attr"`
	Mode     string `xml:"mode,attr"`
	Level    string `xml:"level,attr"`
	Output   string `xml:"output,attr"`
	Inherits string `xml:"inherits,attr"`

	Defines []*DefineNode `xml:"define"`
}

func (t *JsTargetNode) ApplyInherits() error {
	if t.Name == "" {
		return app.Errorf("The name of the target is required")
	}

	if t.Inherits == "" {
		return nil
	}

	for _, parent := range globalConf.Js.Targets {
		if parent.Name == t.Name {
			return app.Errorf("Inherits should reference a previous target: %s", t.Name)
		}
		if parent.Name != t.Inherits {
			continue
		}

		if t.Mode == "" {
			t.Mode = parent.Mode
		}
		if t.Level == "" {
			t.Level = parent.Level
		}
		if t.Output == "" {
			t.Output = parent.Output
		}

		for _, d := range parent.Defines {
			if !t.HasDefine(d.Name) {
				t.Defines = append(t.Defines, d.Clone())
			}
		}

		return nil
	}

	panic("not reached")
}

func (t *JsTargetNode) HasDefine(name string) bool {
	for _, d := range t.Defines {
		if d.Name == name {
			return true
		}
	}
	return false
}

// ==================================================================

type DefineNode struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

func (d *DefineNode) Clone() *DefineNode {
	return &DefineNode{
		Name:  d.Name,
		Value: d.Value,
	}
}

// ==================================================================

type MapNode struct {
	File string `xml:"file,attr"`
}

// ==================================================================

type InputNode struct {
	File string `xml:"file,attr"`
}

// ==================================================================

type ExternNode struct {
	File string `xml:"file,attr"`
}

// ==================================================================

type GssNode struct {
	Compiler string `xml:"compiler,attr"`

	Targets []*GssTargetNode `xml:"target"`
	Funcs   []*FuncNode      `xml:"func"`
	Inputs  []*InputNode     `xml:"input"`
}

func (n *GssNode) CurTarget() *GssTargetNode {
	if n == nil {
		return nil
	}

	for _, t := range n.Targets {
		if t.Name == SelectedTarget {
			return t
		}
	}
	return nil
}

// ==================================================================

type GssTargetNode struct {
	Name     string `xml:"name,attr"`
	Rename   string `xml:"rename,attr"`
	Output   string `xml:"output,attr"`
	Inherits string `xml:"inherits,attr"`

	Defines []*DefineNode `xml:"define"`
}

func (t *GssTargetNode) ApplyInherits() error {
	if t.Name == "" {
		return app.Errorf("The name of the target is required")
	}

	if t.Inherits == "" {
		return nil
	}

	for _, parent := range globalConf.Gss.Targets {
		if parent.Name == t.Name {
			return app.Errorf("Inherits should reference a previous target: %s", t.Name)
		}
		if parent.Name != t.Inherits {
			continue
		}

		if t.Rename == "" {
			t.Rename = parent.Rename
		}
		if t.Output == "" {
			t.Output = parent.Output
		}

		for _, d := range parent.Defines {
			if !t.HasDefine(d.Name) {
				t.Defines = append(t.Defines, d.Clone())
			}
		}

		return nil
	}

	panic("not reached")
}

func (t *GssTargetNode) HasDefine(name string) bool {
	for _, d := range t.Defines {
		if d.Name == name {
			return true
		}
	}
	return false
}

// ==================================================================

type SoyNode struct {
	Root     string `xml:"root,attr"`
	Compiler string `xml:"compiler,attr"`
}

// ==================================================================

type LibraryNode struct {
	Root string `xml:"root,attr"`
}

// ==================================================================

type FuncNode struct {
	Name string `xml:"name,attr"`
}

// ==================================================================

type PrependNode struct {
	File string `xml:"file,attr"`
}
