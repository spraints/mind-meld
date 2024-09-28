package fetch

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spraints/mind-meld/appcmd"
	"github.com/spraints/mind-meld/lmsp"
)

type Target interface {
	Open() (TargetInstance, error)
}

type TargetInstance interface {
	Add(name string, data []byte) error
	Finish() (string, error)
}

func Run(app appcmd.App, target Target) error {
	t, err := target.Open()
	if err != nil {
		return err
	}

	projects, err := listProjects(app)

	for _, project := range projects {
		data, err := readPythonProject(project)
		if err != nil {
			return fmt.Errorf("%s: %w", project.RelPath, err)
		}

		if data != nil {
			if err := t.Add(pyName(project), data); err != nil {
				return fmt.Errorf("%s: %w", project.RelPath, err)
			}
		}
	}

	msg, err := t.Finish()
	if err != nil {
		return fmt.Errorf("error finishing fetch: %w", err)
	}

	fmt.Printf("%s.\n", msg)

	return nil
}

func pyName(p project) string {
	ext := filepath.Ext(p.RelPath)
	bareRelPath := p.RelPath[:len(p.RelPath)-len(ext)]
	return bareRelPath + ".py"
}

type project struct {
	// RelPath is the dirs + filename, relative to the root of mindstorms's storage dir.
	RelPath string
	// Path is the original path to the file.
	Path string
}

func listProjects(app appcmd.App) ([]project, error) {
	for _, d := range app.ProjectDirs() {
		found, err := walkProjectDir(d, "", nil)
		if err == nil {
			return found, nil
		}
	}
	return nil, fmt.Errorf("no project dir found (checked %v)", app.ProjectDirs())
}

func walkProjectDir(dirname string, relPrefix string, result []project) ([]project, error) {
	entries, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			p, err := walkProjectDir(filepath.Join(dirname, e.Name()), relPrefix+e.Name()+"/", result)
			if err != nil {
				return nil, err
			}
			result = p
			continue
		}

		if e.Type().IsRegular() {
			result = append(result, project{
				RelPath: relPrefix + e.Name(),
				Path:    filepath.Join(dirname, e.Name()),
			})
		}
	}
	return result, nil
}

func readPythonProject(proj project) ([]byte, error) {
	f, err := os.Open(proj.Path)
	if err != nil {
		return nil, err
	}

	l, err := lmsp.ReadFile(f)
	if err != nil {
		return nil, err
	}

	man, err := l.Manifest()
	if err != nil {
		return nil, err
	}

	if man.Type != "python" {
		fmt.Printf("%s: skip %s program\n", proj.RelPath, man.Type)
		return nil, nil
	}

	program, err := l.Python()
	if err != nil {
		return nil, err
	}
	return []byte(program), nil
}
