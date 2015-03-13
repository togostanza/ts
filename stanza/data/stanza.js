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
      },
      grouping: function(rows /* , ...keys */) {
        var _this = this;

        var normalizedKeys = Array.prototype.slice.call(arguments, 1).reduce(function(acc, key) {
          if (key instanceof Object) {
            return acc.concat(
              Object.keys(key).map(function(k) {
                return [k, key[k]];
              })
            );
          } else {
            return acc.concat([[key, key]]);
          }
        }, []);

        return (function(rows, keys) {
          var callee = arguments.callee;
          var key1   = keys[0];
          var k1     = key1[0];
          var a1     = key1[1];

          if (keys.length === 1) return rows.map(function(row) { return row[k1] });

          var key2 = keys[1];
          var a2   = key2[1];

          return _this.groupBy(rows, function(row) {
            if (k1 instanceof Array) {
              return k1.reduce(function(acc, k) {
                acc[k] = row[k];
              }, {});
            } else {
              return row[k1];
            }
          }).map(function(i) {
            var ret = {};

            ret[a1] = i[0];
            ret[a2] = callee(i[1], keys.slice(1))

            return ret;
          });
        })(rows, normalizedKeys);
      },
      groupBy: function(array, func) {
        var ret = {};

        array.forEach(function(item) {
          var key  = func(item);
          ret[key] = ret[key] || [];

          ret[key].push(item);
        });

        return Object.keys(ret).map(function(key) {
          return [key, ret[key]];
        });
      },
      unwrapValueFromBinding: function(queryResult) {
        var bindings = queryResult.results.bindings;

        return bindings.map(function(binding) {
          var ret = {};

          Object.keys(binding).forEach(function(key) {
            ret[key] = binding[key].value;
          });

          return ret;
        });
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
