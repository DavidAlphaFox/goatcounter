// Copyright © 2019 Martin Tournoij <martin@arp242.net>
// This file is part of GoatCounter and published under the terms of the EUPL
// v1.2, which can be found in the LICENSE file or at http://eupl12.zgo.at

package cron_test

import (
	"fmt"
	"testing"
	"time"

	"zgo.at/goatcounter"
	"zgo.at/goatcounter/gctest"
)

func TestBrowserStats(t *testing.T) {
	ctx, clean := gctest.DB(t)
	defer clean()

	site := goatcounter.MustGetSite(ctx)
	now := time.Date(2019, 8, 31, 14, 42, 0, 0, time.UTC)

	gctest.StoreHits(ctx, t, []goatcounter.Hit{
		{Site: site.ID, CreatedAt: now, Browser: "Firefox/68.0", FirstVisit: true},
		{Site: site.ID, CreatedAt: now, Browser: "Chrome/77.0.123.666"},
		{Site: site.ID, CreatedAt: now, Browser: "Firefox/69.0"},
		{Site: site.ID, CreatedAt: now, Browser: "Firefox/69.0"},
	}...)

	return

	var stats goatcounter.Stats
	total, err := stats.ListBrowsers(ctx, now, now)
	if err != nil {
		t.Fatal(err)
	}

	want := `4 -> [{Firefox 3 1} {Chrome 1 0}]`
	out := fmt.Sprintf("%d -> %v", total, stats)
	if want != out {
		t.Errorf("\nwant: %s\nout:  %s", want, out)
	}

	// Update existing.
	gctest.StoreHits(ctx, t, []goatcounter.Hit{
		{Site: site.ID, CreatedAt: now, Browser: "Firefox/69.0", FirstVisit: true},
		{Site: site.ID, CreatedAt: now, Browser: "Firefox/69.0"},
		{Site: site.ID, CreatedAt: now, Browser: "Firefox/70.0"},
		{Site: site.ID, CreatedAt: now, Browser: "Firefox/70.0"},
	}...)

	stats = goatcounter.Stats{}
	total, err = stats.ListBrowsers(ctx, now, now)
	if err != nil {
		t.Fatal(err)
	}

	want = `8 -> [{Firefox 7 2} {Chrome 1 0}]`
	out = fmt.Sprintf("%d -> %v", total, stats)
	if want != out {
		t.Errorf("\nwant: %s\nout:  %s", want, out)
	}

	// List just Firefox.
	stats = goatcounter.Stats{}
	total, err = stats.ListBrowser(ctx, "Firefox", now, now)
	if err != nil {
		t.Fatal(err)
	}

	want = `7 -> [{Firefox 69 4 1} {Firefox 70 2 0} {Firefox 68 1 1}]`
	out = fmt.Sprintf("%d -> %v", total, stats)
	if want != out {
		t.Errorf("\nwant: %s\nout:  %s", want, out)
	}
}
