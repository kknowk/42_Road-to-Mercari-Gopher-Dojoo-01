package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func ftCat(reader io.Reader, writer io.Writer) error {
	buf := bufio.NewReader(reader)
	for {
		// read until newline
		line, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		_, writeErr := writer.Write([]byte(line))
		if writeErr != nil {
			return writeErr
		}

		// if push ^D
		if err == io.EOF {
			break
		}
	}
	return nil
}

func main() {
	return_Value := 0
	switch {
	case len(os.Args) == 1:
		// read from stdin
		if err := ftCat(os.Stdin, os.Stdout); err != nil {
			io.WriteString(os.Stderr, err.Error()+"\n")
			os.Exit(1)
		}

	default:
		for _, filePath := range os.Args[1:] {
			switch filePath {
			case "-":
				// read from stdin
				if err := ftCat(os.Stdin, os.Stdout); err != nil {
					io.WriteString(os.Stderr, err.Error()+"\n")
					os.Exit(1)
				}

			default:
				// read from file
				file, err := os.Open(filePath)
				if err != nil {
					// delete "open " from error message
					errMsg := strings.Replace(err.Error(), "open "+filePath, filePath, 1)
					io.WriteString(os.Stderr, "ft_cat: "+errMsg+"\n")
					return_Value = 1
					continue
				}
				defer file.Close()

				if err := ftCat(file, os.Stdout); err != nil {
					// delete "read " from error message
					errMsg := strings.Replace(err.Error(), "read "+filePath, filePath, 1)
					io.WriteString(os.Stderr, "ft_cat: "+errMsg+"\n")
					return_Value = 1
				}

				
			}
		}
	}
	os.Exit(return_Value)
}
