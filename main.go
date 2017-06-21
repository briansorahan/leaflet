package main

import (
	"html/template"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/oklog/ulid"
)

func main() {
	config, err := ParseConfig()
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(config.OutputFile)
	if err != nil {
		log.Fatal(err)
	}
	funcmap := template.FuncMap{
		"varName": varName,
	}
	tmpl, err := template.New("leaflet").Funcs(funcmap).Parse(indexHTML)
	if err != nil {
		_ = f.Close()
		log.Fatal(err)
	}
	if err := tmpl.Execute(f, config); err != nil {
		_ = f.Close()
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

const indexHTML = `<html>
  <head>
    <script src='https://api.mapbox.com/mapbox-gl-js/{{.Version}}/mapbox-gl.js'></script>
    <script src='https://npmcdn.com/@turf/turf/turf.min.js'></script>
    <link href='https://api.mapbox.com/mapbox-gl-js/{{.Version}}/mapbox-gl.css' rel='stylesheet' />
  </head>
  <body>
    <div id='map' style='width: 100%; height: 100%;'></div>
      <script>
        mapboxgl.accessToken = '{{.Token}}';
        var map = new mapboxgl.Map({
          container: 'map',
          style: '{{.Style}}'
        });

        map.on('load', function() {
          {{range .Layers}}
            var data = {{.Contents}};
            map.addSource('{{.Name}}',{"type":"geojson","data":data});
            map.addLayer({"id":"{{.Name}}","type":"fill","source":"{{.Name}}","paint":{"fill-color":"#000fff","fill-opacity":0.5}});
            map.fitBounds(turf.bbox(data));
          {{end}}
        });
      </script>
  </body>
</html>
`

func varName() string {
	var (
		t       = time.Now()
		entropy = rand.New(rand.NewSource(t.UnixNano()))
	)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
