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

func Run(app appcmd.App) error {
	t, err := GitTarget("refs/lego/robot-inventor").Open()
	if err != nil {
		return err
	}

	projects, err := listProjects(app)

	for _, project := range projects {
		data, err := readPythonProject(project)
		if err != nil {
			return fmt.Errorf("%s: %w", project.Name, err)
		}

		if data != nil {
			if err := t.Add(project.Name+".py", data); err != nil {
				return fmt.Errorf("%s: %w", project.Name, err)
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

type project struct {
	Name string
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

func walkProjectDir(dirname string, prefix string, result []project) ([]project, error) {
	entries, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			p, err := walkProjectDir(filepath.Join(dirname, e.Name()), prefix+e.Name()+"/", result)
			if err != nil {
				return nil, err
			}
			result = p
			continue
		}

		if e.Type().IsRegular() {
			result = append(result, project{
				Name: prefix + e.Name(),
				Path: filepath.Join(dirname, e.Name()),
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
		fmt.Printf("%s: skip %s program\n", proj.Name, man.Type)
		return nil, nil
	}

	program, err := l.Python()
	if err != nil {
		return nil, err
	}
	return []byte(program), nil
}
