package jujusvg

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gc "gopkg.in/check.v1"
	"gopkg.in/juju/charm.v6-unstable"

	"gopkg.in/juju/jujusvg.v1/assets"
)

func Test(t *testing.T) { gc.TestingT(t) }

type newSuite struct{}

var _ = gc.Suite(&newSuite{})

var bundle = `
services:
  mongodb:
    charm: "cs:precise/mongodb-21"
    num_units: 1
    annotations:
      "gui-x": "940.5"
      "gui-y": "388.7698359714502"
    constraints: "mem=2G cpu-cores=1"
  elasticsearch:
    charm: "cs:~charming-devs/precise/elasticsearch-2"
    num_units: 1
    annotations:
      "gui-x": "490.5"
      "gui-y": "369.7698359714502"
    constraints: "mem=2G cpu-cores=1"
  charmworld:
    charm: "cs:~juju-jitsu/precise/charmworld-58"
    num_units: 1
    expose: true
    annotations:
      "gui-x": "813.5"
      "gui-y": "112.23016402854975"
    options:
      charm_import_limit: -1
      source: "lp:~bac/charmworld/ingest-local-charms"
      revno: 511
relations:
  - - "charmworld:essearch"
    - "elasticsearch:essearch"
  - - "charmworld:database"
    - "mongodb:database"
series: precise
`

func iconURL(ref *charm.URL) string {
	return "http://0.1.2.3/" + ref.Path() + ".svg"
}

type emptyFetcher struct{}

func (f *emptyFetcher) FetchIcons(*charm.BundleData) (map[string][]byte, error) {
	return nil, nil
}

type errFetcher string

func (f *errFetcher) FetchIcons(*charm.BundleData) (map[string][]byte, error) {
	return nil, fmt.Errorf("%s", *f)
}

func (s *newSuite) TestNewFromBundle(c *gc.C) {
	b, err := charm.ReadBundleData(strings.NewReader(bundle))
	c.Assert(err, gc.IsNil)
	err = b.Verify(nil)
	c.Assert(err, gc.IsNil)

	cvs, err := NewFromBundle(b, iconURL, nil)
	c.Assert(err, gc.IsNil)

	var buf bytes.Buffer
	cvs.Marshal(&buf)
	c.Logf("%s", buf.String())
	assertXMLEqual(c, buf.Bytes(), []byte(`
<?xml version="1.0"?>
<!-- Generated by SVGo -->
<svg width="639" height="465"
     style="font-family:Ubuntu, sans-serif;" viewBox="0 0 639 465"
     xmlns="http://www.w3.org/2000/svg"
     xmlns:xlink="http://www.w3.org/1999/xlink">
<defs>
<g id="serviceBlock" transform="scale(0.8)" >`+assets.ServiceModule+`
</g>
<g id="healthCircle">
<circle cx="10" cy="10" r="10" style="stroke:#38B44A;fill:none;stroke-width:2px"/>
<circle cx="10" cy="10" r="5" style="fill:#38B44A"/>
</g>
<svg:svg xmlns:svg="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" id="icon-1">
&#x9;&#x9;&#x9;&#x9;&#x9;<svg:image width="96" height="96" xlink:href="http://0.1.2.3/~juju-jitsu/precise/charmworld-58.svg"></svg:image>
&#x9;&#x9;&#x9;&#x9;</svg:svg>
<svg:svg xmlns:svg="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" id="icon-2">
&#x9;&#x9;&#x9;&#x9;&#x9;<svg:image width="96" height="96" xlink:href="http://0.1.2.3/~charming-devs/precise/elasticsearch-2.svg"></svg:image>
&#x9;&#x9;&#x9;&#x9;</svg:svg>
<svg:svg xmlns:svg="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" id="icon-3">
&#x9;&#x9;&#x9;&#x9;&#x9;<svg:image width="96" height="96" xlink:href="http://0.1.2.3/precise/mongodb-21.svg"></svg:image>
&#x9;&#x9;&#x9;&#x9;</svg:svg>
</defs>
<g id="relations">
<line x1="417" y1="189" x2="189" y2="351" stroke="#38B44A" stroke-width="2px" stroke-dasharray="129.85, 20" />
<use x="293" y="260" xlink:href="#healthCircle" />
<line x1="417" y1="189" x2="544" y2="276" stroke="#38B44A" stroke-width="2px" stroke-dasharray="66.97, 20" />
<use x="470" y="222" xlink:href="#healthCircle" />
</g>
<g id="services">
<use x="323" y="0" xlink:href="#serviceBlock" id="charmworld" />
<use x="369" y="46" xlink:href="#icon-1" width="96" height="96" />
<g style="font-size:18px;fill:#505050;text-anchor:middle">
<text x="417" y="31" >charmworld</text>
</g>
<use x="0" y="257" xlink:href="#serviceBlock" id="elasticsearch" />
<use x="46" y="303" xlink:href="#icon-2" width="96" height="96" />
<g style="font-size:18px;fill:#505050;text-anchor:middle">
<text x="94" y="288" >elasticsearch</text>
</g>
<use x="450" y="276" xlink:href="#serviceBlock" id="mongodb" />
<use x="496" y="322" xlink:href="#icon-3" width="96" height="96" />
<g style="font-size:18px;fill:#505050;text-anchor:middle">
<text x="544" y="307" >mongodb</text>
</g>
</g>
</svg>
`))
}

func (s *newSuite) TestNewFromBundleWithUnplacedService(c *gc.C) {
	b, err := charm.ReadBundleData(strings.NewReader(bundle))
	c.Assert(err, gc.IsNil)
	err = b.Verify(nil)
	c.Assert(err, gc.IsNil)
	b.Services["charmworld"].Annotations["gui-x"] = ""
	b.Services["charmworld"].Annotations["gui-y"] = ""

	cvs, err := NewFromBundle(b, iconURL, nil)
	c.Assert(err, gc.IsNil)

	var buf bytes.Buffer
	cvs.Marshal(&buf)
	c.Logf("%s", buf.String())
	assertXMLEqual(c, buf.Bytes(), []byte(`
<?xml version="1.0"?>
<!-- Generated by SVGo -->
<svg width="922" height="302"
     style="font-family:Ubuntu, sans-serif;" viewBox="0 0 922 302"
     xmlns="http://www.w3.org/2000/svg"
     xmlns:xlink="http://www.w3.org/1999/xlink">
<defs>
<g id="serviceBlock" transform="scale(0.8)" >

<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px"
	 width="235.958px" height="235.958px" viewBox="-143.979 -6 235.958 235.958" enable-background="new -143.979 -6 235.958 235.958"
	 xml:space="preserve">
<g id="Layer_1" inkscape:version="0.48.3.1 r9886" sodipodi:docname="service_module.svg" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:cc="http://creativecommons.org/ns#" xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:svg="http://www.w3.org/2000/svg" xmlns:sodipodi="http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd" xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape">
	<g id="g33_1_" transform="translate(-399.571,-251.207)">
		<g id="path46_1_">
			<path fill="#BBBBBB" d="M410.565,479.165h-73.988c-38.324,0-57.56,0-68.272-10.713c-10.712-10.713-10.712-29.949-10.712-68.273
				v-73.986c-0.001-38.324-0.001-57.561,10.711-68.273c10.713-10.713,29.949-10.713,68.274-10.713h73.988
				c38.324,0,57.561,0,68.272,10.713c10.713,10.712,10.713,29.949,10.713,68.273v73.986c0,38.324,0,57.561-10.713,68.273
				C468.126,479.165,448.889,479.165,410.565,479.165z M336.577,257.207c-34.445,0-53.419,0-61.203,7.784
				s-7.783,26.757-7.782,61.202v73.986c0,34.444,0,53.419,7.784,61.202c7.784,7.784,26.757,7.784,61.201,7.784h73.988
				c34.444,0,53.418,0,61.202-7.784c7.783-7.783,7.783-26.758,7.783-61.202v-73.986c0-34.444,0-53.418-7.783-61.202
				c-7.784-7.784-26.758-7.784-61.202-7.784H336.577z"/>
		</g>
		<g id="path59_1_">
			<path fill="#BBBBBB" d="M410.565,479.165h-73.988c-38.324,0-57.56,0-68.272-10.713c-10.712-10.713-10.712-29.949-10.712-68.273
				v-73.986c0-38.324,0-57.561,10.712-68.273c10.713-10.713,29.949-10.713,68.272-10.713h73.988c38.324,0,57.561,0,68.272,10.713
				c10.713,10.712,10.713,29.949,10.713,68.273v73.986c0,38.324,0,57.561-10.713,68.273
				C468.126,479.165,448.889,479.165,410.565,479.165z M336.577,257.207c-34.444,0-53.417,0-61.201,7.784
				s-7.784,26.758-7.784,61.202v73.986c0,34.444,0,53.419,7.784,61.202c7.784,7.784,26.757,7.784,61.201,7.784h73.988
				c34.444,0,53.418,0,61.201-7.784c7.784-7.783,7.784-26.758,7.784-61.202v-73.986c0-34.444,0-53.418-7.784-61.202
				c-7.783-7.784-26.757-7.784-61.201-7.784H336.577z"/>
		</g>
	</g>
	<path fill-rule="evenodd" clip-rule="evenodd" fill="#BBBBBB" d="M-42,219.958h32c2.209,0,4,1.791,4,4v2c0,2.209-1.791,4-4,4h-32
		c-2.209,0-4-1.791-4-4v-2C-46,221.749-44.209,219.958-42,219.958z"/>
	<path fill-rule="evenodd" clip-rule="evenodd" fill="#BBBBBB" d="M-42-6h32c2.209,0,4,1.791,4,4v2c0,2.209-1.791,4-4,4h-32
		c-2.209,0-4-1.791-4-4v-2C-46-4.209-44.209-6-42-6z"/>
	<path fill-rule="evenodd" clip-rule="evenodd" fill="#BBBBBB" d="M81.979,127.979v-32c0-2.209,1.791-4,4-4h2c2.209,0,4,1.791,4,4
		v32c0,2.209-1.791,4-4,4h-2C83.771,131.979,81.979,130.188,81.979,127.979z"/>
	<path fill-rule="evenodd" clip-rule="evenodd" fill="#BBBBBB" d="M-143.979,127.979v-32c0-2.209,1.791-4,4-4h2c2.209,0,4,1.791,4,4
		v32c0,2.209-1.791,4-4,4h-2C-142.188,131.979-143.979,130.188-143.979,127.979z"/>
	<path fill="#FFFFFF" d="M10.994-1h-73.988c-73.987,0-73.987,0-73.985,73.986v73.986c0,73.986,0,73.986,73.985,73.986h73.988
		c73.985,0,73.985,0,73.985-73.986V72.986C84.979-1,84.979-1,10.994-1z"/>
</g>
<g id="Layer_2">
</g>
</svg>
</g>
<g id="healthCircle">
<circle cx="10" cy="10" r="10" style="stroke:#38B44A;fill:none;stroke-width:2px"/>
<circle cx="10" cy="10" r="5" style="fill:#38B44A"/>
</g>
<svg:svg xmlns:svg="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" id="icon-1">
&#x9;&#x9;&#x9;&#x9;&#x9;<svg:image width="96" height="96" xlink:href="http://0.1.2.3/~juju-jitsu/precise/charmworld-58.svg"></svg:image>
&#x9;&#x9;&#x9;&#x9;</svg:svg><svg:svg xmlns:svg="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" id="icon-2">
&#x9;&#x9;&#x9;&#x9;&#x9;<svg:image width="96" height="96" xlink:href="http://0.1.2.3/~charming-devs/precise/elasticsearch-2.svg"></svg:image>
&#x9;&#x9;&#x9;&#x9;</svg:svg><svg:svg xmlns:svg="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" id="icon-3">
&#x9;&#x9;&#x9;&#x9;&#x9;<svg:image width="96" height="96" xlink:href="http://0.1.2.3/precise/mongodb-21.svg"></svg:image>
&#x9;&#x9;&#x9;&#x9;</svg:svg></defs>
<g id="relations">
<line x1="733" y1="207" x2="189" y2="94" stroke="#38B44A" stroke-width="2px" stroke-dasharray="267.81, 20" />
<use x="451" y="140" xlink:href="#healthCircle" />
<line x1="733" y1="207" x2="639" y2="113" stroke="#38B44A" stroke-width="2px" stroke-dasharray="56.47, 20" />
<use x="676" y="150" xlink:href="#healthCircle" />
</g>
<g id="services">
<use x="733" y="113" xlink:href="#serviceBlock" id="charmworld" />
<use x="779" y="159" xlink:href="#icon-1" width="96" height="96" />
<g style="font-size:18px;fill:#505050;text-anchor:middle">
<text x="827" y="144" >charmworld</text>
</g>
<use x="0" y="0" xlink:href="#serviceBlock" id="elasticsearch" />
<use x="46" y="46" xlink:href="#icon-2" width="96" height="96" />
<g style="font-size:18px;fill:#505050;text-anchor:middle">
<text x="94" y="31" >elasticsearch</text>
</g>
<use x="450" y="19" xlink:href="#serviceBlock" id="mongodb" />
<use x="496" y="65" xlink:href="#icon-3" width="96" height="96" />
<g style="font-size:18px;fill:#505050;text-anchor:middle">
<text x="544" y="50" >mongodb</text>
</g>
</g>
</svg>
`))
}

func (s *newSuite) TestWithFetcher(c *gc.C) {
	b, err := charm.ReadBundleData(strings.NewReader(bundle))
	c.Assert(err, gc.IsNil)
	err = b.Verify(nil)
	c.Assert(err, gc.IsNil)

	cvs, err := NewFromBundle(b, iconURL, new(emptyFetcher))
	c.Assert(err, gc.IsNil)

	var buf bytes.Buffer
	cvs.Marshal(&buf)
	c.Logf("%s", buf.String())
	assertXMLEqual(c, buf.Bytes(), []byte(`
<?xml version="1.0"?>
<!-- Generated by SVGo -->
<svg width="639" height="465"
     style="font-family:Ubuntu, sans-serif;" viewBox="0 0 639 465"
     xmlns="http://www.w3.org/2000/svg"
     xmlns:xlink="http://www.w3.org/1999/xlink">
<defs>
<g id="serviceBlock" transform="scale(0.8)" >`+assets.ServiceModule+`
</g>
<g id="healthCircle">
<circle cx="10" cy="10" r="10" style="stroke:#38B44A;fill:none;stroke-width:2px"/>
<circle cx="10" cy="10" r="5" style="fill:#38B44A"/>
</g>
</defs>
<g id="relations">
<line x1="417" y1="189" x2="189" y2="351" stroke="#38B44A" stroke-width="2px" stroke-dasharray="129.85, 20" />
<use x="293" y="260" xlink:href="#healthCircle" />
<line x1="417" y1="189" x2="544" y2="276" stroke="#38B44A" stroke-width="2px" stroke-dasharray="66.97, 20" />
<use x="470" y="222" xlink:href="#healthCircle" />
</g>
<g id="services">
<use x="323" y="0" xlink:href="#serviceBlock" id="charmworld" />
<image x="369" y="46" width="96" height="96" xlink:href="http://0.1.2.3/~juju-jitsu/precise/charmworld-58.svg" />
<g style="font-size:18px;fill:#505050;text-anchor:middle">
<text x="417" y="31" >charmworld</text>
</g>
<use x="0" y="257" xlink:href="#serviceBlock" id="elasticsearch" />
<image x="46" y="303" width="96" height="96" xlink:href="http://0.1.2.3/~charming-devs/precise/elasticsearch-2.svg" />
<g style="font-size:18px;fill:#505050;text-anchor:middle">
<text x="94" y="288" >elasticsearch</text>
</g>
<use x="450" y="276" xlink:href="#serviceBlock" id="mongodb" />
<image x="496" y="322" width="96" height="96" xlink:href="http://0.1.2.3/precise/mongodb-21.svg" />
<g style="font-size:18px;fill:#505050;text-anchor:middle">
<text x="544" y="307" >mongodb</text>
</g>
</g>
</svg>
`))
}

func (s *newSuite) TestDefaultHTTPFetcher(c *gc.C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<svg></svg>")
	}))
	defer ts.Close()

	tsIconUrl := func(ref *charm.URL) string {
		return ts.URL + "/" + ref.Path() + ".svg"
	}

	b, err := charm.ReadBundleData(strings.NewReader(bundle))
	c.Assert(err, gc.IsNil)
	err = b.Verify(nil)
	c.Assert(err, gc.IsNil)

	cvs, err := NewFromBundle(b, tsIconUrl, &HTTPFetcher{IconURL: tsIconUrl})
	c.Assert(err, gc.IsNil)

	var buf bytes.Buffer
	cvs.Marshal(&buf)
	c.Logf("%s", buf.String())
	assertXMLEqual(c, buf.Bytes(), []byte(`
<?xml version="1.0"?>
<!-- Generated by SVGo -->
<svg width="639" height="465"
     style="font-family:Ubuntu, sans-serif;" viewBox="0 0 639 465"
     xmlns="http://www.w3.org/2000/svg"
     xmlns:xlink="http://www.w3.org/1999/xlink">
<defs>
<g id="serviceBlock" transform="scale(0.8)" >`+assets.ServiceModule+`
</g>
<g id="healthCircle">
<circle cx="10" cy="10" r="10" style="stroke:#38B44A;fill:none;stroke-width:2px"/>
<circle cx="10" cy="10" r="5" style="fill:#38B44A"/>
</g>
<svg:svg xmlns:svg="http://www.w3.org/2000/svg" id="icon-1"></svg:svg>
<svg:svg xmlns:svg="http://www.w3.org/2000/svg" id="icon-2"></svg:svg>
<svg:svg xmlns:svg="http://www.w3.org/2000/svg" id="icon-3"></svg:svg>
</defs>
<g id="relations">
<line x1="417" y1="189" x2="189" y2="351" stroke="#38B44A" stroke-width="2px" stroke-dasharray="129.85, 20" />
<use x="293" y="260" xlink:href="#healthCircle" />
<line x1="417" y1="189" x2="544" y2="276" stroke="#38B44A" stroke-width="2px" stroke-dasharray="66.97, 20" />
<use x="470" y="222" xlink:href="#healthCircle" />
</g>
<g id="services">
<use x="323" y="0" xlink:href="#serviceBlock" id="charmworld" />
<use x="369" y="46" xlink:href="#icon-1" width="96" height="96" />
<g style="font-size:18px;fill:#505050;text-anchor:middle">
<text x="417" y="31" >charmworld</text>
</g>
<use x="0" y="257" xlink:href="#serviceBlock" id="elasticsearch" />
<use x="46" y="303" xlink:href="#icon-2" width="96" height="96" />
<g style="font-size:18px;fill:#505050;text-anchor:middle">
<text x="94" y="288" >elasticsearch</text>
</g>
<use x="450" y="276" xlink:href="#serviceBlock" id="mongodb" />
<use x="496" y="322" xlink:href="#icon-3" width="96" height="96" />
<g style="font-size:18px;fill:#505050;text-anchor:middle">
<text x="544" y="307" >mongodb</text>
</g>
</g>
</svg>
`))

}

func (s *newSuite) TestFetcherError(c *gc.C) {
	b, err := charm.ReadBundleData(strings.NewReader(bundle))
	c.Assert(err, gc.IsNil)
	err = b.Verify(nil)
	c.Assert(err, gc.IsNil)

	ef := errFetcher("bad-wolf")
	_, err = NewFromBundle(b, iconURL, &ef)
	c.Assert(err, gc.ErrorMatches, "bad-wolf")
}

func (s *newSuite) TestWithBadBundle(c *gc.C) {
	b, err := charm.ReadBundleData(strings.NewReader(bundle))
	c.Assert(err, gc.IsNil)
	b.Relations[0][0] = "evil-unknown-service"
	cvs, err := NewFromBundle(b, iconURL, nil)
	c.Assert(err, gc.ErrorMatches, "cannot verify bundle: .*")
	c.Assert(cvs, gc.IsNil)
}

func (s *newSuite) TestWithBadPosition(c *gc.C) {
	b, err := charm.ReadBundleData(strings.NewReader(bundle))
	c.Assert(err, gc.IsNil)

	b.Services["charmworld"].Annotations["gui-x"] = "bad"
	cvs, err := NewFromBundle(b, iconURL, nil)
	c.Assert(err, gc.ErrorMatches, `service "charmworld" does not have a valid position`)
	c.Assert(cvs, gc.IsNil)

	b, err = charm.ReadBundleData(strings.NewReader(bundle))
	c.Assert(err, gc.IsNil)

	b.Services["charmworld"].Annotations["gui-y"] = "bad"
	cvs, err = NewFromBundle(b, iconURL, nil)
	c.Assert(err, gc.ErrorMatches, `service "charmworld" does not have a valid position`)
	c.Assert(cvs, gc.IsNil)
}
