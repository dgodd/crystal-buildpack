package supply

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack"
	"github.com/kr/text"
)

type Stager interface {
	BuildDir() string
	CacheDir() string
	DepDir() string
	LinkDirectoryInDepDir(string, string) error
}

type Manifest interface {
	RootDir() string
	AllDependencyVersions(string) []string
	InstallDependency(libbuildpack.Dependency, string) error
	InstallOnlyVersion(string, string) error
}

type Command interface {
	Output(dir string, program string, args ...string) (string, error)
	Run(cmd *exec.Cmd) error
}

type Supplier struct {
	Manifest Manifest
	Stager   Stager
	Command  Command
	Log      *libbuildpack.Logger
	Shard    struct {
		Name           string `yaml:"name"`
		CrystalVersion string `yaml:"crystal"`
	}
}

func (s *Supplier) Run() error {
	s.Log.BeginStep("Supplying crystal")

	if err := s.Setup(); err != nil {
		return fmt.Errorf("Setup: %s", err)
	}

	if err := s.UntarLibevent(); err != nil {
		return fmt.Errorf("Installing libevent: %s", err)
	}

	if err := s.InstallCrystal(); err != nil {
		return fmt.Errorf("Installing Crystal: %s", err)
	}

	if s.Shard.Name != "" {
		if err := s.InstallShards(); err != nil {
			return fmt.Errorf("Installing Shards: %s", err)
		}
		if err := s.BuildApp(); err != nil {
			return fmt.Errorf("Building App: %s", err)
		}
	}

	return nil
}

func (s *Supplier) Setup() error {
	if err := libbuildpack.NewYAML().Load(filepath.Join(s.Stager.BuildDir(), "shard.yml"), &(s.Shard)); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("Reading shard.yml: %s", err)
		}
		s.Shard.Name = ""
	}

	if s.Shard.CrystalVersion == "" {
		s.Shard.CrystalVersion = "x"
	}
	versions := s.Manifest.AllDependencyVersions("crystal")
	if v, err := libbuildpack.FindMatchingVersion(s.Shard.CrystalVersion, versions); err != nil {
		return err
	} else {
		s.Shard.CrystalVersion = v
	}

	return nil
}

func (s *Supplier) UntarLibevent() error {
	if err := s.Manifest.InstallOnlyVersion("libevent", s.Stager.DepDir()); err != nil {
		return err
	}
	if err := s.Stager.LinkDirectoryInDepDir(filepath.Join(s.Stager.DepDir(), "libevent", "lib"), "lib"); err != nil {
		return err
	}
	if err := s.Stager.LinkDirectoryInDepDir(filepath.Join(s.Stager.DepDir(), "libevent", "lib", "pkgconfig"), "pkgconfig"); err != nil {
		return err
	}
	return s.Stager.LinkDirectoryInDepDir(filepath.Join(s.Stager.DepDir(), "libevent", "include"), "include")
}

func (s *Supplier) InstallCrystal() error {
	if err := s.Manifest.InstallDependency(libbuildpack.Dependency{Name: "crystal", Version: s.Shard.CrystalVersion}, s.Stager.DepDir()); err != nil {
		return err
	}
	crystalDir, err := filepath.Glob(filepath.Join(s.Stager.DepDir(), "crystal-*"))
	if err != nil {
		return err
	}
	for _, dir := range []string{"bin", "lib"} {
		if err := s.Stager.LinkDirectoryInDepDir(filepath.Join(crystalDir[0], dir), dir); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		}
	}
	return nil
}

func (s *Supplier) InstallShards() error {
	s.Log.BeginStep("Installing Dependencies")

	cmd := exec.Command("crystal", "deps", "--production")
	cmd.Dir = s.Stager.BuildDir()
	cmd.Stdout = text.NewIndentWriter(os.Stdout, []byte("       "))
	cmd.Stderr = text.NewIndentWriter(os.Stderr, []byte("       "))
	cmd.Env = append(os.Environ(), fmt.Sprintf("SHARDS_INSTALL_PATH=%s/shards_lib", s.Stager.CacheDir()))

	return s.Command.Run(cmd)
}

func (s *Supplier) BuildApp() error {
	s.Log.BeginStep("Compiling src/%s.cr (auto-detected from shard.yml)", s.Shard.Name)

	crystalDirs, err := filepath.Glob(filepath.Join(s.Stager.DepDir(), "crystal-*"))
	if err != nil {
		return err
	}
	crystalDir := crystalDirs[0]
	if found, err := libbuildpack.FileExists(filepath.Join(crystalDir, "share/crystal/src")); err != nil {
		return err
	} else if found {
		crystalDir = filepath.Join(crystalDir, "share/crystal/src")
	} else {
		crystalDir = filepath.Join(crystalDir, "src")
	}
	crystalPath := fmt.Sprintf("CRYSTAL_PATH=%s:%s/shards_lib:src", crystalDir, s.Stager.CacheDir())

	cmd := exec.Command("crystal", "build", fmt.Sprintf("src/%s.cr", s.Shard.Name), "--release", "-o", filepath.Join(s.Stager.DepDir(), "app"))
	cmd.Dir = s.Stager.BuildDir()
	cmd.Stdout = text.NewIndentWriter(os.Stdout, []byte("       "))
	cmd.Stderr = text.NewIndentWriter(os.Stderr, []byte("       "))
	cmd.Env = append(os.Environ(), crystalPath)

	return s.Command.Run(cmd)
}
