// Copyright 2018 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package name

import (
	"fmt"
	"strings"
)

const (
	repositoryChars  = "abcdefghijklmnopqrstuvwxyz0123456789_-./"
	regRepoDelimiter = "/"
)

// Repository stores a docker repository name in a structured form.
type Repository struct {
	Registry
	repository string
}

// RepositoryStr returns the repository component of the Repository.
func (r Repository) RepositoryStr() string {
	return r.repository
}

// Name returns the name from which the Repository was derived.
func (r Repository) Name() string {
	regName := r.Registry.Name()
	if regName != "" {
		return regName + regRepoDelimiter + r.RepositoryStr()
	}
	return r.RepositoryStr()
}

func (r Repository) String() string {
	return r.Name()
}

// Scope returns the scope required to perform the given action on the registry.
// TODO(jonjohnsonjr): consider moving scopes to a separate package.
func (r Repository) Scope(action string) string {
	return fmt.Sprintf("repository:%s:%s", r.RepositoryStr(), action)
}

func checkRepository(repository string) error {
	return checkElement("repository", repository, repositoryChars, 2, 255)
}

// NewRepository returns a new Repository representing the given name, according to the given strictness.
func NewRepository(name string, strict Strictness) (Repository, error) {
	if len(name) == 0 {
		return Repository{}, NewErrBadName("a repository name must be specified")
	}

	var registry string
	repo := name
	parts := strings.SplitN(name, regRepoDelimiter, 2)
	if len(parts) == 2 && (strings.ContainsRune(parts[0], '.') || strings.ContainsRune(parts[0], ':')) {
		// The first part of the repository is treated as the registry domain
		// iff it contains a '.' or ':' character, otherwise it is all repository
		// and the domain defaults to DockerHub.
		registry = parts[0]
		repo = parts[1]
	}

	if err := checkRepository(repo); err != nil {
		return Repository{}, err
	}

	reg, err := NewRegistry(registry, strict)
	if err != nil {
		return Repository{}, err
	}
	return Repository{reg, repo}, nil
}