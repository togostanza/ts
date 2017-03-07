# TIPS

## Using HTTP APIs from stanza

If you want to use some external APIs, use `$.ajax` to call them as follows:

```javascript
Stanza(function(stanza, params) {
  $.ajax({
    method: 'GET',
    url: 'http://example.com/example-api.json',
    data: {
      foo: 'hello, this is a query parameter'
    }
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
