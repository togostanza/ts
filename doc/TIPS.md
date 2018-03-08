# TIPS

## Using HTTP APIs from stanza

If you want to use some external APIs, use `fetch` API to call them as follows:

```javascript
Stanza(function(stanza, params) {
  var params = new URLSearchParams();
  params.set('foo', 'hello, this is a query parameter');
  fetch('http://example.com/example-api.json?' + params.toString(), {
    method: 'GET',
  }).then(function(response) {
    return response.json();
  }).then(function(data) {
    // Now you have `data`
    stanza.render({
      template: "stanza.html",
      parameters: {
        greeting: data, // pass `data` to the template
      }
    });
  });
});
```
