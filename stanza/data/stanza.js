function Stanza(execute) {
  var proto = Object.create(HTMLElement.prototype);

  function createStanzaHelper(element) {
    return {
      query: function(params) {
        var queryTemplate = Handlebars.compile(descriptor.templates[params.template], {noEscape: true});
        var query = queryTemplate(params.parameters);

        return $.ajax({
          url: params.endpoint,
          data: {
            format: "json",
            query: query
          }
        });
      },
      render: function(params) {
        var htmlTemplate = Handlebars.compile(descriptor.templates[params.template]);
        var htmlPartial = htmlTemplate(params.parameters);
        var selector = params.selector || "main";
        $(selector, element.shadowRoot).html(htmlPartial);
      },
      $: function(selector) {
        return $(selector, element.shadowRoot);
      }
    };
  }

  function update(element) {
    var params = {};
    descriptor.parameters.forEach(function(key) {
      params[key] = element.getAttribute(key);
    });
    execute(createStanzaHelper(element), params);
  }

  proto.createdCallback = function() {
    var shadow = this.createShadowRoot();

    var style = document.createElement("style");
    style.appendChild(document.createTextNode(descriptor.stylesheet));
    shadow.appendChild(style);
    var main = document.createElement("main");
    shadow.appendChild(main);

    update(this);
  };

  proto.attributeChangedCallback = function(attrName, oldVal, newVal) {
    update(this);
  };

  document.registerElement(descriptor.elementName, {
    prototype: proto
  });
};
