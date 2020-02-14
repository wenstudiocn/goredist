package dist

import "fmt"

type Version struct {
	major int
	minor int
	revision int
}

func NewVersion(major, minor, revision int) *Version {
	return &Version{
		major: major,
		minor: minor,
		revision: revision,
	}
}

func (self *Version)Major() int {
	return self.major
}

func (self *Version)Minor() int {
	return self.minor
}

func (self *Version) Revision() int {
	return self.revision
}

func (self *Version)String() string {
	return fmt.Sprintf("%v.%v.%v", self.major, self.minor, self.revision)
}