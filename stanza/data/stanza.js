function Stanza(execute) {
  var proto = Object.create(HTMLElement.prototype);
  var development = descriptor.development;

  function template(name) {
    var t = descriptor.templates[name];
    if (!t) {
      throw new Error("template \"" + name + "\" is not found");
    }
    return t;
  }

  function createStanzaHelper(element) {
    return {
      query: function(params) {
        if (development) {
          console.log("query: called", params);
        }
        var t = template(params.template);
        var queryTemplate = Handlebars.compile(t, {noEscape: true});
        var query = queryTemplate(params.parameters);

        if (development) {
          console.log("query: query built:\n" + query);
          console.log("query: sending to", params.endpoint);
        }

        var p = $.ajax({
          url: params.endpoint,
          data: {
            format: "json",
            query: query
          }
        });

        if (development) {
          p.then(function(value, textStatus) {
            console.log("query:", textStatus, "data", value);
          });
        }

        return p;
      },
      render: function(params) {
        if (development) {
          console.log("render: called", params)
        }
        var t = template(params.template);
        var htmlTemplate = Handlebars.compile(t);
        var htmlFragment = htmlTemplate(params.parameters);
        if (development) {
          console.log("render: built:\n", htmlFragment)
        }
        var selector = params.selector || "main";
        $(selector, element.shadowRoot).html(htmlFragment);
        if (development) {
          console.log("render: wrote to \"" + selector + "\"")
        }
      },
      root: element.shadowRoot,
      select: function(selector) {
        return this.root.querySelector(selector);
      },
      selectAll: function(selector) {
        return this.root.querySelectorAll(selector);
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
    var found = false;
    descriptor.parameters.forEach(function(key) {
      if (attrName == key) {
        found = true;
      }
    });
    if (found) {
      update(this);
    }
  };

  document.registerElement(descriptor.elementName, {
    prototype: proto
  });
};
