package main

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/chuangbo/xip/pkg/qqwry"
	"github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"
)

func download(filename string) error {
	key, err := qqwry.GetDownloadKey()
	if err != nil {
		return fmt.Errorf("could not get key: %w", err)
	}

	total, dr, err := qqwry.Download(key)
	if err != nil {
		return fmt.Errorf("could not download db: %w", err)
	}
	defer dr.Close()

	fmt.Printf("Downloading to \"%s\"\n", filename)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	p, bar := progressBar(total, dr)
	defer bar.Close()

	// unzip
	// the reason to move gunzip out of qqwry package is for showing the correct progress
	z, err := zlib.NewReader(bar)
	if err != nil {
		return fmt.Errorf("could not create zlib reader: %w", err)
	}

	// save to file
	if _, err := io.Copy(f, z); err != nil {
		return err
	}
	// wait for our bar to complete and flush
	p.Wait()

	return f.Close()
}

func progressBar(total int64, r io.Reader) (*mpb.Progress, io.ReadCloser) {
	p := mpb.New(
		mpb.WithWidth(60),
		mpb.WithRefreshRate(time.Second),
	)

	bar := p.Add(total,
		mpb.NewBarFiller("[=>-|"),
		mpb.PrependDecorators(
			decor.CountersKibiByte("% .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.OnComplete(decor.AverageETA(decor.ET_STYLE_GO, decor.WC{W: 4}), "done"),
			decor.Name(" ] "),
			decor.AverageSpeed(decor.UnitKiB, "% .2f", decor.WC{W: 4}),
		),
	)

	// create proxy reader
	return p, bar.ProxyReader(r)
}
