package repos

import (
	"github.com/c12s/star/internal/domain"
	"os"
	"path/filepath"
)

type nodeIdFSRepo struct {
	dirPath  string
	fileName string
	filePath string
}

func NewNodeIdFSRepo(dirPath, fileName string) (domain.NodeIdRepo, error) {
	return &nodeIdFSRepo{
		dirPath:  dirPath,
		fileName: fileName,
		filePath: dirPath + string(filepath.Separator) + fileName,
	}, nil
}

func (n nodeIdFSRepo) Get() (*domain.NodeId, error) {
	fileContents, err := os.ReadFile(n.filePath)
	if err != nil {
		return nil, err
	}
	return &domain.NodeId{
		Value: string(fileContents),
	}, nil
}

func (n nodeIdFSRepo) Put(nodeId domain.NodeId) error {
	return os.WriteFile(n.filePath, []byte(nodeId.Value), 0666)
}
