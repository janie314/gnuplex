package server

import (
	"fmt"
	"gnuplex-backend/server"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type dBFSDiff struct {
	inDB bool
	inFS bool
}
