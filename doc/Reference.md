# Reference

## `ts` command

`ts` is a command to build TogoStanza stanzas of JavaScript version.
It generates a new stanza from built-in stanza blueprint, builds stanzas and serves stanzas for development.

### Create a new stanza

```sh
$ ts new <name>
```

Creates a new stanza. `<name>` is used for the directory of the stanza and the URL of the stanza. `<name>` can only contain alphanumeric characters and hyphens. The first character must be an alphabet.

### Build stanzas

```sh
$ ts build
```

Builds stanzas under current working directory. Outputs are written under `dist` directory.

### Serve stanzas for development

```sh
$ ts server [-port port]
```

Starts a web server for development. Automatically rebuilds stanzas into `dist` directory when the source is updated.

NOTE: Do not run `ts server` on a production server. `ts server` is designed only for development.

#### -port port

The port to listen on.

## Stanza structure

Each stanza has the following directory structure:

```
<stanza-name>
├── _header.html
├── assets
├── index.js
├── metadata.json
└── templates
    └── stanza.html
```

### _header.html

The contents of this file are embedded at the top of the stanza. Use this to load external JavaScript libraries.

Example: Load D3.js from `d3js.org`

```html
<!-- _header.html -->
<script src="http://d3js.org/d3.v3.min.js" charset="utf-8"></script>
```

NOTE: You need to be careful to use this feature. The scripts loaded here contaminate globally the page that uses the stanza. The other scripts loaded by other stanzas may conflict each other.

### assets (directory)

Holds assets (e.g., images). When you have `assets/example.png`, the asset is accessible in templates as follows.

```html
<!-- templates/stanza.html -->
<img src="assets/example.png" alt="example">
```

### index.js

Defines behavior of the stanza.

```js
// index.js
Stanza(function(stanza, params) {
  // stanza
});
```

The function passed to `Stanza()` is called when the stanza is embedded. Also it is called when `params` are updated. Issue queries, process and render the contents of the stanza in this function.

The first argument, `stanza`, provides useful methods to build stanzas. See details in [Stanza object](#stanza-object).

The second argument, `parms`, is an object which contains parameters given to the stanza. The parameters must be listed in `stanza:parameter` section in `metadata.json`. See [metatadata.json](#metadatajson).

### metadata.json

Describes the stanza, including the identifier of the stanza, human readable name of the stanza, what the stanza does, parameters, usage, license and author.

Example:

```json
{
  "@context": {
    "stanza": "http://togostanza.org/resource/stanza#"
  },
  "@id": "hello",
  "stanza:label": "Hello Example",
  "stanza:definition": "Greeting.",
  "stanza:parameter": [
  ],
  "stanza:usage": "<togostanza-hello></togostanza-hello>",
  "stanza:type": "Stanza",
  "stanza:context": "",
  "stanza:display": "",
  "stanza:provider": "provider of this stanza",
  "stanza:license": "",
  "stanza:author": "author name",
  "stanza:address": "name@example.org",
  "stanza:contributor": [
  ],
  "stanza:created": "2015-02-19",
  "stanza:updated": "2015-02-19"
}
```

### templates (directory)

Contains SPARQL query templates for `stanza.query()` and HTML templates for `stanza.render()`. The template is specified by the filename.

## Stanza object

### `stanza.query(options)`

Issues SPARQL query. `options` is an object that has the following properties:

<dl>
<dt>endpoint</dt><dd>SPARQL endpoint.</dd>
<dt>template</dt><dd>Template name for the query. Specify by the filename in `templates` directory.</dd>
<dt>parameters<dt><dd>Parameters to pass to the query template.</dd>
</dl>

The template is written in [Handlebars][].

`query()` returns a promise. Use `promise.then()` to wait until the query completed. You can handle errors with `promise.fail()`.

### `stanza.render(options)`

Renders contents from the given template. `options` is an object that has the following properties:

<dl>
<dt>template</dt><dd>Template name. Specify by the filename in `templates` directory.</dd>
<dt>parameters<dt><dd>Parameters to pass to the template.</dd>
<dt>selector<dt><dd>Destination selector. Specify this if you want to update the stanza partially. Default is `main`, which is the main element of the stanza.</dd>
</dl>

The template is written in [Handlebars][].

### `stanza.select(selectors)`

Returns the first element within stanza's shadow DOM that match with `selectors`. This is a shorthand method for `stanza.root.querySelector()`.

### `stanza.selectAll(selectors)`

Returns a list of the elements within stanza's shadow DOM that match with `selectors`. This is a shorthand method for `stanza.root.querySelectorAll()`.

### `stanza.root`

Shadow root of the stanza.

### `stanza.unwrapValueFromBinding(result)`

Unwraps an [SPARQL JSON Results Object][], returns simple Array of key-value objects.

```js
var result = {
  "head": {"vars": ["s", "p", "o"]
  },
  "results": {
    "bindings": [{
      "s": {
        "type": "uri",
        "value": "http://example.com/s"
      },
      "p": {
        "type": "uri",
        "value": "http://example.com/p"
      },
      "o": {
        "type": "uri",
        "value": "http://example.com/o"
      }
    }]
  }
}
unwrapValueFromBinding(result)
//=> [
//     {"s": "http://example.com/s", "p":"http://example.com/p", "o":"http://example.com/o"}
//   ]
```

### `stanza.grouping(ary, key1[, key2, ...])`

Groups an array of objects by specified keys.

Example: Group objects values of `x` then `y`.

```js
var ary = [
  {x: 1, y: 3},
  {x: 1, y: 4},
  {x: 2, y: 5},
  {x: 2, y: 6}
]

grouping(ary, "x", "y")
//=> [
//     {x: 1, y: [3, 4]},
//     {x: 2, y: [5, 6]}
//   ]
```

Example: Use a composite key.

```js
var ary = [
  {x: 1, y: 1, z: 3},
  {x: 1, y: 2, z: 4},
  {x: 2, y: 1, z: 5},
  {x: 2, y: 2, z: 6},
  {x: 1, y: 2, z: 7},
  {x: 2, y: 1, z: 8}
];

stanza.grouping(ary, ['x', 'y'], 'z');
//=> [
//     {x_y: [1, 1], z: [3]},
//     {x_y: [1, 2], z: [4, 7]},
//     {x_y: [2, 1], z: [5, 8]},
//     {x_y: [2, 2], z: [6]}]
//   ]
```

Example: Give an alias.

```js
var ary = [
  {x: 1, y: 3},
  {x: 1, y: 4},
  {x: 2, y: 5},
  {x: 2, y: 6}
];

stanza.grouping(ary, {key: 'x', alias: 'z'}, 'y');
//=> [
//     {z: 1, y: [3, 4]},
//     {z: 2, y: [5, 6]}
//   ]
```

## Embedding stanza

### Build and serve

Run `ts build` builds stanza provider into `dist` directory.

Run production web server (e.g. Apache, Nginx, ...) and serve `dist` directory as its document root.
Assume that we have deployed `dist` to `http://example.org/`.
Now you should have the list of available stanzas at `http://example.org/stanza`.

The help page of stanzas should be located at `http://example.org/stanza/<stanza-name>/help.html`.

NOTE: If you want to use stanzas in other domains than the domain stanza hosted, that is, embedding stanzas provided at `example.org` into `example.com` (not `example.org`), you need to configure your web server (`example.org`, which hosts stanzas) to explicitly allow cross-origin resource sharing (CORS). In order to make your stanzas embeddable into any domains, include `Access-Control-Allow-Origin: *` in HTTP headers of responses from the server.

### Load webcomponents.js

Stanzas generated by `ts` are built on top of [Web Components](http://webcomponents.org/); that is, Custom Elements, HTML Imports and Shadow DOM.
Unfortunately, not all browsers support them at this moment.
Thus we need polyfills, a kind of compatibility layer, provided by https://github.com/webcomponents/webcomponentsjs.

Include the following line in `<head></head>` of your html file:

```html
<script src="https://example.org/stanza/assets/components/webcomponentsjs/webcomponents.min.js"></script>
```

### Import stanza

Before using stanzas, you need to import the stanza. Include the following line in `<head>`.
This must be placed after the `<script>` element for `webcomponents.js` described above.

```html
<link rel="import" href="http://example.com/stanza/[stanza-name]/"></script>
```

Note that you need a trailing `/` for URL specified as `href`.

### Use stanza

Put the following code in `<body></body>`, where you want to place the stanza.

```html
<togostanza-[stanza-name]></togostanza-[stanza-name]>
```

If your stanza takes some parameters, pass them as attributes:

```html
<togostanza-[stanza-name] [parameter 1]=[value 1] [parameter 2]=[value 2]></togostanza-[stanza-name]>
```

  [handlebars]: http://handlebarsjs.com/
  [SPARQL JSON Results Object]: https://www.w3.org/TR/sparql11-results-json/#json-result-object
