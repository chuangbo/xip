package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func download(filename string) error {
	url := "https://unpkg.com/qqwry.ipdb/qqwry.ipdb"
	// url = "https://cdn.jsdelivr.net/npm/qqwry.ipdb/qqwry.ipdb"
	fmt.Printf("Downloading \"%s\" to \"%s\"\n", url, filename)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("could not download \"%s\": %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not download \"%s\": bad status code %s", url, resp.Status)
	}

	totalSize := resp.ContentLength
	if totalSize <= 0 {
		fmt.Println("Could not determine file size from Content-Length header.")
		// Handle cases where Content-Length is missing or invalid
	}

	dir, name := filepath.Split(filename)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("could not create temp dir: %w", err)
	}

	temp, err := os.CreateTemp(dir, name+".*.temp")
	if err != nil {
		return fmt.Errorf("could not create temp file: %w", err)
	}

	p, bar := progressBar(totalSize)

	proxyReader := bar.ProxyReader(resp.Body)
	defer proxyReader.Close()

	// save to file
	if _, err := io.Copy(temp, proxyReader); err != nil {
		return fmt.Errorf("could not download file: %w", err)
	}
	if totalSize <= 0 {
		bar.SetTotal(-1, true) // triggering complete event now
	}
	// wait for our bar to complete and flush
	p.Wait()

	if err := temp.Close(); err != nil {
		return fmt.Errorf("could not save file: %w", err)
	}

	if err := os.Rename(temp.Name(), filename); err != nil {
		return fmt.Errorf("could not move temp file: %w", err)
	}

	return nil
}

func progressBar(total int64) (*mpb.Progress, *mpb.Bar) {
	p := mpb.New(
		mpb.WithWidth(60),
		mpb.WithRefreshRate(time.Second),
	)

	bar := p.AddBar(total,
		mpb.PrependDecorators(
			decor.CountersKibiByte("% .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.OnComplete(decor.AverageETA(decor.ET_STYLE_GO, decor.WC{W: 4}), "done"),
			decor.Name(" ] "),
			decor.AverageSpeed(decor.SizeB1024(0), "% .2f", decor.WC{W: 4}),
		),
	)

	// create proxy reader
	return p, bar
}
