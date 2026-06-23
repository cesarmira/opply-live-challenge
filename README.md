# opply-live-challenge

![Pistachios](https://commons.wikimedia.org/wiki/Special:FilePath/Pistachios.jpg)

## Pistachio

The pistachio (*Pistacia vera*) is a small tree from the cashew family, native
to Central Asia and the Middle East, whose edible seeds are eaten worldwide. The
green, faintly sweet kernel grows inside a hard, beige shell that splits open on
its own as the nut ripens. Pistachios are prized both as a snack and as an
ingredient — in everything from ice cream and baklava to savoury crusts — and
are a good source of protein, healthy fats, and fibre.

![Pistachio kernels](https://commons.wikimedia.org/wiki/Special:FilePath/Pistacia_vera_Kerman.jpg)

## API

An HTTP API that suggests alternative ingredients for a given ingredient. The
ingredient is passed as a query parameter; the response is JSON.

```
GET /suggest?ingredient=butter
```

```json
{
  "ingredient": "butter",
  "alternatives": [
    { "name": "olive oil", "notes": "use ~3/4 the amount for baking" },
    { "name": "coconut oil", "notes": "adds a mild coconut flavour" }
  ]
}
```

A missing `ingredient` parameter returns `400`; an unknown ingredient returns `404`.

### Run it

```
make run                                   # starts the server on :8080
curl 'http://127.0.0.1:8080/suggest?ingredient=pistachio'
```
